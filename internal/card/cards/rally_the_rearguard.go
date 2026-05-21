// Rally the Rearguard — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this ability only
// while this is defending."
//
// Block discards the first card in hand to pay for the printed +3{d}; with nothing to discard
// the bonus is skipped.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// rallyTheRearguardBlock discards a card for the printed +3{d}, skipping the bonus when the
// hand is empty.
func rallyTheRearguardBlock(ge card.GameEngine, self *card.CardState) {
	if _, ok := ge.Discard(); ok {
		self.BonusDefense += 3
	}
}

func (RallyTheRearguardRed) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rallyTheRearguardBlock(ge, self)
}
func (RallyTheRearguardRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (RallyTheRearguardYellow) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rallyTheRearguardBlock(ge, self)
}
func (RallyTheRearguardYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (RallyTheRearguardBlue) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rallyTheRearguardBlock(ge, self)
}
func (RallyTheRearguardBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
