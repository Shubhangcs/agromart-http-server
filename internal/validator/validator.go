package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// Validate validates the given struct using struct field tags.
// It returns a human-readable error string on failure, or nil on success.
func Validate(v any) error {
	if err := validate.Struct(v); err != nil {
		var errs validator.ValidationErrors
		if ok := AsValidationErrors(err, &errs); ok {
			messages := make([]string, 0, len(errs))
			for _, fe := range errs {
				messages = append(messages, fieldError(fe))
			}
			return fmt.Errorf("%s", strings.Join(messages, "; "))
		}
		return err
	}
	return nil
}

// AsValidationErrors unwraps a validator.ValidationErrors from err.
func AsValidationErrors(err error, target *validator.ValidationErrors) bool {
	if ve, ok := err.(validator.ValidationErrors); ok {
		*target = ve
		return true
	}
	return false
}

func fieldError(fe validator.FieldError) string {
	field := toSnakeCase(fe.Field())
	switch fe.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at least %s", field, fe.Param())
	case "max":
		if fe.Kind().String() == "string" {
			return fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
		}
		return fmt.Sprintf("%s must be at most %s", field, fe.Param())
	case "numeric":
		return field + " must contain only digits"
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, strings.ReplaceAll(fe.Param(), " ", ", "))
	default:
		return fmt.Sprintf("%s is invalid (%s)", field, fe.Tag())
	}
}

// toSnakeCase converts a PascalCase or camelCase field name to snake_case.
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r | 0x20) // toLower
	}
	return result.String()
}
