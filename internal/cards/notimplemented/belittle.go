// Belittle — Generic Action - Attack. Cost 1. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Belittle, you may reveal an attack action card with 3 or
// less base {p} from your hand. If you do, search your deck for a card named Minnowism, reveal it,
// put it into your hand, then shuffle your deck. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var belittleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BelittleRed struct{}

func (BelittleRed) ID() ids.CardID          { return ids.BelittleRed }
func (BelittleRed) Name() string            { return "Belittle" }
func (BelittleRed) Cost(*sim.TurnState) int { return 1 }
func (BelittleRed) Pitch() int              { return 1 }
func (BelittleRed) Attack() int             { return 3 }
func (BelittleRed) Defense() int            { return 2 }
func (BelittleRed) Types() card.TypeSet     { return belittleTypes }
func (BelittleRed) GoAgain() bool           { return true }
func (BelittleRed) NotSilverAgeLegal()      {}

// not implemented: Minnowism deck-search tutor (additional-cost reveal)
func (BelittleRed) NotImplemented() {}
func (c BelittleRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type BelittleYellow struct{}

func (BelittleYellow) ID() ids.CardID          { return ids.BelittleYellow }
func (BelittleYellow) Name() string            { return "Belittle" }
func (BelittleYellow) Cost(*sim.TurnState) int { return 1 }
func (BelittleYellow) Pitch() int              { return 2 }
func (BelittleYellow) Attack() int             { return 2 }
func (BelittleYellow) Defense() int            { return 2 }
func (BelittleYellow) Types() card.TypeSet     { return belittleTypes }
func (BelittleYellow) GoAgain() bool           { return true }
func (BelittleYellow) NotSilverAgeLegal()      {}

// not implemented: Minnowism deck-search tutor (additional-cost reveal)
func (BelittleYellow) NotImplemented() {}
func (c BelittleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type BelittleBlue struct{}

func (BelittleBlue) ID() ids.CardID          { return ids.BelittleBlue }
func (BelittleBlue) Name() string            { return "Belittle" }
func (BelittleBlue) Cost(*sim.TurnState) int { return 1 }
func (BelittleBlue) Pitch() int              { return 3 }
func (BelittleBlue) Attack() int             { return 1 }
func (BelittleBlue) Defense() int            { return 2 }
func (BelittleBlue) Types() card.TypeSet     { return belittleTypes }
func (BelittleBlue) GoAgain() bool           { return true }
func (BelittleBlue) NotSilverAgeLegal()      {}

// not implemented: Minnowism deck-search tutor (additional-cost reveal)
func (BelittleBlue) NotImplemented() {}
func (c BelittleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
