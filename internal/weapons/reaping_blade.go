// Reaping Blade — Runeblade Weapon - Sword (2H). Power 3.
// Text: "Once per Turn Action - {r}: Attack. If a hero has more {h} than any other hero, they can't
// gain {h}."
//
// Simulation: modelled as an attack source costing 1 resource, dealing 3 damage. The
// health-symmetry rider is ignored (irrelevant to single-turn damage evaluation).

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var reapingBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type ReapingBlade struct{}

func (ReapingBlade) ID() ids.WeaponID        { return ids.ReapingBladeID }
func (ReapingBlade) Name() string            { return "Reaping Blade" }
func (ReapingBlade) Cost(*sim.TurnState) int { return 1 }
func (ReapingBlade) Pitch() int              { return 0 }
func (ReapingBlade) Attack() int             { return 3 }
func (ReapingBlade) Defense() int            { return 0 }
func (ReapingBlade) Types() card.TypeSet     { return reapingBladeTypes }
func (ReapingBlade) GoAgain() bool           { return false }
func (ReapingBlade) Hands() int              { return 2 }
func (ReapingBlade) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}
