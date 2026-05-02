package sim_test

import (
	"io"
	"os"
	"strings"
	"testing"

	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// captureStdout redirects os.Stdout into a pipe for the duration of fn and returns whatever
// fn wrote. Restores the original stdout on return.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	done := make(chan string, 1)
	go func() {
		b, _ := io.ReadAll(r)
		done <- string(b)
	}()

	fn()
	w.Close()
	return <-done
}

// Tests that OptDebug=true makes Opt print a one-line summary of the outcome to stdout,
// and OptDebug=false stays quiet.
func TestOptDebug_PrintsOnlyWhenSet(t *testing.T) {
	a := testutils.NewStubCard("a")
	b := testutils.NewStubCard("b")
	prev := OptDebug
	defer func() { OptDebug = prev }()

	withOptHero(t, testutils.Hero{
		OptStrategy: func(cards []Card) (top, bottom []Card) {
			return []Card{cards[1]}, []Card{cards[0]} // swap: bottom a, keep b on top
		},
	}, func() {
		// Off by default: no output.
		OptDebug = false
		out := captureStdout(t, func() {
			s := NewTurnState([]Card{a, b}, nil)
			s.Opt(2)
		})
		if out != "" {
			t.Errorf("OptDebug=false produced stdout: %q", out)
		}

		// On: a single line naming inputs, top, and bottom.
		OptDebug = true
		out = captureStdout(t, func() {
			s := NewTurnState([]Card{a, b}, nil)
			s.Opt(2)
		})
		if !strings.Contains(out, "Opt(2)") || !strings.Contains(out, "top=") || !strings.Contains(out, "bottom=") {
			t.Errorf("OptDebug=true output missing expected fragments: %q", out)
		}
	})
}
