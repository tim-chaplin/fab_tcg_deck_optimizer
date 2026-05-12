// Razor Reflex — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "Choose 1; Target dagger or sword weapon attack gets +N{p}. Target attack action
// card with cost 1 or less gets +N{p} and 'When this hits, it gets **go again**.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Mode 1's on-hit go-again rider is modelled eagerly: when g.LikelyToHit on the post-buff
// target returns true, the AR grants 1 AP at Play time. That mirrors the chain runner's
// existing LikelyToHit-based on-hit gate (used by Runechant arcane and OnHit handlers) and
// makes the AP available for the next chain step's gate.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// razorReflexAccepts is the per-mode target predicate. Mode 0 gates on sword weapon
// attack; mode 1 gates on cost-≤1 attack action. The chain runner runs this for the
// chosen Mode and rejects the permutation when it returns false, so razorReflexPlay can
// apply the buff unconditionally.
func razorReflexAccepts(g card.GameEngine, c card.Card, mode int8) bool {
	t := c.Types(nil)
	switch mode {
	case 0:
		return t.Has(card.TypeSword) && t.IsWeaponAttack()
	case 1:
		return t.IsAttackAction() && c.Cost(g) <= 1
	}
	return false
}

// razorReflexPlay applies the chosen mode's effect. The chain runner already validated
// the target via razorReflexAccepts, so the buff lands directly. Mode 1 additionally
// fires the on-hit go-again rider eagerly when the post-buff target is likely to hit.
func razorReflexPlay(g card.GameEngine, l card.Logger, self *card.CardState, n int) {
	target := g.AttackReactionTarget()
	if target == nil {
		return
	}
	self.GrantAttackReactionBuff(g, l, n)
	if self.Mode == 1 && g.LikelyToHit(target) {
		g.AddActionPoints(1)
	}
}

func (RazorReflexRed) Modes() int { return 2 }
func (RazorReflexRed) ARTargetAllowed(g card.GameEngine, c card.Card, mode int8) bool {
	return razorReflexAccepts(g, c, mode)
}
func (RazorReflexRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	razorReflexPlay(g, l, self, 3)
}

func (RazorReflexYellow) Modes() int { return 2 }
func (RazorReflexYellow) ARTargetAllowed(g card.GameEngine, c card.Card, mode int8) bool {
	return razorReflexAccepts(g, c, mode)
}
func (RazorReflexYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	razorReflexPlay(g, l, self, 2)
}

func (RazorReflexBlue) Modes() int { return 2 }
func (RazorReflexBlue) ARTargetAllowed(g card.GameEngine, c card.Card, mode int8) bool {
	return razorReflexAccepts(g, c, mode)
}
func (RazorReflexBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {
	razorReflexPlay(g, l, self, 1)
}
