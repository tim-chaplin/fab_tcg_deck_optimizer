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

func (SigilOfSilphidaeBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	enterDamage := banishAuraFromGraveyard(s)
	s.AddStartOfTurnAura(self, sigilOfSilphidaeAuraHandler, 1)
	if enterDamage > 0 {
		s.AddValue(enterDamage)
		l.AppendPostTrigger(self.Card.DisplayName(), "Banished an aura, dealt 1 arcane damage", enterDamage)
	}
}

// sigilOfSilphidaeAuraHandler runs the leave trigger on the next turn: scans the graveyard
// for an aura to banish, credits 1 arcane damage on a hit, then destroys the aura.
func sigilOfSilphidaeAuraHandler(s card.GameEngine, l card.Logger, a card.Aura) {
	n := banishAuraFromGraveyard(s)
	if n > 0 {
		s.AddValue(n)
		l.AppendPostTrigger(a.SelfName(), "Banished an aura, dealt 1 arcane damage", n)
	}
	a.Destroy(true)
}
