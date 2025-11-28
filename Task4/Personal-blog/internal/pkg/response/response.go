package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code    int         `json:"code"`    // 状态码（200=成功，非200=失败）
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"data"`    // 响应数据（可选）
}

// Success 成功响应
func Success(c *gin.Context, data interface{}, message string) {
	if message == "" {
		message = "操作成功"
	}
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, message string) {
	if message == "" {
		message = "操作失败"
	}
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// 常用错误响应快捷方法
func BadRequest(c *gin.Context, message string) {
	Fail(c, http.StatusBadRequest, message)
}

func InternalError(c *gin.Context, message string) {
	Fail(c, http.StatusInternalServerError, message)
}
