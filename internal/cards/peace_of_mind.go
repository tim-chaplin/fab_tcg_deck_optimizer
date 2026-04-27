// Peace of Mind — Generic Instant. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "The next time you would be dealt {p} damage, prevent 4 of that damage. Create a Ponder
// token."

package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var peaceOfMindTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type PeaceOfMindRed struct{}

func (PeaceOfMindRed) ID() card.ID              { return card.PeaceOfMindRed }
func (PeaceOfMindRed) Name() string             { return "Peace of Mind" }
func (PeaceOfMindRed) Cost(*card.TurnState) int { return 2 }
func (PeaceOfMindRed) Pitch() int               { return 1 }
func (PeaceOfMindRed) Attack() int              { return 0 }
func (PeaceOfMindRed) Defense() int             { return 0 }
func (PeaceOfMindRed) Types() card.TypeSet      { return peaceOfMindTypes }
func (PeaceOfMindRed) GoAgain() bool            { return false }

// not implemented: Instant 'prevent 4 of next {p}-damage hit'; creates a Ponder token
func (PeaceOfMindRed) NotImplemented()                              {}
func (PeaceOfMindRed) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type PeaceOfMindYellow struct{}

func (PeaceOfMindYellow) ID() card.ID              { return card.PeaceOfMindYellow }
func (PeaceOfMindYellow) Name() string             { return "Peace of Mind" }
func (PeaceOfMindYellow) Cost(*card.TurnState) int { return 2 }
func (PeaceOfMindYellow) Pitch() int               { return 2 }
func (PeaceOfMindYellow) Attack() int              { return 0 }
func (PeaceOfMindYellow) Defense() int             { return 0 }
func (PeaceOfMindYellow) Types() card.TypeSet      { return peaceOfMindTypes }
func (PeaceOfMindYellow) GoAgain() bool            { return false }

// not implemented: Instant 'prevent 4 of next {p}-damage hit'; creates a Ponder token
func (PeaceOfMindYellow) NotImplemented()                              {}
func (PeaceOfMindYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type PeaceOfMindBlue struct{}

func (PeaceOfMindBlue) ID() card.ID              { return card.PeaceOfMindBlue }
func (PeaceOfMindBlue) Name() string             { return "Peace of Mind" }
func (PeaceOfMindBlue) Cost(*card.TurnState) int { return 2 }
func (PeaceOfMindBlue) Pitch() int               { return 3 }
func (PeaceOfMindBlue) Attack() int              { return 0 }
func (PeaceOfMindBlue) Defense() int             { return 0 }
func (PeaceOfMindBlue) Types() card.TypeSet      { return peaceOfMindTypes }
func (PeaceOfMindBlue) GoAgain() bool            { return false }

// not implemented: Instant 'prevent 4 of next {p}-damage hit'; creates a Ponder token
func (PeaceOfMindBlue) NotImplemented()                              {}
func (PeaceOfMindBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
