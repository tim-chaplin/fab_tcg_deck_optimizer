// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// Play is a no-op beyond flipping AuraCreated; PlayNextTurn fires when the aura is destroyed at
// the start of the next action phase and creates the Runechant token then.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfDeadwoodTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfDeadwoodBlue struct{}

func (SigilOfDeadwoodBlue) ID() card.ID                { return card.SigilOfDeadwoodBlue }
func (SigilOfDeadwoodBlue) Name() string               { return "Sigil of Deadwood (Blue)" }
func (SigilOfDeadwoodBlue) Cost(*card.TurnState) int   { return 0 }
func (SigilOfDeadwoodBlue) Pitch() int                 { return 3 }
func (SigilOfDeadwoodBlue) Attack() int                { return 0 }
func (SigilOfDeadwoodBlue) Defense() int               { return 2 }
func (SigilOfDeadwoodBlue) Types() card.TypeSet        { return sigilOfDeadwoodTypes }
func (SigilOfDeadwoodBlue) GoAgain() bool              { return true }
func (SigilOfDeadwoodBlue) Play(s *card.TurnState) int {
	s.AuraCreated = true
	return 0
}

// PlayNextTurn creates the Runechant token that fires when the aura leaves the arena at the
// start of the next action phase.
func (SigilOfDeadwoodBlue) PlayNextTurn(s *card.TurnState) card.DelayedPlayResult {
	s.DestroyThis()
	return card.DelayedPlayResult{Damage: s.CreateRunechant()}
}
