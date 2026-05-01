// Shatter Sorcery — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Destroy target aura permanent with Sigil in its name. - Prevent the
// next 1 arcane damage that would be dealt to target hero this turn."

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var shatterSorceryTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type ShatterSorceryBlue struct{}

func (ShatterSorceryBlue) ID() ids.CardID          { return ids.ShatterSorceryBlue }
func (ShatterSorceryBlue) Name() string            { return "Shatter Sorcery" }
func (ShatterSorceryBlue) Cost(*sim.TurnState) int { return 0 }
func (ShatterSorceryBlue) Pitch() int              { return 3 }
func (ShatterSorceryBlue) Attack() int             { return 0 }
func (ShatterSorceryBlue) Defense() int            { return 0 }
func (ShatterSorceryBlue) Types() card.TypeSet     { return shatterSorceryTypes }
func (ShatterSorceryBlue) GoAgain() bool           { return false }

// not implemented: Instant: destroy a Sigil aura, and/or prevent 1 arcane damage
func (ShatterSorceryBlue) NotImplemented()                            {}
func (ShatterSorceryBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
