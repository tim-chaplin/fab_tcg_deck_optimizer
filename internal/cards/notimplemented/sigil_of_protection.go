// Sigil of Protection — Generic Action - Aura. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Ward 4** At the beginning of your action phase, destroy Sigil of Protection."
//
// The aura-created flag is set so same-turn aura-readers see the entry.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: ward (opponent damage prevention)

func (SigilOfProtectionRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: ward (opponent damage prevention)

func (SigilOfProtectionYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

// not implemented: ward (opponent damage prevention)

func (SigilOfProtectionBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	cards.SetAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}
