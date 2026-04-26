// Shatter Sorcery — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Destroy target aura permanent with Sigil in its name. - Prevent the
// next 1 arcane damage that would be dealt to target hero this turn."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var shatterSorceryTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type ShatterSorceryBlue struct{}

func (ShatterSorceryBlue) ID() card.ID                               { return card.ShatterSorceryBlue }
func (ShatterSorceryBlue) Name() string                              { return "Shatter Sorcery" }
func (ShatterSorceryBlue) Cost(*card.TurnState) int                  { return 0 }
func (ShatterSorceryBlue) Pitch() int                                { return 3 }
func (ShatterSorceryBlue) Attack() int                               { return 0 }
func (ShatterSorceryBlue) Defense() int                              { return 0 }
func (ShatterSorceryBlue) Types() card.TypeSet                       { return shatterSorceryTypes }
func (ShatterSorceryBlue) GoAgain() bool                             { return false }
// not implemented: Instant: destroy a Sigil aura, and/or prevent 1 arcane damage
func (ShatterSorceryBlue) NotImplemented()                           {}
func (ShatterSorceryBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
