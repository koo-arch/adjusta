package respond

import (
	"net/http"
	"testing"

	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
)

func TestStatusCodeForKindGoogleReauthorization(t *testing.T) {
	t.Parallel()

	if got := statusCodeForKind(internalErrors.KindGoogleReauth); got != http.StatusConflict {
		t.Fatalf("unexpected status code: %d", got)
	}
}
