// Regurgitating Slog — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "As an additional cost to play Regurgitating Slog, you may banish a card named Sloggism
// from your graveyard. If you do, Regurgitating Slog gains **dominate**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func regurgitatingSlogPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if _, ok := ge.BanishFromGraveyard(isSloggism); ok {
		self.GrantedDominate = true
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished a Sloggism, gained dominate", 0)
	}
}

func isSloggism(_ card.GameEngine, pc *card.CardState) bool { return pc.Card.Name() == "Sloggism" }

func (RegurgitatingSlogRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	regurgitatingSlogPlay(ge, l, self)
}

func (RegurgitatingSlogYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	regurgitatingSlogPlay(ge, l, self)
}

func (RegurgitatingSlogBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	regurgitatingSlogPlay(ge, l, self)
}
