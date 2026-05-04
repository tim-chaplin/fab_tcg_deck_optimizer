// Test-only Weapon stubs. The card pool currently lacks Club and Hammer weapons, but ARs
// like Pummel mode 0 gate on those types — these stubs let e2e tests pin the predicate
// and the buff plumbing end-to-end without waiting on a real Club/Hammer printing.

package testutils

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var clubWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand)

// ClubWeapon is a 1-handed Club weapon: swing cost 0, attack 1, no on-hit rider.
type ClubWeapon struct{}

func (ClubWeapon) ID() ids.WeaponID    { return FakeClubWeapon }
func (ClubWeapon) Name() string        { return "test.ClubWeapon" }
func (ClubWeapon) Types() card.TypeSet { return clubWeaponTypes }
func (ClubWeapon) Hands() int          { return 1 }
func (ClubWeapon) Ability() sim.Card   { return ClubWeaponAbility{} }

var clubWeaponAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeClub, card.TypeOneHand, card.TypeAttack)

// ClubWeaponAbility is the activated-ability Card for ClubWeapon: cost 0, power 1, no rider.
type ClubWeaponAbility struct{}

func (ClubWeaponAbility) ID() ids.CardID          { return FakeClubWeaponAbility }
func (ClubWeaponAbility) Name() string            { return "test.ClubWeapon" }
func (ClubWeaponAbility) Cost(*sim.TurnState) int { return 0 }
func (ClubWeaponAbility) Pitch() int              { return 0 }
func (ClubWeaponAbility) Attack() int             { return 1 }
func (ClubWeaponAbility) Defense() int            { return 0 }
func (ClubWeaponAbility) Types() card.TypeSet     { return clubWeaponAbilityTypes }
func (ClubWeaponAbility) GoAgain() bool           { return false }
func (ClubWeaponAbility) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

var hammerWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeHammer, card.TypeOneHand)

// HammerWeapon is a 1-handed Hammer weapon mirroring ClubWeapon's shape so Hammer-gated
// tests have a parallel target.
type HammerWeapon struct{}

func (HammerWeapon) ID() ids.WeaponID    { return FakeHammerWeapon }
func (HammerWeapon) Name() string        { return "test.HammerWeapon" }
func (HammerWeapon) Types() card.TypeSet { return hammerWeaponTypes }
func (HammerWeapon) Hands() int          { return 1 }
func (HammerWeapon) Ability() sim.Card   { return HammerWeaponAbility{} }

var hammerWeaponAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeHammer, card.TypeOneHand, card.TypeAttack)

// HammerWeaponAbility is the activated-ability Card for HammerWeapon: cost 0, power 1, no rider.
type HammerWeaponAbility struct{}

func (HammerWeaponAbility) ID() ids.CardID          { return FakeHammerWeaponAbility }
func (HammerWeaponAbility) Name() string            { return "test.HammerWeapon" }
func (HammerWeaponAbility) Cost(*sim.TurnState) int { return 0 }
func (HammerWeaponAbility) Pitch() int              { return 0 }
func (HammerWeaponAbility) Attack() int             { return 1 }
func (HammerWeaponAbility) Defense() int            { return 0 }
func (HammerWeaponAbility) Types() card.TypeSet     { return hammerWeaponAbilityTypes }
func (HammerWeaponAbility) GoAgain() bool           { return false }
func (HammerWeaponAbility) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
