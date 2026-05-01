// Look Tuff — Generic Action - Attack. Cost 3, Pitch 1, Power 8, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var lookTuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LookTuffRed struct{}

func (LookTuffRed) ID() ids.CardID          { return ids.LookTuffRed }
func (LookTuffRed) Name() string            { return "Look Tuff" }
func (LookTuffRed) Cost(*sim.TurnState) int { return 3 }
func (LookTuffRed) Pitch() int              { return 1 }
func (LookTuffRed) Attack() int             { return 8 }
func (LookTuffRed) Defense() int            { return 3 }
func (LookTuffRed) Types() card.TypeSet     { return lookTuffTypes }
func (LookTuffRed) GoAgain() bool           { return false }

// not implemented: pay {r} or lose 1{p} resolved as 'always pay'
func (LookTuffRed) NotImplemented() {}
func (c LookTuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
