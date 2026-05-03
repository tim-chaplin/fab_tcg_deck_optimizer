// Look Tuff — Generic Action - Attack. Cost 3, Pitch 1, Power 8, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Mode 0 pays the printed 3{r} for 7{p}; mode 1 spends an extra {r} for the full 8{p}.
// See Bluster Buff for the modal-cost wiring.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var lookTuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func lookTuffPlay(s *sim.TurnState, self *sim.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type LookTuffRed struct{}

func (LookTuffRed) ID() ids.CardID          { return ids.LookTuffRed }
func (LookTuffRed) Name() string            { return "Look Tuff" }
func (LookTuffRed) Cost(*sim.TurnState) int { return 3 }
func (LookTuffRed) Pitch() int              { return 1 }
func (LookTuffRed) Attack() int             { return 8 }
func (LookTuffRed) Defense() int            { return 3 }
func (LookTuffRed) Types() card.TypeSet     { return lookTuffTypes }
func (LookTuffRed) GoAgain() bool           { return false }
func (LookTuffRed) Modes() int              { return 2 }
func (LookTuffRed) ModalCost(mode int8) int { return 3 + int(mode) }
func (LookTuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	lookTuffPlay(s, self)
}
