// Sigil of Solace — Generic Instant. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Gain 3{h}"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfSolaceTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type SigilOfSolaceRed struct{}

func (SigilOfSolaceRed) ID() ids.CardID          { return ids.SigilOfSolaceRed }
func (SigilOfSolaceRed) Name() string            { return "Sigil of Solace" }
func (SigilOfSolaceRed) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSolaceRed) Pitch() int              { return 1 }
func (SigilOfSolaceRed) Attack() int             { return 0 }
func (SigilOfSolaceRed) Defense() int            { return 0 }
func (SigilOfSolaceRed) Types() card.TypeSet     { return sigilOfSolaceTypes }
func (SigilOfSolaceRed) GoAgain() bool           { return false }
func (SigilOfSolaceRed) NotSilverAgeLegal()      {}

// not implemented: 3/2/1{h} gain (also banlisted)
func (SigilOfSolaceRed) NotImplemented()                            {}
func (SigilOfSolaceRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type SigilOfSolaceYellow struct{}

func (SigilOfSolaceYellow) ID() ids.CardID          { return ids.SigilOfSolaceYellow }
func (SigilOfSolaceYellow) Name() string            { return "Sigil of Solace" }
func (SigilOfSolaceYellow) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSolaceYellow) Pitch() int              { return 2 }
func (SigilOfSolaceYellow) Attack() int             { return 0 }
func (SigilOfSolaceYellow) Defense() int            { return 0 }
func (SigilOfSolaceYellow) Types() card.TypeSet     { return sigilOfSolaceTypes }
func (SigilOfSolaceYellow) GoAgain() bool           { return false }
func (SigilOfSolaceYellow) NotSilverAgeLegal()      {}

// not implemented: 3/2/1{h} gain (also banlisted)
func (SigilOfSolaceYellow) NotImplemented()                            {}
func (SigilOfSolaceYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type SigilOfSolaceBlue struct{}

func (SigilOfSolaceBlue) ID() ids.CardID          { return ids.SigilOfSolaceBlue }
func (SigilOfSolaceBlue) Name() string            { return "Sigil of Solace" }
func (SigilOfSolaceBlue) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSolaceBlue) Pitch() int              { return 3 }
func (SigilOfSolaceBlue) Attack() int             { return 0 }
func (SigilOfSolaceBlue) Defense() int            { return 0 }
func (SigilOfSolaceBlue) Types() card.TypeSet     { return sigilOfSolaceTypes }
func (SigilOfSolaceBlue) GoAgain() bool           { return false }
func (SigilOfSolaceBlue) NotSilverAgeLegal()      {}

// not implemented: 3/2/1{h} gain (also banlisted)
func (SigilOfSolaceBlue) NotImplemented()                            {}
func (SigilOfSolaceBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
