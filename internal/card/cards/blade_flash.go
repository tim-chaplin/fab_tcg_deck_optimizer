// Blade Flash — Generic Attack Reaction. Cost 1. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target sword attack gains **go again**."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.
// The go-again grant is modelled by bumping ActionPoints eagerly at Play time.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (BladeFlashBlue) ARTargetAllowed(_ card.GameEngine, target *card.CardState, _ int8) bool {
	t := target.Card.Types(nil)
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (BladeFlashBlue) Play(ge card.GameEngine, l card.Logger, _ *card.CardState) {
	if ge.AttackReactionTarget() == nil {
		return
	}
	ge.AddActionPoints(1)
}
