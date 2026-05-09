// Scour the Battlescape — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on the bottom of your deck. If you do, draw a card. If
// Scour the Battlescape is played from arsenal, it gains **go again**."
//
// Modelling: hand-cycle isn't modelled. Standard played-from-arsenal go-again
// (docs/dev-standards.md).

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func scourTheBattlescapePlay(s *sim.TurnState, self *sim.CardState) {
	self.GrantGoAgainIfFromArsenal()
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

// not implemented: hand-cycle rider (put a card on bottom of deck, draw)

func (ScourTheBattlescapeRed) Play(s *sim.TurnState, self *sim.CardState) {
	scourTheBattlescapePlay(s, self)
}

// not implemented: hand-cycle rider (put a card on bottom of deck, draw)

func (ScourTheBattlescapeYellow) Play(s *sim.TurnState, self *sim.CardState) {
	scourTheBattlescapePlay(s, self)
}

// not implemented: hand-cycle rider (put a card on bottom of deck, draw)

func (ScourTheBattlescapeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	scourTheBattlescapePlay(s, self)
}
