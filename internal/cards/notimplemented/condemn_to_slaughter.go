// Condemn to Slaughter — Runeblade Action. Cost 1, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gets +N{p}. You may destroy an aura you control. If
// you do, each opponent destroys an aura permanent they control. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// not implemented: aura-trade rider and opponent-aura destruction clause; only same-turn
// Runeblade-attack +N{p} is modelled

func (CondemnToSlaughterRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	condemnToSlaughterApplySideEffect(ge, 3)
}

// not implemented: aura-trade rider and opponent-aura destruction clause; only same-turn
// Runeblade-attack +N{p} is modelled

func (CondemnToSlaughterYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	condemnToSlaughterApplySideEffect(ge, 2)
}

// not implemented: aura-trade rider and opponent-aura destruction clause; only same-turn
// Runeblade-attack +N{p} is modelled

func (CondemnToSlaughterBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	condemnToSlaughterApplySideEffect(ge, 1)
}

// condemnToSlaughterApplySideEffect grants +n to the first scheduled Runeblade attack (attack
// action card or weapon swing) via pc.BonusAttack so the buffed attack's EffectiveAttack folds
// the bonus into LikelyToHit and the chain credit lands on the target's slot. Condemn's own
// contribution is zero.
func condemnToSlaughterApplySideEffect(ge card.GameEngine, n int) {
	for _, pc := range ge.CardsRemaining() {
		if pc.Card.Types(ge).IsRunebladeAttack() {
			pc.BonusAttack += n
			return
		}
	}
}
