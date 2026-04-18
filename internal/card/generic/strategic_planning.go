// Strategic Planning — Generic Action. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Put an action card with cost 2 or less from a graveyard on the bottom of its owner's deck.
// At the beginning of the end phase, draw a card. **Go again**"
//
// Simplification: Graveyard recovery and end-phase draw aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var strategicPlanningTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type StrategicPlanningRed struct{}

func (StrategicPlanningRed) ID() card.ID                 { return card.StrategicPlanningRed }
func (StrategicPlanningRed) Name() string                { return "Strategic Planning (Red)" }
func (StrategicPlanningRed) Cost() int                   { return 1 }
func (StrategicPlanningRed) Pitch() int                  { return 1 }
func (StrategicPlanningRed) Attack() int                 { return 0 }
func (StrategicPlanningRed) Defense() int                { return 2 }
func (StrategicPlanningRed) Types() card.TypeSet         { return strategicPlanningTypes }
func (StrategicPlanningRed) GoAgain() bool               { return true }
func (StrategicPlanningRed) Play(s *card.TurnState) int { return 0 }

type StrategicPlanningYellow struct{}

func (StrategicPlanningYellow) ID() card.ID                 { return card.StrategicPlanningYellow }
func (StrategicPlanningYellow) Name() string                { return "Strategic Planning (Yellow)" }
func (StrategicPlanningYellow) Cost() int                   { return 1 }
func (StrategicPlanningYellow) Pitch() int                  { return 2 }
func (StrategicPlanningYellow) Attack() int                 { return 0 }
func (StrategicPlanningYellow) Defense() int                { return 2 }
func (StrategicPlanningYellow) Types() card.TypeSet         { return strategicPlanningTypes }
func (StrategicPlanningYellow) GoAgain() bool               { return true }
func (StrategicPlanningYellow) Play(s *card.TurnState) int { return 0 }

type StrategicPlanningBlue struct{}

func (StrategicPlanningBlue) ID() card.ID                 { return card.StrategicPlanningBlue }
func (StrategicPlanningBlue) Name() string                { return "Strategic Planning (Blue)" }
func (StrategicPlanningBlue) Cost() int                   { return 1 }
func (StrategicPlanningBlue) Pitch() int                  { return 3 }
func (StrategicPlanningBlue) Attack() int                 { return 0 }
func (StrategicPlanningBlue) Defense() int                { return 2 }
func (StrategicPlanningBlue) Types() card.TypeSet         { return strategicPlanningTypes }
func (StrategicPlanningBlue) GoAgain() bool               { return true }
func (StrategicPlanningBlue) Play(s *card.TurnState) int { return 0 }
