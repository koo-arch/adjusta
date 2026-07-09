package queryparser

import (
	"fmt"

	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

func (qp *QueryParser) ParseSearchEventQuery() (*usecaseEvents.SearchDraftQuery, error) {

	title, err := qp.ParseString("title")
	if err != nil {
		return nil, fmt.Errorf("failed to parse title: %w", err)
	}

	location, err := qp.ParseString("location")
	if err != nil {
		return nil, fmt.Errorf("failed to parse location: %w", err)
	}

	description, err := qp.ParseString("description")
	if err != nil {
		return nil, fmt.Errorf("failed to parse description: %w", err)
	}

	status, err := qp.ParseString("status")
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}
	eventStatus, err := qp.validateStatus(status)
	if err != nil {
		return nil, fmt.Errorf("failed to validate status: %w", err)
	}

	startTimeGTE, err := qp.ParseTime("start_time_gte")
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_time: %w", err)
	}

	startTimeLTE, err := qp.ParseTime("start_time_lte")
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_time: %w", err)
	}

	endTimeGTE, err := qp.ParseTime("end_time_gte")
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_time: %w", err)
	}

	endTimeLTE, err := qp.ParseTime("end_time_lte")
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_time: %w", err)
	}
	page, perPage, err := qp.ParsePagination()
	if err != nil {
		return nil, fmt.Errorf("failed to parse pagination: %w", err)
	}
	sortBy, sortOrder, err := qp.ParseEventSort()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sort: %w", err)
	}

	options := usecaseEvents.SearchDraftQuery{
		Title:        title,
		Location:     location,
		Description:  description,
		Status:       eventStatus,
		StartTimeGTE: startTimeGTE,
		StartTimeLTE: startTimeLTE,
		EndTimeGTE:   endTimeGTE,
		EndTimeLTE:   endTimeLTE,
		SortBy:       sortBy,
		SortOrder:    sortOrder,
		Page:         page,
		PerPage:      perPage,
	}

	return &options, nil
}

func (qp *QueryParser) ParseEventListQuery() (*usecaseEvents.SearchDraftQuery, error) {
	page, perPage, err := qp.ParsePagination()
	if err != nil {
		return nil, fmt.Errorf("failed to parse pagination: %w", err)
	}
	sortBy, sortOrder, err := qp.ParseEventSort()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sort: %w", err)
	}

	return &usecaseEvents.SearchDraftQuery{
		SortBy:    sortBy,
		SortOrder: sortOrder,
		Page:      page,
		PerPage:   perPage,
	}, nil
}

func (qp *QueryParser) ParseEventSort() (string, string, error) {
	sortBy, err := qp.ParseDefaultString("sort_by", "created_at")
	if err != nil {
		return "", "", err
	}
	sortOrder, err := qp.ParseDefaultString("sort_order", "desc")
	if err != nil {
		return "", "", err
	}

	switch *sortBy {
	case "created_at", "updated_at", "title", "status":
	default:
		return "", "", fmt.Errorf("invalid sort_by: %s", *sortBy)
	}

	switch *sortOrder {
	case "asc", "desc":
	default:
		return "", "", fmt.Errorf("invalid sort_order: %s", *sortOrder)
	}

	return *sortBy, *sortOrder, nil
}

func (qp *QueryParser) validateStatus(status *string) (*value.EventStatus, error) {
	if status == nil {
		return nil, nil
	}

	var result value.EventStatus

	switch *status {
	case "draft":
		result = value.StatusDraft
	case "active", "pending":
		result = value.StatusActive
	case "confirmed":
		result = value.StatusConfirmed
	case "cancelled":
		result = value.StatusCancelled
	default:
		return nil, fmt.Errorf("invalid status: %s", *status)
	}

	return &result, nil
}
