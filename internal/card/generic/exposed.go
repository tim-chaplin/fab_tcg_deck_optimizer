// Exposed — Generic Attack Reaction. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "If you are **marked**, you can't play this. Target attack gets +1{p}. **Mark** the
// defending hero."
//
// Stub only — marked NotImplemented so the optimizer skips it. The printed effect isn't modelled;
// Play returns 0.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var exposedTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type ExposedBlue struct{}

func (ExposedBlue) ID() card.ID                               { return card.ExposedBlue }
func (ExposedBlue) Name() string                              { return "Exposed (Blue)" }
func (ExposedBlue) Cost(*card.TurnState) int                  { return 0 }
func (ExposedBlue) Pitch() int                                { return 3 }
func (ExposedBlue) Attack() int                               { return 0 }
func (ExposedBlue) Defense() int                              { return 0 }
func (ExposedBlue) Types() card.TypeSet                       { return exposedTypes }
func (ExposedBlue) GoAgain() bool                             { return false }
func (ExposedBlue) NotImplemented()                           {}
func (ExposedBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
