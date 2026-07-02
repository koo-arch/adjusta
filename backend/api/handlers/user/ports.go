package user

import (
	"context"

	"github.com/google/uuid"
	usecaseAccount "github.com/koo-arch/adjusta-backend/internal/usecase/account"
)

type ProfileUsecase interface {
	FetchGoogleProfile(ctx context.Context, userID uuid.UUID) (*usecaseAccount.GoogleProfile, error)
}
