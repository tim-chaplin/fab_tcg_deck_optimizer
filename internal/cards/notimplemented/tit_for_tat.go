// Tit for Tat — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "{t} target hero. {u} another target hero. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var titForTatTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type TitForTatBlue struct{}

func (TitForTatBlue) ID() ids.CardID          { return ids.TitForTatBlue }
func (TitForTatBlue) Name() string            { return "Tit for Tat" }
func (TitForTatBlue) Cost(*sim.TurnState) int { return 0 }
func (TitForTatBlue) Pitch() int              { return 3 }
func (TitForTatBlue) Attack() int             { return 0 }
func (TitForTatBlue) Defense() int            { return 2 }
func (TitForTatBlue) Types() card.TypeSet     { return titForTatTypes }
func (TitForTatBlue) GoAgain() bool           { return true }

// not implemented: freeze/unfreeze
func (TitForTatBlue) NotImplemented()                            {}
func (TitForTatBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
