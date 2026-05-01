// Cash In — Generic Action. Cost 4, Pitch 2, Defense 2. Only printed in Yellow.
//
// Text: "You may destroy 4 Coppers, 2 Silvers, or 1 Gold you control rather than pay Cash In's {r}
// cost. Draw 2 cards. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var cashInTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type CashInYellow struct{}

func (CashInYellow) ID() ids.CardID          { return ids.CashInYellow }
func (CashInYellow) Name() string            { return "Cash In" }
func (CashInYellow) Cost(*sim.TurnState) int { return 4 }
func (CashInYellow) Pitch() int              { return 2 }
func (CashInYellow) Attack() int             { return 0 }
func (CashInYellow) Defense() int            { return 2 }
func (CashInYellow) Types() card.TypeSet     { return cashInTypes }
func (CashInYellow) GoAgain() bool           { return true }
func (CashInYellow) NotSilverAgeLegal()      {}

// not implemented: gold/silver/copper tokens, card draw
func (CashInYellow) NotImplemented()                            {}
func (CashInYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }
