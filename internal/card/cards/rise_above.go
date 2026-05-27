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

// riseAboveAlternativeCost reports the alt branch as (0, ok=true) when a Held card remains
// in hand to push onto the deck top.
func riseAboveAlternativeCost(ge card.GameEngine) (int, bool) {
	if ge != nil && ge.HeldHandSize() > 0 {
		return 0, true
	}
	return 0, false
}

func riseAbovePlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.PaidAlternativeCost {
		ge.DiscardToTopOfDeck(self.Card.DisplayName())
	}
}

func (RiseAboveRed) AlternativeCost(ge card.GameEngine) (int, bool) {
	return riseAboveAlternativeCost(ge)
}
func (RiseAboveRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	riseAbovePlay(ge, l, self)
}

func (RiseAboveYellow) AlternativeCost(ge card.GameEngine) (int, bool) {
	return riseAboveAlternativeCost(ge)
}
func (RiseAboveYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	riseAbovePlay(ge, l, self)
}

func (RiseAboveBlue) AlternativeCost(ge card.GameEngine) (int, bool) {
	return riseAboveAlternativeCost(ge)
}
func (RiseAboveBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	riseAbovePlay(ge, l, self)
}
