// Bloodspill Invocation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When an attack action card you control hits, destroy Bloodspill Invocation
// then create N Runechant tokens. When your hero is dealt damage, destroy Bloodspill
// Invocation." (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (BloodspillInvocationRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateHitOrDamageTakenAura(self, bloodspillInvocationAuraHandler, 3, card.TypeSet.IsAttackAction)
}

func (BloodspillInvocationYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateHitOrDamageTakenAura(self, bloodspillInvocationAuraHandler, 2, card.TypeSet.IsAttackAction)
}

func (BloodspillInvocationBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateHitOrDamageTakenAura(self, bloodspillInvocationAuraHandler, 1, card.TypeSet.IsAttackAction)
}

// bloodspillInvocationAuraHandler runs when Bloodspill Invocation leaves the arena. The
// IsAttackAction filter gates the Hit path to attack action cards (a weapon swing doesn't
// pop it); the DamageTaken path has no triggering card and so isn't gated. It creates
// a.Count() Runechants only on the hit path — a triggering card is present there but not
// on DamageTaken — then destroys the aura.
func bloodspillInvocationAuraHandler(ge card.GameEngine, l card.Logger, a card.Aura) {
	if ge.TriggeringCard() != nil {
		n := a.Count()
		ge.CreateRunechants(n)
		l.AppendPostTriggerf(a.CardName(), n, "Created %d runechants", n)
	}
	a.Destroy(true)
}
