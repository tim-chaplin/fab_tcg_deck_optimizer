// Razor Reflex — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "Choose 1; Target dagger or sword weapon attack gets +N{p}. Target attack action
// card with cost 1 or less gets +N{p} and 'When this hits, it gets **go again**.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Mode 1's on-hit go-again rider is modelled eagerly: when sim.LikelyToHit on the post-buff
// target returns true, the AR grants 1 AP at Play time. That mirrors the chain runner's
// existing LikelyToHit-based on-hit gate (used by Runechant arcane and OnHit handlers) and
// makes the AP available for the next chain step's gate.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var razorReflexTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

// razorReflexAccepts is the per-mode target predicate. Mode 0 gates on sword weapon
// attack; mode 1 gates on cost-≤1 attack action. The chain runner runs this for the
// chosen Mode and rejects the permutation when it returns false, so razorReflexPlay can
// apply the buff unconditionally.
//
// Reads Cost against an empty TurnState; no variable-cost cost-≤1 attack actions exist
// in the pool.
func razorReflexAccepts(c sim.Card, mode int8) bool {
	t := c.Types()
	switch mode {
	case 0:
		return t.Has(card.TypeSword) && t.IsWeaponAttack()
	case 1:
		return t.IsAttackAction() && c.Cost(&sim.TurnState{}) <= 1
	}
	return false
}

// razorReflexPlay applies the chosen mode's effect. The chain runner already validated
// the target via razorReflexAccepts, so the buff lands directly. Mode 1 additionally
// fires the on-hit go-again rider eagerly when the post-buff target is likely to hit.
func razorReflexPlay(s *sim.TurnState, self *sim.CardState, n int) {
	target := s.AttackReactionTarget()
	if target == nil {
		return
	}
	sim.GrantAttackReactionBuff(s, self, n)
	if self.Mode == 1 && sim.LikelyToHit(target) {
		s.ActionPoints++
	}
}

type RazorReflexRed struct{}

func (RazorReflexRed) ID() ids.CardID          { return ids.RazorReflexRed }
func (RazorReflexRed) Name() string            { return "Razor Reflex" }
func (RazorReflexRed) Cost(*sim.TurnState) int { return 1 }
func (RazorReflexRed) Pitch() int              { return 1 }
func (RazorReflexRed) Attack() int             { return 0 }
func (RazorReflexRed) Defense() int            { return 2 }
func (RazorReflexRed) Types() card.TypeSet     { return razorReflexTypes }
func (RazorReflexRed) GoAgain() bool           { return false }
func (RazorReflexRed) Modes() int              { return 2 }
func (RazorReflexRed) ARTargetAllowed(c sim.Card, mode int8) bool {
	return razorReflexAccepts(c, mode)
}
func (RazorReflexRed) Play(s *sim.TurnState, self *sim.CardState) {
	razorReflexPlay(s, self, 3)
}

type RazorReflexYellow struct{}

func (RazorReflexYellow) ID() ids.CardID          { return ids.RazorReflexYellow }
func (RazorReflexYellow) Name() string            { return "Razor Reflex" }
func (RazorReflexYellow) Cost(*sim.TurnState) int { return 1 }
func (RazorReflexYellow) Pitch() int              { return 2 }
func (RazorReflexYellow) Attack() int             { return 0 }
func (RazorReflexYellow) Defense() int            { return 2 }
func (RazorReflexYellow) Types() card.TypeSet     { return razorReflexTypes }
func (RazorReflexYellow) GoAgain() bool           { return false }
func (RazorReflexYellow) Modes() int              { return 2 }
func (RazorReflexYellow) ARTargetAllowed(c sim.Card, mode int8) bool {
	return razorReflexAccepts(c, mode)
}
func (RazorReflexYellow) Play(s *sim.TurnState, self *sim.CardState) {
	razorReflexPlay(s, self, 2)
}

type RazorReflexBlue struct{}

func (RazorReflexBlue) ID() ids.CardID          { return ids.RazorReflexBlue }
func (RazorReflexBlue) Name() string            { return "Razor Reflex" }
func (RazorReflexBlue) Cost(*sim.TurnState) int { return 1 }
func (RazorReflexBlue) Pitch() int              { return 3 }
func (RazorReflexBlue) Attack() int             { return 0 }
func (RazorReflexBlue) Defense() int            { return 2 }
func (RazorReflexBlue) Types() card.TypeSet     { return razorReflexTypes }
func (RazorReflexBlue) GoAgain() bool           { return false }
func (RazorReflexBlue) Modes() int              { return 2 }
func (RazorReflexBlue) ARTargetAllowed(c sim.Card, mode int8) bool {
	return razorReflexAccepts(c, mode)
}
func (RazorReflexBlue) Play(s *sim.TurnState, self *sim.CardState) {
	razorReflexPlay(s, self, 1)
}
