package main

import (
	"regexp"

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
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// ------------------ Custom Validation Functions ------------------

var (
	rePhoneIN        = regexp.MustCompile(`^[6-9]\d{9}$`)
	rePasswordStrong = regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d).{8,}$`)
)

// validatePhoneIN checks if the phone number is a valid Indian mobile number
func validatePhoneIN(fl validator.FieldLevel) bool {
	return rePhoneIN.MatchString(fl.Field().String())
}

// validatePasswordStrong ensures password has upper, lower, and numeric characters
func validatePasswordStrong(fl validator.FieldLevel) bool {
	return rePasswordStrong.MatchString(fl.Field().String())
}
