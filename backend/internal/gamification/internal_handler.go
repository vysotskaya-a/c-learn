package gamification

import (
	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/response"
	"github.com/gin-gonic/gin"
)

type InternalHandler struct {
	svc *Service
}

func NewInternalHandler(svc *Service) *InternalHandler {
	return &InternalHandler{svc: svc}
}

func (h *InternalHandler) AwardXP(c *gin.Context) {
	var req models.XPAwardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	result, appErr := h.svc.AwardXP(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, result)
}
