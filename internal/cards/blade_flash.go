// Blade Flash — Generic Attack Reaction. Cost 1. Printed pitch variants: Blue 3. Defense 2.
//
// Text: "Target sword attack gains **go again**."
//
// Predicate is "sword attack" (no "action card" qualifier), so Sword weapons qualify too.
// The go-again grant is modelled by bumping ActionPoints eagerly at Play time.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var bladeFlashTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAttackReaction)

type BladeFlashBlue struct{}

func (BladeFlashBlue) ID() ids.CardID          { return ids.BladeFlashBlue }
func (BladeFlashBlue) Name() string            { return "Blade Flash" }
func (BladeFlashBlue) Cost(*sim.TurnState) int { return 1 }
func (BladeFlashBlue) Pitch() int              { return 3 }
func (BladeFlashBlue) Attack() int             { return 0 }
func (BladeFlashBlue) Defense() int            { return 2 }
func (BladeFlashBlue) Types() card.TypeSet     { return bladeFlashTypes }
func (BladeFlashBlue) GoAgain() bool           { return false }
func (BladeFlashBlue) ARTargetAllowed(c sim.Card, _ int8) bool {
	t := c.Types()
	return t.Has(card.TypeSword) && t.IsAttack()
}
func (BladeFlashBlue) Play(s *sim.TurnState, _ *sim.CardState) {
	if s.AttackReactionTarget() == nil {
		return
	}
	s.ActionPoints++
}
