// Chest Puff — Generic Action - Attack. Cost 2, Pitch 1, Power 7, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var chestPuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type ChestPuffRed struct{}

func (ChestPuffRed) ID() ids.CardID          { return ids.ChestPuffRed }
func (ChestPuffRed) Name() string            { return "Chest Puff" }
func (ChestPuffRed) Cost(*sim.TurnState) int { return 2 }
func (ChestPuffRed) Pitch() int              { return 1 }
func (ChestPuffRed) Attack() int             { return 7 }
func (ChestPuffRed) Defense() int            { return 3 }
func (ChestPuffRed) Types() card.TypeSet     { return chestPuffTypes }
func (ChestPuffRed) GoAgain() bool           { return false }
func (ChestPuffRed) Unplayable()             {}
func (c ChestPuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
