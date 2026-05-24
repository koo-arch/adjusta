package events

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

func (uc *Usecase) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftCreation) (*models.EventDraftDetail, error) {
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	entCalendar, err := uc.calendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	convEvent := uc.calendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

	entEvent, err := uc.eventRepo.Create(ctx, tx, convEvent, entCalendar)
	if err != nil {
		log.Printf("failed to create event for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	entDates, err := uc.dateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates, entEvent)
	if err != nil {
		log.Printf("failed to create proposed dates for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	eventDates := make([]models.ProposedDate, 0, len(entDates))
	for _, date := range entDates {
		eventDates = append(eventDates, models.ProposedDate{
			ID:       &date.ID,
			Start:    &date.StartTime,
			End:      &date.EndTime,
			Priority: date.Priority,
		})
	}

	return &models.EventDraftDetail{
		ID:            entEvent.ID,
		Title:         entEvent.Summary,
		Location:      entEvent.Location,
		Description:   entEvent.Description,
		Status:        models.EventStatus(entEvent.Status),
		ProposedDates: eventDates,
	}, nil
}
