// Amulet of Assertiveness — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Attack Reaction** - Destroy Amulet of Assertiveness: Target attack gains
// "When this hits, banish the top card of your deck. If it's an attack action card, you may play it
// this turn." Activate this ability only if you have 4 or more cards in hand."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (AmuletOfAssertivenessYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
