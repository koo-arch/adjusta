package events

import (
	repoEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	repoProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
)

type EventMutation = repoEvent.EventUpdateOptions
type ProposedDateMutation = repoProposedDate.ProposedDateUpdateOptions
