// Ransack and Raze — Generic Action. Cost 0. Printed pitch variants: Blue 3. Defense 3. Go again.
//
// Text: "Destroy target landmark with cost X. Create X Gold tokens. **Go again**"
//
// Unplayable: requires a landmark target. The Silver Age format the sim models doesn't ship
// any landmark cards, so the destroy clause never has a target and the Gold mint never
// fires. The card is dead code without a landmark pool.

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (RansackAndRazeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
