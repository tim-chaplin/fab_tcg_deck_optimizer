// Regain Composure — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next attack this turn gets +1{p} and "When this hits, {u} your hero." **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var regainComposureTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RegainComposureBlue struct{}

func (RegainComposureBlue) ID() ids.CardID          { return ids.RegainComposureBlue }
func (RegainComposureBlue) Name() string            { return "Regain Composure" }
func (RegainComposureBlue) Cost(*sim.TurnState) int { return 0 }
func (RegainComposureBlue) Pitch() int              { return 3 }
func (RegainComposureBlue) Attack() int             { return 0 }
func (RegainComposureBlue) Defense() int            { return 2 }
func (RegainComposureBlue) Types() card.TypeSet     { return regainComposureTypes }
func (RegainComposureBlue) GoAgain() bool           { return true }

// not implemented: on-hit unfreeze rider (freeze/unfreeze state not tracked)
func (RegainComposureBlue) NotImplemented() {}
func (RegainComposureBlue) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 1)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
