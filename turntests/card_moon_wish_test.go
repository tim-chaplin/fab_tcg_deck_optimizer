package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests printed Cost() == 2 always (irrespective of hand state) and AlternativeCost
// reports (0, true) when a Held card is available, (0, false) when not.
func TestMoonWish_PrintedAndAlternativeCost(t *testing.T) {
	cases := []card.Card{cards.MoonWishRed{}, cards.MoonWishYellow{}, cards.MoonWishBlue{}}
	for _, c := range cases {
		if got := c.Cost(); got != 2 {
			t.Errorf("%s: Cost() = %d, want 2 (printed)", c.Name(), got)
		}
		ac, ok := c.(card.AlternativeCost)
		if !ok {
			t.Fatalf("%s: missing card.AlternativeCost", c.Name())
		}
		held := gameengine.New()
		held.SetHand([]card.Card{testutils.FakeRedAttack()})
		if cost, available := ac.AlternativeCost(held); cost != 0 || !available {
			t.Errorf("%s: AlternativeCost(hand=1) = (%d, %v), want (0, true)", c.Name(), cost, available)
		}
		empty := gameengine.New()
		if _, available := ac.AlternativeCost(empty); available {
			t.Errorf("%s: AlternativeCost(empty hand) reports available, want unavailable", c.Name())
		}
	}
}

// Tests that Play pops a hand card onto the deck top only when PaidAlternativeCost is set.
func TestMoonWish_AltCostMovesHandCardToDeckTop(t *testing.T) {
	dr := testutils.FakeRedAttack().WithName("dr")
	other := testutils.FakeRedAttack().WithName("deckTop")
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{other}).Build()}
	ge.SetHand([]card.Card{dr})
	pc := &card.CardState{
		Card:      cards.MoonWishYellow{},
		Ephemeral: card.Ephemeral{PaidAlternativeCost: true},
	}
	ge.ResolveAttackStep(ge.Logger(), pc)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
	if h := ge.Hand(); len(h) != 0 {
		t.Errorf("Hand = %d entries, want 0 (alt cost should pop the only hand card)", len(h))
	}
	if got := ge.Deck().Size(); got != 2 {
		t.Errorf("Deck size = %d, want 2 (alt-cost'd card prepended onto existing top)", got)
	}
	if top := ge.Deck().PeekTop(); top == nil || top.(card.Card).Name() != "dr" {
		t.Errorf("Deck top = %v, want %q (alt-cost'd card moved to top)", top, "dr")
	}
}

// Tests that Play leaves the hand and deck alone when PaidAlternativeCost is NOT set
// (the printed-cost branch — no alt-cost side effect runs).
func TestMoonWish_PrintedCostBranchSkipsAltCostSideEffect(t *testing.T) {
	dr := testutils.FakeRedAttack().WithName("dr")
	other := testutils.FakeRedAttack().WithName("deckTop")
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCards([]card.Card{other}).Build()}
	ge.SetHand([]card.Card{dr})
	pc := &card.CardState{Card: cards.MoonWishYellow{}}
	ge.ResolveAttackStep(ge.Logger(), pc)
	if h := ge.Hand(); len(h) != 1 {
		t.Errorf("Hand size = %d, want 1 (printed-cost branch should not pop the hand card)", len(h))
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
			pc := &card.CardState{Card: cards.MoonWishYellow{}}
			ge.ResolveAttackStep(ge.Logger(), pc)
			testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
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
		pc := &card.CardState{Card: cards.MoonWishYellow{}}
		ge.ResolveAttackStep(ge.Logger(), pc)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
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
		pc := &card.CardState{Card: cards.MoonWishYellow{}, Ephemeral: card.Ephemeral{BonusAttack: -4}}
		ge.ResolveAttackStep(ge.Logger(), pc)
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
		pc := &card.CardState{Card: cards.MoonWishYellow{}, Ephemeral: card.Ephemeral{GrantedGoAgain: true}}
		ge.ResolveAttackStep(ge.Logger(), pc)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
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
		pc := &card.CardState{Card: cards.MoonWishYellow{}}
		ge.ResolveAttackStep(ge.Logger(), pc)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), pc)
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
