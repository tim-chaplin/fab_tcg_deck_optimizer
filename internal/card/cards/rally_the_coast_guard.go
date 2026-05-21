// Rally the Coast Guard — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Once per Turn Instant** - Discard a card: This gets +3{d}. Activate this only while
// this card is defending."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

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
