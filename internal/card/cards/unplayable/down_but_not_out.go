// Down But Not Out — Generic Action - Attack. Cost 3. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "When this attacks a hero, if you have less {h} and control fewer equipment and tokens
// than them, this gets +3{p}, **overpower**, and \"When this hits, create an Agility, Might, and
// Vigor token.\""

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (DownButNotOutRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (DownButNotOutYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (DownButNotOutBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
