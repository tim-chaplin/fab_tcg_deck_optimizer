// Pilfer the Tomb — Generic Instant. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "Choose 1 or both; - Banish target instant from an opposing hero's graveyard.
// - Banish target yellow card from an opposing hero's graveyard."

package unplayable

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var pilferTheTombTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type PilferTheTombBlue struct{}

func (PilferTheTombBlue) ID() ids.CardID          { return ids.PilferTheTombBlue }
func (PilferTheTombBlue) Name() string            { return "Pilfer the Tomb" }
func (PilferTheTombBlue) Cost(*sim.TurnState) int { return 0 }
func (PilferTheTombBlue) Pitch() int              { return 3 }
func (PilferTheTombBlue) Attack() int             { return 0 }
func (PilferTheTombBlue) Defense() int            { return 0 }
func (PilferTheTombBlue) Types() card.TypeSet     { return pilferTheTombTypes }
func (PilferTheTombBlue) GoAgain() bool           { return false }
func (PilferTheTombBlue) Unplayable()             {}
func (PilferTheTombBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.Log(self, 0)
}
