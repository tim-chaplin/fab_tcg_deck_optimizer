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
// matching extra action point is credited via AddActionPoints so the next card in the
// permutation can spend it (the buffed attack already passed its own EffectiveGoAgain
// check before the AR step, so the grant arrives too late to retroactively bank that
// attack's AP — AddActionPoints back-fills it). Value falls out naturally from whatever
// card spends the extra AP, so the handler credits no Value of its own.
func featherfootFire(ge card.GameEngine, l card.Logger, self card.Item, triggeringCard *card.CardState, _ triggertype.Type) {
	triggeringCard.GrantedGoAgain = true
	ge.AddActionPoints(1)
	self.Destroy(true)
	l.AppendPostTrigger(triggeringCard.Card.DisplayName(),
		"Talisman of Featherfoot destroyed to grant go again on a +1 reaction-step buff", 0)
}

func (TalismanOfFeatherfootYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItem(self.Card, triggertype.AttackBuffedByReaction, featherfootFire, false, nil)
}
