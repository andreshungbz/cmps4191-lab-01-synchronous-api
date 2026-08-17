package validator

import (
	"regexp"
	"slices"
)

var (
	// https://html.spec.whatwg.org/#valid-e-mail-address
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

// Validator holds multiple errors for validating JSON values to enforce
// business rules. These errors are typically returned to the client as a
// JSON response.
type Validator struct {
	Errors map[string]string
}

// New is a factory function returning a new Validator.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid checks if there are any errors recorded in the Validator.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError records an error key and message to the Validator.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check evaluates a boolean expression, recording an error in Validator
// if it evaluates to false.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValue checks if a value is contained in a list of permitted
// values.
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches checks if a given string conforms to a regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique checks a given slice to see if all values in it are unique.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
