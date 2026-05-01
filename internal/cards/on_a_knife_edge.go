// On a Knife Edge — Generic Action. Cost 0, Pitch 2, Defense 2. Only printed in Yellow.
//
// Text: "Your next sword attack this turn gains **go again**. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var onAKnifeEdgeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type OnAKnifeEdgeYellow struct{}

func (OnAKnifeEdgeYellow) ID() ids.CardID          { return ids.OnAKnifeEdgeYellow }
func (OnAKnifeEdgeYellow) Name() string            { return "On a Knife Edge" }
func (OnAKnifeEdgeYellow) Cost(*sim.TurnState) int { return 0 }
func (OnAKnifeEdgeYellow) Pitch() int              { return 2 }
func (OnAKnifeEdgeYellow) Attack() int             { return 0 }
func (OnAKnifeEdgeYellow) Defense() int            { return 2 }
func (OnAKnifeEdgeYellow) Types() card.TypeSet     { return onAKnifeEdgeTypes }
func (OnAKnifeEdgeYellow) GoAgain() bool           { return true }

// not implemented: next-sword-attack go-again grant (weapon chain not scanned)
func (OnAKnifeEdgeYellow) NotImplemented()                            {}
func (OnAKnifeEdgeYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
