// Muscle Mutt — Generic Action - Attack. Cost 3, Pitch 2, Power 6, Defense 2. Only printed in
// Yellow.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var muscleMuttTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type MuscleMuttYellow struct{}

func (MuscleMuttYellow) ID() ids.CardID          { return ids.MuscleMuttYellow }
func (MuscleMuttYellow) Name() string            { return "Muscle Mutt" }
func (MuscleMuttYellow) Cost(*sim.TurnState) int { return 3 }
func (MuscleMuttYellow) Pitch() int              { return 2 }
func (MuscleMuttYellow) Attack() int             { return 6 }
func (MuscleMuttYellow) Defense() int            { return 2 }
func (MuscleMuttYellow) Types() card.TypeSet     { return muscleMuttTypes }
func (MuscleMuttYellow) GoAgain() bool           { return false }
func (c MuscleMuttYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
