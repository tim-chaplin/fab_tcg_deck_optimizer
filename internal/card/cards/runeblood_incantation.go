// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Count tracks verse counters: each start-of-turn fire removes one for a Runechant; the
// first fire that finds none left destroys the aura instead.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (RunebloodIncantationRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, runebloodAuraHandler, 3)
}

func (RunebloodIncantationYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, runebloodAuraHandler, 2)
}

func (RunebloodIncantationBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateStartOfTurnAura(self, runebloodAuraHandler, 1)
}

// runebloodAuraHandler removes a verse counter each fire to create a runechant; a fire
// that finds no counter left destroys the aura and graveyards the card.
func runebloodAuraHandler(ge card.GameEngine, l card.Logger, a card.Aura) {
	if a.Count() <= 0 {
		a.Destroy(true)
		return
	}
	a.DecrementCount()
	ge.CreateRunechants(1)
	l.AppendPostTrigger(a.CardName(), "Created a runechant (verse counter)", 1)
}
