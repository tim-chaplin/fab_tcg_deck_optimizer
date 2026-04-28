// Bloodspill Invocation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When an attack action card you control hits, destroy Bloodspill Invocation
// then create N Runechant tokens. When your hero is dealt damage, destroy Bloodspill
// Invocation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelled as a fragile aura (fragile_aura.go). Only attack action cards qualify for the
// same-turn pop (weapons don't trigger Bloodspill), so Play passes attackActionOnly=true.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var bloodspillInvocationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type BloodspillInvocationRed struct{}

func (BloodspillInvocationRed) ID() ids.CardID           { return ids.BloodspillInvocationRed }
func (BloodspillInvocationRed) Name() string             { return "Bloodspill Invocation" }
func (BloodspillInvocationRed) Cost(*card.TurnState) int { return 1 }
func (BloodspillInvocationRed) Pitch() int               { return 1 }
func (BloodspillInvocationRed) Attack() int              { return 0 }
func (BloodspillInvocationRed) Defense() int             { return 2 }
func (BloodspillInvocationRed) Types() card.TypeSet      { return bloodspillInvocationTypes }
func (BloodspillInvocationRed) GoAgain() bool            { return true }
func (BloodspillInvocationRed) Play(s *card.TurnState, self *card.CardState) {
	fragileAuraPlay(s, self, 3, true)
}

type BloodspillInvocationYellow struct{}

func (BloodspillInvocationYellow) ID() ids.CardID           { return ids.BloodspillInvocationYellow }
func (BloodspillInvocationYellow) Name() string             { return "Bloodspill Invocation" }
func (BloodspillInvocationYellow) Cost(*card.TurnState) int { return 1 }
func (BloodspillInvocationYellow) Pitch() int               { return 2 }
func (BloodspillInvocationYellow) Attack() int              { return 0 }
func (BloodspillInvocationYellow) Defense() int             { return 2 }
func (BloodspillInvocationYellow) Types() card.TypeSet      { return bloodspillInvocationTypes }
func (BloodspillInvocationYellow) GoAgain() bool            { return true }
func (BloodspillInvocationYellow) Play(s *card.TurnState, self *card.CardState) {
	fragileAuraPlay(s, self, 2, true)
}

type BloodspillInvocationBlue struct{}

func (BloodspillInvocationBlue) ID() ids.CardID           { return ids.BloodspillInvocationBlue }
func (BloodspillInvocationBlue) Name() string             { return "Bloodspill Invocation" }
func (BloodspillInvocationBlue) Cost(*card.TurnState) int { return 1 }
func (BloodspillInvocationBlue) Pitch() int               { return 3 }
func (BloodspillInvocationBlue) Attack() int              { return 0 }
func (BloodspillInvocationBlue) Defense() int             { return 2 }
func (BloodspillInvocationBlue) Types() card.TypeSet      { return bloodspillInvocationTypes }
func (BloodspillInvocationBlue) GoAgain() bool            { return true }
func (BloodspillInvocationBlue) Play(s *card.TurnState, self *card.CardState) {
	fragileAuraPlay(s, self, 1, true)
}
