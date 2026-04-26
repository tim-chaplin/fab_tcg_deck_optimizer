// Toughen Up — Generic Defense Reaction. Cost 2, Pitch 3, Defense 4. Only printed in Blue.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type ToughenUpBlue struct{}

func (ToughenUpBlue) ID() card.ID              { return card.ToughenUpBlue }
func (ToughenUpBlue) Name() string             { return "Toughen Up" }
func (ToughenUpBlue) Cost(*card.TurnState) int { return 2 }
func (ToughenUpBlue) Pitch() int               { return 3 }
func (ToughenUpBlue) Attack() int              { return 0 }
func (ToughenUpBlue) Defense() int             { return 4 }
func (ToughenUpBlue) Types() card.TypeSet      { return defenseReactionTypes }
func (ToughenUpBlue) GoAgain() bool            { return false }
func (ToughenUpBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}
