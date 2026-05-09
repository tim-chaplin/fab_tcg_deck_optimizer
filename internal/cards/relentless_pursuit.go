// Relentless Pursuit — Generic Action. Cost 0, Pitch 3, Defense 3. Only printed in Blue.
//
// Text: "**Mark** target opposing hero. If you've attacked them this turn, put this on the bottom
// of its owner's deck. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (RelentlessPursuitBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.OpponentMarked = true
	recycled := s.HasPlayedType(card.TypeAttack)
	if recycled {
		s.RecycleToDeckBottom(self)
	}
	s.Log(self, 0)
	if recycled {
		s.LogRider(self, 0, "Recycled to bottom of deck")
	}
}
