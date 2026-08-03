package googlecalendar

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	domainEvent "github.com/koo-arch/adjusta-backend/internal/domain/event"
	domainOutbox "github.com/koo-arch/adjusta-backend/internal/domain/outboxmessage"
	domainProposedDate "github.com/koo-arch/adjusta-backend/internal/domain/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/domain/value"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	"github.com/koo-arch/adjusta-backend/internal/repoerr"
)

type proposedDateSyncError struct {
	proposedDateID uuid.UUID
	err            error
}

func (e *proposedDateSyncError) Error() string { return e.err.Error() }
func (e *proposedDateSyncError) Unwrap() error { return e.err }

type SyncUsecase struct {
	repos   Repositories
	tx      Transaction
	gateway EventGateway
}

func NewSyncUsecase(repos Repositories, tx Transaction, gateway EventGateway) *SyncUsecase {
	return &SyncUsecase{repos: repos, tx: tx, gateway: gateway}
}

func (uc *SyncUsecase) ProcessMessage(ctx context.Context, messageID uuid.UUID) error {
	message, err := uc.repos.Message.Read(ctx, messageID)
	if err != nil {
		if repoerr.IsNotFound(err) {
			return nil
		}
		return err
	}
	if message.ProcessedAt != nil {
		return nil
	}
	if message.MessageType != domainOutbox.MessageTypeCalendarEventSync || message.AggregateType != domainOutbox.AggregateTypeEvent {
		return uc.markOutboxMessageProcessed(ctx, messageID)
	}

	storedEvent, err := uc.repos.Event.Read(ctx, message.AggregateID, domainEvent.EventReadOptions{
		WithProposedDates: true,
		IncludeDeleted:    true,
	})
	if err != nil {
		if repoerr.IsNotFound(err) {
			return uc.markOutboxMessageProcessed(ctx, messageID)
		}
		return err
	}

	if storedEvent.DeletedAt != nil {
		err = uc.syncDeletedEvent(ctx, storedEvent)
	} else if storedEvent.Status == value.StatusConfirmed {
		err = uc.syncConfirmedEvent(ctx, storedEvent)
	} else {
		err = uc.syncDraftEvent(ctx, storedEvent)
	}
	if err != nil {
		return uc.recordSyncFailure(ctx, messageID, storedEvent.ID, err)
	}
	return uc.markOutboxMessageProcessed(ctx, messageID)
}

func (uc *SyncUsecase) syncDraftEvent(ctx context.Context, storedEvent *domainEvent.Event) error {
	candidateCalendar, err := uc.loadCandidateCalendar(ctx, storedEvent.UserID)
	if err != nil {
		return err
	}
	if candidateCalendar == nil || !candidateCalendar.syncProposedDates || candidateCalendar.googleCalendarID == "" {
		return uc.markDraftEventNotSynced(ctx, storedEvent)
	}

	proposedDates := slices.Clone(storedEvent.ProposedDates)
	sort.SliceStable(proposedDates, func(i, j int) bool {
		return proposedDates[i] != nil && (proposedDates[j] == nil || proposedDates[i].Priority > proposedDates[j].Priority)
	})

	rank := 0
	for _, proposedDate := range proposedDates {
		if proposedDate == nil {
			continue
		}
		if proposedDate.DeletedAt != nil {
			if proposedDate.GoogleEventID != nil {
				if err := uc.gateway.DeleteEvent(ctx, storedEvent.UserID, candidateCalendar.googleCalendarID, *proposedDate.GoogleEventID); err != nil {
					return &proposedDateSyncError{proposedDateID: proposedDate.ID, err: err}
				}
			}
			if err := uc.markProposedDateSynced(ctx, proposedDate.ID); err != nil {
				return err
			}
			continue
		}
		rank++
		existingGoogleEventID := proposedDate.GoogleEventID
		if existingGoogleEventID == nil || *existingGoogleEventID == "" {
			id := stableGoogleEventID("adjustap", proposedDate.ID)
			existingGoogleEventID = &id
		}
		googleEventID, err := uc.gateway.UpsertEvent(
			ctx,
			storedEvent.UserID,
			candidateCalendar.googleCalendarID,
			existingGoogleEventID,
			fmt.Sprintf("%s【第%d候補】", storedEvent.Title, rank),
			storedEvent.Location,
			storedEvent.Description,
			proposedDate.StartTime,
			proposedDate.EndTime,
		)
		if err != nil {
			return &proposedDateSyncError{proposedDateID: proposedDate.ID, err: err}
		}
		if err := uc.markProposedDateSyncedWithGoogleEvent(ctx, proposedDate.ID, googleEventID); err != nil {
			return err
		}
	}

	return uc.markEventSynced(ctx, storedEvent.ID)
}

func (uc *SyncUsecase) markDraftEventNotSynced(ctx context.Context, storedEvent *domainEvent.Event) error {
	return uc.tx.Do(ctx, func(repos Repositories) error {
		for _, proposedDate := range storedEvent.ProposedDates {
			if proposedDate == nil {
				continue
			}
			if _, err := repos.ProposedDate.Update(ctx, proposedDate.ID, proposedDateNotSyncedUpdate()); err != nil {
				return err
			}
		}
		_, err := repos.Event.Update(ctx, storedEvent.ID, eventNotSyncedUpdate())
		return err
	})
}

func (uc *SyncUsecase) syncConfirmedEvent(ctx context.Context, storedEvent *domainEvent.Event) error {
	primaryCalendar, err := uc.repos.Calendar.Read(ctx, storedEvent.PrimaryCalendarID)
	if err != nil {
		return err
	}
	confirmedDate := findProposedDate(storedEvent.ProposedDates, storedEvent.ConfirmedDateID)
	if confirmedDate == nil {
		return errors.New("confirmed proposed date not found")
	}

	candidateCalendar, candidateErr := uc.loadCandidateCalendar(ctx, storedEvent.UserID)
	if candidateErr != nil {
		return candidateErr
	}
	if candidateCalendar != nil && candidateCalendar.googleCalendarID != "" {
		for _, proposedDate := range storedEvent.ProposedDates {
			if proposedDate == nil || proposedDate.GoogleEventID == nil {
				continue
			}
			suffix := "未選択"
			if proposedDate.ID == confirmedDate.ID {
				suffix = "確定済み"
			}
			_, err := uc.gateway.UpsertEvent(
				ctx,
				storedEvent.UserID,
				candidateCalendar.googleCalendarID,
				proposedDate.GoogleEventID,
				fmt.Sprintf("%s【%s】", storedEvent.Title, suffix),
				storedEvent.Location,
				storedEvent.Description,
				proposedDate.StartTime,
				proposedDate.EndTime,
			)
			if err != nil {
				return &proposedDateSyncError{proposedDateID: proposedDate.ID, err: err}
			}
		}
	}

	existingGoogleEventID := storedEvent.ConfirmedGoogleEventID
	if existingGoogleEventID == nil || *existingGoogleEventID == "" {
		id := stableGoogleEventID("adjustae", storedEvent.ID)
		existingGoogleEventID = &id
	}
	googleEventID, err := uc.gateway.UpsertEvent(
		ctx,
		storedEvent.UserID,
		primaryCalendar.GoogleCalendarID,
		existingGoogleEventID,
		storedEvent.Title,
		storedEvent.Location,
		storedEvent.Description,
		confirmedDate.StartTime,
		confirmedDate.EndTime,
	)
	if err != nil {
		return &proposedDateSyncError{proposedDateID: confirmedDate.ID, err: err}
	}

	now := time.Now()
	return uc.tx.Do(ctx, func(repos Repositories) error {
		for _, proposedDate := range storedEvent.ProposedDates {
			if proposedDate == nil {
				continue
			}
			if _, err := repos.ProposedDate.Update(ctx, proposedDate.ID, proposedDateSyncedUpdate(now)); err != nil {
				return err
			}
		}
		_, err := repos.Event.Update(ctx, storedEvent.ID, confirmedEventSyncedUpdate(confirmedDate.ID, googleEventID, now))
		return err
	})
}

func (uc *SyncUsecase) syncDeletedEvent(ctx context.Context, storedEvent *domainEvent.Event) error {
	if storedEvent.ConfirmedGoogleEventID != nil {
		primaryCalendar, err := uc.repos.Calendar.Read(ctx, storedEvent.PrimaryCalendarID)
		if err != nil {
			return err
		}
		if err := uc.gateway.DeleteEvent(ctx, storedEvent.UserID, primaryCalendar.GoogleCalendarID, *storedEvent.ConfirmedGoogleEventID); err != nil {
			return err
		}
	}

	candidateCalendar, err := uc.loadCandidateCalendar(ctx, storedEvent.UserID)
	if err != nil {
		return err
	}
	if candidateCalendar != nil {
		for _, proposedDate := range storedEvent.ProposedDates {
			if proposedDate != nil && proposedDate.GoogleEventID != nil {
				if err := uc.gateway.DeleteEvent(ctx, storedEvent.UserID, candidateCalendar.googleCalendarID, *proposedDate.GoogleEventID); err != nil {
					return &proposedDateSyncError{proposedDateID: proposedDate.ID, err: err}
				}
			}
		}
	}
	return uc.markDeletedEventSynced(ctx, storedEvent)
}

func (uc *SyncUsecase) markProposedDateSyncedWithGoogleEvent(ctx context.Context, id uuid.UUID, googleEventID string) error {
	return uc.tx.Do(ctx, func(repos Repositories) error {
		_, err := repos.ProposedDate.Update(ctx, id, proposedDateSyncedWithGoogleEventUpdate(googleEventID, time.Now()))
		return err
	})
}

func (uc *SyncUsecase) markProposedDateSynced(ctx context.Context, id uuid.UUID) error {
	return uc.tx.Do(ctx, func(repos Repositories) error {
		_, err := repos.ProposedDate.Update(ctx, id, proposedDateSyncedUpdate(time.Now()))
		return err
	})
}

func (uc *SyncUsecase) markDeletedEventSynced(ctx context.Context, storedEvent *domainEvent.Event) error {
	return uc.tx.Do(ctx, func(repos Repositories) error {
		now := time.Now()
		for _, proposedDate := range storedEvent.ProposedDates {
			if proposedDate == nil {
				continue
			}
			if _, err := repos.ProposedDate.Update(ctx, proposedDate.ID, proposedDateSyncedUpdate(now)); err != nil {
				return err
			}
		}
		_, err := repos.Event.Update(ctx, storedEvent.ID, eventSyncedUpdate(now))
		return err
	})
}

func (uc *SyncUsecase) markEventSynced(ctx context.Context, id uuid.UUID) error {
	return uc.tx.Do(ctx, func(repos Repositories) error {
		_, err := repos.Event.Update(ctx, id, eventSyncedUpdate(time.Now()))
		return err
	})
}

func (uc *SyncUsecase) recordSyncFailure(ctx context.Context, messageID, eventID uuid.UUID, syncErr error) error {
	terminal := isTerminalSyncError(syncErr)
	err := uc.tx.Do(ctx, func(repos Repositories) error {
		var proposedErr *proposedDateSyncError
		if errors.As(syncErr, &proposedErr) {
			if _, err := repos.ProposedDate.Update(ctx, proposedErr.proposedDateID, proposedDateFailedUpdate(syncErr)); err != nil {
				return err
			}
		}
		if _, err := repos.Event.Update(ctx, eventID, eventFailedUpdate(syncErr)); err != nil {
			return err
		}
		if terminal {
			return repos.Message.MarkProcessed(ctx, messageID, time.Now())
		}
		return nil
	})
	if err != nil {
		return err
	}
	if terminal {
		return nil
	}
	return syncErr
}

func (uc *SyncUsecase) markOutboxMessageProcessed(ctx context.Context, messageID uuid.UUID) error {
	return uc.tx.Do(ctx, func(repos Repositories) error {
		return repos.Message.MarkProcessed(ctx, messageID, time.Now())
	})
}

func findProposedDate(dates []*domainProposedDate.ProposedDate, id uuid.UUID) *domainProposedDate.ProposedDate {
	for _, date := range dates {
		if date != nil && date.ID == id {
			return date
		}
	}
	return nil
}

func stableGoogleEventID(prefix string, id uuid.UUID) string {
	return prefix + strings.ReplaceAll(id.String(), "-", "")
}

func isTerminalSyncError(err error) bool {
	var apiErr *internalErrors.APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.Kind == internalErrors.KindGoogleReauth ||
		apiErr.Kind == internalErrors.KindBadRequest ||
		apiErr.Kind == internalErrors.KindForbidden ||
		apiErr.Kind == internalErrors.KindNotFound
}
