package events

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
)

func (uc *Usecase) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq *models.EventDraftCreation) (*models.EventDraftDetail, error) {
	var response *models.EventDraftDetail

	err := uc.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		entCalendar, err := uc.findPrimaryCalendar(ctx, repos, userID, email)
		if err != nil {
			return err
		}

		convEvent := uc.calendarApp.ConvertToCalendarEvent(nil, eventReq.Title, eventReq.Location, eventReq.Description, eventReq.SelectedDates[0].Start, eventReq.SelectedDates[0].End)

		storedEvent, err := repos.Event.Create(ctx, convEvent, entCalendar.ID)
		if err != nil {
			log.Printf("failed to create event for account: %s, error: %v", email, err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		storedDates, err := repos.ProposedDate.CreateBulk(ctx, eventReq.SelectedDates, storedEvent.ID)
		if err != nil {
			log.Printf("failed to create proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		response = &models.EventDraftDetail{
			ID:            storedEvent.ID,
			Title:         storedEvent.Summary,
			Location:      storedEvent.Location,
			Description:   storedEvent.Description,
			Status:        storedEvent.Status,
			ProposedDates: buildProposedDates(storedDates),
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running create drafted event transaction: %v", err)
		return nil, normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return response, nil
}
