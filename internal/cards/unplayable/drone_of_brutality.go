// Drone of Brutality — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Drone of Brutality would be put into your graveyard from anywhere, instead put it on
// the bottom of your deck."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (DroneOfBrutalityRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}

func (DroneOfBrutalityYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}

func (DroneOfBrutalityBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {}
