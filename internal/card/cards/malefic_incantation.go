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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func (MaleficIncantationRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self.Card, triggertype.CardOrAbility, maleficAuraHandler, 3, true, card.TypeSet.IsAttackAction)
}

func (MaleficIncantationYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self.Card, triggertype.CardOrAbility, maleficAuraHandler, 2, true, card.TypeSet.IsAttackAction)
}

func (MaleficIncantationBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateAura(self.Card, triggertype.CardOrAbility, maleficAuraHandler, 1, true, card.TypeSet.IsAttackAction)
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
func maleficAuraHandler(ge card.GameEngine, l card.Logger, a card.Aura, triggeringCard *card.CardState, _ triggertype.Type) {
	cardID := a.CardID()
	lastVerse := a.DecrementCount() <= 0
	ge.CreateRunechants(1)
	l.AppendPostTrigger(triggeringCard.Card.DisplayName(), maleficCreatedRunechantText[cardID], 1)
	if lastVerse {
		a.Destroy(true)
	}
}
