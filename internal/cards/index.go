// Package cards is the master registry of every implemented card. The canonical ID type and
// constants live in package card; this package maps IDs to concrete Card values and provides
// lookup / iteration helpers useful for random deck generation, serialization, and compact
// equality checks.
//
// Weapons aren't ID-indexed here (decks are built from cards; weapons are equipment) but the full
// roster is exposed via AllWeapons for convenience.
package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/fake"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/runeblade"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// AllWeapons lists every implemented weapon. Used by deck-search code to enumerate loadouts.
var AllWeapons = []weapon.Weapon{
	weapon.NebulaBlade{},
	weapon.ReapingBlade{},
	weapon.ScepterOfPain{},
}

// ID aliases card.ID so callers of this package don't need two imports just to hold IDs.
type ID = card.ID

// Invalid aliases card.Invalid — the sentinel zero value.
const Invalid = card.Invalid

// byID is indexed directly by ID. Index 0 (Invalid) is nil.
var byID = []card.Card{
	card.Invalid: nil,

	card.AetherSlashRed:    runeblade.AetherSlashRed{},
	card.AetherSlashYellow: runeblade.AetherSlashYellow{},
	card.AetherSlashBlue:   runeblade.AetherSlashBlue{},

	card.AmplifyTheArknightRed:    runeblade.AmplifyTheArknightRed{},
	card.AmplifyTheArknightYellow: runeblade.AmplifyTheArknightYellow{},
	card.AmplifyTheArknightBlue:   runeblade.AmplifyTheArknightBlue{},

	card.ArcaneCussingRed:    runeblade.ArcaneCussingRed{},
	card.ArcaneCussingYellow: runeblade.ArcaneCussingYellow{},
	card.ArcaneCussingBlue:   runeblade.ArcaneCussingBlue{},

	card.ArcanicCrackleRed:    runeblade.ArcanicCrackleRed{},
	card.ArcanicCrackleYellow: runeblade.ArcanicCrackleYellow{},
	card.ArcanicCrackleBlue:   runeblade.ArcanicCrackleBlue{},

	card.ArcanicSpikeRed:    runeblade.ArcanicSpikeRed{},
	card.ArcanicSpikeYellow: runeblade.ArcanicSpikeYellow{},
	card.ArcanicSpikeBlue:   runeblade.ArcanicSpikeBlue{},

	card.BlessingOfOccultRed:    runeblade.BlessingOfOccultRed{},
	card.BlessingOfOccultYellow: runeblade.BlessingOfOccultYellow{},
	card.BlessingOfOccultBlue:   runeblade.BlessingOfOccultBlue{},

	card.BloodspillInvocationRed:    runeblade.BloodspillInvocationRed{},
	card.BloodspillInvocationYellow: runeblade.BloodspillInvocationYellow{},
	card.BloodspillInvocationBlue:   runeblade.BloodspillInvocationBlue{},

	card.CondemnToSlaughterRed:    runeblade.CondemnToSlaughterRed{},
	card.CondemnToSlaughterYellow: runeblade.CondemnToSlaughterYellow{},
	card.CondemnToSlaughterBlue:   runeblade.CondemnToSlaughterBlue{},

	card.ConsumingVolitionRed:    runeblade.ConsumingVolitionRed{},
	card.ConsumingVolitionYellow: runeblade.ConsumingVolitionYellow{},
	card.ConsumingVolitionBlue:   runeblade.ConsumingVolitionBlue{},

	card.DeathlyDuetRed:    runeblade.DeathlyDuetRed{},
	card.DeathlyDuetYellow: runeblade.DeathlyDuetYellow{},
	card.DeathlyDuetBlue:   runeblade.DeathlyDuetBlue{},

	card.DrawnToTheDarkDimensionRed:    runeblade.DrawnToTheDarkDimensionRed{},
	card.DrawnToTheDarkDimensionYellow: runeblade.DrawnToTheDarkDimensionYellow{},
	card.DrawnToTheDarkDimensionBlue:   runeblade.DrawnToTheDarkDimensionBlue{},

	card.DrowningDireRed:    runeblade.DrowningDireRed{},
	card.DrowningDireYellow: runeblade.DrowningDireYellow{},
	card.DrowningDireBlue:   runeblade.DrowningDireBlue{},

	card.HitTheHighNotesRed:    runeblade.HitTheHighNotesRed{},
	card.HitTheHighNotesYellow: runeblade.HitTheHighNotesYellow{},
	card.HitTheHighNotesBlue:   runeblade.HitTheHighNotesBlue{},

	card.HocusPocusRed:    runeblade.HocusPocusRed{},
	card.HocusPocusYellow: runeblade.HocusPocusYellow{},
	card.HocusPocusBlue:   runeblade.HocusPocusBlue{},

	card.MaleficIncantationRed:    runeblade.MaleficIncantationRed{},
	card.MaleficIncantationYellow: runeblade.MaleficIncantationYellow{},
	card.MaleficIncantationBlue:   runeblade.MaleficIncantationBlue{},

	card.MauvrionSkiesRed:    runeblade.MauvrionSkiesRed{},
	card.MauvrionSkiesYellow: runeblade.MauvrionSkiesYellow{},
	card.MauvrionSkiesBlue:   runeblade.MauvrionSkiesBlue{},

	card.MeatAndGreetRed:    runeblade.MeatAndGreetRed{},
	card.MeatAndGreetYellow: runeblade.MeatAndGreetYellow{},
	card.MeatAndGreetBlue:   runeblade.MeatAndGreetBlue{},

	card.OathOfTheArknightRed:    runeblade.OathOfTheArknightRed{},
	card.OathOfTheArknightYellow: runeblade.OathOfTheArknightYellow{},
	card.OathOfTheArknightBlue:   runeblade.OathOfTheArknightBlue{},

	card.ReadTheRunesRed:    runeblade.ReadTheRunesRed{},
	card.ReadTheRunesYellow: runeblade.ReadTheRunesYellow{},
	card.ReadTheRunesBlue:   runeblade.ReadTheRunesBlue{},

	card.ReduceToRunechantRed:    runeblade.ReduceToRunechantRed{},
	card.ReduceToRunechantYellow: runeblade.ReduceToRunechantYellow{},
	card.ReduceToRunechantBlue:   runeblade.ReduceToRunechantBlue{},

	card.ReekOfCorruptionRed:    runeblade.ReekOfCorruptionRed{},
	card.ReekOfCorruptionYellow: runeblade.ReekOfCorruptionYellow{},
	card.ReekOfCorruptionBlue:   runeblade.ReekOfCorruptionBlue{},

	card.RuneFlashRed:    runeblade.RuneFlashRed{},
	card.RuneFlashYellow: runeblade.RuneFlashYellow{},
	card.RuneFlashBlue:   runeblade.RuneFlashBlue{},

	card.RunebloodIncantationRed:    runeblade.RunebloodIncantationRed{},
	card.RunebloodIncantationYellow: runeblade.RunebloodIncantationYellow{},
	card.RunebloodIncantationBlue:   runeblade.RunebloodIncantationBlue{},

	card.RuneragerSwarmRed:    runeblade.RuneragerSwarmRed{},
	card.RuneragerSwarmYellow: runeblade.RuneragerSwarmYellow{},
	card.RuneragerSwarmBlue:   runeblade.RuneragerSwarmBlue{},

	card.RunicFellingsongRed:    runeblade.RunicFellingsongRed{},
	card.RunicFellingsongYellow: runeblade.RunicFellingsongYellow{},
	card.RunicFellingsongBlue:   runeblade.RunicFellingsongBlue{},

	card.RunicReapingRed:    runeblade.RunicReapingRed{},
	card.RunicReapingYellow: runeblade.RunicReapingYellow{},
	card.RunicReapingBlue:   runeblade.RunicReapingBlue{},

	card.ShrillOfSkullformRed:    runeblade.ShrillOfSkullformRed{},
	card.ShrillOfSkullformYellow: runeblade.ShrillOfSkullformYellow{},
	card.ShrillOfSkullformBlue:   runeblade.ShrillOfSkullformBlue{},

	card.SigilOfDeadwoodBlue: runeblade.SigilOfDeadwoodBlue{},

	card.SigilOfSilphidaeBlue: runeblade.SigilOfSilphidaeBlue{},

	card.SigilOfSufferingRed:    runeblade.SigilOfSufferingRed{},
	card.SigilOfSufferingYellow: runeblade.SigilOfSufferingYellow{},
	card.SigilOfSufferingBlue:   runeblade.SigilOfSufferingBlue{},

	card.SigilOfTheArknightBlue: runeblade.SigilOfTheArknightBlue{},

	card.SingeingSteelbladeRed:    runeblade.SingeingSteelbladeRed{},
	card.SingeingSteelbladeYellow: runeblade.SingeingSteelbladeYellow{},
	card.SingeingSteelbladeBlue:   runeblade.SingeingSteelbladeBlue{},

	card.SkyFireLanternsRed:    runeblade.SkyFireLanternsRed{},
	card.SkyFireLanternsYellow: runeblade.SkyFireLanternsYellow{},
	card.SkyFireLanternsBlue:   runeblade.SkyFireLanternsBlue{},

	card.SpellbladeAssaultRed:    runeblade.SpellbladeAssaultRed{},
	card.SpellbladeAssaultYellow: runeblade.SpellbladeAssaultYellow{},
	card.SpellbladeAssaultBlue:   runeblade.SpellbladeAssaultBlue{},

	card.SpellbladeStrikeRed:    runeblade.SpellbladeStrikeRed{},
	card.SpellbladeStrikeYellow: runeblade.SpellbladeStrikeYellow{},
	card.SpellbladeStrikeBlue:   runeblade.SpellbladeStrikeBlue{},

	card.SplinteringDeadwoodRed:    runeblade.SplinteringDeadwoodRed{},
	card.SplinteringDeadwoodYellow: runeblade.SplinteringDeadwoodYellow{},
	card.SplinteringDeadwoodBlue:   runeblade.SplinteringDeadwoodBlue{},

	card.SutcliffesResearchNotesRed:    runeblade.SutcliffesResearchNotesRed{},
	card.SutcliffesResearchNotesYellow: runeblade.SutcliffesResearchNotesYellow{},
	card.SutcliffesResearchNotesBlue:   runeblade.SutcliffesResearchNotesBlue{},

	card.VantagePointRed:    runeblade.VantagePointRed{},
	card.VantagePointYellow: runeblade.VantagePointYellow{},
	card.VantagePointBlue:   runeblade.VantagePointBlue{},

	card.VexingMaliceRed:    runeblade.VexingMaliceRed{},
	card.VexingMaliceYellow: runeblade.VexingMaliceYellow{},
	card.VexingMaliceBlue:   runeblade.VexingMaliceBlue{},

	card.WeepingBattlegroundRed:    runeblade.WeepingBattlegroundRed{},
	card.WeepingBattlegroundYellow: runeblade.WeepingBattlegroundYellow{},
	card.WeepingBattlegroundBlue:   runeblade.WeepingBattlegroundBlue{},

	card.DodgeBlue: generic.DodgeBlue{},

	card.EvasiveLeapRed:    generic.EvasiveLeapRed{},
	card.EvasiveLeapYellow: generic.EvasiveLeapYellow{},
	card.EvasiveLeapBlue:   generic.EvasiveLeapBlue{},

	card.FateForeseenRed:    generic.FateForeseenRed{},
	card.FateForeseenYellow: generic.FateForeseenYellow{},
	card.FateForeseenBlue:   generic.FateForeseenBlue{},

	card.LayLowYellow: generic.LayLowYellow{},

	card.PutInContextBlue: generic.PutInContextBlue{},

	card.RiseAboveRed:    generic.RiseAboveRed{},
	card.RiseAboveYellow: generic.RiseAboveYellow{},
	card.RiseAboveBlue:   generic.RiseAboveBlue{},

	card.SinkBelowRed:    generic.SinkBelowRed{},
	card.SinkBelowYellow: generic.SinkBelowYellow{},
	card.SinkBelowBlue:   generic.SinkBelowBlue{},

	card.SpringboardSomersaultYellow: generic.SpringboardSomersaultYellow{},

	card.ToughenUpBlue: generic.ToughenUpBlue{},

	card.UnmovableRed:    generic.UnmovableRed{},
	card.UnmovableYellow: generic.UnmovableYellow{},
	card.UnmovableBlue:   generic.UnmovableBlue{},

	card.FakeRedAttack:  fake.RedAttack{},
	card.FakeBlueAttack: fake.BlueAttack{},
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
	if id == card.Invalid || int(id) >= len(byID) || byID[id] == nil {
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
		if byID[id] == nil {
			continue
		}
		out = append(out, ID(id))
	}
	return out
}

// Count is the number of registered cards (excluding Invalid).
func Count() int {
	n := 0
	for id := 1; id < len(byID); id++ {
		if byID[id] != nil {
			n++
		}
	}
	return n
}

// Deckable returns every registered card ID that's legal to put in a real deck — i.e. every
// registered card except the test-only fakes. Freshly allocated; safe to mutate.
func Deckable() []ID {
	out := make([]ID, 0, len(byID)-1)
	for id := 1; id < len(byID); id++ {
		if byID[id] == nil {
			continue
		}
		switch ID(id) {
		case card.FakeRedAttack, card.FakeBlueAttack:
			continue
		}
		out = append(out, ID(id))
	}
	return out
}
