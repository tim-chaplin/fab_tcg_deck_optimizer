// Out Muscle — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "While Out Muscle isn't defended by a card with equal or greater {p}, it has **go again**."
//
// Conservative model: the conditional go-again is dropped — assume the defender always
// blocks with an equal-or-greater {p} card.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func outMusclePlay(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
}

func (OutMuscleRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	outMusclePlay(s, l, self)
}

func (OutMuscleYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	outMusclePlay(s, l, self)
}

func (OutMuscleBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	outMusclePlay(s, l, self)
}
