// Drone of Brutality — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Drone of Brutality would be put into your graveyard from anywhere, instead put it on
// the bottom of your deck."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (DroneOfBrutalityRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}

func (DroneOfBrutalityYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}

func (DroneOfBrutalityBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {}
