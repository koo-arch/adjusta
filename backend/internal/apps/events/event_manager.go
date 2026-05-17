package events

import (
	"context"
	"net/http"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/ent"
	appCalendar "github.com/koo-arch/adjusta-backend/internal/apps/calendar"
	"github.com/koo-arch/adjusta-backend/internal/auth"
	internalErrors "github.com/koo-arch/adjusta-backend/internal/errors"
	customCalendar "github.com/koo-arch/adjusta-backend/internal/google/calendar"
	"github.com/koo-arch/adjusta-backend/internal/models"
	repoCalendar "github.com/koo-arch/adjusta-backend/internal/repo/calendar"
	"github.com/koo-arch/adjusta-backend/internal/repo/event"
	"github.com/koo-arch/adjusta-backend/internal/repo/googlecalendarinfo"
	"github.com/koo-arch/adjusta-backend/internal/repo/proposeddate"
	"github.com/koo-arch/adjusta-backend/internal/transaction"
	"github.com/koo-arch/adjusta-backend/utils"
)

type EventManager struct {
	Client       		*ent.Client
	AuthManager  		*auth.AuthManager
	CalendarRepo 		repoCalendar.CalendarRepository
	GoogleCalendarRepo  googlecalendarinfo.GoogleCalendarInfoRepository
	EventRepo    		event.EventRepository
	DateRepo     		proposeddate.ProposedDateRepository
	CalendarApp 		*appCalendar.GoogleCalendarManager
}

func NewEventManager(
	client *ent.Client,
	authManager *auth.AuthManager,
	calendarRepo repoCalendar.CalendarRepository,
	googleCalendarRepo googlecalendarinfo.GoogleCalendarInfoRepository,
	eventRepo event.EventRepository,
	dateRepo proposeddate.ProposedDateRepository,
	calendarApp *appCalendar.GoogleCalendarManager,
) *EventManager {
	return &EventManager{
		Client:       		client,
		AuthManager:  		authManager,
		CalendarRepo: 		calendarRepo,
		GoogleCalendarRepo: googleCalendarRepo,
		EventRepo:    		eventRepo,
		DateRepo:     		dateRepo,
		CalendarApp:  		calendarApp,
	}
}

func (em *EventManager) FinalizeProposedDate(ctx context.Context, userID uuid.UUID, slug, email string, eventReq *models.ConfirmEvent) error {
	tx, err := em.Client.Tx(ctx)
	if err != nil {
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// OAuthトークンを検証
	token, err := em.AuthManager.VerifyOAuthToken(ctx, userID)
	if err != nil {
		log.Printf("failed to verify oauth token for account: %s, error: %v", email, err)
		apiErr := utils.GetAPIError(err, "サーバーでエラーが発生しました")
		return apiErr
	}

	// Google Calendarサービスを作成
	calendarService, err := customCalendar.NewCalendar(ctx, token)
	if err != nil {
		log.Printf("failed to create calendar service for account: %s, error: %v", email, err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, "Googleカレンダー接続に失敗しました")
	}

	// トランザクションをデファーで処理
	defer transaction.HandleTransaction(tx, &err)

	entEvent, err := em.EventRepo.FindBySlugAndUser(ctx, tx, userID, slug, event.EventQueryOptions{})
	if err != nil {
		log.Printf("failed to get event for account: %s, error: %v", email, err)
		if ent.IsNotFound(err) {
			return internalErrors.NewAPIError(http.StatusNotFound, "イベントが見つかりませんでした")
		}	
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
	googleEventID, err := em.HandleGoogleEvent(calendarService, entEvent, eventReq)
	if err != nil {
		log.Printf("failed to handle google event for account: %s, error: %v", email, err)
		apiErr := utils.HandleGoogleAPIError(err)
		return apiErr
	}

	// いずれかの日程候補を確定
	err = em.ConfirmEventDate(ctx, tx, calendarService, googleEventID, eventReq, entEvent)
	if err != nil {
		log.Printf("failed to confirm event date for account: %s, error: %v", email, err)
		return internalErrors.NewAPIError(http.StatusInternalServerError, internalErrors.InternalErrorMessage)
	}

	// トランザクションをコミット
	return nil
}

func (em *EventManager) HandleGoogleEvent(calendarService *customCalendar.Calendar, entEvent *ent.Event, eventReq *models.ConfirmEvent) (*string, error) {
	// Googleカレンダーイベントの新規登録または既存イベントのIDチェック
	var googleEventID *string
	if eventReq.ConfirmDate.ID == nil || entEvent.GoogleEventID == "" {
		// 登録するイベントの情報を作成
		eventDraftCreate := models.EventDraftCreation{
			Title:       entEvent.Summary,
			Location:    entEvent.Location,
			Description: entEvent.Description,
			SelectedDates: []models.SelectedDate{
				{
					Start: *eventReq.ConfirmDate.Start,
					End:   *eventReq.ConfirmDate.End,
				},
			},
		}
		// googleカレンダーにイベントを登録
		googleEvents, err := em.CalendarApp.CreateGoogleEvents(calendarService, &eventDraftCreate)
		if err != nil {
			return nil, fmt.Errorf("failed to insert events, error: %w", err)
		}
		googleEventID = &googleEvents[0].Id

	} else {
		// 既存のGoogleカレンダーイベントIDを使用
		println("日時の変更")
		// Googleカレンダーイベントの更新
		convertGoogleEvent := em.CalendarApp.ConvertToCalendarEvent(&entEvent.GoogleEventID, entEvent.Summary, entEvent.Location, entEvent.Description, *eventReq.ConfirmDate.Start, *eventReq.ConfirmDate.End)
		googleEvent, err := em.CalendarApp.UpdateOrCreateGoogleEvent(calendarService, convertGoogleEvent)
		if err != nil {
			return nil, fmt.Errorf("failed to update events, error: %w", err)
		}
		
		googleEventID = &googleEvent.Id
	}

	return googleEventID, nil
}

func (em *EventManager) ConfirmEventDate(ctx context.Context, tx *ent.Tx, calendarService *customCalendar.Calendar, googleEventID *string, eventReq *models.ConfirmEvent, entEvent *ent.Event) error {
	priority := 1
	dateOptions := proposeddate.ProposedDateQueryOptions{
		Priority:      &priority,
	}

	// 確定日程がDBに存在しない場合は新規作成
	confirmDateID := eventReq.ConfirmDate.ID
	if eventReq.ConfirmDate.ID == nil {
		// オプションに確定日程を設定
		dateOptions.StartTime = eventReq.ConfirmDate.Start
		dateOptions.EndTime = eventReq.ConfirmDate.End

		entDate, err := em.DateRepo.Create(ctx, tx, dateOptions, entEvent)
		if err != nil {
			return fmt.Errorf("failed to create proposed date error: %w", err)
		}
		confirmDateID = &entDate.ID

		// 他の日程候補の優先度を下げる
		err = em.DateRepo.DecrementPriorityExceptID(ctx, tx, entEvent.ID, entDate.ID)
		if err != nil {
			return fmt.Errorf("failed to decrement priority error: %w", err)
		}
	}

	// 確定日程がDBに存在する場合は更新
	if eventReq.ConfirmDate.ID != nil {
		// オプションの優先順位を一度0に設定してから更新
		zero := 0
		dateOptions.Priority = &zero

		entDate, err := em.DateRepo.Update(ctx, tx, *eventReq.ConfirmDate.ID, dateOptions)
		if err != nil {
			return fmt.Errorf("failed to update proposed date error: %w", err)
		}

		// Priorityを振り直す
		err = em.DateRepo.ReorderPriority(ctx, tx, entDate.ID)
		if err != nil {
			return fmt.Errorf("failed to reorder priority error: %w", err)
		}
	}

	// イベントステータスと確定日程IDを更新
	status := models.StatusConfirmed
	eventOptions := event.EventQueryOptions{
		Status: &status,
		ConfirmedDateID: confirmDateID,
		GoogleEventID: googleEventID,
	}
	_, err := em.EventRepo.Update(ctx, tx, entEvent.ID, eventOptions)
	if err != nil {
		return fmt.Errorf("failed to update event status error: %w", err)
	}

	return nil
}
