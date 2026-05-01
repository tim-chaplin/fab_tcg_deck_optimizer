// Even Bigger Than That! — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue
// 3.
//
// Text: "Play Even Bigger Than That! only if you've dealt {p} this turn. **Opt 3**, then reveal the
// top card of your deck. If it has {p} greater than the amount of damage you've dealt this turn,
// create a Quicken token and draw a card."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var evenBiggerThanThatTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type EvenBiggerThanThatRed struct{}

func (EvenBiggerThanThatRed) ID() ids.CardID          { return ids.EvenBiggerThanThatRed }
func (EvenBiggerThanThatRed) Name() string            { return "Even Bigger Than That!" }
func (EvenBiggerThanThatRed) Cost(*sim.TurnState) int { return 0 }
func (EvenBiggerThanThatRed) Pitch() int              { return 1 }
func (EvenBiggerThanThatRed) Attack() int             { return 0 }
func (EvenBiggerThanThatRed) Defense() int            { return 0 }
func (EvenBiggerThanThatRed) Types() card.TypeSet     { return evenBiggerThanThatTypes }
func (EvenBiggerThanThatRed) GoAgain() bool           { return false }

// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn
func (EvenBiggerThanThatRed) NotImplemented()                            {}
func (EvenBiggerThanThatRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type EvenBiggerThanThatYellow struct{}

func (EvenBiggerThanThatYellow) ID() ids.CardID          { return ids.EvenBiggerThanThatYellow }
func (EvenBiggerThanThatYellow) Name() string            { return "Even Bigger Than That!" }
func (EvenBiggerThanThatYellow) Cost(*sim.TurnState) int { return 0 }
func (EvenBiggerThanThatYellow) Pitch() int              { return 2 }
func (EvenBiggerThanThatYellow) Attack() int             { return 0 }
func (EvenBiggerThanThatYellow) Defense() int            { return 0 }
func (EvenBiggerThanThatYellow) Types() card.TypeSet     { return evenBiggerThanThatTypes }
func (EvenBiggerThanThatYellow) GoAgain() bool           { return false }

// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn
func (EvenBiggerThanThatYellow) NotImplemented()                            {}
func (EvenBiggerThanThatYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type EvenBiggerThanThatBlue struct{}

func (EvenBiggerThanThatBlue) ID() ids.CardID          { return ids.EvenBiggerThanThatBlue }
func (EvenBiggerThanThatBlue) Name() string            { return "Even Bigger Than That!" }
func (EvenBiggerThanThatBlue) Cost(*sim.TurnState) int { return 0 }
func (EvenBiggerThanThatBlue) Pitch() int              { return 3 }
func (EvenBiggerThanThatBlue) Attack() int             { return 0 }
func (EvenBiggerThanThatBlue) Defense() int            { return 0 }
func (EvenBiggerThanThatBlue) Types() card.TypeSet     { return evenBiggerThanThatTypes }
func (EvenBiggerThanThatBlue) GoAgain() bool           { return false }

// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn
func (EvenBiggerThanThatBlue) NotImplemented()                            {}
func (EvenBiggerThanThatBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
