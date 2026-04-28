// Flock of the Feather Walkers — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4,
// Blue 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Flock of the Feather Walkers, reveal a card in your hand
// with cost 1 or less. When you attack with Flock of the Feather Walkers, create a Quicken token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var flockOfTheFeatherWalkersTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type FlockOfTheFeatherWalkersRed struct{}

func (FlockOfTheFeatherWalkersRed) ID() ids.CardID           { return ids.FlockOfTheFeatherWalkersRed }
func (FlockOfTheFeatherWalkersRed) Name() string             { return "Flock of the Feather Walkers" }
func (FlockOfTheFeatherWalkersRed) Cost(*card.TurnState) int { return 1 }
func (FlockOfTheFeatherWalkersRed) Pitch() int               { return 1 }
func (FlockOfTheFeatherWalkersRed) Attack() int              { return 5 }
func (FlockOfTheFeatherWalkersRed) Defense() int             { return 2 }
func (FlockOfTheFeatherWalkersRed) Types() card.TypeSet      { return flockOfTheFeatherWalkersTypes }
func (FlockOfTheFeatherWalkersRed) GoAgain() bool            { return false }

// not implemented: additional reveal cost, quicken tokens
func (FlockOfTheFeatherWalkersRed) NotImplemented() {}
func (c FlockOfTheFeatherWalkersRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type FlockOfTheFeatherWalkersYellow struct{}

func (FlockOfTheFeatherWalkersYellow) ID() ids.CardID           { return ids.FlockOfTheFeatherWalkersYellow }
func (FlockOfTheFeatherWalkersYellow) Name() string             { return "Flock of the Feather Walkers" }
func (FlockOfTheFeatherWalkersYellow) Cost(*card.TurnState) int { return 1 }
func (FlockOfTheFeatherWalkersYellow) Pitch() int               { return 2 }
func (FlockOfTheFeatherWalkersYellow) Attack() int              { return 4 }
func (FlockOfTheFeatherWalkersYellow) Defense() int             { return 2 }
func (FlockOfTheFeatherWalkersYellow) Types() card.TypeSet      { return flockOfTheFeatherWalkersTypes }
func (FlockOfTheFeatherWalkersYellow) GoAgain() bool            { return false }

// not implemented: additional reveal cost, quicken tokens
func (FlockOfTheFeatherWalkersYellow) NotImplemented() {}
func (c FlockOfTheFeatherWalkersYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type FlockOfTheFeatherWalkersBlue struct{}

func (FlockOfTheFeatherWalkersBlue) ID() ids.CardID           { return ids.FlockOfTheFeatherWalkersBlue }
func (FlockOfTheFeatherWalkersBlue) Name() string             { return "Flock of the Feather Walkers" }
func (FlockOfTheFeatherWalkersBlue) Cost(*card.TurnState) int { return 1 }
func (FlockOfTheFeatherWalkersBlue) Pitch() int               { return 3 }
func (FlockOfTheFeatherWalkersBlue) Attack() int              { return 3 }
func (FlockOfTheFeatherWalkersBlue) Defense() int             { return 2 }
func (FlockOfTheFeatherWalkersBlue) Types() card.TypeSet      { return flockOfTheFeatherWalkersTypes }
func (FlockOfTheFeatherWalkersBlue) GoAgain() bool            { return false }

// not implemented: additional reveal cost, quicken tokens
func (FlockOfTheFeatherWalkersBlue) NotImplemented() {}
func (c FlockOfTheFeatherWalkersBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
