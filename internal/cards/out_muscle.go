// Out Muscle — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Out Muscle isn't defended by a card with equal or greater {p}, it has **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var outMuscleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type OutMuscleRed struct{}

func (OutMuscleRed) ID() ids.CardID          { return ids.OutMuscleRed }
func (OutMuscleRed) Name() string            { return "Out Muscle" }
func (OutMuscleRed) Cost(*sim.TurnState) int { return 3 }
func (OutMuscleRed) Pitch() int              { return 1 }
func (OutMuscleRed) Attack() int             { return 6 }
func (OutMuscleRed) Defense() int            { return 2 }
func (OutMuscleRed) Types() card.TypeSet     { return outMuscleTypes }
func (OutMuscleRed) GoAgain() bool           { return false }

// not implemented: defended-by-equal-or-greater-power go-again gate
func (OutMuscleRed) NotImplemented() {}
func (c OutMuscleRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type OutMuscleYellow struct{}

func (OutMuscleYellow) ID() ids.CardID          { return ids.OutMuscleYellow }
func (OutMuscleYellow) Name() string            { return "Out Muscle" }
func (OutMuscleYellow) Cost(*sim.TurnState) int { return 3 }
func (OutMuscleYellow) Pitch() int              { return 2 }
func (OutMuscleYellow) Attack() int             { return 5 }
func (OutMuscleYellow) Defense() int            { return 2 }
func (OutMuscleYellow) Types() card.TypeSet     { return outMuscleTypes }
func (OutMuscleYellow) GoAgain() bool           { return false }

// not implemented: defended-by-equal-or-greater-power go-again gate
func (OutMuscleYellow) NotImplemented() {}
func (c OutMuscleYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type OutMuscleBlue struct{}

func (OutMuscleBlue) ID() ids.CardID          { return ids.OutMuscleBlue }
func (OutMuscleBlue) Name() string            { return "Out Muscle" }
func (OutMuscleBlue) Cost(*sim.TurnState) int { return 3 }
func (OutMuscleBlue) Pitch() int              { return 3 }
func (OutMuscleBlue) Attack() int             { return 4 }
func (OutMuscleBlue) Defense() int            { return 2 }
func (OutMuscleBlue) Types() card.TypeSet     { return outMuscleTypes }
func (OutMuscleBlue) GoAgain() bool           { return false }

// not implemented: defended-by-equal-or-greater-power go-again gate
func (OutMuscleBlue) NotImplemented() {}
func (c OutMuscleBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
