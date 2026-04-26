package hand

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
)

// TestBest_MoonWishAltCostTutorsSunKissToHand drives Moon Wish's full implementation through
// hand.Best on the canonical "no go-again" scenario:
//
//   - Hand: Moon Wish (Yellow), Weeping Battleground (Red — a red Defense Reaction).
//   - Deck: a single Sun Kiss (Red).
//
// Expected line: Moon Wish plays via its alt cost (Weeping Battleground is the consumed
// Held card → top of deck). Moon Wish hits (printed 4 power lands in the 1/4/7 window),
// tutors Sun Kiss out of the deck, and — without go-again — appends Sun Kiss to the drawn
// pile. Value = 4 (Moon Wish's printed attack only; Sun Kiss isn't played this turn).
//
// Post-hoc arsenal: Weeping Battleground was Held in BestLine but consumed by alt cost
// (counted out of the BestLine candidate pool via HeldConsumed); Sun Kiss is the only
// remaining Held candidate, so it gets promoted to Arsenal.
func TestBest_MoonWishAltCostTutorsSunKissToHand(t *testing.T) {
	hand := []card.Card{generic.MoonWishYellow{}, runeblade.WeepingBattlegroundRed{}}
	deck := []card.Card{generic.SunKissRed{}}

	got := Best(stubHero, nil, hand, 0, deck, 0, nil)

	if got.Value != 4 {
		t.Errorf("Value = %d, want 4 (Moon Wish printed attack; Sun Kiss not played this turn)",
			got.Value)
	}
	if len(got.HeldConsumed) != 1 || got.HeldConsumed[0].ID() != card.WeepingBattlegroundRed {
		t.Errorf("HeldConsumed = %v, want [Weeping Battleground (Red)]", cardNames(got.HeldConsumed))
	}
	if len(got.Drawn) != 1 || got.Drawn[0].Card.ID() != card.SunKissRed {
		t.Errorf("Drawn = %v, want [Sun Kiss (Red)]", drawnNames(got.Drawn))
	}
	if got.Drawn[0].Role != Arsenal {
		t.Errorf("Drawn[0].Role = %s, want ARSENAL (only Held candidate after alt cost consumed the DR)",
			got.Drawn[0].Role)
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.SunKissRed {
		t.Errorf("ArsenalCard = %v, want Sun Kiss (Red)", got.ArsenalCard)
	}
	if len(got.AttackChain) != 1 || got.AttackChain[0].Card.ID() != card.MoonWishYellow {
		t.Errorf("AttackChain = %v, want [Moon Wish (Yellow)]", chainNames(got.AttackChain))
	}
	mwRole, drRole := bestLineRoles(got.BestLine, card.MoonWishYellow, card.WeepingBattlegroundRed)
	if mwRole != Attack {
		t.Errorf("Moon Wish role = %s, want ATTACK", mwRole)
	}
	if drRole != Held {
		t.Errorf("Weeping Battleground role = %s, want HELD (consumed via HeldConsumed, not flipped)", drRole)
	}
}

// TestBest_MoonWishAltCostTutorFizzlesWithoutSunKiss is the negative-tutor counterpart: same
// hand and budget as the prior test but the deck holds no Sun Kiss. Alt cost still fires
// (Weeping Battleground is consumed) and Moon Wish still hits, but the tutor finds no Sun
// Kiss and bails. No drawn card lands; arsenal stays empty (the only Held candidate was
// consumed by alt cost).
func TestBest_MoonWishAltCostTutorFizzlesWithoutSunKiss(t *testing.T) {
	hand := []card.Card{generic.MoonWishYellow{}, runeblade.WeepingBattlegroundRed{}}
	// Deck deliberately holds no Sun Kiss; one filler card so DrawOne (if it ever fired) has
	// something to pull, surfacing accidental draw-rider regressions.
	deck := []card.Card{runeblade.ArcanicSpikeRed{}}

	got := Best(stubHero, nil, hand, 0, deck, 0, nil)

	if got.Value != 4 {
		t.Errorf("Value = %d, want 4 (Moon Wish printed attack; tutor fizzles)", got.Value)
	}
	if len(got.HeldConsumed) != 1 || got.HeldConsumed[0].ID() != card.WeepingBattlegroundRed {
		t.Errorf("HeldConsumed = %v, want [Weeping Battleground (Red)]", cardNames(got.HeldConsumed))
	}
	if len(got.Drawn) != 0 {
		t.Errorf("Drawn = %v, want [] (tutor found no Sun Kiss; no DrawOne fired)", drawnNames(got.Drawn))
	}
	if got.ArsenalCard != nil {
		t.Errorf("ArsenalCard = %v, want nil (DR was the only hand-Held and got consumed)", got.ArsenalCard)
	}
}

// TestBest_MoonWishWithFlyingHighPlaysTutoredSunKiss covers the go-again branch: Flying High
// (Red) plays first in the chain and grants Moon Wish go-again. Moon Wish then alt-costs the
// Held DR, hits, tutors Sun Kiss, sees self.EffectiveGoAgain() = true, and plays Sun Kiss
// immediately. Sun Kiss's synergy fires (Moon Wish appears in CardsPlayed via the transient
// pre-append), draws the alt-cost'd Weeping Battleground off top of deck, and grants
// (irrelevant) go-again to its synthetic CardState before heading to the graveyard.
//
// Value = 7 (Moon Wish 4 + Sun Kiss 3). The drawn card is Weeping Battleground — the same DR
// the alt cost moved to deck top — and it post-hoc promotes to Arsenal as the only Held
// candidate left after consumption.
func TestBest_MoonWishWithFlyingHighPlaysTutoredSunKiss(t *testing.T) {
	hand := []card.Card{
		generic.FlyingHighRed{},
		generic.MoonWishYellow{},
		runeblade.WeepingBattlegroundRed{},
	}
	deck := []card.Card{generic.SunKissRed{}, runeblade.ArcanicSpikeRed{}}

	got := Best(stubHero, nil, hand, 0, deck, 0, nil)

	if got.Value != 7 {
		t.Errorf("Value = %d, want 7 (Moon Wish 4 + Sun Kiss 3 via Flying High go-again)",
			got.Value)
	}
	if len(got.HeldConsumed) != 1 || got.HeldConsumed[0].ID() != card.WeepingBattlegroundRed {
		t.Errorf("HeldConsumed = %v, want [Weeping Battleground (Red)]", cardNames(got.HeldConsumed))
	}
	if len(got.Drawn) != 1 || got.Drawn[0].Card.ID() != card.WeepingBattlegroundRed {
		t.Errorf("Drawn = %v, want [Weeping Battleground (Red)] (Sun Kiss's DrawOne pulled the alt-cost'd DR off deck top)",
			drawnNames(got.Drawn))
	}
	if got.ArsenalCard == nil || got.ArsenalCard.ID() != card.WeepingBattlegroundRed {
		t.Errorf("ArsenalCard = %v, want Weeping Battleground (Red)", got.ArsenalCard)
	}
	// AttackChain is the chain proper — Sun Kiss was tutored mid-Play and isn't in the chain.
	if len(got.AttackChain) != 2 {
		t.Fatalf("AttackChain length = %d, want 2 (Flying High + Moon Wish; Sun Kiss tutored is mid-Play)",
			len(got.AttackChain))
	}
	if got.AttackChain[0].Card.ID() != card.FlyingHighRed {
		t.Errorf("AttackChain[0] = %s, want Flying High (Red) (must precede Moon Wish to grant go-again)",
			got.AttackChain[0].Card.Name())
	}
	if got.AttackChain[1].Card.ID() != card.MoonWishYellow {
		t.Errorf("AttackChain[1] = %s, want Moon Wish (Yellow)", got.AttackChain[1].Card.Name())
	}
}

// drawnNames renders Drawn entries for failure-message rendering.
func drawnNames(d []CardAssignment) []string {
	out := make([]string, len(d))
	for i, a := range d {
		out[i] = a.Card.Name()
	}
	return out
}

// chainNames renders AttackChain entries for failure-message rendering.
func chainNames(ch []AttackChainEntry) []string {
	out := make([]string, len(ch))
	for i, e := range ch {
		out[i] = e.Card.Name()
	}
	return out
}

// bestLineRoles returns the BestLine roles assigned to the cards with mwID and drID. Defaults
// to Held on either when the card isn't in line — a missing card surfaces as an unexpected
// role mismatch in the calling assertion rather than a panic.
func bestLineRoles(line []CardAssignment, mwID, drID card.ID) (Role, Role) {
	var mw, dr Role = Held, Held
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
