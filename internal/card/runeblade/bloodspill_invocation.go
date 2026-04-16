// Bloodspill Invocation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When an attack action card you control hits, destroy Bloodspill Invocation then
// create N Runechant tokens. When your hero is dealt damage, destroy Bloodspill Invocation."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Simplification: assume an attack will land before we take damage, so the aura pays out fully.
// Play returns N (Red=3, Yellow=2, Blue=1).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var bloodspillInvocationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type BloodspillInvocationRed struct{}

func (BloodspillInvocationRed) Name() string              { return "Bloodspill Invocation (Red)" }
func (BloodspillInvocationRed) Cost() int                 { return 1 }
func (BloodspillInvocationRed) Pitch() int                { return 1 }
func (BloodspillInvocationRed) Attack() int               { return 0 }
func (BloodspillInvocationRed) Defense() int              { return 2 }
func (BloodspillInvocationRed) Types() card.TypeSet    { return bloodspillInvocationTypes }
func (BloodspillInvocationRed) GoAgain() bool             { return true }
func (BloodspillInvocationRed) Play(*card.TurnState) int  { return 3 }

type BloodspillInvocationYellow struct{}

func (BloodspillInvocationYellow) Name() string             { return "Bloodspill Invocation (Yellow)" }
func (BloodspillInvocationYellow) Cost() int                { return 1 }
func (BloodspillInvocationYellow) Pitch() int               { return 2 }
func (BloodspillInvocationYellow) Attack() int              { return 0 }
func (BloodspillInvocationYellow) Defense() int             { return 2 }
func (BloodspillInvocationYellow) Types() card.TypeSet   { return bloodspillInvocationTypes }
func (BloodspillInvocationYellow) GoAgain() bool            { return true }
func (BloodspillInvocationYellow) Play(*card.TurnState) int { return 2 }

type BloodspillInvocationBlue struct{}

func (BloodspillInvocationBlue) Name() string             { return "Bloodspill Invocation (Blue)" }
func (BloodspillInvocationBlue) Cost() int                { return 1 }
func (BloodspillInvocationBlue) Pitch() int               { return 3 }
func (BloodspillInvocationBlue) Attack() int              { return 0 }
func (BloodspillInvocationBlue) Defense() int             { return 2 }
func (BloodspillInvocationBlue) Types() card.TypeSet   { return bloodspillInvocationTypes }
func (BloodspillInvocationBlue) GoAgain() bool            { return true }
func (BloodspillInvocationBlue) Play(*card.TurnState) int { return 1 }
