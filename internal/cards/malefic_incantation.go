// Malefic Incantation — Runeblade Action - Aura. Cost 0, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "This enters the arena with N verse counters. When it has none, destroy it. Once per
// turn, when you play an attack action card, remove a verse counter from this. If you do,
// create a Runechant token." (Red N=3, Yellow N=2, Blue N=1.)
//
// AttackAction trigger with Count=N and OncePerTurn=true: each turn's first attack action
// creates 1 Runechant and burns one verse counter.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var maleficTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type MaleficIncantationRed struct{}

func (MaleficIncantationRed) ID() ids.CardID          { return ids.MaleficIncantationRed }
func (MaleficIncantationRed) Name() string            { return "Malefic Incantation" }
func (MaleficIncantationRed) Cost(*sim.TurnState) int { return 0 }
func (MaleficIncantationRed) Pitch() int              { return 1 }
func (MaleficIncantationRed) Attack() int             { return 0 }
func (MaleficIncantationRed) Defense() int            { return 2 }
func (MaleficIncantationRed) Types() card.TypeSet     { return maleficTypes }
func (MaleficIncantationRed) GoAgain() bool           { return true }
func (MaleficIncantationRed) AddsFutureValue()        {}
func (c MaleficIncantationRed) Play(s *sim.TurnState, self *sim.CardState) {
	maleficPlay(s, self, c, 3)
}

type MaleficIncantationYellow struct{}

func (MaleficIncantationYellow) ID() ids.CardID          { return ids.MaleficIncantationYellow }
func (MaleficIncantationYellow) Name() string            { return "Malefic Incantation" }
func (MaleficIncantationYellow) Cost(*sim.TurnState) int { return 0 }
func (MaleficIncantationYellow) Pitch() int              { return 2 }
func (MaleficIncantationYellow) Attack() int             { return 0 }
func (MaleficIncantationYellow) Defense() int            { return 2 }
func (MaleficIncantationYellow) Types() card.TypeSet     { return maleficTypes }
func (MaleficIncantationYellow) GoAgain() bool           { return true }
func (MaleficIncantationYellow) AddsFutureValue()        {}
func (c MaleficIncantationYellow) Play(s *sim.TurnState, self *sim.CardState) {
	maleficPlay(s, self, c, 2)
}

type MaleficIncantationBlue struct{}

func (MaleficIncantationBlue) ID() ids.CardID          { return ids.MaleficIncantationBlue }
func (MaleficIncantationBlue) Name() string            { return "Malefic Incantation" }
func (MaleficIncantationBlue) Cost(*sim.TurnState) int { return 0 }
func (MaleficIncantationBlue) Pitch() int              { return 3 }
func (MaleficIncantationBlue) Attack() int             { return 0 }
func (MaleficIncantationBlue) Defense() int            { return 2 }
func (MaleficIncantationBlue) Types() card.TypeSet     { return maleficTypes }
func (MaleficIncantationBlue) GoAgain() bool           { return true }
func (MaleficIncantationBlue) AddsFutureValue()        {}
func (c MaleficIncantationBlue) Play(s *sim.TurnState, self *sim.CardState) {
	maleficPlay(s, self, c, 1)
}

// maleficPlay registers the attack-action once-per-turn trigger and emits the same-turn
// chain step. Each trigger fire creates one Runechant — the trigger handler authors a
// post-trigger log line so it groups beneath the triggering attack-action chain step. n
// is the printed counter count carried on the trigger so the handler can stay a top-level
// function (the Count field on the trigger drains as fires consume verse counters; the
// payload N field stays the printed value for log attribution if ever needed).
func maleficPlay(s *sim.TurnState, selfState *sim.CardState, selfCard sim.Card, n int) {
	s.AddAuraTrigger(sim.AuraTrigger{
		Self:        selfCard,
		Type:        sim.TriggerAttackAction,
		Count:       n,
		OncePerTurn: true,
		Handler:     maleficAuraHandler,
	})
	s.LogPlay(selfState)
}

// maleficAuraHandler is the once-per-turn attack-action trigger handler shared across
// Malefic Incantation variants. Reads t.Self for log attribution so the handler is a
// top-level function with no per-Play closure allocation.
func maleficAuraHandler(s *sim.TurnState, t *sim.AuraTrigger) int {
	created := s.CreateRunechants(1)
	// SkipLog discards both the post-trigger entry and its damage credit's text; avoid
	// the two DisplayName lookups + string concat when nothing will read them.
	if s.SkipLog {
		return created
	}
	return s.AddPostTriggerLogEntry(
		sim.DisplayName(t.Self)+" created a runechant",
		sim.DisplayName(s.TriggeringCard),
		created,
	)
}
