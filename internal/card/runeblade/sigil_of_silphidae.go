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

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSilphidaeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfSilphidaeBlue struct{}

func (SigilOfSilphidaeBlue) ID() card.ID              { return card.SigilOfSilphidaeBlue }
func (SigilOfSilphidaeBlue) Name() string             { return "Sigil of Silphidae" }
func (SigilOfSilphidaeBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfSilphidaeBlue) Pitch() int               { return 3 }
func (SigilOfSilphidaeBlue) Attack() int              { return 0 }
func (SigilOfSilphidaeBlue) Defense() int             { return 3 }
func (SigilOfSilphidaeBlue) Types() card.TypeSet      { return sigilOfSilphidaeTypes }
func (SigilOfSilphidaeBlue) GoAgain() bool            { return true }
func (SigilOfSilphidaeBlue) AddsFutureValue()         {}
func (c SigilOfSilphidaeBlue) Play(s *card.TurnState, self *card.CardState) {
	enterDamage := banishAuraFromGraveyard(s)
	s.AddAuraTrigger(card.AuraTrigger{
		Self:  c,
		Type:  card.TriggerStartOfTurn,
		Count: 1,
		Handler: func(s *card.TurnState) int {
			return banishAuraFromGraveyard(s)
		},
	})
	s.ApplyAndLogEffectiveAttack(self)
	if enterDamage > 0 {
		s.LogRiderOnPlay(self, "Banished an aura, dealt 1 arcane damage", enterDamage)
	}
}
