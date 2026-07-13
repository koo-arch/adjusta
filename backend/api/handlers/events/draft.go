package events

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/koo-arch/adjusta-backend/api/dto"
	"github.com/koo-arch/adjusta-backend/api/queryparser"
	"github.com/koo-arch/adjusta-backend/api/requestctx"
	"github.com/koo-arch/adjusta-backend/api/respond"
	"github.com/koo-arch/adjusta-backend/api/validation"
)

func (eh *Handler) FetchAllEventDraftListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		queryparser := queryparser.NewQueryParser(c)
		query, err := queryparser.ParseEventListQuery()
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

		eventUsecase := eh.draftUsecase

		draftedEvents, err := eventUsecase.FetchDraftedEventsPage(ctx, userid, email, *query)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toEventDraftListResponse(draftedEvents))
	}
}

func (eh *Handler) SearchEventDraftHandler() gin.HandlerFunc {
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

		eventUsecase := eh.draftUsecase

		draftedEvents, err := eventUsecase.SearchDraftedEventsPage(ctx, userid, email, *query)
		if err != nil {
			log.Printf("failed to fetch events: %v", err)
			respond.Error(c, err, "イベントの取得に失敗しました")
			return
		}

		respond.OK(c, toEventDraftListResponse(draftedEvents))
	}
}

func (eh *Handler) CreateEventDraftHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		userid, email, err := requestctx.UserIDAndEmail(c)
		if err != nil {
			respond.Error(c, err, extractErrorMessage)
			return
		}

		var eventDraft dto.EventDraftCreation
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}

		if err := validation.CreateEventValidation(&eventDraft); err != nil {
			respond.Error(c, err, "イベントの作成に失敗しました")
			return
		}

		eventUsecase := eh.draftUsecase

		response, err := eventUsecase.CreateDraftedEvents(ctx, userid, email, toDraftCreationRequest(&eventDraft))
		if err != nil {
			log.Printf("failed to create events: %v", err)
			respond.Error(c, err, "イベントの作成に失敗しました")
			return
		}

		respond.OK(c, toEventDraftDetailResponse(response))
	}
}

func (eh *Handler) UpdateEventDraftHandler() gin.HandlerFunc {
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

		var eventDraft dto.EventDraftUpdate
		if err := c.ShouldBindJSON(&eventDraft); err != nil {
			respond.BadRequest(c, "リクエストのデータ形式が不正です")
			return
		}

		if err := validation.UpdateEventValidation(&eventDraft); err != nil {
			respond.Error(c, err, "イベントの更新に失敗しました")
			return
		}

		eventUsecase := eh.draftUsecase

		err = eventUsecase.UpdateDraftedEvents(ctx, userid, eventID, email, toDraftUpdateRequest(&eventDraft))
		if err != nil {
			log.Printf("failed to update events: %v", err)
			respond.Error(c, err, "イベントの更新に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}

func (eh *Handler) DeleteEventDraftHandler() gin.HandlerFunc {
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

		eventUsecase := eh.draftUsecase

		err = eventUsecase.DeleteDraftedEvents(ctx, userid, email, eventID)
		if err != nil {
			log.Printf("failed to delete events: %v", err)
			respond.Error(c, err, "イベントの削除に失敗しました")
			return
		}

		respond.OKMessage(c, "success")
	}
}
