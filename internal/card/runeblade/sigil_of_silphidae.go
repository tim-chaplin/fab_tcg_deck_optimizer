// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Arcane 1. Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Play credits the enter trigger; PlayNextTurn fires when the aura leaves at the start of the
// next action phase and credits the leave trigger. Both triggers currently still assume an aura
// is available to banish — TODO: scan the graveyard so the clauses fizzle correctly.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSilphidaeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfSilphidaeBlue struct{}

func (SigilOfSilphidaeBlue) ID() card.ID               { return card.SigilOfSilphidaeBlue }
func (SigilOfSilphidaeBlue) Name() string              { return "Sigil of Silphidae (Blue)" }
func (SigilOfSilphidaeBlue) Cost(*card.TurnState) int  { return 0 }
func (SigilOfSilphidaeBlue) Pitch() int                { return 3 }
func (SigilOfSilphidaeBlue) Attack() int               { return 0 }
func (SigilOfSilphidaeBlue) Defense() int              { return 3 }
func (SigilOfSilphidaeBlue) Types() card.TypeSet       { return sigilOfSilphidaeTypes }
func (SigilOfSilphidaeBlue) GoAgain() bool             { return true }
func (SigilOfSilphidaeBlue) Play(s *card.TurnState) int {
	s.AuraCreated = true
	s.ArcaneDamageDealt = true // enter-trigger's banish-an-aura-for-1-arcane
	return 1
}

// PlayNextTurn destroys the aura at the start of the next action phase and credits the leave
// trigger's 1 arcane. TODO: gate the 1 arcane on an actual aura being present in the graveyard.
func (SigilOfSilphidaeBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.DestroyThis()
	return card.DelayedPlayResult{Damage: 1}
}
