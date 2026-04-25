// Destructive Tendencies — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Remove all counters from target item token. - Remove all counters from
// target aura token."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var destructiveTendenciesTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type DestructiveTendenciesBlue struct{}

func (DestructiveTendenciesBlue) ID() card.ID                               { return card.DestructiveTendenciesBlue }
func (DestructiveTendenciesBlue) Name() string                              { return "Destructive Tendencies (Blue)" }
func (DestructiveTendenciesBlue) Cost(*card.TurnState) int                  { return 0 }
func (DestructiveTendenciesBlue) Pitch() int                                { return 3 }
func (DestructiveTendenciesBlue) Attack() int                               { return 0 }
func (DestructiveTendenciesBlue) Defense() int                              { return 0 }
func (DestructiveTendenciesBlue) Types() card.TypeSet                       { return destructiveTendenciesTypes }
func (DestructiveTendenciesBlue) GoAgain() bool                             { return false }
// not implemented: Instant remove counters from item / aura tokens
func (DestructiveTendenciesBlue) NotImplemented()                           {}
func (DestructiveTendenciesBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
