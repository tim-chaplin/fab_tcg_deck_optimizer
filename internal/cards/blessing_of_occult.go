// Blessing of Occult — Runeblade Action - Aura. Cost 1, Defense 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "At the start of your turn, destroy Blessing of Occult then create N Runechant tokens."
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates N Runechants next turn.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// blessingOfOccultTriggerText pre-formats the trigger log text for each Runechant count
// (1 = Blue, 2 = Yellow, 3 = Red). The text is captured into a per-play closure on every
// cast, so a constant lookup avoids a Sprintf alloc per cast.
var blessingOfOccultTriggerText = [...]string{
	1: "Created a runechant",
	2: "Created 2 runechants",
	3: "Created 3 runechants",
}

func (BlessingOfOccultRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, blessingOfOccultHandler, 3)
}

func (BlessingOfOccultYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, blessingOfOccultHandler, 2)
}

func (BlessingOfOccultBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, blessingOfOccultHandler, 1)
}

// blessingOfOccultHandler creates a.Count() Runechants and destroys the aura. Count
// carries the per-variant rune count (R=3 / Y=2 / B=1) — the handler is one-shot, so
// Count's "fires remaining" interpretation collapses to "runechants to create on the
// only fire".
func blessingOfOccultHandler(g card.GameEngine, l card.Logger, a card.Aura) {
	n := a.Count()
	g.CreateRunechants(n)
	l.AppendPostTrigger(a.CardName(), blessingOfOccultTriggerText[n], n)
	a.Destroy(true)
}
