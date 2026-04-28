// Pummel — Generic Attack Reaction. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Choose 1; - Target club or hammer weapon attack gains +4{p}. - Target attack action card
// with cost 2 or more gets +4{p} and "When this hits a hero, they discard a card.""

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var pummelTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type PummelRed struct{}

func (PummelRed) ID() ids.CardID          { return ids.PummelRed }
func (PummelRed) Name() string            { return "Pummel" }
func (PummelRed) Cost(*sim.TurnState) int { return 2 }
func (PummelRed) Pitch() int              { return 1 }
func (PummelRed) Attack() int             { return 0 }
func (PummelRed) Defense() int            { return 2 }
func (PummelRed) Types() card.TypeSet     { return pummelTypes }
func (PummelRed) GoAgain() bool           { return false }

// not implemented: modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action
// (on-hit discard)
func (PummelRed) NotImplemented()                            {}
func (PummelRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type PummelYellow struct{}

func (PummelYellow) ID() ids.CardID          { return ids.PummelYellow }
func (PummelYellow) Name() string            { return "Pummel" }
func (PummelYellow) Cost(*sim.TurnState) int { return 2 }
func (PummelYellow) Pitch() int              { return 2 }
func (PummelYellow) Attack() int             { return 0 }
func (PummelYellow) Defense() int            { return 2 }
func (PummelYellow) Types() card.TypeSet     { return pummelTypes }
func (PummelYellow) GoAgain() bool           { return false }

// not implemented: modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action
// (on-hit discard)
func (PummelYellow) NotImplemented()                            {}
func (PummelYellow) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

type PummelBlue struct{}

func (PummelBlue) ID() ids.CardID          { return ids.PummelBlue }
func (PummelBlue) Name() string            { return "Pummel" }
func (PummelBlue) Cost(*sim.TurnState) int { return 2 }
func (PummelBlue) Pitch() int              { return 3 }
func (PummelBlue) Attack() int             { return 0 }
func (PummelBlue) Defense() int            { return 2 }
func (PummelBlue) Types() card.TypeSet     { return pummelTypes }
func (PummelBlue) GoAgain() bool           { return false }

// not implemented: modal AR +4{p}: club/hammer weapon attack OR cost-2+ attack action
// (on-hit discard)
func (PummelBlue) NotImplemented()                            {}
func (PummelBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
