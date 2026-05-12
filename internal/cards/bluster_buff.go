// Bluster Buff — Generic Action - Attack. Cost 1, Pitch 1, Power 6, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Mode 0 pays the printed 1{r} for 5{p}; mode 1 spends an extra {r} for the full 6{p}.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func blusterBuffPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
}

func (BlusterBuffRed) Modes() int              { return 2 }
func (BlusterBuffRed) ModalCost(mode int8) int { return 1 + int(mode) }
func (BlusterBuffRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	blusterBuffPlay(g, l, self)
}
