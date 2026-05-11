// Relentless Pursuit — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "**Mark** target opposing hero. If you've attacked them this turn, put this on the bottom
// of its owner's deck. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (RelentlessPursuitBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.SetOpponentMarked(true)
	recycled := s.HasPlayedType(card.TypeAttack)
	if recycled {
		s.RecycleToDeckBottom(self)
	}
	if recycled {
		l.AppendPostTrigger(self.Card.DisplayName(), "Recycled to bottom of deck", 0)
	}
}
