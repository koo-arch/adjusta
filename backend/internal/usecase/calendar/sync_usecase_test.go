package calendar

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type fakeSyncTransaction struct {
	store SyncStore
}

func (t *fakeSyncTransaction) Do(ctx context.Context, fn func(store SyncStore) error) error {
	return fn(t.store)
}

type fakeSyncStore struct {
	findCalendarByGoogleCalendarIDFn    func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error)
	findAnyCalendarByGoogleCalendarIDFn func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error)
	createCalendarFn                    func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error)
	updateCalendarFn                    func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error)
	ensureUserCalendarRelationFn        func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) error
	listUserCalendarRelationsFn         func(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error)
	softDeleteUserCalendarRelationFn    func(ctx context.Context, userID, calendarID uuid.UUID) error
}

func (s *fakeSyncStore) FindCalendarByGoogleCalendarID(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
	return s.findCalendarByGoogleCalendarIDFn(ctx, userID, googleCalendarID)
}

func (s *fakeSyncStore) FindAnyCalendarByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
	return s.findAnyCalendarByGoogleCalendarIDFn(ctx, googleCalendarID)
}

func (s *fakeSyncStore) CreateCalendar(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
	return s.createCalendarFn(ctx, googleCalendarID, summary)
}

func (s *fakeSyncStore) UpdateCalendar(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
	return s.updateCalendarFn(ctx, id, googleCalendarID, summary)
}

func (s *fakeSyncStore) EnsureUserCalendarRelation(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) (*repoUserCalendar.UserCalendar, error) {
	if err := s.ensureUserCalendarRelationFn(ctx, userID, calendarID, role); err != nil {
		return nil, err
	}
	return &repoUserCalendar.UserCalendar{
		UserID:     userID,
		CalendarID: calendarID,
		Role:       role,
	}, nil
}

func (s *fakeSyncStore) ListUserCalendarRelations(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
	return s.listUserCalendarRelationsFn(ctx, userID)
}

func (s *fakeSyncStore) SoftDeleteUserCalendarRelation(ctx context.Context, userID, calendarID uuid.UUID) error {
	return s.softDeleteUserCalendarRelationFn(ctx, userID, calendarID)
}

func TestSyncCalendarAssignsExternalRoles(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	primaryCalendarID := uuid.New()
	referenceCalendarID := uuid.New()

	var ensuredRoles []domainvalue.UserCalendarRole

	uc := NewSyncUsecase(nil, nil, nil, &fakeSyncTransaction{
		store: &fakeSyncStore{
			findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
				return nil, repoerr.ErrNotFound
			},
			findAnyCalendarByGoogleCalendarIDFn: func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
				return nil, repoerr.ErrNotFound
			},
			createCalendarFn: func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				id := referenceCalendarID
				if googleCalendarID == "primary-cal" {
					id = primaryCalendarID
				}
				return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			updateCalendarFn: func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) error {
				ensuredRoles = append(ensuredRoles, role)
				return nil
			},
			listUserCalendarRelationsFn: func(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
				return nil, nil
			},
			softDeleteUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID) error {
				t.Fatalf("soft delete should not be called")
				return nil
			},
		},
	})

	err := uc.syncCalendar(ctx, []*customCalendar.CalendarList{
		{CalendarID: "primary-cal", Summary: "Primary", Primary: true},
		{CalendarID: "reference-cal", Summary: "Reference", Primary: false},
	}, &repoUser.User{ID: userID, Email: "user@example.com"})
	if err != nil {
		t.Fatalf("syncCalendar returned error: %v", err)
	}

	if len(ensuredRoles) != 2 {
		t.Fatalf("expected 2 ensured roles, got %d", len(ensuredRoles))
	}
	if ensuredRoles[0] != domainvalue.UserCalendarRolePrimary {
		t.Fatalf("expected primary role, got %s", ensuredRoles[0])
	}
	if ensuredRoles[1] != domainvalue.UserCalendarRoleReference {
		t.Fatalf("expected reference role, got %s", ensuredRoles[1])
	}
}

func TestSyncCalendarDoesNotDeleteAdjustaCandidateRelation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	currentCalendarID := uuid.New()
	missingPrimaryCalendarID := uuid.New()
	candidateCalendarID := uuid.New()

	var deletedCalendarIDs []uuid.UUID

	uc := NewSyncUsecase(nil, nil, nil, &fakeSyncTransaction{
		store: &fakeSyncStore{
			findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
				if googleCalendarID != "current-primary" {
					return nil, errors.New("unexpected google calendar id")
				}
				return &repoCalendar.Calendar{ID: currentCalendarID, GoogleCalendarID: googleCalendarID, Summary: "Current Primary"}, nil
			},
			findAnyCalendarByGoogleCalendarIDFn: func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
				t.Fatalf("find any calendar should not be called")
				return nil, nil
			},
			createCalendarFn: func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				t.Fatalf("create calendar should not be called")
				return nil, nil
			},
			updateCalendarFn: func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) error {
				return nil
			},
			listUserCalendarRelationsFn: func(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
				return []*UserCalendarRelationRecord{
					{
						CalendarID:       currentCalendarID,
						GoogleCalendarID: "current-primary",
						Role:             domainvalue.UserCalendarRolePrimary,
					},
					{
						CalendarID:       missingPrimaryCalendarID,
						GoogleCalendarID: "missing-primary",
						Role:             domainvalue.UserCalendarRolePrimary,
					},
					{
						CalendarID:        candidateCalendarID,
						GoogleCalendarID:  "adjusta-candidate",
						Role:              domainvalue.UserCalendarRoleAdjustaCandidate,
						SyncProposedDates: true,
					},
				}, nil
			},
			softDeleteUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID) error {
				deletedCalendarIDs = append(deletedCalendarIDs, calendarID)
				return nil
			},
		},
	})

	err := uc.syncCalendar(ctx, []*customCalendar.CalendarList{
		{CalendarID: "current-primary", Summary: "Current Primary", Primary: true},
	}, &repoUser.User{ID: userID, Email: "user@example.com"})
	if err != nil {
		t.Fatalf("syncCalendar returned error: %v", err)
	}

	if len(deletedCalendarIDs) != 1 {
		t.Fatalf("expected 1 deleted relation, got %d", len(deletedCalendarIDs))
	}
	if deletedCalendarIDs[0] != missingPrimaryCalendarID {
		t.Fatalf("expected missing primary relation to be deleted, got %s", deletedCalendarIDs[0])
	}
}
