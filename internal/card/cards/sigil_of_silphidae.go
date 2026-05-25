// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Play resolves the enter trigger directly via banishAuraFromGraveyard; OnLeavesArena runs
// the leave trigger.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (SigilOfSilphidaeBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	banishAuraFromGraveyard(ge, l, self.Card.DisplayName())
	installSigilSelfDestructAura(ge, self.Card)
}

// OnLeavesArena runs the "when this leaves the arena" clause: banish another aura from the
// graveyard and, if one was banished, deal 1 arcane damage. It runs before Sigil of
// Silphidae itself reaches the graveyard, so the "another aura" restriction holds.
func (c SigilOfSilphidaeBlue) OnLeavesArena(g card.GameEngine, l card.Logger) {
	banishAuraFromGraveyard(g, l, c.DisplayName())
}
