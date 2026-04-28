package cards

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

func TestHitTheHighNotes_NoAuraReturnsBase(t *testing.T) {
	// Neither an aura played nor one created this turn → no bonus, just printed power.
	cases := []struct {
		c    sim.Card
		base int
	}{
		{HitTheHighNotesRed{}, 4},
		{HitTheHighNotesYellow{}, 3},
		{HitTheHighNotesBlue{}, 2},
	}
	for _, tc := range cases {
		var s sim.TurnState
		tc.c.Play(&s, &sim.CardState{Card: tc.c})
		if got := s.Value; got != tc.base {
			t.Errorf("%s: Play() = %d, want %d", tc.c.Name(), got, tc.base)
		}
	}
}

func TestHitTheHighNotes_AuraPlayedTriggersBonus(t *testing.T) {
	// An Aura-typed card earlier in the turn's CardsPlayed → +2 power.
	s := sim.TurnState{CardsPlayed: []sim.Card{testutils.Aura{}}}
	(HitTheHighNotesRed{}).Play(&s, &sim.CardState{Card: HitTheHighNotesRed{}})
	if got := s.Value; got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 aura bonus)", got)
	}
}

func TestHitTheHighNotes_AuraCreatedTriggersBonus(t *testing.T) {
	// AuraCreated flag set earlier in the chain (e.g. Runechant creation) → +2 power, even
	// without an Aura-typed card in CardsPlayed.
	s := sim.TurnState{AuraCreated: true}
	(HitTheHighNotesRed{}).Play(&s, &sim.CardState{Card: HitTheHighNotesRed{}})
	if got := s.Value; got != 6 {
		t.Errorf("Play() = %d, want 6 (base 4 + 2 AuraCreated bonus)", got)
	}
}

// TestHitTheHighNotes_BonusFlowsThroughBonusAttack: the +2{p} rider is a power buff, not a
// damage rider — it must land on self.BonusAttack so EffectiveAttack and LikelyToHit see
// the buffed power. A 4-power Red with the rider becomes a 6-power attack, which falls
// outside the {1,4,7} likely-to-hit window; on-hit triggers from sibling cards (Mauvrion
// Skies's "if this hits, create N runechants") must not fire on a 6-power attack.
func TestHitTheHighNotes_BonusFlowsThroughBonusAttack(t *testing.T) {
	s := sim.TurnState{AuraCreated: true}
	self := &sim.CardState{Card: HitTheHighNotesRed{}}
	(HitTheHighNotesRed{}).Play(&s, self)
	if got := self.EffectiveAttack(); got != 6 {
		t.Errorf("EffectiveAttack() = %d, want 6 (base 4 + 2 power buff)", got)
	}
	if sim.LikelyToHit(self) {
		t.Errorf("LikelyToHit = true at EffectiveAttack 6; want false (6 ∉ {1,4,7})")
	}
}
