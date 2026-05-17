// Prime the Crowd — Generic Action. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next attack action card you play this turn gets +N{p}. **The crowd cheers** each
// Revered hero. **The crowd boos** each Reviled hero. **Go again**" (Red N=4, Yellow N=3, Blue
// N=2.)

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card/cards"
)

// not implemented: Crowd cheers / Crowd boos keywords dropped

func (PrimeTheCrowdRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cards.GrantNextCardBonusAttack(ge, 4, card.IsAttackAction)
}

// not implemented: Crowd cheers / Crowd boos keywords dropped

func (PrimeTheCrowdYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cards.GrantNextCardBonusAttack(ge, 3, card.IsAttackAction)
}

// not implemented: Crowd cheers / Crowd boos keywords dropped

func (PrimeTheCrowdBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	cards.GrantNextCardBonusAttack(ge, 2, card.IsAttackAction)
}
