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
	contribs, total, revealed, _ := runDelayedPlays([]card.Card{a, b}, nil, nil)
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
	contribs, total, revealed, _ := runDelayedPlays(nil, nil, nil)
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

// TestRunDelayedPlays_RevealsAttackActionIntoHand verifies Sigil of the Arknight's reveal:
// the top card (an attack action) comes back in the revealed slice, and the contribution has
// ToHand set for the formatter.
func TestRunDelayedPlays_RevealsAttackActionIntoHand(t *testing.T) {
	sigil := runeblade.SigilOfTheArknightBlue{}
	slash := runeblade.AetherSlashRed{}
	contribs, total, revealed, _ := runDelayedPlays([]card.Card{sigil}, []card.Card{slash}, nil)
	if total != 0 {
		t.Errorf("total = %d, want 0 (reveal contributes via hand, not damage)", total)
	}
	if len(revealed) != 1 || revealed[0] != slash {
		t.Errorf("revealed = %v, want [%v]", revealed, slash)
	}
	if len(contribs) != 1 || contribs[0].ToHand != slash {
		t.Errorf("contribs[0].ToHand = %v, want %v", contribs[0].ToHand, slash)
	}
}

// TestRunDelayedPlays_CascadingReveals: two sigils in a row each reveal the current top, so the
// second one sees the NEW top after the first pops its card off the front of the deck view.
func TestRunDelayedPlays_CascadingReveals(t *testing.T) {
	sigil := runeblade.SigilOfTheArknightBlue{}
	first := runeblade.AetherSlashRed{}
	second := runeblade.ConsumingVolitionRed{}
	_, _, revealed, _ := runDelayedPlays([]card.Card{sigil, sigil}, []card.Card{first, second}, nil)
	if len(revealed) != 2 {
		t.Fatalf("len(revealed) = %d, want 2 (two cascading reveals)", len(revealed))
	}
	if revealed[0] != first || revealed[1] != second {
		t.Errorf("revealed = %v, want [%v, %v]", revealed, first, second)
	}
}

// TestRunDelayedPlays_NonAttackTopSkipsReveal: sigil peeks a non-attack top → no reveal, no
// damage. The top card stays on the deck in the real game.
func TestRunDelayedPlays_NonAttackTopSkipsReveal(t *testing.T) {
	sigil := runeblade.SigilOfTheArknightBlue{}
	// Sigil itself is an Aura (non-attack) — use it as a convenient non-attack top.
	_, total, revealed, _ := runDelayedPlays([]card.Card{sigil}, []card.Card{sigil}, nil)
	if total != 0 {
		t.Errorf("total = %d, want 0 (non-attack top, no credit)", total)
	}
	if revealed != nil {
		t.Errorf("revealed = %v, want nil (non-attack tops aren't moved)", revealed)
	}
}

// TestEvalOneTurn_SigilOfTheArknightRevealsIntoHand is the end-to-end 2-turn check: turn 1
// starts with a Sigil of the Arknight as the ONLY card in hand. The solver plays it (the
// beatsBest tiebreaker prefers playing DelayedPlay cards at equal Value over Held → arsenal
// promotion, crediting their hidden next-turn payoff). The sigil queues its PlayNextTurn
// callback; on turn 2 the callback peeks the top of the post-draw deck — an attack action —
// and moves it into the hand. The returned turn-2 hand should have 5 cards: 4 normal refills
// plus the revealed Aether Slash appended at the tail.
func TestEvalOneTurn_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := runeblade.SigilOfTheArknightBlue{}
	reveal := runeblade.AetherSlashRed{}
	// Deck layout: positions 0..3 are turn 2's normal refill (Blues), position 4 is the reveal
	// target at the post-draw top, positions 5+ are unused filler.
	deckCards := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		reveal,
		fake.BlueAttack{},
	}
	d := New(hero.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{sigil})

	// Assert sigil played: find it as Role=Attack in turn 1's BestLine.
	sigilPlayed := false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == card.SigilOfTheArknightBlue && a.Role == hand.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", state.PrevTurnBestLine)
	}

	// Turn 2: 4 normal draws + 1 revealed = 5 cards. deckCards[0..3] refill turn 2's hand;
	// deckCards[4] is the reveal target appended at the tail.
	wantHand := []card.Card{
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		fake.BlueAttack{},
		reveal,
	}
	if len(state.Hand) != len(wantHand) {
		t.Fatalf("turn 2 hand size = %d, want %d (4 normal draws + 1 revealed)", len(state.Hand), len(wantHand))
	}
	for i, want := range wantHand {
		if state.Hand[i] != want {
			t.Errorf("turn 2 hand[%d] = %v, want %v", i, state.Hand[i], want)
		}
	}
}

// TestEvaluate_DelayedFromLastTurnSurfacesInBest runs a full Evaluate with Sigil of the
// Arknight in the deck and asserts the PlayNextTurn callback lands a DelayedFromLastTurn
// entry on at least some hand's TurnSummary. Uses enough copies + runs that the shuffle
// reliably plays Sigil before the best turn is recorded.
func TestEvaluate_DelayedFromLastTurnSurfacesInBest(t *testing.T) {
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
