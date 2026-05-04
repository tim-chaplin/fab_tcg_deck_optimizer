// Cash In — Generic Action. Cost 4, Pitch 2, Defense 2. Yellow only.
// Text: "You may destroy 4 Coppers, 2 Silvers, or 1 Gold you control rather than pay
// Cash In's {r} cost. Draw 2 cards. **Go again**"
//
// Token-destroy alt-cost isn't modelled; always resolves the printed {4} pay path.

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

func (CashInYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.DrawOne()
	s.DrawOne()
	s.Log(self, 0)
	s.LogRider(self, 0, "Drew 2 cards")
}
