// Thrust — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1. Defense 2.
//
// Text: "Target sword attack gains +3{p}."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.

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
func (ThrustRed) ARTargetAllowed(c sim.Card) bool {
	t := c.Types()
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (ThrustRed) Play(s *sim.TurnState, self *sim.CardState) {
	sim.GrantAttackReactionBuff(s, ThrustRed{}.ARTargetAllowed, 3)
	s.Log(self, 0)
}
