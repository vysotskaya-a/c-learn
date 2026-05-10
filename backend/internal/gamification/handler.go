package gamification

import (
	"github.com/c-learn/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		response.JSON(c, 401, gin.H{"error": "unauthorized", "message": "missing user identity"})
		return
	}

	profile, appErr := h.svc.GetProfile(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, profile)
}

func (h *Handler) GetLeaderboard(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	board, appErr := h.svc.GetLeaderboard(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, board)
}
