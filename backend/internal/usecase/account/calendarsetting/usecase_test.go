package calendarsetting

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type fakeSettingsStore struct {
	userCalendars     []*repoUserCalendar.UserCalendar
	calendars         map[uuid.UUID]*repoCalendar.Calendar
	updatedIDs        []uuid.UUID
	updatedOpts       []repoUserCalendar.UserCalendarQueryOptions
	readCalendarCount int
}

type fakeSettingsUserCalendarRepository struct {
	store *fakeSettingsStore
}

func (r *fakeSettingsUserCalendarRepository) FindByIDAndUser(ctx context.Context, userID, id uuid.UUID) (*repoUserCalendar.UserCalendar, error) {
	for _, userCalendar := range r.store.userCalendars {
		if userCalendar.ID == id && userCalendar.UserID == userID {
			return userCalendar, nil
		}
	}
	return nil, repoerr.ErrNotFound
}

func (r *fakeSettingsUserCalendarRepository) FindByRole(ctx context.Context, userID uuid.UUID, role value.UserCalendarRole) (*repoUserCalendar.UserCalendar, error) {
	for _, userCalendar := range r.store.userCalendars {
		if userCalendar.UserID == userID && userCalendar.Role == role {
			return userCalendar, nil
		}
	}
	return nil, repoerr.ErrNotFound
}

func (r *fakeSettingsUserCalendarRepository) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoUserCalendar.UserCalendar, error) {
	filtered := make([]*repoUserCalendar.UserCalendar, 0, len(r.store.userCalendars))
	for _, userCalendar := range r.store.userCalendars {
		if userCalendar.UserID == userID {
			filtered = append(filtered, userCalendar)
		}
	}
	return filtered, nil
}

func (r *fakeSettingsUserCalendarRepository) Ensure(ctx context.Context, userID, calendarID uuid.UUID, opt repoUserCalendar.UserCalendarQueryOptions) (*repoUserCalendar.UserCalendar, error) {
	return nil, errors.New("unexpected user calendar ensure")
}

func (r *fakeSettingsUserCalendarRepository) Update(ctx context.Context, userID, id uuid.UUID, opt repoUserCalendar.UserCalendarQueryOptions) (*repoUserCalendar.UserCalendar, error) {
	for _, userCalendar := range r.store.userCalendars {
		if userCalendar.ID != id || userCalendar.UserID != userID {
			continue
		}
		if opt.Role != nil {
			userCalendar.Role = *opt.Role
		}
		if opt.IsVisible != nil {
			userCalendar.IsVisible = *opt.IsVisible
		}
		if opt.SyncProposedDates != nil {
			userCalendar.SyncProposedDates = *opt.SyncProposedDates
		}
		r.store.updatedIDs = append(r.store.updatedIDs, id)
		r.store.updatedOpts = append(r.store.updatedOpts, opt)
		return userCalendar, nil
	}
	return nil, errors.New("user calendar not found")
}

func (r *fakeSettingsUserCalendarRepository) SoftDeleteByUserAndCalendar(ctx context.Context, userID, calendarID uuid.UUID) error {
	return errors.New("unexpected user calendar soft delete")
}

type fakeSettingsCalendarRepository struct {
	store *fakeSettingsStore
}

func (r *fakeSettingsCalendarRepository) Read(ctx context.Context, id uuid.UUID) (*repoCalendar.Calendar, error) {
	r.store.readCalendarCount++
	calendar, ok := r.store.calendars[id]
	if !ok {
		return nil, errors.New("calendar not found")
	}
	return calendar, nil
}

func (r *fakeSettingsCalendarRepository) FilterByIDs(ctx context.Context, ids []uuid.UUID) ([]*repoCalendar.Calendar, error) {
	calendars := make([]*repoCalendar.Calendar, 0, len(ids))
	for _, id := range ids {
		calendar, ok := r.store.calendars[id]
		if !ok {
			continue
		}
		calendars = append(calendars, calendar)
	}
	return calendars, nil
}

func (r *fakeSettingsCalendarRepository) FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar filter by user id")
}

func (r *fakeSettingsCalendarRepository) FindByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) (*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar find by fields")
}

func (r *fakeSettingsCalendarRepository) FindByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar find by google calendar id")
}

func (r *fakeSettingsCalendarRepository) FilterByFields(ctx context.Context, userID uuid.UUID, opt repoCalendar.CalendarQueryOptions) ([]*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar filter by fields")
}

func (r *fakeSettingsCalendarRepository) Create(ctx context.Context, opt repoCalendar.CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar create")
}

func (r *fakeSettingsCalendarRepository) Update(ctx context.Context, id uuid.UUID, opt repoCalendar.CalendarMutationOptions) (*repoCalendar.Calendar, error) {
	return nil, errors.New("unexpected calendar update")
}

func (r *fakeSettingsCalendarRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected calendar delete")
}

func (r *fakeSettingsCalendarRepository) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected calendar soft delete")
}

func (r *fakeSettingsCalendarRepository) Restore(ctx context.Context, id uuid.UUID) error {
	return errors.New("unexpected calendar restore")
}

type fakeSettingsTransaction struct {
	store *fakeSettingsStore
}

func (t *fakeSettingsTransaction) DoCalendarSettings(ctx context.Context, fn func(repos CalendarSettingsRepositories) error) error {
	return fn(newFakeSettingsRepositories(t.store))
}

func newFakeSettingsRepositories(store *fakeSettingsStore) CalendarSettingsRepositories {
	return CalendarSettingsRepositories{
		Calendar:     &fakeSettingsCalendarRepository{store: store},
		UserCalendar: &fakeSettingsUserCalendarRepository{store: store},
	}
}

func newSettingsStore(userID uuid.UUID) (*fakeSettingsStore, *repoUserCalendar.UserCalendar, *repoUserCalendar.UserCalendar, *repoUserCalendar.UserCalendar) {
	primaryCalendarID := uuid.New()
	candidateCalendarID := uuid.New()
	referenceCalendarID := uuid.New()

	primary := &repoUserCalendar.UserCalendar{
		ID:         uuid.New(),
		UserID:     userID,
		CalendarID: primaryCalendarID,
		Role:       value.UserCalendarRolePrimary,
		IsVisible:  true,
	}
	candidate := &repoUserCalendar.UserCalendar{
		ID:                uuid.New(),
		UserID:            userID,
		CalendarID:        candidateCalendarID,
		Role:              value.UserCalendarRoleAdjustaCandidate,
		IsVisible:         true,
		SyncProposedDates: false,
	}
	reference := &repoUserCalendar.UserCalendar{
		ID:         uuid.New(),
		UserID:     userID,
		CalendarID: referenceCalendarID,
		Role:       value.UserCalendarRoleReference,
		IsVisible:  false,
	}

	store := &fakeSettingsStore{
		userCalendars: []*repoUserCalendar.UserCalendar{primary, candidate, reference},
		calendars: map[uuid.UUID]*repoCalendar.Calendar{
			primaryCalendarID:   {ID: primaryCalendarID, GoogleCalendarID: "primary@example.com", Summary: "メイン"},
			candidateCalendarID: {ID: candidateCalendarID, GoogleCalendarID: "candidate@example.com", Summary: "Adjusta 候補日程"},
			referenceCalendarID: {ID: referenceCalendarID, GoogleCalendarID: "reference@example.com", Summary: "祝日"},
		},
	}
	return store, primary, candidate, reference
}

func newSettingsUsecase(store *fakeSettingsStore, enabler CandidateCalendarEnabler) *Usecase {
	return NewUsecase(
		newFakeSettingsRepositories(store),
		&fakeSettingsTransaction{store: store},
		enabler,
	)
}

func TestListCalendarSettingsReturnsJoinedSettings(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, primary, _, _ := newSettingsStore(userID)
	uc := newSettingsUsecase(store, nil)

	settings, err := uc.ListCalendarSettings(context.Background(), userID, "user@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(settings) != 3 {
		t.Fatalf("expected 3 settings, got %d", len(settings))
	}
	if settings[0].ID != primary.ID {
		t.Fatalf("expected first setting id %s, got %s", primary.ID, settings[0].ID)
	}
	if settings[0].Summary != "メイン" || settings[0].GoogleCalendarID != "primary@example.com" {
		t.Fatalf("expected calendar fields to be joined, got %+v", settings[0])
	}
	if settings[0].Role != value.UserCalendarRolePrimary {
		t.Fatalf("expected primary role, got %s", settings[0].Role)
	}
	if store.readCalendarCount != 0 {
		t.Fatalf("expected list to avoid per-calendar reads, got %d reads", store.readCalendarCount)
	}
}

func TestGetCandidateSyncSettingReturnsDisabledBeforeCalendarCreation(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store := &fakeSettingsStore{calendars: map[uuid.UUID]*repoCalendar.Calendar{}}
	setting, err := newSettingsUsecase(store, nil).GetCandidateSyncSetting(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if setting.Enabled || setting.Calendar != nil {
		t.Fatalf("expected disabled setting without calendar, got %+v", setting)
	}
}

func TestSetCandidateSyncSettingCreatesCalendarBeforeEnabling(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store := &fakeSettingsStore{calendars: map[uuid.UUID]*repoCalendar.Calendar{}}
	calendarID := uuid.New()
	uc := newSettingsUsecase(store, CandidateCalendarEnablerFunc(func(ctx context.Context, gotUserID uuid.UUID, email string) error {
		if gotUserID != userID {
			t.Fatalf("expected user %s, got %s", userID, gotUserID)
		}
		store.calendars[calendarID] = &repoCalendar.Calendar{ID: calendarID, GoogleCalendarID: "candidate@example.com", Summary: "Adjusta 候補日程"}
		store.userCalendars = append(store.userCalendars, &repoUserCalendar.UserCalendar{
			ID: uuid.New(), UserID: userID, CalendarID: calendarID,
			Role: value.UserCalendarRoleAdjustaCandidate, IsVisible: true, SyncProposedDates: true,
		})
		return nil
	}))

	setting, err := uc.SetCandidateSyncSetting(context.Background(), userID, "user@example.com", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !setting.Enabled || setting.Calendar == nil {
		t.Fatalf("expected enabled setting with calendar, got %+v", setting)
	}
}

func TestUpdateCalendarSettingDemotesExistingPrimary(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, primary, _, reference := newSettingsStore(userID)
	uc := newSettingsUsecase(store, nil)

	newRole := value.UserCalendarRolePrimary
	updated, err := uc.UpdateCalendarSetting(context.Background(), userID, reference.ID, "user@example.com", CalendarSettingUpdateRequest{
		Role: &newRole,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.Role != value.UserCalendarRolePrimary {
		t.Fatalf("expected updated role primary, got %s", updated.Role)
	}
	if primary.Role != value.UserCalendarRoleReference {
		t.Fatalf("expected existing primary to be demoted to reference, got %s", primary.Role)
	}
	if len(store.updatedIDs) != 2 || store.updatedIDs[0] != primary.ID || store.updatedIDs[1] != reference.ID {
		t.Fatalf("expected demotion before promotion, got updated ids %v", store.updatedIDs)
	}
}

func TestUpdateCalendarSettingRejectsSyncOnNonAdjustaCandidate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, primary, _, _ := newSettingsStore(userID)
	uc := newSettingsUsecase(store, nil)

	enabled := true
	_, err := uc.UpdateCalendarSetting(context.Background(), userID, primary.ID, "user@example.com", CalendarSettingUpdateRequest{
		SyncProposedDates: &enabled,
	})
	if !internalErrors.IsKind(err, internalErrors.KindValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if len(store.updatedIDs) != 0 {
		t.Fatalf("expected no update, got updated ids %v", store.updatedIDs)
	}
}

func TestUpdateCalendarSettingResyncsWhenSyncEnabled(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, _, candidate, _ := newSettingsStore(userID)

	resyncCalls := 0
	uc := newSettingsUsecase(store, CandidateCalendarEnablerFunc(func(ctx context.Context, resyncUserID uuid.UUID, email string) error {
		resyncCalls++
		if resyncUserID != userID {
			t.Fatalf("expected resync for user %s, got %s", userID, resyncUserID)
		}
		return nil
	}))

	enabled := true
	updated, err := uc.UpdateCalendarSetting(context.Background(), userID, candidate.ID, "user@example.com", CalendarSettingUpdateRequest{
		SyncProposedDates: &enabled,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated.SyncProposedDates {
		t.Fatal("expected sync_proposed_dates to be enabled")
	}
	if resyncCalls != 1 {
		t.Fatalf("expected 1 resync call, got %d", resyncCalls)
	}
}

func TestUpdateCalendarSettingSkipsResyncWithoutTransition(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, _, candidate, _ := newSettingsStore(userID)
	candidate.SyncProposedDates = true

	resyncCalls := 0
	uc := newSettingsUsecase(store, CandidateCalendarEnablerFunc(func(ctx context.Context, resyncUserID uuid.UUID, email string) error {
		resyncCalls++
		return nil
	}))

	// true → true(変更なし)は再同期しない
	enabled := true
	if _, err := uc.UpdateCalendarSetting(context.Background(), userID, candidate.ID, "user@example.com", CalendarSettingUpdateRequest{
		SyncProposedDates: &enabled,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// true → false(無効化)も再同期しない
	disabled := false
	if _, err := uc.UpdateCalendarSetting(context.Background(), userID, candidate.ID, "user@example.com", CalendarSettingUpdateRequest{
		SyncProposedDates: &disabled,
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resyncCalls != 0 {
		t.Fatalf("expected no resync calls, got %d", resyncCalls)
	}
}

func TestUpdateCalendarSettingKeepsDisabledWhenCandidateEnableFails(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, _, candidate, _ := newSettingsStore(userID)

	uc := newSettingsUsecase(store, CandidateCalendarEnablerFunc(func(ctx context.Context, resyncUserID uuid.UUID, email string) error {
		return errors.New("google api unavailable")
	}))

	enabled := true
	_, err := uc.UpdateCalendarSetting(context.Background(), userID, candidate.ID, "user@example.com", CalendarSettingUpdateRequest{
		SyncProposedDates: &enabled,
	})
	if err == nil {
		t.Fatal("expected candidate enable failure")
	}
	if candidate.SyncProposedDates {
		t.Fatal("expected sync_proposed_dates to remain disabled")
	}
}

func TestUpdateCalendarSettingReturnsNotFoundForUnknownID(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, _, _, _ := newSettingsStore(userID)
	uc := newSettingsUsecase(store, nil)

	visible := true
	_, err := uc.UpdateCalendarSetting(context.Background(), userID, uuid.New(), "user@example.com", CalendarSettingUpdateRequest{
		IsVisible: &visible,
	})
	if !internalErrors.IsKind(err, internalErrors.KindNotFound) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestUpdateCalendarSettingReturnsNotFoundForOtherUsersCalendar(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	store, primary, _, _ := newSettingsStore(userID)
	uc := newSettingsUsecase(store, nil)

	visible := true
	_, err := uc.UpdateCalendarSetting(context.Background(), uuid.New(), primary.ID, "other@example.com", CalendarSettingUpdateRequest{
		IsVisible: &visible,
	})
	if !internalErrors.IsKind(err, internalErrors.KindNotFound) {
		t.Fatalf("expected not found error for other user's calendar, got %v", err)
	}
}
