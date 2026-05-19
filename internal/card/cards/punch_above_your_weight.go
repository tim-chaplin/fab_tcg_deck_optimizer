// Punch Above Your Weight — Generic Action - Attack. Cost 0. Printed power: Red 2, Yellow 2,
// Blue 2. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, you may pay {r}{r}{r}. If you do, this gets +5{p}."
//
// Mode 0 skips the optional cost for 2{p}; mode 1 pays {r}{r}{r} for 7{p}.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func punchAboveYourWeightPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 1 {
		self.BonusAttack += 5
	}
}

func (PunchAboveYourWeightRed) Modes() int              { return 2 }
func (PunchAboveYourWeightRed) ModalCost(mode int8) int { return 3 * int(mode) }
func (PunchAboveYourWeightRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	punchAboveYourWeightPlay(ge, l, self)
}

func (PunchAboveYourWeightYellow) Modes() int              { return 2 }
func (PunchAboveYourWeightYellow) ModalCost(mode int8) int { return 3 * int(mode) }
func (PunchAboveYourWeightYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	punchAboveYourWeightPlay(ge, l, self)
}

func (PunchAboveYourWeightBlue) Modes() int              { return 2 }
func (PunchAboveYourWeightBlue) ModalCost(mode int8) int { return 3 * int(mode) }
func (PunchAboveYourWeightBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	punchAboveYourWeightPlay(ge, l, self)
}
