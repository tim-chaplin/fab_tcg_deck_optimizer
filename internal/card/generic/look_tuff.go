// Look Tuff — Generic Action - Attack. Cost 3, Pitch 1, Power 8, Defense 3. Only printed in Red.
//
// Text: "When this attacks, it gets -1{p} unless you pay {r}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var lookTuffTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type LookTuffRed struct{}

func (LookTuffRed) ID() card.ID              { return card.LookTuffRed }
func (LookTuffRed) Name() string             { return "Look Tuff" }
func (LookTuffRed) Cost(*card.TurnState) int { return 3 }
func (LookTuffRed) Pitch() int               { return 1 }
func (LookTuffRed) Attack() int              { return 8 }
func (LookTuffRed) Defense() int             { return 3 }
func (LookTuffRed) Types() card.TypeSet      { return lookTuffTypes }
func (LookTuffRed) GoAgain() bool            { return false }

// not implemented: pay {r} or lose 1{p} resolved as 'always pay'
func (LookTuffRed) NotImplemented() {}
func (c LookTuffRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
