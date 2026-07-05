package calendar

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	"github.com/koo-arch/adjusta-backend/internal/google"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type fakeCalendarUserRepository struct {
	readFn func(ctx context.Context, id uuid.UUID) (*repoUser.User, error)
}

func (r *fakeCalendarUserRepository) Read(ctx context.Context, id uuid.UUID) (*repoUser.User, error) {
	return r.readFn(ctx, id)
}

func (r *fakeCalendarUserRepository) FindByEmail(ctx context.Context, email string) (*repoUser.User, error) {
	return nil, errors.New("unexpected user find by email")
}

func (r *fakeCalendarUserRepository) Create(ctx context.Context, email string, opt repoUser.UserMutationOptions) (*repoUser.User, error) {
	return nil, errors.New("unexpected user create")
}

func (r *fakeCalendarUserRepository) Update(ctx context.Context, id uuid.UUID, opt repoUser.UserMutationOptions) (*repoUser.User, error) {
	return nil, errors.New("unexpected user update")
}

func (r *fakeCalendarUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected user delete")
}

func (r *fakeCalendarUserRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected user soft delete")
}

func (r *fakeCalendarUserRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected user restore")
}

type fakeCalendarTokenProvider struct{}

func (p *fakeCalendarTokenProvider) GetToken(ctx context.Context, userID uuid.UUID) (*google.AuthToken, error) {
	return &google.AuthToken{}, nil
}

type fakeCalendarServiceFactory struct {
	service CalendarService
}

func (f *fakeCalendarServiceFactory) New(ctx context.Context, token *google.AuthToken) (CalendarService, error) {
	return f.service, nil
}

type fakeListCalendarService struct {
	calendars []*ExternalCalendar
}

func (s *fakeListCalendarService) FetchCalendarList() ([]*ExternalCalendar, error) {
	return s.calendars, nil
}

func (s *fakeListCalendarService) CreateCalendar(summary string) (*ExternalCalendar, error) {
	return nil, errors.New("unexpected calendar create")
}

type fakeCalendarCache struct {
	calendars  []*ExternalCalendar
	found      bool
	setUserID  uuid.UUID
	setValue   []*ExternalCalendar
	setCalled  bool
	invalidate bool
}

func (c *fakeCalendarCache) Get(userID uuid.UUID) ([]*ExternalCalendar, bool) {
	return c.calendars, c.found
}

func (c *fakeCalendarCache) Set(userID uuid.UUID, calendars []*ExternalCalendar) {
	c.setCalled = true
	c.setUserID = userID
	c.setValue = calendars
}

func (c *fakeCalendarCache) Invalidate(userID uuid.UUID) {
	c.invalidate = true
}

func TestSyncGoogleCalendarsReturnsCachedCalendars(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	cached := []*ExternalCalendar{{CalendarID: "cached", Summary: "Cached"}}
	cache := &fakeCalendarCache{calendars: cached, found: true}

	uc := NewSyncUsecase(
		&fakeCalendarUserRepository{
			readFn: func(ctx context.Context, id uuid.UUID) (*repoUser.User, error) {
				return nil, errors.New("user repository should not be called on cache hit")
			},
		},
		&fakeCalendarTokenProvider{},
		&fakeCalendarServiceFactory{},
		nil,
		cache,
	)

	got, err := uc.SyncGoogleCalendars(ctx, userID, "user@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 || got[0].CalendarID != "cached" {
		t.Fatalf("expected cached calendars, got %#v", got)
	}
	if cache.setCalled {
		t.Fatalf("expected cache set not to be called on cache hit")
	}
}

func TestSyncGoogleCalendarsStoresSyncedCalendars(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	calendarID := uuid.New()
	cache := &fakeCalendarCache{}
	incoming := []*ExternalCalendar{{CalendarID: "primary-cal", Summary: "Primary", Primary: true}}

	uc := NewSyncUsecase(
		&fakeCalendarUserRepository{
			readFn: func(ctx context.Context, id uuid.UUID) (*repoUser.User, error) {
				return &repoUser.User{ID: id, Email: "user@example.com"}, nil
			},
		},
		&fakeCalendarTokenProvider{},
		&fakeCalendarServiceFactory{service: &fakeListCalendarService{calendars: incoming}},
		&fakeSyncTransaction{
			store: &fakeSyncStore{
				findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
					return nil, repoerr.ErrNotFound
				},
				findAnyCalendarByGoogleCalendarIDFn: func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
					return nil, repoerr.ErrNotFound
				},
				createCalendarFn: func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
					return &repoCalendar.Calendar{ID: calendarID, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
				},
				updateCalendarFn: func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
					return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
				},
				ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role value.UserCalendarRole, syncProposedDates *bool) error {
					return nil
				},
				listUserCalendarRelationsFn: func(ctx context.Context, userID uuid.UUID) ([]*CalendarRelation, error) {
					return nil, nil
				},
				softDeleteUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID) error {
					return nil
				},
			},
		},
		cache,
	)

	got, err := uc.SyncGoogleCalendars(ctx, userID, "user@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(got) != 1 || got[0].CalendarID != "primary-cal" {
		t.Fatalf("expected synced calendars, got %#v", got)
	}
	if !cache.setCalled {
		t.Fatalf("expected cache set to be called")
	}
	if cache.setUserID != userID {
		t.Fatalf("expected cache set user id %s, got %s", userID, cache.setUserID)
	}
	if len(cache.setValue) != 1 || cache.setValue[0].CalendarID != "primary-cal" {
		t.Fatalf("expected synced calendars to be cached, got %#v", cache.setValue)
	}
}
