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

func (uc *Usecase) UpdateDraftedEvents(ctx context.Context, userID uuid.UUID, eventID uuid.UUID, email string, eventReq DraftUpdateRequest) error {
	var committedErr error

	err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
		if _, err := uc.loadPrimaryCalendar(ctx, repos, userID, email); err != nil {
			return err
		}
		candidateCalendar, err := uc.loadAdjustaCandidateCalendar(ctx, repos, userID, email)
		if err != nil {
			return err
		}

		storedEvent, err := repos.Event.FindByIDAndUser(ctx, userID, eventID, domainEvent.EventReadOptions{
			WithProposedDates: false,
		})
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		eventOptions := EventMutation{
			Title:       &eventReq.Title,
			Location:    &eventReq.Location,
			Description: &eventReq.Description,
		}
		if eventReq.Status != domainvalue.StatusConfirmed {
			if candidateCalendar.SyncProposedDates {
				eventOptions = mergeEventChange(eventOptions, domainEvent.NewPendingEventChange(&eventReq.Status))
			} else {
				eventOptions = mergeEventChange(eventOptions, domainEvent.NewNotSyncedEventChange(&eventReq.Status))
			}
		}
		storedEvent, err = repos.Event.Update(ctx, storedEvent.ID, eventOptions)
		if err != nil {
			log.Printf("failed to update event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewNotFoundError("イベントが見つかりませんでした")
			}
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		existingDates, err := repos.ProposedDate.FilterByEventID(ctx, storedEvent.ID)
		if err != nil {
			log.Printf("failed to get proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		if eventReq.Status == domainvalue.StatusConfirmed {
			storedCalendar, err := repos.Calendar.Read(ctx, storedEvent.PrimaryCalendarID)
			if err != nil {
				log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}

			confirmDate := eventReq.ProposedDates[0]
			googleEventID := ""
			if confirmDate.GoogleEventID != nil {
				googleEventID = *confirmDate.GoogleEventID
			}
			confirmation := ConfirmationRequest{
				ID:            confirmDate.ID,
				GoogleEventID: googleEventID,
				Start:         confirmDate.Start,
				End:           confirmDate.End,
				Priority:      confirmDate.Priority,
			}

			confirmedGoogleEventID, err := uc.handleGoogleEvent(ctx, userID, storedCalendar.GoogleCalendarID, storedEvent, confirmation)
			if err != nil {
				log.Printf("failed to handle google event for account: %s, error: %v", email, err)
				if syncErr := uc.recordEventSyncFailure(ctx, repos, storedEvent.ID, err); syncErr != nil {
					log.Printf("failed to mark sync failure for account: %s, error: %v", email, syncErr)
					return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
				}
				committedErr = internalErrors.NormalizeAPIError(err, "Googleカレンダー更新時にエラーが発生しました")
				return nil
			}

			if err := uc.confirmEventDate(ctx, repos, confirmedGoogleEventID, confirmation, storedEvent); err != nil {
				log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
				return mapUsecaseError(err, internalErrors.InternalErrorMessage)
			}
		}

		if err := uc.updateProposedDates(ctx, repos, eventReq, storedEvent, existingDates, candidateCalendar.SyncProposedDates); err != nil {
			return mapUsecaseError(err, internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running update drafted event transaction: %v", err)
		return mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}
	if committedErr != nil {
		return committedErr
	}

	return nil
}

func (uc *Usecase) updateProposedDates(ctx context.Context, repos EventTxRepositories, eventReq DraftUpdateRequest, storedEvent *EventRecord, existingDates []*ProposedDateRecord, syncProposedDates bool) error {
	requestedDates, err := toDomainDraftDateList(eventReq.ProposedDates)
	if err != nil {
		return err
	}

	changeSet := domainEvent.PlanProposedDateChanges(requestedDates, toDomainExistingDateList(existingDates))
	buildChange := func(start, end *time.Time, priority *int, status *domainvalue.ProposedDateStatus) domainEvent.ProposedDateChange {
		if syncProposedDates {
			return domainEvent.NewPendingProposedDateChange(start, end, priority, status)
		}
		return domainEvent.NewNotSyncedProposedDateChange(start, end, priority, status)
	}

	for _, date := range changeSet.Updates {
		dateOptions := buildProposedDateMutation(buildChange(&date.Start, &date.End, &date.Priority, nil))
		if _, err := repos.ProposedDate.Update(ctx, date.ID, dateOptions); err != nil {
			return fmt.Errorf("failed to update proposed date for account: %s, error: %w", date.ID, err)
		}
	}

	for _, dateID := range changeSet.Deletes {
		if _, err := repos.ProposedDate.Update(ctx, dateID, buildProposedDateMutation(buildChange(nil, nil, nil, nil))); err != nil {
			return fmt.Errorf("failed to update proposed date sync state for account: %s, error: %w", dateID, err)
		}
		if err := repos.ProposedDate.SoftDelete(ctx, dateID); err != nil {
			return fmt.Errorf("failed to delete proposed date for account: %s, error: %w", dateID, err)
		}
	}

	for _, date := range changeSet.Creates {
		dateOptions := buildProposedDateMutation(buildChange(&date.Start, &date.End, &date.Priority, nil))
		createOptions, err := toProposedDateCreateOptions(dateOptions)
		if err != nil {
			return err
		}
		if _, err := repos.ProposedDate.Create(ctx, createOptions, storedEvent.ID); err != nil {
			return fmt.Errorf("failed to create proposed date for event: %s, error: %w", storedEvent.ID, err)
		}
	}

	return nil
}
