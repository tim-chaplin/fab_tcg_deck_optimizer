// Bluster Buff — Generic Action - Attack. Cost 1, Pitch 1, Power 6, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var blusterBuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BlusterBuffRed struct{}

func (BlusterBuffRed) ID() ids.CardID          { return ids.BlusterBuffRed }
func (BlusterBuffRed) Name() string            { return "Bluster Buff" }
func (BlusterBuffRed) Cost(*sim.TurnState) int { return 1 }
func (BlusterBuffRed) Pitch() int              { return 1 }
func (BlusterBuffRed) Attack() int             { return 6 }
func (BlusterBuffRed) Defense() int            { return 3 }
func (BlusterBuffRed) Types() card.TypeSet     { return blusterBuffTypes }
func (BlusterBuffRed) GoAgain() bool           { return false }

// not implemented: pay {r} or lose 1{p} resolved as 'always pay'
func (BlusterBuffRed) NotImplemented() {}
func (c BlusterBuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
