// Lay Low — Generic Defense Reaction. Cost 0, Pitch 2, Defense 3. Only printed in Yellow.
// Text: "If you are marked, you can't play this. If the defending hero is marked, their next attack
// this turn gets -1{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

type LayLowYellow struct{}

func (LayLowYellow) ID() ids.CardID          { return ids.LayLowYellow }
func (LayLowYellow) Name() string            { return "Lay Low" }
func (LayLowYellow) Cost(*sim.TurnState) int { return 0 }
func (LayLowYellow) Pitch() int              { return 2 }
func (LayLowYellow) Attack() int             { return 0 }
func (LayLowYellow) Defense() int            { return 3 }
func (LayLowYellow) Types() card.TypeSet     { return defenseReactionTypes }
func (LayLowYellow) GoAgain() bool           { return false }

// not implemented: marked-defender state not tracked; treated as always legal and the -1{p}
// attacker debuff is dropped
func (LayLowYellow) NotImplemented() {}
func (LayLowYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
