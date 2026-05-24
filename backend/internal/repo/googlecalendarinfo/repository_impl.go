package googlecalendarinfo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/ent/user"
	"github.com/koo-arch/adjusta-backend/internal/models"
	"github.com/koo-arch/adjusta-backend/internal/repo/infraerr"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
)

type GoogleCalendarInfoImpl struct {
	client *ent.Client
}

func NewGoogleCalendarInfoRepository(client *ent.Client) *GoogleCalendarInfoImpl {
	return &GoogleCalendarInfoImpl{
		client: client,
	}
}

func (r *GoogleCalendarInfoImpl) WithTx(tx transaction.Tx) GoogleCalendarInfoRepository {
	return &GoogleCalendarInfoImpl{client: tx.Client()}
}

func (r *GoogleCalendarInfoImpl) Read(ctx context.Context, id uuid.UUID) (*models.GoogleCalendarInfo, error) {
	entity, err := r.client.GoogleCalendarInfo.Get(ctx, id)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelGoogleCalendarInfo(entity), nil
}

func (r *GoogleCalendarInfoImpl) FindByFields(ctx context.Context, opt GoogleCalendarInfoQueryOptions) (*models.GoogleCalendarInfo, error) {
	findGoogleCalendarInfo := r.client.GoogleCalendarInfo.Query()

	if opt.GoogleCalendarID != nil {
		findGoogleCalendarInfo = findGoogleCalendarInfo.Where(googlecalendarinfo.GoogleCalendarIDEQ(*opt.GoogleCalendarID))
	}
	if opt.Summary != nil {
		findGoogleCalendarInfo = findGoogleCalendarInfo.Where(googlecalendarinfo.SummaryEQ(*opt.Summary))
	}
	if opt.IsPrimary != nil {
		findGoogleCalendarInfo = findGoogleCalendarInfo.Where(googlecalendarinfo.IsPrimaryEQ(*opt.IsPrimary))
	}

	entity, err := findGoogleCalendarInfo.Only(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelGoogleCalendarInfo(entity), nil
}

func (r *GoogleCalendarInfoImpl) ListByUser(ctx context.Context, userID uuid.UUID) ([]*models.GoogleCalendarInfo, error) {
	entities, err := r.client.GoogleCalendarInfo.Query().
		Where(googlecalendarinfo.HasCalendarsWith(calendar.HasUserWith(user.IDEQ(userID)))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	modelsInfo := make([]*models.GoogleCalendarInfo, 0, len(entities))
	for _, entity := range entities {
		modelsInfo = append(modelsInfo, toModelGoogleCalendarInfo(entity))
	}
	return modelsInfo, nil
}

func (r *GoogleCalendarInfoImpl) Create(ctx context.Context, opt GoogleCalendarInfoQueryOptions, calendarID uuid.UUID) (*models.GoogleCalendarInfo, error) {
	googleCalendarInfoCreate := r.client.GoogleCalendarInfo.Create()

	entity, err := googleCalendarInfoCreate.
		SetGoogleCalendarID(*opt.GoogleCalendarID).
		SetSummary(*opt.Summary).
		SetIsPrimary(*opt.IsPrimary).
		AddCalendarIDs(calendarID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return toModelGoogleCalendarInfo(entity), nil
}

func (r *GoogleCalendarInfoImpl) Update(ctx context.Context, id uuid.UUID, opt GoogleCalendarInfoQueryOptions, calendarID *uuid.UUID) (*models.GoogleCalendarInfo, error) {
	googleCalendarInfoUpdate := r.client.GoogleCalendarInfo.UpdateOneID(id)

	if opt.GoogleCalendarID != nil {
		googleCalendarInfoUpdate.SetGoogleCalendarID(*opt.GoogleCalendarID)
	}
	if opt.Summary != nil {
		googleCalendarInfoUpdate.SetSummary(*opt.Summary)
	}
	if opt.IsPrimary != nil {
		googleCalendarInfoUpdate.SetIsPrimary(*opt.IsPrimary)
	}

	if calendarID != nil {
		googleCalendarInfoUpdate = googleCalendarInfoUpdate.AddCalendarIDs(*calendarID)
	}

	entity, err := googleCalendarInfoUpdate.Save(ctx)
	if err != nil {
		return nil, infraerr.MapNotFound(err)
	}
	return toModelGoogleCalendarInfo(entity), nil
}

func (r *GoogleCalendarInfoImpl) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.client.GoogleCalendarInfo.DeleteOneID(id).Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *GoogleCalendarInfoImpl) SoftDelete(ctx context.Context, id uuid.UUID) error {
	err := r.client.GoogleCalendarInfo.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func (r *GoogleCalendarInfoImpl) Restore(ctx context.Context, id uuid.UUID) error {
	err := r.client.GoogleCalendarInfo.UpdateOneID(id).
		SetNillableDeletedAt(nil).
		Exec(ctx)
	return infraerr.MapNotFound(err)
}

func toModelGoogleCalendarInfo(entity *ent.GoogleCalendarInfo) *models.GoogleCalendarInfo {
	if entity == nil {
		return nil
	}

	return &models.GoogleCalendarInfo{
		ID:               entity.ID,
		GoogleCalendarID: entity.GoogleCalendarID,
		Summary:          entity.Summary,
		IsPrimary:        entity.IsPrimary,
	}
}
