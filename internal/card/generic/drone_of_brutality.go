// Drone of Brutality — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Drone of Brutality would be put into your graveyard from anywhere, instead put it on
// the bottom of your deck."
//
// Simplification: Graveyard-replacement-to-deck trigger isn't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var droneOfBrutalityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DroneOfBrutalityRed struct{}

func (DroneOfBrutalityRed) ID() card.ID                 { return card.DroneOfBrutalityRed }
func (DroneOfBrutalityRed) Name() string                { return "Drone of Brutality (Red)" }
func (DroneOfBrutalityRed) Cost() int                   { return 2 }
func (DroneOfBrutalityRed) Pitch() int                  { return 1 }
func (DroneOfBrutalityRed) Attack() int                 { return 6 }
func (DroneOfBrutalityRed) Defense() int                { return 2 }
func (DroneOfBrutalityRed) Types() card.TypeSet         { return droneOfBrutalityTypes }
func (DroneOfBrutalityRed) GoAgain() bool               { return false }
func (DroneOfBrutalityRed) NotSilverAgeLegal()           {}
func (c DroneOfBrutalityRed) Play(s *card.TurnState) int { return c.Attack() }

type DroneOfBrutalityYellow struct{}

func (DroneOfBrutalityYellow) ID() card.ID                 { return card.DroneOfBrutalityYellow }
func (DroneOfBrutalityYellow) Name() string                { return "Drone of Brutality (Yellow)" }
func (DroneOfBrutalityYellow) Cost() int                   { return 2 }
func (DroneOfBrutalityYellow) Pitch() int                  { return 2 }
func (DroneOfBrutalityYellow) Attack() int                 { return 5 }
func (DroneOfBrutalityYellow) Defense() int                { return 2 }
func (DroneOfBrutalityYellow) Types() card.TypeSet         { return droneOfBrutalityTypes }
func (DroneOfBrutalityYellow) GoAgain() bool               { return false }
func (DroneOfBrutalityYellow) NotSilverAgeLegal()           {}
func (c DroneOfBrutalityYellow) Play(s *card.TurnState) int { return c.Attack() }

type DroneOfBrutalityBlue struct{}

func (DroneOfBrutalityBlue) ID() card.ID                 { return card.DroneOfBrutalityBlue }
func (DroneOfBrutalityBlue) Name() string                { return "Drone of Brutality (Blue)" }
func (DroneOfBrutalityBlue) Cost() int                   { return 2 }
func (DroneOfBrutalityBlue) Pitch() int                  { return 3 }
func (DroneOfBrutalityBlue) Attack() int                 { return 4 }
func (DroneOfBrutalityBlue) Defense() int                { return 2 }
func (DroneOfBrutalityBlue) Types() card.TypeSet         { return droneOfBrutalityTypes }
func (DroneOfBrutalityBlue) GoAgain() bool               { return false }
func (DroneOfBrutalityBlue) NotSilverAgeLegal()           {}
func (c DroneOfBrutalityBlue) Play(s *card.TurnState) int { return c.Attack() }
