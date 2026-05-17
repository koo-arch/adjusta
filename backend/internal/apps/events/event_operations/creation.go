package event_operations

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type EventCreationManager struct {
	event *events.EventManager
}

func NewEventCreationManager(event *events.EventManager) *EventCreationManager {
	return &EventCreationManager{
		event: event,
	}
}

func (ecm *EventCreationManager) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftCreation) (*models.EventDraftDetail, error) {

	tx, err := ecm.event.Client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}
	entCalendar, err := ecm.event.CalendarRepo.FindByFields(ctx, tx, userID, findOptions)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	convEvent := ecm.event.CalendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

	entEvent, err := ecm.event.EventRepo.Create(ctx, tx, convEvent, entCalendar)
	if err != nil {
		log.Printf("failed to create event for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	entDate, err := ecm.event.DateRepo.CreateBulk(ctx, tx, eventReq.SelectedDates, entEvent)
	if err != nil {
		log.Printf("failed to create proposed dates for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	eventDates := make([]models.ProposedDate, 0)
	for _, date := range entDate {
		eventDates = append(eventDates, models.ProposedDate{
			ID:        &date.ID,
			Start:     &date.StartTime,
			End:       &date.EndTime,
			Priority: date.Priority,
		})
	}
	response := &models.EventDraftDetail{
		ID:            entEvent.ID,
		Title:         entEvent.Summary,
		Location:      entEvent.Location,
		Description:   entEvent.Description,
		Status:	  	   models.EventStatus(entEvent.Status),
		ProposedDates: eventDates,
	}

	// トランザクションをコミット
	return response, nil
}
