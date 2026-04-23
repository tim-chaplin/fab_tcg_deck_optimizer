// Sift — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "Put up to 4 cards from your hand on the bottom of your deck, then draw that many cards.
// **Go again**"
//
// Simplification: Hand cycling isn't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var siftTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type SiftRed struct{}

func (SiftRed) ID() card.ID                 { return card.SiftRed }
func (SiftRed) Name() string                { return "Sift (Red)" }
func (SiftRed) Cost(*card.TurnState) int                   { return 0 }
func (SiftRed) Pitch() int                  { return 1 }
func (SiftRed) Attack() int                 { return 0 }
func (SiftRed) Defense() int                { return 3 }
func (SiftRed) Types() card.TypeSet         { return siftTypes }
func (SiftRed) GoAgain() bool               { return true }
func (SiftRed) Play(s *card.TurnState, _ *card.CardState) int { return 0 }

type SiftYellow struct{}

func (SiftYellow) ID() card.ID                 { return card.SiftYellow }
func (SiftYellow) Name() string                { return "Sift (Yellow)" }
func (SiftYellow) Cost(*card.TurnState) int                   { return 0 }
func (SiftYellow) Pitch() int                  { return 2 }
func (SiftYellow) Attack() int                 { return 0 }
func (SiftYellow) Defense() int                { return 3 }
func (SiftYellow) Types() card.TypeSet         { return siftTypes }
func (SiftYellow) GoAgain() bool               { return true }
func (SiftYellow) Play(s *card.TurnState, _ *card.CardState) int { return 0 }

type SiftBlue struct{}

func (SiftBlue) ID() card.ID                 { return card.SiftBlue }
func (SiftBlue) Name() string                { return "Sift (Blue)" }
func (SiftBlue) Cost(*card.TurnState) int                   { return 0 }
func (SiftBlue) Pitch() int                  { return 3 }
func (SiftBlue) Attack() int                 { return 0 }
func (SiftBlue) Defense() int                { return 3 }
func (SiftBlue) Types() card.TypeSet         { return siftTypes }
func (SiftBlue) GoAgain() bool               { return true }
func (SiftBlue) Play(s *card.TurnState, _ *card.CardState) int { return 0 }
