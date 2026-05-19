// Rise Above — Generic Defense Reaction. Cost 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
//
// Text: "You may put a card from your hand on top of your deck rather than pay Rise Above's
// {r} cost."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// riseAbovePrintedCost is the un-discounted resource cost (also the VariableCost MaxCost bound).
const riseAbovePrintedCost = 2

func riseAboveCost(ge card.GameEngine) int {
	if ge != nil && len(ge.Hand()) > 0 {
		return 0
	}
	return riseAbovePrintedCost
}

func riseAbovePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if len(ge.Hand()) == 0 {
		return
	}
	returned := ge.PopHandAt(0)
	ge.PrependToDeck(returned)
	l.AppendPostTriggerf(self.Card.DisplayName(), 0, "Returned %s to top of deck (alt cost)", returned.DisplayName())
}

func (RiseAboveRed) Cost(ge card.GameEngine) int { return riseAboveCost(ge) }
func (RiseAboveRed) MinCost() int                { return 0 }
func (RiseAboveRed) MaxCost() int                { return riseAbovePrintedCost }
func (RiseAboveRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	riseAbovePlay(ge, l, self)
}

func (RiseAboveYellow) Cost(ge card.GameEngine) int { return riseAboveCost(ge) }
func (RiseAboveYellow) MinCost() int                { return 0 }
func (RiseAboveYellow) MaxCost() int                { return riseAbovePrintedCost }
func (RiseAboveYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	riseAbovePlay(ge, l, self)
}

func (RiseAboveBlue) Cost(ge card.GameEngine) int { return riseAboveCost(ge) }
func (RiseAboveBlue) MinCost() int                { return 0 }
func (RiseAboveBlue) MaxCost() int                { return riseAbovePrintedCost }
func (RiseAboveBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	riseAbovePlay(ge, l, self)
}
