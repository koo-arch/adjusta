package events

import (
	"context"
	"fmt"
	"log"
	"slices"
	"sort"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) (*EventDraftDetailOutput, error) {
	if uc.tx == nil {
		storedEvent, err := uc.loadDraftedEventDetailRecord(ctx, uc.repos, userID, email, eventID)
		if err != nil {
			return nil, err
		}
		return buildEventDraftDetailOutput(storedEvent)
	}

	var response *EventDraftDetailOutput

	err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
		storedEvent, err := uc.loadDraftedEventDetailWithSync(ctx, repos, userID, email, eventID)
		if err != nil {
			return err
		}

		response, err = buildEventDraftDetailOutput(storedEvent)
		return err
	})
	if err != nil {
		log.Printf("failed running fetch drafted event detail transaction: %v", err)
		return nil, mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return response, nil
}

func (uc *Usecase) loadDraftedEventDetailRecord(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string, eventID uuid.UUID) (*domainEvent.Event, error) {
	storedEvent, err := repos.Event.FindByIDAndUser(ctx, userID, eventID, domainEvent.EventReadOptions{
		WithProposedDates: true,
	})
	if err != nil {
		log.Printf("failed to get event detail for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedEvent, nil
}

func (uc *Usecase) loadDraftedEventDetailWithSync(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email string, eventID uuid.UUID) (*domainEvent.Event, error) {
	storedEvent, err := uc.loadDraftedEventDetailRecord(ctx, repos, userID, email, eventID)
	if err != nil {
		return nil, err
	}

	candidateCalendar, err := uc.loadAdjustaCandidateCalendar(ctx, repos, userID, email)
	if err != nil {
		return nil, err
	}
	if !candidateCalendar.SyncProposedDates || candidateCalendar.GoogleCalendarID == "" {
		return storedEvent, nil
	}

	if err := uc.syncProposedDatesOnDetail(ctx, repos, userID, email, candidateCalendar.GoogleCalendarID, storedEvent); err != nil {
		return nil, err
	}

	return uc.loadDraftedEventDetailRecord(ctx, repos, userID, email, eventID)
}

func (uc *Usecase) syncProposedDatesOnDetail(ctx context.Context, repos EventTxRepositories, userID uuid.UUID, email, calendarID string, storedEvent *domainEvent.Event) error {
	var (
		attemptedSync bool
		lastSyncErr   error
	)

	proposedDates := slices.Clone(storedEvent.ProposedDates)
	sort.SliceStable(proposedDates, func(i, j int) bool {
		if proposedDates[i] == nil {
			return false
		}
		if proposedDates[j] == nil {
			return true
		}
		return proposedDates[i].Priority > proposedDates[j].Priority
	})

	candidateRank := 0
	for _, proposedDate := range proposedDates {
		if proposedDate == nil {
			continue
		}
		candidateRank++

		if storedEvent.SyncStatus == value.SyncStatusSynced && proposedDate.SyncStatus == value.SyncStatusSynced {
			continue
		}
		attemptedSync = true

		googleEventID, err := uc.googleCalendar.UpsertEvent(
			ctx,
			userID,
			calendarID,
			proposedDate.GoogleEventID,
			fmt.Sprintf("%s【第%d候補】", storedEvent.Title, candidateRank),
			storedEvent.Location,
			storedEvent.Description,
			proposedDate.StartTime,
			proposedDate.EndTime,
		)
		if err != nil {
			lastSyncErr = err
			log.Printf("failed to sync proposed date on detail for account: %s, event: %s, proposed date: %s, error: %v", email, storedEvent.ID, proposedDate.ID, err)

			if _, updateErr := repos.ProposedDate.Update(ctx, proposedDate.ID, buildProposedDateMutation(domainEvent.NewFailedProposedDateChange(err))); updateErr != nil {
				return fmt.Errorf("failed to mark proposed date sync failure: %w", updateErr)
			}
			continue
		}

		if _, err := repos.ProposedDate.Update(ctx, proposedDate.ID, buildProposedDateMutation(domainEvent.NewSyncedProposedDateChange(googleEventID, time.Now()))); err != nil {
			return fmt.Errorf("failed to update proposed date sync success: %w", err)
		}
	}

	if !attemptedSync || storedEvent.Status == value.StatusConfirmed {
		return nil
	}
	if lastSyncErr != nil {
		if err := uc.recordEventSyncFailure(ctx, repos, storedEvent.ID, lastSyncErr); err != nil {
			return fmt.Errorf("failed to mark event sync failure: %w", err)
		}
		return nil
	}

	if _, err := repos.Event.Update(ctx, storedEvent.ID, mergeEventChange(EventMutation{}, domainEvent.NewSyncedEventSyncChange(time.Now()))); err != nil {
		return fmt.Errorf("failed to update event sync success: %w", err)
	}

	return nil
}
