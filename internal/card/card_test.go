package card

import "testing"

// stubCard is a minimal Card implementation for exercising TurnState helpers. Tests only care
// about identity; all the Card interface methods return zero values.
type stubCard struct{ name string }

func (c stubCard) ID() ID                     { return Invalid }
func (c stubCard) Name() string               { return c.name }
func (stubCard) Cost(*TurnState) int          { return 0 }
func (stubCard) Pitch() int                   { return 0 }
func (stubCard) Attack() int                  { return 0 }
func (stubCard) Defense() int                 { return 0 }
func (stubCard) Types() TypeSet               { return 0 }
func (stubCard) GoAgain() bool                { return false }
func (stubCard) Play(*TurnState, *CardState) int { return 0 }

// TestDrawOne_AppendsTopAndAdvancesDeck: DrawOne moves the top card from Deck into Drawn and
// preserves draw order for the caller.
func TestDrawOne_AppendsTopAndAdvancesDeck(t *testing.T) {
	a, b, c := stubCard{"a"}, stubCard{"b"}, stubCard{"c"}
	s := &TurnState{Deck: []Card{a, b, c}}

	s.DrawOne()
	if got := len(s.Deck); got != 2 {
		t.Fatalf("Deck len = %d, want 2", got)
	}
	if s.Deck[0] != b {
		t.Errorf("Deck[0] = %v, want b (top advanced past a)", s.Deck[0])
	}
	if len(s.Drawn) != 1 || s.Drawn[0] != a {
		t.Errorf("Drawn = %v, want [a]", s.Drawn)
	}

	s.DrawOne()
	if len(s.Drawn) != 2 || s.Drawn[1] != b {
		t.Errorf("Drawn after second draw = %v, want [a, b]", s.Drawn)
	}
}

// TestDrawOne_EmptyDeckIsNoOp: with an empty deck the helper returns silently; Drawn stays
// nil so callers don't see a spurious zero-value card.
func TestDrawOne_EmptyDeckIsNoOp(t *testing.T) {
	s := &TurnState{}
	s.DrawOne()
	if len(s.Drawn) != 0 {
		t.Errorf("Drawn = %v, want empty on no-deck draw", s.Drawn)
	}
	if s.Deck != nil {
		t.Errorf("Deck = %v, want nil", s.Deck)
	}
}
