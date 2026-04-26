// Regain Composure — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next attack this turn gets +1{p} and "When this hits, {u} your hero." **Go again**"

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var regainComposureTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type RegainComposureBlue struct{}

func (RegainComposureBlue) ID() card.ID              { return card.RegainComposureBlue }
func (RegainComposureBlue) Name() string             { return "Regain Composure" }
func (RegainComposureBlue) Cost(*card.TurnState) int { return 0 }
func (RegainComposureBlue) Pitch() int               { return 3 }
func (RegainComposureBlue) Attack() int              { return 0 }
func (RegainComposureBlue) Defense() int             { return 2 }
func (RegainComposureBlue) Types() card.TypeSet      { return regainComposureTypes }
func (RegainComposureBlue) GoAgain() bool            { return true }

// not implemented: on-hit unfreeze rider (freeze/unfreeze state not tracked)
func (RegainComposureBlue) NotImplemented() {}
func (RegainComposureBlue) Play(s *card.TurnState, self *card.CardState) {
	grantNextAttackActionBonus(s, 1)
	s.ApplyAndLogEffectiveAttack(self)
}
