// Drowning Dire — Runeblade Action - Attack. Cost 2, Defense 3. Has Dominate.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 5, Yellow 4, Blue 3.
// Text: "If you have played or created an aura this turn, Drowning Dire gains **dominate**."
//
// Simplification: Dominate (opposing hero blocks with at most 1 card) isn't modelled — the
// optimizer doesn't simulate defender blocks, so Dominate currently adds no value. Damage returned
// is the printed attack. AuraCreated is not set.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var drowningDireTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type DrowningDireRed struct{}

func (DrowningDireRed) ID() card.ID                 { return card.DrowningDireRed }
func (DrowningDireRed) Name() string                 { return "Drowning Dire (Red)" }
func (DrowningDireRed) Cost(*card.TurnState) int                    { return 2 }
func (DrowningDireRed) Pitch() int                   { return 1 }
func (DrowningDireRed) Attack() int                  { return 5 }
func (DrowningDireRed) Defense() int                 { return 3 }
func (DrowningDireRed) Types() card.TypeSet       { return drowningDireTypes }
func (DrowningDireRed) GoAgain() bool                { return false }
func (c DrowningDireRed) Play(*card.TurnState, *card.CardState) int   { return c.Attack() }

type DrowningDireYellow struct{}

func (DrowningDireYellow) ID() card.ID                 { return card.DrowningDireYellow }
func (DrowningDireYellow) Name() string                 { return "Drowning Dire (Yellow)" }
func (DrowningDireYellow) Cost(*card.TurnState) int                    { return 2 }
func (DrowningDireYellow) Pitch() int                   { return 2 }
func (DrowningDireYellow) Attack() int                  { return 4 }
func (DrowningDireYellow) Defense() int                 { return 3 }
func (DrowningDireYellow) Types() card.TypeSet       { return drowningDireTypes }
func (DrowningDireYellow) GoAgain() bool                { return false }
func (c DrowningDireYellow) Play(*card.TurnState, *card.CardState) int   { return c.Attack() }

type DrowningDireBlue struct{}

func (DrowningDireBlue) ID() card.ID                 { return card.DrowningDireBlue }
func (DrowningDireBlue) Name() string                 { return "Drowning Dire (Blue)" }
func (DrowningDireBlue) Cost(*card.TurnState) int                    { return 2 }
func (DrowningDireBlue) Pitch() int                   { return 3 }
func (DrowningDireBlue) Attack() int                  { return 3 }
func (DrowningDireBlue) Defense() int                 { return 3 }
func (DrowningDireBlue) Types() card.TypeSet       { return drowningDireTypes }
func (DrowningDireBlue) GoAgain() bool                { return false }
func (c DrowningDireBlue) Play(*card.TurnState, *card.CardState) int   { return c.Attack() }
