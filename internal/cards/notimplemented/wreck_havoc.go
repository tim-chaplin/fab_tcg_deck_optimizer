// Wreck Havoc — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "Defense reactions can't be played to this chain link. When this hits a hero, you may turn
// a card in their arsenal face up, then destroy a defense reaction in their arsenal."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: defense-reaction lockout, on-hit arsenal banish

func (WreckHavocRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}

// not implemented: defense-reaction lockout, on-hit arsenal banish

func (WreckHavocYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}

// not implemented: defense-reaction lockout, on-hit arsenal banish

func (WreckHavocBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
}
