package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/token"
)

// Tests that the discard rider stays dormant without an aura played or created.
func TestReekOfCorruption_NoAuraReturnsBaseAttack(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ReekOfCorruptionRed{}, 4},
		{cards.ReekOfCorruptionYellow{}, 3},
		{cards.ReekOfCorruptionBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		cs := &card.CardState{Card: tc.c}
		ge.ResolveChainStep(ge.Logger(), cs)
		testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s: Play() = %d, want %d (base attack, no aura)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that the discard rider fires with AuraCreated set on a likely-hit attack.
func TestReekOfCorruption_LikelyToHitWithAuraCreatedTriggersDiscard(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
	c := cards.ReekOfCorruptionRed{}
	cs := &card.CardState{Card: c}
	ge.ResolveChainStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	if got := ge.Value(); got != 4+3 {
		t.Errorf("Red with AuraCreated: Play() = %d, want 7 (base 4 likely to hit + 3 discard)", got)
	}
}

// Tests that an aura earlier in CardsPlayed satisfies the rider precondition.
func TestReekOfCorruption_AuraPlayedTriggersDiscard(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCardsPlayed([]card.Card{testutils.FakeRedAura()}).
		SetAuraCreated(true).
		Build()}
	c := cards.ReekOfCorruptionRed{}
	cs := &card.CardState{Card: c}
	ge.ResolveChainStep(ge.Logger(), cs)
	testutils.FireOnHitIfLikely(ge, ge.Logger(), cs)
	if got := ge.Value(); got != 4+3 {
		t.Errorf("Play() = %d, want %d (aura earlier in chain triggers rider)", got, 4+3)
	}
}

// Tests that the discard rider doesn't fire on blockable variants even with AuraCreated.
func TestReekOfCorruption_BlockableBaseSuppressesDiscard(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.ReekOfCorruptionYellow{}, 3},
		{cards.ReekOfCorruptionBlue{}, 2},
	}
	for _, tc := range cases {
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.want {
			t.Errorf("%s with AuraCreated: Play() = %d, want %d (blockable, no rider)", tc.c.Name(), got, tc.want)
		}
	}
}

// Tests that co-firing runechants don't rescue a blockable variant — "this hits" reads only
// this card's own damage.
func TestReekOfCorruption_RunechantsDontRescue(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetAuraCreated(true).
		AddAura(token.NewRunechant(1)).
		Build()}
	c := cards.ReekOfCorruptionYellow{}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: c})
	if got := ge.Value(); got != 3 {
		t.Errorf("Yellow with 1 Runechant: Play() = %d, want 3 (runechant isn't 'this' damage)", got)
	}
}
