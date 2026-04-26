package deck

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/stubs"
)

// moonWishHero is the no-op hero used by the Moon Wish e2e tests so the assertions on Value
// / Hand / Deck / ArsenalCard reflect Moon Wish's plumbing alone (no Viserai-style
// OnCardPlayed runechant credit perturbing the numbers). Intel matches a typical adult hero
// hand size.
var moonWishHero = stubs.Hero{Intel: 4}

// TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck is the canonical
// no-go-again scenario for Moon Wish run end-to-end through one turn so the buf-mutation
// fix can be observed in the next-turn state:
//
//   - Hand: Moon Wish (Yellow), Weeping Battleground (Red — a red Defense Reaction).
//   - Deck: Sun Kiss (Red) on top of five filler cards.
//
// Expected line: Moon Wish plays via its alt cost — Weeping Battleground is the consumed
// Held card and lands on top of the deck buffer. Moon Wish hits, tutors Sun Kiss out of the
// deck, and (without go-again) appends Sun Kiss to s.Drawn for the post-hoc Arsenal
// promotion.
//
// Cross-turn assertions:
//   - Turn-2 Hand[0] is Weeping Battleground — proves the alt-cost'd card is actually on
//     top of the deck at start of next turn (not just rewritten in a per-Play s.Deck slice).
//   - Turn-2 Deck holds zero Sun Kiss copies — proves the tutor patched the underlying buf.
//   - ArsenalCard is Sun Kiss (only Held candidate after alt cost consumed the DR).
func TestEvalOneTurn_MoonWishAltCostTutorsSunKissAndConsumesDeck(t *testing.T) {
	deckCards := []card.Card{
		generic.SunKissRed{},
		fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{},
		fake.RedAttack{}, fake.RedAttack{},
	}
	d := New(moonWishHero, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{
		generic.MoonWishYellow{},
		runeblade.WeepingBattlegroundRed{},
	})

	if state.PrevTurnValue != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; Sun Kiss tutored, not played)",
			state.PrevTurnValue)
	}
	if state.ArsenalCard == nil || state.ArsenalCard.ID() != card.SunKissRed {
		t.Errorf("ArsenalCard = %v, want Sun Kiss (Red)", state.ArsenalCard)
	}
	if len(state.Hand) == 0 || state.Hand[0].ID() != card.WeepingBattlegroundRed {
		t.Errorf("turn-2 Hand[0] = %v, want Weeping Battleground (Red) (alt-cost'd card on deck top)",
			state.Hand)
	}
	for i, c := range state.Deck {
		if c.ID() == card.SunKissRed {
			t.Errorf("turn-2 Deck[%d] is Sun Kiss (Red); want 0 copies after tutor patched buf", i)
		}
	}
}

// TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss is the negative-tutor variant:
// alt cost still places Weeping Battleground on top of the deck buffer, but with no Sun
// Kiss in the deck the tutor finds nothing and bails. No drawn card lands; arsenal stays
// empty (the only Held candidate was consumed by alt cost).
//
// Cross-turn assertion mirrors the prior test: turn-2 Hand[0] is the alt-cost'd card.
func TestEvalOneTurn_MoonWishAltCostTutorFizzlesWithoutSunKiss(t *testing.T) {
	deckCards := []card.Card{
		fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{},
		fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{},
	}
	d := New(moonWishHero, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{
		generic.MoonWishYellow{},
		runeblade.WeepingBattlegroundRed{},
	})

	if state.PrevTurnValue != 4 {
		t.Errorf("turn-1 Value = %d, want 4 (Moon Wish base attack; tutor fizzles)",
			state.PrevTurnValue)
	}
	if state.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (DR was the only Held; alt cost consumed it)",
			state.ArsenalCard)
	}
	if len(state.Hand) == 0 || state.Hand[0].ID() != card.WeepingBattlegroundRed {
		t.Errorf("turn-2 Hand[0] = %v, want Weeping Battleground (Red) (alt-cost'd, no tutor)",
			state.Hand)
	}
}

// TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss covers the go-again branch:
// Flying High (Red) plays first in the chain and grants Moon Wish go-again. Moon Wish then
// alt-costs the Held DR, hits, tutors Sun Kiss, sees self.EffectiveGoAgain() = true, and
// plays Sun Kiss immediately. Sun Kiss's synergy fires (Moon Wish appears in CardsPlayed
// via the transient pre-append), draws the alt-cost'd Weeping Battleground off top of deck,
// and grants (irrelevant) go-again to its synthetic CardState before heading to the
// graveyard.
//
// Sun Kiss sits at deck position 2 (not the top) so the buf-mutation fix is actually
// exercised: a pre-fix head advance would silently consume buf[0] instead of patching out
// the specific Sun Kiss slot, leaving Sun Kiss to surface in turn-2 Hand or Deck.
//
// Cross-turn assertions:
//   - turn-1 Value = 7 (Moon Wish 4 + Sun Kiss 3) — confirms Sun Kiss actually played.
//   - Sun Kiss appears in NEITHER turn-2 Hand NOR Deck NOR Arsenal — it was tutored AND
//     played, so it's in the graveyard.
//   - ArsenalCard is Weeping Battleground — pulled by Sun Kiss's DrawOne off the alt-cost'd
//     deck top, then Held in Drawn[] and promoted to Arsenal as the only candidate.
func TestEvalOneTurn_MoonWishWithFlyingHighPlaysTutoredSunKiss(t *testing.T) {
	// Sun Kiss at index 2 (not 0) so a pre-fix head++ wouldn't accidentally consume its
	// slot; verifying buf-removal logic actually patches the specific tutored card out.
	deckCards := []card.Card{
		fake.RedAttack{}, fake.RedAttack{},
		generic.SunKissRed{},
		fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{}, fake.RedAttack{},
	}
	d := New(moonWishHero, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{
		generic.FlyingHighRed{},
		generic.MoonWishYellow{},
		runeblade.WeepingBattlegroundRed{},
	})

	if state.PrevTurnValue != 7 {
		t.Errorf("turn-1 Value = %d, want 7 (Moon Wish 4 + Sun Kiss 3 via Flying High go-again)",
			state.PrevTurnValue)
	}
	if state.ArsenalCard == nil || state.ArsenalCard.ID() != card.WeepingBattlegroundRed {
		t.Errorf("ArsenalCard = %v, want Weeping Battleground (Red) "+
			"(Sun Kiss's DrawOne pulled the alt-cost'd DR off deck top)",
			state.ArsenalCard)
	}
	// Sun Kiss was tutored and played — it should be in the graveyard, NOT in any of the
	// next-turn surfaces (Hand, Deck, Arsenal). A pre-fix run would leak Sun Kiss into one
	// of these because the head++ accounting can't pinpoint the tutored slot.
	for i, c := range state.Hand {
		if c.ID() == card.SunKissRed {
			t.Errorf("turn-2 Hand[%d] is Sun Kiss (Red); should be in graveyard (tutored and played)", i)
		}
	}
	for i, c := range state.Deck {
		if c.ID() == card.SunKissRed {
			t.Errorf("turn-2 Deck[%d] is Sun Kiss (Red); should be in graveyard (tutored and played)", i)
		}
	}
	if state.ArsenalCard != nil && state.ArsenalCard.ID() == card.SunKissRed {
		t.Error("ArsenalCard is Sun Kiss; should be in graveyard (tutored and played)")
	}
	// Confirm the chain choice via the surfaced BestLine.
	mw, dr := bestLineRoles(state.PrevTurnBestLine, card.MoonWishYellow, card.WeepingBattlegroundRed)
	if mw != hand.Attack {
		t.Errorf("Moon Wish role = %s, want ATTACK", mw)
	}
	if dr != hand.Held {
		t.Errorf("Weeping Battleground role = %s, want HELD (consumed via HeldConsumed, not flipped)", dr)
	}
}

// bestLineRoles returns the BestLine roles assigned to the cards with mwID and drID. Returns
// hand.Held for either when the card is missing — a missing card surfaces as an unexpected
// role mismatch in the calling assertion rather than a panic.
func bestLineRoles(line []hand.CardAssignment, mwID, drID card.ID) (hand.Role, hand.Role) {
	var mw, dr hand.Role = hand.Held, hand.Held
	for _, a := range line {
		switch a.Card.ID() {
		case mwID:
			mw = a.Role
		case drID:
			dr = a.Role
		}
	}
	return mw, dr
}
