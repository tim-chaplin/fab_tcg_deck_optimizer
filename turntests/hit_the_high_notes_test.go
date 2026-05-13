package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
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
		s := gameengine.NewFromState(nil)
		s.ResolveChainStep(s.Logger(), &card.CardState{Card: tc.c})
		if got := s.Value(); got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
	}
}

func TestHitTheHighNotes_AuraPlayedTriggersBonus(t *testing.T) {
	// An Aura-typed card earlier in the turn's CardsPlayed → +2 power.
	s := gameengine.NewFromSpec(gameengine.Spec{CardsPlayed: []card.Card{testutils.Aura{}}, AuraCreated: true})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.HitTheHighNotesRed{}})
	if got := s.Value(); got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 aura bonus)", got)
	}
}

func TestHitTheHighNotes_AuraCreatedTriggersBonus(t *testing.T) {
	// AuraCreated flag set earlier in the chain (e.g. Runechant creation) → +2 power, even
	// without an Aura-typed card in CardsPlayed.
	s := gameengine.NewFromSpec(gameengine.Spec{AuraCreated: true})
	s.ResolveChainStep(s.Logger(), &card.CardState{Card: cards.HitTheHighNotesRed{}})
	if got := s.Value(); got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 AuraCreated bonus)", got)
	}
}

// Tests that the +2{p} rider flows through self.BonusAttack so EffectiveAttack and
// LikelyToHit see the buffed power.
func TestHitTheHighNotes_BonusFlowsThroughBonusAttack(t *testing.T) {
	s := gameengine.NewFromSpec(gameengine.Spec{AuraCreated: true})
	self := &card.CardState{Card: cards.HitTheHighNotesRed{}}
	s.ResolveChainStep(s.Logger(), self)
	if got := self.EffectiveAttack(); got != 6 {
		t.Errorf("EffectiveAttack() = %d, want 6 (base 4 + 2 power buff)", got)
	}
	if s.LikelyToHit(self) {
		t.Errorf("LikelyToHit = true at EffectiveAttack 6; want false (6 ∉ {1,4,7})")
	}
}
