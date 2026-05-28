// Talisman of Featherfoot — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When an attack you control gains exactly +1{p} from an effect during the
// reaction step, destroy Talisman of Featherfoot and the attack gains **go again**."
//
// Subscribes to triggertype.AttackBuffedByReaction, which fires for every positive
// reaction-step buff. The "exactly +1{p}" gate is enforced here by reading ctx.BuffDelta —
// the gate is Featherfoot-specific so it doesn't belong in the general
// GrantAttackReactionBuff helper that raises the trigger.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// featherfootFire grants the buffed attack go again and self-destructs when the firing
// buff is exactly +1{p}. GrantedGoAgain records the keyword grant on the buffed attack's
// CardState; the attack's own EffectiveGoAgain check (deferred until after all attack
// reactions resolve) banks the corresponding action point.
func featherfootFire(_ card.GameEngine, l card.Logger, self card.Item, ctx card.FireContext) {
	if ctx.BuffDelta != 1 {
		return
	}
	ctx.TriggeringCard.GrantedGoAgain = true
	self.Destroy(true)
	l.AppendPostTrigger(ctx.TriggeringCard.Card.DisplayName(),
		"Talisman of Featherfoot destroyed to grant go again on a +1 reaction-step buff", 0)
}

func (TalismanOfFeatherfootYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItem(self.Card, triggertype.AttackBuffedByReaction, featherfootFire, false, nil)
}
