package events

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.GoogleEvent, error) {
	now := time.Now()
	startTime := now.AddDate(0, -2, 0)
	endTime := now.AddDate(1, 0, 0)

	googleCalendars, err := uc.reader.ListGoogleCalendarInfosByUser(ctx, userID)
	if err != nil {
		log.Printf("failed to get google calendars for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	result, err := uc.googleCalendar.FetchEvents(ctx, userID, googleCalendars, startTime, endTime)
	if result == nil {
		return nil, internalErrors.NormalizeAPIError(err, "認証エラーが発生しました")
	}
	if len(result.FailedCalendars) > 0 {
		log.Printf("failed to fetch events from calendars: %v", result.FailedCalendars)
		failedCalendarsMap := map[string][]string{
			"failed_calendars": result.FailedCalendars,
		}

		return result.Events, internalErrors.NewPartialContentError(
			"一部のカレンダーからイベントを取得できませんでした",
			failedCalendarsMap,
		)
	}
	if err != nil && len(result.Events) == 0 {
		log.Printf("failed to fetch events from Google Calendar: %v", err)
		return nil, internalErrors.FromGoogleAPIError(err)
	}

	return result.Events, nil
}

func (uc *Usecase) FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.EventDraftDetail, error) {
	storedCalendar, err := uc.findPrimaryCalendar(ctx, uc.reader, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := EventSearchOptions{
		WithProposedDates: true,
	}
	storedEvents, err := uc.reader.SearchEvents(ctx, userID, storedCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get events for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	draftedEvents := make([]*models.EventDraftDetail, 0, len(storedEvents))
	for _, storedEvent := range storedEvents {
		draft, err := buildEventDraftDetail(storedEvent)
		if err != nil {
			return nil, err
		}
		draftedEvents = append(draftedEvents, draft)
	}

	return draftedEvents, nil
}

func (uc *Usecase) SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query SearchDraftQuery) ([]*models.EventDraftDetail, error) {
	storedCalendar, err := uc.findPrimaryCalendar(ctx, uc.reader, userID, email)
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
	storedEvents, err := uc.reader.SearchEvents(ctx, userID, storedCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	searchResult := make([]*models.EventDraftDetail, 0, len(storedEvents))
	for _, storedEvent := range storedEvents {
		draft, err := buildEventDraftDetail(storedEvent)
		if err != nil {
			return nil, err
		}
		searchResult = append(searchResult, draft)
	}

	return searchResult, nil
}

func (uc *Usecase) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, slug string) (*models.EventDraftDetail, error) {
	storedEvent, err := uc.reader.FindEventBySlug(ctx, userID, slug, true)
	if err != nil {
		log.Printf("failed to get event detail for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return buildEventDraftDetail(storedEvent)
}

func (uc *Usecase) FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.UpcomingEvent, error) {
	storedCalendar, err := uc.findPrimaryCalendar(ctx, uc.reader, userID, email)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, daysBefore)
	confirmed := models.StatusConfirmed
	eventOptions := EventSearchOptions{
		WithProposedDates: true,
		Status:            &confirmed,
		StartTimeGTE:      &currentTime,
		StartTimeLTE:      &startTime,
	}

	storedEvents, err := uc.reader.SearchEvents(ctx, userID, storedCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		return nil, internalErrors.NewInternalError("イベント取得時にエラーが発生しました")
	}

	upcomingEvents := make([]models.UpcomingEvent, 0)
	for _, storedEvent := range storedEvents {
		if storedEvent.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		for _, storedDate := range storedEvent.ProposedDates {
			if storedEvent.ConfirmedDateID == storedDate.ID {
				upcomingEvents = append(upcomingEvents, models.UpcomingEvent{
					ID:              storedEvent.ID,
					Title:           storedEvent.Summary,
					Location:        storedEvent.Location,
					Description:     storedEvent.Description,
					Status:          storedEvent.Status,
					ConfirmedDateID: storedEvent.ConfirmedDateID,
					GoogleEventID:   storedEvent.GoogleEventID,
					Slug:            storedEvent.Slug,
					Start:           storedDate.StartTime,
					End:             storedDate.EndTime,
				})
				break
			}
		}
	}

	sort.Slice(upcomingEvents, func(i, j int) bool {
		return upcomingEvents[i].Start.Before(upcomingEvents[j].Start)
	})

	return upcomingEvents, nil
}

func (uc *Usecase) FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.NeedsActionDraft, error) {
	storedCalendar, err := uc.findPrimaryCalendar(ctx, uc.reader, userID, email)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, daysBefore)
	pending := models.StatusPending
	eventOptions := EventSearchOptions{
		WithProposedDates: true,
		Status:            &pending,
		StartTimeLTE:      &startTime,
		SortBy:            "ProposedDatePriority",
		SortOrder:         "asc",
	}

	storedEvents, err := uc.reader.SearchEvents(ctx, userID, storedCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	needsActionDrafts := make([]models.NeedsActionDraft, 0)
	for _, storedEvent := range storedEvents {
		if storedEvent.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		for _, storedDate := range storedEvent.ProposedDates {
			needsActionDrafts = append(needsActionDrafts, models.NeedsActionDraft{
				ID:             storedEvent.ID,
				Title:          storedEvent.Summary,
				Location:       storedEvent.Location,
				Description:    storedEvent.Description,
				Status:         storedEvent.Status,
				Slug:           storedEvent.Slug,
				Start:          storedDate.StartTime,
				End:            storedDate.EndTime,
				NeedsAttention: currentTime.After(storedDate.StartTime),
			})
			break
		}
	}

	return needsActionDrafts, nil
}
