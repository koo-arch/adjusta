package user

import (
	"context"
	"fmt"
	"time"

	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/event"
	"github.com/koo-arch/adjusta-backend/ent/proposeddate"
	"github.com/koo-arch/adjusta-backend/ent/oauthtoken"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/google/uuid"
)

type UserRepositoryImpl struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		client: client,
	}
}

func (r *UserRepositoryImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt UserQueryOptions) (*ent.User, error) {
	findQuery := r.client.User.Query()
	if tx != nil {
		findQuery = tx.User.Query()
	}

	if opt.WithOAuthToken {
		findQuery = findQuery.WithOauthToken()
	}

	return findQuery.
		Where(user.IDEQ(id)).
		Only(ctx)
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, tx *ent.Tx, email string, opt UserQueryOptions) (*ent.User, error) {
	findUser := r.client.User.Query()
	if tx != nil {
		findUser = tx.User.Query()
	}

	if opt.WithOAuthToken {
		findUser = findUser.WithOauthToken()
	}
	
	return findUser.
		Where(user.EmailEQ(email)).
		Only(ctx)
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *ent.Tx, email string, jwtToken *models.JWTToken) (*ent.User, error) {
	userCreate := r.client.User.Create()
	if tx != nil {
		userCreate = tx.User.Create()
	}
	return userCreate.
		SetEmail(email).
		SetNillableRefreshToken(&jwtToken.RefreshToken).
		SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
		Save(ctx)
}

func (r *UserRepositoryImpl) Update(ctx context.Context,tx *ent.Tx, id uuid.UUID, jwtToken *models.JWTToken) (*ent.User, error) {
	userUpdate := r.client.User.UpdateOneID(id)
	if tx != nil {
		userUpdate = tx.User.UpdateOneID(id)
	}
	return userUpdate.
		SetNillableRefreshToken(&jwtToken.RefreshToken).
		SetNillableRefreshTokenExpiry(&jwtToken.RefreshExpiration).
		Save(ctx)
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.User.DeleteOneID(id).Exec(ctx)
	}
	return r.client.User.DeleteOneID(id).Exec(ctx)
}

func (r *UserRepositoryImpl) SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	softDeleteUser := r.client.User.UpdateOneID(id)
	if tx != nil {
		softDeleteUser = tx.User.UpdateOneID(id)
	}
	return softDeleteUser.
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *UserRepositoryImpl) Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	restoreUser := r.client.User.UpdateOneID(id)
	if tx != nil {
		restoreUser = tx.User.UpdateOneID(id)
	}
	return restoreUser.
		SetNillableDeletedAt(nil).
		Exec(ctx)
}

func (r *UserRepositoryImpl) SoftDeleteWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	// ユーザーを論理削除
	if err := r.SoftDelete(ctx, tx, id); err != nil {
		return err
	}

	// OAuthTokenを論理削除
	if err := tx.OAuthToken.Update().
		Where(oauthtoken.HasUserWith(user.IDEQ(id))).
		SetDeletedAt(time.Now()).
		Exec(ctx); err != nil {
		return err
	}

	// カレンダーの論理削除
	calendarIDs, err := tx.Calendar.
		Query().
		Where(calendar.HasUserWith(user.IDEQ(id))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query calendars: %w", err)
	}
	if len(calendarIDs) > 0 {
		if err := tx.Calendar.Update().
			Where(calendar.IDIn(calendarIDs...)).
			SetDeletedAt(time.Now()).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to soft delete calendars: %w", err)
		}

		// 関連するイベントIDを取得
		eventIDs, err := tx.Event.
			Query().
			Where(event.HasCalendarWith(calendar.IDIn(calendarIDs...))).
			IDs(ctx)
		if err != nil {
			return fmt.Errorf("failed to query events: %w", err)
		}
		if len(eventIDs) > 0 {
			// 関連するイベントを論理削除
			if err := tx.Event.Update().
				Where(event.IDIn(eventIDs...)).
				SetDeletedAt(time.Now()).
				Exec(ctx); err != nil {
				return fmt.Errorf("failed to soft delete events: %w", err)
			}

			// 関連する提案日を論理削除
			if err := tx.ProposedDate.Update().
				Where(proposeddate.HasEventWith(event.IDIn(eventIDs...))).
				SetDeletedAt(time.Now()).
				Exec(ctx); err != nil {
				return fmt.Errorf("failed to soft delete proposed dates: %w", err)
			}

		}

	}

	return nil
}

func (r *UserRepositoryImpl) RestoreWithRelations(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx == nil {
		return fmt.Errorf("transaction is nil")
	}

	// ユーザーを復元
	if err := r.Restore(ctx, tx, id); err != nil {
		return fmt.Errorf("failed to restore user: %w", err)
	}

	// OAuthToken を復元
	oauthTokenIDs, err := tx.OAuthToken.Query().
		Where(oauthtoken.HasUserWith(user.IDEQ(id))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query oauth tokens: %w", err)
	}
	if len(oauthTokenIDs) > 0 {
		if err := tx.OAuthToken.Update().
			Where(oauthtoken.IDIn(oauthTokenIDs...)).
			SetNillableDeletedAt(nil).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to restore OAuthToken: %w", err)
		}
	}

	// カレンダーを復元
	calendarIDs, err := tx.Calendar.Query().
		Where(calendar.HasUserWith(user.IDEQ(id))).
		IDs(ctx)
	if err != nil {
		return fmt.Errorf("failed to query calendars: %w", err)
	}
	if len(calendarIDs) > 0 {
		if err := tx.Calendar.Update().
			Where(calendar.IDIn(calendarIDs...)).
			SetNillableDeletedAt(nil).
			Exec(ctx); err != nil {
			return fmt.Errorf("failed to restore calendars: %w", err)
		}

		// 関連するイベントを復元
		eventIDs, err := tx.Event.Query().
			Where(event.HasCalendarWith(calendar.IDIn(calendarIDs...))).
			IDs(ctx)
		if err != nil {
			return fmt.Errorf("failed to query events: %w", err)
		}
		if len(eventIDs) > 0 {
			if err := tx.Event.Update().
				Where(event.IDIn(eventIDs...)).
				SetNillableDeletedAt(nil).
				Exec(ctx); err != nil {
				return fmt.Errorf("failed to restore events: %w", err)
			}

			// 関連する提案日を復元
			proposedDateIDs, err := tx.ProposedDate.Query().
				Where(proposeddate.HasEventWith(event.IDIn(eventIDs...))).
				IDs(ctx)
			if err != nil {
				return fmt.Errorf("failed to query proposed dates: %w", err)
			}
			if len(proposedDateIDs) > 0 {
				if err := tx.ProposedDate.Update().
					Where(proposeddate.IDIn(proposedDateIDs...)).
					SetNillableDeletedAt(nil).
					Exec(ctx); err != nil {
					return fmt.Errorf("failed to restore proposed dates: %w", err)
				}
			}
		}
	}

	return nil
}