package validator

import "testing"

// TestValidatorFlow verifies the functionality of the validator.
func TestValidatorFlow(t *testing.T) {
	v := New()

	if !v.Valid() {
		t.Fatal("new validator should start valid")
	}

	v.Check(false, "name", "must be provided")
	if v.Valid() {
		t.Fatal("validator should be invalid after failed check")
	}

	if got := v.Errors["name"]; got != "must be provided" {
		t.Fatalf("unexpected error message: got %q", got)
	}

	v.AddError("name", "should not overwrite")
	if got := v.Errors["name"]; got != "must be provided" {
		t.Fatalf("error message was overwritten: got %q", got)
	}
}
