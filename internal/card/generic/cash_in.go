// Cash In — Generic Action. Cost 4, Pitch 2, Defense 2. Only printed in Yellow.
//
// Text: "You may destroy 4 Coppers, 2 Silvers, or 1 Gold you control rather than pay Cash In's {r}
// cost. Draw 2 cards. **Go again**"
//
// Simplification: Activated Gold/Silver/Copper economy and draws aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var cashInTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type CashInYellow struct{}

func (CashInYellow) ID() card.ID                 { return card.CashInYellow }
func (CashInYellow) Name() string                { return "Cash In (Yellow)" }
func (CashInYellow) Cost(*card.TurnState) int                   { return 4 }
func (CashInYellow) Pitch() int                  { return 2 }
func (CashInYellow) Attack() int                 { return 0 }
func (CashInYellow) Defense() int                { return 2 }
func (CashInYellow) Types() card.TypeSet         { return cashInTypes }
func (CashInYellow) GoAgain() bool               { return true }
func (CashInYellow) NotSilverAgeLegal()           {}
func (CashInYellow) Play(s *card.TurnState, _ *card.CardState) int { return 0 }
