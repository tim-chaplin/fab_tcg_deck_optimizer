// Regurgitating Slog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Regurgitating Slog, you may banish a card named Sloggism
// from your graveyard. If you do, Regurgitating Slog gains **dominate**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func regurgitatingSlogPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	if _, ok := s.BanishFromGraveyard(isSloggism); ok {
		self.GrantedDominate = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a Sloggism, gained dominate", 0)
	}
}

func isSloggism(c sim.Card) bool { return c.Name() == "Sloggism" }

func (RegurgitatingSlogRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	regurgitatingSlogPlay(s, l, self)
}

func (RegurgitatingSlogYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	regurgitatingSlogPlay(s, l, self)
}

func (RegurgitatingSlogBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	regurgitatingSlogPlay(s, l, self)
}
