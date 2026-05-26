package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// Tests that the discard rider stays dormant when ArcaneDamageDealt is false.
func TestConsumingVolition_ArcaneDamageNotDealtReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ConsumingVolitionRed{}, 4},
		{cards.ConsumingVolitionYellow{}, 3},
		{cards.ConsumingVolitionBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		cs := &card.CardState{Card: tc.c}
		ge.ResolveAttackStep(ge.Logger(), cs)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, ArcaneDamageDealt=false)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that the discard rider fires when ArcaneDamageDealt is set and the attack is likely
// to hit.
func TestConsumingVolition_LikelyToHitAndArcaneTriggersDiscard(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetArcaneDamageDealt(true).Build()}
	c := cards.ConsumingVolitionRed{}
	cs := &card.CardState{Card: c}
	ge.ResolveAttackStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	if got := ge.Value(); got != 4+3 {
		t.Errorf("Red with ArcaneDamageDealt: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// Tests that the discard rider doesn't fire on blockable variants even with
// ArcaneDamageDealt set.
func TestConsumingVolition_BlockableBaseSuppressesDiscard(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ConsumingVolitionYellow{}, 3},
		{cards.ConsumingVolitionBlue{}, 2},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetArcaneDamageDealt(true).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s with ArcaneDamageDealt: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that co-firing runechants don't rescue a blockable variant — "this hits" reads only
// this card's own damage.
func TestConsumingVolition_RunechantsDontRescue(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetArcaneDamageDealt(true).
		AddAura(token.NewRunechant(1)).
		Build()}
	c := cards.ConsumingVolitionYellow{}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
	if got := ge.Value(); got != 3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 3 (runechant isn't 'this' damage)", got)
	}
}
