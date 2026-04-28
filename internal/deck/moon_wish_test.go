package deck

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// moonWishHero is the no-op hero used by the Moon Wish e2e tests so the assertions on Value
// / Hand / Deck / Arsenal reflect Moon Wish's plumbing alone (no Viserai-style
// OnCardPlayed runechant credit perturbing the numbers). Intel matches a typical adult hero
// hand size.
var moonWishHero = testutils.Hero{Intel: 4}

// Moon Wish's printed text shuffles the deck after the on-hit Sun Kiss tutor; assertions
// here therefore avoid pinning specific deck positions or which card lands in arsenal off a
// post-tutor DrawOne (those are random per the printed shuffle even though our model
// currently skips it). The valid checks are: Value, the surviving copy count of each card
// across turn-2 surfaces (Hand + Deck + Arsenal), and which BestLine roles got assigned.

// TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck is the canonical
// no-go-again scenario for Moon Wish run end-to-end through one turn:
//
//   - Hand: Moon Wish [Y], Weeping Battleground (Red — a red Defense Reaction).
//   - Deck: Sun Kiss [R] plus filler.
//
// Expected line: Moon Wish plays via its alt cost — Weeping Battleground is the consumed
// hand card and gets returned to the deck. Moon Wish hits, tutors Sun Kiss out of the deck,
// and (without go-again) appends Sun Kiss to s.Hand for the post-hoc Arsenal promotion.
//
// Cross-turn assertions:
//   - turn-1 Value = 4 (Moon Wish base only; Sun Kiss tutored, not played).
//   - Sun Kiss appears in NO turn-2 surface (Hand / Deck / Arsenal) directly — but is
//     promoted to Arsenal as the only Held candidate after alt cost consumed the DR.
//   - Weeping Battleground exists exactly once across turn-2 surfaces — the alt-cost
//     mechanism didn't drop the card on the floor (a pre-fix run lost it entirely).
func TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck(t *testing.T) {
	deckCards := []card.Card{
		cards.SunKissRed{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := New(moonWishHero, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if state.PrevTurnValue != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; Sun Kiss tutored, not played)",
			state.PrevTurnValue)
	}
	if state.ArsenalCard == nil || state.ArsenalCard.ID() != card.SunKissRed {
		t.Errorf("ArsenalCard = %v, want Sun Kiss [R] (post-hoc promoted from State.Hand)",
			state.ArsenalCard)
	}
	if got := countAcrossSurfaces(state, card.SunKissRed); got != 1 {
		t.Errorf("Sun Kiss [R] total across turn-2 Hand/Deck/Arsenal = %d, want 1 (in Arsenal)",
			got)
	}
	if got := countAcrossSurfaces(state, card.WeepingBattlegroundRed); got != 1 {
		t.Errorf("Weeping Battleground [R] total across turn-2 surfaces = %d, want 1 "+
			"(alt cost returned it to deck — should still exist somewhere)",
			got)
	}
}

// TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss is the negative-tutor variant:
// alt cost still returns Weeping Battleground to the deck, but with no Sun Kiss in the deck
// the tutor finds nothing and bails. No drawn card lands; arsenal stays empty (the only
// Held candidate was consumed by alt cost).
func TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss(t *testing.T) {
	deckCards := []card.Card{
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := New(moonWishHero, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if state.PrevTurnValue != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; tutor fizzles)",
			state.PrevTurnValue)
	}
	if state.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (DR was the only Held; alt cost consumed it)",
			state.ArsenalCard)
	}
	if got := countAcrossSurfaces(state, card.WeepingBattlegroundRed); got != 1 {
		t.Errorf("Weeping Battleground [R] total across turn-2 surfaces = %d, want 1 "+
			"(alt cost returned it to deck even when the tutor fizzled)",
			got)
	}
}

// TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss covers the go-again branch:
// Flying High [R] plays first in the chain and grants Moon Wish go-again. Moon Wish then
// alt-costs the Held DR, hits, tutors Sun Kiss, sees self.EffectiveGoAgain() = true, and
// plays Sun Kiss immediately. Sun Kiss heads to the graveyard.
//
// Sun Kiss sits at deck position 2 (not the top) so the buf-mutation fix is actually
// exercised: a pre-fix head advance would silently consume buf[0] instead of patching out
// the specific Sun Kiss slot, leaving Sun Kiss to surface in turn-2 Hand or Deck.
//
// Cross-turn assertions:
//   - turn-1 Value = 7 (Moon Wish 4 + Sun Kiss 3) — confirms Sun Kiss actually played.
//   - PrevTurnGraveyard contains Sun Kiss — proves it was both tutored AND played, not
//     stuck in the deck or some intermediate state.
//   - ArsenalCard is non-nil — Sun Kiss's DrawOne pulled some card off the (per printed
//     rules, shuffled) deck top; the specific identity is random so we don't pin it.
func TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss(t *testing.T) {
	// Sun Kiss at index 2 (not 0) so a pre-fix head++ wouldn't accidentally consume its
	// slot; verifying buf-removal logic actually patches the specific tutored card out.
	deckCards := []card.Card{
		testutils.RedAttack{}, testutils.RedAttack{},
		cards.SunKissRed{},
		testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{}, testutils.RedAttack{},
	}
	d := New(moonWishHero, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{
		cards.FlyingHighRed{},
		cards.MoonWishYellow{},
		cards.WeepingBattlegroundRed{},
	})

	if state.PrevTurnValue != 7 {
		t.Errorf("turn-1 Value = %d, want 7 (Moon Wish 4 + Sun Kiss 3 via Flying High go-again)",
			state.PrevTurnValue)
	}
	skInGraveyard := false
	for _, c := range state.PrevTurnGraveyard {
		if c.ID() == card.SunKissRed {
			skInGraveyard = true
			break
		}
	}
	if !skInGraveyard {
		t.Errorf("Sun Kiss [R] not in turn-1 Graveyard %v; want it there (tutored and played)",
			cardNames(state.PrevTurnGraveyard))
	}
	if got := countAcrossSurfaces(state, card.SunKissRed); got != 0 {
		t.Errorf("Sun Kiss [R] total across turn-2 surfaces = %d, want 0 (it's in the graveyard)",
			got)
	}
	if state.ArsenalCard == nil {
		t.Error("ArsenalCard = nil; want any card (Sun Kiss's DrawOne pulled one card into " +
			"State.Hand → Arsenal promotion is the only candidate)")
	}
}

// cardNames returns the names of cs in order — handy for error messages so a failure shows
// what's actually in a slice instead of opaque %v formatting.
func cardNames(cs []card.Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}

// countAcrossSurfaces totals the occurrences of id across turn-2 Hand, Deck, and Arsenal —
// the surfaces TurnStartState exposes. Used by tests that need to assert "this card exists /
// doesn't exist" without pinning a specific position (positions are randomised by Moon
// Wish's printed-shuffle, even when our current model skips it).
func countAcrossSurfaces(state TurnStartState, id card.ID) int {
	n := 0
	for _, c := range state.Hand {
		if c.ID() == id {
			n++
		}
	}
	for _, c := range state.Deck {
		if c.ID() == id {
			n++
		}
	}
	if state.ArsenalCard != nil && state.ArsenalCard.ID() == id {
		n++
	}
	return n
}
