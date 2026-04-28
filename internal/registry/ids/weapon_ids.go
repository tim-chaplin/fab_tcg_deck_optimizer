package ids

// WeaponID identifies a weapon printing.
//
// KLUDGE: aliased to CardID and anchored after the last fake card ID. Ideally weapons
// would have their own number space (starting at 1), but every weapon swing flows through
// the same chain-runner pipeline as deck cards: Weapon's methods structurally satisfy
// sim.Card so weapons can sit alongside cards in the chain's permutation slice without an
// adapter, and the per-card caches (chain step text, display name, attacker meta) are
// keyed by ids.CardID. With distinct types we'd need either a sim.Card-shaped wrapper or
// a deeper refactor that branches the chain runner per slot kind. See TODO.md → "Weapon
// IDs share the CardID space".
type WeaponID = CardID

// InvalidWeapon is the sentinel zero value.
const InvalidWeapon WeaponID = 0

// Weapon IDs. Anchored after the last fake card so weapons don't share cache slots with
// cards in the shared CardID space.
const (
	AnnalsOfSutcliffeID WeaponID = FakeHugeAttack + iota + 1
	NebulaBladeID
	ReapingBladeID
	RosettaThornID
	ScepterOfPainID
	TalisharID
)
