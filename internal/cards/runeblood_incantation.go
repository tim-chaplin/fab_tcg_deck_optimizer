// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates 1 Runechant per fire; Count=N gives N total fires before the aura dies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (RunebloodIncantationRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, runebloodAuraHandler, 3)
}

func (RunebloodIncantationYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, runebloodAuraHandler, 2)
}

func (RunebloodIncantationBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateStartOfTurnAura(self, runebloodAuraHandler, 1)
}

// runebloodAuraHandler creates 1 runechant per fire and decrements the verse counter.
// When the last verse fires, destroys the aura and graveyards the card.
func runebloodAuraHandler(g card.GameEngine, l card.Logger, a card.Aura) {
	name := a.CardName()
	lastVerse := a.DecrementCount() <= 0
	g.CreateRunechants(1)
	l.AppendPostTrigger(name, "Created a runechant (verse counter)", 1)
	if lastVerse {
		a.Destroy(true)
	}
}
