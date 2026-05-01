// Brothers in Arms — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends, you may pay {r}. If you do, it gets +2{d}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var brothersInArmsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BrothersInArmsRed struct{}

func (BrothersInArmsRed) ID() ids.CardID          { return ids.BrothersInArmsRed }
func (BrothersInArmsRed) Name() string            { return "Brothers in Arms" }
func (BrothersInArmsRed) Cost(*sim.TurnState) int { return 2 }
func (BrothersInArmsRed) Pitch() int              { return 1 }
func (BrothersInArmsRed) Attack() int             { return 6 }
func (BrothersInArmsRed) Defense() int            { return 2 }
func (BrothersInArmsRed) Types() card.TypeSet     { return brothersInArmsTypes }
func (BrothersInArmsRed) GoAgain() bool           { return false }

// not implemented: pay-{r}-for-+2{d} defence rider (defence-side costs aren't solved)
func (BrothersInArmsRed) NotImplemented() {}
func (c BrothersInArmsRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type BrothersInArmsYellow struct{}

func (BrothersInArmsYellow) ID() ids.CardID          { return ids.BrothersInArmsYellow }
func (BrothersInArmsYellow) Name() string            { return "Brothers in Arms" }
func (BrothersInArmsYellow) Cost(*sim.TurnState) int { return 2 }
func (BrothersInArmsYellow) Pitch() int              { return 2 }
func (BrothersInArmsYellow) Attack() int             { return 5 }
func (BrothersInArmsYellow) Defense() int            { return 2 }
func (BrothersInArmsYellow) Types() card.TypeSet     { return brothersInArmsTypes }
func (BrothersInArmsYellow) GoAgain() bool           { return false }

// not implemented: pay-{r}-for-+2{d} defence rider (defence-side costs aren't solved)
func (BrothersInArmsYellow) NotImplemented() {}
func (c BrothersInArmsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type BrothersInArmsBlue struct{}

func (BrothersInArmsBlue) ID() ids.CardID          { return ids.BrothersInArmsBlue }
func (BrothersInArmsBlue) Name() string            { return "Brothers in Arms" }
func (BrothersInArmsBlue) Cost(*sim.TurnState) int { return 2 }
func (BrothersInArmsBlue) Pitch() int              { return 3 }
func (BrothersInArmsBlue) Attack() int             { return 4 }
func (BrothersInArmsBlue) Defense() int            { return 2 }
func (BrothersInArmsBlue) Types() card.TypeSet     { return brothersInArmsTypes }
func (BrothersInArmsBlue) GoAgain() bool           { return false }

// not implemented: pay-{r}-for-+2{d} defence rider (defence-side costs aren't solved)
func (BrothersInArmsBlue) NotImplemented() {}
func (c BrothersInArmsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
