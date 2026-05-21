// Restvine Elixir — Generic Action. Cost 1, Defense 3. Printed only in Red (pitch 1).
// Text: "Your next attack this turn gets +3{p}. You may destroy an Inertia token you
// control. If you do, gain 1{h}. **Go again**"
//
// The "destroy an Inertia token" clause is not modelled — that token is placed only by
// specific opponent cards, a sideboard/matchup consideration the single-turn model omits.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (RestvineElixirRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 3, card.IsAttack)
}
