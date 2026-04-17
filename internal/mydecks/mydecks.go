// Package mydecks resolves user-supplied deck names to paths under the local mydecks/ directory,
// where both fabsim's outputs and fabrary imports live. Centralising this means the name-validation
// rules (path-traversal, Windows-reserved characters) stay consistent across every command that
// takes a deck name from the user.
package mydecks

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Dir is the directory every resolved deck path is rooted under. Relative so commands behave the
// same regardless of where the user runs them from (matching fabsim's default -out behaviour).
const Dir = "mydecks"

// Path returns Dir/<name>.json for a user-supplied deck name. A trailing ".json" on `name` is
// stripped before the join so users can type either "viserai-v2" or "viserai-v2.json"
// interchangeably — whichever is in their muscle memory.
//
// Returns an error if `name` contains path separators or any character that would escape Dir, is
// empty, or is a reserved dot-name. Callers that need the raw-name validator directly can use
// ValidateName.
func Path(name string) (string, error) {
	name = strings.TrimSuffix(name, ".json")
	if err := ValidateName(name); err != nil {
		return "", err
	}
	return filepath.Join(Dir, name+".json"), nil
}

// ValidateName rejects names that would produce an unusable or unsafe path under Dir. Conservative
// by design — anything unusual the user actually wants can be passed as an explicit -out path.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("deck name is empty")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("deck name %q is reserved", name)
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return fmt.Errorf("deck name %q contains an invalid character (one of /\\:*?\"<>|)", name)
	}
	return nil
}
