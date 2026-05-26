package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// Tests that Play with no qualifying next attack returns zero.
func TestPrimeTheCrowd_NoAttackReturnsZero(t *testing.T) {
	ge := gameengine.New()
	for _, c := range []card.Card{cards.PrimeTheCrowdRed{}, cards.PrimeTheCrowdYellow{}, cards.PrimeTheCrowdBlue{}} {
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0", c.Name(), got)
		}
	}
}

// Tests that a non-attack action in CardsRemaining fails the rider predicate.
func TestPrimeTheCrowd_NonAttackInRemainingFizzles(t *testing.T) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{{Card: testutils.FakeRedAction()}}).Build()}
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.PrimeTheCrowdRed{}})
	if got := ge.Value(); got != 0 {
		t.Errorf("Play() = %d, want 0 (non-attack skipped)", got)
	}
}

// Tests that CrowdCheer fires the CrowdCheer trigger so subscribed "whenever the crowd
// cheers you" handlers run. CrowdBoo is verified the same way through Prime the Crowd's
// Reviled-hero path.
func TestPrimeTheCrowd_FiresCrowdCheerTrigger(t *testing.T) {
	cheers := 0
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4, TypeSet: card.NewTypeSet(card.TypeRevered)}).Build()}
	ge.CreateTrigger(testutils.FakeRedAction(), triggertype.CrowdCheer, func(_ card.GameEngine, _ card.Logger, _ card.EphemeralTrigger) {
		cheers++
	}, nil)
	ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.PrimeTheCrowdRed{}})
	if cheers != 1 {
		t.Errorf("CrowdCheer trigger fired %d times, want 1", cheers)
	}
}

// Tests that Prime the Crowd fires CrowdCheer when the active hero is Revered, CrowdBoo
// when Reviled, and skips both when the hero has neither type.
func TestPrimeTheCrowd_FiresCrowdReactionsByHeroType(t *testing.T) {
	cases := []struct {
		name      string
		heroTypes card.TypeSet
		wantCheer bool
		wantBoo   bool
	}{
		{"neither", 0, false, false},
		{"revered", card.NewTypeSet(card.TypeRevered), true, false},
		{"reviled", card.NewTypeSet(card.TypeReviled), false, true},
		{"both", card.NewTypeSet(card.TypeRevered, card.TypeReviled), true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetHero(testutils.Hero{Intel: 4, TypeSet: tc.heroTypes}).Build()}
			ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: cards.PrimeTheCrowdRed{}})
			if got := ge.HasCrowdCheered(); got != tc.wantCheer {
				t.Errorf("HasCrowdCheered = %v, want %v", got, tc.wantCheer)
			}
			if got := ge.HasCrowdBooed(); got != tc.wantBoo {
				t.Errorf("HasCrowdBooed = %v, want %v", got, tc.wantBoo)
			}
		})
	}
}

// Tests that the first attack-action card in CardsRemaining receives the per-variant bonus.
func TestPrimeTheCrowd_NextAttackReturnsBonus(t *testing.T) {
	cases := []struct {
		c    card.Card
		want int
	}{
		{cards.PrimeTheCrowdRed{}, 4},
		{cards.PrimeTheCrowdYellow{}, 3},
		{cards.PrimeTheCrowdBlue{}, 2},
	}
	for _, tc := range cases {
		target := &card.CardState{Card: testutils.FakeRedAttack()}
		ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetCardsRemaining([]*card.CardState{target}).Build()}
		ge.ResolveAttackStep(ge.Logger(), &card.CardState{Card: tc.c})
		if got := ge.Value(); got != 0 {
			t.Errorf("%s: Play() = %d, want 0 (granter returns 0; +N rides on target's BonusAttack)", tc.c.Name(), got)
		}
		if target.BonusAttack != tc.want {
			t.Errorf("%s: target BonusAttack = %d, want %d", tc.c.Name(), target.BonusAttack, tc.want)
		}
	}
}
