package events

import (
	"context"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	internalRepo "github.com/koo-arch/adjusta-backend/internal/repo"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	"github.com/koo-arch/adjusta-backend/utils"
)

func (uc *Usecase) FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.GoogleEvent, error) {
	calendarService, err := uc.getGoogleCalendarService(ctx, userID, email)
	if err != nil {
		return nil, utils.GetAPIError(err, "認証エラーが発生しました")
	}

	now := time.Now()
	startTime := now.AddDate(0, -2, 0)
	endTime := now.AddDate(1, 0, 0)

	googleCalendars, err := uc.repos.GoogleCalendarInfo.ListByUser(ctx, userID)
	if err != nil {
		log.Printf("failed to get google calendars for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	result, err := uc.calendarApp.FetchEventsFromCalendars(calendarService, googleCalendars, startTime, endTime)
	if len(result.FailedCalendars) > 0 {
		log.Printf("failed to fetch events from calendars: %v", result.FailedCalendars)
		failedCalendarsMap := map[string][]string{
			"failed_calendars": result.FailedCalendars,
		}

		return result.Events, internalErrors.NewAPIErrorWithDetails(
			http.StatusPartialContent,
			"一部のカレンダーからイベントを取得できませんでした",
			failedCalendarsMap,
		)
	}
	if err != nil && len(result.Events) == 0 {
		log.Printf("failed to fetch events from Google Calendar: %v", err)
		return nil, utils.HandleGoogleAPIError(err)
	}

	return result.Events, nil
}

func (uc *Usecase) FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.EventDraftDetail, error) {
	entCalendar, err := uc.findPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := event.EventQueryOptions{
		WithProposedDates: true,
	}
	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get events for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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
	entCalendar, err := uc.findPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := event.EventQueryOptions{
		WithProposedDates:    true,
		Summary:              query.Title,
		Location:             query.Location,
		Description:          query.Description,
		Status:               query.Status,
		ProposedDateStartGTE: query.StartTimeGTE,
		ProposedDateStartLTE: query.StartTimeLTE,
		ProposedDateEndGTE:   query.EndTimeGTE,
		ProposedDateEndLTE:   query.EndTimeLTE,
	}
	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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
	var draftedEvent *models.EventDraftDetail

	err := uc.uow.Do(ctx, func(repos internalRepo.Repositories) error {
		queryOpt := event.EventQueryOptions{
			WithProposedDates: true,
		}
		storedEvent, err := repos.Event.FindBySlugAndUser(ctx, userID, slug, queryOpt)
		if err != nil {
			log.Printf("failed to get event for account: %s, error: %v", email, err)
			if repoerr.IsNotFound(err) {
				return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
			}
			return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		draftedEvent, err = buildEventDraftDetail(storedEvent)
		return err
	})
	if err != nil {
		log.Printf("failed running fetch drafted event detail transaction: %v", err)
		return nil, normalizeUsecaseError(err, internalErrors.InternalErrorMessage)
	}

	return draftedEvent, nil
}

func (uc *Usecase) FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.UpcomingEvent, error) {
	entCalendar, err := uc.findPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, daysBefore)
	confirmed := models.StatusConfirmed
	eventOptions := event.EventQueryOptions{
		WithProposedDates:    true,
		Status:               &confirmed,
		ProposedDateStartGTE: &currentTime,
		ProposedDateStartLTE: &startTime,
	}

	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "イベント取得時にエラーが発生しました")
	}

	upcomingEvents := make([]models.UpcomingEvent, 0)
	for _, storedEvent := range storedEvents {
		if storedEvent.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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
	entCalendar, err := uc.findPrimaryCalendar(ctx, uc.repos, userID, email)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, daysBefore)
	pending := models.StatusPending
	eventOptions := event.EventQueryOptions{
		WithProposedDates:    true,
		Status:               &pending,
		ProposedDateStartLTE: &startTime,
		SortBy:               "ProposedDatePriority",
		SortOrder:            "asc",
	}

	storedEvents, err := uc.repos.Event.SearchEvents(ctx, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	needsActionDrafts := make([]models.NeedsActionDraft, 0)
	for _, storedEvent := range storedEvents {
		if storedEvent.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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
