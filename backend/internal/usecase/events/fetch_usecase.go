package events

import (
	"context"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
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

	calendarOptions := repoCalendarOptionsWithGoogleInfo()
	calendars, err := uc.calendarRepo.FilterByFields(ctx, nil, userID, calendarOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	entGoogleCalendars := make([]*ent.GoogleCalendarInfo, 0)
	for _, cal := range calendars {
		if cal.Edges.GoogleCalendarInfos != nil {
			entGoogleCalendars = append(entGoogleCalendars, cal.Edges.GoogleCalendarInfos...)
		}
	}

	result, err := uc.calendarApp.FetchEventsFromCalendars(calendarService, entGoogleCalendars, startTime, endTime)
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
	isPrimary := true
	findOptions := repoCalendarOptionsWithEvents(isPrimary)

	entCalendar, err := uc.calendarRepo.FindByFields(ctx, nil, userID, findOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "カレンダー取得時にエラーが発生しました")
	}

	if entCalendar.Edges.Events == nil {
		log.Printf("No association found between calendar and event")
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	draftedEvents := make([]*models.EventDraftDetail, 0, len(entCalendar.Edges.Events))
	for _, entEvent := range entCalendar.Edges.Events {
		draft, err := buildEventDraftDetail(entEvent)
		if err != nil {
			return nil, err
		}
		draftedEvents = append(draftedEvents, draft)
	}

	return draftedEvents, nil
}

func (uc *Usecase) SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query event.EventQueryOptions) ([]*models.EventDraftDetail, error) {
	entCalendar, err := uc.findPrimaryCalendar(ctx, nil, userID, email)
	if err != nil {
		return nil, err
	}

	eventOptions := event.EventQueryOptions{
		WithProposedDates:    true,
		Summary:              query.Summary,
		Location:             query.Location,
		Description:          query.Description,
		Status:               query.Status,
		ProposedDateStartGTE: query.ProposedDateStartGTE,
		ProposedDateEndLTE:   query.ProposedDateEndLTE,
	}
	entEvents, err := uc.eventRepo.SearchEvents(ctx, nil, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	searchResult := make([]*models.EventDraftDetail, 0, len(entEvents))
	for _, entEvent := range entEvents {
		draft, err := buildEventDraftDetail(entEvent)
		if err != nil {
			return nil, err
		}
		searchResult = append(searchResult, draft)
	}

	return searchResult, nil
}

func (uc *Usecase) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, slug string) (*models.EventDraftDetail, error) {
	tx, err := uc.client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	queryOpt := event.EventQueryOptions{
		WithProposedDates: true,
	}
	entEvent, err := uc.eventRepo.FindBySlugAndUser(ctx, tx, userID, slug, queryOpt)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	return buildEventDraftDetail(entEvent)
}

func (uc *Usecase) FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.UpcomingEvent, error) {
	entCalendar, err := uc.findPrimaryCalendar(ctx, nil, userID, email)
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

	entEvents, err := uc.eventRepo.SearchEvents(ctx, nil, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "イベント取得時にエラーが発生しました")
	}

	upcomingEvents := make([]models.UpcomingEvent, 0)
	for _, entEvent := range entEvents {
		if entEvent.Edges.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		for _, entDate := range entEvent.Edges.ProposedDates {
			if entEvent.ConfirmedDateID == entDate.ID {
				upcomingEvents = append(upcomingEvents, models.UpcomingEvent{
					ID:              entEvent.ID,
					Title:           entEvent.Summary,
					Location:        entEvent.Location,
					Description:     entEvent.Description,
					Status:          models.EventStatus(entEvent.Status),
					ConfirmedDateID: entEvent.ConfirmedDateID,
					GoogleEventID:   entEvent.GoogleEventID,
					Slug:            entEvent.Slug,
					Start:           entDate.StartTime,
					End:             entDate.EndTime,
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
	entCalendar, err := uc.findPrimaryCalendar(ctx, nil, userID, email)
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

	entEvents, err := uc.eventRepo.SearchEvents(ctx, nil, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	needsActionDrafts := make([]models.NeedsActionDraft, 0)
	for _, entEvent := range entEvents {
		if entEvent.Edges.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		for _, entDate := range entEvent.Edges.ProposedDates {
			needsActionDrafts = append(needsActionDrafts, models.NeedsActionDraft{
				ID:             entEvent.ID,
				Title:          entEvent.Summary,
				Location:       entEvent.Location,
				Description:    entEvent.Description,
				Status:         models.EventStatus(entEvent.Status),
				Slug:           entEvent.Slug,
				Start:          entDate.StartTime,
				End:            entDate.EndTime,
				NeedsAttention: currentTime.After(entDate.StartTime),
			})
			break
		}
	}

	return needsActionDrafts, nil
}

func repoCalendarOptionsWithGoogleInfo() repoCalendar.CalendarQueryOptions {
	return repoCalendar.CalendarQueryOptions{
		WithGoogleCalendarInfo: true,
	}
}

func repoCalendarOptionsWithEvents(isPrimary bool) repoCalendar.CalendarQueryOptions {
	return repoCalendar.CalendarQueryOptions{
		IsPrimary:         &isPrimary,
		WithEvents:        true,
		WithProposedDates: true,
	}
}
