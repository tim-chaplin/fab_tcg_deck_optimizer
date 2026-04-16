// Sigil of Deadwood — Runeblade Action - Aura. Cost 0, Pitch 3, Defense 2, Go again.
// Only printed in Blue.
// Text: "Go again. At the beginning of your action phase, destroy this. When this leaves the
// arena, create a Runechant token."
//
// Simplification: assume the aura resolves next turn and produces a Runechant. Play returns 1.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfDeadwoodTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)

type SigilOfDeadwoodBlue struct{}

func (SigilOfDeadwoodBlue) Name() string             { return "Sigil of Deadwood (Blue)" }
func (SigilOfDeadwoodBlue) Cost() int                { return 0 }
func (SigilOfDeadwoodBlue) Pitch() int               { return 3 }
func (SigilOfDeadwoodBlue) Attack() int              { return 0 }
func (SigilOfDeadwoodBlue) Defense() int             { return 2 }
func (SigilOfDeadwoodBlue) Types() card.TypeSet      { return sigilOfDeadwoodTypes }
func (SigilOfDeadwoodBlue) GoAgain() bool            { return true }
func (SigilOfDeadwoodBlue) Play(*card.TurnState) int { return 1 }
