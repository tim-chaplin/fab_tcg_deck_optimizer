// Fervent Forerunner — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If Fervent Forerunner hits, **opt 2**. If Fervent Forerunner is played from arsenal, it
// gains **go again**."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func ferventForerunnerPlay(g card.GameEngine, l card.Logger, self *card.CardState) {
	self.GrantGoAgainIfFromArsenal()
	self.RegisterOnHit(ferventForerunnerOnHit)
}

// ferventForerunnerOnHit fires the printed "If this hits, opt 2" rider.
func ferventForerunnerOnHit(g card.GameEngine, l card.Logger, _ *card.CardState, _ *card.OnHitHandler) {
	g.Opt(l, 2)
}

func (FerventForerunnerRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	ferventForerunnerPlay(g, l, self)
}

func (FerventForerunnerYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	ferventForerunnerPlay(g, l, self)
}

func (FerventForerunnerBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	ferventForerunnerPlay(g, l, self)
}
