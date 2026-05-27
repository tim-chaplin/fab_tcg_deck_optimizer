// Put in Context — Generic Defense Reaction. Cost 0. Printed pitch variants: Blue 3. Defense 3.
//
// Text: "This can only defend an attack with 3 or less base {p}."
//
// Unplayable: the base-power-≤3 attack-side gate is matchup-dependent. Across most
// matchups, attack-action cards routinely print 4+ power, so the block fizzles too often
// to model meaningfully without modelling the opponent's attack-power distribution.

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (PutInContextBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
