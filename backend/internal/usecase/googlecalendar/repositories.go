package googlecalendar

import (
	"context"
	"time"

	"github.com/google/uuid"
	domainCalendar "github.com/koo-arch/adjusta-backend/internal/domain/calendar"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainOutbox "github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	domainUserCalendar "github.com/koo-arch/adjusta-backend/internal/domain/usercalendar"
)

type CalendarRepository interface {
	Read(ctx context.Context, id uuid.UUID) (*domainCalendar.Calendar, error)
}

type EventRepository interface {
	Read(ctx context.Context, id uuid.UUID, opt domainEvent.EventReadOptions) (*domainEvent.Event, error)
	Update(ctx context.Context, id uuid.UUID, opt domainEvent.EventUpdateOptions) (*domainEvent.Event, error)
}

type MessageRepository interface {
	Read(ctx context.Context, id uuid.UUID) (*domainOutbox.Message, error)
	MarkProcessed(ctx context.Context, id uuid.UUID, processedAt time.Time) error
}

type ProposedDateRepository interface {
	Update(ctx context.Context, id uuid.UUID, opt domainProposedDate.ProposedDateUpdateOptions) (*domainProposedDate.ProposedDate, error)
}

type UserCalendarRepository interface {
	FilterByUserID(ctx context.Context, userID uuid.UUID) ([]*domainUserCalendar.UserCalendar, error)
}

type Repositories struct {
	Calendar     CalendarRepository
	Event        EventRepository
	Message      MessageRepository
	ProposedDate ProposedDateRepository
	UserCalendar UserCalendarRepository
}
