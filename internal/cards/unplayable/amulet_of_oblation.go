// Amulet of Oblation — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** **Instant** - Destroy Amulet of Oblation: Until end of turn, target attack
// action gains "If this would be put into a graveyard, instead put it on the bottom of its owner's
// deck." Activate this ability only if a card has entered a graveyard this turn."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (AmuletOfOblationBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}
