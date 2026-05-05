// Starting Stake — Generic Action. Cost 0, Pitch 2, Defense 3. Yellow only.
// Text: "If you control no Gold tokens, create a Gold token."

package cards

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

func (StartingStakeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	if s.Gold() == 0 {
		s.CreateGold(1)
		s.LogRider(self, 0, "Created a gold token")
	}
	s.Log(self, 0)
}
