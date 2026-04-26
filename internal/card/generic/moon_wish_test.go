package generic

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestMoonWish_VariableCost: Cost reads len(s.Held). With any Held card the alt cost fires
// and the card costs 0; without one we fall back to the printed 2. Bounds are static (Min=0,
// Max=2) for the solver pre-screens.
func TestMoonWish_VariableCost(t *testing.T) {
	cases := []card.Card{MoonWishRed{}, MoonWishYellow{}, MoonWishBlue{}}
	for _, c := range cases {
		held := card.TurnState{Held: []card.Card{stubGenericAttack(0, 0)}}
		if got := c.Cost(&held); got != 0 {
			t.Errorf("%s: Cost(Held) = %d, want 0", c.Name(), got)
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

// TestMoonWish_AltCostMovesHeldToDeckTop: when Play fires the alt cost it pops the first
// Held card, prepends it to s.Deck, and records it on s.HeldConsumed. Pins both the
// state-mutation contract and the top-of-deck placement.
func TestMoonWish_AltCostMovesHeldToDeckTop(t *testing.T) {
	dr := stubGenericAttack(0, 0)
	dr.name = "dr"
	other := stubGenericAttack(0, 0)
	other.name = "deckTop"
	s := card.TurnState{
		Held: []card.Card{dr},
		Deck: []card.Card{other},
	}
	self := &card.CardState{Card: MoonWishYellow{}}
	_ = MoonWishYellow{}.Play(&s, self)
	if len(s.Held) != 0 {
		t.Errorf("Held = %d entries, want 0 (alt cost should pop the only Held card)", len(s.Held))
	}
	if len(s.HeldConsumed) != 1 || s.HeldConsumed[0].Name() != "dr" {
		t.Errorf("HeldConsumed = %v, want [dr]", s.HeldConsumed)
	}
	if len(s.Deck) != 2 || s.Deck[0].Name() != "dr" || s.Deck[1].Name() != "deckTop" {
		t.Errorf("Deck = %v, want [dr, deckTop] (alt-cost'd card on top)",
			[]string{s.Deck[0].Name(), s.Deck[1].Name()})
	}
}

// TestMoonWish_TutorPrefersRedSunKissThenYellowThenBlue: when multiple Sun Kiss variants are
// in deck the tutor picks the lowest-pitch (most flexible) printing first.
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
// the printed 4 power Moon Wish (Yellow) sits in the hit window and tutors; with a -4
// BonusAttack it doesn't.
func TestMoonWish_TutorRequiresHit(t *testing.T) {
	deck := []card.Card{SunKissRed{}}
	{
		s := card.TurnState{Deck: deck}
		self := &card.CardState{Card: MoonWishYellow{}}
		_ = MoonWishYellow{}.Play(&s, self)
		if len(s.Drawn) != 1 || s.Drawn[0].ID() != card.SunKissRed {
			t.Errorf("base hit: Drawn = %v, want [Sun Kiss (Red)]", s.Drawn)
		}
	}
	{
		s := card.TurnState{Deck: deck}
		// Drive EffectiveAttack down so LikelyToHit fails (4 - 4 = 0, clamped, not in window).
		self := &card.CardState{Card: MoonWishYellow{}, BonusAttack: -4}
		_ = MoonWishYellow{}.Play(&s, self)
		if len(s.Drawn) != 0 {
			t.Errorf("dampened: Drawn = %v, want [] (no hit, no tutor)", s.Drawn)
		}
	}
}

// TestMoonWish_GoAgainPlaysSunKissImmediately: with self.GrantedGoAgain set, the tutored
// Sun Kiss plays immediately (added to graveyard, damage folded into the return) and does
// NOT land in s.Drawn. Without go-again it stays in s.Drawn for the next hand.
func TestMoonWish_GoAgainPlaysSunKissImmediately(t *testing.T) {
	deck := []card.Card{SunKissRed{}}
	{
		s := card.TurnState{Deck: deck}
		self := &card.CardState{Card: MoonWishYellow{}, GrantedGoAgain: true}
		dmg := MoonWishYellow{}.Play(&s, self)
		if dmg != 4+3 {
			t.Errorf("with go-again: damage = %d, want 7 (Moon Wish 4 + Sun Kiss 3)", dmg)
		}
		if len(s.Drawn) != 0 {
			t.Errorf("with go-again: Drawn = %v, want [] (Sun Kiss played, not tutored to hand)", s.Drawn)
		}
		if len(s.Graveyard) != 1 || s.Graveyard[0].ID() != card.SunKissRed {
			t.Errorf("with go-again: Graveyard = %v, want [Sun Kiss (Red)]", s.Graveyard)
		}
	}
	{
		s := card.TurnState{Deck: deck}
		self := &card.CardState{Card: MoonWishYellow{}}
		dmg := MoonWishYellow{}.Play(&s, self)
		if dmg != 4 {
			t.Errorf("no go-again: damage = %d, want 4 (Sun Kiss not played)", dmg)
		}
		if len(s.Drawn) != 1 || s.Drawn[0].ID() != card.SunKissRed {
			t.Errorf("no go-again: Drawn = %v, want [Sun Kiss (Red)]", s.Drawn)
		}
	}
}
