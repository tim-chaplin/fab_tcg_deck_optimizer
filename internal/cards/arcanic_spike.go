// Arcanic Spike — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you've dealt arcane damage this turn, this gets +2{p}."
//
// Rider reads TurnState.ArcaneDamageDealt: when set at Play time, +2{p}; otherwise printed
// attack alone.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var arcanicSpikeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// arcaneDamageBonus is the +2{p} gained when the "dealt arcane damage this turn" clause is live.
const arcaneDamageBonus = 2

// arcanicSpikeBonus returns the +2{p} power buff when ArcaneDamageDealt is set, else 0.
func arcanicSpikeBonus(s *sim.TurnState) int {
	if s != nil && s.ArcaneDamageDealt {
		return arcaneDamageBonus
	}
	return 0
}

type ArcanicSpikeRed struct{}

func (ArcanicSpikeRed) ID() ids.CardID          { return ids.ArcanicSpikeRed }
func (ArcanicSpikeRed) Name() string            { return "Arcanic Spike" }
func (ArcanicSpikeRed) Cost(*sim.TurnState) int { return 2 }
func (ArcanicSpikeRed) Pitch() int              { return 1 }
func (ArcanicSpikeRed) Attack() int             { return 5 }
func (ArcanicSpikeRed) Defense() int            { return 3 }
func (ArcanicSpikeRed) Types() card.TypeSet     { return arcanicSpikeTypes }
func (ArcanicSpikeRed) GoAgain() bool           { return false }
func (ArcanicSpikeRed) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += arcanicSpikeBonus(s)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type ArcanicSpikeYellow struct{}

func (ArcanicSpikeYellow) ID() ids.CardID          { return ids.ArcanicSpikeYellow }
func (ArcanicSpikeYellow) Name() string            { return "Arcanic Spike" }
func (ArcanicSpikeYellow) Cost(*sim.TurnState) int { return 2 }
func (ArcanicSpikeYellow) Pitch() int              { return 2 }
func (ArcanicSpikeYellow) Attack() int             { return 4 }
func (ArcanicSpikeYellow) Defense() int            { return 3 }
func (ArcanicSpikeYellow) Types() card.TypeSet     { return arcanicSpikeTypes }
func (ArcanicSpikeYellow) GoAgain() bool           { return false }
func (ArcanicSpikeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += arcanicSpikeBonus(s)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type ArcanicSpikeBlue struct{}

func (ArcanicSpikeBlue) ID() ids.CardID          { return ids.ArcanicSpikeBlue }
func (ArcanicSpikeBlue) Name() string            { return "Arcanic Spike" }
func (ArcanicSpikeBlue) Cost(*sim.TurnState) int { return 2 }
func (ArcanicSpikeBlue) Pitch() int              { return 3 }
func (ArcanicSpikeBlue) Attack() int             { return 3 }
func (ArcanicSpikeBlue) Defense() int            { return 3 }
func (ArcanicSpikeBlue) Types() card.TypeSet     { return arcanicSpikeTypes }
func (ArcanicSpikeBlue) GoAgain() bool           { return false }
func (ArcanicSpikeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	self.BonusAttack += arcanicSpikeBonus(s)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
