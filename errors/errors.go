package errors

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorType string

const (
	CompressionError ErrorType = "COMPRESSION_ERROR"
	EncryptionError  ErrorType = "ENCRYPTION_ERROR"
	NetworkError     ErrorType = "NETWORK_ERROR"
	ValidationError  ErrorType = "VALIDATION_ERROR"
	TimeoutError     ErrorType = "TIMEOUT_ERROR"
	InternalError    ErrorType = "INTERNAL_ERROR"
)

type SDKError struct {
	Type      ErrorType `json:"type"`
	Message   string    `json:"message"`
	Code      string    `json:"code"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Stack     string    `json:"stack,omitempty"`
	Cause     error     `json:"-"`
}

func (e *SDKError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *SDKError) Unwrap() error {
	return e.Cause
}

func NewError(errorType ErrorType, message string, code string) *SDKError {
	return &SDKError{
		Type:    errorType,
		Message: message,
		Code:    code,
		Stack:   getStackTrace(),
	}
}

func NewErrorWithCause(errorType ErrorType, message string, code string, cause error) *SDKError {
	return &SDKError{
		Type:    errorType,
		Message: message,
		Code:    code,
		Cause:   cause,
		Stack:   getStackTrace(),
	}
}

func (e *SDKError) WithDetails(details map[string]interface{}) *SDKError {
	e.Details = details
	return e
}

func (e *SDKError) WithDetail(key string, value interface{}) *SDKError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

func getStackTrace() string {
	buf := make([]byte, 1024)
	n := runtime.Stack(buf, false)
	stack := string(buf[:n])
	
	// Clean up the stack trace
	lines := strings.Split(stack, "\n")
	if len(lines) > 10 {
		lines = lines[:10] // Limit stack trace length
	}
	
	return strings.Join(lines, "\n")
}

// Error type constructors
func NewCompressionError(message string, cause error) *SDKError {
	return NewErrorWithCause(CompressionError, message, "COMP_001", cause)
}

func NewEncryptionError(message string, cause error) *SDKError {
	return NewErrorWithCause(EncryptionError, message, "ENC_001", cause)
}

func NewNetworkError(message string, cause error) *SDKError {
	return NewErrorWithCause(NetworkError, message, "NET_001", cause)
}

func NewValidationError(message string) *SDKError {
	return NewError(ValidationError, message, "VAL_001")
}

func NewTimeoutError(message string) *SDKError {
	return NewError(TimeoutError, message, "TIME_001")
}

func NewInternalError(message string, cause error) *SDKError {
	return NewErrorWithCause(InternalError, message, "INT_001", cause)
}
