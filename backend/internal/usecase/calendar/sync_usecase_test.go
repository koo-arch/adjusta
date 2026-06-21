package calendar

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type fakeSyncTransaction struct {
	store *fakeSyncStore
}

func (t *fakeSyncTransaction) Do(ctx context.Context, fn func(repos SyncTxRepositories) error) error {
	calendarRepo := &fakeSyncCalendarRepository{
		store:     t.store,
		calendars: map[uuid.UUID]*repoCalendar.Calendar{},
	}
	return fn(SyncTxRepositories{
		Calendar:     calendarRepo,
		UserCalendar: &fakeSyncUserCalendarRepository{store: t.store, calendarRepo: calendarRepo},
	})
}

type fakeCalendarService struct {
	createCalendarFn func(summary string) (*CalendarRecord, error)
}

func (s *fakeCalendarService) FetchCalendarList() ([]*CalendarRecord, error) {
	return nil, nil
}

func (s *fakeCalendarService) CreateCalendar(summary string) (*CalendarRecord, error) {
	if s.createCalendarFn == nil {
		return nil, errors.New("unexpected create calendar call")
	}
	return s.createCalendarFn(summary)
}

type fakeSyncStore struct {
	findCalendarByGoogleCalendarIDFn    func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error)
	findAnyCalendarByGoogleCalendarIDFn func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error)
	createCalendarFn                    func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error)
	updateCalendarFn                    func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error)
	ensureUserCalendarRelationFn        func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) error
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

func (s *fakeSyncStore) EnsureUserCalendarRelation(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) (*repoUserCalendar.UserCalendar, error) {
	if err := s.ensureUserCalendarRelationFn(ctx, userID, calendarID, role, syncProposedDates); err != nil {
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

type fakeSyncCalendarRepository struct {
	store     *fakeSyncStore
	calendars map[uuid.UUID]*repoCalendar.Calendar
}

func (r *fakeSyncCalendarRepository) Read(ctx context.Context, id uuid.UUID) (*repoCalendar.Calendar, error) {
	if calendar, ok := r.calendars[id]; ok {
		return calendar, nil
	}
	return nil, errors.New("unexpected calendar read")
}

func (r *fakeSyncCalendarRepository) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar filter by user id")
}

func (r *fakeSyncCalendarRepository) FindByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) (*repoCalendar.Calendar, error) {
	if opt.GoogleCalendarID == nil {
		return nil, errors.New("unexpected calendar find fields")
	}
	return r.store.findCalendarByGoogleCalendarIDFn(ctx, userID, *opt.GoogleCalendarID)
}

func (r *fakeSyncCalendarRepository) FindByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
	return r.store.findAnyCalendarByGoogleCalendarIDFn(ctx, googleCalendarID)
}

func (r *fakeSyncCalendarRepository) FilterByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) ([]*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar filter fields")
}

func (r *fakeSyncCalendarRepository) Create(ctx context.Context, opt repoCalendar.CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	if opt.GoogleCalendarID == nil || opt.Summary == nil {
		return nil, errors.New("google calendar id and summary are required")
	}
	calendar, err := r.store.createCalendarFn(ctx, *opt.GoogleCalendarID, *opt.Summary)
	if err == nil && calendar != nil {
		r.calendars[calendar.ID] = calendar
	}
	return calendar, err
}

func (r *fakeSyncCalendarRepository) Update(ctx context.Context, id uuid.UUID, opt repoCalendar.CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	if opt.GoogleCalendarID == nil || opt.Summary == nil {
		return nil, errors.New("google calendar id and summary are required")
	}
	calendar, err := r.store.updateCalendarFn(ctx, id, *opt.GoogleCalendarID, *opt.Summary)
	if err == nil && calendar != nil {
		r.calendars[calendar.ID] = calendar
	}
	return calendar, err
}

func (r *fakeSyncCalendarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected calendar delete")
}

func (r *fakeSyncCalendarRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected calendar soft delete")
}

func (r *fakeSyncCalendarRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected calendar restore")
}

type fakeSyncUserCalendarRepository struct {
	store        *fakeSyncStore
	calendarRepo *fakeSyncCalendarRepository
}

func (r *fakeSyncUserCalendarRepository) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoUserCalendar.UserCalendar, error) {
	relations, err := r.store.listUserCalendarRelationsFn(ctx, userID)
	if err != nil {
		return nil, err
	}

	userCalendars := make([]*repoUserCalendar.UserCalendar, 0, len(relations))
	for _, relation := range relations {
		r.calendarRepo.calendars[relation.CalendarID] = &repoCalendar.Calendar{
			ID:               relation.CalendarID,
			GoogleCalendarID: relation.GoogleCalendarID,
		}
		userCalendars = append(userCalendars, &repoUserCalendar.UserCalendar{
			UserID:            userID,
			CalendarID:        relation.CalendarID,
			Role:              relation.Role,
			SyncProposedDates: relation.SyncProposedDates,
		})
	}
	return userCalendars, nil
}

func (r *fakeSyncUserCalendarRepository) Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt repoUserCalendar.UserCalendarQueryOptions) (*repoUserCalendar.UserCalendar, error) {
	if opt.Role == nil {
		return nil, errors.New("role is required")
	}
	if err := r.store.ensureUserCalendarRelationFn(ctx, userID, calendarID, *opt.Role, opt.SyncProposedDates); err != nil {
		return nil, err
	}
	return &repoUserCalendar.UserCalendar{
		UserID:            userID,
		CalendarID:        calendarID,
		Role:              *opt.Role,
		SyncProposedDates: opt.SyncProposedDates != nil && *opt.SyncProposedDates,
	}, nil
}

func (r *fakeSyncUserCalendarRepository) SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error {
	return r.store.softDeleteUserCalendarRelationFn(ctx, userID, calendarID)
}

func TestSyncCalendarAssignsExternalRolesWithoutCreatingAdjustaCandidate(t *testing.T) {
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
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) error {
				ensuredRoles = append(ensuredRoles, role)
				if role == domainvalue.UserCalendarRoleAdjustaCandidate && syncProposedDates != nil {
					t.Fatalf("expected new adjusta candidate relation to use default sync setting")
				}
				if role != domainvalue.UserCalendarRoleAdjustaCandidate && syncProposedDates != nil {
					t.Fatalf("expected externally synced roles not to override sync setting")
				}
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

	_, err := uc.syncCalendar(ctx, &fakeCalendarService{
		createCalendarFn: func(summary string) (*CalendarRecord, error) {
			t.Fatalf("create calendar should not be called")
			return nil, nil
		},
	}, []*CalendarRecord{
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

func TestSyncCalendarRecreatesMissingAdjustaCandidateRelation(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	currentCalendarID := uuid.New()
	missingPrimaryCalendarID := uuid.New()
	candidateCalendarID := uuid.New()
	recreatedCandidateCalendarID := uuid.New()

	var deletedCalendarIDs []uuid.UUID
	currentRelations := []*UserCalendarRelationRecord{
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
	}

	uc := NewSyncUsecase(nil, nil, nil, &fakeSyncTransaction{
		store: &fakeSyncStore{
			findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
				switch googleCalendarID {
				case "current-primary":
					return &repoCalendar.Calendar{ID: currentCalendarID, GoogleCalendarID: googleCalendarID, Summary: "Current Primary"}, nil
				case "recreated-adjusta-candidate":
					return nil, repoerr.ErrNotFound
				default:
					return nil, errors.New("unexpected google calendar id")
				}
			},
			findAnyCalendarByGoogleCalendarIDFn: func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
				return nil, repoerr.ErrNotFound
			},
			createCalendarFn: func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				if googleCalendarID != "recreated-adjusta-candidate" {
					t.Fatalf("unexpected calendar creation for %s", googleCalendarID)
				}
				return &repoCalendar.Calendar{ID: recreatedCandidateCalendarID, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			updateCalendarFn: func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) error {
				if role == domainvalue.UserCalendarRoleAdjustaCandidate {
					if syncProposedDates == nil || !*syncProposedDates {
						t.Fatal("expected recreated adjusta candidate relation to preserve sync setting")
					}
					currentRelations = append(currentRelations, &UserCalendarRelationRecord{
						CalendarID:        calendarID,
						GoogleCalendarID:  "recreated-adjusta-candidate",
						Role:              role,
						SyncProposedDates: *syncProposedDates,
					})
				}
				return nil
			},
			listUserCalendarRelationsFn: func(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
				return currentRelations, nil
			},
			softDeleteUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID) error {
				deletedCalendarIDs = append(deletedCalendarIDs, calendarID)
				next := make([]*UserCalendarRelationRecord, 0, len(currentRelations))
				for _, relation := range currentRelations {
					if relation.CalendarID != calendarID {
						next = append(next, relation)
					}
				}
				currentRelations = next
				return nil
			},
		},
	})

	_, err := uc.syncCalendar(ctx, &fakeCalendarService{
		createCalendarFn: func(summary string) (*CalendarRecord, error) {
			if summary != domainUserCalendar.AdjustaCandidateCalendarSummary {
				t.Fatalf("unexpected candidate calendar summary: %s", summary)
			}
			return &CalendarRecord{
				CalendarID: "recreated-adjusta-candidate",
				Summary:    summary,
				Primary:    false,
			}, nil
		},
	}, []*CalendarRecord{
		{CalendarID: "current-primary", Summary: "Current Primary", Primary: true},
	}, &repoUser.User{ID: userID, Email: "user@example.com"})
	if err != nil {
		t.Fatalf("syncCalendar returned error: %v", err)
	}

	if len(deletedCalendarIDs) != 2 {
		t.Fatalf("expected 2 deleted relations, got %d", len(deletedCalendarIDs))
	}
	if deletedCalendarIDs[0] != candidateCalendarID {
		t.Fatalf("expected old candidate relation to be replaced first, got %s", deletedCalendarIDs[0])
	}
	if deletedCalendarIDs[1] != missingPrimaryCalendarID {
		t.Fatalf("expected missing primary relation to be deleted, got %s", deletedCalendarIDs[1])
	}

	foundRecreatedCandidate := false
	for _, relation := range currentRelations {
		if relation.Role == domainvalue.UserCalendarRoleAdjustaCandidate && relation.CalendarID == recreatedCandidateCalendarID {
			foundRecreatedCandidate = true
			break
		}
	}
	if !foundRecreatedCandidate {
		t.Fatal("expected recreated adjusta candidate relation to remain")
	}
}

func TestSyncCalendarDoesNotRecreateMissingAdjustaCandidateWhenSyncDisabled(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	primaryCalendarID := uuid.New()
	candidateCalendarID := uuid.New()

	var ensuredRoles []domainvalue.UserCalendarRole
	var deletedCalendarIDs []uuid.UUID

	uc := NewSyncUsecase(nil, nil, nil, &fakeSyncTransaction{
		store: &fakeSyncStore{
			findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
				if googleCalendarID != "primary-cal" {
					t.Fatalf("unexpected google calendar id: %s", googleCalendarID)
				}
				return &repoCalendar.Calendar{ID: primaryCalendarID, GoogleCalendarID: googleCalendarID, Summary: "Primary"}, nil
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
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) error {
				ensuredRoles = append(ensuredRoles, role)
				if role == domainvalue.UserCalendarRoleAdjustaCandidate {
					t.Fatalf("adjusta candidate relation should not be recreated while sync is disabled")
				}
				return nil
			},
			listUserCalendarRelationsFn: func(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
				return []*UserCalendarRelationRecord{
					{
						CalendarID:        candidateCalendarID,
						GoogleCalendarID:  "missing-adjusta-candidate",
						Role:              domainvalue.UserCalendarRoleAdjustaCandidate,
						SyncProposedDates: false,
					},
				}, nil
			},
			softDeleteUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID) error {
				deletedCalendarIDs = append(deletedCalendarIDs, calendarID)
				return nil
			},
		},
	})

	calendars, err := uc.syncCalendar(ctx, &fakeCalendarService{
		createCalendarFn: func(summary string) (*CalendarRecord, error) {
			t.Fatalf("create calendar should not be called")
			return nil, nil
		},
	}, []*CalendarRecord{
		{CalendarID: "primary-cal", Summary: "Primary", Primary: true},
	}, &repoUser.User{ID: userID, Email: "user@example.com"})
	if err != nil {
		t.Fatalf("syncCalendar returned error: %v", err)
	}

	if len(calendars) != 1 {
		t.Fatalf("expected returned calendars length to remain unchanged, got %d", len(calendars))
	}
	if len(ensuredRoles) != 1 {
		t.Fatalf("expected only primary role to be ensured, got %d", len(ensuredRoles))
	}
	if ensuredRoles[0] != domainvalue.UserCalendarRolePrimary {
		t.Fatalf("expected primary role, got %s", ensuredRoles[0])
	}
	if len(deletedCalendarIDs) != 0 {
		t.Fatalf("expected no deleted relations, got %d", len(deletedCalendarIDs))
	}
}

func TestSyncCalendarReusesIncomingAdjustaCandidateCalendar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	primaryCalendarID := uuid.New()
	candidateCalendarID := uuid.New()

	var ensuredRolePairs [][2]interface{}

	uc := NewSyncUsecase(nil, nil, nil, &fakeSyncTransaction{
		store: &fakeSyncStore{
			findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
				return nil, repoerr.ErrNotFound
			},
			findAnyCalendarByGoogleCalendarIDFn: func(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
				return nil, repoerr.ErrNotFound
			},
			createCalendarFn: func(ctx context.Context, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				id := candidateCalendarID
				if googleCalendarID == "primary-cal" {
					id = primaryCalendarID
				}
				return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			updateCalendarFn: func(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repoCalendar.Calendar, error) {
				return &repoCalendar.Calendar{ID: id, GoogleCalendarID: googleCalendarID, Summary: summary}, nil
			},
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) error {
				ensuredRolePairs = append(ensuredRolePairs, [2]interface{}{calendarID, role})
				if role == domainvalue.UserCalendarRoleAdjustaCandidate && syncProposedDates != nil {
					t.Fatalf("expected incoming adjusta candidate relation without stored setting to use default sync setting")
				}
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

	calendars, err := uc.syncCalendar(ctx, &fakeCalendarService{
		createCalendarFn: func(summary string) (*CalendarRecord, error) {
			t.Fatalf("create calendar should not be called")
			return nil, nil
		},
	}, []*CalendarRecord{
		{CalendarID: "primary-cal", Summary: "Primary", Primary: true},
		{CalendarID: "managed-candidate", Summary: domainUserCalendar.AdjustaCandidateCalendarSummary, Primary: false},
	}, &repoUser.User{ID: userID, Email: "user@example.com"})
	if err != nil {
		t.Fatalf("syncCalendar returned error: %v", err)
	}

	if len(calendars) != 2 {
		t.Fatalf("expected returned calendars length to stay unchanged, got %d", len(calendars))
	}
	if len(ensuredRolePairs) != 2 {
		t.Fatalf("expected 2 ensured relations, got %d", len(ensuredRolePairs))
	}
	if ensuredRolePairs[0][1] != domainvalue.UserCalendarRoleAdjustaCandidate {
		t.Fatalf("expected adjusta candidate relation first, got %v", ensuredRolePairs[0][1])
	}
	if ensuredRolePairs[1][1] != domainvalue.UserCalendarRolePrimary {
		t.Fatalf("expected primary relation second, got %v", ensuredRolePairs[1][1])
	}
	if ensuredRolePairs[0][0] == ensuredRolePairs[1][0] {
		t.Fatal("expected candidate calendar to be excluded from external role sync")
	}
}

func TestSyncCalendarPreservesAdjustaCandidateSyncSetting(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := uuid.New()
	candidateCalendarID := uuid.New()

	var ensuredSyncProposedDates *bool

	uc := NewSyncUsecase(nil, nil, nil, &fakeSyncTransaction{
		store: &fakeSyncStore{
			findCalendarByGoogleCalendarIDFn: func(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repoCalendar.Calendar, error) {
				t.Fatalf("find calendar should not be called")
				return nil, nil
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
			ensureUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole, syncProposedDates *bool) error {
				if role == domainvalue.UserCalendarRoleAdjustaCandidate {
					ensuredSyncProposedDates = syncProposedDates
				}
				return nil
			},
			listUserCalendarRelationsFn: func(ctx context.Context, userID uuid.UUID) ([]*UserCalendarRelationRecord, error) {
				return []*UserCalendarRelationRecord{
					{
						CalendarID:        candidateCalendarID,
						GoogleCalendarID:  "managed-candidate",
						Role:              domainvalue.UserCalendarRoleAdjustaCandidate,
						SyncProposedDates: false,
					},
				}, nil
			},
			softDeleteUserCalendarRelationFn: func(ctx context.Context, userID, calendarID uuid.UUID) error {
				t.Fatalf("soft delete should not be called")
				return nil
			},
		},
	})

	_, err := uc.syncCalendar(ctx, &fakeCalendarService{
		createCalendarFn: func(summary string) (*CalendarRecord, error) {
			t.Fatalf("create calendar should not be called")
			return nil, nil
		},
	}, []*CalendarRecord{
		{CalendarID: "managed-candidate", Summary: domainUserCalendar.AdjustaCandidateCalendarSummary, Primary: false},
	}, &repoUser.User{ID: userID, Email: "user@example.com"})
	if err != nil {
		t.Fatalf("syncCalendar returned error: %v", err)
	}

	if ensuredSyncProposedDates == nil {
		t.Fatal("expected adjusta candidate sync setting to be forwarded")
	}
	if *ensuredSyncProposedDates {
		t.Fatal("expected adjusta candidate sync setting to remain false")
	}
}
