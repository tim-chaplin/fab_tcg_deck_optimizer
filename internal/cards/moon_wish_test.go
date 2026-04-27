package cards

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMoonWish_VariableCost: Cost reads len(s.Hand). With any hand card the alt cost fires
// and the card costs 0; without one we fall back to the printed 2. Bounds are static (Min=0,
// Max=2) for the solver pre-screens.
func TestMoonWish_VariableCost(t *testing.T) {
	cases := []card.Card{MoonWishRed{}, MoonWishYellow{}, MoonWishBlue{}}
	for _, c := range cases {
		held := card.TurnState{Hand: []card.Card{stubGenericAttack(0, 0)}}
		if got := c.Cost(&held); got != 0 {
			t.Errorf("%s: Cost(Hand) = %d, want 0", c.Name(), got)
		}
		empty := card.TurnState{}
		if got := c.Cost(&empty); got != 2 {
			t.Errorf("%s: Cost(empty) = %d, want 2", c.Name(), got)
		}
		vc, ok := c.(card.VariableCost)
		if !ok {
			t.Errorf("%s: missing card.VariableCost", c.Name())
			continue
		}
		if vc.MinCost() != 0 || vc.MaxCost() != 2 {
			t.Errorf("%s: bounds = [%d, %d], want [0, 2]", c.Name(), vc.MinCost(), vc.MaxCost())
		}
	}
}

// TestMoonWish_AltCostMovesHandCardToDeckTop: when Play fires the alt cost it pops the first
// hand card and prepends it to s.Deck. Pins both the state-mutation contract and the
// top-of-deck placement, plus the post-trigger "returned X to top of deck" log line that
// names the moved card under Moon Wish's chain entry.
func TestMoonWish_AltCostMovesHandCardToDeckTop(t *testing.T) {
	dr := stubGenericAttack(0, 0)
	dr.name = "dr"
	other := stubGenericAttack(0, 0)
	other.name = "deckTop"
	s := card.TurnState{
		Hand: []card.Card{dr},
		Deck: []card.Card{other},
	}
	self := &card.CardState{Card: MoonWishYellow{}}
	MoonWishYellow{}.Play(&s, self)
	if len(s.Hand) != 0 {
		t.Errorf("Hand = %d entries, want 0 (alt cost should pop the only hand card)", len(s.Hand))
	}
	if len(s.Deck) != 2 || s.Deck[0].Name() != "dr" || s.Deck[1].Name() != "deckTop" {
		t.Errorf("Deck = %v, want [dr, deckTop] (alt-cost'd card on top)",
			[]string{s.Deck[0].Name(), s.Deck[1].Name()})
	}
	// One of the post-trigger log entries should name the returned card.
	wantSuffix := "returned " + card.DisplayName(dr) + " to top of deck"
	found := false
	for _, e := range s.Log {
		if e.Source == "Moon Wish [Y]" && strings.HasSuffix(e.Text, wantSuffix) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a Moon Wish post-trigger log line ending in %q; log = %+v", wantSuffix, s.Log)
	}
}

// TestMoonWish_TutorPrefersRedSunKissThenYellowThenBlue: when multiple Sun Kiss variants are
// in deck the tutor picks the highest-power printing first — Red heals 3, Yellow 2, Blue 1.
func TestMoonWish_TutorPrefersRedSunKissThenYellowThenBlue(t *testing.T) {
	cases := []struct {
		name string
		deck []card.Card
		want card.ID
	}{
		{"red beats yellow and blue", []card.Card{SunKissBlue{}, SunKissYellow{}, SunKissRed{}}, card.SunKissRed},
		{"yellow beats blue", []card.Card{SunKissBlue{}, SunKissYellow{}}, card.SunKissYellow},
		{"blue alone wins", []card.Card{SunKissBlue{}}, card.SunKissBlue},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := bestSunKissInDeck(tc.deck)
			if got == nil || got.ID() != tc.want {
				t.Errorf("bestSunKissInDeck = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestMoonWish_TutorRequiresHit: Moon Wish's hit check (LikelyToHit) gates the tutor. With
// the printed 4 power Moon Wish [Y] sits in the hit window and tutors Sun Kiss into
// hand; with a -4 BonusAttack it doesn't and the deck stays intact.
func TestMoonWish_TutorRequiresHit(t *testing.T) {
	{
		s := card.TurnState{Deck: []card.Card{SunKissRed{}}}
		self := &card.CardState{Card: MoonWishYellow{}}
		MoonWishYellow{}.Play(&s, self)
		if len(s.Hand) != 1 || s.Hand[0].ID() != card.SunKissRed {
			t.Errorf("base hit: Hand = %v, want [Sun Kiss [R]]", s.Hand)
		}
		if len(s.Deck) != 0 {
			t.Errorf("base hit: Deck = %v, want [] (tutor removed Sun Kiss)", s.Deck)
		}
	}
	{
		s := card.TurnState{Deck: []card.Card{SunKissRed{}}}
		// Drive EffectiveAttack down so LikelyToHit fails (4 - 4 = 0, clamped, not in window).
		self := &card.CardState{Card: MoonWishYellow{}, BonusAttack: -4}
		MoonWishYellow{}.Play(&s, self)
		if len(s.Hand) != 0 {
			t.Errorf("dampened: Hand = %v, want [] (no hit, no tutor)", s.Hand)
		}
		if len(s.Deck) != 1 || s.Deck[0].ID() != card.SunKissRed {
			t.Errorf("dampened: Deck = %v, want [Sun Kiss [R]] (untouched)", s.Deck)
		}
	}
}

// TestMoonWish_GoAgainPlaysSunKissImmediately: with self.GrantedGoAgain set, the tutored
// Sun Kiss plays immediately (added to graveyard, damage folded into the return) and does
// NOT land in s.Hand. Without go-again it stays in s.Hand for the next turn.
func TestMoonWish_GoAgainPlaysSunKissImmediately(t *testing.T) {
	{
		s := card.TurnState{Deck: []card.Card{SunKissRed{}}}
		self := &card.CardState{Card: MoonWishYellow{}, GrantedGoAgain: true}
		MoonWishYellow{}.Play(&s, self)
		dmg := s.Value
		if dmg != 4+3 {
			t.Errorf("with go-again: damage = %d, want 7 (Moon Wish 4 + Sun Kiss 3)", dmg)
		}
		if len(s.Hand) != 0 {
			t.Errorf("with go-again: Hand = %v, want [] (Sun Kiss played, not tutored to hand)", s.Hand)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != card.SunKissRed {
			t.Errorf("with go-again: Graveyard = %v, want [Sun Kiss [R]]", s.Graveyard)
		}
	}
	{
		s := card.TurnState{Deck: []card.Card{SunKissRed{}}}
		self := &card.CardState{Card: MoonWishYellow{}}
		MoonWishYellow{}.Play(&s, self)
		dmg := s.Value
		if dmg != 4 {
			t.Errorf("no go-again: damage = %d, want 4 (Sun Kiss not played)", dmg)
		}
		if len(s.Hand) != 1 || s.Hand[0].ID() != card.SunKissRed {
			t.Errorf("no go-again: Hand = %v, want [Sun Kiss [R]]", s.Hand)
		}
	}
}
