// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Enter trigger fires on Play: banishAuraFromGraveyard scans s.Graveyard for an aura and
// credits 1 arcane if one lands in s.Banish. At the start of next turn PlayNextTurn destroys
// this (AddToGraveyard) and runs the leave trigger — another banishAuraFromGraveyard pass,
// skipping the sigil itself so the "another aura" restriction holds.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSilphidaeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfSilphidaeBlue struct{}

func (SigilOfSilphidaeBlue) ID() card.ID              { return card.SigilOfSilphidaeBlue }
func (SigilOfSilphidaeBlue) Name() string             { return "Sigil of Silphidae (Blue)" }
func (SigilOfSilphidaeBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfSilphidaeBlue) Pitch() int               { return 3 }
func (SigilOfSilphidaeBlue) Attack() int              { return 0 }
func (SigilOfSilphidaeBlue) Defense() int             { return 3 }
func (SigilOfSilphidaeBlue) Types() card.TypeSet      { return sigilOfSilphidaeTypes }
func (SigilOfSilphidaeBlue) GoAgain() bool            { return true }
func (SigilOfSilphidaeBlue) NoMemo()                  {}
func (c SigilOfSilphidaeBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AuraCreated = true
	return banishAuraFromGraveyard(s, c)
}
func (c SigilOfSilphidaeBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{Damage: banishAuraFromGraveyard(s, c)}
}
