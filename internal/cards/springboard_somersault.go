// Springboard Somersault — Generic Defense Reaction. Cost 0, Pitch 2, Defense 2. Only printed in
// Yellow.
// Text: "If Springboard Somersault is played from arsenal, it gains +2{d}."
//
// +2{d} when played from arsenal via card.ArsenalDefenseBonus (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SpringboardSomersaultYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
}
func (SpringboardSomersaultYellow) ArsenalDefenseBonus() int { return 2 }
