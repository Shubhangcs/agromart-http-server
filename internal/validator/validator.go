package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	rePhone        = regexp.MustCompile(`^\d{7,15}$`)
	reAadhaar      = regexp.MustCompile(`^\d{12}$`)
	rePAN          = regexp.MustCompile(`^[A-Z]{5}[0-9]{4}[A-Z]$`)
	reExportImport = regexp.MustCompile(`^[A-Z0-9]{10}$`)
	reMSME         = regexp.MustCompile(`^UDYAM-[A-Z]{2}-\d{2}-\d{7}$`)
	reGST          = regexp.MustCompile(`^\d{2}[A-Z]{5}\d{4}[A-Z][A-Z\d]Z[A-Z\d]$`)
	reFassi        = regexp.MustCompile(`^\d{14}$`)
)

var validate = validator.New()

func init() {
	validate.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		return rePhone.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("aadhaar", func(fl validator.FieldLevel) bool {
		return reAadhaar.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("pan", func(fl validator.FieldLevel) bool {
		return rePAN.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("export_import", func(fl validator.FieldLevel) bool {
		return reExportImport.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("msme", func(fl validator.FieldLevel) bool {
		return reMSME.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("gst", func(fl validator.FieldLevel) bool {
		return reGST.MatchString(fl.Field().String())
	})
	validate.RegisterValidation("fassi", func(fl validator.FieldLevel) bool {
		return reFassi.MatchString(fl.Field().String())
	})
}

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
	case "phone":
		return field + " must be a valid phone number (7-15 digits)"
	case "aadhaar":
		return field + " must be a valid 12-digit Aadhaar number"
	case "pan":
		return field + " must be a valid PAN (e.g. ABCDE1234F)"
	case "export_import":
		return field + " must be a valid 10-character IEC code"
	case "msme":
		return field + " must be a valid UDYAM number (e.g. UDYAM-MH-01-0000001)"
	case "gst":
		return field + " must be a valid 15-character GST number"
	case "fassi":
		return field + " must be a valid 14-digit FSSAI number"
	default:
		return field + " is invalid"
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r | 0x20)
	}
	return result.String()
}
