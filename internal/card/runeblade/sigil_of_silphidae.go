// Sigil of Silphidae — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 3, Arcane 1. Go again.
// Only printed in Blue.
// Text: "When this enters or leaves the arena, you may banish another aura from your graveyard.
// If you do, deal 1 arcane damage to target hero. At the beginning of your action phase, destroy
// this."
//
// Simplifications: assume we always have an aura in the graveyard to banish on both the enter
// and leave triggers, so Sigil of Silphidae is worth 2 damage when played. Cross-turn
// aura-persistence isn't modelled; we collapse the effect to an immediate 2 value on play.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSilphidaeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfSilphidaeBlue struct{}

func (SigilOfSilphidaeBlue) ID() card.ID                 { return card.SigilOfSilphidaeBlue }
func (SigilOfSilphidaeBlue) Name() string           { return "Sigil of Silphidae (Blue)" }
func (SigilOfSilphidaeBlue) Cost(*card.TurnState) int              { return 0 }
func (SigilOfSilphidaeBlue) Pitch() int             { return 3 }
func (SigilOfSilphidaeBlue) Attack() int            { return 0 }
func (SigilOfSilphidaeBlue) Defense() int           { return 3 }
func (SigilOfSilphidaeBlue) Types() card.TypeSet    { return sigilOfSilphidaeTypes }
func (SigilOfSilphidaeBlue) GoAgain() bool          { return true }
func (SigilOfSilphidaeBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AuraCreated = true
	s.ArcaneDamageDealt = true // the aura-banish riders deal 1 arcane each (enter + leave)
	return 2
}
