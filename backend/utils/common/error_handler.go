package common

import (
	"github.com/gin-gonic/gin"
)

// APIError represents a standardized API error
type APIError struct {
	Message string      `json:"message"`
	Type    string      `json:"type"`
	Param   interface{} `json:"param,omitempty"`
	Code    interface{} `json:"code,omitempty"`
}

// APIErrorResponse represents the standardized API error response
type APIErrorResponse struct {
	Error APIError `json:"error"`
}

// NewAPIError creates a new APIError instance
func NewAPIError(message, errType string, param, code interface{}) *APIError {
	return &APIError{
		Message: message,
		Type:    errType,
		Param:   param,
		Code:    code,
	}
}

// SendAPIError sends a standardized API error response
func SendAPIError(c *gin.Context, statusCode int, message, errType string) {
	c.JSON(statusCode, APIErrorResponse{
		Error: APIError{
			Message: message,
			Type:    errType,
			Param:   nil,
			Code:    nil,
		},
	})
}

// SendAPIErrorWithDetails sends a standardized API error response with additional details
func SendAPIErrorWithDetails(c *gin.Context, statusCode int, err *APIError) {
	c.JSON(statusCode, APIErrorResponse{
		Error: *err,
	})
}

// Error constants for common error types
const (
	ErrTypeInvalidRequest  = "invalid_request_error"
	ErrTypeNotFound        = "not_found_error"
	ErrTypeBadRequest      = "bad_request_error"
	ErrTypeInternalServer  = "internal_server_error"
	ErrTypeUnauthorized    = "unauthorized_error"
	ErrTypeForbidden       = "forbidden_error"
	ErrTypeConflict        = "conflict_error"
	ErrTypeValidation      = "validation_error"
	ErrTypeForward         = "forward_error"
	ErrTypeInvalidProvider = "invalid_provider_error"
	ErrTypeAccountNotFound = "account_not_found_error"
	ErrTypeModelNotFound   = "model_not_found_error"
)

// Common error messages
const (
	ErrMsgInvalidID          = "Invalid ID parameter"
	ErrMsgAccountNotFound    = "Account not found"
	ErrMsgModelNotFound      = "Model not found"
	ErrMsgFailedToReadBody   = "Failed to read request body"
	ErrMsgFailedToParseBody  = "Failed to parse request body"
	ErrMsgModelMissing       = "model '' is missing"
	ErrMsgInvalidProvider    = "Invalid provider"
	ErrMsgFailedToDelete     = "Failed to delete resource"
	ErrMsgFailedToToggle     = "Failed to toggle resource"
	ErrMsgFailedToUpdate     = "Failed to retrieve updated resource"
	ErrMsgAllAttemptsFailed  = "All account attempts failed"
	ErrMsgNoAccountsFound    = "No accounts found for model '%s'"
	ErrMsgNoModelsFound      = "No actual models found for model '%s'"
	ErrMsgFailedToUpdateBody = "Failed to update request body"
)
