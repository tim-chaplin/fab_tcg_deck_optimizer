// Thrust — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1. Defense 2.
//
// Text: "Target sword attack gains +3{p}."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var thrustTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type ThrustRed struct{}

func (ThrustRed) ID() ids.CardID          { return ids.ThrustRed }
func (ThrustRed) Name() string            { return "Thrust" }
func (ThrustRed) Cost(*sim.TurnState) int { return 1 }
func (ThrustRed) Pitch() int              { return 1 }
func (ThrustRed) Attack() int             { return 0 }
func (ThrustRed) Defense() int            { return 2 }
func (ThrustRed) Types() card.TypeSet     { return thrustTypes }
func (ThrustRed) GoAgain() bool           { return false }

// not implemented: AR +3{p} buff to a target sword attack
func (ThrustRed) NotImplemented()                            {}
func (ThrustRed) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
