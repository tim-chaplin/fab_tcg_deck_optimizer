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

// ClubWeapon is a 1-handed Club weapon: cost 0 to swing, attack 1, no on-hit rider.
// Sized for the simplest e2e exercise of a Club-typed weapon attack target.
type ClubWeapon struct{}

func (ClubWeapon) ID() ids.WeaponID        { return FakeClubWeapon }
func (ClubWeapon) Name() string            { return "test.ClubWeapon" }
func (ClubWeapon) Cost(*sim.TurnState) int { return 0 }
func (ClubWeapon) Pitch() int              { return 0 }
func (ClubWeapon) Attack() int             { return 1 }
func (ClubWeapon) Defense() int            { return 0 }
func (ClubWeapon) Types() card.TypeSet     { return clubWeaponTypes }
func (ClubWeapon) GoAgain() bool           { return false }
func (ClubWeapon) Hands() int              { return 1 }
func (ClubWeapon) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

var hammerWeaponTypes = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeHammer, card.TypeOneHand)

// HammerWeapon is a 1-handed Hammer weapon stub mirroring ClubWeapon's shape — same cost,
// power, and lack of on-hit rider — so tests gating on Hammer have a parallel target.
type HammerWeapon struct{}

func (HammerWeapon) ID() ids.WeaponID        { return FakeHammerWeapon }
func (HammerWeapon) Name() string            { return "test.HammerWeapon" }
func (HammerWeapon) Cost(*sim.TurnState) int { return 0 }
func (HammerWeapon) Pitch() int              { return 0 }
func (HammerWeapon) Attack() int             { return 1 }
func (HammerWeapon) Defense() int            { return 0 }
func (HammerWeapon) Types() card.TypeSet     { return hammerWeaponTypes }
func (HammerWeapon) GoAgain() bool           { return false }
func (HammerWeapon) Hands() int              { return 1 }
func (HammerWeapon) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
