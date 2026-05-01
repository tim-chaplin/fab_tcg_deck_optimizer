// Pursue to the Pits of Despair — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var pursueToThePitsOfDespairTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PursueToThePitsOfDespairRed struct{}

func (PursueToThePitsOfDespairRed) ID() ids.CardID          { return ids.PursueToThePitsOfDespairRed }
func (PursueToThePitsOfDespairRed) Name() string            { return "Pursue to the Pits of Despair" }
func (PursueToThePitsOfDespairRed) Cost(*sim.TurnState) int { return 1 }
func (PursueToThePitsOfDespairRed) Pitch() int              { return 1 }
func (PursueToThePitsOfDespairRed) Attack() int             { return 5 }
func (PursueToThePitsOfDespairRed) Defense() int            { return 3 }
func (PursueToThePitsOfDespairRed) Types() card.TypeSet     { return pursueToThePitsOfDespairTypes }
func (PursueToThePitsOfDespairRed) GoAgain() bool           { return false }

// not implemented: on-hit mark
func (PursueToThePitsOfDespairRed) NotImplemented() {}
func (PursueToThePitsOfDespairRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
