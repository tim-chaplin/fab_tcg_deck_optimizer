// Gravekeeping — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, you may banish a card from their graveyard."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var gravekeepingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type GravekeepingRed struct{}

func (GravekeepingRed) ID() ids.CardID          { return ids.GravekeepingRed }
func (GravekeepingRed) Name() string            { return "Gravekeeping" }
func (GravekeepingRed) Cost(*sim.TurnState) int { return 1 }
func (GravekeepingRed) Pitch() int              { return 1 }
func (GravekeepingRed) Attack() int             { return 5 }
func (GravekeepingRed) Defense() int            { return 2 }
func (GravekeepingRed) Types() card.TypeSet     { return gravekeepingTypes }
func (GravekeepingRed) GoAgain() bool           { return false }

// not implemented: opponent-graveyard banish rider
func (GravekeepingRed) NotImplemented() {}
func (c GravekeepingRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type GravekeepingYellow struct{}

func (GravekeepingYellow) ID() ids.CardID          { return ids.GravekeepingYellow }
func (GravekeepingYellow) Name() string            { return "Gravekeeping" }
func (GravekeepingYellow) Cost(*sim.TurnState) int { return 1 }
func (GravekeepingYellow) Pitch() int              { return 2 }
func (GravekeepingYellow) Attack() int             { return 4 }
func (GravekeepingYellow) Defense() int            { return 2 }
func (GravekeepingYellow) Types() card.TypeSet     { return gravekeepingTypes }
func (GravekeepingYellow) GoAgain() bool           { return false }

// not implemented: opponent-graveyard banish rider
func (GravekeepingYellow) NotImplemented() {}
func (c GravekeepingYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type GravekeepingBlue struct{}

func (GravekeepingBlue) ID() ids.CardID          { return ids.GravekeepingBlue }
func (GravekeepingBlue) Name() string            { return "Gravekeeping" }
func (GravekeepingBlue) Cost(*sim.TurnState) int { return 1 }
func (GravekeepingBlue) Pitch() int              { return 3 }
func (GravekeepingBlue) Attack() int             { return 3 }
func (GravekeepingBlue) Defense() int            { return 2 }
func (GravekeepingBlue) Types() card.TypeSet     { return gravekeepingTypes }
func (GravekeepingBlue) GoAgain() bool           { return false }

// not implemented: opponent-graveyard banish rider
func (GravekeepingBlue) NotImplemented() {}
func (c GravekeepingBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
