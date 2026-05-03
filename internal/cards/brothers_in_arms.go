// Brothers in Arms — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this defends, you may pay {r}. If you do, it gets +2{d}."
//
// Two block-time modes via sim.ModalCard + sim.BlockCost: mode 0 spends nothing for the
// printed 2{d}; mode 1 spends 1{r} for 4{d}. The chain runner enumerates both and picks
// the higher-defense mode that fits the partition's spare defense budget.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var brothersInArmsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func brothersInArmsPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// brothersInArmsBlock fires +2{d} on mode 1 (caller already deducted 1{r} from the spare
// defense budget). Mode 0 is the printed default with no modification.
func brothersInArmsBlock(_ *sim.TurnState, self *sim.CardState) {
	if self.Mode == 1 {
		self.BonusDefense += 2
	}
}

type BrothersInArmsRed struct{}

func (BrothersInArmsRed) ID() ids.CardID          { return ids.BrothersInArmsRed }
func (BrothersInArmsRed) Name() string            { return "Brothers in Arms" }
func (BrothersInArmsRed) Cost(*sim.TurnState) int { return 2 }
func (BrothersInArmsRed) Pitch() int              { return 1 }
func (BrothersInArmsRed) Attack() int             { return 6 }
func (BrothersInArmsRed) Defense() int            { return 2 }
func (BrothersInArmsRed) Types() card.TypeSet     { return brothersInArmsTypes }
func (BrothersInArmsRed) GoAgain() bool           { return false }
func (BrothersInArmsRed) Modes() int              { return 2 }
func (BrothersInArmsRed) BlockCost(mode int8) int { return int(mode) }
func (BrothersInArmsRed) Block(s *sim.TurnState, self *sim.CardState) {
	brothersInArmsBlock(s, self)
}
func (BrothersInArmsRed) Play(s *sim.TurnState, self *sim.CardState) {
	brothersInArmsPlay(s, self)
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
func (BrothersInArmsYellow) Modes() int              { return 2 }
func (BrothersInArmsYellow) BlockCost(mode int8) int { return int(mode) }
func (BrothersInArmsYellow) Block(s *sim.TurnState, self *sim.CardState) {
	brothersInArmsBlock(s, self)
}
func (BrothersInArmsYellow) Play(s *sim.TurnState, self *sim.CardState) {
	brothersInArmsPlay(s, self)
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
func (BrothersInArmsBlue) Modes() int              { return 2 }
func (BrothersInArmsBlue) BlockCost(mode int8) int { return int(mode) }
func (BrothersInArmsBlue) Block(s *sim.TurnState, self *sim.CardState) {
	brothersInArmsBlock(s, self)
}
func (BrothersInArmsBlue) Play(s *sim.TurnState, self *sim.CardState) {
	brothersInArmsPlay(s, self)
}
