package events

import (
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func toDomainDraftDate(date ProposedDateRequest) (domainEvent.DraftProposedDate, error) {
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

func toDomainDraftDateList(dates []ProposedDateRequest) ([]domainEvent.DraftProposedDate, error) {
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

func assignSelectedDatePriorities(dates []SelectedDate) []SelectedDate {
	draftDates := make([]domainEvent.DraftProposedDate, 0, len(dates))
	for _, date := range dates {
		draftDates = append(draftDates, domainEvent.DraftProposedDate{
			Start: date.Start,
			End:   date.End,
		})
	}

	assignedDraftDates := domainEvent.AssignPrioritiesByOrder(draftDates)
	assigned := make([]SelectedDate, 0, len(dates))
	for _, date := range assignedDraftDates {
		assigned = append(assigned, SelectedDate{
			Start:    date.Start,
			End:      date.End,
			Priority: date.Priority,
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
