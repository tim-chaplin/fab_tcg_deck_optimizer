// Pilfer the Tomb — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Banish target instant from an opposing hero's graveyard. - Banish
// target yellow card from an opposing hero's graveyard."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var pilferTheTombTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type PilferTheTombBlue struct{}

func (PilferTheTombBlue) ID() card.ID                               { return card.PilferTheTombBlue }
func (PilferTheTombBlue) Name() string                              { return "Pilfer the Tomb" }
func (PilferTheTombBlue) Cost(*card.TurnState) int                  { return 0 }
func (PilferTheTombBlue) Pitch() int                                { return 3 }
func (PilferTheTombBlue) Attack() int                               { return 0 }
func (PilferTheTombBlue) Defense() int                              { return 0 }
func (PilferTheTombBlue) Types() card.TypeSet                       { return pilferTheTombTypes }
func (PilferTheTombBlue) GoAgain() bool                             { return false }
// not implemented: Instant banish from an opposing graveyard / aura
func (PilferTheTombBlue) NotImplemented()                           {}
func (PilferTheTombBlue) Play(*card.TurnState, *card.CardState) int { return 0 }
