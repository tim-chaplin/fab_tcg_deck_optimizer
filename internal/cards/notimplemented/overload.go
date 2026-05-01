// Overload — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Dominate** If Overload hits, it gains **go again**."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var overloadTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OverloadRed struct{}

func (OverloadRed) ID() ids.CardID          { return ids.OverloadRed }
func (OverloadRed) Name() string            { return "Overload" }
func (OverloadRed) Cost(*sim.TurnState) int { return 0 }
func (OverloadRed) Pitch() int              { return 1 }
func (OverloadRed) Attack() int             { return 3 }
func (OverloadRed) Defense() int            { return 2 }
func (OverloadRed) Types() card.TypeSet     { return overloadTypes }
func (OverloadRed) GoAgain() bool           { return false }
func (OverloadRed) Dominate()               {}

// not implemented: on-hit go-again rider
func (OverloadRed) NotImplemented() {}
func (c OverloadRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type OverloadYellow struct{}

func (OverloadYellow) ID() ids.CardID          { return ids.OverloadYellow }
func (OverloadYellow) Name() string            { return "Overload" }
func (OverloadYellow) Cost(*sim.TurnState) int { return 0 }
func (OverloadYellow) Pitch() int              { return 2 }
func (OverloadYellow) Attack() int             { return 2 }
func (OverloadYellow) Defense() int            { return 2 }
func (OverloadYellow) Types() card.TypeSet     { return overloadTypes }
func (OverloadYellow) GoAgain() bool           { return false }
func (OverloadYellow) Dominate()               {}

// not implemented: on-hit go-again rider
func (OverloadYellow) NotImplemented() {}
func (c OverloadYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type OverloadBlue struct{}

func (OverloadBlue) ID() ids.CardID          { return ids.OverloadBlue }
func (OverloadBlue) Name() string            { return "Overload" }
func (OverloadBlue) Cost(*sim.TurnState) int { return 0 }
func (OverloadBlue) Pitch() int              { return 3 }
func (OverloadBlue) Attack() int             { return 1 }
func (OverloadBlue) Defense() int            { return 2 }
func (OverloadBlue) Types() card.TypeSet     { return overloadTypes }
func (OverloadBlue) GoAgain() bool           { return false }
func (OverloadBlue) Dominate()               {}

// not implemented: on-hit go-again rider
func (OverloadBlue) NotImplemented() {}
func (c OverloadBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
