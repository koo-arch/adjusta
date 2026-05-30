package queryparser

import (
	"fmt"

	"github.com/koo-arch/adjusta-backend/internal/models"
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
	eventStatus, err := qp.vaildateStatus(status)
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

	options := usecaseEvents.SearchDraftQuery{
		Title:        title,
		Location:     location,
		Description:  description,
		Status:       eventStatus,
		StartTimeGTE: startTimeGTE,
		StartTimeLTE: startTimeLTE,
		EndTimeGTE:   endTimeGTE,
		EndTimeLTE:   endTimeLTE,
	}

	return &options, nil
}

func (qp *QueryParser) vaildateStatus(status *string) (*models.EventStatus, error) {
	if status == nil {
		return nil, nil
	}

	var result models.EventStatus

	switch *status {
	case "pending":
		result = models.StatusPending
	case "confirmed":
		result = models.StatusConfirmed
	case "cancelled":
		result = models.StatusCancelled
	default:
		return nil, fmt.Errorf("invalid status: %s", *status)
	}

	return &result, nil
}
