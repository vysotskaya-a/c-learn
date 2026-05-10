package auth

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

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	user, appErr := h.svc.Register(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.Created(c, user)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	tokens, appErr := h.svc.Login(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, tokens)
}

func (h *Handler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	tokens, appErr := h.svc.Refresh(c.Request.Context(), req)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, tokens)
}

func (h *Handler) Me(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		response.JSON(c, 401, gin.H{"error": "unauthorized", "message": "missing user identity"})
		return
	}

	user, appErr := h.svc.GetMe(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, user)
}
