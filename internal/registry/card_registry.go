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

	ids.CondemnToSlaughterRed:    cards.CondemnToSlaughterRed{},
	ids.CondemnToSlaughterYellow: cards.CondemnToSlaughterYellow{},
	ids.CondemnToSlaughterBlue:   cards.CondemnToSlaughterBlue{},

	ids.ConsumingVolitionRed:    cards.ConsumingVolitionRed{},
	ids.ConsumingVolitionYellow: cards.ConsumingVolitionYellow{},
	ids.ConsumingVolitionBlue:   cards.ConsumingVolitionBlue{},

	ids.DeathlyDuetRed:    cards.DeathlyDuetRed{},
	ids.DeathlyDuetYellow: cards.DeathlyDuetYellow{},
	ids.DeathlyDuetBlue:   cards.DeathlyDuetBlue{},

	ids.DrawnToTheDarkDimensionRed:    cards.DrawnToTheDarkDimensionRed{},
	ids.DrawnToTheDarkDimensionYellow: cards.DrawnToTheDarkDimensionYellow{},
	ids.DrawnToTheDarkDimensionBlue:   cards.DrawnToTheDarkDimensionBlue{},

	ids.DrowningDireRed:    cards.DrowningDireRed{},
	ids.DrowningDireYellow: cards.DrowningDireYellow{},
	ids.DrowningDireBlue:   cards.DrowningDireBlue{},

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

	ids.SplinteringDeadwoodRed:    cards.SplinteringDeadwoodRed{},
	ids.SplinteringDeadwoodYellow: cards.SplinteringDeadwoodYellow{},
	ids.SplinteringDeadwoodBlue:   cards.SplinteringDeadwoodBlue{},

	ids.SutcliffesResearchNotesRed:    cards.SutcliffesResearchNotesRed{},
	ids.SutcliffesResearchNotesYellow: cards.SutcliffesResearchNotesYellow{},
	ids.SutcliffesResearchNotesBlue:   cards.SutcliffesResearchNotesBlue{},

	ids.VantagePointRed:    cards.VantagePointRed{},
	ids.VantagePointYellow: cards.VantagePointYellow{},
	ids.VantagePointBlue:   cards.VantagePointBlue{},

	ids.VexingMaliceRed:    cards.VexingMaliceRed{},
	ids.VexingMaliceYellow: cards.VexingMaliceYellow{},
	ids.VexingMaliceBlue:   cards.VexingMaliceBlue{},

	ids.WeepingBattlegroundRed:    cards.WeepingBattlegroundRed{},
	ids.WeepingBattlegroundYellow: cards.WeepingBattlegroundYellow{},
	ids.WeepingBattlegroundBlue:   cards.WeepingBattlegroundBlue{},

	ids.AdrenalineRushRed:    cards.AdrenalineRushRed{},
	ids.AdrenalineRushYellow: cards.AdrenalineRushYellow{},
	ids.AdrenalineRushBlue:   cards.AdrenalineRushBlue{},

	ids.AmuletOfAssertivenessYellow: cards.AmuletOfAssertivenessYellow{},

	ids.AmuletOfEchoesBlue: cards.AmuletOfEchoesBlue{},

	ids.AmuletOfHavencallBlue: cards.AmuletOfHavencallBlue{},

	ids.AmuletOfIgnitionYellow: cards.AmuletOfIgnitionYellow{},

	ids.AmuletOfInterventionBlue: cards.AmuletOfInterventionBlue{},

	ids.AmuletOfOblationBlue: cards.AmuletOfOblationBlue{},

	ids.ArcanePolarityRed:    cards.ArcanePolarityRed{},
	ids.ArcanePolarityYellow: cards.ArcanePolarityYellow{},
	ids.ArcanePolarityBlue:   cards.ArcanePolarityBlue{},

	ids.BackAlleyBreaklineRed:    cards.BackAlleyBreaklineRed{},
	ids.BackAlleyBreaklineYellow: cards.BackAlleyBreaklineYellow{},
	ids.BackAlleyBreaklineBlue:   cards.BackAlleyBreaklineBlue{},

	ids.BarragingBrawnhideRed:    cards.BarragingBrawnhideRed{},
	ids.BarragingBrawnhideYellow: cards.BarragingBrawnhideYellow{},
	ids.BarragingBrawnhideBlue:   cards.BarragingBrawnhideBlue{},

	ids.BattlefrontBastionRed:    cards.BattlefrontBastionRed{},
	ids.BattlefrontBastionYellow: cards.BattlefrontBastionYellow{},
	ids.BattlefrontBastionBlue:   cards.BattlefrontBastionBlue{},

	ids.BelittleRed:    cards.BelittleRed{},
	ids.BelittleYellow: cards.BelittleYellow{},
	ids.BelittleBlue:   cards.BelittleBlue{},

	ids.BladeFlashBlue: cards.BladeFlashBlue{},

	ids.BlanchRed:    cards.BlanchRed{},
	ids.BlanchYellow: cards.BlanchYellow{},
	ids.BlanchBlue:   cards.BlanchBlue{},

	ids.BlowForABlowRed: cards.BlowForABlowRed{},

	ids.BlusterBuffRed: cards.BlusterBuffRed{},

	ids.BrandishRed:    cards.BrandishRed{},
	ids.BrandishYellow: cards.BrandishYellow{},
	ids.BrandishBlue:   cards.BrandishBlue{},

	ids.BrothersInArmsRed:    cards.BrothersInArmsRed{},
	ids.BrothersInArmsYellow: cards.BrothersInArmsYellow{},
	ids.BrothersInArmsBlue:   cards.BrothersInArmsBlue{},

	ids.BrushOffRed:    cards.BrushOffRed{},
	ids.BrushOffYellow: cards.BrushOffYellow{},
	ids.BrushOffBlue:   cards.BrushOffBlue{},

	ids.BrutalAssaultRed:    cards.BrutalAssaultRed{},
	ids.BrutalAssaultYellow: cards.BrutalAssaultYellow{},
	ids.BrutalAssaultBlue:   cards.BrutalAssaultBlue{},

	ids.CadaverousContrabandRed:    cards.CadaverousContrabandRed{},
	ids.CadaverousContrabandYellow: cards.CadaverousContrabandYellow{},
	ids.CadaverousContrabandBlue:   cards.CadaverousContrabandBlue{},

	ids.CalmingBreezeRed: cards.CalmingBreezeRed{},

	ids.CaptainsCallRed:    cards.CaptainsCallRed{},
	ids.CaptainsCallYellow: cards.CaptainsCallYellow{},
	ids.CaptainsCallBlue:   cards.CaptainsCallBlue{},

	ids.CashInYellow: cards.CashInYellow{},

	ids.ChestPuffRed: cards.ChestPuffRed{},

	ids.ClapEmInIronsBlue: cards.ClapEmInIronsBlue{},

	ids.ClarityPotionBlue: cards.ClarityPotionBlue{},

	ids.ClearwaterElixirRed: cards.ClearwaterElixirRed{},

	ids.ComeToFightRed:    cards.ComeToFightRed{},
	ids.ComeToFightYellow: cards.ComeToFightYellow{},
	ids.ComeToFightBlue:   cards.ComeToFightBlue{},

	ids.CountYourBlessingsRed:    cards.CountYourBlessingsRed{},
	ids.CountYourBlessingsYellow: cards.CountYourBlessingsYellow{},
	ids.CountYourBlessingsBlue:   cards.CountYourBlessingsBlue{},

	ids.CrackedBaubleYellow: cards.CrackedBaubleYellow{},

	ids.CrashDownTheGatesRed:    cards.CrashDownTheGatesRed{},
	ids.CrashDownTheGatesYellow: cards.CrashDownTheGatesYellow{},
	ids.CrashDownTheGatesBlue:   cards.CrashDownTheGatesBlue{},

	ids.CriticalStrikeRed:    cards.CriticalStrikeRed{},
	ids.CriticalStrikeYellow: cards.CriticalStrikeYellow{},
	ids.CriticalStrikeBlue:   cards.CriticalStrikeBlue{},

	ids.CutDownToSizeRed:    cards.CutDownToSizeRed{},
	ids.CutDownToSizeYellow: cards.CutDownToSizeYellow{},
	ids.CutDownToSizeBlue:   cards.CutDownToSizeBlue{},

	ids.DemolitionCrewRed:    cards.DemolitionCrewRed{},
	ids.DemolitionCrewYellow: cards.DemolitionCrewYellow{},
	ids.DemolitionCrewBlue:   cards.DemolitionCrewBlue{},

	ids.DestructiveDeliberationRed:    cards.DestructiveDeliberationRed{},
	ids.DestructiveDeliberationYellow: cards.DestructiveDeliberationYellow{},
	ids.DestructiveDeliberationBlue:   cards.DestructiveDeliberationBlue{},

	ids.DestructiveTendenciesBlue: cards.DestructiveTendenciesBlue{},

	ids.DodgeBlue: cards.DodgeBlue{},

	ids.DownButNotOutRed:    cards.DownButNotOutRed{},
	ids.DownButNotOutYellow: cards.DownButNotOutYellow{},
	ids.DownButNotOutBlue:   cards.DownButNotOutBlue{},

	ids.DragDownRed:    cards.DragDownRed{},
	ids.DragDownYellow: cards.DragDownYellow{},
	ids.DragDownBlue:   cards.DragDownBlue{},

	ids.DroneOfBrutalityRed:    cards.DroneOfBrutalityRed{},
	ids.DroneOfBrutalityYellow: cards.DroneOfBrutalityYellow{},
	ids.DroneOfBrutalityBlue:   cards.DroneOfBrutalityBlue{},

	ids.EirinasPrayerRed:    cards.EirinasPrayerRed{},
	ids.EirinasPrayerYellow: cards.EirinasPrayerYellow{},
	ids.EirinasPrayerBlue:   cards.EirinasPrayerBlue{},

	ids.EmissaryOfMoonRed: cards.EmissaryOfMoonRed{},

	ids.EmissaryOfTidesRed: cards.EmissaryOfTidesRed{},

	ids.EmissaryOfWindRed: cards.EmissaryOfWindRed{},

	ids.EnchantingMelodyRed:    cards.EnchantingMelodyRed{},
	ids.EnchantingMelodyYellow: cards.EnchantingMelodyYellow{},
	ids.EnchantingMelodyBlue:   cards.EnchantingMelodyBlue{},

	ids.EnergyPotionBlue: cards.EnergyPotionBlue{},

	ids.EvasiveLeapRed:    cards.EvasiveLeapRed{},
	ids.EvasiveLeapYellow: cards.EvasiveLeapYellow{},
	ids.EvasiveLeapBlue:   cards.EvasiveLeapBlue{},

	ids.EvenBiggerThanThatRed:    cards.EvenBiggerThanThatRed{},
	ids.EvenBiggerThanThatYellow: cards.EvenBiggerThanThatYellow{},
	ids.EvenBiggerThanThatBlue:   cards.EvenBiggerThanThatBlue{},

	ids.ExposedBlue: cards.ExposedBlue{},

	ids.FactFindingMissionRed:    cards.FactFindingMissionRed{},
	ids.FactFindingMissionYellow: cards.FactFindingMissionYellow{},
	ids.FactFindingMissionBlue:   cards.FactFindingMissionBlue{},

	ids.FateForeseenRed:    cards.FateForeseenRed{},
	ids.FateForeseenYellow: cards.FateForeseenYellow{},
	ids.FateForeseenBlue:   cards.FateForeseenBlue{},

	ids.FeistyLocalsRed:    cards.FeistyLocalsRed{},
	ids.FeistyLocalsYellow: cards.FeistyLocalsYellow{},
	ids.FeistyLocalsBlue:   cards.FeistyLocalsBlue{},

	ids.FerventForerunnerRed:    cards.FerventForerunnerRed{},
	ids.FerventForerunnerYellow: cards.FerventForerunnerYellow{},
	ids.FerventForerunnerBlue:   cards.FerventForerunnerBlue{},

	ids.FiddlersGreenRed:    cards.FiddlersGreenRed{},
	ids.FiddlersGreenYellow: cards.FiddlersGreenYellow{},
	ids.FiddlersGreenBlue:   cards.FiddlersGreenBlue{},

	ids.FlexRed:    cards.FlexRed{},
	ids.FlexYellow: cards.FlexYellow{},
	ids.FlexBlue:   cards.FlexBlue{},

	ids.FlockOfTheFeatherWalkersRed:    cards.FlockOfTheFeatherWalkersRed{},
	ids.FlockOfTheFeatherWalkersYellow: cards.FlockOfTheFeatherWalkersYellow{},
	ids.FlockOfTheFeatherWalkersBlue:   cards.FlockOfTheFeatherWalkersBlue{},

	ids.FlyingHighRed:    cards.FlyingHighRed{},
	ids.FlyingHighYellow: cards.FlyingHighYellow{},
	ids.FlyingHighBlue:   cards.FlyingHighBlue{},

	ids.FoolsGoldYellow: cards.FoolsGoldYellow{},

	ids.ForceSightRed:    cards.ForceSightRed{},
	ids.ForceSightYellow: cards.ForceSightYellow{},
	ids.ForceSightBlue:   cards.ForceSightBlue{},

	ids.FreewheelingRenegadesRed:    cards.FreewheelingRenegadesRed{},
	ids.FreewheelingRenegadesYellow: cards.FreewheelingRenegadesYellow{},
	ids.FreewheelingRenegadesBlue:   cards.FreewheelingRenegadesBlue{},

	ids.FrontlineScoutRed:    cards.FrontlineScoutRed{},
	ids.FrontlineScoutYellow: cards.FrontlineScoutYellow{},
	ids.FrontlineScoutBlue:   cards.FrontlineScoutBlue{},

	ids.FyendalsFightingSpiritRed:    cards.FyendalsFightingSpiritRed{},
	ids.FyendalsFightingSpiritYellow: cards.FyendalsFightingSpiritYellow{},
	ids.FyendalsFightingSpiritBlue:   cards.FyendalsFightingSpiritBlue{},

	ids.GravekeepingRed:    cards.GravekeepingRed{},
	ids.GravekeepingYellow: cards.GravekeepingYellow{},
	ids.GravekeepingBlue:   cards.GravekeepingBlue{},

	ids.HandBehindThePenRed: cards.HandBehindThePenRed{},

	ids.HealingBalmRed:    cards.HealingBalmRed{},
	ids.HealingBalmYellow: cards.HealingBalmYellow{},
	ids.HealingBalmBlue:   cards.HealingBalmBlue{},

	ids.HealingPotionBlue: cards.HealingPotionBlue{},

	ids.HighStrikerRed:    cards.HighStrikerRed{},
	ids.HighStrikerYellow: cards.HighStrikerYellow{},
	ids.HighStrikerBlue:   cards.HighStrikerBlue{},

	ids.HumbleRed:    cards.HumbleRed{},
	ids.HumbleYellow: cards.HumbleYellow{},
	ids.HumbleBlue:   cards.HumbleBlue{},

	ids.ImperialSealOfCommandRed: cards.ImperialSealOfCommandRed{},

	ids.InfectiousHostRed:    cards.InfectiousHostRed{},
	ids.InfectiousHostYellow: cards.InfectiousHostYellow{},
	ids.InfectiousHostBlue:   cards.InfectiousHostBlue{},

	ids.JackBeNimbleRed: cards.JackBeNimbleRed{},

	ids.JackBeQuickRed: cards.JackBeQuickRed{},

	ids.LayLowYellow: cards.LayLowYellow{},

	ids.LeadTheChargeRed:    cards.LeadTheChargeRed{},
	ids.LeadTheChargeYellow: cards.LeadTheChargeYellow{},
	ids.LeadTheChargeBlue:   cards.LeadTheChargeBlue{},

	ids.LifeForALifeRed:    cards.LifeForALifeRed{},
	ids.LifeForALifeYellow: cards.LifeForALifeYellow{},
	ids.LifeForALifeBlue:   cards.LifeForALifeBlue{},

	ids.LifeOfThePartyRed:    cards.LifeOfThePartyRed{},
	ids.LifeOfThePartyYellow: cards.LifeOfThePartyYellow{},
	ids.LifeOfThePartyBlue:   cards.LifeOfThePartyBlue{},

	ids.LookingForAScrapRed:    cards.LookingForAScrapRed{},
	ids.LookingForAScrapYellow: cards.LookingForAScrapYellow{},
	ids.LookingForAScrapBlue:   cards.LookingForAScrapBlue{},

	ids.LookTuffRed: cards.LookTuffRed{},

	ids.LungingPressBlue: cards.LungingPressBlue{},

	ids.MemorialGroundRed:    cards.MemorialGroundRed{},
	ids.MemorialGroundYellow: cards.MemorialGroundYellow{},
	ids.MemorialGroundBlue:   cards.MemorialGroundBlue{},

	ids.MinnowismRed:    cards.MinnowismRed{},
	ids.MinnowismYellow: cards.MinnowismYellow{},
	ids.MinnowismBlue:   cards.MinnowismBlue{},

	ids.MoneyOrYourLifeRed:    cards.MoneyOrYourLifeRed{},
	ids.MoneyOrYourLifeYellow: cards.MoneyOrYourLifeYellow{},
	ids.MoneyOrYourLifeBlue:   cards.MoneyOrYourLifeBlue{},

	ids.MoneyWhereYaMouthIsRed:    cards.MoneyWhereYaMouthIsRed{},
	ids.MoneyWhereYaMouthIsYellow: cards.MoneyWhereYaMouthIsYellow{},
	ids.MoneyWhereYaMouthIsBlue:   cards.MoneyWhereYaMouthIsBlue{},

	ids.MoonWishRed:    cards.MoonWishRed{},
	ids.MoonWishYellow: cards.MoonWishYellow{},
	ids.MoonWishBlue:   cards.MoonWishBlue{},

	ids.MuscleMuttYellow: cards.MuscleMuttYellow{},

	ids.NimbleStrikeRed:    cards.NimbleStrikeRed{},
	ids.NimbleStrikeYellow: cards.NimbleStrikeYellow{},
	ids.NimbleStrikeBlue:   cards.NimbleStrikeBlue{},

	ids.NimblismRed:    cards.NimblismRed{},
	ids.NimblismYellow: cards.NimblismYellow{},
	ids.NimblismBlue:   cards.NimblismBlue{},

	ids.NimbyRed:    cards.NimbyRed{},
	ids.NimbyYellow: cards.NimbyYellow{},
	ids.NimbyBlue:   cards.NimbyBlue{},

	ids.NipAtTheHeelsBlue: cards.NipAtTheHeelsBlue{},

	ids.OasisRespiteRed:    cards.OasisRespiteRed{},
	ids.OasisRespiteYellow: cards.OasisRespiteYellow{},
	ids.OasisRespiteBlue:   cards.OasisRespiteBlue{},

	ids.OnAKnifeEdgeYellow: cards.OnAKnifeEdgeYellow{},

	ids.OnTheHorizonRed:    cards.OnTheHorizonRed{},
	ids.OnTheHorizonYellow: cards.OnTheHorizonYellow{},
	ids.OnTheHorizonBlue:   cards.OnTheHorizonBlue{},

	ids.OutedRed: cards.OutedRed{},

	ids.OutMuscleRed:    cards.OutMuscleRed{},
	ids.OutMuscleYellow: cards.OutMuscleYellow{},
	ids.OutMuscleBlue:   cards.OutMuscleBlue{},

	ids.OverloadRed:    cards.OverloadRed{},
	ids.OverloadYellow: cards.OverloadYellow{},
	ids.OverloadBlue:   cards.OverloadBlue{},

	ids.PeaceOfMindRed:    cards.PeaceOfMindRed{},
	ids.PeaceOfMindYellow: cards.PeaceOfMindYellow{},
	ids.PeaceOfMindBlue:   cards.PeaceOfMindBlue{},

	ids.PerformanceBonusRed:    cards.PerformanceBonusRed{},
	ids.PerformanceBonusYellow: cards.PerformanceBonusYellow{},
	ids.PerformanceBonusBlue:   cards.PerformanceBonusBlue{},

	ids.PickACardAnyCardRed:    cards.PickACardAnyCardRed{},
	ids.PickACardAnyCardYellow: cards.PickACardAnyCardYellow{},
	ids.PickACardAnyCardBlue:   cards.PickACardAnyCardBlue{},

	ids.PilferTheTombBlue: cards.PilferTheTombBlue{},

	ids.PlunderRunRed:    cards.PlunderRunRed{},
	ids.PlunderRunYellow: cards.PlunderRunYellow{},
	ids.PlunderRunBlue:   cards.PlunderRunBlue{},

	ids.PotionOfDejaVuBlue: cards.PotionOfDejaVuBlue{},

	ids.PotionOfIronhideBlue: cards.PotionOfIronhideBlue{},

	ids.PotionOfLuckBlue: cards.PotionOfLuckBlue{},

	ids.PotionOfSeeingBlue: cards.PotionOfSeeingBlue{},

	ids.PotionOfStrengthBlue: cards.PotionOfStrengthBlue{},

	ids.PoundForPoundRed:    cards.PoundForPoundRed{},
	ids.PoundForPoundYellow: cards.PoundForPoundYellow{},
	ids.PoundForPoundBlue:   cards.PoundForPoundBlue{},

	ids.PrimeTheCrowdRed:    cards.PrimeTheCrowdRed{},
	ids.PrimeTheCrowdYellow: cards.PrimeTheCrowdYellow{},
	ids.PrimeTheCrowdBlue:   cards.PrimeTheCrowdBlue{},

	ids.PromiseOfPlentyRed:    cards.PromiseOfPlentyRed{},
	ids.PromiseOfPlentyYellow: cards.PromiseOfPlentyYellow{},
	ids.PromiseOfPlentyBlue:   cards.PromiseOfPlentyBlue{},

	ids.PublicBountyRed:    cards.PublicBountyRed{},
	ids.PublicBountyYellow: cards.PublicBountyYellow{},
	ids.PublicBountyBlue:   cards.PublicBountyBlue{},

	ids.PummelRed:    cards.PummelRed{},
	ids.PummelYellow: cards.PummelYellow{},
	ids.PummelBlue:   cards.PummelBlue{},

	ids.PunchAboveYourWeightRed:    cards.PunchAboveYourWeightRed{},
	ids.PunchAboveYourWeightYellow: cards.PunchAboveYourWeightYellow{},
	ids.PunchAboveYourWeightBlue:   cards.PunchAboveYourWeightBlue{},

	ids.PursueToTheEdgeOfOblivionRed: cards.PursueToTheEdgeOfOblivionRed{},

	ids.PursueToThePitsOfDespairRed: cards.PursueToThePitsOfDespairRed{},

	ids.PushThePointRed:    cards.PushThePointRed{},
	ids.PushThePointYellow: cards.PushThePointYellow{},
	ids.PushThePointBlue:   cards.PushThePointBlue{},

	ids.PutInContextBlue: cards.PutInContextBlue{},

	ids.RagingOnslaughtRed:    cards.RagingOnslaughtRed{},
	ids.RagingOnslaughtYellow: cards.RagingOnslaughtYellow{},
	ids.RagingOnslaughtBlue:   cards.RagingOnslaughtBlue{},

	ids.RallyTheCoastGuardRed:    cards.RallyTheCoastGuardRed{},
	ids.RallyTheCoastGuardYellow: cards.RallyTheCoastGuardYellow{},
	ids.RallyTheCoastGuardBlue:   cards.RallyTheCoastGuardBlue{},

	ids.RallyTheRearguardRed:    cards.RallyTheRearguardRed{},
	ids.RallyTheRearguardYellow: cards.RallyTheRearguardYellow{},
	ids.RallyTheRearguardBlue:   cards.RallyTheRearguardBlue{},

	ids.RansackAndRazeBlue: cards.RansackAndRazeBlue{},

	ids.RavenousRabbleRed:    cards.RavenousRabbleRed{},
	ids.RavenousRabbleYellow: cards.RavenousRabbleYellow{},
	ids.RavenousRabbleBlue:   cards.RavenousRabbleBlue{},

	ids.RazorReflexRed:    cards.RazorReflexRed{},
	ids.RazorReflexYellow: cards.RazorReflexYellow{},
	ids.RazorReflexBlue:   cards.RazorReflexBlue{},

	ids.RegainComposureBlue: cards.RegainComposureBlue{},

	ids.RegurgitatingSlogRed:    cards.RegurgitatingSlogRed{},
	ids.RegurgitatingSlogYellow: cards.RegurgitatingSlogYellow{},
	ids.RegurgitatingSlogBlue:   cards.RegurgitatingSlogBlue{},

	ids.ReinforceTheLineRed:    cards.ReinforceTheLineRed{},
	ids.ReinforceTheLineYellow: cards.ReinforceTheLineYellow{},
	ids.ReinforceTheLineBlue:   cards.ReinforceTheLineBlue{},

	ids.RelentlessPursuitBlue: cards.RelentlessPursuitBlue{},

	ids.RestvineElixirRed: cards.RestvineElixirRed{},

	ids.RiftingRed:    cards.RiftingRed{},
	ids.RiftingYellow: cards.RiftingYellow{},
	ids.RiftingBlue:   cards.RiftingBlue{},

	ids.RightBehindYouRed:    cards.RightBehindYouRed{},
	ids.RightBehindYouYellow: cards.RightBehindYouYellow{},
	ids.RightBehindYouBlue:   cards.RightBehindYouBlue{},

	ids.RiseAboveRed:    cards.RiseAboveRed{},
	ids.RiseAboveYellow: cards.RiseAboveYellow{},
	ids.RiseAboveBlue:   cards.RiseAboveBlue{},

	ids.SapwoodElixirRed: cards.SapwoodElixirRed{},

	ids.ScarForAScarRed:    cards.ScarForAScarRed{},
	ids.ScarForAScarYellow: cards.ScarForAScarYellow{},
	ids.ScarForAScarBlue:   cards.ScarForAScarBlue{},

	ids.ScourTheBattlescapeRed:    cards.ScourTheBattlescapeRed{},
	ids.ScourTheBattlescapeYellow: cards.ScourTheBattlescapeYellow{},
	ids.ScourTheBattlescapeBlue:   cards.ScourTheBattlescapeBlue{},

	ids.ScoutThePeripheryRed:    cards.ScoutThePeripheryRed{},
	ids.ScoutThePeripheryYellow: cards.ScoutThePeripheryYellow{},
	ids.ScoutThePeripheryBlue:   cards.ScoutThePeripheryBlue{},

	ids.SeekHorizonRed:    cards.SeekHorizonRed{},
	ids.SeekHorizonYellow: cards.SeekHorizonYellow{},
	ids.SeekHorizonBlue:   cards.SeekHorizonBlue{},

	ids.ShatterSorceryBlue: cards.ShatterSorceryBlue{},

	ids.SiftRed:    cards.SiftRed{},
	ids.SiftYellow: cards.SiftYellow{},
	ids.SiftBlue:   cards.SiftBlue{},

	ids.SigilOfCyclesBlue: cards.SigilOfCyclesBlue{},

	ids.SigilOfFyendalBlue: cards.SigilOfFyendalBlue{},

	ids.SigilOfProtectionRed:    cards.SigilOfProtectionRed{},
	ids.SigilOfProtectionYellow: cards.SigilOfProtectionYellow{},
	ids.SigilOfProtectionBlue:   cards.SigilOfProtectionBlue{},

	ids.SigilOfSolaceRed:    cards.SigilOfSolaceRed{},
	ids.SigilOfSolaceYellow: cards.SigilOfSolaceYellow{},
	ids.SigilOfSolaceBlue:   cards.SigilOfSolaceBlue{},

	ids.SinkBelowRed:    cards.SinkBelowRed{},
	ids.SinkBelowYellow: cards.SinkBelowYellow{},
	ids.SinkBelowBlue:   cards.SinkBelowBlue{},

	ids.SirensOfSafeHarborRed:    cards.SirensOfSafeHarborRed{},
	ids.SirensOfSafeHarborYellow: cards.SirensOfSafeHarborYellow{},
	ids.SirensOfSafeHarborBlue:   cards.SirensOfSafeHarborBlue{},

	ids.SloggismRed:    cards.SloggismRed{},
	ids.SloggismYellow: cards.SloggismYellow{},
	ids.SloggismBlue:   cards.SloggismBlue{},

	ids.SmashingGoodTimeRed:    cards.SmashingGoodTimeRed{},
	ids.SmashingGoodTimeYellow: cards.SmashingGoodTimeYellow{},
	ids.SmashingGoodTimeBlue:   cards.SmashingGoodTimeBlue{},

	ids.SmashUpRed: cards.SmashUpRed{},

	ids.SnatchRed:    cards.SnatchRed{},
	ids.SnatchYellow: cards.SnatchYellow{},
	ids.SnatchBlue:   cards.SnatchBlue{},

	ids.SoundTheAlarmRed: cards.SoundTheAlarmRed{},

	ids.SpringboardSomersaultYellow: cards.SpringboardSomersaultYellow{},

	ids.SpringLoadRed:    cards.SpringLoadRed{},
	ids.SpringLoadYellow: cards.SpringLoadYellow{},
	ids.SpringLoadBlue:   cards.SpringLoadBlue{},

	ids.StartingStakeYellow: cards.StartingStakeYellow{},

	ids.StonyWoottonhogRed:    cards.StonyWoottonhogRed{},
	ids.StonyWoottonhogYellow: cards.StonyWoottonhogYellow{},
	ids.StonyWoottonhogBlue:   cards.StonyWoottonhogBlue{},

	ids.StrategicPlanningRed:    cards.StrategicPlanningRed{},
	ids.StrategicPlanningYellow: cards.StrategicPlanningYellow{},
	ids.StrategicPlanningBlue:   cards.StrategicPlanningBlue{},

	ids.StrikeGoldRed:    cards.StrikeGoldRed{},
	ids.StrikeGoldYellow: cards.StrikeGoldYellow{},
	ids.StrikeGoldBlue:   cards.StrikeGoldBlue{},

	ids.SunKissRed:    cards.SunKissRed{},
	ids.SunKissYellow: cards.SunKissYellow{},
	ids.SunKissBlue:   cards.SunKissBlue{},

	ids.SurgingMilitiaRed:    cards.SurgingMilitiaRed{},
	ids.SurgingMilitiaYellow: cards.SurgingMilitiaYellow{},
	ids.SurgingMilitiaBlue:   cards.SurgingMilitiaBlue{},

	ids.TalismanOfBalanceBlue: cards.TalismanOfBalanceBlue{},

	ids.TalismanOfCremationBlue: cards.TalismanOfCremationBlue{},

	ids.TalismanOfDousingYellow: cards.TalismanOfDousingYellow{},

	ids.TalismanOfFeatherfootYellow: cards.TalismanOfFeatherfootYellow{},

	ids.TalismanOfRecompenseYellow: cards.TalismanOfRecompenseYellow{},

	ids.TalismanOfTithesBlue: cards.TalismanOfTithesBlue{},

	ids.TalismanOfWarfareYellow: cards.TalismanOfWarfareYellow{},

	ids.TestOfStrengthRed: cards.TestOfStrengthRed{},

	ids.ThrustRed: cards.ThrustRed{},

	ids.TimesnapPotionBlue: cards.TimesnapPotionBlue{},

	ids.TipOffRed:    cards.TipOffRed{},
	ids.TipOffYellow: cards.TipOffYellow{},
	ids.TipOffBlue:   cards.TipOffBlue{},

	ids.TitaniumBaubleBlue: cards.TitaniumBaubleBlue{},

	ids.TitForTatBlue: cards.TitForTatBlue{},

	ids.TongueTiedRed: cards.TongueTiedRed{},

	ids.ToughenUpBlue: cards.ToughenUpBlue{},

	ids.TradeInRed:    cards.TradeInRed{},
	ids.TradeInYellow: cards.TradeInYellow{},
	ids.TradeInBlue:   cards.TradeInBlue{},

	ids.TremorOfIArathaelRed:    cards.TremorOfIArathaelRed{},
	ids.TremorOfIArathaelYellow: cards.TremorOfIArathaelYellow{},
	ids.TremorOfIArathaelBlue:   cards.TremorOfIArathaelBlue{},

	ids.TrotAlongBlue: cards.TrotAlongBlue{},

	ids.UnmovableRed:    cards.UnmovableRed{},
	ids.UnmovableYellow: cards.UnmovableYellow{},
	ids.UnmovableBlue:   cards.UnmovableBlue{},

	ids.VigorRushRed:    cards.VigorRushRed{},
	ids.VigorRushYellow: cards.VigorRushYellow{},
	ids.VigorRushBlue:   cards.VigorRushBlue{},

	ids.VisitTheBlacksmithBlue: cards.VisitTheBlacksmithBlue{},

	ids.WageGoldRed:    cards.WageGoldRed{},
	ids.WageGoldYellow: cards.WageGoldYellow{},
	ids.WageGoldBlue:   cards.WageGoldBlue{},

	ids.WalkThePlankRed:    cards.WalkThePlankRed{},
	ids.WalkThePlankYellow: cards.WalkThePlankYellow{},
	ids.WalkThePlankBlue:   cards.WalkThePlankBlue{},

	ids.WarmongersRecitalRed:    cards.WarmongersRecitalRed{},
	ids.WarmongersRecitalYellow: cards.WarmongersRecitalYellow{},
	ids.WarmongersRecitalBlue:   cards.WarmongersRecitalBlue{},

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

	ids.WreckHavocRed:    cards.WreckHavocRed{},
	ids.WreckHavocYellow: cards.WreckHavocYellow{},
	ids.WreckHavocBlue:   cards.WreckHavocBlue{},

	ids.YintiYantiRed:    cards.YintiYantiRed{},
	ids.YintiYantiYellow: cards.YintiYantiYellow{},
	ids.YintiYantiBlue:   cards.YintiYantiBlue{},

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

// GetCard returns the card for the given ID. Panics if id is Invalid or out of range —
// callers should only pass IDs they got from this package.
func GetCard(id CardID) sim.Card {
	if id == ids.InvalidCard || int(id) >= len(cardsByID) || cardsByID[id] == nil {
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
// Freshly allocated; safe to mutate. The fake card IDs (ids.FakeRedAttack, …) are
// deliberately not in the registry, so this is just AllCards() under a different name —
// kept distinct so callers who want "deck-legal cards" stay readable even if the registry
// ever holds non-deckable entries again.
func DeckableCards() []CardID {
	return AllCards()
}
