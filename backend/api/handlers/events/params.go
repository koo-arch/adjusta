package events

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/koo-arch/adjusta-backend/api/respond"
)

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
