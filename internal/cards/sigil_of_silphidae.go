// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Play resolves the enter trigger directly via banishAuraFromGraveyard. The start-of-turn
// handler runs the leave trigger.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (SigilOfSilphidaeBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	banishAuraFromGraveyard(g, l, self.Card.DisplayName())
	g.CreateStartOfTurnAura(self, sigilOfSilphidaeAuraHandler, 1)
}

// sigilOfSilphidaeAuraHandler runs the leave trigger on the next turn: scans the graveyard
// for an aura to banish and deals 1 arcane damage on success, then destroys the aura.
func sigilOfSilphidaeAuraHandler(g card.GameEngine, l card.Logger, a card.Aura) {
	banishAuraFromGraveyard(g, l, a.CardName())
	a.Destroy(true)
}
