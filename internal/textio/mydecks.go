package textio

// Path resolution for the local mydecks/ directory, where fabsim outputs JSON-persisted
// decks and fabrary imports both land. Centralises name validation (path-traversal,
// Windows-reserved characters) so every caller gets the same rules.

import (
	"fmt"
	"path/filepath"
	"strings"
)

// MydecksDir is the directory every resolved deck path is rooted under. Relative so
// callers behave the same regardless of the current working directory.
const MydecksDir = "mydecks"

// MydecksPath returns MydecksDir/<name>.json for a user-supplied deck name. A trailing
// ".json" on name is stripped before the join so users can type either "viserai-v2" or
// "viserai-v2.json".
//
// Returns an error if name contains path separators or any character that would escape
// MydecksDir, is empty, or is a reserved dot-name.
func MydecksPath(name string) (string, error) {
	name = strings.TrimSuffix(name, ".json")
	if err := ValidateMydecksName(name); err != nil {
		return "", err
	}
	return filepath.Join(MydecksDir, name+".json"), nil
}

// ValidateMydecksName rejects names that would produce an unusable or unsafe path under
// MydecksDir. Conservative by design: unusual names can be handled via an explicit -out
// path instead.
func ValidateMydecksName(name string) error {
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
