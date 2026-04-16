package card

// ID uniquely identifies a printed card. The zero value (Invalid) is reserved so that a zero-
// valued ID in other data structures can be detected as "unset".
//
// IDs are stable within a single build but are NOT a persistence format: adding or removing cards
// may renumber existing entries. Treat IDs as opaque in-process handles.
//
// Each pitch variant (Red / Yellow / Blue) of a printed card is a distinct card and gets its own
// ID. Weapons get IDs too so that every Card implementation has a unique, non-zero identifier.
type ID uint16

// Invalid is the sentinel zero value. Valid IDs start at 1.
const Invalid ID = 0

// Runeblade card IDs. Ordered alphabetically by card name, Red → Yellow → Blue within each family.
// Sigil of Deadwood only has a Blue variant (no R/Y printings).
const (
	AetherSlashRed ID = iota + 1
	AetherSlashYellow
	AetherSlashBlue
	AmplifyTheArknightRed
	AmplifyTheArknightYellow
	AmplifyTheArknightBlue
	ArcaneCussingRed
	ArcaneCussingYellow
	ArcaneCussingBlue
	ArcanicCrackleRed
	ArcanicCrackleYellow
	ArcanicCrackleBlue
	ArcanicSpikeRed
	ArcanicSpikeYellow
	ArcanicSpikeBlue
	BlessingOfOccultRed
	BlessingOfOccultYellow
	BlessingOfOccultBlue
	BloodspillInvocationRed
	BloodspillInvocationYellow
	BloodspillInvocationBlue
	CondemnToSlaughterRed
	CondemnToSlaughterYellow
	CondemnToSlaughterBlue
	ConsumingVolitionRed
	ConsumingVolitionYellow
	ConsumingVolitionBlue
	DeathlyDuetRed
	DeathlyDuetYellow
	DeathlyDuetBlue
	DrawnToTheDarkDimensionRed
	DrawnToTheDarkDimensionYellow
	DrawnToTheDarkDimensionBlue
	DrowningDireRed
	DrowningDireYellow
	DrowningDireBlue
	HitTheHighNotesRed
	HitTheHighNotesYellow
	HitTheHighNotesBlue
	HocusPocusRed
	HocusPocusYellow
	HocusPocusBlue
	MaleficIncantationRed
	MaleficIncantationYellow
	MaleficIncantationBlue
	MauvrionSkiesRed
	MauvrionSkiesYellow
	MauvrionSkiesBlue
	MeatAndGreetRed
	MeatAndGreetYellow
	MeatAndGreetBlue
	OathOfTheArknightRed
	OathOfTheArknightYellow
	OathOfTheArknightBlue
	ReadTheRunesRed
	ReadTheRunesYellow
	ReadTheRunesBlue
	ReduceToRunechantRed
	ReduceToRunechantYellow
	ReduceToRunechantBlue
	ReekOfCorruptionRed
	ReekOfCorruptionYellow
	ReekOfCorruptionBlue
	RuneFlashRed
	RuneFlashYellow
	RuneFlashBlue
	RunebloodIncantationRed
	RunebloodIncantationYellow
	RunebloodIncantationBlue
	RuneragerSwarmRed
	RuneragerSwarmYellow
	RuneragerSwarmBlue
	RunicFellingsongRed
	RunicFellingsongYellow
	RunicFellingsongBlue
	RunicReapingRed
	RunicReapingYellow
	RunicReapingBlue
	ShrillOfSkullformRed
	ShrillOfSkullformYellow
	ShrillOfSkullformBlue
	SigilOfDeadwoodBlue
	SigilOfSilphidaeBlue
	SigilOfSufferingRed
	SigilOfSufferingYellow
	SigilOfSufferingBlue
	SigilOfTheArknightBlue
	SingeingSteelbladeRed
	SingeingSteelbladeYellow
	SingeingSteelbladeBlue
	SkyFireLanternsRed
	SkyFireLanternsYellow
	SkyFireLanternsBlue
	SpellbladeAssaultRed
	SpellbladeAssaultYellow
	SpellbladeAssaultBlue
	SpellbladeStrikeRed
	SpellbladeStrikeYellow
	SpellbladeStrikeBlue
	SplinteringDeadwoodRed
	SplinteringDeadwoodYellow
	SplinteringDeadwoodBlue
	SutcliffesResearchNotesRed
	SutcliffesResearchNotesYellow
	SutcliffesResearchNotesBlue
	VantagePointRed
	VantagePointYellow
	VantagePointBlue
	VexingMaliceRed
	VexingMaliceYellow
	VexingMaliceBlue
	WeepingBattlegroundRed
	WeepingBattlegroundYellow
	WeepingBattlegroundBlue

	// Generic card IDs. Ordered alphabetically by card name, Red → Yellow → Blue within each family.
	DodgeBlue
	EvasiveLeapRed
	EvasiveLeapYellow
	EvasiveLeapBlue
	FateForeseenRed
	FateForeseenYellow
	FateForeseenBlue
	LayLowYellow
	PutInContextBlue
	RiseAboveRed
	RiseAboveYellow
	RiseAboveBlue
	SinkBelowRed
	SinkBelowYellow
	SinkBelowBlue
	SpringboardSomersaultYellow
	ToughenUpBlue
	UnmovableRed
	UnmovableYellow
	UnmovableBlue

	// Weapon IDs. Weapons aren't in the card registry (decks don't hold weapons) but each gets a
	// unique ID so every Card implementation has a non-zero ID.
	NebulaBladeID
	ReapingBladeID
	ScepterOfPainID

	// Test-only synthetic card IDs. Registered so that hand.Best's cache key lookup doesn't panic
	// on them. Not real FaB cards and should not be used in production decks.
	FakeRedAttack
	FakeBlueAttack
)
