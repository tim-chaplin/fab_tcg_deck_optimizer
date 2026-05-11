// Bluster Buff — Generic Action - Attack. Cost 1, Pitch 1, Power 6, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Two modes via sim.ModalCard + sim.ModalCost: mode 0 pays the printed 1{r} for 5{p};
// mode 1 spends an extra {r} for the full 6{p}. The chain runner enumerates both and
// picks the higher-Value tuple per partition.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func blusterBuffPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)
}

func (BlusterBuffRed) Modes() int              { return 2 }
func (BlusterBuffRed) ModalCost(mode int8) int { return 1 + int(mode) }
func (BlusterBuffRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	blusterBuffPlay(s, l, self)
}
