// Blanch — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, cards they own lose all colors until the end of their next turn."
//
// The on-hit color-strip debuff isn't modelled: the single-turn solver doesn't simulate the
// opponent's next turn, so stripping colors from their cards has no measurable single-turn
// effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (BlanchRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    {}
func (BlanchYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (BlanchBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   {}
