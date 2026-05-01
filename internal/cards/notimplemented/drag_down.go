// Drag Down — Generic Defense Reaction. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 0.
//
// Text: "When this defends an attack, it gets -3{p}."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

type DragDownRed struct{}

func (DragDownRed) ID() ids.CardID          { return ids.DragDownRed }
func (DragDownRed) Name() string            { return "Drag Down" }
func (DragDownRed) Cost(*sim.TurnState) int { return 0 }
func (DragDownRed) Pitch() int              { return 1 }
func (DragDownRed) Attack() int             { return 0 }
func (DragDownRed) Defense() int            { return 0 }
func (DragDownRed) Types() card.TypeSet     { return cards.DefenseReactionTypes }
func (DragDownRed) GoAgain() bool           { return false }

// not implemented: -3{p} attacker debuff (defender-side power reduction not exposed)
func (DragDownRed) NotImplemented() {}
func (DragDownRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

type DragDownYellow struct{}

func (DragDownYellow) ID() ids.CardID          { return ids.DragDownYellow }
func (DragDownYellow) Name() string            { return "Drag Down" }
func (DragDownYellow) Cost(*sim.TurnState) int { return 0 }
func (DragDownYellow) Pitch() int              { return 2 }
func (DragDownYellow) Attack() int             { return 0 }
func (DragDownYellow) Defense() int            { return 0 }
func (DragDownYellow) Types() card.TypeSet     { return cards.DefenseReactionTypes }
func (DragDownYellow) GoAgain() bool           { return false }

// not implemented: -3{p} attacker debuff (defender-side power reduction not exposed)
func (DragDownYellow) NotImplemented() {}
func (DragDownYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

type DragDownBlue struct{}

func (DragDownBlue) ID() ids.CardID          { return ids.DragDownBlue }
func (DragDownBlue) Name() string            { return "Drag Down" }
func (DragDownBlue) Cost(*sim.TurnState) int { return 0 }
func (DragDownBlue) Pitch() int              { return 3 }
func (DragDownBlue) Attack() int             { return 0 }
func (DragDownBlue) Defense() int            { return 0 }
func (DragDownBlue) Types() card.TypeSet     { return cards.DefenseReactionTypes }
func (DragDownBlue) GoAgain() bool           { return false }

// not implemented: -3{p} attacker debuff (defender-side power reduction not exposed)
func (DragDownBlue) NotImplemented() {}
func (DragDownBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}
