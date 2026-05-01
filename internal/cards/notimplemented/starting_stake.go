// Starting Stake — Generic Action. Cost 0, Pitch 2, Defense 3. Only printed in Yellow.
//
// Text: "If you control no Gold tokens, create a Gold token."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var startingStakeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type StartingStakeYellow struct{}

func (StartingStakeYellow) ID() ids.CardID          { return ids.StartingStakeYellow }
func (StartingStakeYellow) Name() string            { return "Starting Stake" }
func (StartingStakeYellow) Cost(*sim.TurnState) int { return 0 }
func (StartingStakeYellow) Pitch() int              { return 2 }
func (StartingStakeYellow) Attack() int             { return 0 }
func (StartingStakeYellow) Defense() int            { return 3 }
func (StartingStakeYellow) Types() card.TypeSet     { return startingStakeTypes }
func (StartingStakeYellow) GoAgain() bool           { return false }

// not implemented: gold tokens
func (StartingStakeYellow) NotImplemented()                            {}
func (StartingStakeYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
