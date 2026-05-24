package events

import (
	"context"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) findPrimaryCalendar(ctx context.Context, tx *ent.Tx, userID uuid.UUID, email string) (*ent.Calendar, error) {
	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	entCalendar, err := uc.calendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
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

func buildProposedDates(entDates []*ent.ProposedDate) []models.ProposedDate {
	proposedDates := make([]models.ProposedDate, 0, len(entDates))
	for _, entDate := range entDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:       &entDate.ID,
			Start:    &entDate.StartTime,
			End:      &entDate.EndTime,
			Priority: entDate.Priority,
		})
	}

	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority < proposedDates[j].Priority
	})

	return proposedDates
}

func buildEventDraftDetail(entEvent *ent.Event) (*models.EventDraftDetail, error) {
	if entEvent.Edges.ProposedDates == nil {
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return &models.EventDraftDetail{
		ID:              entEvent.ID,
		Title:           entEvent.Summary,
		Location:        entEvent.Location,
		Description:     entEvent.Description,
		Status:          models.EventStatus(entEvent.Status),
		ConfirmedDateID: &entEvent.ConfirmedDateID,
		GoogleEventID:   entEvent.GoogleEventID,
		Slug:            entEvent.Slug,
		ProposedDates:   buildProposedDates(entEvent.Edges.ProposedDates),
	}, nil
}
