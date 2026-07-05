package event

import (
	"slices"

	"github.com/google/uuid"
)

const PriorityStep = 1024

func AssignPrioritiesByOrder(requested []DraftProposedDate) []DraftProposedDate {
	assigned := make([]DraftProposedDate, 0, len(requested))
	for i, date := range requested {
		date.Priority = PriorityForOrder(i, len(requested))
		assigned = append(assigned, date)
	}
	return assigned
}

func PriorityForOrder(index, total int) int {
	if total <= 0 {
		return 0
	}
	return (total - index) * PriorityStep
}

func reconcileRequestedPriorities(requested []DraftProposedDate, existing []ExistingProposedDate) []DraftProposedDate {
	if len(requested) == 0 {
		return nil
	}
	if len(existing) == 0 {
		return AssignPrioritiesByOrder(requested)
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

	anchorIDs := findStableOrderAnchors(currentOrderIDs, requestedExistingIDs)
	if len(anchorIDs) == 0 {
		return AssignPrioritiesByOrder(requested)
	}

	anchorSet := make(map[uuid.UUID]struct{}, len(anchorIDs))
	for _, id := range anchorIDs {
		anchorSet[id] = struct{}{}
	}

	assigned := make([]DraftProposedDate, 0, len(requested))
	buffer := make([]DraftProposedDate, 0)
	var previousAnchorPriority *int

	flushBufferedDates := func(nextAnchorPriority *int) bool {
		if len(buffer) == 0 {
			return true
		}

		priorities, ok := distributeSparsePriorities(previousAnchorPriority, nextAnchorPriority, len(buffer))
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
		if !flushBufferedDates(&nextAnchorPriority) {
			return AssignPrioritiesByOrder(requested)
		}

		date.Priority = existingDate.Priority
		assigned = append(assigned, date)
		previousAnchorPriority = &date.Priority
	}

	if !flushBufferedDates(nil) {
		return AssignPrioritiesByOrder(requested)
	}

	return assigned
}

func nextHighestPriorityExcluding(existing []ExistingProposedDate, excludeID *uuid.UUID) int {
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

func distributeSparsePriorities(higher, lower *int, count int) ([]int, bool) {
	if count == 0 {
		return nil, true
	}

	switch {
	case higher == nil && lower == nil:
		priorities := make([]int, 0, count)
		for i := 0; i < count; i++ {
			priorities = append(priorities, PriorityForOrder(i, count))
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

func findStableOrderAnchors(currentOrderIDs, requestedExistingIDs []uuid.UUID) []uuid.UUID {
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
