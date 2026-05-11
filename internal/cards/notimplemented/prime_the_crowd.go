// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gets +N{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**" (Red N=4, Yellow N=3, Blue
// N=2.)

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
)

// not implemented: Crowd cheers / Crowd boos keywords dropped

func (PrimeTheCrowdRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	cards.GrantNextCardBonusAttack(s, 4, cards.IsAttackAction)
}

// not implemented: Crowd cheers / Crowd boos keywords dropped

func (PrimeTheCrowdYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	cards.GrantNextCardBonusAttack(s, 3, cards.IsAttackAction)
}

// not implemented: Crowd cheers / Crowd boos keywords dropped

func (PrimeTheCrowdBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	cards.GrantNextCardBonusAttack(s, 2, cards.IsAttackAction)
}
