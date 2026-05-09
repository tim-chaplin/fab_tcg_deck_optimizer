// Vantage Point — Runeblade Action - Attack.
//
// Text: "If you've played or created an aura this turn, this gets **overpower**."
//
// Credits sim.OverpowerValue (0) for the granted Overpower; flag still flips on s.Overpower
// for any future consumer.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// vantagePointPlay flips s.Overpower when an aura has been played or created this turn, then
// emits the chain step.
func vantagePointPlay(s *sim.TurnState, self *sim.CardState) {
	if s.HasPlayedOrCreatedAura() {
		s.Overpower = true
		s.AddValue(sim.OverpowerValue)
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

func (VantagePointRed) Play(s *sim.TurnState, self *sim.CardState) {
	vantagePointPlay(s, self)
}

func (VantagePointYellow) Play(s *sim.TurnState, self *sim.CardState) {
	vantagePointPlay(s, self)
}

func (VantagePointBlue) Play(s *sim.TurnState, self *sim.CardState) {
	vantagePointPlay(s, self)
}
