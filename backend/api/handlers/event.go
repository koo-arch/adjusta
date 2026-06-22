package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/queryparser"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/validation"
	"github.com/koo-arch/adjusta-backend/internal/errors"
	usecaseEvents "github.com/koo-arch/adjusta-backend/internal/usecase/events"
)

type EventHandler struct {
	eventUsecase api.EventService
}

func NewEventHandler(eventUsecase api.EventService) *EventHandler {
	return &EventHandler{eventUsecase: eventUsecase}
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

func toProposedDateResponse(date usecaseEvents.ProposedDateOutput) dto.ProposedDate {
	return dto.ProposedDate{
		ID:            date.ID,
		GoogleEventID: date.GoogleEventID,
		Start:         date.Start,
		End:           date.End,
		Priority:      date.Priority,
		Status:        date.Status,
		SyncStatus:    date.SyncStatus,
		LastSyncedAt:  date.LastSyncedAt,
		LastSyncError: date.LastSyncError,
	}
}

func toProposedDateResponses(dates []usecaseEvents.ProposedDateOutput) []dto.ProposedDate {
	responses := make([]dto.ProposedDate, 0, len(dates))
	for _, date := range dates {
		responses = append(responses, toProposedDateResponse(date))
	}
	return responses
}

func toEventDraftDetailResponse(event *usecaseEvents.EventDraftDetailOutput) *dto.EventDraftDetail {
	if event == nil {
		return nil
	}

	return &dto.EventDraftDetail{
		ID:                     event.ID,
		Title:                  event.Title,
		Location:               event.Location,
		Description:            event.Description,
		Status:                 event.Status,
		SyncStatus:             event.SyncStatus,
		ConfirmedDateID:        event.ConfirmedDateID,
		GoogleEventID:          event.GoogleEventID,
		ConfirmedGoogleEventID: event.ConfirmedGoogleEventID,
		LastSyncedAt:           event.LastSyncedAt,
		LastSyncError:          event.LastSyncError,
		ProposedDates:          toProposedDateResponses(event.ProposedDates),
	}
}

func toEventDraftDetailResponses(events []*usecaseEvents.EventDraftDetailOutput) []*dto.EventDraftDetail {
	responses := make([]*dto.EventDraftDetail, 0, len(events))
	for _, event := range events {
		responses = append(responses, toEventDraftDetailResponse(event))
	}
	return responses
}

func toUpcomingEventResponse(event usecaseEvents.UpcomingEventOutput) dto.UpcomingEvent {
	return dto.UpcomingEvent{
		ID:                     event.ID,
		Title:                  event.Title,
		Location:               event.Location,
		Description:            event.Description,
		Status:                 event.Status,
		SyncStatus:             event.SyncStatus,
		ConfirmedDateID:        event.ConfirmedDateID,
		GoogleEventID:          event.GoogleEventID,
		ConfirmedGoogleEventID: event.ConfirmedGoogleEventID,
		LastSyncedAt:           event.LastSyncedAt,
		LastSyncError:          event.LastSyncError,
		Start:                  event.Start,
		End:                    event.End,
	}
}

func toUpcomingEventResponses(events []usecaseEvents.UpcomingEventOutput) []dto.UpcomingEvent {
	responses := make([]dto.UpcomingEvent, 0, len(events))
	for _, event := range events {
		responses = append(responses, toUpcomingEventResponse(event))
	}
	return responses
}

func toNeedsActionDraftResponse(event usecaseEvents.NeedsActionDraftOutput) dto.NeedsActionDraft {
	return dto.NeedsActionDraft{
		ID:             event.ID,
		Title:          event.Title,
		Location:       event.Location,
		Description:    event.Description,
		Status:         event.Status,
		Start:          event.Start,
		End:            event.End,
		NeedsAttention: event.NeedsAttention,
	}
}

func toNeedsActionDraftResponses(events []usecaseEvents.NeedsActionDraftOutput) []dto.NeedsActionDraft {
	responses := make([]dto.NeedsActionDraft, 0, len(events))
	for _, event := range events {
		responses = append(responses, toNeedsActionDraftResponse(event))
	}
	return responses
}

func toGoogleEventResponse(event *usecaseEvents.FetchedGoogleEvent) *dto.GoogleEvent {
	if event == nil {
		return nil
	}

	return &dto.GoogleEvent{
		ID:          event.ID,
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		ColorID:     event.ColorID,
		Start:       event.Start,
		End:         event.End,
	}
}

func toGoogleEventResponses(events []*usecaseEvents.FetchedGoogleEvent) []*dto.GoogleEvent {
	responses := make([]*dto.GoogleEvent, 0, len(events))
	for _, event := range events {
		responses = append(responses, toGoogleEventResponse(event))
	}
	return responses
}

func (eh *EventHandler) FetchEventListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := eh.eventUsecase

		accountsEvents, err := eventUsecase.FetchAllGoogleEvents(ctx, userid, email)
		if err, ok := err.(*errors.APIError); ok && err.Kind == errors.KindPartial {
			respond.Partial(c, gin.H{
				"events":  toGoogleEventResponses(accountsEvents),
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
			"events": toGoogleEventResponses(accountsEvents),
		})
	}
}

func (eh *EventHandler) FetchAllEventDraftListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := eh.eventUsecase

		draftedEvents, err := eventUsecase.FetchAllDraftedEvents(ctx, userid, email)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toEventDraftDetailResponses(draftedEvents))
	}
}

func (eh *EventHandler) SearchEventDraftHandler() gin.HandlerFunc {
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

		eventUsecase := eh.eventUsecase

		draftedEvents, err := eventUsecase.SearchDraftedEvents(ctx, userid, email, *query)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toEventDraftDetailResponses(draftedEvents))
	}
}

func (eh *EventHandler) FetchUpcomingEventsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := eh.eventUsecase

		daysBefore := 3
		upcomingEvents, err := eventUsecase.FetchUpcomingEvents(ctx, userid, email, daysBefore)
		if err != nil {
			log.Printf("failed to fetch upcoming events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toUpcomingEventResponses(upcomingEvents))
	}
}

func (eh *EventHandler) FetchNeedsActionDraftsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		eventUsecase := eh.eventUsecase

		daysBefore := 3
		needsActionDrafts, err := eventUsecase.FetchNeedsActionDrafts(ctx, userid, email, daysBefore)
		if err != nil {
			log.Printf("failed to fetch needs action events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toNeedsActionDraftResponses(needsActionDrafts))
	}
}

func (eh *EventHandler) FetchEventDraftDetailHandler() gin.HandlerFunc {
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

		eventUsecase := eh.eventUsecase

		draftedEvent, err := eventUsecase.FetchDraftedEventDetail(ctx, userid, email, eventID)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベント詳細の取得に失敗しました")
			return
		}

		respond.OK(c, toEventDraftDetailResponse(draftedEvent))
	}
}

func (eh *EventHandler) CreateEventDraftHandler() gin.HandlerFunc {
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

		eventUsecase := eh.eventUsecase

		response, err := eventUsecase.CreateDraftedEvents(ctx, userid, email, toDraftCreationRequest(eventDraft))
		if err != nil {
			log.Printf("failed to create events: %v", err)
			respond.Error(c, err, "イベントの作成に失敗しました")
			return
		}

		respond.OK(c, toEventDraftDetailResponse(response))
	}
}

func (eh *EventHandler) EventFinalizeHandler() gin.HandlerFunc {
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

		eventUsecase := eh.eventUsecase

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

func (eh *EventHandler) UpdateEventDraftHandler() gin.HandlerFunc {
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

		eventUsecase := eh.eventUsecase

		err = eventUsecase.UpdateDraftedEvents(ctx, userid, eventID, email, toDraftUpdateRequest(eventDraft))
		if err != nil {
			log.Printf("failed to update events: %v", err)
			respond.Error(c, err, "イベントの更新に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}

func (eh *EventHandler) DeleteEventDraftHandler() gin.HandlerFunc {
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

		eventUsecase := eh.eventUsecase

		err = eventUsecase.DeleteDraftedEvents(ctx, userid, email, eventID)
		if err != nil {
			log.Printf("failed to delete events: %v", err)
			respond.Error(c, err, "イベントの削除に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}
