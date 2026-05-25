package ids

// WeaponID identifies a weapon printing.
//
// KLUDGE: aliased to CardID. Weapon swings flow through the same chain-runner pipeline as
// deck cards — Weapon structurally satisfies sim.Card so weapons sit alongside cards in
// the chain's permutation slice without an adapter, and per-card caches key by CardID.
// See TODO.md → "Weapon IDs share the CardID space".
type WeaponID = CardID

// InvalidWeapon is the sentinel zero value.
const InvalidWeapon WeaponID = 0

// Weapon IDs. Anchored past the last real card so weapons get distinct cache slots in the
// shared CardID space. Test fakes return InvalidCard / InvalidWeapon; per-ID caches bypass
// the Invalid slot rather than consume a production ID.
const (
	AnnalsOfSutcliffeID WeaponID = ZealousBeltingBlue + iota + 1
	NebulaBladeID
	ReapingBladeID
	RosettaThornID
	ScepterOfPainID
	TalisharID
)

// Weapon ability IDs. Anchored after the weapon-permanent IDs in the shared CardID space;
// cardMetaCache keys off the ability ID since the chain runner enqueues the ability.
const (
	AnnalsOfSutcliffeAbilityID CardID = TalisharID + iota + 1
	NebulaBladeAbilityID
	ReapingBladeAbilityID
	RosettaThornAbilityID
	ScepterOfPainAbilityID
	TalisharAbilityID
)
