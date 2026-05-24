package events

import (
	"context"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) findPrimaryCalendar(ctx context.Context, repos internalRepo.Repositories, userID uuid.UUID, email string) (*models.StoredCalendar, error) {
	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	entCalendar, err := repos.Calendar.FindByFields(ctx, userID, findOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return entCalendar, nil
}

func (uc *Usecase) getGoogleCalendarService(ctx context.Context, userID uuid.UUID, email string) (*customCalendar.Calendar, error) {
	token, err := uc.googleTokenManager.GetToken(ctx, userID)
	if err != nil {
		log.Printf("failed to verify token for account: %s, error: %v", email, err)
		return nil, utils.GetAPIError(err, "OAuthトークンの認証に失敗しました")
	}

	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		log.Printf("failed to create calendar service for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "Googleカレンダー接続に失敗しました")
	}

	return calendarService, nil
}

func buildProposedDates(storedDates []*models.StoredProposedDate) []models.ProposedDate {
	proposedDates := make([]models.ProposedDate, 0, len(storedDates))
	for _, storedDate := range storedDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:       &storedDate.ID,
			Start:    &storedDate.StartTime,
			End:      &storedDate.EndTime,
			Priority: storedDate.Priority,
		})
	}

	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority < proposedDates[j].Priority
	})

	return proposedDates
}

func buildEventDraftDetail(storedEvent *models.StoredEvent) (*models.EventDraftDetail, error) {
	if storedEvent.ProposedDates == nil {
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return &models.EventDraftDetail{
		ID:              storedEvent.ID,
		Title:           storedEvent.Summary,
		Location:        storedEvent.Location,
		Description:     storedEvent.Description,
		Status:          storedEvent.Status,
		ConfirmedDateID: &storedEvent.ConfirmedDateID,
		GoogleEventID:   storedEvent.GoogleEventID,
		Slug:            storedEvent.Slug,
		ProposedDates:   buildProposedDates(storedEvent.ProposedDates),
	}, nil
}

func normalizeUsecaseError(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}

	if apiErr, ok := err.(*internalErrors.APIError); ok {
		return apiErr
	}

	return internalErrors.NewAPIError(http.StatusInternalServerError, fallbackMessage)
}
