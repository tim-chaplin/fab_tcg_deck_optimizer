// Condemn to Slaughter — Runeblade Action. Cost 1, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gets +N{p}. You may destroy an aura you
// control. If you do, each opponent destroys an aura permanent they control. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// The "each opponent destroys an aura permanent" follow-up is not modelled: there is no
// opponent board. The aura-sacrifice itself cashes a controlled payoff aura via
// SacrificePayoffAura.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// condemnToSlaughterPlay buffs the next Runeblade attack +n and sacrifices a controlled
// payoff aura.
func condemnToSlaughterPlay(ge card.GameEngine, n int) {
	GrantNextCardBonusAttack(ge, n, card.IsRunebladeAttack)
	ge.SacrificePayoffAura()
}

func (CondemnToSlaughterRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	condemnToSlaughterPlay(ge, 3)
}

func (CondemnToSlaughterYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	condemnToSlaughterPlay(ge, 2)
}

func (CondemnToSlaughterBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	condemnToSlaughterPlay(ge, 1)
}
