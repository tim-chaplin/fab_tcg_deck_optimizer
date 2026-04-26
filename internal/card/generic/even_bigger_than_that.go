// Even Bigger Than That! — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue
// 3.
//
// Text: "Play Even Bigger Than That! only if you've dealt {p} this turn. **Opt 3**, then reveal the
// top card of your deck. If it has {p} greater than the amount of damage you've dealt this turn,
// create a Quicken token and draw a card."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var evenBiggerThanThatTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type EvenBiggerThanThatRed struct{}

func (EvenBiggerThanThatRed) ID() card.ID                               { return card.EvenBiggerThanThatRed }
func (EvenBiggerThanThatRed) Name() string                              { return "Even Bigger Than That! (Red)" }
func (EvenBiggerThanThatRed) Cost(*card.TurnState) int                  { return 0 }
func (EvenBiggerThanThatRed) Pitch() int                                { return 1 }
func (EvenBiggerThanThatRed) Attack() int                               { return 0 }
func (EvenBiggerThanThatRed) Defense() int                              { return 0 }
func (EvenBiggerThanThatRed) Types() card.TypeSet                       { return evenBiggerThanThatTypes }
func (EvenBiggerThanThatRed) GoAgain() bool                             { return false }
// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn
func (EvenBiggerThanThatRed) NotImplemented()                           {}
func (EvenBiggerThanThatRed) Play(*card.TurnState, *card.CardState) int { return 0 }

type EvenBiggerThanThatYellow struct{}

func (EvenBiggerThanThatYellow) ID() card.ID                               { return card.EvenBiggerThanThatYellow }
func (EvenBiggerThanThatYellow) Name() string                              { return "Even Bigger Than That! (Yellow)" }
func (EvenBiggerThanThatYellow) Cost(*card.TurnState) int                  { return 0 }
func (EvenBiggerThanThatYellow) Pitch() int                                { return 2 }
func (EvenBiggerThanThatYellow) Attack() int                               { return 0 }
func (EvenBiggerThanThatYellow) Defense() int                              { return 0 }
func (EvenBiggerThanThatYellow) Types() card.TypeSet                       { return evenBiggerThanThatTypes }
func (EvenBiggerThanThatYellow) GoAgain() bool                             { return false }
// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn
func (EvenBiggerThanThatYellow) NotImplemented()                           {}
func (EvenBiggerThanThatYellow) Play(*card.TurnState, *card.CardState) int { return 0 }

type EvenBiggerThanThatBlue struct{}

func (EvenBiggerThanThatBlue) ID() card.ID                               { return card.EvenBiggerThanThatBlue }
func (EvenBiggerThanThatBlue) Name() string                              { return "Even Bigger Than That! (Blue)" }
func (EvenBiggerThanThatBlue) Cost(*card.TurnState) int                  { return 0 }
func (EvenBiggerThanThatBlue) Pitch() int                                { return 3 }
func (EvenBiggerThanThatBlue) Attack() int                               { return 0 }
func (EvenBiggerThanThatBlue) Defense() int                              { return 0 }
func (EvenBiggerThanThatBlue) Types() card.TypeSet                       { return evenBiggerThanThatTypes }
func (EvenBiggerThanThatBlue) GoAgain() bool                             { return false }
// not implemented: Opt + reveal-and-Quicken trigger; gated on damage dealt this turn
func (EvenBiggerThanThatBlue) NotImplemented()                           {}
func (EvenBiggerThanThatBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
