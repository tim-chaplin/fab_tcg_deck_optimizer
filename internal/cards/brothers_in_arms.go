// Brothers in Arms — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends, you may pay {r}. If you do, it gets +2{d}."
//
// Two block-time modes via sim.ModalCard + sim.BlockCost: mode 0 spends nothing for the
// printed 2{d}; mode 1 spends 1{r} for 4{d}. The chain runner enumerates both and picks
// the higher-defense mode that fits the partition's spare defense budget.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// brothersInArmsBlock fires +2{d} on mode 1 (caller already deducted 1{r} from the spare
// defense budget). Mode 0 is the printed default with no modification.
func brothersInArmsBlock(_ card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 1 {
		self.BonusDefense += 2
	}
}

func (BrothersInArmsRed) Modes() int              { return 2 }
func (BrothersInArmsRed) BlockCost(mode int8) int { return int(mode) }
func (BrothersInArmsRed) Block(s card.GameEngine, l card.Logger, self *card.CardState) {
	brothersInArmsBlock(s, l, self)
}
func (BrothersInArmsRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}

func (BrothersInArmsYellow) Modes() int              { return 2 }
func (BrothersInArmsYellow) BlockCost(mode int8) int { return int(mode) }
func (BrothersInArmsYellow) Block(s card.GameEngine, l card.Logger, self *card.CardState) {
	brothersInArmsBlock(s, l, self)
}
func (BrothersInArmsYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}

func (BrothersInArmsBlue) Modes() int              { return 2 }
func (BrothersInArmsBlue) BlockCost(mode int8) int { return int(mode) }
func (BrothersInArmsBlue) Block(s card.GameEngine, l card.Logger, self *card.CardState) {
	brothersInArmsBlock(s, l, self)
}
func (BrothersInArmsBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}
