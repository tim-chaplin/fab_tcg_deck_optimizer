// Clap 'Em in Irons — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When this enters the arena, {t} target Pirate hero or ally. It can't {u}
// while this is in the arena. At the start of your turn, destroy this."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// not implemented: passive tap-target Pirate; can't unfreeze; self-destroys at start of turn

func (ClapEmInIronsBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) { self.Log(l, 0) }
