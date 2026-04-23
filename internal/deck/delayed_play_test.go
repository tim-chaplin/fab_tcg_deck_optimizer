package deck

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
)

// stubDelayed implements card.DelayedPlay and records each PlayNextTurn call so tests can
// assert the queue was processed exactly once per turn boundary. Kept around until Phase 4
// deletes the legacy DelayedPlay mechanism; real cards have all migrated to AuraTrigger.
type stubDelayed struct {
	damage int
	calls  *int
}

func (s stubDelayed) ID() card.ID              { return card.Invalid }
func (s stubDelayed) Name() string             { return "StubDelayed" }
func (s stubDelayed) Cost(*card.TurnState) int { return 0 }
func (s stubDelayed) Pitch() int               { return 1 }
func (s stubDelayed) Attack() int              { return 0 }
func (s stubDelayed) Defense() int             { return 0 }
func (s stubDelayed) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (s stubDelayed) GoAgain() bool                              { return true }
func (s stubDelayed) Play(*card.TurnState, *card.CardState) int  { return 0 }
func (s stubDelayed) PlayNextTurn(*card.TurnState) card.DelayedPlayResult {
	*s.calls++
	return card.DelayedPlayResult{Damage: s.damage}
}

// TestCollectDelayedPlays_OnlyAttackRole verifies only Role==Attack entries are queued:
// pitched / defended / held / arsenaled copies don't have their aura in the arena so their
// next-turn trigger shouldn't fire.
func TestCollectDelayedPlays_OnlyAttackRole(t *testing.T) {
	var calls int
	d := stubDelayed{damage: 1, calls: &calls}
	plain := fake.RedAttack{}
	line := []hand.CardAssignment{
		{Card: d, Role: hand.Attack},
		{Card: d, Role: hand.Pitch},
		{Card: d, Role: hand.Held},
		{Card: d, Role: hand.Arsenal},
		{Card: plain, Role: hand.Attack},
	}
	got := collectDelayedPlays(line, nil)
	if len(got) != 1 {
		t.Fatalf("len(queued) = %d, want 1 (only Role=Attack DelayedPlay qualifies)", len(got))
	}
	if got[0] != d {
		t.Errorf("queued[0] = %v, want %v", got[0], d)
	}
}

// TestRunDelayedPlays_FiresEachQueuedCardOnce verifies every queued card's PlayNextTurn is
// invoked exactly once per pass and the contributions/total are reported.
func TestRunDelayedPlays_FiresEachQueuedCardOnce(t *testing.T) {
	var callsA, callsB int
	a := stubDelayed{damage: 2, calls: &callsA}
	b := stubDelayed{damage: 3, calls: &callsB}
	contribs, total, _, revealed := runDelayedPlays([]card.Card{a, b}, nil)
	if total != 5 {
		t.Errorf("total = %d, want 5 (2+3)", total)
	}
	if len(contribs) != 2 {
		t.Fatalf("len(contribs) = %d, want 2", len(contribs))
	}
	if contribs[0].Damage != 2 {
		t.Errorf("contribs[0].Damage = %d, want 2", contribs[0].Damage)
	}
	if contribs[1].Damage != 3 {
		t.Errorf("contribs[1].Damage = %d, want 3", contribs[1].Damage)
	}
	if len(revealed) != 0 {
		t.Errorf("revealed = %v, want nil (damage-only stubs don't reveal)", revealed)
	}
	if callsA != 1 || callsB != 1 {
		t.Errorf("PlayNextTurn call counts = (%d, %d), want (1, 1)", callsA, callsB)
	}
}

// TestRunDelayedPlays_EmptyQueue short-circuits: no contribs, no allocation, zero total.
func TestRunDelayedPlays_EmptyQueue(t *testing.T) {
	contribs, total, _, revealed := runDelayedPlays(nil, nil)
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if contribs != nil {
		t.Errorf("contribs = %v, want nil", contribs)
	}
	if revealed != nil {
		t.Errorf("revealed = %v, want nil", revealed)
	}
}
