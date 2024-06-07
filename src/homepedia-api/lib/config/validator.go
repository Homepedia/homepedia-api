package config

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	v := validator.New()
	v.RegisterValidation("containsUppercase", ContainsAtLeastOneUppercase)
	v.RegisterValidation("containsSpecialCharacter", ContainsAtLeastOneSpecialCharacter)
	return &CustomValidator{validator: v}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func ContainsAtLeastOneUppercase(fl validator.FieldLevel) bool {
	for _, r := range fl.Field().String() {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

func ContainsAtLeastOneSpecialCharacter(fl validator.FieldLevel) bool {
	specialChars := "!@#$%^&*()-_=+[]{}|;:'\",.<>?/`~\\"
	for _, r := range fl.Field().String() {
		if strings.ContainsRune(specialChars, r) {
			return true
		}
	}
	return false
}
