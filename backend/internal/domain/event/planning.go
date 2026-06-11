package event

import (
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/internal/domainvalue"
)

var ErrInvalidDateRange = errors.New("invalid date range")

type DraftProposedDate struct {
	ID       *uuid.UUID
	Start    time.Time
	End      time.Time
	Priority int
}

const PriorityStep = 1024

type ExistingProposedDate struct {
	ID       uuid.UUID
	Start    time.Time
	End      time.Time
	Priority int
}

type ProposedDateCreate struct {
	Start    time.Time
	End      time.Time
	Priority int
}

type ProposedDateUpdate struct {
	ID       uuid.UUID
	Start    time.Time
	End      time.Time
	Priority int
}

type ProposedDateChangeSet struct {
	Creates []ProposedDateCreate
	Updates []ProposedDateUpdate
	Deletes []uuid.UUID
}

type ConfirmationPlan struct {
	Status domainvalue.EventStatus
	Create *ProposedDateCreate
	Update *ProposedDateUpdate
}

func BuildProposedDateChangeSet(requested []DraftProposedDate, existing []ExistingProposedDate) ProposedDateChangeSet {
	changeSet := ProposedDateChangeSet{
		Creates: make([]ProposedDateCreate, 0),
		Updates: make([]ProposedDateUpdate, 0),
		Deletes: make([]uuid.UUID, 0),
	}

	existingByID := make(map[uuid.UUID]ExistingProposedDate, len(existing))
	for _, date := range existing {
		existingByID[date.ID] = date
	}

	requested = assignRequestedPriorities(requested, existing)
	requestedByID := make(map[uuid.UUID]DraftProposedDate, len(requested))
	for _, date := range requested {
		if date.ID == nil {
			changeSet.Creates = append(changeSet.Creates, ProposedDateCreate{
				Start:    date.Start,
				End:      date.End,
				Priority: date.Priority,
			})
			continue
		}

		requestedByID[*date.ID] = date
	}

	for _, date := range existing {
		requestedDate, ok := requestedByID[date.ID]
		if !ok {
			changeSet.Deletes = append(changeSet.Deletes, date.ID)
			continue
		}

		if date.Start.Equal(requestedDate.Start) && date.End.Equal(requestedDate.End) && date.Priority == requestedDate.Priority {
			continue
		}

		changeSet.Updates = append(changeSet.Updates, ProposedDateUpdate{
			ID:       date.ID,
			Start:    requestedDate.Start,
			End:      requestedDate.End,
			Priority: requestedDate.Priority,
		})
	}

	// If an ID was requested but no longer exists in DB, treat it as a create fallback.
	for id, date := range requestedByID {
		if _, ok := existingByID[id]; ok {
			continue
		}
		changeSet.Creates = append(changeSet.Creates, ProposedDateCreate{
			Start:    date.Start,
			End:      date.End,
			Priority: date.Priority,
		})
	}

	return changeSet
}

func NormalizeDraftProposedDatesByOrder(requested []DraftProposedDate) []DraftProposedDate {
	normalized := make([]DraftProposedDate, 0, len(requested))
	for i, date := range requested {
		date.Priority = PriorityValueForOrder(i, len(requested))
		normalized = append(normalized, date)
	}
	return normalized
}

func PriorityValueForOrder(index, total int) int {
	if total <= 0 {
		return 0
	}
	return (total - index) * PriorityStep
}

func BuildConfirmationPlan(confirmDate DraftProposedDate, existing []ExistingProposedDate) (*ConfirmationPlan, error) {
	if !confirmDate.Start.Before(confirmDate.End) {
		return nil, ErrInvalidDateRange
	}

	plan := &ConfirmationPlan{
		Status: domainvalue.StatusConfirmed,
	}
	highestPriority := nextHighestPriority(existing, confirmDate.ID)

	if confirmDate.ID == nil {
		plan.Create = &ProposedDateCreate{
			Start:    confirmDate.Start,
			End:      confirmDate.End,
			Priority: highestPriority,
		}
		return plan, nil
	}

	plan.Update = &ProposedDateUpdate{
		ID:       *confirmDate.ID,
		Start:    confirmDate.Start,
		End:      confirmDate.End,
		Priority: highestPriority,
	}

	return plan, nil
}

func nextHighestPriority(existing []ExistingProposedDate, excludeID *uuid.UUID) int {
	highest := 0
	for _, date := range existing {
		if excludeID != nil && date.ID == *excludeID {
			continue
		}
		if date.Priority > highest {
			highest = date.Priority
		}
	}
	return highest + PriorityStep
}

func assignRequestedPriorities(requested []DraftProposedDate, existing []ExistingProposedDate) []DraftProposedDate {
	if len(requested) == 0 {
		return nil
	}
	if len(existing) == 0 {
		return NormalizeDraftProposedDatesByOrder(requested)
	}

	existingByID := make(map[uuid.UUID]ExistingProposedDate, len(existing))
	currentOrderIDs := make([]uuid.UUID, 0, len(existing))
	currentOrder := slices.Clone(existing)
	slices.SortFunc(currentOrder, func(a, b ExistingProposedDate) int {
		switch {
		case a.Priority > b.Priority:
			return -1
		case a.Priority < b.Priority:
			return 1
		default:
			return 0
		}
	})
	for _, date := range currentOrder {
		existingByID[date.ID] = date
		currentOrderIDs = append(currentOrderIDs, date.ID)
	}

	requestedExistingIDs := make([]uuid.UUID, 0, len(requested))
	for _, date := range requested {
		if date.ID == nil {
			continue
		}
		if _, ok := existingByID[*date.ID]; ok {
			requestedExistingIDs = append(requestedExistingIDs, *date.ID)
		}
	}

	anchorIDs := longestCommonSubsequence(currentOrderIDs, requestedExistingIDs)
	if len(anchorIDs) == 0 {
		return NormalizeDraftProposedDatesByOrder(requested)
	}

	anchorSet := make(map[uuid.UUID]struct{}, len(anchorIDs))
	for _, id := range anchorIDs {
		anchorSet[id] = struct{}{}
	}

	assigned := make([]DraftProposedDate, 0, len(requested))
	buffer := make([]DraftProposedDate, 0)
	var previousAnchorPriority *int

	flushBuffer := func(nextAnchorPriority *int) bool {
		if len(buffer) == 0 {
			return true
		}

		priorities, ok := distributePriorities(previousAnchorPriority, nextAnchorPriority, len(buffer))
		if !ok {
			return false
		}

		for i, date := range buffer {
			date.Priority = priorities[i]
			assigned = append(assigned, date)
		}
		buffer = buffer[:0]
		return true
	}

	for _, date := range requested {
		if date.ID == nil {
			buffer = append(buffer, date)
			continue
		}

		existingDate, ok := existingByID[*date.ID]
		if !ok {
			buffer = append(buffer, date)
			continue
		}

		if _, ok := anchorSet[existingDate.ID]; !ok {
			buffer = append(buffer, date)
			continue
		}

		nextAnchorPriority := existingDate.Priority
		if !flushBuffer(&nextAnchorPriority) {
			return NormalizeDraftProposedDatesByOrder(requested)
		}

		date.Priority = existingDate.Priority
		assigned = append(assigned, date)
		previousAnchorPriority = &date.Priority
	}

	if !flushBuffer(nil) {
		return NormalizeDraftProposedDatesByOrder(requested)
	}

	return assigned
}

func distributePriorities(higher, lower *int, count int) ([]int, bool) {
	if count == 0 {
		return nil, true
	}

	switch {
	case higher == nil && lower == nil:
		priorities := make([]int, 0, count)
		for i := 0; i < count; i++ {
			priorities = append(priorities, PriorityValueForOrder(i, count))
		}
		return priorities, true
	case higher == nil:
		priorities := make([]int, 0, count)
		for i := 0; i < count; i++ {
			priorities = append(priorities, *lower+PriorityStep*(count-i))
		}
		return priorities, true
	case lower == nil:
		if *higher <= count {
			return nil, false
		}
		priorities := make([]int, 0, count)
		for i := 1; i <= count; i++ {
			priority := (*higher * (count + 1 - i)) / (count + 1)
			if priority <= 0 {
				return nil, false
			}
			if len(priorities) > 0 && priorities[len(priorities)-1] <= priority {
				return nil, false
			}
			priorities = append(priorities, priority)
		}
		return priorities, true
	default:
		gap := *higher - *lower
		if gap <= count {
			return nil, false
		}
		priorities := make([]int, 0, count)
		for i := 1; i <= count; i++ {
			priority := *lower + (gap*(count+1-i))/(count+1)
			if priority >= *higher || priority <= *lower {
				return nil, false
			}
			if len(priorities) > 0 && priorities[len(priorities)-1] <= priority {
				return nil, false
			}
			priorities = append(priorities, priority)
		}
		return priorities, true
	}
}

func longestCommonSubsequence(currentOrderIDs, requestedExistingIDs []uuid.UUID) []uuid.UUID {
	if len(currentOrderIDs) == 0 || len(requestedExistingIDs) == 0 {
		return nil
	}

	dp := make([][]int, len(currentOrderIDs)+1)
	for i := range dp {
		dp[i] = make([]int, len(requestedExistingIDs)+1)
	}

	for i := len(currentOrderIDs) - 1; i >= 0; i-- {
		for j := len(requestedExistingIDs) - 1; j >= 0; j-- {
			if currentOrderIDs[i] == requestedExistingIDs[j] {
				dp[i][j] = dp[i+1][j+1] + 1
				continue
			}
			if dp[i+1][j] >= dp[i][j+1] {
				dp[i][j] = dp[i+1][j]
			} else {
				dp[i][j] = dp[i][j+1]
			}
		}
	}

	result := make([]uuid.UUID, 0, dp[0][0])
	i, j := 0, 0
	for i < len(currentOrderIDs) && j < len(requestedExistingIDs) {
		if currentOrderIDs[i] == requestedExistingIDs[j] {
			result = append(result, currentOrderIDs[i])
			i++
			j++
			continue
		}
		if dp[i+1][j] > dp[i][j+1] {
			i++
		} else {
			j++
		}
	}

	return result
}
