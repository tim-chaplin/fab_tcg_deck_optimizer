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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards/notimplemented"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards/unplayable"
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

	ids.CondemnToSlaughterRed:    notimplemented.CondemnToSlaughterRed{},
	ids.CondemnToSlaughterYellow: notimplemented.CondemnToSlaughterYellow{},
	ids.CondemnToSlaughterBlue:   notimplemented.CondemnToSlaughterBlue{},

	ids.ConsumingVolitionRed:    cards.ConsumingVolitionRed{},
	ids.ConsumingVolitionYellow: cards.ConsumingVolitionYellow{},
	ids.ConsumingVolitionBlue:   cards.ConsumingVolitionBlue{},

	ids.DeathlyDuetRed:    cards.DeathlyDuetRed{},
	ids.DeathlyDuetYellow: cards.DeathlyDuetYellow{},
	ids.DeathlyDuetBlue:   cards.DeathlyDuetBlue{},

	ids.DrawnToTheDarkDimensionRed:    cards.DrawnToTheDarkDimensionRed{},
	ids.DrawnToTheDarkDimensionYellow: cards.DrawnToTheDarkDimensionYellow{},
	ids.DrawnToTheDarkDimensionBlue:   cards.DrawnToTheDarkDimensionBlue{},

	ids.DrowningDireRed:    notimplemented.DrowningDireRed{},
	ids.DrowningDireYellow: notimplemented.DrowningDireYellow{},
	ids.DrowningDireBlue:   notimplemented.DrowningDireBlue{},

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

	ids.SplinteringDeadwoodRed:    notimplemented.SplinteringDeadwoodRed{},
	ids.SplinteringDeadwoodYellow: notimplemented.SplinteringDeadwoodYellow{},
	ids.SplinteringDeadwoodBlue:   notimplemented.SplinteringDeadwoodBlue{},

	ids.SutcliffesResearchNotesRed:    cards.SutcliffesResearchNotesRed{},
	ids.SutcliffesResearchNotesYellow: cards.SutcliffesResearchNotesYellow{},
	ids.SutcliffesResearchNotesBlue:   cards.SutcliffesResearchNotesBlue{},

	ids.VantagePointRed:    notimplemented.VantagePointRed{},
	ids.VantagePointYellow: notimplemented.VantagePointYellow{},
	ids.VantagePointBlue:   notimplemented.VantagePointBlue{},

	ids.VexingMaliceRed:    cards.VexingMaliceRed{},
	ids.VexingMaliceYellow: cards.VexingMaliceYellow{},
	ids.VexingMaliceBlue:   cards.VexingMaliceBlue{},

	ids.WeepingBattlegroundRed:    cards.WeepingBattlegroundRed{},
	ids.WeepingBattlegroundYellow: cards.WeepingBattlegroundYellow{},
	ids.WeepingBattlegroundBlue:   cards.WeepingBattlegroundBlue{},

	ids.AdrenalineRushRed:    cards.AdrenalineRushRed{},
	ids.AdrenalineRushYellow: cards.AdrenalineRushYellow{},
	ids.AdrenalineRushBlue:   cards.AdrenalineRushBlue{},

	ids.AmuletOfAssertivenessYellow: unplayable.AmuletOfAssertivenessYellow{},

	ids.AmuletOfEchoesBlue: unplayable.AmuletOfEchoesBlue{},

	ids.AmuletOfHavencallBlue: unplayable.AmuletOfHavencallBlue{},

	ids.AmuletOfIgnitionYellow: unplayable.AmuletOfIgnitionYellow{},

	ids.AmuletOfInterventionBlue: unplayable.AmuletOfInterventionBlue{},

	ids.AmuletOfOblationBlue: unplayable.AmuletOfOblationBlue{},

	ids.ArcanePolarityRed:    notimplemented.ArcanePolarityRed{},
	ids.ArcanePolarityYellow: notimplemented.ArcanePolarityYellow{},
	ids.ArcanePolarityBlue:   notimplemented.ArcanePolarityBlue{},

	ids.BackAlleyBreaklineRed:    notimplemented.BackAlleyBreaklineRed{},
	ids.BackAlleyBreaklineYellow: notimplemented.BackAlleyBreaklineYellow{},
	ids.BackAlleyBreaklineBlue:   notimplemented.BackAlleyBreaklineBlue{},

	ids.BarragingBrawnhideRed:    notimplemented.BarragingBrawnhideRed{},
	ids.BarragingBrawnhideYellow: notimplemented.BarragingBrawnhideYellow{},
	ids.BarragingBrawnhideBlue:   notimplemented.BarragingBrawnhideBlue{},

	ids.BattlefrontBastionRed:    notimplemented.BattlefrontBastionRed{},
	ids.BattlefrontBastionYellow: notimplemented.BattlefrontBastionYellow{},
	ids.BattlefrontBastionBlue:   notimplemented.BattlefrontBastionBlue{},

	ids.BelittleRed:    notimplemented.BelittleRed{},
	ids.BelittleYellow: notimplemented.BelittleYellow{},
	ids.BelittleBlue:   notimplemented.BelittleBlue{},

	ids.BladeFlashBlue: notimplemented.BladeFlashBlue{},

	ids.BlanchRed:    notimplemented.BlanchRed{},
	ids.BlanchYellow: notimplemented.BlanchYellow{},
	ids.BlanchBlue:   notimplemented.BlanchBlue{},

	ids.BlowForABlowRed: cards.BlowForABlowRed{},

	ids.BlusterBuffRed: notimplemented.BlusterBuffRed{},

	ids.BrandishRed:    notimplemented.BrandishRed{},
	ids.BrandishYellow: notimplemented.BrandishYellow{},
	ids.BrandishBlue:   notimplemented.BrandishBlue{},

	ids.BrothersInArmsRed:    notimplemented.BrothersInArmsRed{},
	ids.BrothersInArmsYellow: notimplemented.BrothersInArmsYellow{},
	ids.BrothersInArmsBlue:   notimplemented.BrothersInArmsBlue{},

	ids.BrushOffRed:    notimplemented.BrushOffRed{},
	ids.BrushOffYellow: notimplemented.BrushOffYellow{},
	ids.BrushOffBlue:   notimplemented.BrushOffBlue{},

	ids.BrutalAssaultRed:    cards.BrutalAssaultRed{},
	ids.BrutalAssaultYellow: cards.BrutalAssaultYellow{},
	ids.BrutalAssaultBlue:   cards.BrutalAssaultBlue{},

	ids.CadaverousContrabandRed:    notimplemented.CadaverousContrabandRed{},
	ids.CadaverousContrabandYellow: notimplemented.CadaverousContrabandYellow{},
	ids.CadaverousContrabandBlue:   notimplemented.CadaverousContrabandBlue{},

	ids.CalmingBreezeRed: notimplemented.CalmingBreezeRed{},

	ids.CaptainsCallRed:    notimplemented.CaptainsCallRed{},
	ids.CaptainsCallYellow: notimplemented.CaptainsCallYellow{},
	ids.CaptainsCallBlue:   notimplemented.CaptainsCallBlue{},

	ids.CashInYellow: notimplemented.CashInYellow{},

	ids.ChestPuffRed: notimplemented.ChestPuffRed{},

	ids.ClapEmInIronsBlue: notimplemented.ClapEmInIronsBlue{},

	ids.ClarityPotionBlue: cards.ClarityPotionBlue{},

	ids.ClearwaterElixirRed: notimplemented.ClearwaterElixirRed{},

	ids.ComeToFightRed:    cards.ComeToFightRed{},
	ids.ComeToFightYellow: cards.ComeToFightYellow{},
	ids.ComeToFightBlue:   cards.ComeToFightBlue{},

	ids.CountYourBlessingsRed:    notimplemented.CountYourBlessingsRed{},
	ids.CountYourBlessingsYellow: notimplemented.CountYourBlessingsYellow{},
	ids.CountYourBlessingsBlue:   notimplemented.CountYourBlessingsBlue{},

	ids.CrackedBaubleYellow: notimplemented.CrackedBaubleYellow{},

	ids.CrashDownTheGatesRed:    notimplemented.CrashDownTheGatesRed{},
	ids.CrashDownTheGatesYellow: notimplemented.CrashDownTheGatesYellow{},
	ids.CrashDownTheGatesBlue:   notimplemented.CrashDownTheGatesBlue{},

	ids.CriticalStrikeRed:    cards.CriticalStrikeRed{},
	ids.CriticalStrikeYellow: cards.CriticalStrikeYellow{},
	ids.CriticalStrikeBlue:   cards.CriticalStrikeBlue{},

	ids.CutDownToSizeRed:    notimplemented.CutDownToSizeRed{},
	ids.CutDownToSizeYellow: notimplemented.CutDownToSizeYellow{},
	ids.CutDownToSizeBlue:   notimplemented.CutDownToSizeBlue{},

	ids.DemolitionCrewRed:    cards.DemolitionCrewRed{},
	ids.DemolitionCrewYellow: cards.DemolitionCrewYellow{},
	ids.DemolitionCrewBlue:   cards.DemolitionCrewBlue{},

	ids.DestructiveDeliberationRed:    notimplemented.DestructiveDeliberationRed{},
	ids.DestructiveDeliberationYellow: notimplemented.DestructiveDeliberationYellow{},
	ids.DestructiveDeliberationBlue:   notimplemented.DestructiveDeliberationBlue{},

	ids.DestructiveTendenciesBlue: notimplemented.DestructiveTendenciesBlue{},

	ids.DodgeBlue: cards.DodgeBlue{},

	ids.DownButNotOutRed:    notimplemented.DownButNotOutRed{},
	ids.DownButNotOutYellow: notimplemented.DownButNotOutYellow{},
	ids.DownButNotOutBlue:   notimplemented.DownButNotOutBlue{},

	ids.DragDownRed:    notimplemented.DragDownRed{},
	ids.DragDownYellow: notimplemented.DragDownYellow{},
	ids.DragDownBlue:   notimplemented.DragDownBlue{},

	ids.DroneOfBrutalityRed:    notimplemented.DroneOfBrutalityRed{},
	ids.DroneOfBrutalityYellow: notimplemented.DroneOfBrutalityYellow{},
	ids.DroneOfBrutalityBlue:   notimplemented.DroneOfBrutalityBlue{},

	ids.EirinasPrayerRed:    notimplemented.EirinasPrayerRed{},
	ids.EirinasPrayerYellow: notimplemented.EirinasPrayerYellow{},
	ids.EirinasPrayerBlue:   notimplemented.EirinasPrayerBlue{},

	ids.EmissaryOfMoonRed: notimplemented.EmissaryOfMoonRed{},

	ids.EmissaryOfTidesRed: notimplemented.EmissaryOfTidesRed{},

	ids.EmissaryOfWindRed: notimplemented.EmissaryOfWindRed{},

	ids.EnchantingMelodyRed:    notimplemented.EnchantingMelodyRed{},
	ids.EnchantingMelodyYellow: notimplemented.EnchantingMelodyYellow{},
	ids.EnchantingMelodyBlue:   notimplemented.EnchantingMelodyBlue{},

	ids.EnergyPotionBlue: notimplemented.EnergyPotionBlue{},

	ids.EvasiveLeapRed:    cards.EvasiveLeapRed{},
	ids.EvasiveLeapYellow: cards.EvasiveLeapYellow{},
	ids.EvasiveLeapBlue:   cards.EvasiveLeapBlue{},

	ids.EvenBiggerThanThatRed:    notimplemented.EvenBiggerThanThatRed{},
	ids.EvenBiggerThanThatYellow: notimplemented.EvenBiggerThanThatYellow{},
	ids.EvenBiggerThanThatBlue:   notimplemented.EvenBiggerThanThatBlue{},

	ids.ExposedBlue: notimplemented.ExposedBlue{},

	ids.FactFindingMissionRed:    notimplemented.FactFindingMissionRed{},
	ids.FactFindingMissionYellow: notimplemented.FactFindingMissionYellow{},
	ids.FactFindingMissionBlue:   notimplemented.FactFindingMissionBlue{},

	ids.FateForeseenRed:    cards.FateForeseenRed{},
	ids.FateForeseenYellow: cards.FateForeseenYellow{},
	ids.FateForeseenBlue:   cards.FateForeseenBlue{},

	ids.FeistyLocalsRed:    notimplemented.FeistyLocalsRed{},
	ids.FeistyLocalsYellow: notimplemented.FeistyLocalsYellow{},
	ids.FeistyLocalsBlue:   notimplemented.FeistyLocalsBlue{},

	ids.FerventForerunnerRed:    cards.FerventForerunnerRed{},
	ids.FerventForerunnerYellow: cards.FerventForerunnerYellow{},
	ids.FerventForerunnerBlue:   cards.FerventForerunnerBlue{},

	ids.FiddlersGreenRed:    cards.FiddlersGreenRed{},
	ids.FiddlersGreenYellow: cards.FiddlersGreenYellow{},
	ids.FiddlersGreenBlue:   cards.FiddlersGreenBlue{},

	ids.FlexRed:    notimplemented.FlexRed{},
	ids.FlexYellow: notimplemented.FlexYellow{},
	ids.FlexBlue:   notimplemented.FlexBlue{},

	ids.FlockOfTheFeatherWalkersRed:    notimplemented.FlockOfTheFeatherWalkersRed{},
	ids.FlockOfTheFeatherWalkersYellow: notimplemented.FlockOfTheFeatherWalkersYellow{},
	ids.FlockOfTheFeatherWalkersBlue:   notimplemented.FlockOfTheFeatherWalkersBlue{},

	ids.FlyingHighRed:    cards.FlyingHighRed{},
	ids.FlyingHighYellow: cards.FlyingHighYellow{},
	ids.FlyingHighBlue:   cards.FlyingHighBlue{},

	ids.FoolsGoldYellow: notimplemented.FoolsGoldYellow{},

	ids.ForceSightRed:    cards.ForceSightRed{},
	ids.ForceSightYellow: cards.ForceSightYellow{},
	ids.ForceSightBlue:   cards.ForceSightBlue{},

	ids.FreewheelingRenegadesRed:    notimplemented.FreewheelingRenegadesRed{},
	ids.FreewheelingRenegadesYellow: notimplemented.FreewheelingRenegadesYellow{},
	ids.FreewheelingRenegadesBlue:   notimplemented.FreewheelingRenegadesBlue{},

	ids.FrontlineScoutRed:    notimplemented.FrontlineScoutRed{},
	ids.FrontlineScoutYellow: notimplemented.FrontlineScoutYellow{},
	ids.FrontlineScoutBlue:   notimplemented.FrontlineScoutBlue{},

	ids.FyendalsFightingSpiritRed:    cards.FyendalsFightingSpiritRed{},
	ids.FyendalsFightingSpiritYellow: cards.FyendalsFightingSpiritYellow{},
	ids.FyendalsFightingSpiritBlue:   cards.FyendalsFightingSpiritBlue{},

	ids.GravekeepingRed:    notimplemented.GravekeepingRed{},
	ids.GravekeepingYellow: notimplemented.GravekeepingYellow{},
	ids.GravekeepingBlue:   notimplemented.GravekeepingBlue{},

	ids.HandBehindThePenRed: notimplemented.HandBehindThePenRed{},

	ids.HealingBalmRed:    cards.HealingBalmRed{},
	ids.HealingBalmYellow: cards.HealingBalmYellow{},
	ids.HealingBalmBlue:   cards.HealingBalmBlue{},

	ids.HealingPotionBlue: notimplemented.HealingPotionBlue{},

	ids.HighStrikerRed:    notimplemented.HighStrikerRed{},
	ids.HighStrikerYellow: notimplemented.HighStrikerYellow{},
	ids.HighStrikerBlue:   notimplemented.HighStrikerBlue{},

	ids.HumbleRed:    notimplemented.HumbleRed{},
	ids.HumbleYellow: notimplemented.HumbleYellow{},
	ids.HumbleBlue:   notimplemented.HumbleBlue{},

	ids.ImperialSealOfCommandRed: notimplemented.ImperialSealOfCommandRed{},

	ids.InfectiousHostRed:    notimplemented.InfectiousHostRed{},
	ids.InfectiousHostYellow: notimplemented.InfectiousHostYellow{},
	ids.InfectiousHostBlue:   notimplemented.InfectiousHostBlue{},

	ids.JackBeNimbleRed: notimplemented.JackBeNimbleRed{},

	ids.JackBeQuickRed: notimplemented.JackBeQuickRed{},

	ids.LayLowYellow: notimplemented.LayLowYellow{},

	ids.LeadTheChargeRed:    notimplemented.LeadTheChargeRed{},
	ids.LeadTheChargeYellow: notimplemented.LeadTheChargeYellow{},
	ids.LeadTheChargeBlue:   notimplemented.LeadTheChargeBlue{},

	ids.LifeForALifeRed:    cards.LifeForALifeRed{},
	ids.LifeForALifeYellow: cards.LifeForALifeYellow{},
	ids.LifeForALifeBlue:   cards.LifeForALifeBlue{},

	ids.LifeOfThePartyRed:    notimplemented.LifeOfThePartyRed{},
	ids.LifeOfThePartyYellow: notimplemented.LifeOfThePartyYellow{},
	ids.LifeOfThePartyBlue:   notimplemented.LifeOfThePartyBlue{},

	ids.LookingForAScrapRed:    notimplemented.LookingForAScrapRed{},
	ids.LookingForAScrapYellow: notimplemented.LookingForAScrapYellow{},
	ids.LookingForAScrapBlue:   notimplemented.LookingForAScrapBlue{},

	ids.LookTuffRed: notimplemented.LookTuffRed{},

	ids.LungingPressBlue: notimplemented.LungingPressBlue{},

	ids.MemorialGroundRed:    notimplemented.MemorialGroundRed{},
	ids.MemorialGroundYellow: notimplemented.MemorialGroundYellow{},
	ids.MemorialGroundBlue:   notimplemented.MemorialGroundBlue{},

	ids.MinnowismRed:    cards.MinnowismRed{},
	ids.MinnowismYellow: cards.MinnowismYellow{},
	ids.MinnowismBlue:   cards.MinnowismBlue{},

	ids.MoneyOrYourLifeRed:    notimplemented.MoneyOrYourLifeRed{},
	ids.MoneyOrYourLifeYellow: notimplemented.MoneyOrYourLifeYellow{},
	ids.MoneyOrYourLifeBlue:   notimplemented.MoneyOrYourLifeBlue{},

	ids.MoneyWhereYaMouthIsRed:    notimplemented.MoneyWhereYaMouthIsRed{},
	ids.MoneyWhereYaMouthIsYellow: notimplemented.MoneyWhereYaMouthIsYellow{},
	ids.MoneyWhereYaMouthIsBlue:   notimplemented.MoneyWhereYaMouthIsBlue{},

	ids.MoonWishRed:    cards.MoonWishRed{},
	ids.MoonWishYellow: cards.MoonWishYellow{},
	ids.MoonWishBlue:   cards.MoonWishBlue{},

	ids.MuscleMuttYellow: cards.MuscleMuttYellow{},

	ids.NimbleStrikeRed:    notimplemented.NimbleStrikeRed{},
	ids.NimbleStrikeYellow: notimplemented.NimbleStrikeYellow{},
	ids.NimbleStrikeBlue:   notimplemented.NimbleStrikeBlue{},

	ids.NimblismRed:    cards.NimblismRed{},
	ids.NimblismYellow: cards.NimblismYellow{},
	ids.NimblismBlue:   cards.NimblismBlue{},

	ids.NimbyRed:    notimplemented.NimbyRed{},
	ids.NimbyYellow: notimplemented.NimbyYellow{},
	ids.NimbyBlue:   notimplemented.NimbyBlue{},

	ids.NipAtTheHeelsBlue: notimplemented.NipAtTheHeelsBlue{},

	ids.OasisRespiteRed:    notimplemented.OasisRespiteRed{},
	ids.OasisRespiteYellow: notimplemented.OasisRespiteYellow{},
	ids.OasisRespiteBlue:   notimplemented.OasisRespiteBlue{},

	ids.OnAKnifeEdgeYellow: notimplemented.OnAKnifeEdgeYellow{},

	ids.OnTheHorizonRed:    cards.OnTheHorizonRed{},
	ids.OnTheHorizonYellow: cards.OnTheHorizonYellow{},
	ids.OnTheHorizonBlue:   cards.OnTheHorizonBlue{},

	ids.OutedRed: notimplemented.OutedRed{},

	ids.OutMuscleRed:    notimplemented.OutMuscleRed{},
	ids.OutMuscleYellow: notimplemented.OutMuscleYellow{},
	ids.OutMuscleBlue:   notimplemented.OutMuscleBlue{},

	ids.OverloadRed:    notimplemented.OverloadRed{},
	ids.OverloadYellow: notimplemented.OverloadYellow{},
	ids.OverloadBlue:   notimplemented.OverloadBlue{},

	ids.PeaceOfMindRed:    notimplemented.PeaceOfMindRed{},
	ids.PeaceOfMindYellow: notimplemented.PeaceOfMindYellow{},
	ids.PeaceOfMindBlue:   notimplemented.PeaceOfMindBlue{},

	ids.PerformanceBonusRed:    notimplemented.PerformanceBonusRed{},
	ids.PerformanceBonusYellow: notimplemented.PerformanceBonusYellow{},
	ids.PerformanceBonusBlue:   notimplemented.PerformanceBonusBlue{},

	ids.PickACardAnyCardRed:    notimplemented.PickACardAnyCardRed{},
	ids.PickACardAnyCardYellow: notimplemented.PickACardAnyCardYellow{},
	ids.PickACardAnyCardBlue:   notimplemented.PickACardAnyCardBlue{},

	ids.PilferTheTombBlue: notimplemented.PilferTheTombBlue{},

	ids.PlunderRunRed:    notimplemented.PlunderRunRed{},
	ids.PlunderRunYellow: notimplemented.PlunderRunYellow{},
	ids.PlunderRunBlue:   notimplemented.PlunderRunBlue{},

	ids.PotionOfDejaVuBlue: notimplemented.PotionOfDejaVuBlue{},

	ids.PotionOfIronhideBlue: notimplemented.PotionOfIronhideBlue{},

	ids.PotionOfLuckBlue: notimplemented.PotionOfLuckBlue{},

	ids.PotionOfSeeingBlue: unplayable.PotionOfSeeingBlue{},

	ids.PotionOfStrengthBlue: notimplemented.PotionOfStrengthBlue{},

	ids.PoundForPoundRed:    cards.PoundForPoundRed{},
	ids.PoundForPoundYellow: cards.PoundForPoundYellow{},
	ids.PoundForPoundBlue:   cards.PoundForPoundBlue{},

	ids.PrimeTheCrowdRed:    notimplemented.PrimeTheCrowdRed{},
	ids.PrimeTheCrowdYellow: notimplemented.PrimeTheCrowdYellow{},
	ids.PrimeTheCrowdBlue:   notimplemented.PrimeTheCrowdBlue{},

	ids.PromiseOfPlentyRed:    notimplemented.PromiseOfPlentyRed{},
	ids.PromiseOfPlentyYellow: notimplemented.PromiseOfPlentyYellow{},
	ids.PromiseOfPlentyBlue:   notimplemented.PromiseOfPlentyBlue{},

	ids.PublicBountyRed:    notimplemented.PublicBountyRed{},
	ids.PublicBountyYellow: notimplemented.PublicBountyYellow{},
	ids.PublicBountyBlue:   notimplemented.PublicBountyBlue{},

	ids.PummelRed:    notimplemented.PummelRed{},
	ids.PummelYellow: notimplemented.PummelYellow{},
	ids.PummelBlue:   notimplemented.PummelBlue{},

	ids.PunchAboveYourWeightRed:    notimplemented.PunchAboveYourWeightRed{},
	ids.PunchAboveYourWeightYellow: notimplemented.PunchAboveYourWeightYellow{},
	ids.PunchAboveYourWeightBlue:   notimplemented.PunchAboveYourWeightBlue{},

	ids.PursueToTheEdgeOfOblivionRed: notimplemented.PursueToTheEdgeOfOblivionRed{},

	ids.PursueToThePitsOfDespairRed: notimplemented.PursueToThePitsOfDespairRed{},

	ids.PushThePointRed:    notimplemented.PushThePointRed{},
	ids.PushThePointYellow: notimplemented.PushThePointYellow{},
	ids.PushThePointBlue:   notimplemented.PushThePointBlue{},

	ids.PutInContextBlue: notimplemented.PutInContextBlue{},

	ids.RagingOnslaughtRed:    cards.RagingOnslaughtRed{},
	ids.RagingOnslaughtYellow: cards.RagingOnslaughtYellow{},
	ids.RagingOnslaughtBlue:   cards.RagingOnslaughtBlue{},

	ids.RallyTheCoastGuardRed:    notimplemented.RallyTheCoastGuardRed{},
	ids.RallyTheCoastGuardYellow: notimplemented.RallyTheCoastGuardYellow{},
	ids.RallyTheCoastGuardBlue:   notimplemented.RallyTheCoastGuardBlue{},

	ids.RallyTheRearguardRed:    notimplemented.RallyTheRearguardRed{},
	ids.RallyTheRearguardYellow: notimplemented.RallyTheRearguardYellow{},
	ids.RallyTheRearguardBlue:   notimplemented.RallyTheRearguardBlue{},

	ids.RansackAndRazeBlue: notimplemented.RansackAndRazeBlue{},

	ids.RavenousRabbleRed:    cards.RavenousRabbleRed{},
	ids.RavenousRabbleYellow: cards.RavenousRabbleYellow{},
	ids.RavenousRabbleBlue:   cards.RavenousRabbleBlue{},

	ids.RazorReflexRed:    notimplemented.RazorReflexRed{},
	ids.RazorReflexYellow: notimplemented.RazorReflexYellow{},
	ids.RazorReflexBlue:   notimplemented.RazorReflexBlue{},

	ids.RegainComposureBlue: notimplemented.RegainComposureBlue{},

	ids.RegurgitatingSlogRed:    notimplemented.RegurgitatingSlogRed{},
	ids.RegurgitatingSlogYellow: notimplemented.RegurgitatingSlogYellow{},
	ids.RegurgitatingSlogBlue:   notimplemented.RegurgitatingSlogBlue{},

	ids.ReinforceTheLineRed:    notimplemented.ReinforceTheLineRed{},
	ids.ReinforceTheLineYellow: notimplemented.ReinforceTheLineYellow{},
	ids.ReinforceTheLineBlue:   notimplemented.ReinforceTheLineBlue{},

	ids.RelentlessPursuitBlue: notimplemented.RelentlessPursuitBlue{},

	ids.RestvineElixirRed: notimplemented.RestvineElixirRed{},

	ids.RiftingRed:    notimplemented.RiftingRed{},
	ids.RiftingYellow: notimplemented.RiftingYellow{},
	ids.RiftingBlue:   notimplemented.RiftingBlue{},

	ids.RightBehindYouRed:    notimplemented.RightBehindYouRed{},
	ids.RightBehindYouYellow: notimplemented.RightBehindYouYellow{},
	ids.RightBehindYouBlue:   notimplemented.RightBehindYouBlue{},

	ids.RiseAboveRed:    notimplemented.RiseAboveRed{},
	ids.RiseAboveYellow: notimplemented.RiseAboveYellow{},
	ids.RiseAboveBlue:   notimplemented.RiseAboveBlue{},

	ids.SapwoodElixirRed: notimplemented.SapwoodElixirRed{},

	ids.ScarForAScarRed:    cards.ScarForAScarRed{},
	ids.ScarForAScarYellow: cards.ScarForAScarYellow{},
	ids.ScarForAScarBlue:   cards.ScarForAScarBlue{},

	ids.ScourTheBattlescapeRed:    notimplemented.ScourTheBattlescapeRed{},
	ids.ScourTheBattlescapeYellow: notimplemented.ScourTheBattlescapeYellow{},
	ids.ScourTheBattlescapeBlue:   notimplemented.ScourTheBattlescapeBlue{},

	ids.ScoutThePeripheryRed:    cards.ScoutThePeripheryRed{},
	ids.ScoutThePeripheryYellow: cards.ScoutThePeripheryYellow{},
	ids.ScoutThePeripheryBlue:   cards.ScoutThePeripheryBlue{},

	ids.SeekHorizonRed:    notimplemented.SeekHorizonRed{},
	ids.SeekHorizonYellow: notimplemented.SeekHorizonYellow{},
	ids.SeekHorizonBlue:   notimplemented.SeekHorizonBlue{},

	ids.ShatterSorceryBlue: notimplemented.ShatterSorceryBlue{},

	ids.SiftRed:    notimplemented.SiftRed{},
	ids.SiftYellow: notimplemented.SiftYellow{},
	ids.SiftBlue:   notimplemented.SiftBlue{},

	ids.SigilOfCyclesBlue: notimplemented.SigilOfCyclesBlue{},

	ids.SigilOfFyendalBlue: cards.SigilOfFyendalBlue{},

	ids.SigilOfProtectionRed:    notimplemented.SigilOfProtectionRed{},
	ids.SigilOfProtectionYellow: notimplemented.SigilOfProtectionYellow{},
	ids.SigilOfProtectionBlue:   notimplemented.SigilOfProtectionBlue{},

	ids.SigilOfSolaceRed:    notimplemented.SigilOfSolaceRed{},
	ids.SigilOfSolaceYellow: notimplemented.SigilOfSolaceYellow{},
	ids.SigilOfSolaceBlue:   notimplemented.SigilOfSolaceBlue{},

	ids.SinkBelowRed:    notimplemented.SinkBelowRed{},
	ids.SinkBelowYellow: notimplemented.SinkBelowYellow{},
	ids.SinkBelowBlue:   notimplemented.SinkBelowBlue{},

	ids.SirensOfSafeHarborRed:    cards.SirensOfSafeHarborRed{},
	ids.SirensOfSafeHarborYellow: cards.SirensOfSafeHarborYellow{},
	ids.SirensOfSafeHarborBlue:   cards.SirensOfSafeHarborBlue{},

	ids.SloggismRed:    cards.SloggismRed{},
	ids.SloggismYellow: cards.SloggismYellow{},
	ids.SloggismBlue:   cards.SloggismBlue{},

	ids.SmashingGoodTimeRed:    notimplemented.SmashingGoodTimeRed{},
	ids.SmashingGoodTimeYellow: notimplemented.SmashingGoodTimeYellow{},
	ids.SmashingGoodTimeBlue:   notimplemented.SmashingGoodTimeBlue{},

	ids.SmashUpRed: notimplemented.SmashUpRed{},

	ids.SnatchRed:    cards.SnatchRed{},
	ids.SnatchYellow: cards.SnatchYellow{},
	ids.SnatchBlue:   cards.SnatchBlue{},

	ids.SoundTheAlarmRed: notimplemented.SoundTheAlarmRed{},

	ids.SpringboardSomersaultYellow: cards.SpringboardSomersaultYellow{},

	ids.SpringLoadRed:    cards.SpringLoadRed{},
	ids.SpringLoadYellow: cards.SpringLoadYellow{},
	ids.SpringLoadBlue:   cards.SpringLoadBlue{},

	ids.StartingStakeYellow: notimplemented.StartingStakeYellow{},

	ids.StonyWoottonhogRed:    notimplemented.StonyWoottonhogRed{},
	ids.StonyWoottonhogYellow: notimplemented.StonyWoottonhogYellow{},
	ids.StonyWoottonhogBlue:   notimplemented.StonyWoottonhogBlue{},

	ids.StrategicPlanningRed:    notimplemented.StrategicPlanningRed{},
	ids.StrategicPlanningYellow: notimplemented.StrategicPlanningYellow{},
	ids.StrategicPlanningBlue:   notimplemented.StrategicPlanningBlue{},

	ids.StrikeGoldRed:    notimplemented.StrikeGoldRed{},
	ids.StrikeGoldYellow: notimplemented.StrikeGoldYellow{},
	ids.StrikeGoldBlue:   notimplemented.StrikeGoldBlue{},

	ids.SunKissRed:    cards.SunKissRed{},
	ids.SunKissYellow: cards.SunKissYellow{},
	ids.SunKissBlue:   cards.SunKissBlue{},

	ids.SurgingMilitiaRed:    notimplemented.SurgingMilitiaRed{},
	ids.SurgingMilitiaYellow: notimplemented.SurgingMilitiaYellow{},
	ids.SurgingMilitiaBlue:   notimplemented.SurgingMilitiaBlue{},

	ids.TalismanOfBalanceBlue: notimplemented.TalismanOfBalanceBlue{},

	ids.TalismanOfCremationBlue: notimplemented.TalismanOfCremationBlue{},

	ids.TalismanOfDousingYellow: notimplemented.TalismanOfDousingYellow{},

	ids.TalismanOfFeatherfootYellow: notimplemented.TalismanOfFeatherfootYellow{},

	ids.TalismanOfRecompenseYellow: notimplemented.TalismanOfRecompenseYellow{},

	ids.TalismanOfTithesBlue: notimplemented.TalismanOfTithesBlue{},

	ids.TalismanOfWarfareYellow: notimplemented.TalismanOfWarfareYellow{},

	ids.TestOfStrengthRed: notimplemented.TestOfStrengthRed{},

	ids.ThrustRed: notimplemented.ThrustRed{},

	ids.TimesnapPotionBlue: notimplemented.TimesnapPotionBlue{},

	ids.TipOffRed:    notimplemented.TipOffRed{},
	ids.TipOffYellow: notimplemented.TipOffYellow{},
	ids.TipOffBlue:   notimplemented.TipOffBlue{},

	ids.TitaniumBaubleBlue: notimplemented.TitaniumBaubleBlue{},

	ids.TitForTatBlue: notimplemented.TitForTatBlue{},

	ids.TongueTiedRed: notimplemented.TongueTiedRed{},

	ids.ToughenUpBlue: cards.ToughenUpBlue{},

	ids.TradeInRed:    notimplemented.TradeInRed{},
	ids.TradeInYellow: notimplemented.TradeInYellow{},
	ids.TradeInBlue:   notimplemented.TradeInBlue{},

	ids.TremorOfIArathaelRed:    notimplemented.TremorOfIArathaelRed{},
	ids.TremorOfIArathaelYellow: notimplemented.TremorOfIArathaelYellow{},
	ids.TremorOfIArathaelBlue:   notimplemented.TremorOfIArathaelBlue{},

	ids.TrotAlongBlue: cards.TrotAlongBlue{},

	ids.UnmovableRed:    cards.UnmovableRed{},
	ids.UnmovableYellow: cards.UnmovableYellow{},
	ids.UnmovableBlue:   cards.UnmovableBlue{},

	ids.VigorRushRed:    cards.VigorRushRed{},
	ids.VigorRushYellow: cards.VigorRushYellow{},
	ids.VigorRushBlue:   cards.VigorRushBlue{},

	ids.VisitTheBlacksmithBlue: notimplemented.VisitTheBlacksmithBlue{},

	ids.WageGoldRed:    notimplemented.WageGoldRed{},
	ids.WageGoldYellow: notimplemented.WageGoldYellow{},
	ids.WageGoldBlue:   notimplemented.WageGoldBlue{},

	ids.WalkThePlankRed:    notimplemented.WalkThePlankRed{},
	ids.WalkThePlankYellow: notimplemented.WalkThePlankYellow{},
	ids.WalkThePlankBlue:   notimplemented.WalkThePlankBlue{},

	ids.WarmongersRecitalRed:    notimplemented.WarmongersRecitalRed{},
	ids.WarmongersRecitalYellow: notimplemented.WarmongersRecitalYellow{},
	ids.WarmongersRecitalBlue:   notimplemented.WarmongersRecitalBlue{},

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

	ids.WreckHavocRed:    notimplemented.WreckHavocRed{},
	ids.WreckHavocYellow: notimplemented.WreckHavocYellow{},
	ids.WreckHavocBlue:   notimplemented.WreckHavocBlue{},

	ids.YintiYantiRed:    notimplemented.YintiYantiRed{},
	ids.YintiYantiYellow: notimplemented.YintiYantiYellow{},
	ids.YintiYantiBlue:   notimplemented.YintiYantiBlue{},

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
// Freshly allocated; safe to mutate. Test-only synthetic IDs (testutils.FakeRedAttack, …)
// live in the testutils package and aren't registered here, so this is just AllCards()
// under a different name — kept distinct so callers who want "deck-legal cards" stay
// readable even if the registry ever holds non-deckable entries again.
func DeckableCards() []CardID {
	return AllCards()
}
