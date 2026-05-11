// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates 1 Runechant per fire; Count=N gives N total fires before the aura dies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (c RunebloodIncantationRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	runebloodPlay(s, l, self, c, 3)
}

func (c RunebloodIncantationYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	runebloodPlay(s, l, self, c, 2)
}

func (c RunebloodIncantationBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	runebloodPlay(s, l, self, c, 1)
}

// runebloodAuraHandler creates 1 runechant per fire and decrements the verse counter.
// When the last verse fires, destroys the aura and graveyards the card.
func runebloodAuraHandler(s *sim.TurnState, l card.Logger, _ *sim.Trigger, a *sim.Aura) {
	name := a.Self.DisplayName()
	a.Count--
	lastVerse := a.Count <= 0
	s.CreateRunechants(1)
	l.AppendPostTrigger(name, "Created a runechant (verse counter)", 1)
	if lastVerse {
		s.DestroyAura(a, true)
	}
}

// runebloodPlay registers a start-of-turn trigger with Count=n and emits the same-turn
// chain step (no value contribution; every rune is credited at its future-turn fire).
func runebloodPlay(s card.GameEngine, l card.Logger, selfState *card.CardState, selfCard sim.Card, n int) {
	s.AddAura(sim.Aura{
		Trigger: sim.Trigger{TriggerType: sim.TriggerStartOfTurn, Handler: runebloodAuraHandler},
		Self:    sim.CardOrTokenType{Card: selfCard},
		Count:   n,
	})
}
