// Reinforce the Line — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Target defending attack action card gains +4{d}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var reinforceTheLineTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type ReinforceTheLineRed struct{}

func (ReinforceTheLineRed) ID() ids.CardID          { return ids.ReinforceTheLineRed }
func (ReinforceTheLineRed) Name() string            { return "Reinforce the Line" }
func (ReinforceTheLineRed) Cost(*sim.TurnState) int { return 0 }
func (ReinforceTheLineRed) Pitch() int              { return 1 }
func (ReinforceTheLineRed) Attack() int             { return 0 }
func (ReinforceTheLineRed) Defense() int            { return 0 }
func (ReinforceTheLineRed) Types() card.TypeSet     { return reinforceTheLineTypes }
func (ReinforceTheLineRed) GoAgain() bool           { return false }

// not implemented: Instant +N{d} grant to a defending attack action card
func (ReinforceTheLineRed) NotImplemented()                            {}
func (ReinforceTheLineRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }

type ReinforceTheLineYellow struct{}

func (ReinforceTheLineYellow) ID() ids.CardID          { return ids.ReinforceTheLineYellow }
func (ReinforceTheLineYellow) Name() string            { return "Reinforce the Line" }
func (ReinforceTheLineYellow) Cost(*sim.TurnState) int { return 0 }
func (ReinforceTheLineYellow) Pitch() int              { return 2 }
func (ReinforceTheLineYellow) Attack() int             { return 0 }
func (ReinforceTheLineYellow) Defense() int            { return 0 }
func (ReinforceTheLineYellow) Types() card.TypeSet     { return reinforceTheLineTypes }
func (ReinforceTheLineYellow) GoAgain() bool           { return false }

// not implemented: Instant +N{d} grant to a defending attack action card
func (ReinforceTheLineYellow) NotImplemented()                            {}
func (ReinforceTheLineYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }

type ReinforceTheLineBlue struct{}

func (ReinforceTheLineBlue) ID() ids.CardID          { return ids.ReinforceTheLineBlue }
func (ReinforceTheLineBlue) Name() string            { return "Reinforce the Line" }
func (ReinforceTheLineBlue) Cost(*sim.TurnState) int { return 0 }
func (ReinforceTheLineBlue) Pitch() int              { return 3 }
func (ReinforceTheLineBlue) Attack() int             { return 0 }
func (ReinforceTheLineBlue) Defense() int            { return 0 }
func (ReinforceTheLineBlue) Types() card.TypeSet     { return reinforceTheLineTypes }
func (ReinforceTheLineBlue) GoAgain() bool           { return false }

// not implemented: Instant +N{d} grant to a defending attack action card
func (ReinforceTheLineBlue) NotImplemented()                            {}
func (ReinforceTheLineBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogChain(self, 0) }
