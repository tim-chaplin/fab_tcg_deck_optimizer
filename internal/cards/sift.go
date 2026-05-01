// Sift — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "Put up to 4 cards from your hand on the bottom of your deck, then draw that many cards.
// **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var siftTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SiftRed struct{}

func (SiftRed) ID() ids.CardID          { return ids.SiftRed }
func (SiftRed) Name() string            { return "Sift" }
func (SiftRed) Cost(*sim.TurnState) int { return 0 }
func (SiftRed) Pitch() int              { return 1 }
func (SiftRed) Attack() int             { return 0 }
func (SiftRed) Defense() int            { return 3 }
func (SiftRed) Types() card.TypeSet     { return siftTypes }
func (SiftRed) GoAgain() bool           { return true }

// not implemented: hand cycling
func (SiftRed) NotImplemented()                            {}
func (SiftRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type SiftYellow struct{}

func (SiftYellow) ID() ids.CardID          { return ids.SiftYellow }
func (SiftYellow) Name() string            { return "Sift" }
func (SiftYellow) Cost(*sim.TurnState) int { return 0 }
func (SiftYellow) Pitch() int              { return 2 }
func (SiftYellow) Attack() int             { return 0 }
func (SiftYellow) Defense() int            { return 3 }
func (SiftYellow) Types() card.TypeSet     { return siftTypes }
func (SiftYellow) GoAgain() bool           { return true }

// not implemented: hand cycling
func (SiftYellow) NotImplemented()                            {}
func (SiftYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type SiftBlue struct{}

func (SiftBlue) ID() ids.CardID          { return ids.SiftBlue }
func (SiftBlue) Name() string            { return "Sift" }
func (SiftBlue) Cost(*sim.TurnState) int { return 0 }
func (SiftBlue) Pitch() int              { return 3 }
func (SiftBlue) Attack() int             { return 0 }
func (SiftBlue) Defense() int            { return 3 }
func (SiftBlue) Types() card.TypeSet     { return siftTypes }
func (SiftBlue) GoAgain() bool           { return true }

// not implemented: hand cycling
func (SiftBlue) NotImplemented()                            {}
func (SiftBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
