package googleapierror

import (
	"testing"

	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"google.golang.org/api/googleapi"
)

func TestNormalizeReturnsGoogleReauthorizationForAuthzFailures(t *testing.T) {
	t.Parallel()

	for _, statusCode := range []int{401, 403} {
		err := Normalize(&googleapi.Error{Code: statusCode})
		apiErr, ok := err.(*internalErrors.APIError)
		if !ok {
			t.Fatalf("expected APIError for status %d, got %T", statusCode, err)
		}
		if apiErr.Kind != internalErrors.KindGoogleReauth {
			t.Fatalf("unexpected kind for status %d: %s", statusCode, apiErr.Kind)
		}
	}
}
