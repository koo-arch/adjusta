package googlecalendarinfo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	"github.com/koo-arch/adjusta-backend/ent/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/ent/calendar"
	"github.com/koo-arch/adjusta-backend/ent/user"
)

type GoogleCalendarInfoImpl struct {
	client *ent.Client
}

func NewGoogleCalendarInfoRepository(client *ent.Client) *GoogleCalendarInfoImpl {
	return &GoogleCalendarInfoImpl{
		client: client,
	}
}

func (r *GoogleCalendarInfoImpl) Read(ctx context.Context, tx *ent.Tx, id uuid.UUID) (*ent.GoogleCalendarInfo, error) {
	if tx != nil {
		return tx.GoogleCalendarInfo.Get(ctx, id)
	}
	return r.client.GoogleCalendarInfo.Get(ctx, id)
}

func (r *GoogleCalendarInfoImpl) FindByFields(ctx context.Context, tx *ent.Tx, opt GoogleCalendarInfoQueryOptions) (*ent.GoogleCalendarInfo, error) {
	findGoogleCalendarInfo := r.client.GoogleCalendarInfo.Query()
	if tx != nil {
		findGoogleCalendarInfo = tx.GoogleCalendarInfo.Query()
	}

	if opt.GoogleCalendarID != nil {
		findGoogleCalendarInfo = findGoogleCalendarInfo.Where(googlecalendarinfo.GoogleCalendarIDEQ(*opt.GoogleCalendarID))
	}
	if opt.Summary != nil {
		findGoogleCalendarInfo = findGoogleCalendarInfo.Where(googlecalendarinfo.SummaryEQ(*opt.Summary))
	}
	if opt.IsPrimary != nil {
		findGoogleCalendarInfo = findGoogleCalendarInfo.Where(googlecalendarinfo.IsPrimaryEQ(*opt.IsPrimary))
	}

	return findGoogleCalendarInfo.Only(ctx)
}

func (r *GoogleCalendarInfoImpl) ListByUser(ctx context.Context, tx *ent.Tx, userID uuid.UUID) ([]*ent.GoogleCalendarInfo, error) {
	findGoogleCalendarInfo := r.client.GoogleCalendarInfo.Query()
	if tx != nil {
		findGoogleCalendarInfo = tx.GoogleCalendarInfo.Query()
	}

	return findGoogleCalendarInfo.
		Where(googlecalendarinfo.HasCalendarsWith(calendar.HasUserWith(user.IDEQ(userID)))).
		All(ctx)
}

func (r *GoogleCalendarInfoImpl) Create(ctx context.Context, tx *ent.Tx, opt GoogleCalendarInfoQueryOptions, entCalendar *ent.Calendar) (*ent.GoogleCalendarInfo, error) {
	googleCalendarInfoCreate := r.client.GoogleCalendarInfo.Create()
	if tx != nil {
		googleCalendarInfoCreate = tx.GoogleCalendarInfo.Create()
	}

	return googleCalendarInfoCreate.
		SetGoogleCalendarID(*opt.GoogleCalendarID).
		SetSummary(*opt.Summary).
		SetIsPrimary(*opt.IsPrimary).
		AddCalendars(entCalendar).
		Save(ctx)
}

func (r *GoogleCalendarInfoImpl) Update(ctx context.Context, tx *ent.Tx, id uuid.UUID, opt GoogleCalendarInfoQueryOptions, entCalendar *ent.Calendar) (*ent.GoogleCalendarInfo, error) {
	googleCalendarInfoUpdate := r.client.GoogleCalendarInfo.UpdateOneID(id)
	if tx != nil {
		googleCalendarInfoUpdate = tx.GoogleCalendarInfo.UpdateOneID(id)
	}

	if opt.GoogleCalendarID != nil {
		googleCalendarInfoUpdate.SetGoogleCalendarID(*opt.GoogleCalendarID)
	}
	if opt.Summary != nil {
		googleCalendarInfoUpdate.SetSummary(*opt.Summary)
	}
	if opt.IsPrimary != nil {
		googleCalendarInfoUpdate.SetIsPrimary(*opt.IsPrimary)
	}

	if entCalendar != nil {
		googleCalendarInfoUpdate = googleCalendarInfoUpdate.AddCalendars(entCalendar)
	}

	return googleCalendarInfoUpdate.Save(ctx)
}

func (r *GoogleCalendarInfoImpl) Delete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	if tx != nil {
		return tx.GoogleCalendarInfo.DeleteOneID(id).Exec(ctx)
	}
	return r.client.GoogleCalendarInfo.DeleteOneID(id).Exec(ctx)
}

func (r *GoogleCalendarInfoImpl) SoftDelete(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	softDeleteGoogleCalendarInfo := r.client.GoogleCalendarInfo.UpdateOneID(id)
	if tx != nil {
		softDeleteGoogleCalendarInfo = tx.GoogleCalendarInfo.UpdateOneID(id)
	}
	return softDeleteGoogleCalendarInfo.
		SetDeletedAt(time.Now()).
		Exec(ctx)
}

func (r *GoogleCalendarInfoImpl) Restore(ctx context.Context, tx *ent.Tx, id uuid.UUID) error {
	restoreGoogleCalendarInfo := r.client.GoogleCalendarInfo.UpdateOneID(id)
	if tx != nil {
		restoreGoogleCalendarInfo = tx.GoogleCalendarInfo.UpdateOneID(id)
	}
	return restoreGoogleCalendarInfo.
		SetNillableDeletedAt(nil).
		Exec(ctx)
}