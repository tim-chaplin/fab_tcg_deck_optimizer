// Pound for Pound — Generic Action - Attack. Cost 3. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When you play Pound for Pound, if you have less {h} than an opposing hero, it gains
// **dominate**."
//
// Modelling: the "less {h} than an opposing hero" clause is treated as a hero attribute — the
// Dominate grant fires for heroes that implement card.LowerHealthWanter (via
// sim.HeroWantsLowerHealth) and never fires otherwise, a coarse proxy that skips per-turn
// life tracking. Standard self.GrantedDominate wiring (docs/dev-standards.md).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// poundForPoundPlay grants self Dominate when the current hero opts into LowerHealthWanter,
// then emits the chain step.
func poundForPoundPlay(s card.GameEngine, l card.Logger, self *card.CardState) {
	if sim.HeroWantsLowerHealth() {
		self.GrantedDominate = true
	}
}

func (PoundForPoundRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	poundForPoundPlay(s, l, self)
}

func (PoundForPoundYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	poundForPoundPlay(s, l, self)
}

func (PoundForPoundBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	poundForPoundPlay(s, l, self)
}
