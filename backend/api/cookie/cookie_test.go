package cookie

import (
	"net/http"
	"testing"
)

func TestManagerUsesSameSiteLax(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		secure bool
	}{
		{name: "local", secure: false},
		{name: "production", secure: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			manager := NewManager("", tt.secure)
			if got := manager.Options().SameSite; got != http.SameSiteLaxMode {
				t.Fatalf("session cookie SameSite = %v, want %v", got, http.SameSiteLaxMode)
			}
			if got := manager.Expired(SessionCookieName).SameSite; got != http.SameSiteLaxMode {
				t.Fatalf("expired cookie SameSite = %v, want %v", got, http.SameSiteLaxMode)
			}
		})
	}
}
