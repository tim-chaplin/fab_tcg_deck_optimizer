// Chest Puff — Generic Action - Attack. Cost 2, Pitch 1, Power 7, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Mode 0 pays the printed 2{r} for 6{p}; mode 1 spends an extra {r} for the full 7{p}.
// See Bluster Buff for the modal-cost wiring.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var chestPuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func chestPuffPlay(s *sim.TurnState, self *sim.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type ChestPuffRed struct{}

func (ChestPuffRed) ID() ids.CardID          { return ids.ChestPuffRed }
func (ChestPuffRed) Name() string            { return "Chest Puff" }
func (ChestPuffRed) Cost(*sim.TurnState) int { return 2 }
func (ChestPuffRed) Pitch() int              { return 1 }
func (ChestPuffRed) Attack() int             { return 7 }
func (ChestPuffRed) Defense() int            { return 3 }
func (ChestPuffRed) Types() card.TypeSet     { return chestPuffTypes }
func (ChestPuffRed) GoAgain() bool           { return false }
func (ChestPuffRed) Modes() int              { return 2 }
func (ChestPuffRed) ModalCost(mode int8) int { return 2 + int(mode) }
func (ChestPuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	chestPuffPlay(s, self)
}
