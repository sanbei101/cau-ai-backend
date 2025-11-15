package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func Success[T any](c *gin.Context, data T) {
	c.JSON(http.StatusOK, Response[T]{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func Fail(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response[any]{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func SuccessWithPage[T any](c *gin.Context, list []T, total int64, page, pageSize int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	resp := PageResponse[T]{
		List:     list,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	Success(c, resp)
}
