// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// The rune fires at next turn's upkeep, so this card implements card.DelayedPlay —
// PlayNextTurn destroys the aura and creates the Runechant on next turn's starting state.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfDeadwoodTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfDeadwoodBlue struct{}

func (SigilOfDeadwoodBlue) ID() card.ID              { return card.SigilOfDeadwoodBlue }
func (SigilOfDeadwoodBlue) Name() string             { return "Sigil of Deadwood (Blue)" }
func (SigilOfDeadwoodBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfDeadwoodBlue) Pitch() int               { return 3 }
func (SigilOfDeadwoodBlue) Attack() int              { return 0 }
func (SigilOfDeadwoodBlue) Defense() int             { return 2 }
func (SigilOfDeadwoodBlue) Types() card.TypeSet      { return sigilOfDeadwoodTypes }
func (SigilOfDeadwoodBlue) GoAgain() bool            { return true }
func (SigilOfDeadwoodBlue) AddsFutureValue()         {}
func (c SigilOfDeadwoodBlue) Play(s *card.TurnState, _ *card.CardState) int {
	s.AuraCreated = true
	return 0
}
func (c SigilOfDeadwoodBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.AddToGraveyard(c)
	return card.DelayedPlayResult{Damage: s.CreateRunechants(1)}
}
