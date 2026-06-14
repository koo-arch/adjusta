package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *appmodel.ConfirmEvent) error {
	var committedErr error

	err := uc.tx.Do(ctx, func(store EventTxStore) error {
		storedEvent, err := store.FindEventBySlug(ctx, userID, slug, false)
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		storedCalendar, err := store.ReadCalendar(ctx, storedEvent.PrimaryCalendarID)
		if err != nil {
			log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		googleEventID, err := uc.handleGoogleEvent(ctx, userID, storedCalendar.GoogleCalendarID, storedEvent, eventReq)
		if err != nil {
			log.Printf("failed to handle google event for account: %s, error: %v", email, err)
			if syncErr := uc.markEventSyncFailed(ctx, store, storedEvent.ID, err); syncErr != nil {
				log.Printf("failed to mark sync failure for account: %s, error: %v", email, syncErr)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}
			committedErr = internalErrors.NormalizeAPIError(err, "サーバーでエラーが発生しました")
			return nil
		}

		if err := uc.confirmEventDate(ctx, store, googleEventID, eventReq, storedEvent); err != nil {
			log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
			return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running finalize proposed date transaction: %v", err)
		return normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}
	if committedErr != nil {
		return committedErr
	}

	return nil
}

func (uc *Usecase) handleGoogleEvent(ctx context.Context, userID uuid.UUID, calendarID string, storedEvent *EventRecord, eventReq *appmodel.ConfirmEvent) (*string, error) {
	var existingGoogleEventID *string
	if eventReq.ConfirmDate.ID != nil {
		if storedEvent.ConfirmedGoogleEventID != nil && *storedEvent.ConfirmedGoogleEventID != "" {
			existingGoogleEventID = storedEvent.ConfirmedGoogleEventID
		} else if eventReq.ConfirmDate.GoogleEventID != "" {
			existingGoogleEventID = &eventReq.ConfirmDate.GoogleEventID
		} else if storedEvent.GoogleEventID != "" {
			existingGoogleEventID = &storedEvent.GoogleEventID
		}
	}

	googleEventID, err := uc.googleCalendar.UpsertEvent(
		ctx,
		userID,
		calendarID,
		existingGoogleEventID,
		storedEvent.Title,
		storedEvent.Location,
		storedEvent.Description,
		*eventReq.ConfirmDate.Start,
		*eventReq.ConfirmDate.End,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert google event: %w", err)
	}

	return &googleEventID, nil
}

func (uc *Usecase) confirmEventDate(ctx context.Context, store EventTxStore, googleEventID *string, eventReq *appmodel.ConfirmEvent, storedEvent *EventRecord) error {
	confirmDate, err := toDomainConfirmationDate(eventReq.ConfirmDate)
	if err != nil {
		return err
	}

	existingDates, err := store.ListProposedDatesByEvent(ctx, storedEvent.ID)
	if err != nil {
		return fmt.Errorf("failed to list proposed dates error: %w", err)
	}

	plan, err := domainEvent.BuildConfirmationPlan(confirmDate, toDomainExistingProposedDates(existingDates))
	if err != nil {
		return internalErrors.NewBadRequestError("確定候補日程が不正です")
	}

	confirmDateID := eventReq.ConfirmDate.ID
	if eventReq.ConfirmDate.ID == nil {
		dateOptions := withPendingProposedDateSync(ProposedDateMutation{
			Start:    &plan.Create.Start,
			End:      &plan.Create.End,
			Priority: &plan.Create.Priority,
			Status:   proposedDateStatusPtr(domainvalue.ProposedDateStatusConfirmed),
		})

		storedDate, err := store.CreateProposedDate(ctx, dateOptions, storedEvent.ID)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &storedDate.ID
	}

	if eventReq.ConfirmDate.ID != nil {
		dateOptions := withPendingProposedDateSync(ProposedDateMutation{
			Priority: &plan.Update.Priority,
			Status:   proposedDateStatusPtr(domainvalue.ProposedDateStatusConfirmed),
		})

		if _, err := store.UpdateProposedDate(ctx, *eventReq.ConfirmDate.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}
	}

	if err := uc.markUnselectedProposedDates(ctx, store, existingDates, confirmDateID); err != nil {
		return fmt.Errorf("failed to update proposed date statuses: %w", err)
	}

	now := time.Now()
	eventOptions := withSyncedEventSync(EventMutation{
		Status:                 &plan.Status,
		ConfirmedDateID:        confirmDateID,
		GoogleEventID:          googleEventID,
		ConfirmedGoogleEventID: googleEventID,
		LastSyncedAt:           &now,
	})
	if _, err := store.UpdateEvent(ctx, storedEvent.ID, eventOptions); err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}

func (uc *Usecase) markUnselectedProposedDates(ctx context.Context, store EventTxStore, storedDates []*ProposedDateRecord, confirmedDateID *uuid.UUID) error {
	if confirmedDateID == nil {
		return nil
	}

	notSelected := domainvalue.ProposedDateStatusNotSelected
	for _, storedDate := range storedDates {
		if storedDate.ID == *confirmedDateID {
			continue
		}

		if _, err := store.UpdateProposedDate(ctx, storedDate.ID, withPendingProposedDateSync(ProposedDateMutation{
			Status: &notSelected,
		})); err != nil {
			return err
		}
	}

	return nil
}

func proposedDateStatusPtr(status domainvalue.ProposedDateStatus) *domainvalue.ProposedDateStatus {
	return &status
}
