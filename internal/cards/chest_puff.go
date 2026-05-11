// Chest Puff — Generic Action - Attack. Cost 2, Pitch 1, Power 7, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Mode 0 pays the printed 2{r} for 6{p}; mode 1 spends an extra {r} for the full 7{p}.
// See Bluster Buff for the modal-cost wiring.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func chestPuffPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (ChestPuffRed) Modes() int              { return 2 }
func (ChestPuffRed) ModalCost(mode int8) int { return 2 + int(mode) }
func (ChestPuffRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	chestPuffPlay(s, l, self)
}
