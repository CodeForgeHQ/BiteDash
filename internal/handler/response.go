package handler

import (
	"log/slog"
	"net/http"

	"bitedash/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func respondError(c *gin.Context, status int, message string) {
	c.JSON(status, dto.ErrorResponse{Error: message})
}

func respondBadRequest(c *gin.Context, message string) {
	respondError(c, http.StatusBadRequest, message)
}

func respondUnauthorized(c *gin.Context, message string) {
	respondError(c, http.StatusUnauthorized, message)
}

func respondInternal(c *gin.Context, message string) {
	respondError(c, http.StatusInternalServerError, message)
}

func respondInternalError(c *gin.Context, message string, err error) {
	slog.Error(message, "error", err, "request_id", c.GetString("requestID"))
	respondInternal(c, message)
}

func authenticatedUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, err := uuid.Parse(c.GetString("userID"))
	if err != nil || userID == uuid.Nil {
		respondUnauthorized(c, "invalid user identity")
		return uuid.Nil, false
	}
	return userID, true
}
