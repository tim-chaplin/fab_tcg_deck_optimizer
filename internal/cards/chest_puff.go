// Chest Puff — Generic Action - Attack. Cost 2, Pitch 1, Power 7, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."

package cards

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

// not implemented: pay {r} or lose 1{p} resolved as 'always pay'
func (ChestPuffRed) NotImplemented() {}
func (c ChestPuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
