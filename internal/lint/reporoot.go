package lint

import (
	"os"
	"path/filepath"
	"testing"
)

// RepoRoot walks up from the test's working directory to the directory holding go.mod.
// Repo-wide lint tests use this to anchor a filepath.WalkDir over the whole tree
// regardless of which package's test binary is running.
func RepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found walking up from the test directory")
		}
		dir = parent
	}
}
