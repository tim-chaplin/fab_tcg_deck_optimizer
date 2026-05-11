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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (BloodspillInvocationRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	fragileAuraPlay(s, l, self, 3, true)
}

func (BloodspillInvocationYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	fragileAuraPlay(s, l, self, 2, true)
}

func (BloodspillInvocationBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	fragileAuraPlay(s, l, self, 1, true)
}
