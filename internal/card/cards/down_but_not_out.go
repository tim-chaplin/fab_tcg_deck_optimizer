// Down But Not Out — Generic Action - Attack. Cost 3. Printed power: Red 5, Yellow 4, Blue 3
// (pitch 1/2/3). Defense 3.
// Text: "When this attacks a hero, if you have less {h} and control fewer equipment and tokens
// than them, this gets +3{p}, **overpower**, and \"When this hits, create an Agility, Might, and
// Vigor token.\""
//
// The conditional block (+3{p}, overpower, and the on-hit Agility/Might/Vigor token rider) is not
// modelled: it is gated on a comparison of our life and board against the opponent's, which the
// single-turn model doesn't track, so the optimiser's worst-case opponent assumption leaves it
// inactive. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (DownButNotOutRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (DownButNotOutYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}

func (DownButNotOutBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
