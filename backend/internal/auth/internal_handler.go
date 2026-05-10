package auth

import (
	"github.com/c-learn/pkg/response"
	"github.com/gin-gonic/gin"
)

type InternalHandler struct {
	svc *Service
}

func NewInternalHandler(svc *Service) *InternalHandler {
	return &InternalHandler{svc: svc}
}

func (h *InternalHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	info, appErr := h.svc.GetUserInfo(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, info)
}
