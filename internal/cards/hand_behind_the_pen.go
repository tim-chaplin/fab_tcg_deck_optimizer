// Hand Behind the Pen — Generic Action - Attack. Cost 2, Pitch 1, Power 6, Defense 2. Only printed
// in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish a non-attack
// action card from their arsenal."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var handBehindThePenTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type HandBehindThePenRed struct{}

func (HandBehindThePenRed) ID() ids.CardID          { return ids.HandBehindThePenRed }
func (HandBehindThePenRed) Name() string            { return "Hand Behind the Pen" }
func (HandBehindThePenRed) Cost(*sim.TurnState) int { return 2 }
func (HandBehindThePenRed) Pitch() int              { return 1 }
func (HandBehindThePenRed) Attack() int             { return 6 }
func (HandBehindThePenRed) Defense() int            { return 2 }
func (HandBehindThePenRed) Types() card.TypeSet     { return handBehindThePenTypes }
func (HandBehindThePenRed) GoAgain() bool           { return false }

// not implemented: on-hit opponent-arsenal manipulation rider
func (HandBehindThePenRed) NotImplemented() {}
func (HandBehindThePenRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
