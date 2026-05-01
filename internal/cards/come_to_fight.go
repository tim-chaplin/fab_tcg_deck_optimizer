// Come to Fight — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense
// 3.
//
// Text: "The next attack action card you play this turn gains +N{p}. **Go again**" (Red N=3,
// Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var comeToFightTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type ComeToFightRed struct{}

func (ComeToFightRed) ID() ids.CardID          { return ids.ComeToFightRed }
func (ComeToFightRed) Name() string            { return "Come to Fight" }
func (ComeToFightRed) Cost(*sim.TurnState) int { return 1 }
func (ComeToFightRed) Pitch() int              { return 1 }
func (ComeToFightRed) Attack() int             { return 0 }
func (ComeToFightRed) Defense() int            { return 3 }
func (ComeToFightRed) Types() card.TypeSet     { return comeToFightTypes }
func (ComeToFightRed) GoAgain() bool           { return true }
func (ComeToFightRed) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 3)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type ComeToFightYellow struct{}

func (ComeToFightYellow) ID() ids.CardID          { return ids.ComeToFightYellow }
func (ComeToFightYellow) Name() string            { return "Come to Fight" }
func (ComeToFightYellow) Cost(*sim.TurnState) int { return 1 }
func (ComeToFightYellow) Pitch() int              { return 2 }
func (ComeToFightYellow) Attack() int             { return 0 }
func (ComeToFightYellow) Defense() int            { return 3 }
func (ComeToFightYellow) Types() card.TypeSet     { return comeToFightTypes }
func (ComeToFightYellow) GoAgain() bool           { return true }
func (ComeToFightYellow) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 2)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type ComeToFightBlue struct{}

func (ComeToFightBlue) ID() ids.CardID          { return ids.ComeToFightBlue }
func (ComeToFightBlue) Name() string            { return "Come to Fight" }
func (ComeToFightBlue) Cost(*sim.TurnState) int { return 1 }
func (ComeToFightBlue) Pitch() int              { return 3 }
func (ComeToFightBlue) Attack() int             { return 0 }
func (ComeToFightBlue) Defense() int            { return 3 }
func (ComeToFightBlue) Types() card.TypeSet     { return comeToFightTypes }
func (ComeToFightBlue) GoAgain() bool           { return true }
func (ComeToFightBlue) Play(s *sim.TurnState, self *sim.CardState) {
	grantNextAttackActionBonus(s, 1)
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
