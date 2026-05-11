// Arcane Polarity — Generic Instant.
//
// Text: "Gain 1{h} If you've been dealt arcane damage this turn, instead gain N{h}."
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Gates the alternate gain on s.ArcaneIncomingDamage > 0 (seeded from the matchup's
// -arcane-incoming). Life gain is credited 1-to-1 with damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// arcanePolarityPlay credits the conditional life gain as the chain step.
func arcanePolarityPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState, arcaneGain int) {
	gain := 1
	if s.ArcaneIncomingDamage > 0 {
		gain = arcaneGain
	}
	s.AddValue(gain)
	self.Log(l, gain)
}

func (ArcanePolarityRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	arcanePolarityPlay(s, l, self, 4)
}

func (ArcanePolarityYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	arcanePolarityPlay(s, l, self, 3)
}

func (ArcanePolarityBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	arcanePolarityPlay(s, l, self, 2)
}
