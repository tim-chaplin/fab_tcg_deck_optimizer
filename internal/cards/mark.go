package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"

// markOpponentOnHit fires the printed "When this hits a hero, mark them" rider.
func markOpponentOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	s.OpponentMarked = true
	s.LogRider(self, 0, "Marked the opposing hero")
}
