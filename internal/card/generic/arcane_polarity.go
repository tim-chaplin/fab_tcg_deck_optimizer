// Arcane Polarity — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 1{h} If you've been dealt arcane damage this turn, instead gain 4{h}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var arcanePolarityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type ArcanePolarityRed struct{}

func (ArcanePolarityRed) ID() card.ID                               { return card.ArcanePolarityRed }
func (ArcanePolarityRed) Name() string                              { return "Arcane Polarity (Red)" }
func (ArcanePolarityRed) Cost(*card.TurnState) int                  { return 0 }
func (ArcanePolarityRed) Pitch() int                                { return 1 }
func (ArcanePolarityRed) Attack() int                               { return 0 }
func (ArcanePolarityRed) Defense() int                              { return 0 }
func (ArcanePolarityRed) Types() card.TypeSet                       { return arcanePolarityTypes }
func (ArcanePolarityRed) GoAgain() bool                             { return false }
// not implemented: 1{h} gain (4/3/2{h} if dealt arcane damage this turn)
func (ArcanePolarityRed) NotImplemented()                           {}
func (ArcanePolarityRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type ArcanePolarityYellow struct{}

func (ArcanePolarityYellow) ID() card.ID                               { return card.ArcanePolarityYellow }
func (ArcanePolarityYellow) Name() string                              { return "Arcane Polarity (Yellow)" }
func (ArcanePolarityYellow) Cost(*card.TurnState) int                  { return 0 }
func (ArcanePolarityYellow) Pitch() int                                { return 2 }
func (ArcanePolarityYellow) Attack() int                               { return 0 }
func (ArcanePolarityYellow) Defense() int                              { return 0 }
func (ArcanePolarityYellow) Types() card.TypeSet                       { return arcanePolarityTypes }
func (ArcanePolarityYellow) GoAgain() bool                             { return false }
// not implemented: 1{h} gain (4/3/2{h} if dealt arcane damage this turn)
func (ArcanePolarityYellow) NotImplemented()                           {}
func (ArcanePolarityYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type ArcanePolarityBlue struct{}

func (ArcanePolarityBlue) ID() card.ID                               { return card.ArcanePolarityBlue }
func (ArcanePolarityBlue) Name() string                              { return "Arcane Polarity (Blue)" }
func (ArcanePolarityBlue) Cost(*card.TurnState) int                  { return 0 }
func (ArcanePolarityBlue) Pitch() int                                { return 3 }
func (ArcanePolarityBlue) Attack() int                               { return 0 }
func (ArcanePolarityBlue) Defense() int                              { return 0 }
func (ArcanePolarityBlue) Types() card.TypeSet                       { return arcanePolarityTypes }
func (ArcanePolarityBlue) GoAgain() bool                             { return false }
// not implemented: 1{h} gain (4/3/2{h} if dealt arcane damage this turn)
func (ArcanePolarityBlue) NotImplemented()                           {}
func (ArcanePolarityBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
