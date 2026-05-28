package ids

// Weapon-permanent card IDs. Weapons are cards: the platonic weapon card lives in the card
// universe, goes to the graveyard when its equipped object is destroyed, and is counted by
// card conservation — so it carries a CardID like any other card. Anchored past the last
// real deck card so weapons get distinct cache slots in the shared CardID space.
const (
	AnnalsOfSutcliffeID CardID = ZealousBeltingBlue + iota + 1
	NebulaBladeID
	ReapingBladeID
	RosettaThornID
	ScepterOfPainID
	TalisharID
)

// Weapon ability IDs. Anchored after the weapon-permanent IDs in the shared CardID space;
// cardMetaCache keys off the ability ID since the attack-turn runner enqueues the ability.
const (
	AnnalsOfSutcliffeAbilityID CardID = TalisharID + iota + 1
	NebulaBladeAbilityID
	ReapingBladeAbilityID
	RosettaThornAbilityID
	ScepterOfPainAbilityID
	TalisharAbilityID
)
