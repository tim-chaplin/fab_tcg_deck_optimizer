package testutils

// Test-only synthetic card IDs. Anchored past every real card and weapon ID so
// cardMetaCache (keyed by ids.CardID) gets a distinct slot per test stub. Production
// code never sees these — they exist only so test fakes don't share a cache slot with
// real cards or with each other.

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// FakeRedAttack and friends are the Fake-prefix synthetic IDs the StubCard, GenericAttack,
// and named-fixture fakes (BluePitch, BlueAttack, RedAttack, …) hand back from ID().
// Anchored past the last weapon-ability ID so fakes don't collide with any production
// printing.
const (
	FakeRedAttack ids.CardID = ids.TalisharAbilityID + iota + 1
	FakeBlueAttack
	FakeYellowAttack
	FakeCostlyDraw
	FakeCostlyAttack
	FakePitchOneDR
	FakeHugeAttack
	FakeBluePitch
	FakeInstant
	FakeNoGoAgainAttack
	FakeClubWeapon
	FakeHammerWeapon
	FakeRedPitch
	FakeClubWeaponAbility
	FakeHammerWeaponAbility
)
