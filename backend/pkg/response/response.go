package response

import (
	"github.com/c-learn/pkg/errs"
	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, code int, data any) {
	c.JSON(code, data)
}

func Error(c *gin.Context, err *errs.AppError) {
	c.JSON(err.Code, gin.H{
		"error":   err.Error,
		"message": err.Message,
		"details": err.Details,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(201, data)
}

func OK(c *gin.Context, data any) {
	c.JSON(200, data)
}

func NoContent(c *gin.Context) {
	c.Status(204)
}
