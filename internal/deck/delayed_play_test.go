package deck

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
)

// stubDelayed implements card.DelayedPlay and records each PlayNextTurn call so tests can
// assert the queue was processed exactly once per turn boundary.
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
func (s stubDelayed) GoAgain() bool            { return true }
func (s stubDelayed) Play(*card.TurnState) int { return 0 }
func (s stubDelayed) PlayNextTurn(*card.TurnState) int {
	*s.calls++
	return s.damage
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
	contribs, total := runDelayedPlays([]card.Card{a, b}, nil)
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
	if callsA != 1 || callsB != 1 {
		t.Errorf("PlayNextTurn call counts = (%d, %d), want (1, 1)", callsA, callsB)
	}
}

// TestRunDelayedPlays_EmptyQueue short-circuits: no contribs, no allocation, zero total.
func TestRunDelayedPlays_EmptyQueue(t *testing.T) {
	contribs, total := runDelayedPlays(nil, nil)
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if contribs != nil {
		t.Errorf("contribs = %v, want nil", contribs)
	}
}

// TestRunDelayedPlays_PassesPostDrawDeck verifies the callback sees the deck slice it was
// given — important for Sigil of the Arknight, which peeks Deck[0] as the next-turn reveal.
func TestRunDelayedPlays_PassesPostDrawDeck(t *testing.T) {
	sigil := runeblade.SigilOfTheArknightBlue{}
	attackTop := []card.Card{runeblade.AetherSlashRed{}}
	_, total := runDelayedPlays([]card.Card{sigil}, attackTop)
	if total != card.DrawValue {
		t.Errorf("total = %d, want %d (top card is an attack action → DrawValue credited)", total, card.DrawValue)
	}
	// Sigil itself is an Aura, not an action-attack — reveal fails.
	_, total = runDelayedPlays([]card.Card{sigil}, []card.Card{sigil})
	if total != 0 {
		t.Errorf("total = %d, want 0 (top card is aura, not an attack action)", total)
	}
}

// TestEvaluate_DelayedFromLastTurnSurfacesInBest runs a full Evaluate with Sigil of the
// Arknight in the deck and asserts the PlayNextTurn callback lands a DelayedFromLastTurn
// entry on at least some hand's TurnSummary. Uses enough copies + runs that the shuffle
// reliably plays Sigil before the best turn is recorded.
func TestEvaluate_DelayedFromLastTurnSurfacesInBest(t *testing.T) {
	// Pack the deck: lots of Sigils so every hand plays one; Aether Slash as the reveal target
	// so PlayNextTurn credits DrawValue; filler Blue attacks to round out hand value. Minimum
	// deck size for Viserai (Int 4): 20.
	sigil := runeblade.SigilOfTheArknightBlue{}
	slash := runeblade.AetherSlashRed{}
	deckCards := make([]card.Card, 0, 20)
	for i := 0; i < 8; i++ {
		deckCards = append(deckCards, sigil)
	}
	for i := 0; i < 6; i++ {
		deckCards = append(deckCards, slash)
	}
	for i := 0; i < 6; i++ {
		deckCards = append(deckCards, fake.BlueAttack{})
	}
	d := New(hero.Viserai{}, nil, deckCards)
	rng := rand.New(rand.NewSource(42))
	d.Evaluate(20, 0, rng)

	if d.Stats.PerCard[card.SigilOfTheArknightBlue].Plays == 0 {
		t.Fatal("Sigil was never played — test fixture not provoking the code path")
	}
	// Across 20 runs * multiple turns each, the best-value turn almost certainly had a Sigil
	// queued from the prior turn. Failing here means the delayed bookkeeping never reached
	// Stats.Best.
	if len(d.Stats.Best.Summary.DelayedFromLastTurn) == 0 {
		t.Errorf("Stats.Best.Summary.DelayedFromLastTurn is empty; Best.Value=%d; Sigils played=%d",
			d.Stats.Best.Summary.Value, d.Stats.PerCard[card.SigilOfTheArknightBlue].Plays)
	}
}
