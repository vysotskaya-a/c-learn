package lms

import (
	"strconv"

	"github.com/c-learn/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetCourseTree(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	tree, appErr := h.svc.GetCourseTree(c.Request.Context(), userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, gin.H{"modules": tree})
}

func (h *Handler) GetLesson(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	lessonID := c.Param("id")
	lesson, appErr := h.svc.GetLesson(c.Request.Context(), lessonID, userID)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, lesson)
}

func (h *Handler) ListSolutions(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sols, appErr := h.svc.ListSolutions(c.Request.Context(), userID, limit, offset)
	if appErr != nil {
		response.Error(c, appErr)
		return
	}
	response.OK(c, sols)
}
