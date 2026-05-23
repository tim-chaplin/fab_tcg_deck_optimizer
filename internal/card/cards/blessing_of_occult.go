// Blessing of Occult — Runeblade Action - Aura. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "At the start of your turn, destroy Blessing of Occult then create N Runechant tokens."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates N Runechants next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// blessingOfOccultTriggerText pre-formats the trigger log text for each Runechant count
// (1 = Blue, 2 = Yellow, 3 = Red).
var blessingOfOccultTriggerText = [...]string{
	1: "Created a runechant",
	2: "Created 2 runechants",
	3: "Created 3 runechants",
}

func (BlessingOfOccultRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.StartOfTurn, blessingOfOccultHandler, 3, false, nil)
}

func (BlessingOfOccultYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.StartOfTurn, blessingOfOccultHandler, 2, false, nil)
}

func (BlessingOfOccultBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self, triggertype.StartOfTurn, blessingOfOccultHandler, 1, false, nil)
}

// blessingOfOccultHandler creates a.Count() Runechants and destroys the aura. Count carries
// the per-variant rune count (R=3 / Y=2 / B=1).
func blessingOfOccultHandler(ge card.GameEngine, l card.Logger, a card.Aura) {
	n := a.Count()
	ge.CreateRunechants(n)
	l.AppendPostTrigger(a.CardName(), blessingOfOccultTriggerText[n], n)
	a.Destroy(true)
}
