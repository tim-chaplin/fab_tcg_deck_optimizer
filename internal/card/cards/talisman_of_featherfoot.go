// Talisman of Featherfoot — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When an attack you control gains exactly +1{p} from an effect during the
// reaction step, destroy Talisman of Featherfoot and the attack gains **go again**."
//
// Subscribes to triggertype.AttackBuffedByReaction, which fires for every positive
// reaction-step buff. The "exactly +1{p}" gate is enforced inside this handler by reading
// the buff size from ge.AttackReactionBuffDelta — the gate is Featherfoot-specific so it
// doesn't belong in the general GrantAttackReactionBuff helper.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// featherfootFire grants the buffed attack go again and self-destructs when the firing
// buff is exactly +1{p}. GrantedGoAgain records the keyword grant on the buffed attack's
// CardState; the attack's own EffectiveGoAgain check (deferred until after all attack
// reactions resolve) banks the corresponding action point.
func featherfootFire(ge card.GameEngine, l card.Logger, self card.Item, triggeringCard *card.CardState, _ triggertype.Type) {
	if ge.AttackReactionBuffDelta() != 1 {
		return
	}
	triggeringCard.GrantedGoAgain = true
	self.Destroy(true)
	l.AppendPostTrigger(triggeringCard.Card.DisplayName(),
		"Talisman of Featherfoot destroyed to grant go again on a +1 reaction-step buff", 0)
}

func (TalismanOfFeatherfootYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItem(self.Card, triggertype.AttackBuffedByReaction, featherfootFire, false, nil)
}
