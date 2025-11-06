package main

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// CustomValidator holds the validator instance
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator initializes the validator and registers custom rules
func newValidator() *CustomValidator {
	v := validator.New()

	// Register custom validation functions
	_ = v.RegisterValidation("phoneIN", validatePhoneIN)
	_ = v.RegisterValidation("passwordStrong", validatePasswordStrong)

	return &CustomValidator{validator: v}
}

// Validate implements the Echo Validator interface or can be used manually
func (cv *CustomValidator) Validate(i any) error {
	return cv.validator.Struct(i)
}

// ------------------ Custom Validation Functions ------------------

var (
	rePhoneIN = regexp.MustCompile(`^[6-9]\d{9}$`)
)

// validatePhoneIN checks if the phone number is a valid Indian mobile number
func validatePhoneIN(fl validator.FieldLevel) bool {
	return rePhoneIN.MatchString(fl.Field().String())
}

// validatePasswordStrong: ≥8 chars, at least one lower, one upper, one digit
func validatePasswordStrong(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	if len([]rune(s)) < 8 { // rune-length in case of unicode input
		return false
	}

	var hasLower, hasUpper, hasDigit bool
	for _, r := range s {
		switch {
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
		if hasLower && hasUpper && hasDigit {
			return true
		}
	}
	return false
}
