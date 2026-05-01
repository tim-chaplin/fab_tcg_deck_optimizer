// Package cards is the master registry of every implemented card. The canonical ID type and
// constants live in package card; this package maps IDs to concrete Card values and provides
// lookup / iteration helpers useful for random deck generation, serialization, and compact
// equality checks.
//
// Weapons aren't ID-indexed — they're equipment, not deck cards. The weapon roster lives in
// package weapon alongside the Weapon implementations (weapon.All, weapon.ByName).
package registry

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/optimizations"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// CardID aliases ids.CardID so callers of this package don't need two imports just to hold IDs.
type CardID = ids.CardID

// cardsByID is indexed directly by CardID. Index 0 (Invalid) is nil.
var cardsByID = []sim.Card{
	ids.InvalidCard: nil,

	ids.AetherSlashRed:    cards.AetherSlashRed{},
	ids.AetherSlashYellow: cards.AetherSlashYellow{},
	ids.AetherSlashBlue:   cards.AetherSlashBlue{},

	ids.AmplifyTheArknightRed:    cards.AmplifyTheArknightRed{},
	ids.AmplifyTheArknightYellow: cards.AmplifyTheArknightYellow{},
	ids.AmplifyTheArknightBlue:   cards.AmplifyTheArknightBlue{},

	ids.ArcaneCussingRed:    cards.ArcaneCussingRed{},
	ids.ArcaneCussingYellow: cards.ArcaneCussingYellow{},
	ids.ArcaneCussingBlue:   cards.ArcaneCussingBlue{},

	ids.ArcanicCrackleRed:    cards.ArcanicCrackleRed{},
	ids.ArcanicCrackleYellow: cards.ArcanicCrackleYellow{},
	ids.ArcanicCrackleBlue:   cards.ArcanicCrackleBlue{},

	ids.ArcanicSpikeRed:    cards.ArcanicSpikeRed{},
	ids.ArcanicSpikeYellow: cards.ArcanicSpikeYellow{},
	ids.ArcanicSpikeBlue:   cards.ArcanicSpikeBlue{},

	ids.BlessingOfOccultRed:    cards.BlessingOfOccultRed{},
	ids.BlessingOfOccultYellow: cards.BlessingOfOccultYellow{},
	ids.BlessingOfOccultBlue:   cards.BlessingOfOccultBlue{},

	ids.BloodspillInvocationRed:    cards.BloodspillInvocationRed{},
	ids.BloodspillInvocationYellow: cards.BloodspillInvocationYellow{},
	ids.BloodspillInvocationBlue:   cards.BloodspillInvocationBlue{},

	ids.ConsumingVolitionRed:    cards.ConsumingVolitionRed{},
	ids.ConsumingVolitionYellow: cards.ConsumingVolitionYellow{},
	ids.ConsumingVolitionBlue:   cards.ConsumingVolitionBlue{},

	ids.DeathlyDuetRed:    cards.DeathlyDuetRed{},
	ids.DeathlyDuetYellow: cards.DeathlyDuetYellow{},
	ids.DeathlyDuetBlue:   cards.DeathlyDuetBlue{},

	ids.DrawnToTheDarkDimensionRed:    cards.DrawnToTheDarkDimensionRed{},
	ids.DrawnToTheDarkDimensionYellow: cards.DrawnToTheDarkDimensionYellow{},
	ids.DrawnToTheDarkDimensionBlue:   cards.DrawnToTheDarkDimensionBlue{},

	ids.HitTheHighNotesRed:    cards.HitTheHighNotesRed{},
	ids.HitTheHighNotesYellow: cards.HitTheHighNotesYellow{},
	ids.HitTheHighNotesBlue:   cards.HitTheHighNotesBlue{},

	ids.HocusPocusRed:    cards.HocusPocusRed{},
	ids.HocusPocusYellow: cards.HocusPocusYellow{},
	ids.HocusPocusBlue:   cards.HocusPocusBlue{},

	ids.MaleficIncantationRed:    cards.MaleficIncantationRed{},
	ids.MaleficIncantationYellow: cards.MaleficIncantationYellow{},
	ids.MaleficIncantationBlue:   cards.MaleficIncantationBlue{},

	ids.MauvrionSkiesRed:    cards.MauvrionSkiesRed{},
	ids.MauvrionSkiesYellow: cards.MauvrionSkiesYellow{},
	ids.MauvrionSkiesBlue:   cards.MauvrionSkiesBlue{},

	ids.MeatAndGreetRed:    cards.MeatAndGreetRed{},
	ids.MeatAndGreetYellow: cards.MeatAndGreetYellow{},
	ids.MeatAndGreetBlue:   cards.MeatAndGreetBlue{},

	ids.OathOfTheArknightRed:    cards.OathOfTheArknightRed{},
	ids.OathOfTheArknightYellow: cards.OathOfTheArknightYellow{},
	ids.OathOfTheArknightBlue:   cards.OathOfTheArknightBlue{},

	ids.ReadTheRunesRed:    cards.ReadTheRunesRed{},
	ids.ReadTheRunesYellow: cards.ReadTheRunesYellow{},
	ids.ReadTheRunesBlue:   cards.ReadTheRunesBlue{},

	ids.ReduceToRunechantRed:    cards.ReduceToRunechantRed{},
	ids.ReduceToRunechantYellow: cards.ReduceToRunechantYellow{},
	ids.ReduceToRunechantBlue:   cards.ReduceToRunechantBlue{},

	ids.ReekOfCorruptionRed:    cards.ReekOfCorruptionRed{},
	ids.ReekOfCorruptionYellow: cards.ReekOfCorruptionYellow{},
	ids.ReekOfCorruptionBlue:   cards.ReekOfCorruptionBlue{},

	ids.RuneFlashRed:    cards.RuneFlashRed{},
	ids.RuneFlashYellow: cards.RuneFlashYellow{},
	ids.RuneFlashBlue:   cards.RuneFlashBlue{},

	ids.RunebloodIncantationRed:    cards.RunebloodIncantationRed{},
	ids.RunebloodIncantationYellow: cards.RunebloodIncantationYellow{},
	ids.RunebloodIncantationBlue:   cards.RunebloodIncantationBlue{},

	ids.RuneragerSwarmRed:    cards.RuneragerSwarmRed{},
	ids.RuneragerSwarmYellow: cards.RuneragerSwarmYellow{},
	ids.RuneragerSwarmBlue:   cards.RuneragerSwarmBlue{},

	ids.RunicFellingsongRed:    cards.RunicFellingsongRed{},
	ids.RunicFellingsongYellow: cards.RunicFellingsongYellow{},
	ids.RunicFellingsongBlue:   cards.RunicFellingsongBlue{},

	ids.RunicReapingRed:    cards.RunicReapingRed{},
	ids.RunicReapingYellow: cards.RunicReapingYellow{},
	ids.RunicReapingBlue:   cards.RunicReapingBlue{},

	ids.ShrillOfSkullformRed:    cards.ShrillOfSkullformRed{},
	ids.ShrillOfSkullformYellow: cards.ShrillOfSkullformYellow{},
	ids.ShrillOfSkullformBlue:   cards.ShrillOfSkullformBlue{},

	ids.SigilOfDeadwoodBlue: cards.SigilOfDeadwoodBlue{},

	ids.SigilOfSilphidaeBlue: cards.SigilOfSilphidaeBlue{},

	ids.SigilOfSufferingRed:    cards.SigilOfSufferingRed{},
	ids.SigilOfSufferingYellow: cards.SigilOfSufferingYellow{},
	ids.SigilOfSufferingBlue:   cards.SigilOfSufferingBlue{},

	ids.SigilOfTheArknightBlue: cards.SigilOfTheArknightBlue{},

	ids.SingeingSteelbladeRed:    cards.SingeingSteelbladeRed{},
	ids.SingeingSteelbladeYellow: cards.SingeingSteelbladeYellow{},
	ids.SingeingSteelbladeBlue:   cards.SingeingSteelbladeBlue{},

	ids.SkyFireLanternsRed:    cards.SkyFireLanternsRed{},
	ids.SkyFireLanternsYellow: cards.SkyFireLanternsYellow{},
	ids.SkyFireLanternsBlue:   cards.SkyFireLanternsBlue{},

	ids.SpellbladeAssaultRed:    cards.SpellbladeAssaultRed{},
	ids.SpellbladeAssaultYellow: cards.SpellbladeAssaultYellow{},
	ids.SpellbladeAssaultBlue:   cards.SpellbladeAssaultBlue{},

	ids.SpellbladeStrikeRed:    cards.SpellbladeStrikeRed{},
	ids.SpellbladeStrikeYellow: cards.SpellbladeStrikeYellow{},
	ids.SpellbladeStrikeBlue:   cards.SpellbladeStrikeBlue{},

	ids.SutcliffesResearchNotesRed:    cards.SutcliffesResearchNotesRed{},
	ids.SutcliffesResearchNotesYellow: cards.SutcliffesResearchNotesYellow{},
	ids.SutcliffesResearchNotesBlue:   cards.SutcliffesResearchNotesBlue{},

	ids.VexingMaliceRed:    cards.VexingMaliceRed{},
	ids.VexingMaliceYellow: cards.VexingMaliceYellow{},
	ids.VexingMaliceBlue:   cards.VexingMaliceBlue{},

	ids.WeepingBattlegroundRed:    cards.WeepingBattlegroundRed{},
	ids.WeepingBattlegroundYellow: cards.WeepingBattlegroundYellow{},
	ids.WeepingBattlegroundBlue:   cards.WeepingBattlegroundBlue{},

	ids.AdrenalineRushRed:    cards.AdrenalineRushRed{},
	ids.AdrenalineRushYellow: cards.AdrenalineRushYellow{},
	ids.AdrenalineRushBlue:   cards.AdrenalineRushBlue{},

	ids.BlowForABlowRed: cards.BlowForABlowRed{},

	ids.BrutalAssaultRed:    cards.BrutalAssaultRed{},
	ids.BrutalAssaultYellow: cards.BrutalAssaultYellow{},
	ids.BrutalAssaultBlue:   cards.BrutalAssaultBlue{},

	ids.ClarityPotionBlue: cards.ClarityPotionBlue{},

	ids.ComeToFightRed:    cards.ComeToFightRed{},
	ids.ComeToFightYellow: cards.ComeToFightYellow{},
	ids.ComeToFightBlue:   cards.ComeToFightBlue{},

	ids.CriticalStrikeRed:    cards.CriticalStrikeRed{},
	ids.CriticalStrikeYellow: cards.CriticalStrikeYellow{},
	ids.CriticalStrikeBlue:   cards.CriticalStrikeBlue{},

	ids.DemolitionCrewRed:    cards.DemolitionCrewRed{},
	ids.DemolitionCrewYellow: cards.DemolitionCrewYellow{},
	ids.DemolitionCrewBlue:   cards.DemolitionCrewBlue{},

	ids.DodgeBlue: cards.DodgeBlue{},

	ids.EvasiveLeapRed:    cards.EvasiveLeapRed{},
	ids.EvasiveLeapYellow: cards.EvasiveLeapYellow{},
	ids.EvasiveLeapBlue:   cards.EvasiveLeapBlue{},

	ids.FateForeseenRed:    cards.FateForeseenRed{},
	ids.FateForeseenYellow: cards.FateForeseenYellow{},
	ids.FateForeseenBlue:   cards.FateForeseenBlue{},

	ids.FerventForerunnerRed:    cards.FerventForerunnerRed{},
	ids.FerventForerunnerYellow: cards.FerventForerunnerYellow{},
	ids.FerventForerunnerBlue:   cards.FerventForerunnerBlue{},

	ids.FiddlersGreenRed:    cards.FiddlersGreenRed{},
	ids.FiddlersGreenYellow: cards.FiddlersGreenYellow{},
	ids.FiddlersGreenBlue:   cards.FiddlersGreenBlue{},

	ids.FlyingHighRed:    cards.FlyingHighRed{},
	ids.FlyingHighYellow: cards.FlyingHighYellow{},
	ids.FlyingHighBlue:   cards.FlyingHighBlue{},

	ids.ForceSightRed:    cards.ForceSightRed{},
	ids.ForceSightYellow: cards.ForceSightYellow{},
	ids.ForceSightBlue:   cards.ForceSightBlue{},

	ids.FyendalsFightingSpiritRed:    cards.FyendalsFightingSpiritRed{},
	ids.FyendalsFightingSpiritYellow: cards.FyendalsFightingSpiritYellow{},
	ids.FyendalsFightingSpiritBlue:   cards.FyendalsFightingSpiritBlue{},

	ids.HealingBalmRed:    cards.HealingBalmRed{},
	ids.HealingBalmYellow: cards.HealingBalmYellow{},
	ids.HealingBalmBlue:   cards.HealingBalmBlue{},

	ids.LifeForALifeRed:    cards.LifeForALifeRed{},
	ids.LifeForALifeYellow: cards.LifeForALifeYellow{},
	ids.LifeForALifeBlue:   cards.LifeForALifeBlue{},

	ids.MinnowismRed:    cards.MinnowismRed{},
	ids.MinnowismYellow: cards.MinnowismYellow{},
	ids.MinnowismBlue:   cards.MinnowismBlue{},

	ids.MoonWishRed:    cards.MoonWishRed{},
	ids.MoonWishYellow: cards.MoonWishYellow{},
	ids.MoonWishBlue:   cards.MoonWishBlue{},

	ids.MuscleMuttYellow: cards.MuscleMuttYellow{},

	ids.NimblismRed:    cards.NimblismRed{},
	ids.NimblismYellow: cards.NimblismYellow{},
	ids.NimblismBlue:   cards.NimblismBlue{},

	ids.OnTheHorizonRed:    cards.OnTheHorizonRed{},
	ids.OnTheHorizonYellow: cards.OnTheHorizonYellow{},
	ids.OnTheHorizonBlue:   cards.OnTheHorizonBlue{},

	ids.PoundForPoundRed:    cards.PoundForPoundRed{},
	ids.PoundForPoundYellow: cards.PoundForPoundYellow{},
	ids.PoundForPoundBlue:   cards.PoundForPoundBlue{},

	ids.RagingOnslaughtRed:    cards.RagingOnslaughtRed{},
	ids.RagingOnslaughtYellow: cards.RagingOnslaughtYellow{},
	ids.RagingOnslaughtBlue:   cards.RagingOnslaughtBlue{},

	ids.RavenousRabbleRed:    cards.RavenousRabbleRed{},
	ids.RavenousRabbleYellow: cards.RavenousRabbleYellow{},
	ids.RavenousRabbleBlue:   cards.RavenousRabbleBlue{},

	ids.ScarForAScarRed:    cards.ScarForAScarRed{},
	ids.ScarForAScarYellow: cards.ScarForAScarYellow{},
	ids.ScarForAScarBlue:   cards.ScarForAScarBlue{},

	ids.ScoutThePeripheryRed:    cards.ScoutThePeripheryRed{},
	ids.ScoutThePeripheryYellow: cards.ScoutThePeripheryYellow{},
	ids.ScoutThePeripheryBlue:   cards.ScoutThePeripheryBlue{},

	ids.SigilOfFyendalBlue: cards.SigilOfFyendalBlue{},

	ids.SirensOfSafeHarborRed:    cards.SirensOfSafeHarborRed{},
	ids.SirensOfSafeHarborYellow: cards.SirensOfSafeHarborYellow{},
	ids.SirensOfSafeHarborBlue:   cards.SirensOfSafeHarborBlue{},

	ids.SloggismRed:    cards.SloggismRed{},
	ids.SloggismYellow: cards.SloggismYellow{},
	ids.SloggismBlue:   cards.SloggismBlue{},

	ids.SnatchRed:    cards.SnatchRed{},
	ids.SnatchYellow: cards.SnatchYellow{},
	ids.SnatchBlue:   cards.SnatchBlue{},

	ids.SpringboardSomersaultYellow: cards.SpringboardSomersaultYellow{},

	ids.SpringLoadRed:    cards.SpringLoadRed{},
	ids.SpringLoadYellow: cards.SpringLoadYellow{},
	ids.SpringLoadBlue:   cards.SpringLoadBlue{},

	ids.SunKissRed:    cards.SunKissRed{},
	ids.SunKissYellow: cards.SunKissYellow{},
	ids.SunKissBlue:   cards.SunKissBlue{},

	ids.ToughenUpBlue: cards.ToughenUpBlue{},

	ids.TrotAlongBlue: cards.TrotAlongBlue{},

	ids.UnmovableRed:    cards.UnmovableRed{},
	ids.UnmovableYellow: cards.UnmovableYellow{},
	ids.UnmovableBlue:   cards.UnmovableBlue{},

	ids.VigorRushRed:    cards.VigorRushRed{},
	ids.VigorRushYellow: cards.VigorRushYellow{},
	ids.VigorRushBlue:   cards.VigorRushBlue{},

	ids.WaterTheSeedsRed:    cards.WaterTheSeedsRed{},
	ids.WaterTheSeedsYellow: cards.WaterTheSeedsYellow{},
	ids.WaterTheSeedsBlue:   cards.WaterTheSeedsBlue{},

	ids.WhisperOfTheOracleRed:    cards.WhisperOfTheOracleRed{},
	ids.WhisperOfTheOracleYellow: cards.WhisperOfTheOracleYellow{},
	ids.WhisperOfTheOracleBlue:   cards.WhisperOfTheOracleBlue{},

	ids.WoundedBullRed:    cards.WoundedBullRed{},
	ids.WoundedBullYellow: cards.WoundedBullYellow{},
	ids.WoundedBullBlue:   cards.WoundedBullBlue{},

	ids.WoundingBlowRed:    cards.WoundingBlowRed{},
	ids.WoundingBlowYellow: cards.WoundingBlowYellow{},
	ids.WoundingBlowBlue:   cards.WoundingBlowBlue{},

	ids.ZealousBeltingRed:    cards.ZealousBeltingRed{},
	ids.ZealousBeltingYellow: cards.ZealousBeltingYellow{},
	ids.ZealousBeltingBlue:   cards.ZealousBeltingBlue{},
}

// init eagerly populates package sim's chain-step text and DisplayName caches so the
// per-Play hot path is pure cache reads, and wires the sim → registry forward-declared
// hooks (sim.GetCard / sim.DeckableCards / sim.AllWeapons) so sim's deck builder can
// reach the registry without importing it (would cycle through cards → sim → registry).
// Done at registration time because the registry is the only place that knows the full
// card set, and the caches are sized for the full ID space.
func init() {
	optimizations.WarmChainStepCache(cardsByID)
	optimizations.WarmDisplayNameCache(cardsByID)
	sim.GetCard = GetCard
	sim.DeckableCards = func() []ids.CardID { return DeckableCards() }
	sim.AllWeapons = AllWeapons
}

// cardsByName maps sim.DisplayName(c) → CardID for reverse lookup. Built once at init. Keyed
// on DisplayName (not bare Name) so each pitch variant gets a distinct entry — Card.Name()
// collapses all three printings to the same base string, so it's not a unique key.
var cardsByName = func() map[string]CardID {
	m := make(map[string]CardID, len(cardsByID)-1)
	for id, c := range cardsByID {
		if c == nil {
			continue
		}
		m[sim.DisplayName(c)] = CardID(id)
	}
	return m
}()

// GetCard returns the card for the given ID, or nil when the ID has no registered card —
// the case for NotImplemented and Unplayable IDs whose subpackages aren't imported by the
// registry. Panics if id is the Invalid sentinel or out of the registered ID range, since
// those indicate a programming error (an out-of-range or sentinel ID). Callers iterating
// IDs should null-check the result.
func GetCard(id CardID) sim.Card {
	if id == ids.InvalidCard || int(id) >= len(cardsByID) {
		panic("cardindex: invalid card ID")
	}
	return cardsByID[id]
}

// CardByName looks up an ID by the card's DisplayName ("Aether Slash [R]"). Returns
// (Invalid, false) if no such card.
func CardByName(name string) (CardID, bool) {
	id, ok := cardsByName[name]
	return id, ok
}

// AllCards returns every valid card ID in registration order. The returned slice is freshly
// allocated and safe for the caller to mutate (e.g. shuffle for random deck generation).
func AllCards() []CardID {
	out := make([]CardID, 0, len(cardsByID)-1)
	for id := 1; id < len(cardsByID); id++ {
		if cardsByID[id] == nil {
			continue
		}
		out = append(out, CardID(id))
	}
	return out
}

// DeckableCards returns every registered card ID that's legal to put in a real deck.
// Freshly allocated; safe to mutate. Test-only synthetic IDs (testutils.FakeRedAttack, …)
// live in the testutils package and aren't registered here, so this is just AllCards()
// under a different name — kept distinct so callers who want "deck-legal cards" stay
// readable even if the registry ever holds non-deckable entries again.
func DeckableCards() []CardID {
	return AllCards()
}
