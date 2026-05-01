// Oasis Respite — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Prevent the next 4 damage that would be dealt to target hero this turn by a source of your
// choice. If they have less life than each other hero, they may gain 1{h}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var oasisRespiteTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type OasisRespiteRed struct{}

func (OasisRespiteRed) ID() ids.CardID          { return ids.OasisRespiteRed }
func (OasisRespiteRed) Name() string            { return "Oasis Respite" }
func (OasisRespiteRed) Cost(*sim.TurnState) int { return 1 }
func (OasisRespiteRed) Pitch() int              { return 1 }
func (OasisRespiteRed) Attack() int             { return 0 }
func (OasisRespiteRed) Defense() int            { return 0 }
func (OasisRespiteRed) Types() card.TypeSet     { return oasisRespiteTypes }
func (OasisRespiteRed) GoAgain() bool           { return false }

// not implemented: Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}
func (OasisRespiteRed) NotImplemented()                            {}
func (OasisRespiteRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type OasisRespiteYellow struct{}

func (OasisRespiteYellow) ID() ids.CardID          { return ids.OasisRespiteYellow }
func (OasisRespiteYellow) Name() string            { return "Oasis Respite" }
func (OasisRespiteYellow) Cost(*sim.TurnState) int { return 1 }
func (OasisRespiteYellow) Pitch() int              { return 2 }
func (OasisRespiteYellow) Attack() int             { return 0 }
func (OasisRespiteYellow) Defense() int            { return 0 }
func (OasisRespiteYellow) Types() card.TypeSet     { return oasisRespiteTypes }
func (OasisRespiteYellow) GoAgain() bool           { return false }

// not implemented: Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}
func (OasisRespiteYellow) NotImplemented()                            {}
func (OasisRespiteYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type OasisRespiteBlue struct{}

func (OasisRespiteBlue) ID() ids.CardID          { return ids.OasisRespiteBlue }
func (OasisRespiteBlue) Name() string            { return "Oasis Respite" }
func (OasisRespiteBlue) Cost(*sim.TurnState) int { return 1 }
func (OasisRespiteBlue) Pitch() int              { return 3 }
func (OasisRespiteBlue) Attack() int             { return 0 }
func (OasisRespiteBlue) Defense() int            { return 0 }
func (OasisRespiteBlue) Types() card.TypeSet     { return oasisRespiteTypes }
func (OasisRespiteBlue) GoAgain() bool           { return false }

// not implemented: Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}
func (OasisRespiteBlue) NotImplemented()                            {}
func (OasisRespiteBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
