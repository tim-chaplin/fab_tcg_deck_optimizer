package textio

import (
	"path/filepath"
	"testing"
)

func TestMydecksPath(t *testing.T) {
	// filepath.Join produces platform-native separators, so expected values go through Join too —
	// that keeps the test portable between Windows and Unix CI.
	cases := []struct {
		in, want string
	}{
		{"viserai-v2", filepath.Join(MydecksDir, "viserai-v2.json")},
		{"viserai-v2.json", filepath.Join(MydecksDir, "viserai-v2.json")},
		{"my best deck", filepath.Join(MydecksDir, "my best deck.json")},
	}
	for _, c := range cases {
		got, err := MydecksPath(c.in)
		if err != nil {
			t.Errorf("MydecksMydecksPath(%q) error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("MydecksMydecksPath(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestValidateMydecksName(t *testing.T) {
	ok := []string{"viserai-v2", "my best deck", "deck_42", "Viserai"}
	bad := []string{"", ".", "..", "../escape", "path/with/slash", "back\\slash", "has*star", "has?mark"}

	for _, name := range ok {
		if err := ValidateMydecksName(name); err != nil {
			t.Errorf("ValidateMydecksName(%q) = %v, want nil", name, err)
		}
	}
	for _, name := range bad {
		if err := ValidateMydecksName(name); err == nil {
			t.Errorf("ValidateMydecksName(%q) = nil, want error", name)
		}
	}
}
