// Right Behind You — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends together with another card from hand, this gets +1{d} and you may look
// at the top card of your deck. You may put it on the bottom."
//
// Block scans ge.Defenders() for a second plain blocker; if at least one is present (DRs
// alongside don't count), the +1{d} fires by bumping self.BonusDefense. The deck-top
// peek/bottom rider is dropped — the optimizer would never bottom a card it'd rather
// draw next.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// rightBehindYouBlock fires the +1{d} together-bonus when at least two plain blockers
// share the defenders slot. Short-circuits on the second non-DR sighting.
func rightBehindYouBlock(ge card.GameEngine, l card.Logger, self *card.CardState) {
	plainCount := 0
	for _, d := range ge.Defenders() {
		if d.Types(nil).IsDefenseReaction() {
			continue
		}
		plainCount++
		if plainCount >= 2 {
			self.BonusDefense += 1
			return
		}
	}
}

func (RightBehindYouRed) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rightBehindYouBlock(ge, l, self)
}
func (RightBehindYouRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (RightBehindYouYellow) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rightBehindYouBlock(ge, l, self)
}
func (RightBehindYouYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (RightBehindYouBlue) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rightBehindYouBlock(ge, l, self)
}
func (RightBehindYouBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
