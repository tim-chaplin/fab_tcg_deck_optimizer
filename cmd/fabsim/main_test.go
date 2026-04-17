package main

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestFabraryPathFor pins the sibling-path derivation. A .json extension is replaced; anything
// else gets .txt appended instead of clobbered so the JSON can't accidentally be overwritten.
func TestFabraryPathFor(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"mydecks/best_deck.json", "mydecks/best_deck.txt"},
		{"best.json", "best.txt"},
		{"deck", "deck.txt"},
		{"deck.data", "deck.data.txt"},
		{"with.dots.json", "with.dots.txt"},
	}
	for _, c := range cases {
		if got := fabraryPathFor(c.in); got != c.want {
			t.Errorf("fabraryPathFor(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestResolveOutPath pins the -deck / -out precedence. -deck wins when set alone; -out (or its
// default) passes through unchanged; mixing -deck with a non-default -out is a hard error so the
// user never has to guess which flag took effect.
func TestResolveOutPath(t *testing.T) {
	viseraiPath := filepath.Join("mydecks", "viserai-v2.json")

	cases := []struct {
		name    string
		out     string
		deck    string
		want    string
		wantErr string // substring match; "" means no error
	}{
		{"defaults", defaultOutPath, "", defaultOutPath, ""},
		{"custom -out, no -deck", "experiments/v1.json", "", "experiments/v1.json", ""},
		{"-deck alone", defaultOutPath, "viserai-v2", viseraiPath, ""},
		{"-deck strips .json suffix", defaultOutPath, "viserai-v2.json", viseraiPath, ""},
		{"conflict: -deck + non-default -out", "experiments/v1.json", "viserai-v2", "", "mutually exclusive"},
		{"invalid deck name", defaultOutPath, "../escape", "", "invalid character"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := resolveOutPath(c.out, c.deck)
			if c.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), c.wantErr) {
					t.Fatalf("got err=%v, want containing %q", err, c.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}
