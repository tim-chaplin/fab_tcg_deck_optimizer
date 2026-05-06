package notimplemented

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestDrowningDire_NoAuraThisTurnDoesNotGrantDominate: without an aura played or created
// this turn, s.HasPlayedOrCreatedAura() is false and Drowning Dire's Dominate clause doesn't fire.
// self.GrantedDominate stays false and the card is not effectively dominating.
func TestDrowningDire_NoAuraThisTurnDoesNotGrantDominate(t *testing.T) {
	cards := []sim.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}}
	for _, c := range cards {
		self := &sim.CardState{Card: c}
		c.Play(&sim.TurnState{}, self)
		if self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = true without aura, want false", c.Name())
		}
		if self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = true without aura, want false", c.Name())
		}
	}
}

// TestDrowningDire_AuraCreatedThisTurnGrantsDominate: when an aura is in play (s.AuraCreated
// flipped by a prior Aura-making card earlier in the chain), the conditional clause fires and
// Play sets self.GrantedDominate.
func TestDrowningDire_AuraCreatedThisTurnGrantsDominate(t *testing.T) {
	cards := []sim.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}}
	for _, c := range cards {
		self := &sim.CardState{Card: c}
		s := sim.TurnState{AuraCreated: true}
		c.Play(&s, self)
		if !self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = false with AuraCreated, want true", c.Name())
		}
		if !self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = false with AuraCreated, want true", c.Name())
		}
	}
}

// TestDrowningDire_AuraPlayedThisTurnGrantsDominate: HasPlayedOrCreatedAura also catches auras
// present via CardsPlayed containing an Aura-typed card, so a prior explicit Aura play (not
// just token creation) also satisfies the clause.
func TestDrowningDire_AuraPlayedThisTurnGrantsDominate(t *testing.T) {
	self := &sim.CardState{Card: DrowningDireRed{}}
	s := sim.TurnState{CardsPlayed: []sim.Card{testutils.Aura{}}}
	(DrowningDireRed{}).Play(&s, self)
	if !self.GrantedDominate {
		t.Error("GrantedDominate = false after aura played this turn, want true")
	}
}

// Tests that every Drowning Dire variant carries sim.NotImplemented (on-hit recycle rider
// is not modelled).
func TestDrowningDire_NotImplemented(t *testing.T) {
	cards := []sim.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}}
	for _, c := range cards {
		if _, ok := c.(sim.NotImplemented); !ok {
			t.Errorf("%s: missing sim.NotImplemented marker for the on-hit recycle rider", c.Name())
		}
	}
}
