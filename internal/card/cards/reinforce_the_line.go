// Reinforce the Line — Generic Instant. Cost 0.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text (Red): "Target defending attack action card gains +4{d}." Yellow +3, Blue +2.
//
// Modelled as conditional damage prevention: blocking is fungible, so +N{d} on a defender
// is N less unblocked damage. Prevents N only when an attack action card is defending.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// reinforceTheLineDefend prevents n damage when a defending attack action card is present.
func reinforceTheLineDefend(ge card.GameEngine, n int, self *card.CardState) {
	for _, d := range ge.Defenders() {
		if d.Types(nil).IsAttackAction() {
			self.BonusDefense += n
			return
		}
	}
}

func (ReinforceTheLineRed) DefensiveInstant() {}
func (ReinforceTheLineRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	reinforceTheLineDefend(ge, 4, self)
}

func (ReinforceTheLineYellow) DefensiveInstant() {}
func (ReinforceTheLineYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	reinforceTheLineDefend(ge, 3, self)
}

func (ReinforceTheLineBlue) DefensiveInstant() {}
func (ReinforceTheLineBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	reinforceTheLineDefend(ge, 2, self)
}
