// Rosetta Thorn — Runeblade Weapon - Sword (2H). Power 2, Arcane 2.
//
// Text: "**Once per Turn Action** - {r}: **Attack**
//
// Whenever you attack with Rosetta Thorn, if you've played an attack action card and a
// 'non-attack' action card this turn, deal 2 arcane damage to target hero."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

var rosettaThornTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type RosettaThorn struct{}

func (RosettaThorn) ID() ids.WeaponID    { return ids.RosettaThornID }
func (RosettaThorn) Name() string        { return "Rosetta Thorn" }
func (RosettaThorn) Types() card.TypeSet { return rosettaThornTypes }
func (RosettaThorn) Hands() int          { return 2 }
func (RosettaThorn) Ability() card.Card  { return rosettaThornAbility }

// Cached at package init — see nebula_blade.go for the alloc-free rationale.
var rosettaThornAbility card.Card = RosettaThornAbility{}

func (RosettaThorn) NotSilverAgeLegal() {}

// not implemented: on-attack 2 arcane damage rider gated on having played an attack action AND
// a non-attack action this turn
func (RosettaThorn) NotImplemented() {}

var rosettaThornAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type RosettaThornAbility struct{}

func (RosettaThornAbility) ID() ids.CardID                     { return ids.RosettaThornAbilityID }
func (RosettaThornAbility) Name() string                       { return "Rosetta Thorn" }
func (RosettaThornAbility) DisplayName() string                { return "Rosetta Thorn" }
func (RosettaThornAbility) Cost(card.GameEngine) int           { return 1 }
func (RosettaThornAbility) Pitch() int                         { return 0 }
func (RosettaThornAbility) Attack() int                        { return 2 }
func (RosettaThornAbility) Defense() int                       { return 0 }
func (RosettaThornAbility) Types(card.GameEngine) card.TypeSet { return rosettaThornAbilityTypes }
func (RosettaThornAbility) GoAgain() bool                      { return false }
func (RosettaThornAbility) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
}
