package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 定义返回结构
type Response struct {
	Code      int         `json:"code"`       // 错误码：0 表示成功，其他表示失败
	Message   string      `json:"message"`    // 提示信息
	RequestID string      `json:"request_id"` // 请求的 ID，便于追踪
	Data      interface{} `json:"data"`       // 返回数据
}

// Codes 定义常见状态码
const (
	SuccessCode = 0 // 成功的状态码
	ErrorCode   = 1 // 失败的状态码
	WarnCode    = 2 // 警告的状态码
)

// Option 响应可选项
type Option struct {
	HTTPStatus int
	Code       int
	Message    string
	Data       interface{}
}

// getRequestID 从 gin.Context 获取 request_id，如果不存在返回空字符串
func getRequestID(c *gin.Context) string {
	if id, ok := c.Get("request_id"); ok {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

// JSON 响应统一出口
func JSON(c *gin.Context, opts Option) {
	if opts.Message == "" {
		switch opts.Code {
		case SuccessCode:
			opts.Message = "success"
		case WarnCode:
			opts.Message = "warning"
		case ErrorCode:
			opts.Message = "error"
		default:
			opts.Message = "info"
		}
	}
	if opts.HTTPStatus == 0 {
		opts.HTTPStatus = http.StatusOK
	}

	c.JSON(opts.HTTPStatus, Response{
		Code:      opts.Code,
		Message:   opts.Message,
		RequestID: getRequestID(c),
		Data:      opts.Data,
	})
}

// Success 返回成功响应
func Success(c *gin.Context, msg string, data interface{}) {
	JSON(c, Option{
		Code:       SuccessCode,
		Message:    msg,
		Data:       data,
		HTTPStatus: http.StatusOK,
	})
}

// Warn 返回警告响应
func Warn(c *gin.Context, msg string, data interface{}) {
	JSON(c, Option{
		Code:       WarnCode,
		Message:    msg,
		Data:       data,
		HTTPStatus: http.StatusOK,
	})
}

// Error 返回错误响应，可以自定义 http 状态码
func Error(c *gin.Context, msg string, data interface{}, httpStatus ...int) {
	status := http.StatusInternalServerError
	if len(httpStatus) > 0 && httpStatus[0] > 0 {
		status = httpStatus[0]
	}
	JSON(c, Option{
		Code:       ErrorCode,
		Message:    msg,
		Data:       data,
		HTTPStatus: status,
	})
}
