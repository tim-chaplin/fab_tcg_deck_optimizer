// Oasis Respite — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Prevent the next 4 damage that would be dealt to target hero this turn by a source of your
// choice. If they have less life than each other hero, they may gain 1{h}."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var oasisRespiteTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type OasisRespiteRed struct{}

func (OasisRespiteRed) ID() card.ID              { return card.OasisRespiteRed }
func (OasisRespiteRed) Name() string             { return "Oasis Respite" }
func (OasisRespiteRed) Cost(*card.TurnState) int { return 1 }
func (OasisRespiteRed) Pitch() int               { return 1 }
func (OasisRespiteRed) Attack() int              { return 0 }
func (OasisRespiteRed) Defense() int             { return 0 }
func (OasisRespiteRed) Types() card.TypeSet      { return oasisRespiteTypes }
func (OasisRespiteRed) GoAgain() bool            { return false }

// not implemented: Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}
func (OasisRespiteRed) NotImplemented()                              {}
func (OasisRespiteRed) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type OasisRespiteYellow struct{}

func (OasisRespiteYellow) ID() card.ID              { return card.OasisRespiteYellow }
func (OasisRespiteYellow) Name() string             { return "Oasis Respite" }
func (OasisRespiteYellow) Cost(*card.TurnState) int { return 1 }
func (OasisRespiteYellow) Pitch() int               { return 2 }
func (OasisRespiteYellow) Attack() int              { return 0 }
func (OasisRespiteYellow) Defense() int             { return 0 }
func (OasisRespiteYellow) Types() card.TypeSet      { return oasisRespiteTypes }
func (OasisRespiteYellow) GoAgain() bool            { return false }

// not implemented: Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}
func (OasisRespiteYellow) NotImplemented()                              {}
func (OasisRespiteYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type OasisRespiteBlue struct{}

func (OasisRespiteBlue) ID() card.ID              { return card.OasisRespiteBlue }
func (OasisRespiteBlue) Name() string             { return "Oasis Respite" }
func (OasisRespiteBlue) Cost(*card.TurnState) int { return 1 }
func (OasisRespiteBlue) Pitch() int               { return 3 }
func (OasisRespiteBlue) Attack() int              { return 0 }
func (OasisRespiteBlue) Defense() int             { return 0 }
func (OasisRespiteBlue) Types() card.TypeSet      { return oasisRespiteTypes }
func (OasisRespiteBlue) GoAgain() bool            { return false }

// not implemented: Instant 'prevent N damage from chosen source to target hero'; conditional 1{h}
func (OasisRespiteBlue) NotImplemented()                              {}
func (OasisRespiteBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
