// Fact-Finding Mission — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, you may look at a face-down card in their arsenal or equipment
// zones."
//
// The on-hit peek rider isn't modelled: the single-turn solver doesn't track the opponent's
// arsenal or equipment, so looking at a face-down card there has no measurable single-turn
// effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (FactFindingMissionRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    {}
func (FactFindingMissionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (FactFindingMissionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   {}
