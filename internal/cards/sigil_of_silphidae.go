// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Play resolves the enter trigger directly via banishAuraFromGraveyard. The start-of-turn
// handler runs the leave trigger — because the sim graveyards Self only after Count hits zero,
// the scan naturally can't pick up Silphidae itself, satisfying "another aura" without an
// explicit skip.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfSilphidaeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfSilphidaeBlue struct{}

func (SigilOfSilphidaeBlue) ID() ids.CardID          { return ids.SigilOfSilphidaeBlue }
func (SigilOfSilphidaeBlue) Name() string            { return "Sigil of Silphidae" }
func (SigilOfSilphidaeBlue) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSilphidaeBlue) Pitch() int              { return 3 }
func (SigilOfSilphidaeBlue) Attack() int             { return 0 }
func (SigilOfSilphidaeBlue) Defense() int            { return 3 }
func (SigilOfSilphidaeBlue) Types() card.TypeSet     { return sigilOfSilphidaeTypes }
func (SigilOfSilphidaeBlue) GoAgain() bool           { return true }
func (SigilOfSilphidaeBlue) AddsFutureValue()        {}
func (c SigilOfSilphidaeBlue) Play(s *sim.TurnState, self *sim.CardState) {
	enterDamage := banishAuraFromGraveyard(s)
	s.RegisterStartOfTurn(c, 1, "Banished an aura, dealt 1 arcane damage", func(s *sim.TurnState) int {
		return banishAuraFromGraveyard(s)
	})
	s.ApplyAndLogEffectiveAttack(self)
	if enterDamage > 0 {
		s.ApplyAndLogRiderOnPlay(self, "Banished an aura, dealt 1 arcane damage", enterDamage)
	}
}
