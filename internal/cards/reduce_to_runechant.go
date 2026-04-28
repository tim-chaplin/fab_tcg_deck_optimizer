// Reduce to Runechant — Runeblade Defense Reaction. Printed cost 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 4, Yellow 3, Blue 2.
// Text: "Reduce to Runechant costs {r} less to play for each Runechant you control. Create a
// Runechant token."
//
// Cost returns max(0, printed - s.Runechants) at play time (sim.VariableCost bounds [0, 1]).
// Play creates one Runechant, crediting +1 at creation. Defense-reaction state is reset
// between reactions so the token itself doesn't carry into next turn's carryover — only its
// damage credit lands.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var reduceToRunechantTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

const reduceToRunechantPrintedCost = 1

func reduceToRunechantCost(s *sim.TurnState) int {
	eff := reduceToRunechantPrintedCost - s.Runechants
	if eff < 0 {
		return 0
	}
	return eff
}

type ReduceToRunechantRed struct{}

func (ReduceToRunechantRed) ID() ids.CardID            { return ids.ReduceToRunechantRed }
func (ReduceToRunechantRed) Name() string              { return "Reduce to Runechant" }
func (ReduceToRunechantRed) Cost(s *sim.TurnState) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantRed) MinCost() int              { return 0 }
func (ReduceToRunechantRed) MaxCost() int              { return reduceToRunechantPrintedCost }
func (ReduceToRunechantRed) Pitch() int                { return 1 }
func (ReduceToRunechantRed) Attack() int               { return 0 }
func (ReduceToRunechantRed) Defense() int              { return 4 }
func (ReduceToRunechantRed) Types() card.TypeSet       { return reduceToRunechantTypes }
func (ReduceToRunechantRed) GoAgain() bool             { return false }
func (ReduceToRunechantRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}

type ReduceToRunechantYellow struct{}

func (ReduceToRunechantYellow) ID() ids.CardID            { return ids.ReduceToRunechantYellow }
func (ReduceToRunechantYellow) Name() string              { return "Reduce to Runechant" }
func (ReduceToRunechantYellow) Cost(s *sim.TurnState) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantYellow) MinCost() int              { return 0 }
func (ReduceToRunechantYellow) MaxCost() int              { return reduceToRunechantPrintedCost }
func (ReduceToRunechantYellow) Pitch() int                { return 2 }
func (ReduceToRunechantYellow) Attack() int               { return 0 }
func (ReduceToRunechantYellow) Defense() int              { return 3 }
func (ReduceToRunechantYellow) Types() card.TypeSet       { return reduceToRunechantTypes }
func (ReduceToRunechantYellow) GoAgain() bool             { return false }
func (ReduceToRunechantYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}

type ReduceToRunechantBlue struct{}

func (ReduceToRunechantBlue) ID() ids.CardID            { return ids.ReduceToRunechantBlue }
func (ReduceToRunechantBlue) Name() string              { return "Reduce to Runechant" }
func (ReduceToRunechantBlue) Cost(s *sim.TurnState) int { return reduceToRunechantCost(s) }
func (ReduceToRunechantBlue) MinCost() int              { return 0 }
func (ReduceToRunechantBlue) MaxCost() int              { return reduceToRunechantPrintedCost }
func (ReduceToRunechantBlue) Pitch() int                { return 3 }
func (ReduceToRunechantBlue) Attack() int               { return 0 }
func (ReduceToRunechantBlue) Defense() int              { return 2 }
func (ReduceToRunechantBlue) Types() card.TypeSet       { return reduceToRunechantTypes }
func (ReduceToRunechantBlue) GoAgain() bool             { return false }
func (ReduceToRunechantBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}
