// Calming Breeze — Generic Instant. Cost 0. Printed pitch variants: Red 1.
//
// Text: "The next 3 times you would be dealt damage this turn, prevent 1 of that damage."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var calmingBreezeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type CalmingBreezeRed struct{}

func (CalmingBreezeRed) ID() ids.CardID          { return ids.CalmingBreezeRed }
func (CalmingBreezeRed) Name() string            { return "Calming Breeze" }
func (CalmingBreezeRed) Cost(*sim.TurnState) int { return 0 }
func (CalmingBreezeRed) Pitch() int              { return 1 }
func (CalmingBreezeRed) Attack() int             { return 0 }
func (CalmingBreezeRed) Defense() int            { return 0 }
func (CalmingBreezeRed) Types() card.TypeSet     { return calmingBreezeTypes }
func (CalmingBreezeRed) GoAgain() bool           { return false }

// not implemented: Instant 'prevent 1 of each of the next 3 damage events'
func (CalmingBreezeRed) NotImplemented()                            {}
func (CalmingBreezeRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
