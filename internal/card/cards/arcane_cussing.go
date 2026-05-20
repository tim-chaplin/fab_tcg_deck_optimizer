// Arcane Cussing — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When you deal or are dealt damage, destroy this. When this leaves the arena
// during your turn, create N Runechants." (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (ArcaneCussingRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateHitOrDamageTakenAura(self, arcaneCussingAuraHandler, 3, nil)
}

func (ArcaneCussingYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateHitOrDamageTakenAura(self, arcaneCussingAuraHandler, 2, nil)
}

func (ArcaneCussingBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateHitOrDamageTakenAura(self, arcaneCussingAuraHandler, 1, nil)
}

// arcaneCussingAuraHandler runs when Arcane Cussing leaves the arena — destroyed by you
// dealing damage (a Hit) or being dealt damage (DamageTaken). The unfiltered Hit lets any
// attacker pop it. It creates a.Count() Runechants only when the aura leaves during your
// turn, then destroys the aura.
func arcaneCussingAuraHandler(ge card.GameEngine, l card.Logger, a card.Aura) {
	if ge.IsMyTurn() {
		n := a.Count()
		ge.CreateRunechants(n)
		l.AppendPostTriggerf(a.CardName(), n, "Created %d runechants", n)
	}
	a.Destroy(true)
}
