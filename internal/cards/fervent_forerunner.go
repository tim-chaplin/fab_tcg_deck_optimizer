// Fervent Forerunner — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Fervent Forerunner hits, **opt 2**. If Fervent Forerunner is played from arsenal, it
// gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func ferventForerunnerPlay(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	self.GrantGoAgainIfFromArsenal()
	self.RegisterOnHit(ferventForerunnerOnHit)
}

// ferventForerunnerOnHit fires the printed "If this hits, opt 2" rider. Top-level so
// registration stays alloc-free.
func ferventForerunnerOnHit(s sim.GameEngine, l sim.Logger, _ *sim.CardState, _ *sim.OnHitHandler) {
	s.Opt(l, 2)
}

func (FerventForerunnerRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	ferventForerunnerPlay(s, l, self)
}

func (FerventForerunnerYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	ferventForerunnerPlay(s, l, self)
}

func (FerventForerunnerBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	ferventForerunnerPlay(s, l, self)
}
