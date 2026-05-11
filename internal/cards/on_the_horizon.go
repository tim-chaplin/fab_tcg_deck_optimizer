// On the Horizon — Generic Block. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Printed
// defense: Red 4, Yellow 3, Blue 2.
//
// Text: "When this defends, look at the top card of your deck."
//
// Block-typed: only legal roles are pitch and plain block, so Play is never invoked by
// the chain runner (the partition enumerator forbids Attack via the Action / Weapon
// gate). The deck-peek defend trigger isn't modelled — it surfaces information for the
// player, not a state change the solver can credit.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (OnTheHorizonRed) Play(card.GameEngine, card.Logger, *card.CardState) {}

func (OnTheHorizonYellow) Play(card.GameEngine, card.Logger, *card.CardState) {}

func (OnTheHorizonBlue) Play(card.GameEngine, card.Logger, *card.CardState) {}
