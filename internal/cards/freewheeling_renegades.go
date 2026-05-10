// Freewheeling Renegades — Generic Action - Attack. Cost 1. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If this is defended by an action card, this has -2{p}."
//
// Conservative model: the -2{p} self-debuff fires unconditionally — assume the defender
// always blocks with an action card.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func freewheelingRenegadesPlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	self.BonusAttack -= 2
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

func (FreewheelingRenegadesRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	freewheelingRenegadesPlay(s, l, self)
}

func (FreewheelingRenegadesYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	freewheelingRenegadesPlay(s, l, self)
}

func (FreewheelingRenegadesBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	freewheelingRenegadesPlay(s, l, self)
}
