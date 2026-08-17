// Package vcs contains functions for reading the running program's build information.
package vcs

import "runtime/debug"

// Version returns the build's module version for the main package.
func Version() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		return bi.Main.Version
	}

	return ""
}
