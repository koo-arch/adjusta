package events

import (
	"context"
	"log"
	"sort"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
	repositorymodel "github.com/koo-arch/adjusta-backend/internal/repositorymodel"
)

func (uc *Usecase) findPrimaryCalendar(ctx context.Context, finder PrimaryCalendarFinder, userID uuid.UUID, email string) (*repositorymodel.StoredCalendar, error) {
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

func buildProposedDates(storedDates []*repositorymodel.StoredProposedDate) []appmodel.ProposedDate {
	proposedDates := make([]appmodel.ProposedDate, 0, len(storedDates))
	for _, storedDate := range storedDates {
		proposedDates = append(proposedDates, appmodel.ProposedDate{
			ID:       &storedDate.ID,
			Start:    &storedDate.StartTime,
			End:      &storedDate.EndTime,
			Priority: storedDate.Priority,
		})
	}

	sort.Slice(proposedDates, func(i, j int) bool {
		return proposedDates[i].Priority < proposedDates[j].Priority
	})

	return proposedDates
}

func buildEventDraftDetail(storedEvent *repositorymodel.StoredEvent) (*appmodel.EventDraftDetail, error) {
	if storedEvent.ProposedDates == nil {
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &appmodel.EventDraftDetail{
		ID:              storedEvent.ID,
		Title:           storedEvent.Summary,
		Location:        storedEvent.Location,
		Description:     storedEvent.Description,
		Status:          storedEvent.Status,
		ConfirmedDateID: &storedEvent.ConfirmedDateID,
		GoogleEventID:   storedEvent.GoogleEventID,
		Slug:            storedEvent.Slug,
		ProposedDates:   buildProposedDates(storedEvent.ProposedDates),
	}, nil
}

func normalizeUsecaseError(err error, fallbackMessage string) error {
	return internalErrors.NormalizeAPIError(err, fallbackMessage)
}

func toDomainDraftProposedDate(date appmodel.ProposedDate) (domainEvent.DraftProposedDate, error) {
	if date.Start == nil || date.End == nil {
		return domainEvent.DraftProposedDate{}, internalErrors.NewBadRequestError("候補日程が不正です")
	}

	return domainEvent.DraftProposedDate{
		ID:       date.ID,
		Start:    *date.Start,
		End:      *date.End,
		Priority: date.Priority,
	}, nil
}

func toDomainDraftProposedDates(dates []appmodel.ProposedDate) ([]domainEvent.DraftProposedDate, error) {
	converted := make([]domainEvent.DraftProposedDate, 0, len(dates))
	for _, date := range dates {
		convertedDate, err := toDomainDraftProposedDate(date)
		if err != nil {
			return nil, err
		}
		converted = append(converted, convertedDate)
	}
	return converted, nil
}

func toDomainConfirmationDate(date appmodel.ConfirmDate) (domainEvent.DraftProposedDate, error) {
	if date.Start == nil || date.End == nil {
		return domainEvent.DraftProposedDate{}, internalErrors.NewBadRequestError("確定候補日程が不正です")
	}

	return domainEvent.DraftProposedDate{
		ID:       date.ID,
		Start:    *date.Start,
		End:      *date.End,
		Priority: date.Priority,
	}, nil
}

func toDomainExistingProposedDates(dates []*repositorymodel.StoredProposedDate) []domainEvent.ExistingProposedDate {
	converted := make([]domainEvent.ExistingProposedDate, 0, len(dates))
	for _, date := range dates {
		converted = append(converted, domainEvent.ExistingProposedDate{
			ID:       date.ID,
			Start:    date.StartTime,
			End:      date.EndTime,
			Priority: date.Priority,
		})
	}
	return converted
}
