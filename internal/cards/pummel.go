// Pummel — Generic Attack Reaction. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Choose 1; Target club or hammer weapon attack gains +N{p}. Target attack action
// card with cost 2 or more gets +N{p} and 'When this hits a hero, they discard a card.'"
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Mode 0 grants +N{p} to a club/hammer weapon attack. Mode 1 grants +N{p} to a cost-≥2
// attack action card and registers an OnHit hero-discard rider via g.OpponentDiscard.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// pummelAccepts is the per-mode target predicate. Mode 0 gates on club/hammer weapon
// attack; mode 1 gates on cost-≥2 attack action. The chain runner runs this for the
// chosen Mode and rejects the permutation when it returns false, so pummelPlay can apply
// the buff unconditionally.
//
// Reads Cost against an empty TurnState; variable-cost cards aren't expected in mode 1's
// gate range.
func pummelAccepts(c card.Card, mode int8) bool {
	t := c.Types(nil)
	switch mode {
	case 0:
		return (t.Has(card.TypeClub) || t.Has(card.TypeHammer)) && t.IsWeaponAttack()
	case 1:
		return t.IsAttackAction() && c.Cost(&sim.TurnState{}) >= 2
	}
	return false
}

// pummelOnHitDiscard fires the printed "when this hits a hero, they discard a card" rider.
func pummelOnHitDiscard(s card.GameEngine, l card.Logger, self *card.CardState, h *card.OnHitHandler) {
	v := s.OpponentDiscard(1)
	l.AppendPostTriggerf(self.Card.DisplayName(), v,
		"%s forced opponent to discard 1", h.Source.DisplayName())
}

// pummelPlay applies the chosen mode's effect. The chain runner already validated the
// target via pummelAccepts, so the buff lands directly. Mode 1 additionally registers the
// on-hit hero-discard rider on the target.
func pummelPlay(s card.GameEngine, l card.Logger, self *card.CardState, n int) {
	target := s.AttackReactionTarget()
	if target == nil {
		return
	}
	card.GrantAttackReactionBuff(s, l, self, n)
	if self.Mode == 1 {
		target.OnHit = append(target.OnHit, card.OnHitHandler{
			Fire:   pummelOnHitDiscard,
			Source: self.Card,
		})
	}
}

func (PummelRed) Modes() int { return 2 }
func (PummelRed) ARTargetAllowed(c card.Card, mode int8) bool {
	return pummelAccepts(c, mode)
}
func (PummelRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	pummelPlay(s, l, self, 4)
}

func (PummelYellow) Modes() int { return 2 }
func (PummelYellow) ARTargetAllowed(c card.Card, mode int8) bool {
	return pummelAccepts(c, mode)
}
func (PummelYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	pummelPlay(s, l, self, 3)
}

func (PummelBlue) Modes() int { return 2 }
func (PummelBlue) ARTargetAllowed(c card.Card, mode int8) bool {
	return pummelAccepts(c, mode)
}
func (PummelBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	pummelPlay(s, l, self, 2)
}
