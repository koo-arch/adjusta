package googlecalendar

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	domainCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainOutbox "github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type fakeCalendarRepository struct {
	calendar *domainCalendar.Calendar
}

func (f *fakeCalendarRepository) Read(context.Context, uuid.UUID) (*domainCalendar.Calendar, error) {
	return f.calendar, nil
}

type fakeEventRepository struct {
	event  *domainEvent.Event
	update domainEvent.EventUpdateOptions
}

func (f *fakeEventRepository) Read(context.Context, uuid.UUID, domainEvent.EventReadOptions) (*domainEvent.Event, error) {
	return f.event, nil
}

func (f *fakeEventRepository) Update(_ context.Context, _ uuid.UUID, update domainEvent.EventUpdateOptions) (*domainEvent.Event, error) {
	f.update = update
	return f.event, nil
}

type fakeMessageRepository struct {
	message   *domainOutbox.Message
	processed bool
}

func (f *fakeMessageRepository) Read(context.Context, uuid.UUID) (*domainOutbox.Message, error) {
	return f.message, nil
}

func (f *fakeMessageRepository) MarkProcessed(context.Context, uuid.UUID, time.Time) error {
	f.processed = true
	return nil
}

type fakeProposedDateRepository struct {
	date   *domainProposedDate.ProposedDate
	update domainProposedDate.ProposedDateUpdateOptions
}

func (f *fakeProposedDateRepository) Update(_ context.Context, _ uuid.UUID, update domainProposedDate.ProposedDateUpdateOptions) (*domainProposedDate.ProposedDate, error) {
	f.update = update
	return f.date, nil
}

type fakeUserCalendarRepository struct {
	userCalendar *domainUserCalendar.UserCalendar
}

func (f *fakeUserCalendarRepository) FilterByUserID(context.Context, uuid.UUID) ([]*domainUserCalendar.UserCalendar, error) {
	return []*domainUserCalendar.UserCalendar{f.userCalendar}, nil
}

type fakeTransaction struct {
	repos Repositories
}

func (f *fakeTransaction) Do(ctx context.Context, fn func(repos Repositories) error) error {
	return fn(f.repos)
}

type fakeEventGateway struct {
	upsert func(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error)
}

func (f *fakeEventGateway) UpsertEvent(ctx context.Context, userID uuid.UUID, calendarID string, existingGoogleEventID *string, title, location, description string, start, end time.Time) (string, error) {
	return f.upsert(ctx, userID, calendarID, existingGoogleEventID, title, location, description, start, end)
}

func (f *fakeEventGateway) DeleteEvent(context.Context, uuid.UUID, string, string) error {
	return errors.New("DeleteEvent should not be called")
}

func TestProcessMessageSynchronizesDraftAndMarksMessageProcessed(t *testing.T) {
	userID := uuid.New()
	eventID := uuid.New()
	messageID := uuid.New()
	dateID := uuid.New()
	calendarID := uuid.New()
	start := time.Now().UTC().Add(time.Hour)
	date := &domainProposedDate.ProposedDate{
		ID:         dateID,
		EventID:    eventID,
		StartTime:  start,
		EndTime:    start.Add(time.Hour),
		Priority:   1024,
		SyncStatus: value.SyncStatusPending,
	}
	event := &domainEvent.Event{
		ID:            eventID,
		UserID:        userID,
		Title:         "Planning",
		Status:        value.StatusActive,
		SyncStatus:    value.SyncStatusPending,
		ProposedDates: []*domainProposedDate.ProposedDate{date},
	}
	messageRepo := &fakeMessageRepository{message: &domainOutbox.Message{
		ID:            messageID,
		MessageType:   domainOutbox.MessageTypeCalendarEventSync,
		AggregateType: domainOutbox.AggregateTypeEvent,
		AggregateID:   eventID,
	}}
	eventRepo := &fakeEventRepository{event: event}
	dateRepo := &fakeProposedDateRepository{date: date}
	repos := Repositories{
		Calendar: &fakeCalendarRepository{calendar: &domainCalendar.Calendar{
			ID:               calendarID,
			GoogleCalendarID: "candidate-calendar",
		}},
		Event:        eventRepo,
		Message:      messageRepo,
		ProposedDate: dateRepo,
		UserCalendar: &fakeUserCalendarRepository{userCalendar: &domainUserCalendar.UserCalendar{
			UserID:            userID,
			CalendarID:        calendarID,
			Role:              value.UserCalendarRoleAdjustaCandidate,
			SyncProposedDates: true,
		}},
	}
	usecase := NewSyncUsecase(repos, &fakeTransaction{repos: repos}, &fakeEventGateway{
		upsert: func(_ context.Context, gotUserID uuid.UUID, gotCalendarID string, existingID *string, title, _, _ string, _, _ time.Time) (string, error) {
			if gotUserID != userID || gotCalendarID != "candidate-calendar" || title != "Planning【第1候補】" {
				t.Fatalf("unexpected event payload: %s %s %s", gotUserID, gotCalendarID, title)
			}
			expectedID := stableGoogleEventID("adjustap", dateID)
			if existingID == nil || *existingID != expectedID {
				t.Fatalf("unexpected stable event id: %#v", existingID)
			}
			return "google-event-id", nil
		},
	})

	if err := usecase.ProcessMessage(context.Background(), messageID); err != nil {
		t.Fatalf("ProcessMessage returned error: %v", err)
	}
	if !messageRepo.processed {
		t.Fatal("expected message to be marked processed")
	}
	if eventRepo.update.SyncStatus == nil || *eventRepo.update.SyncStatus != value.SyncStatusSynced {
		t.Fatalf("unexpected event sync update: %#v", eventRepo.update.SyncStatus)
	}
	if dateRepo.update.GoogleEventID == nil || *dateRepo.update.GoogleEventID != "google-event-id" {
		t.Fatalf("unexpected proposed date update: %#v", dateRepo.update.GoogleEventID)
	}
}
