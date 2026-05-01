// Flex — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you attack or defend with Flex, you may pay {r}{r}. If you do, it gains +2{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var flexTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FlexRed struct{}

func (FlexRed) ID() ids.CardID          { return ids.FlexRed }
func (FlexRed) Name() string            { return "Flex" }
func (FlexRed) Cost(*sim.TurnState) int { return 0 }
func (FlexRed) Pitch() int              { return 1 }
func (FlexRed) Attack() int             { return 4 }
func (FlexRed) Defense() int            { return 2 }
func (FlexRed) Types() card.TypeSet     { return flexTypes }
func (FlexRed) GoAgain() bool           { return false }

// not implemented: pay-{r}{r}-for-+2{p} attack/defence mode
func (FlexRed) NotImplemented() {}
func (c FlexRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type FlexYellow struct{}

func (FlexYellow) ID() ids.CardID          { return ids.FlexYellow }
func (FlexYellow) Name() string            { return "Flex" }
func (FlexYellow) Cost(*sim.TurnState) int { return 0 }
func (FlexYellow) Pitch() int              { return 2 }
func (FlexYellow) Attack() int             { return 3 }
func (FlexYellow) Defense() int            { return 2 }
func (FlexYellow) Types() card.TypeSet     { return flexTypes }
func (FlexYellow) GoAgain() bool           { return false }

// not implemented: pay-{r}{r}-for-+2{p} attack/defence mode
func (FlexYellow) NotImplemented() {}
func (c FlexYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}

type FlexBlue struct{}

func (FlexBlue) ID() ids.CardID          { return ids.FlexBlue }
func (FlexBlue) Name() string            { return "Flex" }
func (FlexBlue) Cost(*sim.TurnState) int { return 0 }
func (FlexBlue) Pitch() int              { return 3 }
func (FlexBlue) Attack() int             { return 2 }
func (FlexBlue) Defense() int            { return 2 }
func (FlexBlue) Types() card.TypeSet     { return flexTypes }
func (FlexBlue) GoAgain() bool           { return false }

// not implemented: pay-{r}{r}-for-+2{p} attack/defence mode
func (FlexBlue) NotImplemented() {}
func (c FlexBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
}
