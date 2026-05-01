// Brush Off — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "The next time you would be dealt 3 or less damage this turn, prevent it."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var brushOffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type BrushOffRed struct{}

func (BrushOffRed) ID() ids.CardID          { return ids.BrushOffRed }
func (BrushOffRed) Name() string            { return "Brush Off" }
func (BrushOffRed) Cost(*sim.TurnState) int { return 0 }
func (BrushOffRed) Pitch() int              { return 1 }
func (BrushOffRed) Attack() int             { return 0 }
func (BrushOffRed) Defense() int            { return 0 }
func (BrushOffRed) Types() card.TypeSet     { return brushOffTypes }
func (BrushOffRed) GoAgain() bool           { return false }

// not implemented: Instant 'prevent next damage of N or less' (3 / 2 / 1 by colour)
func (BrushOffRed) NotImplemented()                            {}
func (BrushOffRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }

type BrushOffYellow struct{}

func (BrushOffYellow) ID() ids.CardID          { return ids.BrushOffYellow }
func (BrushOffYellow) Name() string            { return "Brush Off" }
func (BrushOffYellow) Cost(*sim.TurnState) int { return 0 }
func (BrushOffYellow) Pitch() int              { return 2 }
func (BrushOffYellow) Attack() int             { return 0 }
func (BrushOffYellow) Defense() int            { return 0 }
func (BrushOffYellow) Types() card.TypeSet     { return brushOffTypes }
func (BrushOffYellow) GoAgain() bool           { return false }

// not implemented: Instant 'prevent next damage of N or less' (3 / 2 / 1 by colour)
func (BrushOffYellow) NotImplemented()                            {}
func (BrushOffYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }

type BrushOffBlue struct{}

func (BrushOffBlue) ID() ids.CardID          { return ids.BrushOffBlue }
func (BrushOffBlue) Name() string            { return "Brush Off" }
func (BrushOffBlue) Cost(*sim.TurnState) int { return 0 }
func (BrushOffBlue) Pitch() int              { return 3 }
func (BrushOffBlue) Attack() int             { return 0 }
func (BrushOffBlue) Defense() int            { return 0 }
func (BrushOffBlue) Types() card.TypeSet     { return brushOffTypes }
func (BrushOffBlue) GoAgain() bool           { return false }

// not implemented: Instant 'prevent next damage of N or less' (3 / 2 / 1 by colour)
func (BrushOffBlue) NotImplemented()                            {}
func (BrushOffBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }
