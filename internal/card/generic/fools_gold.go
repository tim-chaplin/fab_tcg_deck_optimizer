// Fool's Gold — Generic Resource. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "When this is discarded, create a Gold token."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var foolsGoldTypes = card.NewTypeSet(card.TypeGeneric)

type FoolsGoldYellow struct{}

func (FoolsGoldYellow) ID() card.ID              { return card.FoolsGoldYellow }
func (FoolsGoldYellow) Name() string             { return "Fool's Gold" }
func (FoolsGoldYellow) Cost(*card.TurnState) int { return 0 }
func (FoolsGoldYellow) Pitch() int               { return 2 }
func (FoolsGoldYellow) Attack() int              { return 0 }
func (FoolsGoldYellow) Defense() int             { return 0 }
func (FoolsGoldYellow) Types() card.TypeSet      { return foolsGoldTypes }
func (FoolsGoldYellow) GoAgain() bool            { return false }

// not implemented: discard trigger creates a Gold token
func (FoolsGoldYellow) NotImplemented()                              {}
func (FoolsGoldYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
