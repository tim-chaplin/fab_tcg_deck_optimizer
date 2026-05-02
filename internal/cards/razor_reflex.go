// Razor Reflex — Generic Attack Reaction. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "Choose 1; Target dagger or sword weapon attack gets +N{p}. Target attack action
// card with cost 1 or less gets +N{p} and 'When this hits, it gets **go again**.'"
// (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelling fudge: mode 1's on-hit go-again rider is dropped — the chain runner's AP gate
// runs at chain-step resolution, before OnHit fires, so a post-hit GrantedGoAgain flip
// can't propagate to AP for the next chain step.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var razorReflexTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

// razorReflexAccepts is the union of mode 0 (sword weapon attack) and mode 1 (cost-≤1 attack
// action) target predicates.
func razorReflexAccepts(c sim.Card) bool {
	t := c.Types()
	if t.Has(card.TypeSword) && t.IsAttack() {
		return true
	}
	return razorReflexMode1Allowed(c)
}

// razorReflexMode1Allowed gates mode 1: cost-≤1 attack action card. Reads Cost against an
// empty TurnState; no variable-cost cost-≤1 attack actions exist in the pool.
func razorReflexMode1Allowed(c sim.Card) bool {
	if !c.Types().IsAttackAction() {
		return false
	}
	return c.Cost(&sim.TurnState{}) <= 1
}

// razorReflexPlay applies the chosen mode's +N{p} buff. Mode 0 fires on sword weapon attacks;
// mode 1 fires on cost-≤1 attack actions. Mismatched target × mode resolves as a zero-Value
// no-op.
func razorReflexPlay(s *sim.TurnState, self *sim.CardState, n int) {
	target := s.AttackReactionTarget()
	if target == nil {
		return
	}
	switch self.Mode {
	case 0:
		t := target.Card.Types()
		if t.Has(card.TypeSword) && t.IsAttack() {
			sim.GrantAttackReactionBuff(s, self, n)
		}
	case 1:
		if razorReflexMode1Allowed(target.Card) {
			sim.GrantAttackReactionBuff(s, self, n)
		}
	}
}

type RazorReflexRed struct{}

func (RazorReflexRed) ID() ids.CardID                  { return ids.RazorReflexRed }
func (RazorReflexRed) Name() string                    { return "Razor Reflex" }
func (RazorReflexRed) Cost(*sim.TurnState) int         { return 1 }
func (RazorReflexRed) Pitch() int                      { return 1 }
func (RazorReflexRed) Attack() int                     { return 0 }
func (RazorReflexRed) Defense() int                    { return 2 }
func (RazorReflexRed) Types() card.TypeSet             { return razorReflexTypes }
func (RazorReflexRed) GoAgain() bool                   { return false }
func (RazorReflexRed) Modes() int                      { return 2 }
func (RazorReflexRed) ARTargetAllowed(c sim.Card) bool { return razorReflexAccepts(c) }
func (RazorReflexRed) Play(s *sim.TurnState, self *sim.CardState) {
	razorReflexPlay(s, self, 3)
}

type RazorReflexYellow struct{}

func (RazorReflexYellow) ID() ids.CardID                  { return ids.RazorReflexYellow }
func (RazorReflexYellow) Name() string                    { return "Razor Reflex" }
func (RazorReflexYellow) Cost(*sim.TurnState) int         { return 1 }
func (RazorReflexYellow) Pitch() int                      { return 2 }
func (RazorReflexYellow) Attack() int                     { return 0 }
func (RazorReflexYellow) Defense() int                    { return 2 }
func (RazorReflexYellow) Types() card.TypeSet             { return razorReflexTypes }
func (RazorReflexYellow) GoAgain() bool                   { return false }
func (RazorReflexYellow) Modes() int                      { return 2 }
func (RazorReflexYellow) ARTargetAllowed(c sim.Card) bool { return razorReflexAccepts(c) }
func (RazorReflexYellow) Play(s *sim.TurnState, self *sim.CardState) {
	razorReflexPlay(s, self, 2)
}

type RazorReflexBlue struct{}

func (RazorReflexBlue) ID() ids.CardID                  { return ids.RazorReflexBlue }
func (RazorReflexBlue) Name() string                    { return "Razor Reflex" }
func (RazorReflexBlue) Cost(*sim.TurnState) int         { return 1 }
func (RazorReflexBlue) Pitch() int                      { return 3 }
func (RazorReflexBlue) Attack() int                     { return 0 }
func (RazorReflexBlue) Defense() int                    { return 2 }
func (RazorReflexBlue) Types() card.TypeSet             { return razorReflexTypes }
func (RazorReflexBlue) GoAgain() bool                   { return false }
func (RazorReflexBlue) Modes() int                      { return 2 }
func (RazorReflexBlue) ARTargetAllowed(c sim.Card) bool { return razorReflexAccepts(c) }
func (RazorReflexBlue) Play(s *sim.TurnState, self *sim.CardState) {
	razorReflexPlay(s, self, 1)
}
