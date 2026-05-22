// Talisman of Recompense — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** Whenever you pitch a card, if you would gain exactly one {r}, instead destroy
// Talisman of Recompense and gain {r}{r}{r}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// recompenseFire is the pitch-trigger handler: a pitch that would yield exactly one
// resource instead yields three, and Talisman of Recompense is destroyed. Pitch values
// above one don't qualify, so the handler no-ops without consuming the talisman.
func recompenseFire(ge card.GameEngine, l card.Logger, self card.Item) {
	if ge.TriggeringCard().Pitch() != 1 {
		return
	}
	ge.AddResourcePoints(2)
	self.Destroy(true)
	l.AppendPostTrigger(self.CardName(), "Destroyed to make a 1-resource pitch yield 3", 2)
}

func (TalismanOfRecompenseYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreatePitchTriggeredItem(self, recompenseFire)
}

// MaxResourcePoints: one qualifying pitch yields {r}{r}{r} instead of {r}, adding two.
func (TalismanOfRecompenseYellow) MaxResourcePoints() int { return 2 }
