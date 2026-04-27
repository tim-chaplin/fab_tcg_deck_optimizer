// Count Your Blessings — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain X{h}, where X is 3 plus the number of Count Your Blessings in your graveyard."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var countYourBlessingsTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type CountYourBlessingsRed struct{}

func (CountYourBlessingsRed) ID() card.ID              { return card.CountYourBlessingsRed }
func (CountYourBlessingsRed) Name() string             { return "Count Your Blessings" }
func (CountYourBlessingsRed) Cost(*card.TurnState) int { return 2 }
func (CountYourBlessingsRed) Pitch() int               { return 1 }
func (CountYourBlessingsRed) Attack() int              { return 0 }
func (CountYourBlessingsRed) Defense() int             { return 0 }
func (CountYourBlessingsRed) Types() card.TypeSet      { return countYourBlessingsTypes }
func (CountYourBlessingsRed) GoAgain() bool            { return false }
func (CountYourBlessingsRed) NotSilverAgeLegal()       {}

// not implemented: graveyard-scaled X{h} gain (also banlisted)
func (CountYourBlessingsRed) NotImplemented()                              {}
func (CountYourBlessingsRed) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type CountYourBlessingsYellow struct{}

func (CountYourBlessingsYellow) ID() card.ID              { return card.CountYourBlessingsYellow }
func (CountYourBlessingsYellow) Name() string             { return "Count Your Blessings" }
func (CountYourBlessingsYellow) Cost(*card.TurnState) int { return 2 }
func (CountYourBlessingsYellow) Pitch() int               { return 2 }
func (CountYourBlessingsYellow) Attack() int              { return 0 }
func (CountYourBlessingsYellow) Defense() int             { return 0 }
func (CountYourBlessingsYellow) Types() card.TypeSet      { return countYourBlessingsTypes }
func (CountYourBlessingsYellow) GoAgain() bool            { return false }
func (CountYourBlessingsYellow) NotSilverAgeLegal()       {}

// not implemented: graveyard-scaled X{h} gain (also banlisted)
func (CountYourBlessingsYellow) NotImplemented()                              {}
func (CountYourBlessingsYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type CountYourBlessingsBlue struct{}

func (CountYourBlessingsBlue) ID() card.ID              { return card.CountYourBlessingsBlue }
func (CountYourBlessingsBlue) Name() string             { return "Count Your Blessings" }
func (CountYourBlessingsBlue) Cost(*card.TurnState) int { return 2 }
func (CountYourBlessingsBlue) Pitch() int               { return 3 }
func (CountYourBlessingsBlue) Attack() int              { return 0 }
func (CountYourBlessingsBlue) Defense() int             { return 0 }
func (CountYourBlessingsBlue) Types() card.TypeSet      { return countYourBlessingsTypes }
func (CountYourBlessingsBlue) GoAgain() bool            { return false }
func (CountYourBlessingsBlue) NotSilverAgeLegal()       {}

// not implemented: graveyard-scaled X{h} gain (also banlisted)
func (CountYourBlessingsBlue) NotImplemented()                              {}
func (CountYourBlessingsBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
