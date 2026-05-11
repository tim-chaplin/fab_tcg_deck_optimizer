// Pick a Card, Any Card — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Look at target opponent's hand then name a card. Choose a random card from their hand and
// reveal it. If it's the named card, create a Silver token. Repeat this process thrice. **Go
// again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: silver tokens, opponent hand inspection

func (PickACardAnyCardRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {}

// not implemented: silver tokens, opponent hand inspection

func (PickACardAnyCardYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
}

// not implemented: silver tokens, opponent hand inspection

func (PickACardAnyCardBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {}
