package events

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
)

type fakeEventTxStore struct {
	t                              *testing.T
	findPrimaryCalendarFn          func(ctx context.Context, userID uuid.UUID) (*EventCalendar, error)
	findAdjustaCandidateCalendarFn func(ctx context.Context, userID uuid.UUID) (*EventCalendar, error)
	findEventByIDFn                func(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error)
	readCalendarFn                 func(ctx context.Context, calendarID uuid.UUID) (*EventCalendar, error)
	createEventFn                  func(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*domainEvent.Event, error)
	updateEventFn                  func(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error)
	softDeleteEventFn              func(ctx context.Context, id uuid.UUID) error
	listProposedDatesByEventFn     func(ctx context.Context, eventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error)
	createProposedDatesFn          func(ctx context.Context, selectedDates []SelectedDate, eventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error)
	updateProposedDateFn           func(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error)
	deleteProposedDateFn           func(ctx context.Context, id uuid.UUID) error
	createProposedDateFn           func(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*domainProposedDate.ProposedDate, error)
}

func (f *fakeEventTxStore) unexpected(method string) {
	f.t.Fatalf("%s should not be called", method)
}

func (f *fakeEventTxStore) FindPrimaryCalendar(ctx context.Context, userID uuid.UUID) (*EventCalendar, error) {
	if f.findPrimaryCalendarFn == nil {
		f.unexpected("FindPrimaryCalendar")
	}
	return f.findPrimaryCalendarFn(ctx, userID)
}

func (f *fakeEventTxStore) FindAdjustaCandidateCalendar(ctx context.Context, userID uuid.UUID) (*EventCalendar, error) {
	if f.findAdjustaCandidateCalendarFn == nil {
		f.unexpected("FindAdjustaCandidateCalendar")
	}
	return f.findAdjustaCandidateCalendarFn(ctx, userID)
}

func (f *fakeEventTxStore) FindEventByID(ctx context.Context, userID, eventID uuid.UUID, withProposedDates bool) (*domainEvent.Event, error) {
	if f.findEventByIDFn == nil {
		f.unexpected("FindEventByID")
	}
	return f.findEventByIDFn(ctx, userID, eventID, withProposedDates)
}

func (f *fakeEventTxStore) ReadCalendar(ctx context.Context, calendarID uuid.UUID) (*EventCalendar, error) {
	if f.readCalendarFn == nil {
		f.unexpected("ReadCalendar")
	}
	return f.readCalendarFn(ctx, calendarID)
}

func (f *fakeEventTxStore) CreateEvent(ctx context.Context, userID, primaryCalendarID uuid.UUID, title, location, description string, start, end time.Time) (*domainEvent.Event, error) {
	if f.createEventFn == nil {
		f.unexpected("CreateEvent")
	}
	return f.createEventFn(ctx, userID, primaryCalendarID, title, location, description, start, end)
}

func (f *fakeEventTxStore) UpdateEvent(ctx context.Context, id uuid.UUID, opt EventMutation) (*domainEvent.Event, error) {
	if f.updateEventFn == nil {
		f.unexpected("UpdateEvent")
	}
	return f.updateEventFn(ctx, id, opt)
}

func (f *fakeEventTxStore) SoftDeleteEvent(ctx context.Context, id uuid.UUID) error {
	if f.softDeleteEventFn == nil {
		f.unexpected("SoftDeleteEvent")
	}
	return f.softDeleteEventFn(ctx, id)
}

func (f *fakeEventTxStore) ListProposedDatesByEvent(ctx context.Context, eventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
	if f.listProposedDatesByEventFn == nil {
		f.unexpected("ListProposedDatesByEvent")
	}
	return f.listProposedDatesByEventFn(ctx, eventID)
}

func (f *fakeEventTxStore) CreateProposedDates(ctx context.Context, selectedDates []SelectedDate, eventID uuid.UUID) ([]*domainProposedDate.ProposedDate, error) {
	if f.createProposedDatesFn == nil {
		f.unexpected("CreateProposedDates")
	}
	return f.createProposedDatesFn(ctx, selectedDates, eventID)
}

func (f *fakeEventTxStore) UpdateProposedDate(ctx context.Context, id uuid.UUID, opt ProposedDateMutation) (*domainProposedDate.ProposedDate, error) {
	if f.updateProposedDateFn == nil {
		f.unexpected("UpdateProposedDate")
	}
	return f.updateProposedDateFn(ctx, id, opt)
}

func (f *fakeEventTxStore) DeleteProposedDate(ctx context.Context, id uuid.UUID) error {
	if f.deleteProposedDateFn == nil {
		f.unexpected("DeleteProposedDate")
	}
	return f.deleteProposedDateFn(ctx, id)
}

func (f *fakeEventTxStore) CreateProposedDate(ctx context.Context, opt ProposedDateMutation, eventID uuid.UUID) (*domainProposedDate.ProposedDate, error) {
	if f.createProposedDateFn == nil {
		f.unexpected("CreateProposedDate")
	}
	return f.createProposedDateFn(ctx, opt, eventID)
}

type fakeEventTransaction struct {
	store *fakeEventTxStore
}

func (f *fakeEventTransaction) DoEvent(ctx context.Context, fn func(repos EventTxRepositories) error) error {
	return fn(fakeReposFromTxStore(f.store))
}
