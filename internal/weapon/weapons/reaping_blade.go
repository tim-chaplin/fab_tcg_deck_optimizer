// Reaping Blade — Runeblade Weapon - Sword (2H). Power 3.
// Text: "Once per Turn Action - {r}: Attack. If a hero has more {h} than any other hero, they can't
// gain {h}."
//
// Simulation: modelled as an attack source costing 1 resource, dealing 3 damage. The
// health-symmetry rider is ignored (irrelevant to single-turn damage evaluation).

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

var reapingBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

// ReapingBlade is the platonic weapon card — the equipped permanent. Its activated ability
// is what the attack-turn runner enqueues each turn.
type ReapingBlade struct{}

func (ReapingBlade) ID() ids.CardID                                     { return ids.ReapingBladeID }
func (ReapingBlade) Name() string                                       { return "Reaping Blade" }
func (ReapingBlade) DisplayName() string                                { return "Reaping Blade" }
func (ReapingBlade) Cost() int                                          { return 0 }
func (ReapingBlade) Pitch() int                                         { return 0 }
func (ReapingBlade) Attack() int                                        { return 0 }
func (ReapingBlade) Defense() int                                       { return 0 }
func (ReapingBlade) Types(card.GameEngine) card.TypeSet                 { return reapingBladeTypes }
func (ReapingBlade) GoAgain(card.GameEngine) bool                       { return false }
func (ReapingBlade) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (ReapingBlade) Hands() int                                         { return 2 }
func (ReapingBlade) Ability() card.Card                                 { return reapingBladeAbility }

var reapingBladeAbility card.Card = ReapingBladeAbility{}

var reapingBladeAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type ReapingBladeAbility struct{}

func (ReapingBladeAbility) ID() ids.CardID                     { return ids.ReapingBladeAbilityID }
func (ReapingBladeAbility) Name() string                       { return "Reaping Blade" }
func (ReapingBladeAbility) DisplayName() string                { return "Reaping Blade" }
func (ReapingBladeAbility) Cost() int                          { return 1 }
func (ReapingBladeAbility) Pitch() int                         { return 0 }
func (ReapingBladeAbility) Attack() int                        { return 3 }
func (ReapingBladeAbility) Defense() int                       { return 0 }
func (ReapingBladeAbility) Types(card.GameEngine) card.TypeSet { return reapingBladeAbilityTypes }
func (ReapingBladeAbility) GoAgain(card.GameEngine) bool       { return false }
func (ReapingBladeAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
