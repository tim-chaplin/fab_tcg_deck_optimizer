// Flex — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you attack or defend with Flex, you may pay {r}{r}. If you do, it gains +2{p}."
//
// Mode 0 skips the optional cost for printed power; mode 1 pays {r}{r} for +2{p}.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func flexPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 1 {
		self.BonusAttack += 2
	}
}

func (FlexRed) Modes() int              { return 2 }
func (FlexRed) ModalCost(mode int8) int { return 2 * int(mode) }
func (FlexRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	flexPlay(ge, l, self)
}

func (FlexYellow) Modes() int              { return 2 }
func (FlexYellow) ModalCost(mode int8) int { return 2 * int(mode) }
func (FlexYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	flexPlay(ge, l, self)
}

func (FlexBlue) Modes() int              { return 2 }
func (FlexBlue) ModalCost(mode int8) int { return 2 * int(mode) }
func (FlexBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	flexPlay(ge, l, self)
}
