package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

// FieldViolation represents an error on a specific field
type FieldViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

// ValidationDetails represents the validation error response structure
type ValidationDetails struct {
	Violations []FieldViolation `json:"violations"`
}

// ValidationErrorResponse represents the final structured error response
type ValidationErrorResponse struct {
	Code    string              `json:"code"`
	Message string              `json:"message"`
	Details []ValidationDetails `json:"details"`
}

// Custom error messages for validation tags
var validationMessages = map[string]string{
	"required": "is required",
	"min":      "is too short",
	"max":      "is too long",
	"alphanum": "must contain only letters and numbers",
	"email":    "must be a valid email address",
	"gte":      "must be greater than or equal to a specific value",
	"lte":      "must be less than or equal to a specific value",
}

// Convert validation errors into a structured response
func FormatValidationErrors(err error) ValidationErrorResponse {
	var violations []FieldViolation

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			fieldName := strings.ToLower(fe.Field()) // Normalize field name
			message, exists := validationMessages[fe.Tag()]
			if !exists {
				message = fmt.Sprintf("failed validation on '%s'", fe.Tag())
			}
			violations = append(violations, FieldViolation{
				Field:       fieldName,
				Description: message,
			})
		}
	} else {
		violations = append(violations, FieldViolation{
			Field:       "general",
			Description: err.Error(),
		})
	}

	return ValidationErrorResponse{
		Code:    "INVALID_ARGUMENT",
		Message: "invalid parameters",
		Details: []ValidationDetails{
			{
				Violations: violations,
			},
		},
	}
}

// CreateFieldViolation to standardize field-level error responses
func CreateFieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}
