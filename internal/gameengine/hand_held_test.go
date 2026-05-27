package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Tests that ge.HeldHand() returns only Held-role entries, skipping Pitch and Attack
// role cards that the partition has scheduled to commit to attack-step costs / play.
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

// Tests that Discard / MoveFromHandToTopOfDeck / MoveFromHandToBottomOfDeck pop only Held-role
// hand entries, skipping Pitch / Attack entries the partition has scheduled to commit
// to attack-step costs. A bug here would let a card like Rise Above or the Emissary cycle
// silently remove a Pitch-role card, leaving its cost unpaid downstream.
func TestDiscard_OnlyPopsHeldRole(t *testing.T) {
	ge := New()
	ge.SetHandStates([]card.CardState{
		{Card: fakeCard{id: 1}, Role: card.Pitch},
		{Card: fakeCard{id: 2}, Role: card.Held},
		{Card: fakeCard{id: 3}, Role: card.Attack},
		{Card: fakeCard{id: 4}, Role: card.Held},
	})
	if !ge.Discard("test") {
		t.Fatal("Discard returned false despite 2 Held entries present")
	}
	remaining := ge.GameState.HandStates()
	if len(remaining) != 3 {
		t.Fatalf("hand size = %d, want 3 (only Held entry removed)", len(remaining))
	}
	for _, s := range remaining {
		if s.Card.ID() == 2 {
			t.Errorf("Held card 2 still present after Discard — Pitch / Attack ahead of it should not have been targeted")
		}
	}
	foundPitch := false
	for _, s := range remaining {
		if s.Card.ID() == 1 && s.Role == card.Pitch {
			foundPitch = true
		}
	}
	if !foundPitch {
		t.Errorf("Pitch entry with ID 1 missing or re-tagged: %+v", remaining)
	}

	// Second discard should pop ID 4 (the other Held).
	if !ge.Discard("test") {
		t.Fatal("second Discard returned false despite 1 Held entry remaining")
	}
	// Third call should return false — no Held entries left.
	if ge.Discard("test") {
		t.Errorf("Discard with no Held cards returned true, want false")
	}
}
