package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errorMessages []string

	for _, err := range err.(validator.ValidationErrors) {
		errorMessages = append(errorMessages, formatValidationError(err))
	}

	return fmt.Errorf("validation failed: %s", strings.Join(errorMessages, ", "))
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
