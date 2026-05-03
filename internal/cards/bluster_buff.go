// Bluster Buff — Generic Action - Attack. Cost 1, Pitch 1, Power 6, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."
//
// Two modes via sim.ModalCard + sim.ModalCost: mode 0 pays the printed 1{r} for 5{p};
// mode 1 spends an extra {r} for the full 6{p}. The chain runner enumerates both and
// picks the higher-Value tuple per partition.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var blusterBuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func blusterBuffPlay(s *sim.TurnState, self *sim.CardState) {
	if self.Mode == 0 {
		self.BonusAttack -= 1
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type BlusterBuffRed struct{}

func (BlusterBuffRed) ID() ids.CardID          { return ids.BlusterBuffRed }
func (BlusterBuffRed) Name() string            { return "Bluster Buff" }
func (BlusterBuffRed) Cost(*sim.TurnState) int { return 1 }
func (BlusterBuffRed) Pitch() int              { return 1 }
func (BlusterBuffRed) Attack() int             { return 6 }
func (BlusterBuffRed) Defense() int            { return 3 }
func (BlusterBuffRed) Types() card.TypeSet     { return blusterBuffTypes }
func (BlusterBuffRed) GoAgain() bool           { return false }
func (BlusterBuffRed) Modes() int              { return 2 }
func (BlusterBuffRed) ModalCost(mode int8) int { return 1 + int(mode) }
func (BlusterBuffRed) Play(s *sim.TurnState, self *sim.CardState) {
	blusterBuffPlay(s, self)
}
