package event_operations

import (
	"net/http"
	"context"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"

	"github.com/koo-arch/adjusta-backend/internal/apps/events"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/utils"
)

type EventFetchingManager struct {
	event *events.EventManager
}

func NewEventFetchingManager(event *events.EventManager) *EventFetchingManager {
	return &EventFetchingManager{
		event: event,
	}
}

func (efm *EventFetchingManager) FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.GoogleEvent, error) {

	token, err := efm.event.AuthManager.VerifyOAuthToken(ctx, userID)
	if err != nil {
		log.Printf("failed to verify token for account: %s, error: %v", email, err)
		apiErr := utils.GetAPIError(err, "認証エラーが発生しました")
		return nil, apiErr
	}

	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		log.Printf("failed to connect to Google Calendar: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, "Googleカレンダーの接続に失敗しました")
	}

	now := time.Now()
	startTime := now.AddDate(0, -2, 0)
	endTime := now.AddDate(1, 0, 0)

	calendarOptions := repoCalendar.CalendarQueryOptions{
		WithGoogleCalendarInfo: true,
	}
	calendars, err := efm.event.CalendarRepo.FilterByFields(ctx, nil, userID, calendarOptions)
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

	result, err := efm.event.CalendarApp.FetchEventsFromCalendars(calendarService, entGoogleCalendars, startTime, endTime)
	if len(result.FailedCalendars) > 0 {
		log.Printf("failed to fetch events from calendars: %v", result.FailedCalendars)

		// 失敗したカレンダーの情報を
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
		apiErr := utils.HandleGoogleAPIError(err)
		return nil, apiErr
	}

	return result.Events, nil
}

func (efm *EventFetchingManager) FetchAllDraftedEvents(ctx context.Context, userID uuid.UUID, email string) ([]*models.EventDraftDetail, error) {
	isPrimary := true
	findOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary:         &isPrimary,
		WithEvents:        true,
		WithProposedDates: true,
	}
	entCalendar, err := efm.event.CalendarRepo.FindByFields(ctx, nil, userID, findOptions)
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

	draftedEvents := make([]*models.EventDraftDetail, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errCh := make(chan error, 1)

	for _, entEvent := range entCalendar.Edges.Events {
		wg.Add(1)

		go func(entEvent *ent.Event) {
			defer wg.Done()
			var proposedDates []models.ProposedDate

			if entEvent.Edges.ProposedDates == nil {
				return
			}
			for _, entDate := range entEvent.Edges.ProposedDates {
				proposedDates = append(proposedDates, models.ProposedDate{
					ID:       &entDate.ID,
					Start:    &entDate.StartTime,
					End:      &entDate.EndTime,
					Priority: entDate.Priority,
				})
			}

			// Priorityに基づいてProposedDatesを昇順にソート
			sort.Slice(proposedDates, func(i, j int) bool {
				return proposedDates[i].Priority < proposedDates[j].Priority
			})

			// 同時に書き込むことがないようにミューテックスを使う
			mu.Lock()
			draftedEvents = append(draftedEvents, &models.EventDraftDetail{
				ID:              entEvent.ID,
				Title:           entEvent.Summary,
				Location:        entEvent.Location,
				Description:     entEvent.Description,
				Status:          models.EventStatus(entEvent.Status),
				ConfirmedDateID: &entEvent.ConfirmedDateID,
				Slug: 		  	 entEvent.Slug,
				GoogleEventID:   entEvent.GoogleEventID,
				ProposedDates:   proposedDates,
			})
			mu.Unlock()
		}(entEvent)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		return nil, <-errCh
	}

	return draftedEvents, nil
}

func (efm *EventFetchingManager) SearchDraftedEvents(ctx context.Context, userID uuid.UUID, email string, query event.EventQueryOptions) ([]*models.EventDraftDetail, error) {
	isPrimary := true
	calendarOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	entCalendar, err := efm.event.CalendarRepo.FindByFields(ctx, nil, userID, calendarOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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
	entEvent, err := efm.event.EventRepo.SearchEvents(ctx, nil, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	searchResult := make([]*models.EventDraftDetail, 0)
	for _, event := range entEvent {
		if event.Edges.ProposedDates == nil {
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		proposedDates := make([]models.ProposedDate, 0)
		for _, entDate := range event.Edges.ProposedDates {
			proposedDates = append(proposedDates, models.ProposedDate{
				ID:       &entDate.ID,
				Start:    &entDate.StartTime,
				End:      &entDate.EndTime,
				Priority: entDate.Priority,
			})
		}

		// Priorityに基づいてProposedDatesを昇順にソート
		sort.Slice(proposedDates, func(i, j int) bool {
			return proposedDates[i].Priority < proposedDates[j].Priority
		})

		searchResult = append(searchResult, &models.EventDraftDetail{
			ID:              event.ID,
			Title:           event.Summary,
			Location:        event.Location,
			Description:     event.Description,
			Status:          models.EventStatus(event.Status),
			ConfirmedDateID: &event.ConfirmedDateID,
			GoogleEventID:   event.GoogleEventID,
			Slug :           event.Slug,
			ProposedDates:   proposedDates,
		})
	}

	return searchResult, nil
}

func (efm *EventFetchingManager) FetchDraftedEventDetail(ctx context.Context, userID uuid.UUID, email string, slug string) (*models.EventDraftDetail, error) {
	tx, err := efm.event.Client.Tx(ctx)
	if err != nil {
		log.Printf("failed starting transaction: %v", err)
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	defer transaction.HandleTransaction(tx, &err)

	queryOpt := event.EventQueryOptions{
		WithProposedDates: true,
	}
	entEvent, err := efm.event.EventRepo.FindBySlugAndUser(ctx, tx, userID, slug, queryOpt)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	if entEvent.Edges.ProposedDates == nil {
		log.Printf("No association found between calendar and event")
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	proposedDates := make([]models.ProposedDate, 0)
	for _, entDate := range entEvent.Edges.ProposedDates {
		proposedDates = append(proposedDates, models.ProposedDate{
			ID:       &entDate.ID,
			Start:    &entDate.StartTime,
			End:      &entDate.EndTime,
			Priority: entDate.Priority,
		})
	}

	// Priorityに基づいてProposedDatesを昇順にソート
	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority < proposedDates[j].Priority
	})

	return &models.EventDraftDetail{
		ID:              entEvent.ID,
		Title:           entEvent.Summary,
		Location:        entEvent.Location,
		Description:     entEvent.Description,
		Status:          models.EventStatus(entEvent.Status),
		ConfirmedDateID: &entEvent.ConfirmedDateID,
		GoogleEventID:   entEvent.GoogleEventID,
		Slug:            entEvent.Slug,
		ProposedDates:   proposedDates,
	}, nil
}

func (efm *EventFetchingManager) FetchUpcomingEvents(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.UpcomingEvent, error) {
	isPrimary := true
	calendarOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	entCalendar, err := efm.event.CalendarRepo.FindByFields(ctx, nil, userID, calendarOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
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

	entEvents, err := efm.event.EventRepo.SearchEvents(ctx, nil, userID, entCalendar.ID, eventOptions)
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

	// 開始日時で昇順にソート
	sort.Slice(upcomingEvents, func(i, j int) bool {
		return upcomingEvents[i].Start.Before(upcomingEvents[j].Start)
	})

	return upcomingEvents, nil
}

func (efm *EventFetchingManager) FetchNeedsActionDrafts(ctx context.Context, userID uuid.UUID, email string, daysBefore int) ([]models.NeedsActionDraft, error) {
	isPrimary := true
	calendarOptions := repoCalendar.CalendarQueryOptions{
		IsPrimary: &isPrimary,
	}

	entCalendar, err := efm.event.CalendarRepo.FindByFields(ctx, nil, userID, calendarOptions)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	currentTime := time.Now()
	startTime := currentTime.AddDate(0, 0, daysBefore)
	draft := models.StatusPending
	eventOptions := event.EventQueryOptions{
		WithProposedDates:    true,
		Status:               &draft,
		ProposedDateStartLTE: &startTime,
		SortBy:               "ProposedDatePriority",
		SortOrder:            "asc",
	}

	entEvents, err := efm.event.EventRepo.SearchEvents(ctx, nil, userID, entCalendar.ID, eventOptions)
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return nil, internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}
		return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	NeedsActionDrafts:= make([]models.NeedsActionDraft, 0)
	for _, entEvent := range entEvents {
		if entEvent.Edges.ProposedDates == nil {
			log.Printf("No association found between calendar and event")
			return nil, internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
		}

		for _, entDate := range entEvent.Edges.ProposedDates {
			// 開始日時が現在時刻よりも前の場合はisPastをtrueにする
			isPast := currentTime.After(entDate.StartTime)
			NeedsActionDrafts = append(NeedsActionDrafts, models.NeedsActionDraft{
				ID:             entEvent.ID,
				Title:          entEvent.Summary,
				Location:       entEvent.Location,
				Description:    entEvent.Description,
				Status:         models.EventStatus(entEvent.Status),
				Slug:           entEvent.Slug,
				Start:          entDate.StartTime,
				End:            entDate.EndTime,
				NeedsAttention: isPast,
			})
			break
		}
	}

	// 開始時刻で昇順にソート
	sort.Slice(NeedsActionDrafts, func(i, j int) bool {
		return NeedsActionDrafts[i].Start.Before(NeedsActionDrafts[j].Start)
	})

	return NeedsActionDrafts, nil
}
