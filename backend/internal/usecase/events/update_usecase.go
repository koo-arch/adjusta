package events

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.EventDraftUpdate) error {
	err := uc.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		if _, err := uc.findPrimaryCalendar(ctx, repos, userID, email); err != nil {
			return err
		}

		entEvent, err := repos.Event.FindBySlugAndUser(ctx, userID, slug, event.EventQueryOptions{})
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
			}
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		eventOptions := event.EventQueryOptions{
			Summary:     &eventReq.Title,
			Location:    &eventReq.Location,
			Description: &eventReq.Description,
			Status:      &eventReq.Status,
		}
		entEvent, err = repos.Event.Update(ctx, entEvent.ID, eventOptions)
		if err != nil {
			log.Printf("failed to update event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
			}
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		existingDates, err := repos.ProposedDate.FilterByEventID(ctx, entEvent.ID)
		if err != nil {
			log.Printf("failed to get proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		if eventReq.Status == models.StatusConfirmed {
			confirmDate := models.ConfirmDate{
				ID:       eventReq.ProposedDates[0].ID,
				Start:    eventReq.ProposedDates[0].Start,
				End:      eventReq.ProposedDates[0].End,
				Priority: eventReq.ProposedDates[0].Priority,
			}
			confirmEvent := models.ConfirmEvent{
				ConfirmDate: confirmDate,
			}

			calendarService, err := uc.getGoogleCalendarService(ctx, userID, email)
			if err != nil {
				return err
			}

			googleEventID, err := uc.handleGoogleEvent(calendarService, entEvent, &confirmEvent)
			if err != nil {
				log.Printf("failed to handle google event for account: %s, error: %v", email, err)
				return utils.HandleGoogleAPIError(err)
			}

			if err := uc.confirmEventDate(ctx, repos, googleEventID, &confirmEvent, entEvent); err != nil {
				log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
				return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
			}
		}

		if err := uc.updateProposedDates(ctx, repos, eventReq, entEvent, existingDates); err != nil {
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running update drafted event transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) updateProposedDates(ctx context.Context, repos internalRepo.Repositories, eventReq *models.EventDraftUpdate, storedEvent *models.StoredEvent, existingDates []*models.StoredProposedDate) error {
	updateDateMap := make(map[uuid.UUID]models.ProposedDate)
	for _, date := range eventReq.ProposedDates {
		if date.ID != nil {
			updateDateMap[*date.ID] = date
		} else {
			updateDateMap[uuid.New()] = date
		}
	}

	for _, date := range existingDates {
		if updateDate, ok := updateDateMap[date.ID]; ok {
			dateOptions := proposeddate.ProposedDateQueryOptions{
				StartTime: updateDate.Start,
				EndTime:   updateDate.End,
				Priority:  &updateDate.Priority,
			}
			if _, err := repos.ProposedDate.Update(ctx, date.ID, dateOptions); err != nil {
				return fmt.Errorf("failed to update proposed date for account: %s, error: %w", updateDate.ID, err)
			}
			delete(updateDateMap, date.ID)
		} else {
			if err := repos.ProposedDate.Delete(ctx, date.ID); err != nil {
				return fmt.Errorf("failed to delete proposed date for account: %s, error: %w", date.ID, err)
			}
		}
	}

	for _, date := range updateDateMap {
		dateOptions := proposeddate.ProposedDateQueryOptions{
			StartTime: date.Start,
			EndTime:   date.End,
			Priority:  &date.Priority,
		}
		if _, err := repos.ProposedDate.Create(ctx, dateOptions, storedEvent.ID); err != nil {
			return fmt.Errorf("failed to create proposed date for account: %s, error: %w", date.ID, err)
		}
	}

	return nil
}
