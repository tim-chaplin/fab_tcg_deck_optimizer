// Clearwater Elixir — Generic Action. Cost 1, Pitch 1, Defense 3. Only printed in Red.
//
// Text: "Your next attack this turn gets +3{p}. You may destroy a Bloodrot Pox token you control.
// If you do, gain 1{h}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var clearwaterElixirTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type ClearwaterElixirRed struct{}

func (ClearwaterElixirRed) ID() ids.CardID          { return ids.ClearwaterElixirRed }
func (ClearwaterElixirRed) Name() string            { return "Clearwater Elixir" }
func (ClearwaterElixirRed) Cost(*sim.TurnState) int { return 1 }
func (ClearwaterElixirRed) Pitch() int              { return 1 }
func (ClearwaterElixirRed) Attack() int             { return 0 }
func (ClearwaterElixirRed) Defense() int            { return 3 }
func (ClearwaterElixirRed) Types() card.TypeSet     { return clearwaterElixirTypes }
func (ClearwaterElixirRed) GoAgain() bool           { return true }

// not implemented: Bloodrot Pox health-gain rider dropped (status tokens not tracked)
func (ClearwaterElixirRed) NotImplemented() {}
func (ClearwaterElixirRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
