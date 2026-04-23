// Bloodspill Invocation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When an attack action card you control hits, destroy Bloodspill Invocation
// then create N Runechant tokens. When your hero is dealt damage, destroy Bloodspill
// Invocation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelled as a fragile aura (fragile_aura.go). Only attack action cards qualify for the
// same-turn pop (weapons don't trigger Bloodspill), so Play passes attackActionOnly=true.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var bloodspillInvocationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type BloodspillInvocationRed struct{}

func (BloodspillInvocationRed) ID() card.ID                { return card.BloodspillInvocationRed }
func (BloodspillInvocationRed) Name() string               { return "Bloodspill Invocation (Red)" }
func (BloodspillInvocationRed) Cost(*card.TurnState) int                  { return 1 }
func (BloodspillInvocationRed) Pitch() int                 { return 1 }
func (BloodspillInvocationRed) Attack() int                { return 0 }
func (BloodspillInvocationRed) Defense() int               { return 2 }
func (BloodspillInvocationRed) Types() card.TypeSet        { return bloodspillInvocationTypes }
func (BloodspillInvocationRed) GoAgain() bool              { return true }
func (BloodspillInvocationRed) Play(s *card.TurnState, _ *card.CardState) int { return fragileAuraValue(s, 3, true) }

type BloodspillInvocationYellow struct{}

func (BloodspillInvocationYellow) ID() card.ID                { return card.BloodspillInvocationYellow }
func (BloodspillInvocationYellow) Name() string               { return "Bloodspill Invocation (Yellow)" }
func (BloodspillInvocationYellow) Cost(*card.TurnState) int                  { return 1 }
func (BloodspillInvocationYellow) Pitch() int                 { return 2 }
func (BloodspillInvocationYellow) Attack() int                { return 0 }
func (BloodspillInvocationYellow) Defense() int               { return 2 }
func (BloodspillInvocationYellow) Types() card.TypeSet        { return bloodspillInvocationTypes }
func (BloodspillInvocationYellow) GoAgain() bool              { return true }
func (BloodspillInvocationYellow) Play(s *card.TurnState, _ *card.CardState) int { return fragileAuraValue(s, 2, true) }

type BloodspillInvocationBlue struct{}

func (BloodspillInvocationBlue) ID() card.ID                { return card.BloodspillInvocationBlue }
func (BloodspillInvocationBlue) Name() string               { return "Bloodspill Invocation (Blue)" }
func (BloodspillInvocationBlue) Cost(*card.TurnState) int                  { return 1 }
func (BloodspillInvocationBlue) Pitch() int                 { return 3 }
func (BloodspillInvocationBlue) Attack() int                { return 0 }
func (BloodspillInvocationBlue) Defense() int               { return 2 }
func (BloodspillInvocationBlue) Types() card.TypeSet        { return bloodspillInvocationTypes }
func (BloodspillInvocationBlue) GoAgain() bool              { return true }
func (BloodspillInvocationBlue) Play(s *card.TurnState, _ *card.CardState) int { return fragileAuraValue(s, 1, true) }
