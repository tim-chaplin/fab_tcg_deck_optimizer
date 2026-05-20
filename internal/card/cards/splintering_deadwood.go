// Splintering Deadwood — Runeblade Action - Attack. Cost 3, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Printed power: Red 7, Yellow 6, Blue 5.
// Text: "When this attacks or hits, you may destroy an aura you control. If you do, create a
// Runechant token."
//
// The "you may destroy an aura you control" leg is optional, so the trade only happens when
// it is worth a Runechant; see GameEngine.DestroyOwnedAura for which auras qualify.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// splinteringDeadwoodPlay wires both halves of "when this attacks or hits": the attack leg
// runs immediately, the hit leg registers for when this attack lands.
func splinteringDeadwoodPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	splinteringDeadwoodTrade(ge, l, self)
	self.RegisterOnHit(splinteringDeadwoodOnHit)
}

// splinteringDeadwoodTrade runs one leg of the rider: destroy an owned aura, and when one
// is destroyed, create the Runechant.
func splinteringDeadwoodTrade(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if !ge.DestroyOwnedAura() {
		return
	}
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Destroyed an aura to create a runechant", 1)
}

func splinteringDeadwoodOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	splinteringDeadwoodTrade(ge, l, self)
}

func (SplinteringDeadwoodRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	splinteringDeadwoodPlay(ge, l, self)
}

func (SplinteringDeadwoodYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	splinteringDeadwoodPlay(ge, l, self)
}

func (SplinteringDeadwoodBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	splinteringDeadwoodPlay(ge, l, self)
}
