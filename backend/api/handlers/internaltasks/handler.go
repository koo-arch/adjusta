package internaltasks

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SyncProcessor interface {
	ProcessMessage(ctx context.Context, messageID uuid.UUID) error
}

type Dispatcher interface {
	DispatchPending(ctx context.Context, limit int) error
}

type AuthorizationVerifier interface {
	VerifyAuthorization(ctx context.Context, authorization string) error
}

type Handler struct {
	processor  SyncProcessor
	dispatcher Dispatcher
	verifier   AuthorizationVerifier
}

func NewHandler(processor SyncProcessor, dispatcher Dispatcher, verifier AuthorizationVerifier) *Handler {
	return &Handler{processor: processor, dispatcher: dispatcher, verifier: verifier}
}

func (h *Handler) ProcessGoogleCalendarSync(c *gin.Context) {
	if !h.authorized(c) {
		return
	}
	var request struct {
		OutboxMessageID uuid.UUID `json:"outbox_message_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task payload"})
		return
	}
	if request.OutboxMessageID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task payload"})
		return
	}
	if err := h.processor.ProcessMessage(c.Request.Context(), request.OutboxMessageID); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "calendar sync failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) DispatchOutbox(c *gin.Context) {
	if !h.authorized(c) {
		return
	}
	if err := h.dispatcher.DispatchPending(c.Request.Context(), 100); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "outbox dispatch failed"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) authorized(c *gin.Context) bool {
	if err := h.verifier.VerifyAuthorization(c.Request.Context(), c.GetHeader("Authorization")); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return false
	}
	return true
}
