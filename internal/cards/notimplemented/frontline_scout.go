// Frontline Scout — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may look at the defending hero's hand. If Frontline Scout is played from arsenal, it
// gains **go again**."
//
// Modelling: hand-peek isn't modelled. Standard played-from-arsenal go-again
// (docs/dev-standards.md).

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func frontlineScoutPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.GrantGoAgainIfFromArsenal()
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

// not implemented: opposing-hand-peek rider

func (FrontlineScoutRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	frontlineScoutPlay(s, l, self)
}

// not implemented: opposing-hand-peek rider

func (FrontlineScoutYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	frontlineScoutPlay(s, l, self)
}

// not implemented: opposing-hand-peek rider

func (FrontlineScoutBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	frontlineScoutPlay(s, l, self)
}
