// Raging Onslaught — Generic Action - Attack. Cost 3. Printed power: Red 7, Yellow 6, Blue 5.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var ragingOnslaughtTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type RagingOnslaughtRed struct{}

func (RagingOnslaughtRed) ID() ids.CardID          { return ids.RagingOnslaughtRed }
func (RagingOnslaughtRed) Name() string            { return "Raging Onslaught" }
func (RagingOnslaughtRed) Cost(*sim.TurnState) int { return 3 }
func (RagingOnslaughtRed) Pitch() int              { return 1 }
func (RagingOnslaughtRed) Attack() int             { return 7 }
func (RagingOnslaughtRed) Defense() int            { return 3 }
func (RagingOnslaughtRed) Types() card.TypeSet     { return ragingOnslaughtTypes }
func (RagingOnslaughtRed) GoAgain() bool           { return false }
func (c RagingOnslaughtRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type RagingOnslaughtYellow struct{}

func (RagingOnslaughtYellow) ID() ids.CardID          { return ids.RagingOnslaughtYellow }
func (RagingOnslaughtYellow) Name() string            { return "Raging Onslaught" }
func (RagingOnslaughtYellow) Cost(*sim.TurnState) int { return 3 }
func (RagingOnslaughtYellow) Pitch() int              { return 2 }
func (RagingOnslaughtYellow) Attack() int             { return 6 }
func (RagingOnslaughtYellow) Defense() int            { return 3 }
func (RagingOnslaughtYellow) Types() card.TypeSet     { return ragingOnslaughtTypes }
func (RagingOnslaughtYellow) GoAgain() bool           { return false }
func (c RagingOnslaughtYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type RagingOnslaughtBlue struct{}

func (RagingOnslaughtBlue) ID() ids.CardID          { return ids.RagingOnslaughtBlue }
func (RagingOnslaughtBlue) Name() string            { return "Raging Onslaught" }
func (RagingOnslaughtBlue) Cost(*sim.TurnState) int { return 3 }
func (RagingOnslaughtBlue) Pitch() int              { return 3 }
func (RagingOnslaughtBlue) Attack() int             { return 5 }
func (RagingOnslaughtBlue) Defense() int            { return 3 }
func (RagingOnslaughtBlue) Types() card.TypeSet     { return ragingOnslaughtTypes }
func (RagingOnslaughtBlue) GoAgain() bool           { return false }
func (c RagingOnslaughtBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
