package runner

import (
	"github.com/c-learn/internal/models"
	"github.com/c-learn/pkg/errs"
	"github.com/c-learn/pkg/response"
	"github.com/c-learn/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Run(c *gin.Context) {
	var req models.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, 400, gin.H{"error": "validation_error", "message": "invalid request body"})
		return
	}

	if err := validator.ValidateSourceCode(req.SourceCode); err != nil {
		response.Error(c, errs.NewValidation(err.Error(), nil))
		return
	}

	result, err := h.svc.Run(c.Request.Context(), req)
	if err != nil {
		response.Error(c, errs.NewInternal(err.Error()))
		return
	}
	response.OK(c, result)
}
