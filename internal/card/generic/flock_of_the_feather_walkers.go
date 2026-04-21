// Flock of the Feather Walkers — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4,
// Blue 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Flock of the Feather Walkers, reveal a card in your hand
// with cost 1 or less. When you attack with Flock of the Feather Walkers, create a Quicken token."
//
// Simplification: Additional reveal cost and Quicken token creation aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var flockOfTheFeatherWalkersTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FlockOfTheFeatherWalkersRed struct{}

func (FlockOfTheFeatherWalkersRed) ID() card.ID                 { return card.FlockOfTheFeatherWalkersRed }
func (FlockOfTheFeatherWalkersRed) Name() string                { return "Flock of the Feather Walkers (Red)" }
func (FlockOfTheFeatherWalkersRed) Cost(*card.TurnState) int                   { return 1 }
func (FlockOfTheFeatherWalkersRed) Pitch() int                  { return 1 }
func (FlockOfTheFeatherWalkersRed) Attack() int                 { return 5 }
func (FlockOfTheFeatherWalkersRed) Defense() int                { return 2 }
func (FlockOfTheFeatherWalkersRed) Types() card.TypeSet         { return flockOfTheFeatherWalkersTypes }
func (FlockOfTheFeatherWalkersRed) GoAgain() bool               { return false }
func (c FlockOfTheFeatherWalkersRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type FlockOfTheFeatherWalkersYellow struct{}

func (FlockOfTheFeatherWalkersYellow) ID() card.ID                 { return card.FlockOfTheFeatherWalkersYellow }
func (FlockOfTheFeatherWalkersYellow) Name() string                { return "Flock of the Feather Walkers (Yellow)" }
func (FlockOfTheFeatherWalkersYellow) Cost(*card.TurnState) int                   { return 1 }
func (FlockOfTheFeatherWalkersYellow) Pitch() int                  { return 2 }
func (FlockOfTheFeatherWalkersYellow) Attack() int                 { return 4 }
func (FlockOfTheFeatherWalkersYellow) Defense() int                { return 2 }
func (FlockOfTheFeatherWalkersYellow) Types() card.TypeSet         { return flockOfTheFeatherWalkersTypes }
func (FlockOfTheFeatherWalkersYellow) GoAgain() bool               { return false }
func (c FlockOfTheFeatherWalkersYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }

type FlockOfTheFeatherWalkersBlue struct{}

func (FlockOfTheFeatherWalkersBlue) ID() card.ID                 { return card.FlockOfTheFeatherWalkersBlue }
func (FlockOfTheFeatherWalkersBlue) Name() string                { return "Flock of the Feather Walkers (Blue)" }
func (FlockOfTheFeatherWalkersBlue) Cost(*card.TurnState) int                   { return 1 }
func (FlockOfTheFeatherWalkersBlue) Pitch() int                  { return 3 }
func (FlockOfTheFeatherWalkersBlue) Attack() int                 { return 3 }
func (FlockOfTheFeatherWalkersBlue) Defense() int                { return 2 }
func (FlockOfTheFeatherWalkersBlue) Types() card.TypeSet         { return flockOfTheFeatherWalkersTypes }
func (FlockOfTheFeatherWalkersBlue) GoAgain() bool               { return false }
func (c FlockOfTheFeatherWalkersBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() }
