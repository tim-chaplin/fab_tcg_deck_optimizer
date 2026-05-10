// Look Tuff — Generic Action - Attack. Cost 3, Pitch 1, Power 8, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Mode 0 pays the printed 3{r} for 7{p}; mode 1 spends an extra {r} for the full 8{p}.
// See Bluster Buff for the modal-cost wiring.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func lookTuffPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

func (LookTuffRed) Modes() int              { return 2 }
func (LookTuffRed) ModalCost(mode int8) int { return 3 + int(mode) }
func (LookTuffRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	lookTuffPlay(s, l, self)
}
