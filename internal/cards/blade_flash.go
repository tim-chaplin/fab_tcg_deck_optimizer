// Blade Flash — Generic Attack Reaction. Cost 1. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target sword attack gains **go again**."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.
// The go-again grant is modelled by bumping ActionPoints eagerly at Play time.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (BladeFlashBlue) ARTargetAllowed(c sim.Card, _ int8) bool {
	t := c.Types()
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (BladeFlashBlue) Play(s *sim.TurnState, l sim.Logger, _ *sim.CardState) {
	if s.AttackReactionTarget() == nil {
		return
	}
	s.ActionPoints++
}
