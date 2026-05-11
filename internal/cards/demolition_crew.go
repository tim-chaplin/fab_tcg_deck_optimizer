// Demolition Crew — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Demolition Crew, reveal a card in your hand with cost 2 or
// greater. **Dominate**"
//
// The reveal is non-consuming — the revealed card stays in hand. The "reveal a cost-2+ card"
// clause is the additional cost, modelled via the PlayPrecondition opt-in: with no eligible
// reveal target the chain runner aborts the permutation as illegal (FaB rule: a card whose
// costs you can't fully pay is illegal to declare).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// demolitionCrewPrecondition is the shared additional-cost check across all 3 pitch
// variants. The chain runner's hand snapshot has already removed the playing card and
// popped this card's pitches by the time PlayPrecondition runs, so the scan only sees
// cards that genuinely remain in hand.
func demolitionCrewPrecondition(s card.GameEngine) bool {
	for _, c := range s.Hand() {
		if c.Cost(s) >= 2 {
			return true
		}
	}
	return false
}

func (DemolitionCrewRed) Dominate() {}
func (DemolitionCrewRed) PlayPrecondition(s card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(s)
}
func (c DemolitionCrewRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DemolitionCrewYellow) Dominate() {}
func (DemolitionCrewYellow) PlayPrecondition(s card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(s)
}
func (c DemolitionCrewYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DemolitionCrewBlue) Dominate() {}
func (DemolitionCrewBlue) PlayPrecondition(s card.GameEngine, _ *card.CardState) bool {
	return demolitionCrewPrecondition(s)
}
func (c DemolitionCrewBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}
