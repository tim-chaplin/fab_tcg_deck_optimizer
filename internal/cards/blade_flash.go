// Blade Flash — Generic Attack Reaction. Cost 1. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target sword attack gains **go again**."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.
// The go-again grant is modelled by bumping ActionPoints eagerly at Play time.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (BladeFlashBlue) ARTargetAllowed(c card.Card, _ int8) bool {
	t := c.Types()
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (BladeFlashBlue) Play(s card.GameEngine, l card.Logger, _ *card.CardState) {
	if s.AttackReactionTarget() == nil {
		return
	}
	s.AddActionPoints(1)
}
