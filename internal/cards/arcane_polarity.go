// Arcane Polarity — Generic Instant.
//
// Text: "Gain 1{h} If you've been dealt arcane damage this turn, instead gain N{h}."
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Gates the alternate gain on g.ArcaneIncomingDamage() > 0 (seeded from the matchup's
// -arcane-incoming). Life gain is credited 1-to-1 with damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// arcanePolarityPlay credits the conditional life gain as the chain step.
func arcanePolarityPlay(g card.GameEngine, l card.Logger, self *card.CardState, arcaneGain int) {
	gain := 1
	if g.ArcaneIncomingDamage() > 0 {
		gain = arcaneGain
	}
	g.AddValue(gain)
}

func (ArcanePolarityRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	arcanePolarityPlay(g, l, self, 4)
}

func (ArcanePolarityYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	arcanePolarityPlay(g, l, self, 3)
}

func (ArcanePolarityBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	arcanePolarityPlay(g, l, self, 2)
}
