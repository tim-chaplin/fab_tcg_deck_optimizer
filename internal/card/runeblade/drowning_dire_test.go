package runeblade

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestDrowningDire_NoAuraThisTurnDoesNotGrantDominate: without an aura played or created
// this turn, s.HasAuraInPlay() is false and Drowning Dire's Dominate clause doesn't fire.
// self.GrantedDominate stays false and the card is not effectively dominating.
func TestDrowningDire_NoAuraThisTurnDoesNotGrantDominate(t *testing.T) {
	cards := []card.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}}
	for _, c := range cards {
		self := &card.CardState{Card: c}
		c.Play(&card.TurnState{}, self)
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
	cards := []card.Card{DrowningDireRed{}, DrowningDireYellow{}, DrowningDireBlue{}}
	for _, c := range cards {
		self := &card.CardState{Card: c}
		s := card.TurnState{AuraCreated: true}
		c.Play(&s, self)
		if !self.GrantedDominate {
			t.Errorf("%s: GrantedDominate = false with AuraCreated, want true", c.Name())
		}
		if !self.EffectiveDominate() {
			t.Errorf("%s: EffectiveDominate = false with AuraCreated, want true", c.Name())
		}
	}
}

// TestDrowningDire_AuraPlayedThisTurnGrantsDominate: HasAuraInPlay also catches auras
// present via CardsPlayed containing an Aura-typed card, so a prior explicit Aura play (not
// just token creation) also satisfies the clause.
func TestDrowningDire_AuraPlayedThisTurnGrantsDominate(t *testing.T) {
	self := &card.CardState{Card: DrowningDireRed{}}
	s := card.TurnState{CardsPlayed: []card.Card{stubAura{}}}
	(DrowningDireRed{}).Play(&s, self)
	if !self.GrantedDominate {
		t.Error("GrantedDominate = false after aura played this turn, want true")
	}
}
