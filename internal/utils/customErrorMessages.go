package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

func CustomErrorMessages(err error) string {
	var errorMsg string
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		return "Invalid input"
	}

	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		switch tag {
		case "required":
			errorMsg = fmt.Sprintf("Field %s is required", field)
		case "min":
			errorMsg = fmt.Sprintf("Field %s min length %s", field, err.Param())
		case "max":
			errorMsg = fmt.Sprintf("Field %s max length %s", field, err.Param())
		case "email":
			errorMsg = fmt.Sprintf("Field %s must be a valid email", field)
		case "at least one field must be present":
			errorMsg = "At least one field is required"
		default:
			errorMsg = fmt.Sprintf("Field %s is invalid", field)
		}
	}
	return errorMsg
}
