// On the Horizon — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Printed
// defense: Red 4, Yellow 3, Blue 2.
//
// Text: "When this defends, look at the top card of your deck."
//
// Block-typed: only legal roles are pitch and plain block, so Play is never invoked by
// the chain runner (the partition enumerator forbids Attack via the Action / Weapon
// gate). The deck-peek defend trigger isn't modelled — it surfaces information for the
// player, not a state change the solver can credit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var onTheHorizonTypes = card.NewTypeSet(card.TypeGeneric, card.TypeBlock)

type OnTheHorizonRed struct{}

func (OnTheHorizonRed) ID() ids.CardID                      { return ids.OnTheHorizonRed }
func (OnTheHorizonRed) Name() string                        { return "On the Horizon" }
func (OnTheHorizonRed) Cost(*sim.TurnState) int             { return 0 }
func (OnTheHorizonRed) Pitch() int                          { return 1 }
func (OnTheHorizonRed) Attack() int                         { return 0 }
func (OnTheHorizonRed) Defense() int                        { return 4 }
func (OnTheHorizonRed) Types() card.TypeSet                 { return onTheHorizonTypes }
func (OnTheHorizonRed) GoAgain() bool                       { return false }
func (OnTheHorizonRed) Play(*sim.TurnState, *sim.CardState) {}

type OnTheHorizonYellow struct{}

func (OnTheHorizonYellow) ID() ids.CardID                      { return ids.OnTheHorizonYellow }
func (OnTheHorizonYellow) Name() string                        { return "On the Horizon" }
func (OnTheHorizonYellow) Cost(*sim.TurnState) int             { return 0 }
func (OnTheHorizonYellow) Pitch() int                          { return 2 }
func (OnTheHorizonYellow) Attack() int                         { return 0 }
func (OnTheHorizonYellow) Defense() int                        { return 3 }
func (OnTheHorizonYellow) Types() card.TypeSet                 { return onTheHorizonTypes }
func (OnTheHorizonYellow) GoAgain() bool                       { return false }
func (OnTheHorizonYellow) Play(*sim.TurnState, *sim.CardState) {}

type OnTheHorizonBlue struct{}

func (OnTheHorizonBlue) ID() ids.CardID                      { return ids.OnTheHorizonBlue }
func (OnTheHorizonBlue) Name() string                        { return "On the Horizon" }
func (OnTheHorizonBlue) Cost(*sim.TurnState) int             { return 0 }
func (OnTheHorizonBlue) Pitch() int                          { return 3 }
func (OnTheHorizonBlue) Attack() int                         { return 0 }
func (OnTheHorizonBlue) Defense() int                        { return 2 }
func (OnTheHorizonBlue) Types() card.TypeSet                 { return onTheHorizonTypes }
func (OnTheHorizonBlue) GoAgain() bool                       { return false }
func (OnTheHorizonBlue) Play(*sim.TurnState, *sim.CardState) {}
