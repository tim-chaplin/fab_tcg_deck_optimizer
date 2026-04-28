// Infectious Host — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, if you control a Frailty token, create a Frailty token under
// their control, then repeat for Inertia and Bloodrot Pox."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var infectiousHostTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type InfectiousHostRed struct{}

func (InfectiousHostRed) ID() ids.CardID          { return ids.InfectiousHostRed }
func (InfectiousHostRed) Name() string            { return "Infectious Host" }
func (InfectiousHostRed) Cost(*sim.TurnState) int { return 0 }
func (InfectiousHostRed) Pitch() int              { return 1 }
func (InfectiousHostRed) Attack() int             { return 4 }
func (InfectiousHostRed) Defense() int            { return 2 }
func (InfectiousHostRed) Types() card.TypeSet     { return infectiousHostTypes }
func (InfectiousHostRed) GoAgain() bool           { return false }

// not implemented: frailty/inertia/bloodrot pox tokens
func (InfectiousHostRed) NotImplemented() {}
func (c InfectiousHostRed) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type InfectiousHostYellow struct{}

func (InfectiousHostYellow) ID() ids.CardID          { return ids.InfectiousHostYellow }
func (InfectiousHostYellow) Name() string            { return "Infectious Host" }
func (InfectiousHostYellow) Cost(*sim.TurnState) int { return 0 }
func (InfectiousHostYellow) Pitch() int              { return 2 }
func (InfectiousHostYellow) Attack() int             { return 3 }
func (InfectiousHostYellow) Defense() int            { return 2 }
func (InfectiousHostYellow) Types() card.TypeSet     { return infectiousHostTypes }
func (InfectiousHostYellow) GoAgain() bool           { return false }

// not implemented: frailty/inertia/bloodrot pox tokens
func (InfectiousHostYellow) NotImplemented() {}
func (c InfectiousHostYellow) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

type InfectiousHostBlue struct{}

func (InfectiousHostBlue) ID() ids.CardID          { return ids.InfectiousHostBlue }
func (InfectiousHostBlue) Name() string            { return "Infectious Host" }
func (InfectiousHostBlue) Cost(*sim.TurnState) int { return 0 }
func (InfectiousHostBlue) Pitch() int              { return 3 }
func (InfectiousHostBlue) Attack() int             { return 2 }
func (InfectiousHostBlue) Defense() int            { return 2 }
func (InfectiousHostBlue) Types() card.TypeSet     { return infectiousHostTypes }
func (InfectiousHostBlue) GoAgain() bool           { return false }

// not implemented: frailty/inertia/bloodrot pox tokens
func (InfectiousHostBlue) NotImplemented() {}
func (c InfectiousHostBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
