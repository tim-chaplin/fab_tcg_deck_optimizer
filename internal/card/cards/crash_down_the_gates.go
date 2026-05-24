// Crash Down the Gates — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, they reveal the top card of their deck. If this has {p} greater
// than the revealed card, this gets +2{p}. When this hits a hero, destroy the top card of their
// deck."
//
// Neither rider is modelled: the single-turn solver doesn't track the opponent's deck, so the
// reveal-comparison +2{p} buff can't be evaluated and the on-hit deck-top destruction has no
// measurable single-turn effect. Modelled as the printed attack.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (CrashDownTheGatesRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState)    {}
func (CrashDownTheGatesYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (CrashDownTheGatesBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState)   {}
