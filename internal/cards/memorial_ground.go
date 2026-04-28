// Memorial Ground — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Put target attack action card with cost 2 or less from your graveyard on top of your
// deck."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var memorialGroundTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type MemorialGroundRed struct{}

func (MemorialGroundRed) ID() ids.CardID          { return ids.MemorialGroundRed }
func (MemorialGroundRed) Name() string            { return "Memorial Ground" }
func (MemorialGroundRed) Cost(*sim.TurnState) int { return 0 }
func (MemorialGroundRed) Pitch() int              { return 1 }
func (MemorialGroundRed) Attack() int             { return 0 }
func (MemorialGroundRed) Defense() int            { return 0 }
func (MemorialGroundRed) Types() card.TypeSet     { return memorialGroundTypes }
func (MemorialGroundRed) GoAgain() bool           { return false }

// not implemented: Instant 'graveyard → top of deck' for low-cost attack action card
func (MemorialGroundRed) NotImplemented()                            {}
func (MemorialGroundRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type MemorialGroundYellow struct{}

func (MemorialGroundYellow) ID() ids.CardID          { return ids.MemorialGroundYellow }
func (MemorialGroundYellow) Name() string            { return "Memorial Ground" }
func (MemorialGroundYellow) Cost(*sim.TurnState) int { return 0 }
func (MemorialGroundYellow) Pitch() int              { return 2 }
func (MemorialGroundYellow) Attack() int             { return 0 }
func (MemorialGroundYellow) Defense() int            { return 0 }
func (MemorialGroundYellow) Types() card.TypeSet     { return memorialGroundTypes }
func (MemorialGroundYellow) GoAgain() bool           { return false }

// not implemented: Instant 'graveyard → top of deck' for low-cost attack action card
func (MemorialGroundYellow) NotImplemented()                            {}
func (MemorialGroundYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type MemorialGroundBlue struct{}

func (MemorialGroundBlue) ID() ids.CardID          { return ids.MemorialGroundBlue }
func (MemorialGroundBlue) Name() string            { return "Memorial Ground" }
func (MemorialGroundBlue) Cost(*sim.TurnState) int { return 0 }
func (MemorialGroundBlue) Pitch() int              { return 3 }
func (MemorialGroundBlue) Attack() int             { return 0 }
func (MemorialGroundBlue) Defense() int            { return 0 }
func (MemorialGroundBlue) Types() card.TypeSet     { return memorialGroundTypes }
func (MemorialGroundBlue) GoAgain() bool           { return false }

// not implemented: Instant 'graveyard → top of deck' for low-cost attack action card
func (MemorialGroundBlue) NotImplemented()                            {}
func (MemorialGroundBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
