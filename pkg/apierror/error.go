package apierror

import (
	"fmt"
	"net/http"
)

// ErrorCode 表示API错误码
type ErrorCode string

// 预定义错误码
const (
	ErrorCodeBadRequest          ErrorCode = "BAD_REQUEST"
	ErrorCodeUnauthorized        ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden           ErrorCode = "FORBIDDEN"
	ErrorCodeNotFound            ErrorCode = "NOT_FOUND"
	ErrorCodeInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrorCodeValidationFailed    ErrorCode = "VALIDATION_FAILED"
	ErrorCodeDuplicateEntity     ErrorCode = "DUPLICATE_ENTITY"
	ErrorCodeWeb3SignatureError  ErrorCode = "WEB3_SIGNATURE_ERROR"
)

// APIError 表示API错误
type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
	Status  int       `json:"-"` // HTTP状态码，不返回给客户端
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewBadRequestError 创建400错误
func NewBadRequestError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeBadRequest,
		Message: message,
		Detail:  detail,
		Status:  http.StatusBadRequest,
	}
}

// NewUnauthorizedError 创建401错误
func NewUnauthorizedError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeUnauthorized,
		Message: message,
		Detail:  detail,
		Status:  http.StatusUnauthorized,
	}
}

// NewForbiddenError 创建403错误
func NewForbiddenError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeForbidden,
		Message: message,
		Detail:  detail,
		Status:  http.StatusForbidden,
	}
}

// NewNotFoundError 创建404错误
func NewNotFoundError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeNotFound,
		Message: message,
		Detail:  detail,
		Status:  http.StatusNotFound,
	}
}

// NewInternalServerError 创建500错误
func NewInternalServerError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeInternalServerError,
		Message: message,
		Detail:  detail,
		Status:  http.StatusInternalServerError,
	}
}

// NewValidationError 创建校验错误
func NewValidationError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeValidationFailed,
		Message: message,
		Detail:  detail,
		Status:  http.StatusBadRequest,
	}
}

// NewDuplicateEntityError 创建重复实体错误
func NewDuplicateEntityError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeDuplicateEntity,
		Message: message,
		Detail:  detail,
		Status:  http.StatusConflict,
	}
}

// NewWeb3SignatureError 创建Web3签名错误
func NewWeb3SignatureError(message string, detail string) *APIError {
	return &APIError{
		Code:    ErrorCodeWeb3SignatureError,
		Message: message,
		Detail:  detail,
		Status:  http.StatusBadRequest,
	}
}
