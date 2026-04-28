// Drone of Brutality — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Drone of Brutality would be put into your graveyard from anywhere, instead put it on
// the bottom of your deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var droneOfBrutalityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DroneOfBrutalityRed struct{}

func (DroneOfBrutalityRed) ID() ids.CardID           { return ids.DroneOfBrutalityRed }
func (DroneOfBrutalityRed) Name() string             { return "Drone of Brutality" }
func (DroneOfBrutalityRed) Cost(*card.TurnState) int { return 2 }
func (DroneOfBrutalityRed) Pitch() int               { return 1 }
func (DroneOfBrutalityRed) Attack() int              { return 6 }
func (DroneOfBrutalityRed) Defense() int             { return 2 }
func (DroneOfBrutalityRed) Types() card.TypeSet      { return droneOfBrutalityTypes }
func (DroneOfBrutalityRed) GoAgain() bool            { return false }
func (DroneOfBrutalityRed) NotSilverAgeLegal()       {}

// not implemented: graveyard-replacement-to-deck trigger
func (DroneOfBrutalityRed) NotImplemented() {}
func (c DroneOfBrutalityRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type DroneOfBrutalityYellow struct{}

func (DroneOfBrutalityYellow) ID() ids.CardID           { return ids.DroneOfBrutalityYellow }
func (DroneOfBrutalityYellow) Name() string             { return "Drone of Brutality" }
func (DroneOfBrutalityYellow) Cost(*card.TurnState) int { return 2 }
func (DroneOfBrutalityYellow) Pitch() int               { return 2 }
func (DroneOfBrutalityYellow) Attack() int              { return 5 }
func (DroneOfBrutalityYellow) Defense() int             { return 2 }
func (DroneOfBrutalityYellow) Types() card.TypeSet      { return droneOfBrutalityTypes }
func (DroneOfBrutalityYellow) GoAgain() bool            { return false }
func (DroneOfBrutalityYellow) NotSilverAgeLegal()       {}

// not implemented: graveyard-replacement-to-deck trigger
func (DroneOfBrutalityYellow) NotImplemented() {}
func (c DroneOfBrutalityYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type DroneOfBrutalityBlue struct{}

func (DroneOfBrutalityBlue) ID() ids.CardID           { return ids.DroneOfBrutalityBlue }
func (DroneOfBrutalityBlue) Name() string             { return "Drone of Brutality" }
func (DroneOfBrutalityBlue) Cost(*card.TurnState) int { return 2 }
func (DroneOfBrutalityBlue) Pitch() int               { return 3 }
func (DroneOfBrutalityBlue) Attack() int              { return 4 }
func (DroneOfBrutalityBlue) Defense() int             { return 2 }
func (DroneOfBrutalityBlue) Types() card.TypeSet      { return droneOfBrutalityTypes }
func (DroneOfBrutalityBlue) GoAgain() bool            { return false }
func (DroneOfBrutalityBlue) NotSilverAgeLegal()       {}

// not implemented: graveyard-replacement-to-deck trigger
func (DroneOfBrutalityBlue) NotImplemented() {}
func (c DroneOfBrutalityBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
