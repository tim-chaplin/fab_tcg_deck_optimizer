package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestHitTheHighNotes_NoAuraReturnsBase(t *testing.T) {
	// Neither an aura played nor one created this turn → no bonus, just printed power.
	cases := []struct {
		c    card.Card
		base int
	}{
		{cards.HitTheHighNotesRed{}, 4},
		{cards.HitTheHighNotesYellow{}, 3},
		{cards.HitTheHighNotesBlue{}, 2},
	}
	for _, tc := range cases {
		ge := gameengine.New()
		ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
	}
}

func TestHitTheHighNotes_AuraPlayedTriggersBonus(t *testing.T) {
	// An Aura-typed card earlier in the turn'ge CardsPlayed → +2 power.
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetCardsPlayed([]card.Card{testutils.FakeRedAura()}).
		SetAuraCreated(true).
		Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.HitTheHighNotesRed{}})
	if got := ge.Value(); got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 aura bonus)", got)
	}
}

func TestHitTheHighNotes_AuraCreatedTriggersBonus(t *testing.T) {
	// AuraCreated flag set earlier in the chain (e.g. Runechant creation) → +2 power, even
	// without an Aura-typed card in CardsPlayed.
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
	ge.ResolveChainStep(ge.Logger(), &card.CardState{Card: cards.HitTheHighNotesRed{}})
	if got := ge.Value(); got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 AuraCreated bonus)", got)
	}
}

// Tests that the +2{p} rider flows through pc.BonusAttack so EffectiveAttack and
// LikelyToHit see the buffed power.
func TestHitTheHighNotes_BonusFlowsThroughBonusAttack(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetAuraCreated(true).Build()}
	pc := &card.CardState{Card: cards.HitTheHighNotesRed{}}
	ge.ResolveChainStep(ge.Logger(), pc)
	if got := pc.EffectiveAttack(); got != 6 {
		t.Errorf("EffectiveAttack() = %d, want 6 (base 4 + 2 power buff)", got)
	}
	if ge.LikelyToHit(pc) {
		t.Errorf("LikelyToHit = true at EffectiveAttack 6; want false (6 ∉ {1,4,7})")
	}
}

// Tests that playing Hit the High Notes with a Malefic Incantation already in play turns on
// its "played or created an aura this turn" rider: Malefic's play-triggered Runechant is
// created before Hit the High Notes resolves. The Runechant must also survive — it is not
// consumed by the very attack whose play triggered its creation.
func TestHitTheHighNotes_SeesRunechantFromTriggeredMalefic(t *testing.T) {
	prior := gameengine.GameStateBuilder().
		CreateAuraFromCard(cards.MaleficIncantationRed{}).
		SetIncomingDamage(0).
		Build()
	d := deck.New(testutils.Hero{Intel: 4}, nil, nil)
	hand := []card.Card{testutils.FakeBlueResource(), cards.HitTheHighNotesRed{}}

	summary := sim.EvalOneTurnForTesting(d, prior, hand)

	if summary.Value != 7 {
		t.Fatalf("Value = %d, want 7 (Hit the High Notes 6 = 4 base + 2 aura rider, plus 1 for Malefic's Runechant)\nBestLine: %s",
			summary.Value, formatBestLine(summary.BestLine))
	}
	if got := summary.State.RunechantCount(); got != 1 {
		t.Fatalf("end-of-turn RunechantCount = %d, want 1 (Malefic's Runechant must survive — not consumed by the attack that triggered it)", got)
	}
}
