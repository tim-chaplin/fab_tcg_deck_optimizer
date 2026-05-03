// Drag Down — Generic Defense Reaction. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 0.
//
// Text: "When this defends an attack, it gets -3{p}."
//
// The -3{p} attacker debuff is modelled as a flat 3/2/1{d} block on R/Y/B respectively —
// equivalent to "block N damage" against an attack chain whose damage is consumed by
// IncomingDamage / BlockTotal. The Red→Blue decay matches the standard pitch-vs-body
// trade-off across pitch variants.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func dragDownPlay(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
}

type DragDownRed struct{}

func (DragDownRed) ID() ids.CardID          { return ids.DragDownRed }
func (DragDownRed) Name() string            { return "Drag Down" }
func (DragDownRed) Cost(*sim.TurnState) int { return 0 }
func (DragDownRed) Pitch() int              { return 1 }
func (DragDownRed) Attack() int             { return 0 }
func (DragDownRed) Defense() int            { return 3 }
func (DragDownRed) Types() card.TypeSet     { return DefenseReactionTypes }
func (DragDownRed) GoAgain() bool           { return false }
func (DragDownRed) Play(s *sim.TurnState, self *sim.CardState) {
	dragDownPlay(s, self)
}

type DragDownYellow struct{}

func (DragDownYellow) ID() ids.CardID          { return ids.DragDownYellow }
func (DragDownYellow) Name() string            { return "Drag Down" }
func (DragDownYellow) Cost(*sim.TurnState) int { return 0 }
func (DragDownYellow) Pitch() int              { return 2 }
func (DragDownYellow) Attack() int             { return 0 }
func (DragDownYellow) Defense() int            { return 2 }
func (DragDownYellow) Types() card.TypeSet     { return DefenseReactionTypes }
func (DragDownYellow) GoAgain() bool           { return false }
func (DragDownYellow) Play(s *sim.TurnState, self *sim.CardState) {
	dragDownPlay(s, self)
}

type DragDownBlue struct{}

func (DragDownBlue) ID() ids.CardID          { return ids.DragDownBlue }
func (DragDownBlue) Name() string            { return "Drag Down" }
func (DragDownBlue) Cost(*sim.TurnState) int { return 0 }
func (DragDownBlue) Pitch() int              { return 3 }
func (DragDownBlue) Attack() int             { return 0 }
func (DragDownBlue) Defense() int            { return 1 }
func (DragDownBlue) Types() card.TypeSet     { return DefenseReactionTypes }
func (DragDownBlue) GoAgain() bool           { return false }
func (DragDownBlue) Play(s *sim.TurnState, self *sim.CardState) {
	dragDownPlay(s, self)
}
