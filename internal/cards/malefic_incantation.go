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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (c MaleficIncantationRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	maleficPlay(s, l, self, c, 3)
}

func (c MaleficIncantationYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	maleficPlay(s, l, self, c, 2)
}

func (c MaleficIncantationBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	maleficPlay(s, l, self, c, 1)
}

// maleficCreatedRunechantText is the precomputed rider line for each Malefic Incantation
// variant. Built once at init() so neither Play nor the per-fire handler does any string
// formatting on the hot anneal path.
var maleficCreatedRunechantText = func() map[ids.CardID]string {
	out := make(map[ids.CardID]string, 3)
	for _, c := range []card.Card{
		MaleficIncantationRed{},
		MaleficIncantationYellow{},
		MaleficIncantationBlue{},
	} {
		out[c.ID()] = c.DisplayName() + " created a runechant"
	}
	return out
}()

// maleficPlay registers the attack-action once-per-turn trigger and emits the same-turn
// chain step. Each trigger fire creates one Runechant — the trigger handler authors a
// post-trigger log line so it groups beneath the triggering attack-action chain step. n
// is the printed counter count carried on the trigger so the handler can stay a top-level
// function.
func maleficPlay(s card.GameEngine, l card.Logger, selfState *card.CardState, selfCard card.Card, n int) {
	s.AddAura(sim.Aura{
		Trigger:     sim.Trigger{TriggerType: sim.TriggerAttackAction, Handler: maleficAuraHandler},
		Self:        sim.CardOrTokenType{Card: selfCard},
		Count:       n,
		OncePerTurn: true,
	})
}

// maleficAuraHandler is the once-per-turn attack-action trigger handler shared across
// Malefic Incantation variants. Per-variant rider text is read off the table by
// aura.Self.CardID() so the hot fire path runs zero string allocations. Decrements
// aura.Count (the verse counter) and destroys the aura when the last verse fires.
func maleficAuraHandler(s card.GameEngine, l card.Logger, _ *sim.Trigger, a *sim.Aura) {
	cardID := a.Self.CardID()
	a.Count--
	lastVerse := a.Count <= 0
	s.CreateRunechants(1)
	l.AppendPostTrigger(s.TriggeringCard().DisplayName(), maleficCreatedRunechantText[cardID], 1)
	if lastVerse {
		s.DestroyAura(a, true)
	}
}
