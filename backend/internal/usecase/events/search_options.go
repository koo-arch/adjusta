package events

import (
	"time"

	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
)

type SearchDraftQuery struct {
	Title        *string
	Location     *string
	Description  *string
	Status       *value.EventStatus
	StartTimeGTE *time.Time
	StartTimeLTE *time.Time
	EndTimeGTE   *time.Time
	EndTimeLTE   *time.Time
	SortBy       string
	SortOrder    string
	Page         int
	PerPage      int
}

type EventSearchOptions struct {
	WithProposedDates bool
	Title             *string
	Location          *string
	Description       *string
	Status            *value.EventStatus
	StartTimeGTE      *time.Time
	StartTimeLTE      *time.Time
	EndTimeGTE        *time.Time
	EndTimeLTE        *time.Time
	SortBy            string
	SortOrder         string
	Page              int
	PerPage           int
}

func toEventSearchOptions(opt EventSearchOptions) repoEvent.EventSearchOptions {
	eventOffset := 0
	eventLimit := 0
	if opt.PerPage > 0 {
		page := opt.Page
		if page < 1 {
			page = 1
		}
		eventOffset = (page - 1) * opt.PerPage
		eventLimit = opt.PerPage
	}

	return repoEvent.EventSearchOptions{
		Title:                opt.Title,
		Location:             opt.Location,
		Description:          opt.Description,
		Status:               opt.Status,
		WithProposedDates:    opt.WithProposedDates,
		ProposedDateStartGTE: opt.StartTimeGTE,
		ProposedDateStartLTE: opt.StartTimeLTE,
		ProposedDateEndGTE:   opt.EndTimeGTE,
		ProposedDateEndLTE:   opt.EndTimeLTE,
		SortBy:               opt.SortBy,
		SortOrder:            opt.SortOrder,
		EventOffset:          eventOffset,
		EventLimit:           eventLimit,
	}
}
