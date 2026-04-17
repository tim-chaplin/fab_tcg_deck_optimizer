package mydecks

import (
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	// filepath.Join produces platform-native separators, so expected values go through Join too —
	// that keeps the test portable between Windows and Unix CI.
	cases := []struct {
		in, want string
	}{
		{"viserai-v2", filepath.Join(Dir, "viserai-v2.json")},
		{"viserai-v2.json", filepath.Join(Dir, "viserai-v2.json")},
		{"my best deck", filepath.Join(Dir, "my best deck.json")},
	}
	for _, c := range cases {
		got, err := Path(c.in)
		if err != nil {
			t.Errorf("Path(%q) error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("Path(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestValidateName(t *testing.T) {
	ok := []string{"viserai-v2", "my best deck", "deck_42", "Viserai"}
	bad := []string{"", ".", "..", "../escape", "path/with/slash", "back\\slash", "has*star", "has?mark"}

	for _, name := range ok {
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q) = %v, want nil", name, err)
		}
	}
	for _, name := range bad {
		if err := ValidateName(name); err == nil {
			t.Errorf("ValidateName(%q) = nil, want error", name)
		}
	}
}
