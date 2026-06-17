package events

import (
	"context"
	"log"

	"github.com/google/uuid"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

func (uc *Usecase) loadPrimaryCalendar(ctx context.Context, finder PrimaryCalendarFinder, userID uuid.UUID, email string) (*CalendarRecord, error) {
	storedCalendar, err := finder.FindPrimaryCalendar(ctx, userID)
	if err != nil {
		log.Printf("failed to get primary calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return nil, internalErrors.NewNotFoundError("カレンダーが見つかりませんでした")
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}

func (uc *Usecase) loadAdjustaCandidateCalendar(ctx context.Context, finder AdjustaCandidateCalendarFinder, userID uuid.UUID, email string) (*CalendarRecord, error) {
	storedCalendar, err := finder.FindAdjustaCandidateCalendar(ctx, userID)
	if err != nil {
		log.Printf("failed to get adjusta candidate calendar for account: %s, error: %v", email, err)
		if repoerr.IsNotFound(err) {
			return &CalendarRecord{SyncProposedDates: false}, nil
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return storedCalendar, nil
}
