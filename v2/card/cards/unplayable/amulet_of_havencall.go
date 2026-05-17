// Amulet of Havencall — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Defense Reaction** - Destroy Amulet of Havencall: Search your deck for a
// card named Rally the Rearguard, add it to this chain link as a defending card, then shuffle.
// Activate this ability only if you have no cards in hand."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (AmuletOfHavencallBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
