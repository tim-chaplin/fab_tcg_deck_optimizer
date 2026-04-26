// Gravekeeping — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, you may banish a card from their graveyard."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var gravekeepingTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type GravekeepingRed struct{}

func (GravekeepingRed) ID() card.ID              { return card.GravekeepingRed }
func (GravekeepingRed) Name() string             { return "Gravekeeping" }
func (GravekeepingRed) Cost(*card.TurnState) int { return 1 }
func (GravekeepingRed) Pitch() int               { return 1 }
func (GravekeepingRed) Attack() int              { return 5 }
func (GravekeepingRed) Defense() int             { return 2 }
func (GravekeepingRed) Types() card.TypeSet      { return gravekeepingTypes }
func (GravekeepingRed) GoAgain() bool            { return false }

// not implemented: opponent-graveyard banish rider
func (GravekeepingRed) NotImplemented() {}
func (c GravekeepingRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type GravekeepingYellow struct{}

func (GravekeepingYellow) ID() card.ID              { return card.GravekeepingYellow }
func (GravekeepingYellow) Name() string             { return "Gravekeeping" }
func (GravekeepingYellow) Cost(*card.TurnState) int { return 1 }
func (GravekeepingYellow) Pitch() int               { return 2 }
func (GravekeepingYellow) Attack() int              { return 4 }
func (GravekeepingYellow) Defense() int             { return 2 }
func (GravekeepingYellow) Types() card.TypeSet      { return gravekeepingTypes }
func (GravekeepingYellow) GoAgain() bool            { return false }

// not implemented: opponent-graveyard banish rider
func (GravekeepingYellow) NotImplemented() {}
func (c GravekeepingYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type GravekeepingBlue struct{}

func (GravekeepingBlue) ID() card.ID              { return card.GravekeepingBlue }
func (GravekeepingBlue) Name() string             { return "Gravekeeping" }
func (GravekeepingBlue) Cost(*card.TurnState) int { return 1 }
func (GravekeepingBlue) Pitch() int               { return 3 }
func (GravekeepingBlue) Attack() int              { return 3 }
func (GravekeepingBlue) Defense() int             { return 2 }
func (GravekeepingBlue) Types() card.TypeSet      { return gravekeepingTypes }
func (GravekeepingBlue) GoAgain() bool            { return false }

// not implemented: opponent-graveyard banish rider
func (GravekeepingBlue) NotImplemented() {}
func (c GravekeepingBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
