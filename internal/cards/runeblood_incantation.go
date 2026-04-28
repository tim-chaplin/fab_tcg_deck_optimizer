// Runeblood Incantation — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. Runeblood Incantation enters the arena with N verse counters on it. At the
// beginning of your action phase, remove a verse counter. If you do, create a Runechant token.
// Otherwise, destroy Runeblood Incantation." (Red N=3, Yellow N=2, Blue N=1.)
//
// Handler creates 1 Runechant per fire; Count=N gives N total fires before the aura dies.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var runebloodIncantationTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type RunebloodIncantationRed struct{}

func (RunebloodIncantationRed) ID() ids.CardID          { return ids.RunebloodIncantationRed }
func (RunebloodIncantationRed) Name() string            { return "Runeblood Incantation" }
func (RunebloodIncantationRed) Cost(*sim.TurnState) int { return 1 }
func (RunebloodIncantationRed) Pitch() int              { return 1 }
func (RunebloodIncantationRed) Attack() int             { return 0 }
func (RunebloodIncantationRed) Defense() int            { return 2 }
func (RunebloodIncantationRed) Types() card.TypeSet     { return runebloodIncantationTypes }
func (RunebloodIncantationRed) GoAgain() bool           { return true }
func (RunebloodIncantationRed) AddsFutureValue()        {}
func (c RunebloodIncantationRed) Play(s *sim.TurnState, self *sim.CardState) {
	runebloodPlay(s, self, c, 3)
}

type RunebloodIncantationYellow struct{}

func (RunebloodIncantationYellow) ID() ids.CardID          { return ids.RunebloodIncantationYellow }
func (RunebloodIncantationYellow) Name() string            { return "Runeblood Incantation" }
func (RunebloodIncantationYellow) Cost(*sim.TurnState) int { return 1 }
func (RunebloodIncantationYellow) Pitch() int              { return 2 }
func (RunebloodIncantationYellow) Attack() int             { return 0 }
func (RunebloodIncantationYellow) Defense() int            { return 2 }
func (RunebloodIncantationYellow) Types() card.TypeSet     { return runebloodIncantationTypes }
func (RunebloodIncantationYellow) GoAgain() bool           { return true }
func (RunebloodIncantationYellow) AddsFutureValue()        {}
func (c RunebloodIncantationYellow) Play(s *sim.TurnState, self *sim.CardState) {
	runebloodPlay(s, self, c, 2)
}

type RunebloodIncantationBlue struct{}

func (RunebloodIncantationBlue) ID() ids.CardID          { return ids.RunebloodIncantationBlue }
func (RunebloodIncantationBlue) Name() string            { return "Runeblood Incantation" }
func (RunebloodIncantationBlue) Cost(*sim.TurnState) int { return 1 }
func (RunebloodIncantationBlue) Pitch() int              { return 3 }
func (RunebloodIncantationBlue) Attack() int             { return 0 }
func (RunebloodIncantationBlue) Defense() int            { return 2 }
func (RunebloodIncantationBlue) Types() card.TypeSet     { return runebloodIncantationTypes }
func (RunebloodIncantationBlue) GoAgain() bool           { return true }
func (RunebloodIncantationBlue) AddsFutureValue()        {}
func (c RunebloodIncantationBlue) Play(s *sim.TurnState, self *sim.CardState) {
	runebloodPlay(s, self, c, 1)
}

// runebloodPlay registers a start-of-turn trigger with Count=n and emits the same-turn
// chain step (no value contribution; every rune is credited at its future-turn fire).
func runebloodPlay(s *sim.TurnState, selfState *sim.CardState, selfCard sim.Card, n int) {
	s.RegisterStartOfTurn(selfCard, n, "Created a runechant (verse counter)", func(s *sim.TurnState) int { return s.CreateRunechants(1) })
	s.LogPlay(selfState)
}
