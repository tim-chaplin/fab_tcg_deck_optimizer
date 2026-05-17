package turntests

import (
	"strings"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
)

// Tests that Cost is 0 when a hand card can pay the alt cost, else printed 2; static
// bounds are [0, 2].
func TestMoonWish_VariableCost(t *testing.T) {
	cases := []card.Card{cards.MoonWishRed{}, cards.MoonWishYellow{}, cards.MoonWishBlue{}}
	for _, c := range cases {
		held := gameengine.New()
		held.SetHand([]card.Card{testutils.GenericAttack(0, 0)})
		if got := c.Cost(held); got != 0 {
			t.Errorf("%s: Cost(Hand) = %d, want 0", c.Name(), got)
		}
		empty := gameengine.New()
		if got := c.Cost(empty); got != 2 {
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

// Tests that the alt cost pops a hand card, prepends it to the deck, and logs the
// "returned X to top of deck" rider.
func TestMoonWish_AltCostMovesHandCardToDeckTop(t *testing.T) {
	dr := testutils.GenericAttack(0, 0).WithName("dr")
	other := testutils.GenericAttack(0, 0).WithName("deckTop")
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{other}).Build()}
	ge.SetHand([]card.Card{dr})
	self := &card.CardState{Card: cards.MoonWishYellow{}}
	ge.ResolveChainStep(ge.Logger(), self)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), self)
	if h := ge.Hand(); len(h) != 0 {
		t.Errorf("Hand = %d entries, want 0 (alt cost should pop the only hand card)", len(h))
	}
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("Deck size = %d, want 2 (alt-cost'd card prepended onto existing top)", got)
	}
	if top := ge.Deck().PeekTop(); top == nil || top.(card.Card).Name() != "dr" {
		t.Errorf("Deck top = %v, want %q (alt-cost'd card moved to top)", top, "dr")
	}
	// One of the post-trigger log entries should name the returned card.
	wantSuffix := "returned " + dr.DisplayName() + " to top of deck"
	found := false
	for _, e := range ge.LogEntries() {
		if e.Source == "Moon Wish [Y]" && strings.HasSuffix(e.Text, wantSuffix) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected a Moon Wish post-trigger log line ending in %q; log = %+v", wantSuffix, ge.LogEntries())
	}
}

// Tests that the Sun Kiss tutor picks the highest-power printing (Red > Yellow > Blue).
func TestMoonWish_TutorPrefersRedSunKissThenYellowThenBlue(t *testing.T) {
	cases := []struct {
		name string
		deck []card.Card
		want ids.CardID
	}{
		{"red beats yellow and blue", []card.Card{cards.SunKissBlue{}, cards.SunKissYellow{}, cards.SunKissRed{}}, ids.SunKissRed},
		{"yellow beats blue", []card.Card{cards.SunKissBlue{}, cards.SunKissYellow{}}, ids.SunKissYellow},
		{"blue alone wins", []card.Card{cards.SunKissBlue{}}, ids.SunKissBlue},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards(append([]card.Card(nil), tc.deck...)).Build()}
			self := &card.CardState{Card: cards.MoonWishYellow{}}
			ge.ResolveChainStep(ge.Logger(), self)
			testutils.FireOnHitIfLikely(ge, ge.Logger(), self)
			h := ge.Hand()
			if len(h) != 1 || h[0].ID() != tc.want {
				t.Errorf("Hand = %v, want first entry to be %v", h, tc.want)
			}
		})
	}
}

// Tests that LikelyToHit gates the Sun Kiss tutor: a -4 BonusAttack drops the hit and
// leaves the deck intact.
func TestMoonWish_TutorRequiresHit(t *testing.T) {
	{
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{cards.SunKissRed{}}).Build()}
		self := &card.CardState{Card: cards.MoonWishYellow{}}
		ge.ResolveChainStep(ge.Logger(), self)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), self)
		h := ge.Hand()
		if len(h) != 1 || h[0].ID() != ids.SunKissRed {
			t.Errorf("base hit: Hand = %v, want [Sun Kiss [R]]", h)
		}
		if d := ge.Deck(); d.Size() != 0 {
			t.Errorf("base hit: Deck = %v, want [] (tutor removed Sun Kiss)", d)
		}
	}
	{
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{cards.SunKissRed{}}).Build()}
		// Drive EffectiveAttack down so LikelyToHit fails (4 - 4 = 0, clamped, not in window).
		self := &card.CardState{Card: cards.MoonWishYellow{}, BonusAttack: -4}
		ge.ResolveChainStep(ge.Logger(), self)
		if h := ge.Hand(); len(h) != 0 {
			t.Errorf("dampened: Hand = %v, want [] (no hit, no tutor)", h)
		}
		if d := ge.Deck(); d.Size() != 1 || d.PeekTop().(card.Card).ID() != ids.SunKissRed {
			t.Errorf("dampened: Deck = %v, want [Sun Kiss [R]] (untouched)", d)
		}
	}
}

// Tests that Sun Kiss plays immediately when self has go-again, otherwise lands in hand.
func TestMoonWish_GoAgainPlaysSunKissImmediately(t *testing.T) {
	{
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{cards.SunKissRed{}}).Build()}
		self := &card.CardState{Card: cards.MoonWishYellow{}, GrantedGoAgain: true}
		ge.ResolveChainStep(ge.Logger(), self)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), self)
		dmg := ge.Value()
		if dmg != 4+3 {
			t.Errorf("with go-again: damage = %d, want 7 (Moon Wish 4 + Sun Kiss 3)", dmg)
		}
		if h := ge.Hand(); len(h) != 0 {
			t.Errorf("with go-again: Hand = %v, want [] (Sun Kiss played, not tutored to hand)", h)
		}
		g := ge.Graveyard()
		if len(g) != 1 || g[0].ID() != ids.SunKissRed {
			t.Errorf("with go-again: Graveyard = %v, want [Sun Kiss [R]]", g)
		}
	}
	{
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{cards.SunKissRed{}}).Build()}
		self := &card.CardState{Card: cards.MoonWishYellow{}}
		ge.ResolveChainStep(ge.Logger(), self)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), self)
		dmg := ge.Value()
		if dmg != 4 {
			t.Errorf("no go-again: damage = %d, want 4 (Sun Kiss not played)", dmg)
		}
		h := ge.Hand()
		if len(h) != 1 || h[0].ID() != ids.SunKissRed {
			t.Errorf("no go-again: Hand = %v, want [Sun Kiss [R]]", h)
		}
	}
}
