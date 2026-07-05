package auth

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	repoSession "github.com/koo-arch/adjusta-backend/internal/domain/session"
	repoUser "github.com/koo-arch/adjusta-backend/internal/domain/user"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/google"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type signInPlan struct {
	userID    uuid.UUID
	accountID uuid.UUID
	mode      signInMode
}

type signInMode int

const (
	signInModeCreateUserAndAccount signInMode = iota
	signInModeUpdateUserAndCreateAccount
	signInModeUpdateUserAndAccount
)

func (am *Authenticator) SignInWithGoogle(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoSession.Session, *repoUser.User, error) {
	plan, err := am.resolveSignInPlan(ctx, userInfo)
	if err != nil {
		return nil, nil, err
	}

	var (
		u             *repoUser.User
		storedSession *repoSession.Session
	)

	sessionToken := uuid.NewString()
	expiresAt := repoSession.ExpiresAtFrom(time.Now(), am.sessionLifetime)

	err = am.tx.Do(ctx, func(repos AuthTxRepositories) error {
		var err error
		u, err = am.persistSignIn(ctx, repos, plan, userInfo, oauthToken)
		if err != nil {
			return err
		}

		storedSession, err = repos.Session.Create(ctx, u.ID, sessionToken, expiresAt)
		if err != nil {
			log.Printf("failed to create session during google sign in: %v", err)
			return internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return storedSession, u, nil
}

func (am *Authenticator) ProcessUserSignIn(ctx context.Context, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoUser.User, error) {
	plan, err := am.resolveSignInPlan(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	var u *repoUser.User
	err = am.tx.Do(ctx, func(repos AuthTxRepositories) error {
		var err error
		u, err = am.persistSignIn(ctx, repos, plan, userInfo, oauthToken)
		return err
	})
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (am *Authenticator) resolveSignInPlan(ctx context.Context, userInfo *google.UserProfile) (*signInPlan, error) {
	u, err := am.repos.User.FindByEmail(ctx, userInfo.Email)
	if err != nil {
		log.Printf("failed to get user by email: %v", err)
		if repoerr.IsNotFound(err) {
			return &signInPlan{mode: signInModeCreateUserAndAccount}, nil
		}
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	accountInfo, err := am.repos.Account.FindByUserID(ctx, u.ID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return &signInPlan{
				userID: u.ID,
				mode:   signInModeUpdateUserAndCreateAccount,
			}, nil
		}
		log.Printf("failed to get account by user id: %v", err)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}

	return &signInPlan{
		userID:    u.ID,
		accountID: accountInfo.ID,
		mode:      signInModeUpdateUserAndAccount,
	}, nil
}

func (am *Authenticator) persistSignIn(ctx context.Context, repos AuthTxRepositories, plan *signInPlan, userInfo *google.UserProfile, oauthToken *google.AuthToken) (*repoUser.User, error) {
	switch plan.mode {
	case signInModeCreateUserAndAccount:
		u, err := repos.User.Create(ctx, userInfo.Email, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to create user: %v", err)
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = repos.Account.Create(ctx, u.ID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account: %v", err)
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return u, nil
	case signInModeUpdateUserAndCreateAccount:
		u, err := repos.User.Update(ctx, plan.userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user during account creation: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("ユーザーが見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = repos.Account.Create(ctx, plan.userID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to create account for existing user: %v", err)
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return u, nil
	case signInModeUpdateUserAndAccount:
		u, err := repos.User.Update(ctx, plan.userID, buildUserMutationOptions(userInfo))
		if err != nil {
			log.Printf("failed to update user login state: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("ユーザーが見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		_, err = repos.Account.Update(ctx, plan.accountID, buildAccountMutationOptions(userInfo, oauthToken))
		if err != nil {
			log.Printf("failed to update account token: %v", err)
			if repoerr.IsNotFound(err) {
				return nil, internalErrors.NewNotFoundError("アカウント情報が見つかりませんでした")
			}
			return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
		}

		return u, nil
	default:
		log.Printf("unexpected sign in mode: %d", plan.mode)
		return nil, internalErrors.NewInternalError(internalErrors.InternalErrorMessage)
	}
}
