// Pummel — Generic Attack Reaction. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Choose 1; Target club or hammer weapon attack gains +N{p}. Target attack action
// card with cost 2 or more gets +N{p} and 'When this hits a hero, they discard a card.'"
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Mode 0 grants +N{p} to a club/hammer weapon attack. Mode 1 grants +N{p} to a cost-≥2
// attack action card and registers an OnHit hero-discard rider via ge.OpponentDiscard.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// pummelAccepts is the per-mode target predicate. Mode 0 gates on club/hammer weapon attack;
// mode 1 gates on cost-≥2 attack action.
func pummelAccepts(ge card.GameEngine, target *card.CardState, mode int8) bool {
	t := target.Card.Types(nil)
	switch mode {
	case 0:
		return (t.Has(card.TypeClub) || t.Has(card.TypeHammer)) && t.IsWeaponAttack()
	case 1:
		return t.IsAttackAction() && target.EffectiveCost(ge) >= 2
	}
	return false
}

// pummelOnHitDiscard fires the printed "when this hits a hero, they discard a card" rider.
func pummelOnHitDiscard(ge card.GameEngine, l card.Logger, self *card.CardState, h *card.OnHitHandler) {
	v := ge.OpponentDiscard(1)
	l.AppendPostTriggerf(self.Card.DisplayName(), v,
		"%s forced opponent to discard 1", h.Source.DisplayName())
}

// pummelPlay applies the chosen mode's effect. Mode 1 additionally registers the on-hit
// hero-discard rider on the target.
func pummelPlay(ge card.GameEngine, l card.Logger, self *card.CardState, n int) {
	target := ge.AttackReactionTarget()
	if target == nil {
		return
	}
	self.GrantAttackReactionBuff(ge, l, n)
	if self.Mode == 1 {
		target.OnHit = append(target.OnHit, card.OnHitHandler{
			Fire:   pummelOnHitDiscard,
			Source: self.Card,
		})
	}
}

func (PummelRed) Modes() int { return 2 }
func (PummelRed) ARTargetAllowed(ge card.GameEngine, target *card.CardState, mode int8) bool {
	return pummelAccepts(ge, target, mode)
}
func (PummelRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pummelPlay(ge, l, self, 4)
}

func (PummelYellow) Modes() int { return 2 }
func (PummelYellow) ARTargetAllowed(ge card.GameEngine, target *card.CardState, mode int8) bool {
	return pummelAccepts(ge, target, mode)
}
func (PummelYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pummelPlay(ge, l, self, 3)
}

func (PummelBlue) Modes() int { return 2 }
func (PummelBlue) ARTargetAllowed(ge card.GameEngine, target *card.CardState, mode int8) bool {
	return pummelAccepts(ge, target, mode)
}
func (PummelBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	pummelPlay(ge, l, self, 2)
}
