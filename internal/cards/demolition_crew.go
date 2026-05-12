// Demolition Crew — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Demolition Crew, reveal a card in your hand with cost 2 or
// greater. **Dominate**"
//
// The reveal is non-consuming — the revealed card stays in hand.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// demolitionCrewPrecondition is the shared additional-cost check across all 3 pitch
// variants.
func demolitionCrewPrecondition(g card.GameEngine) bool {
	for _, c := range g.Hand() {
		if c.Cost(g) >= 2 {
			return true
		}
	}
	return false
}

func (DemolitionCrewRed) Dominate() {}
func (DemolitionCrewRed) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(g)
}
func (c DemolitionCrewRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DemolitionCrewYellow) Dominate() {}
func (DemolitionCrewYellow) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(g)
}
func (c DemolitionCrewYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DemolitionCrewBlue) Dominate() {}
func (DemolitionCrewBlue) PlayPrecondition(g card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(g)
}
func (c DemolitionCrewBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
