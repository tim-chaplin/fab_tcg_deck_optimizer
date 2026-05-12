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

func vantagePointPlay(g card.GameEngine, _ card.Logger, self *card.CardState) {
	if g.AuraCreated() {
		self.GrantedOverpower = true
	}
}

func (VantagePointRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(g, l, self)
}

func (VantagePointYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(g, l, self)
}

func (VantagePointBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(g, l, self)
}
