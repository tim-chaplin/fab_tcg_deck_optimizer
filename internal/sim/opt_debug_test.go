package sim

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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

// Tests that gameengine.OptDebug=true makes Opt print a one-line summary of the outcome to stdout,
// and gameengine.OptDebug=false stays quiet.
func TestOptDebug_PrintsOnlyWhenSet(t *testing.T) {
	a := NewFakeCard("a")
	b := NewFakeCard("b")
	prev := gameengine.OptDebug
	defer func() { gameengine.OptDebug = prev }()

	withOptHero(t, FakeHero{
		OptStrategy: func(cards []card.Card) (top, bottom []card.Card) {
			return []card.Card{cards[1]}, []card.Card{cards[0]} // swap: bottom a, keep b on top
		},
	}, func() {
		// Off by default: no output.
		gameengine.OptDebug = false
		out := captureStdout(t, func() {
			s := gameengine.NewFromCards([]card.Card{a, b}, nil)
			s.Opt(s.Logger(), 2)
		})
		if out != "" {
			t.Errorf("gameengine.OptDebug=false produced stdout: %q", out)
		}

		// On: a single line naming inputs, top, and bottom.
		gameengine.OptDebug = true
		out = captureStdout(t, func() {
			s := gameengine.NewFromCards([]card.Card{a, b}, nil)
			s.Opt(s.Logger(), 2)
		})
		if !strings.Contains(out, "Opt(2)") || !strings.Contains(out, "top=") || !strings.Contains(out, "bottom=") {
			t.Errorf("gameengine.OptDebug=true output missing expected fragments: %q", out)
		}
	})
}
