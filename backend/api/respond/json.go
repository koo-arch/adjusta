package respond

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JSON writes a successful or custom HTTP response payload.
func JSON(c *gin.Context, status int, payload any) {
	c.JSON(status, payload)
}

func AbortJSON(c *gin.Context, status int, payload any) {
	JSON(c, status, payload)
	c.Abort()
}

func OK(c *gin.Context, payload any) {
	JSON(c, http.StatusOK, payload)
}

func Partial(c *gin.Context, payload any) {
	JSON(c, http.StatusPartialContent, payload)
}

func Message(c *gin.Context, status int, message string) {
	JSON(c, status, gin.H{"message": message})
}

func OKMessage(c *gin.Context, message string) {
	Message(c, http.StatusOK, message)
}

// BadRequest, Unauthorized, and Internal are for HTTP-local failures decided in the handler
// or middleware itself, such as invalid input or session save failures.
func BadRequest(c *gin.Context, message string) {
	AbortJSON(c, http.StatusBadRequest, gin.H{"error": message})
}

func Unauthorized(c *gin.Context, message string) {
	AbortJSON(c, http.StatusUnauthorized, gin.H{"error": message})
}

func Internal(c *gin.Context, message string) {
	AbortJSON(c, http.StatusInternalServerError, gin.H{"error": message})
}
