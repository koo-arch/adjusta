package events

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) FetchAllGoogleEvents(ctx context.Context, userID uuid.UUID, email string) ([]*FetchedGoogleEvent, error) {
	now := time.Now()
	startTime := now.AddDate(0, -2, 0)
	endTime := now.AddDate(1, 0, 0)

	googleCalendars, err := uc.repos.Calendar.FilterByUserID(ctx, userID)
	if err != nil {
		log.Printf("failed to get google calendars for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	result, err := uc.googleCalendar.FetchEvents(ctx, userID, toEventCalendars(googleCalendars), startTime, endTime)
	if result == nil {
		return nil, internalErrors.NormalizeAPIError(err, "Googleカレンダーのイベント取得に失敗しました")
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
		return nil, internalErrors.NormalizeAPIError(err, "Googleカレンダーのイベント取得に失敗しました")
	}

	return result.Events, nil
}

func (uc *Usecase) upsertConfirmedGoogleEvent(ctx context.Context, userID uuid.UUID, calendarID string, storedEvent *domainEvent.Event, confirmation ConfirmationRequest) (*string, error) {
	existingGoogleEventID := domainEvent.ResolveReusableGoogleEventID(
		confirmation.ID,
		storedEvent.ConfirmedGoogleEventID,
		confirmation.GoogleEventID,
	)

	googleEventID, err := uc.googleCalendar.UpsertEvent(
		ctx,
		userID,
		calendarID,
		existingGoogleEventID,
		storedEvent.Title,
		storedEvent.Location,
		storedEvent.Description,
		*confirmation.Start,
		*confirmation.End,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert google event: %w", err)
	}

	return &googleEventID, nil
}
