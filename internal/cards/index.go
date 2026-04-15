// Package cards is the master registry of every implemented card. It assigns each printed card a
// stable unique ID and provides lookup / iteration helpers — useful for random deck generation,
// serialization, and compact equality checks.
//
// IDs are stable within a single build but are NOT a persistence format: adding or removing cards
// may renumber existing entries. Treat IDs as opaque in-process handles.
//
// Each pitch variant (Red / Yellow / Blue) of a card is a distinct printed card and gets its own
// ID. Weapons live in the weapon package and are intentionally not indexed here — decks are built
// from cards; weapons are equipment.
package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
)

// ID uniquely identifies a printed card. The zero value (Invalid) is reserved so that a zero-
// valued ID in other data structures can be detected as "unset".
type ID uint16

// Sentinel for "no card". Valid IDs start at 1.
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
	RuneFlashRed
	RuneFlashYellow
	RuneFlashBlue
	RunebloodIncantationRed
	RunebloodIncantationYellow
	RunebloodIncantationBlue
	RuneragerSwarmRed
	RuneragerSwarmYellow
	RuneragerSwarmBlue
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
	SingeingSteelbladeRed
	SingeingSteelbladeYellow
	SingeingSteelbladeBlue
	SpellbladeAssaultRed
	SpellbladeAssaultYellow
	SpellbladeAssaultBlue
	SpellbladeStrikeRed
	SpellbladeStrikeYellow
	SpellbladeStrikeBlue
	SplinteringDeadwoodRed
	SplinteringDeadwoodYellow
	SplinteringDeadwoodBlue
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

	// Test-only synthetic cards. Registered so that hand.Best's cache key lookup doesn't panic on
	// them. Not real FaB cards and should not be used in production decks.
	FakeRedAttack
	FakeBlueAttack
)

// byID is indexed directly by ID. Index 0 (Invalid) is nil.
var byID = []card.Card{
	Invalid: nil,

	AetherSlashRed:    runeblade.AetherSlashRed{},
	AetherSlashYellow: runeblade.AetherSlashYellow{},
	AetherSlashBlue:   runeblade.AetherSlashBlue{},

	AmplifyTheArknightRed:    runeblade.AmplifyTheArknightRed{},
	AmplifyTheArknightYellow: runeblade.AmplifyTheArknightYellow{},
	AmplifyTheArknightBlue:   runeblade.AmplifyTheArknightBlue{},

	ArcaneCussingRed:    runeblade.ArcaneCussingRed{},
	ArcaneCussingYellow: runeblade.ArcaneCussingYellow{},
	ArcaneCussingBlue:   runeblade.ArcaneCussingBlue{},

	ArcanicCrackleRed:    runeblade.ArcanicCrackleRed{},
	ArcanicCrackleYellow: runeblade.ArcanicCrackleYellow{},
	ArcanicCrackleBlue:   runeblade.ArcanicCrackleBlue{},

	ArcanicSpikeRed:    runeblade.ArcanicSpikeRed{},
	ArcanicSpikeYellow: runeblade.ArcanicSpikeYellow{},
	ArcanicSpikeBlue:   runeblade.ArcanicSpikeBlue{},

	BlessingOfOccultRed:    runeblade.BlessingOfOccultRed{},
	BlessingOfOccultYellow: runeblade.BlessingOfOccultYellow{},
	BlessingOfOccultBlue:   runeblade.BlessingOfOccultBlue{},

	BloodspillInvocationRed:    runeblade.BloodspillInvocationRed{},
	BloodspillInvocationYellow: runeblade.BloodspillInvocationYellow{},
	BloodspillInvocationBlue:   runeblade.BloodspillInvocationBlue{},

	CondemnToSlaughterRed:    runeblade.CondemnToSlaughterRed{},
	CondemnToSlaughterYellow: runeblade.CondemnToSlaughterYellow{},
	CondemnToSlaughterBlue:   runeblade.CondemnToSlaughterBlue{},

	ConsumingVolitionRed:    runeblade.ConsumingVolitionRed{},
	ConsumingVolitionYellow: runeblade.ConsumingVolitionYellow{},
	ConsumingVolitionBlue:   runeblade.ConsumingVolitionBlue{},

	DeathlyDuetRed:    runeblade.DeathlyDuetRed{},
	DeathlyDuetYellow: runeblade.DeathlyDuetYellow{},
	DeathlyDuetBlue:   runeblade.DeathlyDuetBlue{},

	DrawnToTheDarkDimensionRed:    runeblade.DrawnToTheDarkDimensionRed{},
	DrawnToTheDarkDimensionYellow: runeblade.DrawnToTheDarkDimensionYellow{},
	DrawnToTheDarkDimensionBlue:   runeblade.DrawnToTheDarkDimensionBlue{},

	HitTheHighNotesRed:    runeblade.HitTheHighNotesRed{},
	HitTheHighNotesYellow: runeblade.HitTheHighNotesYellow{},
	HitTheHighNotesBlue:   runeblade.HitTheHighNotesBlue{},

	HocusPocusRed:    runeblade.HocusPocusRed{},
	HocusPocusYellow: runeblade.HocusPocusYellow{},
	HocusPocusBlue:   runeblade.HocusPocusBlue{},

	MaleficIncantationRed:    runeblade.MaleficIncantationRed{},
	MaleficIncantationYellow: runeblade.MaleficIncantationYellow{},
	MaleficIncantationBlue:   runeblade.MaleficIncantationBlue{},

	MauvrionSkiesRed:    runeblade.MauvrionSkiesRed{},
	MauvrionSkiesYellow: runeblade.MauvrionSkiesYellow{},
	MauvrionSkiesBlue:   runeblade.MauvrionSkiesBlue{},

	MeatAndGreetRed:    runeblade.MeatAndGreetRed{},
	MeatAndGreetYellow: runeblade.MeatAndGreetYellow{},
	MeatAndGreetBlue:   runeblade.MeatAndGreetBlue{},

	OathOfTheArknightRed:    runeblade.OathOfTheArknightRed{},
	OathOfTheArknightYellow: runeblade.OathOfTheArknightYellow{},
	OathOfTheArknightBlue:   runeblade.OathOfTheArknightBlue{},

	ReadTheRunesRed:    runeblade.ReadTheRunesRed{},
	ReadTheRunesYellow: runeblade.ReadTheRunesYellow{},
	ReadTheRunesBlue:   runeblade.ReadTheRunesBlue{},

	ReduceToRunechantRed:    runeblade.ReduceToRunechantRed{},
	ReduceToRunechantYellow: runeblade.ReduceToRunechantYellow{},
	ReduceToRunechantBlue:   runeblade.ReduceToRunechantBlue{},

	RuneFlashRed:    runeblade.RuneFlashRed{},
	RuneFlashYellow: runeblade.RuneFlashYellow{},
	RuneFlashBlue:   runeblade.RuneFlashBlue{},

	RunebloodIncantationRed:    runeblade.RunebloodIncantationRed{},
	RunebloodIncantationYellow: runeblade.RunebloodIncantationYellow{},
	RunebloodIncantationBlue:   runeblade.RunebloodIncantationBlue{},

	RuneragerSwarmRed:    runeblade.RuneragerSwarmRed{},
	RuneragerSwarmYellow: runeblade.RuneragerSwarmYellow{},
	RuneragerSwarmBlue:   runeblade.RuneragerSwarmBlue{},

	RunicReapingRed:    runeblade.RunicReapingRed{},
	RunicReapingYellow: runeblade.RunicReapingYellow{},
	RunicReapingBlue:   runeblade.RunicReapingBlue{},

	ShrillOfSkullformRed:    runeblade.ShrillOfSkullformRed{},
	ShrillOfSkullformYellow: runeblade.ShrillOfSkullformYellow{},
	ShrillOfSkullformBlue:   runeblade.ShrillOfSkullformBlue{},

	SigilOfDeadwoodBlue: runeblade.SigilOfDeadwoodBlue{},

	SigilOfSilphidaeBlue: runeblade.SigilOfSilphidaeBlue{},

	SigilOfSufferingRed:    runeblade.SigilOfSufferingRed{},
	SigilOfSufferingYellow: runeblade.SigilOfSufferingYellow{},
	SigilOfSufferingBlue:   runeblade.SigilOfSufferingBlue{},

	SingeingSteelbladeRed:    runeblade.SingeingSteelbladeRed{},
	SingeingSteelbladeYellow: runeblade.SingeingSteelbladeYellow{},
	SingeingSteelbladeBlue:   runeblade.SingeingSteelbladeBlue{},

	SpellbladeAssaultRed:    runeblade.SpellbladeAssaultRed{},
	SpellbladeAssaultYellow: runeblade.SpellbladeAssaultYellow{},
	SpellbladeAssaultBlue:   runeblade.SpellbladeAssaultBlue{},

	SpellbladeStrikeRed:    runeblade.SpellbladeStrikeRed{},
	SpellbladeStrikeYellow: runeblade.SpellbladeStrikeYellow{},
	SpellbladeStrikeBlue:   runeblade.SpellbladeStrikeBlue{},

	SplinteringDeadwoodRed:    runeblade.SplinteringDeadwoodRed{},
	SplinteringDeadwoodYellow: runeblade.SplinteringDeadwoodYellow{},
	SplinteringDeadwoodBlue:   runeblade.SplinteringDeadwoodBlue{},

	VantagePointRed:    runeblade.VantagePointRed{},
	VantagePointYellow: runeblade.VantagePointYellow{},
	VantagePointBlue:   runeblade.VantagePointBlue{},

	VexingMaliceRed:    runeblade.VexingMaliceRed{},
	VexingMaliceYellow: runeblade.VexingMaliceYellow{},
	VexingMaliceBlue:   runeblade.VexingMaliceBlue{},

	WeepingBattlegroundRed:    runeblade.WeepingBattlegroundRed{},
	WeepingBattlegroundYellow: runeblade.WeepingBattlegroundYellow{},
	WeepingBattlegroundBlue:   runeblade.WeepingBattlegroundBlue{},

	DodgeBlue: generic.DodgeBlue{},

	EvasiveLeapRed:    generic.EvasiveLeapRed{},
	EvasiveLeapYellow: generic.EvasiveLeapYellow{},
	EvasiveLeapBlue:   generic.EvasiveLeapBlue{},

	FateForeseenRed:    generic.FateForeseenRed{},
	FateForeseenYellow: generic.FateForeseenYellow{},
	FateForeseenBlue:   generic.FateForeseenBlue{},

	LayLowYellow: generic.LayLowYellow{},

	PutInContextBlue: generic.PutInContextBlue{},

	RiseAboveRed:    generic.RiseAboveRed{},
	RiseAboveYellow: generic.RiseAboveYellow{},
	RiseAboveBlue:   generic.RiseAboveBlue{},

	SinkBelowRed:    generic.SinkBelowRed{},
	SinkBelowYellow: generic.SinkBelowYellow{},
	SinkBelowBlue:   generic.SinkBelowBlue{},

	SpringboardSomersaultYellow: generic.SpringboardSomersaultYellow{},

	ToughenUpBlue: generic.ToughenUpBlue{},

	UnmovableRed:    generic.UnmovableRed{},
	UnmovableYellow: generic.UnmovableYellow{},
	UnmovableBlue:   generic.UnmovableBlue{},

	FakeRedAttack:  fake.RedAttack{},
	FakeBlueAttack: fake.BlueAttack{},
}

// byName maps Card.Name() → ID for reverse lookup. Built once at init.
var byName = func() map[string]ID {
	m := make(map[string]ID, len(byID)-1)
	for id, c := range byID {
		if c == nil {
			continue
		}
		m[c.Name()] = ID(id)
	}
	return m
}()

// Get returns the card for the given ID. Panics if id is Invalid or out of range — callers should
// only pass IDs they got from this package.
func Get(id ID) card.Card {
	if id == Invalid || int(id) >= len(byID) {
		panic("cardindex: invalid card ID")
	}
	return byID[id]
}

// ByName looks up an ID by the card's Name(). Returns (Invalid, false) if no such card.
func ByName(name string) (ID, bool) {
	id, ok := byName[name]
	return id, ok
}

// All returns every valid card ID in registration order. The returned slice is freshly allocated
// and safe for the caller to mutate (e.g. shuffle for random deck generation).
func All() []ID {
	out := make([]ID, 0, len(byID)-1)
	for id := 1; id < len(byID); id++ {
		out = append(out, ID(id))
	}
	return out
}

// Count is the number of registered cards (excluding Invalid).
func Count() int { return len(byID) - 1 }
