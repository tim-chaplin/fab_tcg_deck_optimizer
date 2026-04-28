// Sapwood Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy a Frailty token you control. If you
// do, gain 1{h}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sapwoodElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SapwoodElixirRed struct{}

func (SapwoodElixirRed) ID() ids.CardID          { return ids.SapwoodElixirRed }
func (SapwoodElixirRed) Name() string            { return "Sapwood Elixir" }
func (SapwoodElixirRed) Cost(*sim.TurnState) int { return 1 }
func (SapwoodElixirRed) Pitch() int              { return 1 }
func (SapwoodElixirRed) Attack() int             { return 0 }
func (SapwoodElixirRed) Defense() int            { return 3 }
func (SapwoodElixirRed) Types() card.TypeSet     { return sapwoodElixirTypes }
func (SapwoodElixirRed) GoAgain() bool           { return true }

// not implemented: Frailty health-gain rider dropped (status tokens not tracked)
func (SapwoodElixirRed) NotImplemented() {}
func (SapwoodElixirRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.ApplyAndLogEffectiveAttack(self)
}
