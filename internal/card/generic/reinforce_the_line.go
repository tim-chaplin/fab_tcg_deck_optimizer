// Reinforce the Line — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Target defending attack action card gains +4{d}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var reinforceTheLineTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type ReinforceTheLineRed struct{}

func (ReinforceTheLineRed) ID() card.ID                               { return card.ReinforceTheLineRed }
func (ReinforceTheLineRed) Name() string                              { return "Reinforce the Line (Red)" }
func (ReinforceTheLineRed) Cost(*card.TurnState) int                  { return 0 }
func (ReinforceTheLineRed) Pitch() int                                { return 1 }
func (ReinforceTheLineRed) Attack() int                               { return 0 }
func (ReinforceTheLineRed) Defense() int                              { return 0 }
func (ReinforceTheLineRed) Types() card.TypeSet                       { return reinforceTheLineTypes }
func (ReinforceTheLineRed) GoAgain() bool                             { return false }
// not implemented: Instant +1{d} grant to a defending card
func (ReinforceTheLineRed) NotImplemented()                           {}
func (ReinforceTheLineRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type ReinforceTheLineYellow struct{}

func (ReinforceTheLineYellow) ID() card.ID                               { return card.ReinforceTheLineYellow }
func (ReinforceTheLineYellow) Name() string                              { return "Reinforce the Line (Yellow)" }
func (ReinforceTheLineYellow) Cost(*card.TurnState) int                  { return 0 }
func (ReinforceTheLineYellow) Pitch() int                                { return 2 }
func (ReinforceTheLineYellow) Attack() int                               { return 0 }
func (ReinforceTheLineYellow) Defense() int                              { return 0 }
func (ReinforceTheLineYellow) Types() card.TypeSet                       { return reinforceTheLineTypes }
func (ReinforceTheLineYellow) GoAgain() bool                             { return false }
// not implemented: Instant +1{d} grant to a defending card
func (ReinforceTheLineYellow) NotImplemented()                           {}
func (ReinforceTheLineYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type ReinforceTheLineBlue struct{}

func (ReinforceTheLineBlue) ID() card.ID                               { return card.ReinforceTheLineBlue }
func (ReinforceTheLineBlue) Name() string                              { return "Reinforce the Line (Blue)" }
func (ReinforceTheLineBlue) Cost(*card.TurnState) int                  { return 0 }
func (ReinforceTheLineBlue) Pitch() int                                { return 3 }
func (ReinforceTheLineBlue) Attack() int                               { return 0 }
func (ReinforceTheLineBlue) Defense() int                              { return 0 }
func (ReinforceTheLineBlue) Types() card.TypeSet                       { return reinforceTheLineTypes }
func (ReinforceTheLineBlue) GoAgain() bool                             { return false }
// not implemented: Instant +1{d} grant to a defending card
func (ReinforceTheLineBlue) NotImplemented()                           {}
func (ReinforceTheLineBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
