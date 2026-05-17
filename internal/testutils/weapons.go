// Test-only Weapon and weapon-ability stubs. The card pool currently lacks Club and Hammer
// printings, but ARs like Pummel mode 0 gate on those types — these stubs let turntests
// pin the predicate and the buff plumbing end-to-end without waiting on a real printing.

package testutils

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

var clubWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand)

// ClubWeapon is a 1-handed Club weapon: swing cost 0, attack 1, no on-hit rider.
type ClubWeapon struct{}

func (ClubWeapon) ID() ids.WeaponID    { return FakeClubWeapon }
func (ClubWeapon) Name() string        { return "test.ClubWeapon" }
func (ClubWeapon) DisplayName() string { return "test.ClubWeapon" }
func (ClubWeapon) Types() card.TypeSet { return clubWeaponTypes }
func (ClubWeapon) Hands() int          { return 1 }
func (ClubWeapon) Ability() card.Card  { return clubWeaponAbility }

// Cached at package init so the chain runner's per-Best w.Ability() lookup is alloc-free.
var clubWeaponAbility card.Card = ClubWeaponAbility{}

var clubWeaponAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand, card.TypeAttack)

// ClubWeaponAbility is the activated-ability Card for ClubWeapon: cost 0, power 1, no rider.
type ClubWeaponAbility struct{}

func (ClubWeaponAbility) ID() ids.CardID                     { return FakeClubWeaponAbility }
func (ClubWeaponAbility) Name() string                       { return "test.ClubWeapon" }
func (ClubWeaponAbility) DisplayName() string                { return "test.ClubWeapon" }
func (ClubWeaponAbility) Cost(card.GameEngine) int           { return 0 }
func (ClubWeaponAbility) Pitch() int                         { return 0 }
func (ClubWeaponAbility) Attack() int                        { return 1 }
func (ClubWeaponAbility) Defense() int                       { return 0 }
func (ClubWeaponAbility) Types(card.GameEngine) card.TypeSet { return clubWeaponAbilityTypes }
func (ClubWeaponAbility) GoAgain(card.GameEngine) bool       { return false }
func (ClubWeaponAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}

var hammerWeaponAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeHammer, card.TypeOneHand, card.TypeAttack)

// HammerWeaponAbility is a stand-alone activated-ability Card with the Hammer + Weapon type
// line so Hammer-gated AR predicates (e.g. Pummel mode 0) can be exercised without a real
// weapon printing to back it.
type HammerWeaponAbility struct{}

func (HammerWeaponAbility) ID() ids.CardID                     { return FakeHammerWeaponAbility }
func (HammerWeaponAbility) Name() string                       { return "test.HammerWeapon" }
func (HammerWeaponAbility) DisplayName() string                { return "test.HammerWeapon" }
func (HammerWeaponAbility) Cost(card.GameEngine) int           { return 0 }
func (HammerWeaponAbility) Pitch() int                         { return 0 }
func (HammerWeaponAbility) Attack() int                        { return 1 }
func (HammerWeaponAbility) Defense() int                       { return 0 }
func (HammerWeaponAbility) Types(card.GameEngine) card.TypeSet { return hammerWeaponAbilityTypes }
func (HammerWeaponAbility) GoAgain(card.GameEngine) bool       { return false }
func (HammerWeaponAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
}
