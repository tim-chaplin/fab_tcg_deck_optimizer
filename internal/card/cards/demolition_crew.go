// Demolition Crew — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Demolition Crew, reveal a card in your hand with cost 2 or
// greater. **Dominate**"
//
// The reveal is non-consuming — the revealed card stays in hand.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// demolitionCrewPrecondition is the shared additional-cost reveal check across all 3
// pitch variants.
func demolitionCrewPrecondition(ge card.GameEngine) bool {
	return ge.HandHasMatching(func(ge card.GameEngine, pc *card.CardState) bool {
		return pc.EffectiveCost(ge) >= 2
	})
}

func (DemolitionCrewRed) Dominate() {}
func (DemolitionCrewRed) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(ge)
}
func (c DemolitionCrewRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DemolitionCrewYellow) Dominate() {}
func (DemolitionCrewYellow) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(ge)
}
func (c DemolitionCrewYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DemolitionCrewBlue) Dominate() {}
func (DemolitionCrewBlue) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(ge)
}
func (c DemolitionCrewBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
