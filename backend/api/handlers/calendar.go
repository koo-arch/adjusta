package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/queryparser"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/validation"
	"github.com/koo-arch/adjusta-backend/internal/errors"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type CalendarHandler struct {
	handler *Handler
}

func NewCalendarHandler(handler *Handler) *CalendarHandler {
	return &CalendarHandler{handler: handler}
}

var extractErrorMessage = "ユーザー情報確認時にエラーが発生しました。"

func parseEventIDParam(c *gin.Context) (uuid.UUID, bool) {
	eventIDParam := c.Param("id")
	if eventIDParam == "" {
		respond.BadRequest(c, "イベントIDがありません")
		return uuid.UUID{}, false
	}

	eventID, err := uuid.Parse(eventIDParam)
	if err != nil {
		respond.BadRequest(c, "イベントIDが不正です")
		return uuid.UUID{}, false
	}

	return eventID, true
}

func toDraftCreationRequest(eventDraft *dto.EventDraftCreation) usecaseEvents.DraftCreationRequest {
	selectedDates := make([]usecaseEvents.SelectedDate, 0, len(eventDraft.SelectedDates))
	for _, date := range eventDraft.SelectedDates {
		selectedDates = append(selectedDates, usecaseEvents.SelectedDate{
			Start:    date.Start,
			End:      date.End,
			Priority: date.Priority,
		})
	}

	return usecaseEvents.DraftCreationRequest{
		Title:         eventDraft.Title,
		Location:      eventDraft.Location,
		Description:   eventDraft.Description,
		SelectedDates: selectedDates,
	}
}

func toDraftUpdateRequest(eventDraft *dto.EventDraftUpdate) usecaseEvents.DraftUpdateRequest {
	proposedDates := make([]usecaseEvents.ProposedDateRequest, 0, len(eventDraft.ProposedDates))
	for _, date := range eventDraft.ProposedDates {
		proposedDates = append(proposedDates, usecaseEvents.ProposedDateRequest{
			ID:            date.ID,
			GoogleEventID: date.GoogleEventID,
			Start:         date.Start,
			End:           date.End,
			Priority:      date.Priority,
		})
	}

	return usecaseEvents.DraftUpdateRequest{
		Title:         eventDraft.Title,
		Location:      eventDraft.Location,
		Description:   eventDraft.Description,
		Status:        eventDraft.Status,
		ProposedDates: proposedDates,
	}
}

func (ch *CalendarHandler) FetchEventListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		accountsEvents, err := eventUsecase.FetchAllGoogleEvents(ctx, userid, email)
		if err, ok := err.(*errors.APIError); ok && err.Kind == errors.KindPartial {
			respond.Partial(c, gin.H{
				"events":  accountsEvents,
				"warning": err.Details,
			})
			return
		}
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "Googleカレンダーのイベント取得に失敗しました")
			return
		}

		respond.OK(c, gin.H{
			"events": accountsEvents,
		})
	}
}

func (ch *CalendarHandler) FetchAllEventDraftListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		draftedEvents, err := eventUsecase.FetchAllDraftedEvents(ctx, userid, email)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, draftedEvents)
	}
}

func (ch *CalendarHandler) SearchEventDraftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// クエリパラメータの取得
		queryparser := queryparser.NewQueryParser(c)

		// クエリパラメータの解析
		query, err := queryparser.ParseSearchEventQuery()
		if err != nil {
			respond.BadRequest(c, "クエリが不正です")
			return
		}

		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		draftedEvents, err := eventUsecase.SearchDraftedEvents(ctx, userid, email, *query)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, draftedEvents)
	}
}

func (ch *CalendarHandler) FetchUpcomingEventsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		daysBefore := 3
		upcomingEvents, err := eventUsecase.FetchUpcomingEvents(ctx, userid, email, daysBefore)
		if err != nil {
			log.Printf("failed to fetch upcoming events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, upcomingEvents)
	}
}

func (ch *CalendarHandler) FetchNeedsActionDraftsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		daysBefore := 3
		upcomingEvents, err := eventUsecase.FetchNeedsActionDrafts(ctx, userid, email, daysBefore)
		if err != nil {
			log.Printf("failed to fetch needs action events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, upcomingEvents)
	}
}

func (ch *CalendarHandler) FetchEventDraftDetailHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventID, ok := parseEventIDParam(c)
		if !ok {
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		draftedEvent, err := eventUsecase.FetchDraftedEventDetail(ctx, userid, email, eventID)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベント詳細の取得に失敗しました")
			return
		}

		respond.OK(c, draftedEvent)
	}
}

func (ch *CalendarHandler) CreateEventDraftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		var eventDraft *dto.EventDraftCreation
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}

		if err := validation.CreateEventValidation(eventDraft); err != nil {
			respond.Error(c, err, "イベントの作成に失敗しました")
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		response, err := eventUsecase.CreateDraftedEvents(ctx, userid, email, toDraftCreationRequest(eventDraft))
		if err != nil {
			log.Printf("failed to create events: %v", err)
			respond.Error(c, err, "イベントの作成に失敗しました")
			return
		}

		respond.OK(c, response)
	}
}

func (ch *CalendarHandler) EventFinalizeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventID, ok := parseEventIDParam(c)
		if !ok {
			return
		}

		var confirmEvent *dto.ConfirmEvent
		if err := c.ShouldBindJSON(&confirmEvent); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}

		if err := validation.FinalizeValidation(confirmEvent); err != nil {
			log.Printf("failed to validate confirm event: %v", err)
			respond.Error(c, err, "イベントの確定に失敗しました")
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		confirmation := usecaseEvents.ConfirmationRequest{
			ID:            confirmEvent.ConfirmDate.ID,
			GoogleEventID: confirmEvent.ConfirmDate.GoogleEventID,
			Start:         confirmEvent.ConfirmDate.Start,
			End:           confirmEvent.ConfirmDate.End,
			Priority:      confirmEvent.ConfirmDate.Priority,
		}

		err = eventUsecase.FinalizeProposedDate(ctx, userid, eventID, email, confirmation)
		if err != nil {
			log.Printf("failed to finalize event: %v", err)
			respond.Error(c, err, "イベントの確定に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}

func (ch *CalendarHandler) UpdateEventDraftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventID, ok := parseEventIDParam(c)
		if !ok {
			return
		}

		var eventDraft *dto.EventDraftUpdate
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}

		if err := validation.UpdateEventValidation(eventDraft); err != nil {
			respond.Error(c, err, "イベントの更新に失敗しました")
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		err = eventUsecase.UpdateDraftedEvents(ctx, userid, eventID, email, toDraftUpdateRequest(eventDraft))
		if err != nil {
			log.Printf("failed to update events: %v", err)
			respond.Error(c, err, "イベントの更新に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}

func (ch *CalendarHandler) DeleteEventDraftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventID, ok := parseEventIDParam(c)
		if !ok {
			return
		}

		eventUsecase := ch.handler.Server.EventUsecase

		err = eventUsecase.DeleteDraftedEvents(ctx, userid, email, eventID)
		if err != nil {
			log.Printf("failed to delete events: %v", err)
			respond.Error(c, err, "イベントの削除に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}
