// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// Modelling: Play flips AuraCreated so same-turn aura-readers see it and registers a
// start-of-turn AuraTrigger with Count=1. Next turn the sim fires the trigger — the handler
// creates one live Runechant on the new turn's state — and graveyards Sigil of Deadwood as
// Count hits zero.
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
	s.AddAuraTrigger(card.AuraTrigger{
		Self:    c,
		Type:    card.TriggerStartOfTurn,
		Count:   1,
		Handler: func(s *card.TurnState) int { return s.CreateRunechants(1) },
	})
	return 0
}
