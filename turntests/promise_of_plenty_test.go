package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// End-to-end: Promise of Plenty (Red, printed power 3) plays with an empty arsenal and a
// non-empty deck — the on-hit rider fires through finalizeActiveAttack's hit gate, and the
// post-attack-turn arsenal is occupied with the popped deck-top card. The fill doesn't flip
// IsCacheable because the engine accessor never observes the card's identity.
func TestPromiseOfPlenty_OnHitFillsEmptyArsenal(t *testing.T) {
	deckTop := testutils.FakeRedAttack().WithName("DeckTop")
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{deckTop})
	// Promise Blue has printed power 1 — squarely in the LikelyDamageDealt hit window
	// (1/4/7 leak past the heuristic blocker), so finalizeActiveAttack fires the OnHit.
	hand := []card.Card{cards.PromiseOfPlentyBlue{}}
	summary := sim.EvalOneTurnForTesting(d, nil, hand)

	if summary.State.Arsenal() == nil {
		t.Errorf("Arsenal still empty after Promise hits; fill rider didn't fire\nBestLine: %s",
			sim.FormatBestLine(summary.BestLine))
	}
}

// End-to-end: Promise of Plenty leaves an already-occupied arsenal alone. Uses the Blue
// variant (power 1, in the LikelyDamageDealt hit window) so the OnHit rider actually
// fires; otherwise the test would pass for the wrong reason (rider never invoked).
func TestPromiseOfPlenty_OccupiedArsenalUntouched(t *testing.T) {
	existing := testutils.FakeRedAttack().WithName("Existing")
	deckTop := testutils.FakeRedAttack().WithName("DeckTop")
	d := deck.New(testutils.Hero{Intel: 4}, nil, []deck.Card{deckTop})
	state := gameengine.GameStateBuilder().
		SetHero(testutils.Hero{Intel: 4}).
		SetArsenal(existing).
		Build()
	hand := []card.Card{cards.PromiseOfPlentyBlue{}}
	summary := sim.EvalOneTurnForTesting(d, state, hand)

	if got := summary.State.Arsenal(); got == nil || got.Name() != "Existing" {
		t.Errorf("Arsenal = %v, want Existing (occupied slot should be untouched)", got)
	}
}

// Tests the played-from-arsenal go-again rider: Promise plays from arsenal → flips
// GrantedGoAgain.
func TestPromiseOfPlenty_FromArsenalGoAgain(t *testing.T) {
	ge := gameengine.New()
	pc := &card.CardState{Card: cards.PromiseOfPlentyRed{}, FromArsenal: true}
	(cards.PromiseOfPlentyRed{}).Play(ge, ge.Logger(), pc)
	if !pc.GrantedGoAgain {
		t.Errorf("GrantedGoAgain = false after Play from arsenal, want true")
	}
}
