// Test-only Weapon fake. The card pool currently lacks Club and Hammer printings, but
// ARs like Pummel mode 0 gate on those types — ClubWeapon lets a turn-level test pin the
// predicate and the buff plumbing end-to-end without waiting on a real printing.

package testutils

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

var clubWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand)

// ClubWeapon is a 1-handed Club weapon card whose swing ability deals 1, no on-hit rider.
// As a platonic weapon card it satisfies weapon.Card (card.Card + Hands + Ability); the
// equip-time builder wraps it in a mutable weapon.Weapon object.
type ClubWeapon struct{}

func (ClubWeapon) ID() ids.CardID                                     { return ids.InvalidCard }
func (ClubWeapon) Name() string                                       { return "test.ClubWeapon" }
func (ClubWeapon) DisplayName() string                                { return "test.ClubWeapon" }
func (ClubWeapon) Cost() int                                          { return 0 }
func (ClubWeapon) Pitch() int                                         { return 0 }
func (ClubWeapon) Attack() int                                        { return 0 }
func (ClubWeapon) Defense() int                                       { return 0 }
func (ClubWeapon) Types(card.GameEngine) card.TypeSet                 { return clubWeaponTypes }
func (ClubWeapon) GoAgain(card.GameEngine) bool                       { return false }
func (ClubWeapon) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (ClubWeapon) Hands() int                                         { return 1 }
func (ClubWeapon) Ability() card.Card                                 { return clubWeaponAbility }

func (ClubWeapon) WeaponTrigger() (triggertype.Type, weapon.Handler) { return 0, nil }

// Cached at package init so the attack-turn runner's per-Best w.Ability() lookup is alloc-free.
// The ability is a FakeWeaponSwing with Club + OneHand sub-types and power 1.
var clubWeaponAbility card.Card = FakeWeaponSwing().
	WithName("test.ClubWeapon").
	WithPower(1).
	WithTypes(card.TypeClub, card.TypeOneHand)

var counterWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand)

// CounterWeapon is a 1-handed weapon whose swing ability puts a counter on its own equipped
// object via self.Weapon. The swing has Go again so both copies of a dual-wield loadout can
// swing in one turn (one starting action point + the first swing's go-again funds the
// second). Two equipped copies share one Ability() card value (like any pair of identical
// weapons), so it pins per-object counter attribution: each swing must bump the specific
// object the runner attributed it to. As a platonic weapon card it satisfies weapon.Card;
// the equip-time builder wraps each copy in its own mutable weapon.Weapon object.
type CounterWeapon struct{}

func (CounterWeapon) ID() ids.CardID                                     { return ids.InvalidCard }
func (CounterWeapon) Name() string                                       { return "test.CounterWeapon" }
func (CounterWeapon) DisplayName() string                                { return "test.CounterWeapon" }
func (CounterWeapon) Cost() int                                          { return 0 }
func (CounterWeapon) Pitch() int                                         { return 0 }
func (CounterWeapon) Attack() int                                        { return 0 }
func (CounterWeapon) Defense() int                                       { return 0 }
func (CounterWeapon) Types(card.GameEngine) card.TypeSet                 { return counterWeaponTypes }
func (CounterWeapon) GoAgain(card.GameEngine) bool                       { return false }
func (CounterWeapon) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (CounterWeapon) Hands() int                                         { return 1 }
func (CounterWeapon) Ability() card.Card                                 { return counterWeaponAbility }

func (CounterWeapon) WeaponTrigger() (triggertype.Type, weapon.Handler) { return 0, nil }

// counterWeaponAbility is the shared swing ability: power 1, Go again, whose Play bumps the
// counter on self.Weapon — the specific equipped object the attack-turn runner resolved for
// this swing. Go again lets both copies of a dual-wield loadout swing in a single turn.
var counterWeaponAbility card.Card = FakeWeaponSwing().
	WithName("test.CounterWeapon").
	WithPower(1).
	WithGoAgain().
	WithTypes(card.TypeClub, card.TypeOneHand).
	WithPlay(func(_ card.GameEngine, _ card.Logger, self *card.CardState) {
		if self.Weapon != nil {
			self.Weapon.SetCount(self.Weapon.Count() + 1)
		}
	})
