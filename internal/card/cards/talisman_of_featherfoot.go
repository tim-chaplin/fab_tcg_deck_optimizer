// Talisman of Featherfoot — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** When an attack you control gains exactly +1{p} from an effect during the
// reaction step, destroy Talisman of Featherfoot and the attack gains **go again**."
//
// Subscribes to triggertype.AttackBuffedByReaction, which the engine raises from
// CardState.GrantAttackReactionBuff whenever a reaction-step buff applies exactly +1{p}
// to an attack. The handler flips the triggering attack's GrantedGoAgain and destroys
// the talisman. Larger buffs (+2 or more from a single source) skip the trigger by
// design; the "exactly +1" gate lives at the fire site, not in this handler.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// featherfootFire grants the buffed attack go again and self-destructs. The fire site
// already verified the +1 gate, so the handler unconditionally applies the payoff.
// GrantedGoAgain records the keyword grant on the buffed attack's CardState; the
// attack's own EffectiveGoAgain check (deferred until after all attack reactions
// resolve) banks the corresponding action point.
func featherfootFire(_ card.GameEngine, l card.Logger, self card.Item, triggeringCard *card.CardState, _ triggertype.Type) {
	triggeringCard.GrantedGoAgain = true
	self.Destroy(true)
	l.AppendPostTrigger(triggeringCard.Card.DisplayName(),
		"Talisman of Featherfoot destroyed to grant go again on a +1 reaction-step buff", 0)
}

func (TalismanOfFeatherfootYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItem(self.Card, triggertype.AttackBuffedByReaction, featherfootFire, false, nil)
}
