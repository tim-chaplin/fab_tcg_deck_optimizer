// Sloggism — Generic Action. Cost 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "The next attack action card with cost 2 or greater you play this turn gains +N{p}. **Go
// again**" (Red N=6, Yellow N=5, Blue N=4.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sloggismTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// sloggismIsTarget gates the rider on attack action cards whose cost is 2 or more.
func sloggismIsTarget(s *sim.TurnState, pc *sim.CardState) bool {
	return pc.Card.Types().IsAttackAction() && pc.Card.Cost(s) >= 2
}

type SloggismRed struct{}

func (SloggismRed) ID() ids.CardID          { return ids.SloggismRed }
func (SloggismRed) Name() string            { return "Sloggism" }
func (SloggismRed) Cost(*sim.TurnState) int { return 3 }
func (SloggismRed) Pitch() int              { return 1 }
func (SloggismRed) Attack() int             { return 0 }
func (SloggismRed) Defense() int            { return 2 }
func (SloggismRed) Types() card.TypeSet     { return sloggismTypes }
func (SloggismRed) GoAgain() bool           { return true }
func (SloggismRed) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 6, sloggismIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type SloggismYellow struct{}

func (SloggismYellow) ID() ids.CardID          { return ids.SloggismYellow }
func (SloggismYellow) Name() string            { return "Sloggism" }
func (SloggismYellow) Cost(*sim.TurnState) int { return 3 }
func (SloggismYellow) Pitch() int              { return 2 }
func (SloggismYellow) Attack() int             { return 0 }
func (SloggismYellow) Defense() int            { return 2 }
func (SloggismYellow) Types() card.TypeSet     { return sloggismTypes }
func (SloggismYellow) GoAgain() bool           { return true }
func (SloggismYellow) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 5, sloggismIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type SloggismBlue struct{}

func (SloggismBlue) ID() ids.CardID          { return ids.SloggismBlue }
func (SloggismBlue) Name() string            { return "Sloggism" }
func (SloggismBlue) Cost(*sim.TurnState) int { return 3 }
func (SloggismBlue) Pitch() int              { return 3 }
func (SloggismBlue) Attack() int             { return 0 }
func (SloggismBlue) Defense() int            { return 2 }
func (SloggismBlue) Types() card.TypeSet     { return sloggismTypes }
func (SloggismBlue) GoAgain() bool           { return true }
func (SloggismBlue) Play(s *sim.TurnState, self *sim.CardState) {
	GrantNextCardBonusAttack(s, 4, sloggismIsTarget)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
