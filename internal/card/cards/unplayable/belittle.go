// Belittle — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Belittle, you may reveal an attack action card with 3 or
// less base {p} from your hand. If you do, search your deck for a card named Minnowism, reveal it,
// put it into your hand, then shuffle your deck. **Go again**"

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (BelittleRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BelittleYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (BelittleBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
