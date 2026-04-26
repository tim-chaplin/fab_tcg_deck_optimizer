// Bluster Buff — Generic Action - Attack. Cost 1, Pitch 1, Power 6, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var blusterBuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type BlusterBuffRed struct{}

func (BlusterBuffRed) ID() card.ID              { return card.BlusterBuffRed }
func (BlusterBuffRed) Name() string             { return "Bluster Buff" }
func (BlusterBuffRed) Cost(*card.TurnState) int { return 1 }
func (BlusterBuffRed) Pitch() int               { return 1 }
func (BlusterBuffRed) Attack() int              { return 6 }
func (BlusterBuffRed) Defense() int             { return 3 }
func (BlusterBuffRed) Types() card.TypeSet      { return blusterBuffTypes }
func (BlusterBuffRed) GoAgain() bool            { return false }

// not implemented: pay {r} or lose 1{p} resolved as 'always pay'
func (BlusterBuffRed) NotImplemented() {}
func (c BlusterBuffRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
