// Vantage Point — Runeblade Action - Attack.
//
// Text: "If you've played or created an aura this turn, this gets **overpower**."
//
// Credits sim.OverpowerValue (0) for the granted Overpower; flag still flips on s.Overpower()
// for any future consumer.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// vantagePointPlay flips s.Overpower() when an aura has been played or created this turn, then
// emits the chain step.
func vantagePointPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if s.HasPlayedOrCreatedAura() {
		s.SetOverpower(true)
		s.AddValue(sim.OverpowerValue)
	}
}

func (VantagePointRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(s, l, self)
}

func (VantagePointYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(s, l, self)
}

func (VantagePointBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	vantagePointPlay(s, l, self)
}
