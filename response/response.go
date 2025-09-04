package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response defines the structure of a standard API response.
type Response struct {
	Code      int         `json:"code"`       // Error code: 0 means success, other values indicate failure.
	Message   string      `json:"message"`    // Prompt or error message.
	RequestID string      `json:"request_id"` // Unique request ID for tracing.
	Data      interface{} `json:"data"`       // Response data payload.
}

// Common response codes.
const (
	SuccessCode = 0 // Success status code.
	ErrorCode   = 1 // Error status code.
	WarnCode    = 2 // Warning status code.
)

// Option defines optional fields for customizing API responses.
type Option struct {
	HTTPStatus int         // HTTP status code.
	Code       int         // Business error code.
	Message    string      // Custom message.
	Data       interface{} // Data payload.
}

// getRequestID retrieves the request ID from gin.Context.
// Returns an empty string if not found.
func getRequestID(c *gin.Context) string {
	if id, ok := c.Get("request_id"); ok {
		if str, ok := id.(string); ok {
			return str
		}
	}
	return ""
}

// JSON sends a standardized API response as JSON.
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

// Success returns a success response with a custom message and data.
func Success(c *gin.Context, msg string, data interface{}) {
	JSON(c, Option{
		Code:       SuccessCode,
		Message:    msg,
		Data:       data,
		HTTPStatus: http.StatusOK,
	})
}

// Warn returns a warning response with a custom message and data.
func Warn(c *gin.Context, msg string, data interface{}) {
	JSON(c, Option{
		Code:       WarnCode,
		Message:    msg,
		Data:       data,
		HTTPStatus: http.StatusOK,
	})
}

// Error returns an error response with a custom message, data, and HTTP status code.
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
