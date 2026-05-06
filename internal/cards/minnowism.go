// Minnowism — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with 3 or less base {p} you play this turn gains +N{p}. **Go
// again**" (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var minnowismTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// minnowismIsTarget gates the rider on attack action cards with printed power 3 or less.
func minnowismIsTarget(_ *sim.TurnState, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttackAction() && pc.Card.Attack() <= 3
}

type MinnowismRed struct{}

func (MinnowismRed) ID() ids.CardID          { return ids.MinnowismRed }
func (MinnowismRed) Name() string            { return "Minnowism" }
func (MinnowismRed) Cost(*sim.TurnState) int { return 0 }
func (MinnowismRed) Pitch() int              { return 1 }
func (MinnowismRed) Attack() int             { return 0 }
func (MinnowismRed) Defense() int            { return 2 }
func (MinnowismRed) Types() card.TypeSet     { return minnowismTypes }
func (MinnowismRed) GoAgain() bool           { return true }
func (MinnowismRed) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 3, minnowismIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type MinnowismYellow struct{}

func (MinnowismYellow) ID() ids.CardID          { return ids.MinnowismYellow }
func (MinnowismYellow) Name() string            { return "Minnowism" }
func (MinnowismYellow) Cost(*sim.TurnState) int { return 0 }
func (MinnowismYellow) Pitch() int              { return 2 }
func (MinnowismYellow) Attack() int             { return 0 }
func (MinnowismYellow) Defense() int            { return 2 }
func (MinnowismYellow) Types() card.TypeSet     { return minnowismTypes }
func (MinnowismYellow) GoAgain() bool           { return true }
func (MinnowismYellow) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 2, minnowismIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type MinnowismBlue struct{}

func (MinnowismBlue) ID() ids.CardID          { return ids.MinnowismBlue }
func (MinnowismBlue) Name() string            { return "Minnowism" }
func (MinnowismBlue) Cost(*sim.TurnState) int { return 0 }
func (MinnowismBlue) Pitch() int              { return 3 }
func (MinnowismBlue) Attack() int             { return 0 }
func (MinnowismBlue) Defense() int            { return 2 }
func (MinnowismBlue) Types() card.TypeSet     { return minnowismTypes }
func (MinnowismBlue) GoAgain() bool           { return true }
func (MinnowismBlue) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 1, minnowismIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
