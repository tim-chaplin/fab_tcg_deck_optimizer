// Walk the Plank — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a Pirate hero, {t} them or an ally they control."
//
// The on-hit pirate-target tap rider isn't modelled: opposing-hero class isn't tracked, and tap
// state has no measurable single-turn effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (WalkThePlankRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    {}
func (WalkThePlankYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (WalkThePlankBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   {}
