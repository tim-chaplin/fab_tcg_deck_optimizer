// Smash Up — Generic Action - Attack. Cost 1, Pitch 1, Power 5, Defense 2. Only printed in Red.
//
// Text: "When this hits a hero, turn a card in their arsenal face-up, then banish an attack action
// card from their arsenal."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var smashUpTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SmashUpRed struct{}

func (SmashUpRed) ID() ids.CardID           { return ids.SmashUpRed }
func (SmashUpRed) Name() string             { return "Smash Up" }
func (SmashUpRed) Cost(*card.TurnState) int { return 1 }
func (SmashUpRed) Pitch() int               { return 1 }
func (SmashUpRed) Attack() int              { return 5 }
func (SmashUpRed) Defense() int             { return 2 }
func (SmashUpRed) Types() card.TypeSet      { return smashUpTypes }
func (SmashUpRed) GoAgain() bool            { return false }

// not implemented: on-hit opponent-arsenal manipulation rider
func (SmashUpRed) NotImplemented() {}
func (SmashUpRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
