// Enchanting Melody — Generic Action - Aura. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue
// 3. Defense 2.
//
// Text: "**Go again** If your hero would be dealt damage, instead destroy Enchanting Melody and
// prevent 4 damage that source would deal. At the beginning of your end phase, destroy Enchanting
// Melody unless you have played a 'non-attack' action card this turn."
//
// Sets s.AuraCreated so same-turn aura-readers see the entry.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: damage-prevention trigger, end-phase destruction clause

func (EnchantingMelodyRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: damage-prevention trigger, end-phase destruction clause

func (EnchantingMelodyYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: damage-prevention trigger, end-phase destruction clause

func (EnchantingMelodyBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
