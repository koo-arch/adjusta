package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*EventDraftDetailOutput, error) {
	storedCalendar, err := uc.loadPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := EventSearchOptions{
		WithProposedDates: true,
	}
	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, storedCalendar.ID, toEventSearchOptions(eventOptions))
	if err != nil {
		log.Printf("failed to get events for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	draftedEvents := make([]*EventDraftDetailOutput, 0, len(storedEvents))
	for _, storedEvent := range storedEvents {
		draft, err := buildEventDraftDetailOutput(storedEvent)
		if err != nil {
			return nil, err
		}
		draftedEvents = append(draftedEvents, draft)
	}

	return draftedEvents, nil
}

func (uc *Usecase) SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query SearchDraftQuery) ([]*EventDraftDetailOutput, error) {
	storedCalendar, err := uc.loadPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := EventSearchOptions{
		WithProposedDates: true,
		Title:             query.Title,
		Location:          query.Location,
		Description:       query.Description,
		Status:            query.Status,
		StartTimeGTE:      query.StartTimeGTE,
		StartTimeLTE:      query.StartTimeLTE,
		EndTimeGTE:        query.EndTimeGTE,
		EndTimeLTE:        query.EndTimeLTE,
	}
	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, storedCalendar.ID, toEventSearchOptions(eventOptions))
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	searchResult := make([]*EventDraftDetailOutput, 0, len(storedEvents))
	for _, storedEvent := range storedEvents {
		draft, err := buildEventDraftDetailOutput(storedEvent)
		if err != nil {
			return nil, err
		}
		searchResult = append(searchResult, draft)
	}

	return searchResult, nil
}

func (uc *Usecase) CreateDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventReq DraftCreationRequest) (*EventDraftDetailOutput, error) {
	var response *EventDraftDetailOutput

	err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
		storedCalendar, err := uc.loadPrimaryCalendar(ctx, repos, userID, email)
		if err != nil {
			return err
		}
		candidateCalendar, err := uc.loadAdjustaCandidateCalendar(ctx, repos, userID, email)
		if err != nil {
			return err
		}

		storedEvent, err := repos.Event.Create(ctx, userID, domainEvent.EventCreateOptions{
			Title:       eventReq.Title,
			Location:    eventReq.Location,
			Description: eventReq.Description,
		}, storedCalendar.ID)
		if err != nil {
			log.Printf("failed to create event for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		storedDates, err := repos.ProposedDate.CreateBulk(ctx, toProposedDateCreateOptionsList(assignSelectedDatePriorities(eventReq.SelectedDates)), storedEvent.ID)
		if err != nil {
			log.Printf("failed to create proposed dates for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		if candidateCalendar.SyncProposedDates {
			storedEvent, err = repos.Event.Update(ctx, storedEvent.ID, mergeEventChange(EventMutation{}, domainEvent.NewPendingEventSyncChange()))
			if err != nil {
				log.Printf("failed to mark event sync pending for account: %s, error: %v", email, err)
				return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
			}

			for i, storedDate := range storedDates {
				storedDates[i], err = repos.ProposedDate.Update(ctx, storedDate.ID, buildProposedDateMutation(domainEvent.NewPendingProposedDateSyncChange()))
				if err != nil {
					log.Printf("failed to mark proposed date sync pending for account: %s, error: %v", email, err)
					return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
				}
			}
		}

		response = &EventDraftDetailOutput{
			ID:            storedEvent.ID,
			Title:         storedEvent.Title,
			Location:      storedEvent.Location,
			Description:   storedEvent.Description,
			Status:        storedEvent.Status,
			SyncStatus:    storedEvent.SyncStatus,
			GoogleEventID: domainEvent.ResolveGoogleEventID(storedEvent.ConfirmedGoogleEventID),
			ProposedDates: buildProposedDateOutputs(storedDates),
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running create drafted event transaction: %v", err)
		return nil, mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return response, nil
}

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
		if eventReq.Status != value.StatusConfirmed {
			eventOptions = mergeEventChange(eventOptions, domainEvent.NewDraftEventChange(&eventReq.Status, candidateCalendar.SyncProposedDates))
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

		if eventReq.Status == value.StatusConfirmed {
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

			confirmedGoogleEventID, err := uc.upsertConfirmedGoogleEvent(ctx, userID, storedCalendar.GoogleCalendarID, storedEvent, confirmation)
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

func (uc *Usecase) DeleteDraftedEvents(ctx context.Context, userID uuid.UUID, email string, eventID uuid.UUID) error {
	err := uc.tx.DoEvent(ctx, func(repos EventTxRepositories) error {
		if _, err := uc.loadPrimaryCalendar(ctx, repos, userID, email); err != nil {
			return err
		}

		if _, err := repos.Event.Update(ctx, eventID, mergeEventChange(EventMutation{}, domainEvent.NewPendingEventSyncChange())); err != nil {
			log.Printf("failed to mark event sync pending for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError("イベント削除時にエラーが発生しました")
		}

		if err := repos.Event.SoftDelete(ctx, eventID); err != nil {
			log.Printf("failed to delete event for account: %s, error: %v", email, err)
			return internalErrors.NewInternalError("イベント削除時にエラーが発生しました")
		}

		return nil
	})
	if err != nil {
		log.Printf("failed running delete drafted event transaction: %v", err)
		return mapUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return nil
}

func (uc *Usecase) updateProposedDates(ctx context.Context, repos EventTxRepositories, eventReq DraftUpdateRequest, storedEvent *domainEvent.Event, existingDates []*domainProposedDate.ProposedDate, syncProposedDates bool) error {
	requestedDates, err := toDomainDraftDateList(eventReq.ProposedDates)
	if err != nil {
		return err
	}

	changeSet := domainEvent.PlanProposedDateChanges(requestedDates, toDomainExistingDateList(existingDates))
	buildChange := func(start, end *time.Time, priority *int, status *value.ProposedDateStatus) domainEvent.ProposedDateChange {
		return domainEvent.NewDraftProposedDateChange(start, end, priority, status, syncProposedDates)
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
