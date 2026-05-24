package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Tests that ge.HeldHand() returns only Held-role entries, skipping Pitch and Attack
// role cards that the partition has scheduled to commit to chain costs / play.
func TestHeldHand_FiltersOutPitchAndAttackRoles(t *testing.T) {
	ge := New()
	ge.SetHandStates([]card.CardState{
		{Card: fakeCard{id: 1}, Role: card.Pitch},
		{Card: fakeCard{id: 2}, Role: card.Held},
		{Card: fakeCard{id: 3}, Role: card.Attack},
		{Card: fakeCard{id: 4}, Role: card.Held},
	})
	got := ge.HeldHand()
	if len(got) != 2 {
		t.Fatalf("HeldHand len = %d, want 2 (entries 2 and 4)", len(got))
	}
	if got[0].ID() != 2 || got[1].ID() != 4 {
		t.Errorf("HeldHand IDs = [%d, %d], want [2, 4]", got[0].ID(), got[1].ID())
	}
	if h := ge.Hand(); len(h) != 4 {
		t.Errorf("Hand still surfaces ALL entries (incl. Pitch / Attack), got len = %d, want 4 — Demolition Crew + Spring Load rely on this", len(h))
	}
}

// Tests that ge.PopHandAt(i) indexes into the Held subset, NOT the raw role-tagged
// slice. Pre-fix, PopHandAt popped raw[i], so cards like Rise Above / the Emissary cycle
// could accidentally remove a Pitch-role card that the partition had scheduled to pay
// for a chain cost — leaving the cost unpaid downstream OR double-counting that card
// (once via the divert effect, once via the end-of-turn pitch recycle).
func TestPopHandAt_OnlyPopsHeldRole(t *testing.T) {
	ge := New()
	ge.SetHandStates([]card.CardState{
		{Card: fakeCard{id: 1}, Role: card.Pitch},
		{Card: fakeCard{id: 2}, Role: card.Held},
		{Card: fakeCard{id: 3}, Role: card.Attack},
		{Card: fakeCard{id: 4}, Role: card.Held},
	})
	// PopHandAt(0) should pop the FIRST Held entry (id=2), not the raw index 0 (id=1, Pitch).
	got := ge.PopHandAt(0)
	if got.ID() != 2 {
		t.Errorf("PopHandAt(0) = ID %d, want 2 (first Held entry, skipping Pitch at raw index 0)", got.ID())
	}
	// The Pitch entry must remain untouched.
	remaining := ge.GameState.HandStates()
	if len(remaining) != 3 {
		t.Fatalf("hand size = %d, want 3 (only Held entry removed)", len(remaining))
	}
	for _, s := range remaining {
		if s.Role == card.Held && s.Card.ID() == 2 {
			t.Errorf("Held card 2 still present after PopHandAt(0)")
		}
	}
	// The remaining Pitch entry must still carry Pitch role (not accidentally re-tagged).
	foundPitch := false
	for _, s := range remaining {
		if s.Card.ID() == 1 && s.Role == card.Pitch {
			foundPitch = true
		}
	}
	if !foundPitch {
		t.Errorf("Pitch entry with ID 1 missing or re-tagged: %+v", remaining)
	}

	// PopHandAt(0) again should pop ID 4 (the other Held). Pitch and Attack stay.
	got = ge.PopHandAt(0)
	if got.ID() != 4 {
		t.Errorf("PopHandAt(0) second call = ID %d, want 4", got.ID())
	}

	// PopHandAt(0) with no Held cards left must panic — callers gate on len(HeldHand()).
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("PopHandAt with no Held cards did not panic")
		}
	}()
	ge.PopHandAt(0)
}
