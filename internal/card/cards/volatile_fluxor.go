// Volatile Fluxor — Lightning Action - Attack. Cost 0, Defense 3, Go again. Common.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "If you've played an instant card this chain link, this gets +N{p}. When this hits,
// create a Lightning Flow token. Go again." (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling:
//   - The "+N if you've played an instant" rider keys on HasPlayedType(Instant) at Play time —
//     the engine's record of cards played this turn. (The printed "this chain link" clause is
//     modelled at turn granularity.) Activating an item ability is not "playing a card", so a
//     cracked Amulet of Lightning doesn't satisfy it.
//   - The on-hit "create a Lightning Flow" is a RegisterOnHit rider; the token is inert fodder.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func fluxorPlay(ge card.GameEngine, self *card.CardState, n int) {
	if ge.HasPlayedType(card.TypeInstant) {
		self.BonusAttack += n
	}
	self.RegisterOnHit(fluxorOnHitCreateLightningFlow)
}

func fluxorOnHitCreateLightningFlow(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.CreateLightningFlow(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "created a Lightning Flow", 1)
}

func (VolatileFluxorRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fluxorPlay(ge, self, 3)
}

func (VolatileFluxorYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fluxorPlay(ge, self, 2)
}

func (VolatileFluxorBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	fluxorPlay(ge, self, 1)
}
