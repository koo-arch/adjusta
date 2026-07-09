package events

import (
	"github.com/koo-arch/adjusta-backend/api/dto"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

func toDraftCreationRequest(eventDraft *dto.EventDraftCreation) usecaseEvents.DraftCreationRequest {
	selectedDates := make([]usecaseEvents.SelectedDate, 0, len(eventDraft.SelectedDates))
	for _, date := range eventDraft.SelectedDates {
		selectedDates = append(selectedDates, usecaseEvents.SelectedDate{
			Start:    date.Start,
			End:      date.End,
			Priority: date.Priority,
		})
	}

	return usecaseEvents.DraftCreationRequest{
		Title:         eventDraft.Title,
		Location:      eventDraft.Location,
		Description:   eventDraft.Description,
		SelectedDates: selectedDates,
	}
}

func toDraftUpdateRequest(eventDraft *dto.EventDraftUpdate) usecaseEvents.DraftUpdateRequest {
	proposedDates := make([]usecaseEvents.ProposedDateRequest, 0, len(eventDraft.ProposedDates))
	for _, date := range eventDraft.ProposedDates {
		proposedDates = append(proposedDates, usecaseEvents.ProposedDateRequest{
			ID:            date.ID,
			GoogleEventID: date.GoogleEventID,
			Start:         date.Start,
			End:           date.End,
			Priority:      date.Priority,
		})
	}

	return usecaseEvents.DraftUpdateRequest{
		Title:         eventDraft.Title,
		Location:      eventDraft.Location,
		Description:   eventDraft.Description,
		Status:        eventDraft.Status,
		ProposedDates: proposedDates,
	}
}

func toProposedDateResponse(date usecaseEvents.ProposedDateOutput) dto.ProposedDate {
	return dto.ProposedDate{
		ID:            date.ID,
		GoogleEventID: date.GoogleEventID,
		Start:         date.Start,
		End:           date.End,
		Priority:      date.Priority,
		Status:        date.Status,
		SyncStatus:    date.SyncStatus,
		LastSyncedAt:  date.LastSyncedAt,
		LastSyncError: date.LastSyncError,
	}
}

func toProposedDateResponses(dates []usecaseEvents.ProposedDateOutput) []dto.ProposedDate {
	responses := make([]dto.ProposedDate, 0, len(dates))
	for _, date := range dates {
		responses = append(responses, toProposedDateResponse(date))
	}
	return responses
}

func toEventDraftDetailResponse(event *usecaseEvents.EventDraftDetailOutput) *dto.EventDraftDetail {
	if event == nil {
		return nil
	}

	return &dto.EventDraftDetail{
		ID:                     event.ID,
		Title:                  event.Title,
		Location:               event.Location,
		Description:            event.Description,
		Status:                 event.Status,
		SyncStatus:             event.SyncStatus,
		ConfirmedDateID:        event.ConfirmedDateID,
		GoogleEventID:          event.GoogleEventID,
		ConfirmedGoogleEventID: event.ConfirmedGoogleEventID,
		LastSyncedAt:           event.LastSyncedAt,
		LastSyncError:          event.LastSyncError,
		ProposedDates:          toProposedDateResponses(event.ProposedDates),
	}
}

func toEventDraftDetailResponses(events []*usecaseEvents.EventDraftDetailOutput) []*dto.EventDraftDetail {
	responses := make([]*dto.EventDraftDetail, 0, len(events))
	for _, event := range events {
		responses = append(responses, toEventDraftDetailResponse(event))
	}
	return responses
}

func toEventDraftListResponse(result *usecaseEvents.EventDraftListOutput) *dto.EventDraftList {
	if result == nil {
		return &dto.EventDraftList{
			Items: []*dto.EventDraftDetail{},
		}
	}

	return &dto.EventDraftList{
		Items: toEventDraftDetailResponses(result.Items),
		Pagination: dto.Pagination{
			Page:       result.Pagination.Page,
			PerPage:    result.Pagination.PerPage,
			TotalItems: result.Pagination.TotalItems,
			TotalPages: result.Pagination.TotalPages,
		},
	}
}

func toUpcomingEventResponse(event usecaseEvents.UpcomingEventOutput) dto.UpcomingEvent {
	return dto.UpcomingEvent{
		ID:                     event.ID,
		Title:                  event.Title,
		Location:               event.Location,
		Description:            event.Description,
		Status:                 event.Status,
		SyncStatus:             event.SyncStatus,
		ConfirmedDateID:        event.ConfirmedDateID,
		GoogleEventID:          event.GoogleEventID,
		ConfirmedGoogleEventID: event.ConfirmedGoogleEventID,
		LastSyncedAt:           event.LastSyncedAt,
		LastSyncError:          event.LastSyncError,
		Start:                  event.Start,
		End:                    event.End,
	}
}

func toUpcomingEventResponses(events []usecaseEvents.UpcomingEventOutput) []dto.UpcomingEvent {
	responses := make([]dto.UpcomingEvent, 0, len(events))
	for _, event := range events {
		responses = append(responses, toUpcomingEventResponse(event))
	}
	return responses
}

func toNeedsActionDraftResponse(event usecaseEvents.NeedsActionDraftOutput) dto.NeedsActionDraft {
	return dto.NeedsActionDraft{
		ID:             event.ID,
		Title:          event.Title,
		Location:       event.Location,
		Description:    event.Description,
		Status:         event.Status,
		Start:          event.Start,
		End:            event.End,
		NeedsAttention: event.NeedsAttention,
	}
}

func toNeedsActionDraftResponses(events []usecaseEvents.NeedsActionDraftOutput) []dto.NeedsActionDraft {
	responses := make([]dto.NeedsActionDraft, 0, len(events))
	for _, event := range events {
		responses = append(responses, toNeedsActionDraftResponse(event))
	}
	return responses
}
