package cards

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestMoonWish_VariableCost: Cost reads len(s.Hand). With any hand card the alt cost fires
// and the card costs 0; without one we fall back to the printed 2. Static bounds: Min=0, Max=2.
func TestMoonWish_VariableCost(t *testing.T) {
	cases := []sim.Card{MoonWishRed{}, MoonWishYellow{}, MoonWishBlue{}}
	for _, c := range cases {
		held := sim.TurnState{Hand: []sim.Card{testutils.GenericAttack(0, 0)}}
		if got := c.Cost(&held); got != 0 {
			t.Errorf("%s: Cost(Hand) = %d, want 0", c.Name(), got)
		}
		empty := sim.TurnState{}
		if got := c.Cost(&empty); got != 2 {
			t.Errorf("%s: Cost(empty) = %d, want 2", c.Name(), got)
		}
		vc, ok := c.(sim.VariableCost)
		if !ok {
			t.Errorf("%s: missing sim.VariableCost", c.Name())
			continue
		}
		if vc.MinCost() != 0 || vc.MaxCost() != 2 {
			t.Errorf("%s: bounds = [%d, %d], want [0, 2]", c.Name(), vc.MinCost(), vc.MaxCost())
		}
	}
}

// TestMoonWish_AltCostMovesHandCardToDeckTop: when Play fires the alt cost it pops the first
// hand card and prepends it to the deck. Pins both the state-mutation contract and the
// top-of-deck placement, plus the post-trigger "returned X to top of deck" log line that
// names the moved card under Moon Wish's chain entry.
func TestMoonWish_AltCostMovesHandCardToDeckTop(t *testing.T) {
	dr := testutils.GenericAttack(0, 0).WithName("dr")
	other := testutils.GenericAttack(0, 0).WithName("deckTop")
	s := sim.NewTurnState([]sim.Card{other}, nil)
	s.Hand = []sim.Card{dr}
	self := &sim.CardState{Card: MoonWishYellow{}}
	MoonWishYellow{}.Play(s, self)
	if len(s.Hand) != 0 {
		t.Errorf("Hand = %d entries, want 0 (alt cost should pop the only hand card)", len(s.Hand))
	}
	d := s.Deck()
	if len(d) != 2 || d[0].Name() != "dr" || d[1].Name() != "deckTop" {
		t.Errorf("Deck = %v, want [dr, deckTop] (alt-cost'd card on top)",
			[]string{d[0].Name(), d[1].Name()})
	}
	// One of the post-trigger log entries should name the returned card.
	wantSuffix := "returned " + sim.DisplayName(dr) + " to top of deck"
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
// Drives the priority through the live Play path so we exercise the same TutorFromDeck call
// production uses.
func TestMoonWish_TutorPrefersRedSunKissThenYellowThenBlue(t *testing.T) {
	cases := []struct {
		name string
		deck []sim.Card
		want ids.CardID
	}{
		{"red beats yellow and blue", []sim.Card{SunKissBlue{}, SunKissYellow{}, SunKissRed{}}, ids.SunKissRed},
		{"yellow beats blue", []sim.Card{SunKissBlue{}, SunKissYellow{}}, ids.SunKissYellow},
		{"blue alone wins", []sim.Card{SunKissBlue{}}, ids.SunKissBlue},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := sim.NewTurnState(append([]sim.Card(nil), tc.deck...), nil)
			self := &sim.CardState{Card: MoonWishYellow{}}
			MoonWishYellow{}.Play(s, self)
			if len(s.Hand) != 1 || s.Hand[0].ID() != tc.want {
				t.Errorf("Hand = %v, want first entry to be %v", s.Hand, tc.want)
			}
		})
	}
}

// TestMoonWish_TutorRequiresHit: Moon Wish's hit check (LikelyToHit) gates the tutor. With
// the printed 4 power Moon Wish [Y] sits in the hit window and tutors Sun Kiss into
// hand; with a -4 BonusAttack it doesn't and the deck stays intact.
func TestMoonWish_TutorRequiresHit(t *testing.T) {
	{
		s := sim.NewTurnState([]sim.Card{SunKissRed{}}, nil)
		self := &sim.CardState{Card: MoonWishYellow{}}
		MoonWishYellow{}.Play(s, self)
		if len(s.Hand) != 1 || s.Hand[0].ID() != ids.SunKissRed {
			t.Errorf("base hit: Hand = %v, want [Sun Kiss [R]]", s.Hand)
		}
		if d := s.Deck(); len(d) != 0 {
			t.Errorf("base hit: Deck = %v, want [] (tutor removed Sun Kiss)", d)
		}
	}
	{
		s := sim.NewTurnState([]sim.Card{SunKissRed{}}, nil)
		// Drive EffectiveAttack down so LikelyToHit fails (4 - 4 = 0, clamped, not in window).
		self := &sim.CardState{Card: MoonWishYellow{}, BonusAttack: -4}
		MoonWishYellow{}.Play(s, self)
		if len(s.Hand) != 0 {
			t.Errorf("dampened: Hand = %v, want [] (no hit, no tutor)", s.Hand)
		}
		if d := s.Deck(); len(d) != 1 || d[0].ID() != ids.SunKissRed {
			t.Errorf("dampened: Deck = %v, want [Sun Kiss [R]] (untouched)", d)
		}
	}
}

// TestMoonWish_GoAgainPlaysSunKissImmediately: with self.GrantedGoAgain set, the tutored
// Sun Kiss plays immediately (added to graveyard, damage folded into the return) and does
// NOT land in s.Hand. Without go-again it stays in s.Hand for the next turn.
func TestMoonWish_GoAgainPlaysSunKissImmediately(t *testing.T) {
	{
		s := sim.NewTurnState([]sim.Card{SunKissRed{}}, nil)
		self := &sim.CardState{Card: MoonWishYellow{}, GrantedGoAgain: true}
		MoonWishYellow{}.Play(s, self)
		dmg := s.Value
		if dmg != 4+3 {
			t.Errorf("with go-again: damage = %d, want 7 (Moon Wish 4 + Sun Kiss 3)", dmg)
		}
		if len(s.Hand) != 0 {
			t.Errorf("with go-again: Hand = %v, want [] (Sun Kiss played, not tutored to hand)", s.Hand)
		}
		g := s.Graveyard()
		if len(g) != 1 || g[0].ID() != ids.SunKissRed {
			t.Errorf("with go-again: Graveyard = %v, want [Sun Kiss [R]]", g)
		}
	}
	{
		s := sim.NewTurnState([]sim.Card{SunKissRed{}}, nil)
		self := &sim.CardState{Card: MoonWishYellow{}}
		MoonWishYellow{}.Play(s, self)
		dmg := s.Value
		if dmg != 4 {
			t.Errorf("no go-again: damage = %d, want 4 (Sun Kiss not played)", dmg)
		}
		if len(s.Hand) != 1 || s.Hand[0].ID() != ids.SunKissRed {
			t.Errorf("no go-again: Hand = %v, want [Sun Kiss [R]]", s.Hand)
		}
	}
}
