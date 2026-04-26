// Muscle Mutt — Generic Action - Attack. Cost 3, Pitch 2, Power 6, Defense 2. Only printed in
// Yellow.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var muscleMuttTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type MuscleMuttYellow struct{}

func (MuscleMuttYellow) ID() card.ID              { return card.MuscleMuttYellow }
func (MuscleMuttYellow) Name() string             { return "Muscle Mutt" }
func (MuscleMuttYellow) Cost(*card.TurnState) int { return 3 }
func (MuscleMuttYellow) Pitch() int               { return 2 }
func (MuscleMuttYellow) Attack() int              { return 6 }
func (MuscleMuttYellow) Defense() int             { return 2 }
func (MuscleMuttYellow) Types() card.TypeSet      { return muscleMuttTypes }
func (MuscleMuttYellow) GoAgain() bool            { return false }
func (c MuscleMuttYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
