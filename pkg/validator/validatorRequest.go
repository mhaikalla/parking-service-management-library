package validator

import (
	"regexp"

	validation "github.com/go-playground/validator/v10"
)

func NewValidator() validation.Validate {
	validators := *validation.New()
	validators.RegisterValidation("alpha_or_numeric", ValidateAlphaOrNumeric)
	return validators
}

func ValidateAlphaOrNumeric(fl validation.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(fl.Field().String())
}
