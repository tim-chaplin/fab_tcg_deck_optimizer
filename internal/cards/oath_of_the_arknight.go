// Oath of the Arknight — Runeblade Action. Cost 2, Defense 3, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Your next Runeblade attack this turn gains +N{p}. Create a Runechant token. Go again."
// (Red N=3, Yellow N=2, Blue N=1.)

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (OathOfTheArknightRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	oathPlay(ge, l, self, 3)
}

func (OathOfTheArknightYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	oathPlay(ge, l, self, 2)
}

func (OathOfTheArknightBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	oathPlay(ge, l, self, 1)
}

// oathPlay grants +n to the first scheduled Runeblade attack, then creates one Runechant.
func oathPlay(ge card.GameEngine, l card.Logger, self *card.CardState, bonus int) {
	GrantNextCardBonusAttack(ge, bonus, card.IsRunebladeAttack)
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
}
