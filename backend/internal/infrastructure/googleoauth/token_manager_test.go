package googleoauth

import (
	"errors"
	"fmt"
	"testing"

	"golang.org/x/oauth2"
)

func TestIsGoogleReauthorizationRequired(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "invalid grant",
			err:  &oauth2.RetrieveError{ErrorCode: "invalid_grant"},
			want: true,
		},
		{
			name: "wrapped invalid grant",
			err:  fmt.Errorf("refresh token: %w", &oauth2.RetrieveError{ErrorCode: "invalid_grant"}),
			want: true,
		},
		{
			name: "invalid token message",
			err:  errors.New(`oauth2: "invalid_token"`),
			want: true,
		},
		{
			name: "temporary failure",
			err:  &oauth2.RetrieveError{ErrorCode: "temporarily_unavailable"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGoogleReauthorizationRequired(tt.err); got != tt.want {
				t.Fatalf("isGoogleReauthorizationRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
