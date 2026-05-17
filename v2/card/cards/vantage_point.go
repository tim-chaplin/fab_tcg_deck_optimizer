// Vantage Point — Runeblade Action - Attack.
//
// Text: "If you've played or created an aura this turn, this gets **overpower**."
//
// Flips self.GrantedOverpower when an aura has been created this turn. The engine doesn't
// currently fold Overpower into a per-step bonus — block allocation is accounted for at the
// partition level.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func vantagePointPlay(ge card.GameEngine, _ card.Logger, self *card.CardState) {
	if ge.AuraCreated() {
		self.GrantedOverpower = true
	}
}

func (VantagePointRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(ge, l, self)
}

func (VantagePointYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(ge, l, self)
}

func (VantagePointBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(ge, l, self)
}
