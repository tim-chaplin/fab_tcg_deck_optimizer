// Regurgitating Slog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Regurgitating Slog, you may banish a card named Sloggism
// from your graveyard. If you do, Regurgitating Slog gains **dominate**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func regurgitatingSlogPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := s.BanishFromGraveyard(isSloggism); ok {
		self.GrantedDominate = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a Sloggism, gained dominate", 0)
	}
}

func isSloggism(c card.Card) bool { return c.Name() == "Sloggism" }

func (RegurgitatingSlogRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	regurgitatingSlogPlay(s, l, self)
}

func (RegurgitatingSlogYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	regurgitatingSlogPlay(s, l, self)
}

func (RegurgitatingSlogBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	regurgitatingSlogPlay(s, l, self)
}
