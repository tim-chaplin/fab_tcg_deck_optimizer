// Arcane Polarity — Generic Instant.
//
// Text: "Gain 1{h} If you've been dealt arcane damage this turn, instead gain N{h}."
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Gates the alternate gain on ge.ArcaneIncomingDamage() > 0 (seeded from the matchup's
// -arcane-incoming). Life gain is credited 1-to-1 with damage.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// arcanePolarityPlay credits the conditional life gain as the chain step.
func arcanePolarityPlay(ge card.GameEngine, l card.Logger, self *card.CardState, arcaneGain int) {
	gain := 1
	if ge.ArcaneIncomingDamage() > 0 {
		gain = arcaneGain
	}
	ge.AddValue(gain)
}

func (ArcanePolarityRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	arcanePolarityPlay(ge, l, self, 4)
}

func (ArcanePolarityYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	arcanePolarityPlay(ge, l, self, 3)
}

func (ArcanePolarityBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	arcanePolarityPlay(ge, l, self, 2)
}
