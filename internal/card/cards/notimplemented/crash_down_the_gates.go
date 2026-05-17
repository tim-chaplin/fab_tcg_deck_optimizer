// Crash Down the Gates — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, they reveal the top card of their deck. If this has {p} greater
// than the revealed card, this gets +2{p}. When this hits a hero, destroy the top card of their
// deck."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// not implemented: deck-reveal comparison + on-hit deck-top destruction

func (CrashDownTheGatesRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: deck-reveal comparison + on-hit deck-top destruction

func (CrashDownTheGatesYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

// not implemented: deck-reveal comparison + on-hit deck-top destruction

func (CrashDownTheGatesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
