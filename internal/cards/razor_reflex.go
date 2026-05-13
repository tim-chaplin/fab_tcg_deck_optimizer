// Razor Reflex — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "Choose 1; Target dagger or sword weapon attack gets +N{p}. Target attack action
// card with cost 1 or less gets +N{p} and 'When this hits, it gets **go again**.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Mode 1's on-hit go-again rider is modelled eagerly: when ge.LikelyToHit on the post-buff
// target returns true, the AR grants 1 AP at Play time so the AP is available for the next
// chain step's gate.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// razorReflexAccepts is the per-mode target predicate. Mode 0 gates on sword weapon attack;
// mode 1 gates on cost-≤1 attack action.
func razorReflexAccepts(ge card.GameEngine, c card.Card, mode int8) bool {
	t := c.Types(nil)
	switch mode {
	case 0:
		return t.Has(card.TypeSword) && t.IsWeaponAttack()
	case 1:
		return t.IsAttackAction() && c.Cost(ge) <= 1
	}
	return false
}

// razorReflexPlay applies the chosen mode's effect. Mode 1 additionally fires the on-hit
// go-again rider eagerly when the post-buff target is likely to hit.
func razorReflexPlay(ge card.GameEngine, l card.Logger, self *card.CardState, n int) {
	target := ge.AttackReactionTarget()
	if target == nil {
		return
	}
	self.GrantAttackReactionBuff(ge, l, n)
	if self.Mode == 1 && ge.LikelyToHit(target) {
		ge.AddActionPoints(1)
	}
}

func (RazorReflexRed) Modes() int { return 2 }
func (RazorReflexRed) ARTargetAllowed(ge card.GameEngine, c card.Card, mode int8) bool {
	return razorReflexAccepts(ge, c, mode)
}
func (RazorReflexRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	razorReflexPlay(ge, l, self, 3)
}

func (RazorReflexYellow) Modes() int { return 2 }
func (RazorReflexYellow) ARTargetAllowed(ge card.GameEngine, c card.Card, mode int8) bool {
	return razorReflexAccepts(ge, c, mode)
}
func (RazorReflexYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	razorReflexPlay(ge, l, self, 2)
}

func (RazorReflexBlue) Modes() int { return 2 }
func (RazorReflexBlue) ARTargetAllowed(ge card.GameEngine, c card.Card, mode int8) bool {
	return razorReflexAccepts(ge, c, mode)
}
func (RazorReflexBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	razorReflexPlay(ge, l, self, 1)
}
