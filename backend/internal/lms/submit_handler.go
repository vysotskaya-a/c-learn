package lms

import (
	"github.com/c-learn/pkg/response"
	"github.com/gin-gonic/gin"
)

type SubmitHandler struct {
	svc *Service
}

func NewSubmitHandler(svc *Service) *SubmitHandler {
	return &SubmitHandler{svc: svc}
}

type RunRequest struct {
	SourceCode string `json:"source_code"`
	Stdin      string `json:"stdin"`
}

type SubmitRequest struct {
	SourceCode string `json:"source_code"`
}

func (h *SubmitHandler) Run(c *gin.Context) {
	taskID := c.Param("id")
	var req RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	result, appErr := h.svc.RunCode(c.Request.Context(), taskID, req.SourceCode, req.Stdin)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, result)
}

func (h *SubmitHandler) Submit(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	taskID := c.Param("id")
	var req SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	result, appErr := h.svc.Submit(c.Request.Context(), userID, taskID, req.SourceCode)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, result)
}
