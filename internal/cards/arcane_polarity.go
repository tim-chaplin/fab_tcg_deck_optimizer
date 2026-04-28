// Arcane Polarity — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 1{h} If you've been dealt arcane damage this turn, instead gain 4{h}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var arcanePolarityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type ArcanePolarityRed struct{}

func (ArcanePolarityRed) ID() ids.CardID          { return ids.ArcanePolarityRed }
func (ArcanePolarityRed) Name() string            { return "Arcane Polarity" }
func (ArcanePolarityRed) Cost(*sim.TurnState) int { return 0 }
func (ArcanePolarityRed) Pitch() int              { return 1 }
func (ArcanePolarityRed) Attack() int             { return 0 }
func (ArcanePolarityRed) Defense() int            { return 0 }
func (ArcanePolarityRed) Types() card.TypeSet     { return arcanePolarityTypes }
func (ArcanePolarityRed) GoAgain() bool           { return false }

// not implemented: 1{h} gain (4/3/2{h} if dealt arcane damage this turn)
func (ArcanePolarityRed) NotImplemented()                            {}
func (ArcanePolarityRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type ArcanePolarityYellow struct{}

func (ArcanePolarityYellow) ID() ids.CardID          { return ids.ArcanePolarityYellow }
func (ArcanePolarityYellow) Name() string            { return "Arcane Polarity" }
func (ArcanePolarityYellow) Cost(*sim.TurnState) int { return 0 }
func (ArcanePolarityYellow) Pitch() int              { return 2 }
func (ArcanePolarityYellow) Attack() int             { return 0 }
func (ArcanePolarityYellow) Defense() int            { return 0 }
func (ArcanePolarityYellow) Types() card.TypeSet     { return arcanePolarityTypes }
func (ArcanePolarityYellow) GoAgain() bool           { return false }

// not implemented: 1{h} gain (4/3/2{h} if dealt arcane damage this turn)
func (ArcanePolarityYellow) NotImplemented()                            {}
func (ArcanePolarityYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type ArcanePolarityBlue struct{}

func (ArcanePolarityBlue) ID() ids.CardID          { return ids.ArcanePolarityBlue }
func (ArcanePolarityBlue) Name() string            { return "Arcane Polarity" }
func (ArcanePolarityBlue) Cost(*sim.TurnState) int { return 0 }
func (ArcanePolarityBlue) Pitch() int              { return 3 }
func (ArcanePolarityBlue) Attack() int             { return 0 }
func (ArcanePolarityBlue) Defense() int            { return 0 }
func (ArcanePolarityBlue) Types() card.TypeSet     { return arcanePolarityTypes }
func (ArcanePolarityBlue) GoAgain() bool           { return false }

// not implemented: 1{h} gain (4/3/2{h} if dealt arcane damage this turn)
func (ArcanePolarityBlue) NotImplemented()                            {}
func (ArcanePolarityBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
