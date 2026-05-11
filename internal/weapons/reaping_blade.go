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

func (ReapingBlade) ID() ids.WeaponID    { return ids.ReapingBladeID }
func (ReapingBlade) Name() string        { return "Reaping Blade" }
func (ReapingBlade) Types() card.TypeSet { return reapingBladeTypes }
func (ReapingBlade) Hands() int          { return 2 }
func (ReapingBlade) Ability() sim.Card   { return reapingBladeAbility }

// Cached at package init so Weapon.Ability() returns a stable interface value (no
// per-call re-box of the zero-size struct on the chain runner's hot path).
var reapingBladeAbility sim.Card = ReapingBladeAbility{}

var reapingBladeAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type ReapingBladeAbility struct{}

func (ReapingBladeAbility) ID() ids.CardID          { return ids.ReapingBladeAbilityID }
func (ReapingBladeAbility) Name() string            { return "Reaping Blade" }
func (ReapingBladeAbility) DisplayName() string     { return "Reaping Blade" }
func (ReapingBladeAbility) Cost(sim.GameEngine) int { return 1 }
func (ReapingBladeAbility) Pitch() int              { return 0 }
func (ReapingBladeAbility) Attack() int             { return 3 }
func (ReapingBladeAbility) Defense() int            { return 0 }
func (ReapingBladeAbility) Types() card.TypeSet     { return reapingBladeAbilityTypes }
func (ReapingBladeAbility) GoAgain() bool           { return false }
func (ReapingBladeAbility) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
}
