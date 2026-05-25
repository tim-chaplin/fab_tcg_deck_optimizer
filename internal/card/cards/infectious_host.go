// Infectious Host — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, if you control a Frailty token, create a Frailty token under
// their control, then repeat for Inertia and Bloodrot Pox."
//
// The spread rider isn't modelled: no implemented card grants the active player Frailty,
// Inertia, or Bloodrot Pox tokens (it's a sideboard answer to opponents that play them),
// so the gates never fire and the card behaves like a vanilla 4/3/2 attack. The opposing-
// side mint plumbing (CreateFrailtyForOpponent / CreateInertiaForOpponent /
// CreateBloodrotPoxForOpponent) exists for cards that mint these tokens directly on the
// opponent.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (InfectiousHostRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    {}
func (InfectiousHostYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (InfectiousHostBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   {}
