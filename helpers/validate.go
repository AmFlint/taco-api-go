package helpers

import validator2 "gopkg.in/validator.v2"

// Validate - Validate a data structure against given scenario
func Validate(v interface{}, tag string) error {
	validator := validator2.NewValidator()

	// If validation tag is provided
	if len(tag) > 0 {
		validator.SetTag(tag)
	}
	return validator.Validate(v)
}