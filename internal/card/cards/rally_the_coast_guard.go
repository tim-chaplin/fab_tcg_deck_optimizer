// Rally the Coast Guard — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this only while
// this card is defending."
//
// Block discards the first card in hand to pay for the printed +3{d}; with nothing to discard
// the bonus is skipped.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// rallyTheCoastGuardBlock discards a card for the printed +3{d}, skipping the bonus when the
// hand is empty.
func rallyTheCoastGuardBlock(ge card.GameEngine, self *card.CardState) {
	if _, ok := ge.Discard(); ok {
		self.BonusDefense += 3
	}
}

func (RallyTheCoastGuardRed) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rallyTheCoastGuardBlock(ge, self)
}
func (RallyTheCoastGuardRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (RallyTheCoastGuardYellow) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rallyTheCoastGuardBlock(ge, self)
}
func (RallyTheCoastGuardYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (RallyTheCoastGuardBlue) Block(ge card.GameEngine, l card.Logger, self *card.CardState) {
	rallyTheCoastGuardBlock(ge, self)
}
func (RallyTheCoastGuardBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
