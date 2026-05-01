// Lead the Charge — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time you play an action card with cost 0 or greater this turn, gain 1 action
// point. **Go again**"

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var leadTheChargeTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

type LeadTheChargeRed struct{}

func (LeadTheChargeRed) ID() ids.CardID          { return ids.LeadTheChargeRed }
func (LeadTheChargeRed) Name() string            { return "Lead the Charge" }
func (LeadTheChargeRed) Cost(*sim.TurnState) int { return 0 }
func (LeadTheChargeRed) Pitch() int              { return 1 }
func (LeadTheChargeRed) Attack() int             { return 0 }
func (LeadTheChargeRed) Defense() int            { return 2 }
func (LeadTheChargeRed) Types() card.TypeSet     { return leadTheChargeTypes }
func (LeadTheChargeRed) GoAgain() bool           { return true }

// not implemented: action point grant
func (LeadTheChargeRed) NotImplemented()                            {}
func (LeadTheChargeRed) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type LeadTheChargeYellow struct{}

func (LeadTheChargeYellow) ID() ids.CardID          { return ids.LeadTheChargeYellow }
func (LeadTheChargeYellow) Name() string            { return "Lead the Charge" }
func (LeadTheChargeYellow) Cost(*sim.TurnState) int { return 0 }
func (LeadTheChargeYellow) Pitch() int              { return 2 }
func (LeadTheChargeYellow) Attack() int             { return 0 }
func (LeadTheChargeYellow) Defense() int            { return 2 }
func (LeadTheChargeYellow) Types() card.TypeSet     { return leadTheChargeTypes }
func (LeadTheChargeYellow) GoAgain() bool           { return true }

// not implemented: action point grant
func (LeadTheChargeYellow) NotImplemented()                            {}
func (LeadTheChargeYellow) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }

type LeadTheChargeBlue struct{}

func (LeadTheChargeBlue) ID() ids.CardID          { return ids.LeadTheChargeBlue }
func (LeadTheChargeBlue) Name() string            { return "Lead the Charge" }
func (LeadTheChargeBlue) Cost(*sim.TurnState) int { return 0 }
func (LeadTheChargeBlue) Pitch() int              { return 3 }
func (LeadTheChargeBlue) Attack() int             { return 0 }
func (LeadTheChargeBlue) Defense() int            { return 2 }
func (LeadTheChargeBlue) Types() card.TypeSet     { return leadTheChargeTypes }
func (LeadTheChargeBlue) GoAgain() bool           { return true }

// not implemented: action point grant
func (LeadTheChargeBlue) NotImplemented()                            {}
func (LeadTheChargeBlue) Play(s *sim.TurnState, self *sim.CardState) { s.Log(self, 0) }
