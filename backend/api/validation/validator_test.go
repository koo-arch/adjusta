package validation

import (
	"testing"

	"github.com/koo-arch/adjusta-backend/internal/appmodel"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func TestFinalizeValidationReturnsValidationAPIErrorForMissingDates(t *testing.T) {
	t.Parallel()

	err := FinalizeValidation(&appmodel.ConfirmEvent{
		ConfirmDate: appmodel.ConfirmDate{},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	apiErr, ok := err.(*internalErrors.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Kind != internalErrors.KindValidation {
		t.Fatalf("expected validation kind, got %s", apiErr.Kind)
	}
	if got := apiErr.Details["confirm_date.start"]; len(got) != 1 || got[0] == "" {
		t.Fatalf("expected confirm_date.start detail, got %#v", apiErr.Details)
	}
	if got := apiErr.Details["confirm_date.end"]; len(got) != 1 || got[0] == "" {
		t.Fatalf("expected confirm_date.end detail, got %#v", apiErr.Details)
	}
}

func TestUpdateEventValidationHandlesNilProposedDateBounds(t *testing.T) {
	t.Parallel()

	err := UpdateEventValidation(&appmodel.EventDraftUpdate{
		Title: "event",
		ProposedDates: []appmodel.ProposedDate{
			{},
		},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	apiErr, ok := err.(*internalErrors.APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.Kind != internalErrors.KindValidation {
		t.Fatalf("expected validation kind, got %s", apiErr.Kind)
	}
	if got := apiErr.Details["proposed_dates"]; len(got) != 1 || got[0] == "" {
		t.Fatalf("expected proposed_dates detail, got %#v", apiErr.Details)
	}
}
