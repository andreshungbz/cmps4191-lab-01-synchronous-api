package vcs

import "testing"

// TestVersionSmoke verifies Version can be called without panicking.
func TestVersionSmoke(t *testing.T) {
	v := Version()

	if v == "" {
		t.Log("version returned an empty string (expected when build metadata is unavailable)")
	}
}
