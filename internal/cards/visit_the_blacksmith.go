// Visit the Blacksmith — Generic Action. Cost 0, Pitch 3, Defense 2. Only printed in Blue.
//
// Text: "Your next sword attack this turn gains +1{p}. **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var visitTheBlacksmithTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type VisitTheBlacksmithBlue struct{}

func (VisitTheBlacksmithBlue) ID() ids.CardID          { return ids.VisitTheBlacksmithBlue }
func (VisitTheBlacksmithBlue) Name() string            { return "Visit the Blacksmith" }
func (VisitTheBlacksmithBlue) Cost(*sim.TurnState) int { return 0 }
func (VisitTheBlacksmithBlue) Pitch() int              { return 3 }
func (VisitTheBlacksmithBlue) Attack() int             { return 0 }
func (VisitTheBlacksmithBlue) Defense() int            { return 2 }
func (VisitTheBlacksmithBlue) Types() card.TypeSet     { return visitTheBlacksmithTypes }
func (VisitTheBlacksmithBlue) GoAgain() bool           { return true }

// not implemented: next-sword-attack +1{p} grant (weapon chain not peeked)
func (VisitTheBlacksmithBlue) NotImplemented()                            {}
func (VisitTheBlacksmithBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
