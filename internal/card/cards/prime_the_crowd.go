// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2. Go again.
//
// Text: "The next attack action card you play this turn gets +N{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**" (Red N=4, Yellow N=3, Blue
// N=2.)
//
// Crowd cheers / Crowd boos are hero-state riders the single-turn solver doesn't model.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (PrimeTheCrowdRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 4, card.IsAttackAction)
}

func (PrimeTheCrowdYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 3, card.IsAttackAction)
}

func (PrimeTheCrowdBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	GrantNextCardBonusAttack(ge, 2, card.IsAttackAction)
}
