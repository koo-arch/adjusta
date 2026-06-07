package calendar

import (
	"context"

	"github.com/google/uuid"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	repoUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type calendarSyncUserReader struct {
	userRepo repoUser.UserRepository
}

func NewCalendarSyncUserReader(userRepo repoUser.UserRepository) usecaseCalendar.UserReader {
	return &calendarSyncUserReader{userRepo: userRepo}
}

func (r *calendarSyncUserReader) GetByID(ctx context.Context, userID uuid.UUID) (*repositorymodel.User, error) {
	return r.userRepo.Read(ctx, userID, repoUser.UserQueryOptions{})
}

type calendarSyncTransaction struct {
	uow infraRepository.UnitOfWork
}

func NewCalendarSyncTransaction(uow infraRepository.UnitOfWork) usecaseCalendar.SyncTransaction {
	return &calendarSyncTransaction{uow: uow}
}

func (t *calendarSyncTransaction) Do(ctx context.Context, fn func(store usecaseCalendar.SyncStore) error) error {
	return t.uow.Do(ctx, func(repos infraRepository.Repositories) error {
		return fn(&calendarSyncStore{repos: repos})
	})
}

type calendarSyncStore struct {
	repos infraRepository.Repositories
}

func (s *calendarSyncStore) FindCalendarByGoogleCalendarID(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*repositorymodel.StoredCalendar, error) {
	return s.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		GoogleCalendarID: &googleCalendarID,
	})
}

func (s *calendarSyncStore) FindAnyCalendarByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*repositorymodel.StoredCalendar, error) {
	return s.repos.Calendar.FindByGoogleCalendarID(ctx, googleCalendarID)
}

func (s *calendarSyncStore) CreateCalendar(ctx context.Context, googleCalendarID, summary string) (*repositorymodel.StoredCalendar, error) {
	return s.repos.Calendar.Create(ctx, repoCalendar.CalendarMutationOptions{
		GoogleCalendarID: &googleCalendarID,
		Summary:          &summary,
	})
}

func (s *calendarSyncStore) UpdateCalendar(ctx context.Context, id uuid.UUID, googleCalendarID, summary string) (*repositorymodel.StoredCalendar, error) {
	return s.repos.Calendar.Update(ctx, id, repoCalendar.CalendarMutationOptions{
		GoogleCalendarID: &googleCalendarID,
		Summary:          &summary,
	})
}

func (s *calendarSyncStore) EnsureUserCalendar(ctx context.Context, userID, calendarID uuid.UUID, role domainvalue.UserCalendarRole) (*repositorymodel.UserCalendar, error) {
	isVisible := true
	syncProposedDates := role == domainvalue.UserCalendarRoleAdjustaCandidate
	return s.repos.UserCalendar.Ensure(ctx, userID, calendarID, repoUserCalendar.UserCalendarQueryOptions{
		Role:              &role,
		IsVisible:         &isVisible,
		SyncProposedDates: &syncProposedDates,
	})
}

func (s *calendarSyncStore) ListCalendarsByUser(ctx context.Context, userID uuid.UUID) ([]*repositorymodel.StoredCalendar, error) {
	return s.repos.Calendar.FilterByUserID(ctx, userID)
}

func (s *calendarSyncStore) SoftDeleteUserCalendar(ctx context.Context, userID, calendarID uuid.UUID) error {
	return s.repos.UserCalendar.SoftDeleteByUserAndCalendar(ctx, userID, calendarID)
}
