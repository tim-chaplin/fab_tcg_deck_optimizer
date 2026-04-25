// Calming Breeze — Generic Instant. Cost 0. Printed pitch variants: Red 1.
//
// Text: "The next 3 times you would be dealt damage this turn, prevent 1 of that damage."
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var calmingBreezeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type CalmingBreezeRed struct{}

func (CalmingBreezeRed) ID() card.ID                               { return card.CalmingBreezeRed }
func (CalmingBreezeRed) Name() string                              { return "Calming Breeze (Red)" }
func (CalmingBreezeRed) Cost(*card.TurnState) int                  { return 0 }
func (CalmingBreezeRed) Pitch() int                                { return 1 }
func (CalmingBreezeRed) Attack() int                               { return 0 }
func (CalmingBreezeRed) Defense() int                              { return 0 }
func (CalmingBreezeRed) Types() card.TypeSet                       { return calmingBreezeTypes }
func (CalmingBreezeRed) GoAgain() bool                             { return false }
func (CalmingBreezeRed) NotImplemented()                           {}
func (CalmingBreezeRed) Play(*card.TurnState, *card.CardState) int { return 0 }
