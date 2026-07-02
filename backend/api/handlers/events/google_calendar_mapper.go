package events

import (
	"github.com/koo-arch/adjusta-backend/api/dto"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

func toGoogleEventResponse(event *usecaseEvents.FetchedGoogleEvent) *dto.GoogleEvent {
	if event == nil {
		return nil
	}

	return &dto.GoogleEvent{
		ID:          event.ID,
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		ColorID:     event.ColorID,
		Start:       event.Start,
		End:         event.End,
	}
}

func toGoogleEventResponses(events []*usecaseEvents.FetchedGoogleEvent) []*dto.GoogleEvent {
	responses := make([]*dto.GoogleEvent, 0, len(events))
	for _, event := range events {
		responses = append(responses, toGoogleEventResponse(event))
	}
	return responses
}
