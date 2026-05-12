// Restvine Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy an Inertia token you control. If
// you do, gain 1{h}. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: Inertia health-gain rider dropped (status tokens not tracked)

func (RestvineElixirRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	cards.GrantNextCardBonusAttack(g, 3, cards.IsAttack)
}
