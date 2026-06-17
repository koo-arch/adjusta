package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type draftedEventDetailFinder interface {
	FindEventBySlug(ctx context.Context, userID uuid.UUID, slug string, withProposedDates bool) (*EventRecord, error)
}

func (uc *Usecase) loadDraftedEventDetailRecord(ctx context.Context, finder draftedEventDetailFinder, userID uuid.UUID, email, slug string) (*EventRecord, error) {
	storedEvent, err := finder.FindEventBySlug(ctx, userID, slug, true)
	if err != nil {
		log.Printf("failed to get event detail for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedEvent, nil
}

func (uc *Usecase) loadDraftedEventDetailWithSync(ctx context.Context, store EventTxStore, userID uuid.UUID, email, slug string) (*EventRecord, error) {
	storedEvent, err := uc.loadDraftedEventDetailRecord(ctx, store, userID, email, slug)
	if err != nil {
		return nil, err
	}

	candidateCalendar, err := uc.loadAdjustaCandidateCalendar(ctx, store, userID, email)
	if err != nil {
		return nil, err
	}
	if !candidateCalendar.SyncProposedDates || candidateCalendar.GoogleCalendarID == "" {
		return storedEvent, nil
	}

	if err := uc.syncProposedDatesOnDetail(ctx, store, userID, email, candidateCalendar.GoogleCalendarID, storedEvent); err != nil {
		return nil, err
	}

	return uc.loadDraftedEventDetailRecord(ctx, store, userID, email, slug)
}

func (uc *Usecase) syncProposedDatesOnDetail(ctx context.Context, store EventTxStore, userID uuid.UUID, email, calendarID string, storedEvent *EventRecord) error {
	var (
		attemptedSync bool
		lastSyncErr   error
	)

	for _, proposedDate := range storedEvent.ProposedDates {
		if proposedDate == nil {
			continue
		}
		attemptedSync = true

		googleEventID, err := uc.googleCalendar.UpsertEvent(
			ctx,
			userID,
			calendarID,
			proposedDate.GoogleEventID,
			storedEvent.Title,
			storedEvent.Location,
			storedEvent.Description,
			proposedDate.StartTime,
			proposedDate.EndTime,
		)
		if err != nil {
			lastSyncErr = err
			log.Printf("failed to sync proposed date on detail for account: %s, event: %s, proposed date: %s, error: %v", email, storedEvent.ID, proposedDate.ID, err)

			if _, updateErr := store.UpdateProposedDate(ctx, proposedDate.ID, buildProposedDateMutation(domainEvent.NewFailedProposedDateChange(err))); updateErr != nil {
				return fmt.Errorf("failed to mark proposed date sync failure: %w", updateErr)
			}
			continue
		}

		if _, err := store.UpdateProposedDate(ctx, proposedDate.ID, buildProposedDateMutation(domainEvent.NewSyncedProposedDateChange(googleEventID, time.Now()))); err != nil {
			return fmt.Errorf("failed to update proposed date sync success: %w", err)
		}
	}

	if !attemptedSync || storedEvent.Status == domainvalue.StatusConfirmed {
		return nil
	}
	if lastSyncErr != nil {
		if err := uc.recordEventSyncFailure(ctx, store, storedEvent.ID, lastSyncErr); err != nil {
			return fmt.Errorf("failed to mark event sync failure: %w", err)
		}
		return nil
	}

	if _, err := store.UpdateEvent(ctx, storedEvent.ID, mergeEventChange(EventMutation{}, domainEvent.NewSyncedEventSyncChange(time.Now()))); err != nil {
		return fmt.Errorf("failed to update event sync success: %w", err)
	}

	return nil
}
