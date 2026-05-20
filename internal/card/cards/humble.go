// Humble — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, they lose all hero card abilities until the end of their next
// turn."
//
// The on-hit hero-ability suppression rider is not modelled: it's an opponent-side debuff with no
// tempo or damage value to the card-vs-deck Value the simulator optimises, so Humble scores as a
// vanilla attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (HumbleRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (HumbleYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (HumbleBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
