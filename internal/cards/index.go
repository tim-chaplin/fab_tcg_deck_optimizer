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

// weaponsByName maps Weapon.Name() → Weapon for reverse lookup. Built once at init.
var weaponsByName = func() map[string]weapon.Weapon {
	m := make(map[string]weapon.Weapon, len(AllWeapons))
	for _, w := range AllWeapons {
		m[w.Name()] = w
	}
	return m
}()

// WeaponByName returns the registered Weapon for the given display name. Returns (nil, false)
// when no such weapon exists — serialization callers surface that to the user.
func WeaponByName(name string) (weapon.Weapon, bool) {
	w, ok := weaponsByName[name]
	return w, ok
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

	card.AdrenalineRushRed: generic.AdrenalineRushRed{},
	card.AdrenalineRushYellow: generic.AdrenalineRushYellow{},
	card.AdrenalineRushBlue: generic.AdrenalineRushBlue{},

	card.BackAlleyBreaklineRed: generic.BackAlleyBreaklineRed{},
	card.BackAlleyBreaklineYellow: generic.BackAlleyBreaklineYellow{},
	card.BackAlleyBreaklineBlue: generic.BackAlleyBreaklineBlue{},

	card.BarragingBrawnhideRed: generic.BarragingBrawnhideRed{},
	card.BarragingBrawnhideYellow: generic.BarragingBrawnhideYellow{},
	card.BarragingBrawnhideBlue: generic.BarragingBrawnhideBlue{},

	card.BattlefrontBastionRed: generic.BattlefrontBastionRed{},
	card.BattlefrontBastionYellow: generic.BattlefrontBastionYellow{},
	card.BattlefrontBastionBlue: generic.BattlefrontBastionBlue{},

	card.BelittleRed: generic.BelittleRed{},
	card.BelittleYellow: generic.BelittleYellow{},
	card.BelittleBlue: generic.BelittleBlue{},

	card.BlanchRed: generic.BlanchRed{},
	card.BlanchYellow: generic.BlanchYellow{},
	card.BlanchBlue: generic.BlanchBlue{},

	card.BlowForABlowRed: generic.BlowForABlowRed{},

	card.BlusterBuffRed: generic.BlusterBuffRed{},

	card.BrandishRed: generic.BrandishRed{},
	card.BrandishYellow: generic.BrandishYellow{},
	card.BrandishBlue: generic.BrandishBlue{},

	card.BrothersInArmsRed: generic.BrothersInArmsRed{},
	card.BrothersInArmsYellow: generic.BrothersInArmsYellow{},
	card.BrothersInArmsBlue: generic.BrothersInArmsBlue{},

	card.BrutalAssaultRed: generic.BrutalAssaultRed{},
	card.BrutalAssaultYellow: generic.BrutalAssaultYellow{},
	card.BrutalAssaultBlue: generic.BrutalAssaultBlue{},

	card.CadaverousContrabandRed: generic.CadaverousContrabandRed{},
	card.CadaverousContrabandYellow: generic.CadaverousContrabandYellow{},
	card.CadaverousContrabandBlue: generic.CadaverousContrabandBlue{},

	card.CaptainsCallRed: generic.CaptainsCallRed{},
	card.CaptainsCallYellow: generic.CaptainsCallYellow{},
	card.CaptainsCallBlue: generic.CaptainsCallBlue{},

	card.CashInYellow: generic.CashInYellow{},

	card.ChestPuffRed: generic.ChestPuffRed{},

	card.ClearwaterElixirRed: generic.ClearwaterElixirRed{},

	card.ComeToFightRed: generic.ComeToFightRed{},
	card.ComeToFightYellow: generic.ComeToFightYellow{},
	card.ComeToFightBlue: generic.ComeToFightBlue{},

	card.CrashDownTheGatesRed: generic.CrashDownTheGatesRed{},
	card.CrashDownTheGatesYellow: generic.CrashDownTheGatesYellow{},
	card.CrashDownTheGatesBlue: generic.CrashDownTheGatesBlue{},

	card.CriticalStrikeRed: generic.CriticalStrikeRed{},
	card.CriticalStrikeYellow: generic.CriticalStrikeYellow{},
	card.CriticalStrikeBlue: generic.CriticalStrikeBlue{},

	card.CutDownToSizeRed: generic.CutDownToSizeRed{},
	card.CutDownToSizeYellow: generic.CutDownToSizeYellow{},
	card.CutDownToSizeBlue: generic.CutDownToSizeBlue{},

	card.DemolitionCrewRed: generic.DemolitionCrewRed{},
	card.DemolitionCrewYellow: generic.DemolitionCrewYellow{},
	card.DemolitionCrewBlue: generic.DemolitionCrewBlue{},

	card.DestructiveDeliberationRed: generic.DestructiveDeliberationRed{},
	card.DestructiveDeliberationYellow: generic.DestructiveDeliberationYellow{},
	card.DestructiveDeliberationBlue: generic.DestructiveDeliberationBlue{},

	card.DodgeBlue: generic.DodgeBlue{},

	card.DownButNotOutRed: generic.DownButNotOutRed{},
	card.DownButNotOutYellow: generic.DownButNotOutYellow{},
	card.DownButNotOutBlue: generic.DownButNotOutBlue{},

	card.DragDownRed: generic.DragDownRed{},
	card.DragDownYellow: generic.DragDownYellow{},
	card.DragDownBlue: generic.DragDownBlue{},

	card.DroneOfBrutalityRed: generic.DroneOfBrutalityRed{},
	card.DroneOfBrutalityYellow: generic.DroneOfBrutalityYellow{},
	card.DroneOfBrutalityBlue: generic.DroneOfBrutalityBlue{},

	card.EmissaryOfMoonRed: generic.EmissaryOfMoonRed{},

	card.EmissaryOfTidesRed: generic.EmissaryOfTidesRed{},

	card.EmissaryOfWindRed: generic.EmissaryOfWindRed{},

	card.EnchantingMelodyRed: generic.EnchantingMelodyRed{},
	card.EnchantingMelodyYellow: generic.EnchantingMelodyYellow{},
	card.EnchantingMelodyBlue: generic.EnchantingMelodyBlue{},

	card.EvasiveLeapRed: generic.EvasiveLeapRed{},
	card.EvasiveLeapYellow: generic.EvasiveLeapYellow{},
	card.EvasiveLeapBlue: generic.EvasiveLeapBlue{},

	card.FactFindingMissionRed: generic.FactFindingMissionRed{},
	card.FactFindingMissionYellow: generic.FactFindingMissionYellow{},
	card.FactFindingMissionBlue: generic.FactFindingMissionBlue{},

	card.FateForeseenRed: generic.FateForeseenRed{},
	card.FateForeseenYellow: generic.FateForeseenYellow{},
	card.FateForeseenBlue: generic.FateForeseenBlue{},

	card.FeistyLocalsRed: generic.FeistyLocalsRed{},
	card.FeistyLocalsYellow: generic.FeistyLocalsYellow{},
	card.FeistyLocalsBlue: generic.FeistyLocalsBlue{},

	card.FerventForerunnerRed: generic.FerventForerunnerRed{},
	card.FerventForerunnerYellow: generic.FerventForerunnerYellow{},
	card.FerventForerunnerBlue: generic.FerventForerunnerBlue{},

	card.FiddlersGreenRed: generic.FiddlersGreenRed{},
	card.FiddlersGreenYellow: generic.FiddlersGreenYellow{},
	card.FiddlersGreenBlue: generic.FiddlersGreenBlue{},

	card.FlexRed: generic.FlexRed{},
	card.FlexYellow: generic.FlexYellow{},
	card.FlexBlue: generic.FlexBlue{},

	card.FlockOfTheFeatherWalkersRed: generic.FlockOfTheFeatherWalkersRed{},
	card.FlockOfTheFeatherWalkersYellow: generic.FlockOfTheFeatherWalkersYellow{},
	card.FlockOfTheFeatherWalkersBlue: generic.FlockOfTheFeatherWalkersBlue{},

	card.FlyingHighRed: generic.FlyingHighRed{},
	card.FlyingHighYellow: generic.FlyingHighYellow{},
	card.FlyingHighBlue: generic.FlyingHighBlue{},

	card.ForceSightRed: generic.ForceSightRed{},
	card.ForceSightYellow: generic.ForceSightYellow{},
	card.ForceSightBlue: generic.ForceSightBlue{},

	card.FreewheelingRenegadesRed: generic.FreewheelingRenegadesRed{},
	card.FreewheelingRenegadesYellow: generic.FreewheelingRenegadesYellow{},
	card.FreewheelingRenegadesBlue: generic.FreewheelingRenegadesBlue{},

	card.FrontlineScoutRed: generic.FrontlineScoutRed{},
	card.FrontlineScoutYellow: generic.FrontlineScoutYellow{},
	card.FrontlineScoutBlue: generic.FrontlineScoutBlue{},

	card.FyendalsFightingSpiritRed: generic.FyendalsFightingSpiritRed{},
	card.FyendalsFightingSpiritYellow: generic.FyendalsFightingSpiritYellow{},
	card.FyendalsFightingSpiritBlue: generic.FyendalsFightingSpiritBlue{},

	card.GravekeepingRed: generic.GravekeepingRed{},
	card.GravekeepingYellow: generic.GravekeepingYellow{},
	card.GravekeepingBlue: generic.GravekeepingBlue{},

	card.HandBehindThePenRed: generic.HandBehindThePenRed{},

	card.HealingBalmRed: generic.HealingBalmRed{},
	card.HealingBalmYellow: generic.HealingBalmYellow{},
	card.HealingBalmBlue: generic.HealingBalmBlue{},

	card.HighStrikerRed: generic.HighStrikerRed{},
	card.HighStrikerYellow: generic.HighStrikerYellow{},
	card.HighStrikerBlue: generic.HighStrikerBlue{},

	card.HumbleRed: generic.HumbleRed{},
	card.HumbleYellow: generic.HumbleYellow{},
	card.HumbleBlue: generic.HumbleBlue{},

	card.InfectiousHostRed: generic.InfectiousHostRed{},
	card.InfectiousHostYellow: generic.InfectiousHostYellow{},
	card.InfectiousHostBlue: generic.InfectiousHostBlue{},

	card.JackBeNimbleRed: generic.JackBeNimbleRed{},

	card.JackBeQuickRed: generic.JackBeQuickRed{},

	card.LayLowYellow: generic.LayLowYellow{},

	card.LeadTheChargeRed: generic.LeadTheChargeRed{},
	card.LeadTheChargeYellow: generic.LeadTheChargeYellow{},
	card.LeadTheChargeBlue: generic.LeadTheChargeBlue{},

	card.LifeForALifeRed: generic.LifeForALifeRed{},
	card.LifeForALifeYellow: generic.LifeForALifeYellow{},
	card.LifeForALifeBlue: generic.LifeForALifeBlue{},

	card.LifeOfThePartyRed: generic.LifeOfThePartyRed{},
	card.LifeOfThePartyYellow: generic.LifeOfThePartyYellow{},
	card.LifeOfThePartyBlue: generic.LifeOfThePartyBlue{},

	card.LookTuffRed: generic.LookTuffRed{},

	card.LookingForAScrapRed: generic.LookingForAScrapRed{},
	card.LookingForAScrapYellow: generic.LookingForAScrapYellow{},
	card.LookingForAScrapBlue: generic.LookingForAScrapBlue{},

	card.MinnowismRed: generic.MinnowismRed{},
	card.MinnowismYellow: generic.MinnowismYellow{},
	card.MinnowismBlue: generic.MinnowismBlue{},

	card.MoneyOrYourLifeRed: generic.MoneyOrYourLifeRed{},
	card.MoneyOrYourLifeYellow: generic.MoneyOrYourLifeYellow{},
	card.MoneyOrYourLifeBlue: generic.MoneyOrYourLifeBlue{},

	card.MoneyWhereYaMouthIsRed: generic.MoneyWhereYaMouthIsRed{},
	card.MoneyWhereYaMouthIsYellow: generic.MoneyWhereYaMouthIsYellow{},
	card.MoneyWhereYaMouthIsBlue: generic.MoneyWhereYaMouthIsBlue{},

	card.MoonWishRed: generic.MoonWishRed{},
	card.MoonWishYellow: generic.MoonWishYellow{},
	card.MoonWishBlue: generic.MoonWishBlue{},

	card.MuscleMuttYellow: generic.MuscleMuttYellow{},

	card.NimbleStrikeRed: generic.NimbleStrikeRed{},
	card.NimbleStrikeYellow: generic.NimbleStrikeYellow{},
	card.NimbleStrikeBlue: generic.NimbleStrikeBlue{},

	card.NimblismRed: generic.NimblismRed{},
	card.NimblismYellow: generic.NimblismYellow{},
	card.NimblismBlue: generic.NimblismBlue{},

	card.NimbyRed: generic.NimbyRed{},
	card.NimbyYellow: generic.NimbyYellow{},
	card.NimbyBlue: generic.NimbyBlue{},

	card.OnAKnifeEdgeYellow: generic.OnAKnifeEdgeYellow{},

	card.OnTheHorizonRed: generic.OnTheHorizonRed{},
	card.OnTheHorizonYellow: generic.OnTheHorizonYellow{},
	card.OnTheHorizonBlue: generic.OnTheHorizonBlue{},

	card.OutMuscleRed: generic.OutMuscleRed{},
	card.OutMuscleYellow: generic.OutMuscleYellow{},
	card.OutMuscleBlue: generic.OutMuscleBlue{},

	card.OutedRed: generic.OutedRed{},

	card.OverloadRed: generic.OverloadRed{},
	card.OverloadYellow: generic.OverloadYellow{},
	card.OverloadBlue: generic.OverloadBlue{},

	card.PerformanceBonusRed: generic.PerformanceBonusRed{},
	card.PerformanceBonusYellow: generic.PerformanceBonusYellow{},
	card.PerformanceBonusBlue: generic.PerformanceBonusBlue{},

	card.PickACardAnyCardRed: generic.PickACardAnyCardRed{},
	card.PickACardAnyCardYellow: generic.PickACardAnyCardYellow{},
	card.PickACardAnyCardBlue: generic.PickACardAnyCardBlue{},

	card.PlunderRunRed: generic.PlunderRunRed{},
	card.PlunderRunYellow: generic.PlunderRunYellow{},
	card.PlunderRunBlue: generic.PlunderRunBlue{},

	card.PoundForPoundRed: generic.PoundForPoundRed{},
	card.PoundForPoundYellow: generic.PoundForPoundYellow{},
	card.PoundForPoundBlue: generic.PoundForPoundBlue{},

	card.PrimeTheCrowdRed: generic.PrimeTheCrowdRed{},
	card.PrimeTheCrowdYellow: generic.PrimeTheCrowdYellow{},
	card.PrimeTheCrowdBlue: generic.PrimeTheCrowdBlue{},

	card.PromiseOfPlentyRed: generic.PromiseOfPlentyRed{},
	card.PromiseOfPlentyYellow: generic.PromiseOfPlentyYellow{},
	card.PromiseOfPlentyBlue: generic.PromiseOfPlentyBlue{},

	card.PublicBountyRed: generic.PublicBountyRed{},
	card.PublicBountyYellow: generic.PublicBountyYellow{},
	card.PublicBountyBlue: generic.PublicBountyBlue{},

	card.PunchAboveYourWeightRed: generic.PunchAboveYourWeightRed{},
	card.PunchAboveYourWeightYellow: generic.PunchAboveYourWeightYellow{},
	card.PunchAboveYourWeightBlue: generic.PunchAboveYourWeightBlue{},

	card.PursueToTheEdgeOfOblivionRed: generic.PursueToTheEdgeOfOblivionRed{},

	card.PursueToThePitsOfDespairRed: generic.PursueToThePitsOfDespairRed{},

	card.PushThePointRed: generic.PushThePointRed{},
	card.PushThePointYellow: generic.PushThePointYellow{},
	card.PushThePointBlue: generic.PushThePointBlue{},

	card.PutInContextBlue: generic.PutInContextBlue{},

	card.RagingOnslaughtRed: generic.RagingOnslaughtRed{},
	card.RagingOnslaughtYellow: generic.RagingOnslaughtYellow{},
	card.RagingOnslaughtBlue: generic.RagingOnslaughtBlue{},

	card.RallyTheCoastGuardRed: generic.RallyTheCoastGuardRed{},
	card.RallyTheCoastGuardYellow: generic.RallyTheCoastGuardYellow{},
	card.RallyTheCoastGuardBlue: generic.RallyTheCoastGuardBlue{},

	card.RallyTheRearguardRed: generic.RallyTheRearguardRed{},
	card.RallyTheRearguardYellow: generic.RallyTheRearguardYellow{},
	card.RallyTheRearguardBlue: generic.RallyTheRearguardBlue{},

	card.RansackAndRazeBlue: generic.RansackAndRazeBlue{},

	card.RavenousRabbleRed: generic.RavenousRabbleRed{},
	card.RavenousRabbleYellow: generic.RavenousRabbleYellow{},
	card.RavenousRabbleBlue: generic.RavenousRabbleBlue{},

	card.RegainComposureBlue: generic.RegainComposureBlue{},

	card.RegurgitatingSlogRed: generic.RegurgitatingSlogRed{},
	card.RegurgitatingSlogYellow: generic.RegurgitatingSlogYellow{},
	card.RegurgitatingSlogBlue: generic.RegurgitatingSlogBlue{},

	card.RelentlessPursuitBlue: generic.RelentlessPursuitBlue{},

	card.RestvineElixirRed: generic.RestvineElixirRed{},

	card.RiftingRed: generic.RiftingRed{},
	card.RiftingYellow: generic.RiftingYellow{},
	card.RiftingBlue: generic.RiftingBlue{},

	card.RightBehindYouRed: generic.RightBehindYouRed{},
	card.RightBehindYouYellow: generic.RightBehindYouYellow{},
	card.RightBehindYouBlue: generic.RightBehindYouBlue{},

	card.RiseAboveRed: generic.RiseAboveRed{},
	card.RiseAboveYellow: generic.RiseAboveYellow{},
	card.RiseAboveBlue: generic.RiseAboveBlue{},

	card.SapwoodElixirRed: generic.SapwoodElixirRed{},

	card.ScarForAScarRed: generic.ScarForAScarRed{},
	card.ScarForAScarYellow: generic.ScarForAScarYellow{},
	card.ScarForAScarBlue: generic.ScarForAScarBlue{},

	card.ScourTheBattlescapeRed: generic.ScourTheBattlescapeRed{},
	card.ScourTheBattlescapeYellow: generic.ScourTheBattlescapeYellow{},
	card.ScourTheBattlescapeBlue: generic.ScourTheBattlescapeBlue{},

	card.ScoutThePeripheryRed: generic.ScoutThePeripheryRed{},
	card.ScoutThePeripheryYellow: generic.ScoutThePeripheryYellow{},
	card.ScoutThePeripheryBlue: generic.ScoutThePeripheryBlue{},

	card.SeekHorizonRed: generic.SeekHorizonRed{},
	card.SeekHorizonYellow: generic.SeekHorizonYellow{},
	card.SeekHorizonBlue: generic.SeekHorizonBlue{},

	card.SiftRed: generic.SiftRed{},
	card.SiftYellow: generic.SiftYellow{},
	card.SiftBlue: generic.SiftBlue{},

	card.SigilOfCyclesBlue: generic.SigilOfCyclesBlue{},

	card.SigilOfFyendalBlue: generic.SigilOfFyendalBlue{},

	card.SigilOfProtectionRed: generic.SigilOfProtectionRed{},
	card.SigilOfProtectionYellow: generic.SigilOfProtectionYellow{},
	card.SigilOfProtectionBlue: generic.SigilOfProtectionBlue{},

	card.SinkBelowRed: generic.SinkBelowRed{},
	card.SinkBelowYellow: generic.SinkBelowYellow{},
	card.SinkBelowBlue: generic.SinkBelowBlue{},

	card.SirensOfSafeHarborRed: generic.SirensOfSafeHarborRed{},
	card.SirensOfSafeHarborYellow: generic.SirensOfSafeHarborYellow{},
	card.SirensOfSafeHarborBlue: generic.SirensOfSafeHarborBlue{},

	card.SloggismRed: generic.SloggismRed{},
	card.SloggismYellow: generic.SloggismYellow{},
	card.SloggismBlue: generic.SloggismBlue{},

	card.SmashUpRed: generic.SmashUpRed{},

	card.SmashingGoodTimeRed: generic.SmashingGoodTimeRed{},
	card.SmashingGoodTimeYellow: generic.SmashingGoodTimeYellow{},
	card.SmashingGoodTimeBlue: generic.SmashingGoodTimeBlue{},

	card.SnatchRed: generic.SnatchRed{},
	card.SnatchYellow: generic.SnatchYellow{},
	card.SnatchBlue: generic.SnatchBlue{},

	card.SoundTheAlarmRed: generic.SoundTheAlarmRed{},

	card.SpringLoadRed: generic.SpringLoadRed{},
	card.SpringLoadYellow: generic.SpringLoadYellow{},
	card.SpringLoadBlue: generic.SpringLoadBlue{},

	card.SpringboardSomersaultYellow: generic.SpringboardSomersaultYellow{},

	card.StartingStakeYellow: generic.StartingStakeYellow{},

	card.StonyWoottonhogRed: generic.StonyWoottonhogRed{},
	card.StonyWoottonhogYellow: generic.StonyWoottonhogYellow{},
	card.StonyWoottonhogBlue: generic.StonyWoottonhogBlue{},

	card.StrategicPlanningRed: generic.StrategicPlanningRed{},
	card.StrategicPlanningYellow: generic.StrategicPlanningYellow{},
	card.StrategicPlanningBlue: generic.StrategicPlanningBlue{},

	card.StrikeGoldRed: generic.StrikeGoldRed{},
	card.StrikeGoldYellow: generic.StrikeGoldYellow{},
	card.StrikeGoldBlue: generic.StrikeGoldBlue{},

	card.SunKissRed: generic.SunKissRed{},
	card.SunKissYellow: generic.SunKissYellow{},
	card.SunKissBlue: generic.SunKissBlue{},

	card.SurgingMilitiaRed: generic.SurgingMilitiaRed{},
	card.SurgingMilitiaYellow: generic.SurgingMilitiaYellow{},
	card.SurgingMilitiaBlue: generic.SurgingMilitiaBlue{},

	card.TestOfStrengthRed: generic.TestOfStrengthRed{},

	card.TipOffRed: generic.TipOffRed{},
	card.TipOffYellow: generic.TipOffYellow{},
	card.TipOffBlue: generic.TipOffBlue{},

	card.TitForTatBlue: generic.TitForTatBlue{},

	card.TongueTiedRed: generic.TongueTiedRed{},

	card.ToughenUpBlue: generic.ToughenUpBlue{},

	card.TradeInRed: generic.TradeInRed{},
	card.TradeInYellow: generic.TradeInYellow{},
	card.TradeInBlue: generic.TradeInBlue{},

	card.TremorOfIArathaelRed: generic.TremorOfIArathaelRed{},
	card.TremorOfIArathaelYellow: generic.TremorOfIArathaelYellow{},
	card.TremorOfIArathaelBlue: generic.TremorOfIArathaelBlue{},

	card.TrotAlongBlue: generic.TrotAlongBlue{},

	card.UnmovableRed: generic.UnmovableRed{},
	card.UnmovableYellow: generic.UnmovableYellow{},
	card.UnmovableBlue: generic.UnmovableBlue{},

	card.VigorRushRed: generic.VigorRushRed{},
	card.VigorRushYellow: generic.VigorRushYellow{},
	card.VigorRushBlue: generic.VigorRushBlue{},

	card.VisitTheBlacksmithBlue: generic.VisitTheBlacksmithBlue{},

	card.WageGoldRed: generic.WageGoldRed{},
	card.WageGoldYellow: generic.WageGoldYellow{},
	card.WageGoldBlue: generic.WageGoldBlue{},

	card.WalkThePlankRed: generic.WalkThePlankRed{},
	card.WalkThePlankYellow: generic.WalkThePlankYellow{},
	card.WalkThePlankBlue: generic.WalkThePlankBlue{},

	card.WarmongersRecitalRed: generic.WarmongersRecitalRed{},
	card.WarmongersRecitalYellow: generic.WarmongersRecitalYellow{},
	card.WarmongersRecitalBlue: generic.WarmongersRecitalBlue{},

	card.WaterTheSeedsRed: generic.WaterTheSeedsRed{},
	card.WaterTheSeedsYellow: generic.WaterTheSeedsYellow{},
	card.WaterTheSeedsBlue: generic.WaterTheSeedsBlue{},

	card.WhisperOfTheOracleRed: generic.WhisperOfTheOracleRed{},
	card.WhisperOfTheOracleYellow: generic.WhisperOfTheOracleYellow{},
	card.WhisperOfTheOracleBlue: generic.WhisperOfTheOracleBlue{},

	card.WoundedBullRed: generic.WoundedBullRed{},
	card.WoundedBullYellow: generic.WoundedBullYellow{},
	card.WoundedBullBlue: generic.WoundedBullBlue{},

	card.WoundingBlowRed: generic.WoundingBlowRed{},
	card.WoundingBlowYellow: generic.WoundingBlowYellow{},
	card.WoundingBlowBlue: generic.WoundingBlowBlue{},

	card.WreckHavocRed: generic.WreckHavocRed{},
	card.WreckHavocYellow: generic.WreckHavocYellow{},
	card.WreckHavocBlue: generic.WreckHavocBlue{},

	card.YintiYantiRed: generic.YintiYantiRed{},
	card.YintiYantiYellow: generic.YintiYantiYellow{},
	card.YintiYantiBlue: generic.YintiYantiBlue{},

	card.ZealousBeltingRed: generic.ZealousBeltingRed{},
	card.ZealousBeltingYellow: generic.ZealousBeltingYellow{},
	card.ZealousBeltingBlue: generic.ZealousBeltingBlue{},

	card.FakeRedAttack:    fake.RedAttack{},
	card.FakeBlueAttack:   fake.BlueAttack{},
	card.FakeYellowAttack: fake.YellowAttack{},
	card.FakeDrawCantrip:  fake.DrawCantrip{},
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

// Deckable returns every registered card ID that's legal to put in a real deck — i.e. every
// registered card except the test-only fakes. Freshly allocated; safe to mutate.
func Deckable() []ID {
	out := make([]ID, 0, len(byID)-1)
	for id := 1; id < len(byID); id++ {
		if byID[id] == nil {
			continue
		}
		switch ID(id) {
		case card.FakeRedAttack, card.FakeBlueAttack, card.FakeYellowAttack, card.FakeDrawCantrip:
			continue
		}
		out = append(out, ID(id))
	}
	return out
}
