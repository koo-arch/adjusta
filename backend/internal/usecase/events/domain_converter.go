package events

import (
	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func toDomainDraftDate(date appmodel.ProposedDate) (domainEvent.DraftProposedDate, error) {
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

func toDomainDraftDateList(dates []appmodel.ProposedDate) ([]domainEvent.DraftProposedDate, error) {
	converted := make([]domainEvent.DraftProposedDate, 0, len(dates))
	for _, date := range dates {
		convertedDate, err := toDomainDraftDate(date)
		if err != nil {
			return nil, err
		}
		converted = append(converted, convertedDate)
	}
	return converted, nil
}

func toConfirmationRequest(date appmodel.ConfirmDate) ConfirmationRequest {
	return ConfirmationRequest{
		ID:            date.ID,
		GoogleEventID: date.GoogleEventID,
		Start:         date.Start,
		End:           date.End,
		Priority:      date.Priority,
	}
}

func toDomainConfirmationRequest(date ConfirmationRequest) (domainEvent.DraftProposedDate, error) {
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

func assignSelectedDatePriorities(dates []appmodel.SelectedDate) []SelectedDate {
	assigned := make([]SelectedDate, 0, len(dates))
	for i, date := range dates {
		assigned = append(assigned, SelectedDate{
			Start:    date.Start,
			End:      date.End,
			Priority: domainEvent.PriorityForOrder(i, len(dates)),
		})
	}
	return assigned
}

func toDomainExistingDateList(dates []*ProposedDateRecord) []domainEvent.ExistingProposedDate {
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
