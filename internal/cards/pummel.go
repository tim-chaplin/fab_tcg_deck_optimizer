// Pummel — Generic Attack Reaction. Cost 2. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "Choose 1; Target club or hammer weapon attack gains +N{p}. Target attack action
// card with cost 2 or more gets +N{p} and 'When this hits a hero, they discard a card.'"
// (Red N=4, Yellow N=3, Blue N=2.)
//
// Mode 0 grants +N{p} to a club/hammer weapon attack. Mode 1 grants +N{p} to a cost-≥2
// attack action card and registers an OnHit hero-discard rider crediting sim.DiscardValue.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var pummelTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

// pummelAccepts is the union of mode 0 (club/hammer weapon attack) and mode 1 (cost-≥2
// attack action) target predicates. Reads Cost against an empty TurnState — variable-cost
// cards aren't expected in mode 1's gate range.
func pummelAccepts(c sim.Card) bool {
	t := c.Types()
	if (t.Has(card.TypeClub) || t.Has(card.TypeHammer)) && t.IsAttack() {
		return true
	}
	return pummelMode1Allowed(c)
}

// pummelMode1Allowed gates mode 1: cost-≥2 attack action card.
func pummelMode1Allowed(c sim.Card) bool {
	if !c.Types().IsAttackAction() {
		return false
	}
	return c.Cost(&sim.TurnState{}) >= 2
}

// pummelOnHitDiscard fires the printed "when this hits a hero, they discard a card" rider.
func pummelOnHitDiscard(s *sim.TurnState, self *sim.CardState, h *sim.OnHitHandler) {
	s.AddValue(sim.DiscardValue)
	s.LogPostTriggerf(sim.DisplayName(self.Card), sim.DiscardValue,
		"%s forced opponent to discard 1", sim.DisplayName(h.Source))
}

// pummelPlay applies the chosen mode's effect. Mode 0 grants +N{p} to a club/hammer weapon
// attack. Mode 1 grants +N{p} to a cost-≥2 attack action and registers the on-hit discard
// rider. Mismatched target × mode resolves as a zero-Value no-op.
func pummelPlay(s *sim.TurnState, self *sim.CardState, n int) {
	target := s.AttackReactionTarget()
	if target == nil {
		return
	}
	switch self.Mode {
	case 0:
		t := target.Card.Types()
		if (t.Has(card.TypeClub) || t.Has(card.TypeHammer)) && t.IsAttack() {
			sim.GrantAttackReactionBuff(s, self, n)
		}
	case 1:
		if pummelMode1Allowed(target.Card) {
			sim.GrantAttackReactionBuff(s, self, n)
			target.OnHit = append(target.OnHit, sim.OnHitHandler{
				Fire:   pummelOnHitDiscard,
				Source: self.Card,
			})
		}
	}
}

type PummelRed struct{}

func (PummelRed) ID() ids.CardID                  { return ids.PummelRed }
func (PummelRed) Name() string                    { return "Pummel" }
func (PummelRed) Cost(*sim.TurnState) int         { return 2 }
func (PummelRed) Pitch() int                      { return 1 }
func (PummelRed) Attack() int                     { return 0 }
func (PummelRed) Defense() int                    { return 2 }
func (PummelRed) Types() card.TypeSet             { return pummelTypes }
func (PummelRed) GoAgain() bool                   { return false }
func (PummelRed) Modes() int                      { return 2 }
func (PummelRed) ARTargetAllowed(c sim.Card) bool { return pummelAccepts(c) }
func (PummelRed) Play(s *sim.TurnState, self *sim.CardState) {
	pummelPlay(s, self, 4)
}

type PummelYellow struct{}

func (PummelYellow) ID() ids.CardID                  { return ids.PummelYellow }
func (PummelYellow) Name() string                    { return "Pummel" }
func (PummelYellow) Cost(*sim.TurnState) int         { return 2 }
func (PummelYellow) Pitch() int                      { return 2 }
func (PummelYellow) Attack() int                     { return 0 }
func (PummelYellow) Defense() int                    { return 2 }
func (PummelYellow) Types() card.TypeSet             { return pummelTypes }
func (PummelYellow) GoAgain() bool                   { return false }
func (PummelYellow) Modes() int                      { return 2 }
func (PummelYellow) ARTargetAllowed(c sim.Card) bool { return pummelAccepts(c) }
func (PummelYellow) Play(s *sim.TurnState, self *sim.CardState) {
	pummelPlay(s, self, 3)
}

type PummelBlue struct{}

func (PummelBlue) ID() ids.CardID                  { return ids.PummelBlue }
func (PummelBlue) Name() string                    { return "Pummel" }
func (PummelBlue) Cost(*sim.TurnState) int         { return 2 }
func (PummelBlue) Pitch() int                      { return 3 }
func (PummelBlue) Attack() int                     { return 0 }
func (PummelBlue) Defense() int                    { return 2 }
func (PummelBlue) Types() card.TypeSet             { return pummelTypes }
func (PummelBlue) GoAgain() bool                   { return false }
func (PummelBlue) Modes() int                      { return 2 }
func (PummelBlue) ARTargetAllowed(c sim.Card) bool { return pummelAccepts(c) }
func (PummelBlue) Play(s *sim.TurnState, self *sim.CardState) {
	pummelPlay(s, self, 2)
}
