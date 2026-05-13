// Sirens of Safe Harbor — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is put into your graveyard from anywhere, gain 1{h}."
//
// Modelling: attacking sends this to the graveyard, firing the 1{h} gain (+1 damage
// equivalent). Pitched copies bypass the rider — they go to the bottom of the deck.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SirensOfSafeHarborRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.AddValue(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Gained 1 health (graveyard trigger)", 1)
}

func (SirensOfSafeHarborYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.AddValue(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Gained 1 health (graveyard trigger)", 1)
}

func (SirensOfSafeHarborBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.AddValue(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "Gained 1 health (graveyard trigger)", 1)
}
