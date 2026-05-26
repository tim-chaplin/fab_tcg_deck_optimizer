// Test-only Weapon fake. The card pool currently lacks Club and Hammer printings, but
// ARs like Pummel mode 0 gate on those types — ClubWeapon lets a turn-level test pin the
// predicate and the buff plumbing end-to-end without waiting on a real printing.

package testutils

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

var clubWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand)

// ClubWeapon is a 1-handed Club weapon whose swing ability deals 1, no on-hit rider.
type ClubWeapon struct{}

func (ClubWeapon) ID() ids.WeaponID    { return ids.InvalidWeapon }
func (ClubWeapon) Name() string        { return "test.ClubWeapon" }
func (ClubWeapon) DisplayName() string { return "test.ClubWeapon" }
func (ClubWeapon) Types() card.TypeSet { return clubWeaponTypes }
func (ClubWeapon) Hands() int          { return 1 }
func (ClubWeapon) Ability() card.Card  { return clubWeaponAbility }

// Cached at package init so the attack-turn runner's per-Best w.Ability() lookup is alloc-free.
// The ability is a FakeWeaponSwing with Club + OneHand sub-types and power 1.
var clubWeaponAbility card.Card = FakeWeaponSwing().
	WithName("test.ClubWeapon").
	WithPower(1).
	WithTypes(card.TypeClub, card.TypeOneHand)
