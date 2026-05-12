// Oath of the Arknight — Runeblade Action. Cost 2, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gains +N{p}. Create a Runechant token. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (OathOfTheArknightRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	oathPlay(g, l, self, 3)
}

func (OathOfTheArknightYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	oathPlay(g, l, self, 2)
}

func (OathOfTheArknightBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	oathPlay(g, l, self, 1)
}

// oathPlay grants +n to the first scheduled Runeblade attack via pc.BonusAttack so the
// buffed attack's EffectiveAttack folds the bonus into LikelyToHit and the chain credit
// lands on the target's slot, not Oath'g. Always creates a Runechant token, which IS
// Oath's own contribution and lands as a sub-line under self's chain entry.
func oathPlay(g card.GameEngine, l card.Logger, self *card.CardState, bonus int) {
	GrantNextCardBonusAttack(g, bonus, IsRunebladeAttack)
	g.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
