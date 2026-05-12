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
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (MaleficIncantationRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateOncePerTurnAttackActionAura(self, maleficAuraHandler, 3)
}

func (MaleficIncantationYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateOncePerTurnAttackActionAura(self, maleficAuraHandler, 2)
}

func (MaleficIncantationBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	g.CreateOncePerTurnAttackActionAura(self, maleficAuraHandler, 1)
}

// maleficCreatedRunechantText is the precomputed rider line for each Malefic Incantation
// variant.
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

// maleficAuraHandler is the once-per-turn attack-action trigger handler shared across
// Malefic Incantation variants. Decrements the verse counter and destroys the aura when
// the last verse fires.
func maleficAuraHandler(g card.GameEngine, l card.Logger, a card.Aura) {
	cardID := a.CardID()
	lastVerse := a.DecrementCount() <= 0
	g.CreateRunechants(1)
	l.AppendPostTrigger(g.TriggeringCard().DisplayName(), maleficCreatedRunechantText[cardID], 1)
	if lastVerse {
		a.Destroy(true)
	}
}
