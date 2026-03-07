package response

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:  code,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, errs interface{}) {
	c.JSON(code, APIResponse{
		Status:  code,
		Message: message,
		Errors:  errs,
	})
}
