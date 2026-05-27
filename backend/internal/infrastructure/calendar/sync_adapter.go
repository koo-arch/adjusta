package calendar

import (
	"context"

	"github.com/google/uuid"
	infraRepository "github.com/koo-arch/adjusta-backend/internal/infrastructure/repository"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	repoGoogleCalendarInfo "github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	repoUser "github.com/koo-arch/adjusta-backend/internal/repo/user"
	usecaseCalendar "github.com/koo-arch/adjusta-backend/internal/usecase/calendar"
)

type calendarSyncUserReader struct {
	userRepo repoUser.UserRepository
}

func NewCalendarSyncUserReader(userRepo repoUser.UserRepository) usecaseCalendar.UserReader {
	return &calendarSyncUserReader{userRepo: userRepo}
}

func (r *calendarSyncUserReader) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
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

func (s *calendarSyncStore) FindCalendarByGoogleCalendarID(ctx context.Context, userID uuid.UUID, googleCalendarID string) (*models.StoredCalendar, error) {
	return s.repos.Calendar.FindByFields(ctx, userID, repoCalendar.CalendarQueryOptions{
		WithGoogleCalendarInfo: true,
		GoogleCalendarID:       &googleCalendarID,
	})
}

func (s *calendarSyncStore) CreateCalendar(ctx context.Context, userID uuid.UUID) (*models.StoredCalendar, error) {
	return s.repos.Calendar.Create(ctx, userID)
}

func (s *calendarSyncStore) FindGoogleCalendarInfoByGoogleCalendarID(ctx context.Context, googleCalendarID string) (*models.GoogleCalendarInfo, error) {
	return s.repos.GoogleCalendarInfo.FindByFields(ctx, repoGoogleCalendarInfo.GoogleCalendarInfoQueryOptions{
		GoogleCalendarID: &googleCalendarID,
	})
}

func (s *calendarSyncStore) CreateGoogleCalendarInfo(ctx context.Context, googleCalendarID, summary string, isPrimary bool, calendarID uuid.UUID) (*models.GoogleCalendarInfo, error) {
	return s.repos.GoogleCalendarInfo.Create(ctx, repoGoogleCalendarInfo.GoogleCalendarInfoQueryOptions{
		GoogleCalendarID: &googleCalendarID,
		Summary:          &summary,
		IsPrimary:        &isPrimary,
	}, calendarID)
}

func (s *calendarSyncStore) LinkGoogleCalendarInfoToCalendar(ctx context.Context, googleCalendarInfoID, calendarID uuid.UUID) error {
	_, err := s.repos.GoogleCalendarInfo.Update(ctx, googleCalendarInfoID, repoGoogleCalendarInfo.GoogleCalendarInfoQueryOptions{}, &calendarID)
	return err
}

func (s *calendarSyncStore) ListGoogleCalendarInfosByUser(ctx context.Context, userID uuid.UUID) ([]*models.GoogleCalendarInfo, error) {
	return s.repos.GoogleCalendarInfo.ListByUser(ctx, userID)
}

func (s *calendarSyncStore) SoftDeleteGoogleCalendarInfo(ctx context.Context, id uuid.UUID) error {
	return s.repos.GoogleCalendarInfo.SoftDelete(ctx, id)
}
