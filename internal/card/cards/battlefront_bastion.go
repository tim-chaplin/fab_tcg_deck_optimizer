// Battlefront Bastion — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends alone, prevent the next 1 damage that would be dealt to you this turn."
//
// Block scans ge.Defenders() for a second plain blocker; if none is present (DRs alongside
// don't count), the +1 prevention fires by bumping self.BonusDefense.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// battlefrontBastionBlock fires the +1 alone-bonus when this is the only plain blocker.
// Iterates Defenders and short-circuits on the second non-DR sighting so the typical
// partition pays at most a few comparisons.
func battlefrontBastionBlock(ge card.GameEngine, l card.Logger, self *card.CardState) {
	plainCount := 0
	for _, d := range ge.Defenders() {
		if d.Types(nil).IsDefenseReaction() {
			continue
		}
		plainCount++
		if plainCount > 1 {
			return
		}
	}
	self.BonusDefense += 1
}

func (BattlefrontBastionRed) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	battlefrontBastionBlock(ge, l, self)
}
func (BattlefrontBastionRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BattlefrontBastionYellow) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	battlefrontBastionBlock(ge, l, self)
}
func (BattlefrontBastionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BattlefrontBastionBlue) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	battlefrontBastionBlock(ge, l, self)
}
func (BattlefrontBastionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
