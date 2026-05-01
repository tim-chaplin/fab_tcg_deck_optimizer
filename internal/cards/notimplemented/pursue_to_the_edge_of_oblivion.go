// Pursue to the Edge of Oblivion — Generic Action - Attack. Cost 0, Pitch 1, Power 4, Defense 3.
// Only printed in Red.
//
// Text: "When this hits a hero, **mark** them."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var pursueToTheEdgeOfOblivionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type PursueToTheEdgeOfOblivionRed struct{}

func (PursueToTheEdgeOfOblivionRed) ID() ids.CardID          { return ids.PursueToTheEdgeOfOblivionRed }
func (PursueToTheEdgeOfOblivionRed) Name() string            { return "Pursue to the Edge of Oblivion" }
func (PursueToTheEdgeOfOblivionRed) Cost(*sim.TurnState) int { return 0 }
func (PursueToTheEdgeOfOblivionRed) Pitch() int              { return 1 }
func (PursueToTheEdgeOfOblivionRed) Attack() int             { return 4 }
func (PursueToTheEdgeOfOblivionRed) Defense() int            { return 3 }
func (PursueToTheEdgeOfOblivionRed) Types() card.TypeSet     { return pursueToTheEdgeOfOblivionTypes }
func (PursueToTheEdgeOfOblivionRed) GoAgain() bool           { return false }

// not implemented: on-hit mark
func (PursueToTheEdgeOfOblivionRed) NotImplemented() {}
func (PursueToTheEdgeOfOblivionRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
