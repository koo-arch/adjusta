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
}

func toEventSearchOptions(opt EventSearchOptions) repoEvent.EventSearchOptions {
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
	}
}
