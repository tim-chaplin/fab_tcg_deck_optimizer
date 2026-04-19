// Visit the Blacksmith — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next sword attack this turn gains +1{p}. **Go again**"
//
// Simplification: Next-sword-attack bonuses aren't applied (weapon chain isn't peeked).
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var visitTheBlacksmithTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type VisitTheBlacksmithBlue struct{}

func (VisitTheBlacksmithBlue) ID() card.ID                 { return card.VisitTheBlacksmithBlue }
func (VisitTheBlacksmithBlue) Name() string                { return "Visit the Blacksmith (Blue)" }
func (VisitTheBlacksmithBlue) Cost(*card.TurnState) int                   { return 0 }
func (VisitTheBlacksmithBlue) Pitch() int                  { return 3 }
func (VisitTheBlacksmithBlue) Attack() int                 { return 0 }
func (VisitTheBlacksmithBlue) Defense() int                { return 2 }
func (VisitTheBlacksmithBlue) Types() card.TypeSet         { return visitTheBlacksmithTypes }
func (VisitTheBlacksmithBlue) GoAgain() bool               { return true }
func (VisitTheBlacksmithBlue) Play(s *card.TurnState) int { return 0 }
