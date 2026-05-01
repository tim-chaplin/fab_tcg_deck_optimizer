// Restvine Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy an Inertia token you control. If
// you do, gain 1{h}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var restvineElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RestvineElixirRed struct{}

func (RestvineElixirRed) ID() ids.CardID          { return ids.RestvineElixirRed }
func (RestvineElixirRed) Name() string            { return "Restvine Elixir" }
func (RestvineElixirRed) Cost(*sim.TurnState) int { return 1 }
func (RestvineElixirRed) Pitch() int              { return 1 }
func (RestvineElixirRed) Attack() int             { return 0 }
func (RestvineElixirRed) Defense() int            { return 3 }
func (RestvineElixirRed) Types() card.TypeSet     { return restvineElixirTypes }
func (RestvineElixirRed) GoAgain() bool           { return true }

// not implemented: Inertia health-gain rider dropped (status tokens not tracked)
func (RestvineElixirRed) NotImplemented() {}
func (RestvineElixirRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
