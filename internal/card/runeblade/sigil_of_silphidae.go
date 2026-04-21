// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Enter trigger fires on Play: banishAuraFromGraveyard scans s.Graveyard for an aura and
// credits 1 arcane if one lands in s.Banish. At the start of next turn PlayNextTurn scans
// the graveyard for the leave trigger FIRST, then adds the sigil to the graveyard — the
// ordering honours the printed "another aura" restriction without any explicit skip.
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
	return banishAuraFromGraveyard(s)
}
func (c SigilOfSilphidaeBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	// Scan BEFORE Silphidae lands in the graveyard so the printed "another aura" restriction
	// is satisfied naturally — the scan can't pick up the sigil itself.
	r := card.DelayedPlayResult{Damage: banishAuraFromGraveyard(s)}
	s.AddToGraveyard(c)
	return r
}
