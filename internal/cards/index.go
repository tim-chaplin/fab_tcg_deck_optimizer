// Package cards is the master registry of every implemented card. The canonical ID type and
// constants live in package card; this package maps IDs to concrete Card values and provides
// lookup / iteration helpers useful for random deck generation, serialization, and compact
// equality checks.
//
// Weapons aren't ID-indexed — they're equipment, not deck cards. The weapon roster lives in
// package weapon alongside the Weapon implementations (weapon.All, weapon.ByName).
package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// ID aliases card.ID so callers of this package don't need two imports just to hold IDs.
type ID = card.ID

// Invalid aliases card.Invalid — the sentinel zero value.
const Invalid = card.Invalid

// byID is indexed directly by ID. Index 0 (Invalid) is nil.
var byID = []card.Card{
	card.Invalid: nil,

	card.AetherSlashRed:    AetherSlashRed{},
	card.AetherSlashYellow: AetherSlashYellow{},
	card.AetherSlashBlue:   AetherSlashBlue{},

	card.AmplifyTheArknightRed:    AmplifyTheArknightRed{},
	card.AmplifyTheArknightYellow: AmplifyTheArknightYellow{},
	card.AmplifyTheArknightBlue:   AmplifyTheArknightBlue{},

	card.ArcaneCussingRed:    ArcaneCussingRed{},
	card.ArcaneCussingYellow: ArcaneCussingYellow{},
	card.ArcaneCussingBlue:   ArcaneCussingBlue{},

	card.ArcanicCrackleRed:    ArcanicCrackleRed{},
	card.ArcanicCrackleYellow: ArcanicCrackleYellow{},
	card.ArcanicCrackleBlue:   ArcanicCrackleBlue{},

	card.ArcanicSpikeRed:    ArcanicSpikeRed{},
	card.ArcanicSpikeYellow: ArcanicSpikeYellow{},
	card.ArcanicSpikeBlue:   ArcanicSpikeBlue{},

	card.BlessingOfOccultRed:    BlessingOfOccultRed{},
	card.BlessingOfOccultYellow: BlessingOfOccultYellow{},
	card.BlessingOfOccultBlue:   BlessingOfOccultBlue{},

	card.BloodspillInvocationRed:    BloodspillInvocationRed{},
	card.BloodspillInvocationYellow: BloodspillInvocationYellow{},
	card.BloodspillInvocationBlue:   BloodspillInvocationBlue{},

	card.CondemnToSlaughterRed:    CondemnToSlaughterRed{},
	card.CondemnToSlaughterYellow: CondemnToSlaughterYellow{},
	card.CondemnToSlaughterBlue:   CondemnToSlaughterBlue{},

	card.ConsumingVolitionRed:    ConsumingVolitionRed{},
	card.ConsumingVolitionYellow: ConsumingVolitionYellow{},
	card.ConsumingVolitionBlue:   ConsumingVolitionBlue{},

	card.DeathlyDuetRed:    DeathlyDuetRed{},
	card.DeathlyDuetYellow: DeathlyDuetYellow{},
	card.DeathlyDuetBlue:   DeathlyDuetBlue{},

	card.DrawnToTheDarkDimensionRed:    DrawnToTheDarkDimensionRed{},
	card.DrawnToTheDarkDimensionYellow: DrawnToTheDarkDimensionYellow{},
	card.DrawnToTheDarkDimensionBlue:   DrawnToTheDarkDimensionBlue{},

	card.DrowningDireRed:    DrowningDireRed{},
	card.DrowningDireYellow: DrowningDireYellow{},
	card.DrowningDireBlue:   DrowningDireBlue{},

	card.HitTheHighNotesRed:    HitTheHighNotesRed{},
	card.HitTheHighNotesYellow: HitTheHighNotesYellow{},
	card.HitTheHighNotesBlue:   HitTheHighNotesBlue{},

	card.HocusPocusRed:    HocusPocusRed{},
	card.HocusPocusYellow: HocusPocusYellow{},
	card.HocusPocusBlue:   HocusPocusBlue{},

	card.MaleficIncantationRed:    MaleficIncantationRed{},
	card.MaleficIncantationYellow: MaleficIncantationYellow{},
	card.MaleficIncantationBlue:   MaleficIncantationBlue{},

	card.MauvrionSkiesRed:    MauvrionSkiesRed{},
	card.MauvrionSkiesYellow: MauvrionSkiesYellow{},
	card.MauvrionSkiesBlue:   MauvrionSkiesBlue{},

	card.MeatAndGreetRed:    MeatAndGreetRed{},
	card.MeatAndGreetYellow: MeatAndGreetYellow{},
	card.MeatAndGreetBlue:   MeatAndGreetBlue{},

	card.OathOfTheArknightRed:    OathOfTheArknightRed{},
	card.OathOfTheArknightYellow: OathOfTheArknightYellow{},
	card.OathOfTheArknightBlue:   OathOfTheArknightBlue{},

	card.ReadTheRunesRed:    ReadTheRunesRed{},
	card.ReadTheRunesYellow: ReadTheRunesYellow{},
	card.ReadTheRunesBlue:   ReadTheRunesBlue{},

	card.ReduceToRunechantRed:    ReduceToRunechantRed{},
	card.ReduceToRunechantYellow: ReduceToRunechantYellow{},
	card.ReduceToRunechantBlue:   ReduceToRunechantBlue{},

	card.ReekOfCorruptionRed:    ReekOfCorruptionRed{},
	card.ReekOfCorruptionYellow: ReekOfCorruptionYellow{},
	card.ReekOfCorruptionBlue:   ReekOfCorruptionBlue{},

	card.RuneFlashRed:    RuneFlashRed{},
	card.RuneFlashYellow: RuneFlashYellow{},
	card.RuneFlashBlue:   RuneFlashBlue{},

	card.RunebloodIncantationRed:    RunebloodIncantationRed{},
	card.RunebloodIncantationYellow: RunebloodIncantationYellow{},
	card.RunebloodIncantationBlue:   RunebloodIncantationBlue{},

	card.RuneragerSwarmRed:    RuneragerSwarmRed{},
	card.RuneragerSwarmYellow: RuneragerSwarmYellow{},
	card.RuneragerSwarmBlue:   RuneragerSwarmBlue{},

	card.RunicFellingsongRed:    RunicFellingsongRed{},
	card.RunicFellingsongYellow: RunicFellingsongYellow{},
	card.RunicFellingsongBlue:   RunicFellingsongBlue{},

	card.RunicReapingRed:    RunicReapingRed{},
	card.RunicReapingYellow: RunicReapingYellow{},
	card.RunicReapingBlue:   RunicReapingBlue{},

	card.ShrillOfSkullformRed:    ShrillOfSkullformRed{},
	card.ShrillOfSkullformYellow: ShrillOfSkullformYellow{},
	card.ShrillOfSkullformBlue:   ShrillOfSkullformBlue{},

	card.SigilOfDeadwoodBlue: SigilOfDeadwoodBlue{},

	card.SigilOfSilphidaeBlue: SigilOfSilphidaeBlue{},

	card.SigilOfSufferingRed:    SigilOfSufferingRed{},
	card.SigilOfSufferingYellow: SigilOfSufferingYellow{},
	card.SigilOfSufferingBlue:   SigilOfSufferingBlue{},

	card.SigilOfTheArknightBlue: SigilOfTheArknightBlue{},

	card.SingeingSteelbladeRed:    SingeingSteelbladeRed{},
	card.SingeingSteelbladeYellow: SingeingSteelbladeYellow{},
	card.SingeingSteelbladeBlue:   SingeingSteelbladeBlue{},

	card.SkyFireLanternsRed:    SkyFireLanternsRed{},
	card.SkyFireLanternsYellow: SkyFireLanternsYellow{},
	card.SkyFireLanternsBlue:   SkyFireLanternsBlue{},

	card.SpellbladeAssaultRed:    SpellbladeAssaultRed{},
	card.SpellbladeAssaultYellow: SpellbladeAssaultYellow{},
	card.SpellbladeAssaultBlue:   SpellbladeAssaultBlue{},

	card.SpellbladeStrikeRed:    SpellbladeStrikeRed{},
	card.SpellbladeStrikeYellow: SpellbladeStrikeYellow{},
	card.SpellbladeStrikeBlue:   SpellbladeStrikeBlue{},

	card.SplinteringDeadwoodRed:    SplinteringDeadwoodRed{},
	card.SplinteringDeadwoodYellow: SplinteringDeadwoodYellow{},
	card.SplinteringDeadwoodBlue:   SplinteringDeadwoodBlue{},

	card.SutcliffesResearchNotesRed:    SutcliffesResearchNotesRed{},
	card.SutcliffesResearchNotesYellow: SutcliffesResearchNotesYellow{},
	card.SutcliffesResearchNotesBlue:   SutcliffesResearchNotesBlue{},

	card.VantagePointRed:    VantagePointRed{},
	card.VantagePointYellow: VantagePointYellow{},
	card.VantagePointBlue:   VantagePointBlue{},

	card.VexingMaliceRed:    VexingMaliceRed{},
	card.VexingMaliceYellow: VexingMaliceYellow{},
	card.VexingMaliceBlue:   VexingMaliceBlue{},

	card.WeepingBattlegroundRed:    WeepingBattlegroundRed{},
	card.WeepingBattlegroundYellow: WeepingBattlegroundYellow{},
	card.WeepingBattlegroundBlue:   WeepingBattlegroundBlue{},

	card.AdrenalineRushRed:    AdrenalineRushRed{},
	card.AdrenalineRushYellow: AdrenalineRushYellow{},
	card.AdrenalineRushBlue:   AdrenalineRushBlue{},

	card.AmuletOfAssertivenessYellow: AmuletOfAssertivenessYellow{},

	card.AmuletOfEchoesBlue: AmuletOfEchoesBlue{},

	card.AmuletOfHavencallBlue: AmuletOfHavencallBlue{},

	card.AmuletOfIgnitionYellow: AmuletOfIgnitionYellow{},

	card.AmuletOfInterventionBlue: AmuletOfInterventionBlue{},

	card.AmuletOfOblationBlue: AmuletOfOblationBlue{},

	card.ArcanePolarityRed:    ArcanePolarityRed{},
	card.ArcanePolarityYellow: ArcanePolarityYellow{},
	card.ArcanePolarityBlue:   ArcanePolarityBlue{},

	card.BackAlleyBreaklineRed:    BackAlleyBreaklineRed{},
	card.BackAlleyBreaklineYellow: BackAlleyBreaklineYellow{},
	card.BackAlleyBreaklineBlue:   BackAlleyBreaklineBlue{},

	card.BarragingBrawnhideRed:    BarragingBrawnhideRed{},
	card.BarragingBrawnhideYellow: BarragingBrawnhideYellow{},
	card.BarragingBrawnhideBlue:   BarragingBrawnhideBlue{},

	card.BattlefrontBastionRed:    BattlefrontBastionRed{},
	card.BattlefrontBastionYellow: BattlefrontBastionYellow{},
	card.BattlefrontBastionBlue:   BattlefrontBastionBlue{},

	card.BelittleRed:    BelittleRed{},
	card.BelittleYellow: BelittleYellow{},
	card.BelittleBlue:   BelittleBlue{},

	card.BladeFlashBlue: BladeFlashBlue{},

	card.BlanchRed:    BlanchRed{},
	card.BlanchYellow: BlanchYellow{},
	card.BlanchBlue:   BlanchBlue{},

	card.BlowForABlowRed: BlowForABlowRed{},

	card.BlusterBuffRed: BlusterBuffRed{},

	card.BrandishRed:    BrandishRed{},
	card.BrandishYellow: BrandishYellow{},
	card.BrandishBlue:   BrandishBlue{},

	card.BrothersInArmsRed:    BrothersInArmsRed{},
	card.BrothersInArmsYellow: BrothersInArmsYellow{},
	card.BrothersInArmsBlue:   BrothersInArmsBlue{},

	card.BrushOffRed:    BrushOffRed{},
	card.BrushOffYellow: BrushOffYellow{},
	card.BrushOffBlue:   BrushOffBlue{},

	card.BrutalAssaultRed:    BrutalAssaultRed{},
	card.BrutalAssaultYellow: BrutalAssaultYellow{},
	card.BrutalAssaultBlue:   BrutalAssaultBlue{},

	card.CadaverousContrabandRed:    CadaverousContrabandRed{},
	card.CadaverousContrabandYellow: CadaverousContrabandYellow{},
	card.CadaverousContrabandBlue:   CadaverousContrabandBlue{},

	card.CalmingBreezeRed: CalmingBreezeRed{},

	card.CaptainsCallRed:    CaptainsCallRed{},
	card.CaptainsCallYellow: CaptainsCallYellow{},
	card.CaptainsCallBlue:   CaptainsCallBlue{},

	card.CashInYellow: CashInYellow{},

	card.ChestPuffRed: ChestPuffRed{},

	card.ClapEmInIronsBlue: ClapEmInIronsBlue{},

	card.ClarityPotionBlue: ClarityPotionBlue{},

	card.ClearwaterElixirRed: ClearwaterElixirRed{},

	card.ComeToFightRed:    ComeToFightRed{},
	card.ComeToFightYellow: ComeToFightYellow{},
	card.ComeToFightBlue:   ComeToFightBlue{},

	card.CountYourBlessingsRed:    CountYourBlessingsRed{},
	card.CountYourBlessingsYellow: CountYourBlessingsYellow{},
	card.CountYourBlessingsBlue:   CountYourBlessingsBlue{},

	card.CrackedBaubleYellow: CrackedBaubleYellow{},

	card.CrashDownTheGatesRed:    CrashDownTheGatesRed{},
	card.CrashDownTheGatesYellow: CrashDownTheGatesYellow{},
	card.CrashDownTheGatesBlue:   CrashDownTheGatesBlue{},

	card.CriticalStrikeRed:    CriticalStrikeRed{},
	card.CriticalStrikeYellow: CriticalStrikeYellow{},
	card.CriticalStrikeBlue:   CriticalStrikeBlue{},

	card.CutDownToSizeRed:    CutDownToSizeRed{},
	card.CutDownToSizeYellow: CutDownToSizeYellow{},
	card.CutDownToSizeBlue:   CutDownToSizeBlue{},

	card.DemolitionCrewRed:    DemolitionCrewRed{},
	card.DemolitionCrewYellow: DemolitionCrewYellow{},
	card.DemolitionCrewBlue:   DemolitionCrewBlue{},

	card.DestructiveDeliberationRed:    DestructiveDeliberationRed{},
	card.DestructiveDeliberationYellow: DestructiveDeliberationYellow{},
	card.DestructiveDeliberationBlue:   DestructiveDeliberationBlue{},

	card.DestructiveTendenciesBlue: DestructiveTendenciesBlue{},

	card.DodgeBlue: DodgeBlue{},

	card.DownButNotOutRed:    DownButNotOutRed{},
	card.DownButNotOutYellow: DownButNotOutYellow{},
	card.DownButNotOutBlue:   DownButNotOutBlue{},

	card.DragDownRed:    DragDownRed{},
	card.DragDownYellow: DragDownYellow{},
	card.DragDownBlue:   DragDownBlue{},

	card.DroneOfBrutalityRed:    DroneOfBrutalityRed{},
	card.DroneOfBrutalityYellow: DroneOfBrutalityYellow{},
	card.DroneOfBrutalityBlue:   DroneOfBrutalityBlue{},

	card.EirinasPrayerRed:    EirinasPrayerRed{},
	card.EirinasPrayerYellow: EirinasPrayerYellow{},
	card.EirinasPrayerBlue:   EirinasPrayerBlue{},

	card.EmissaryOfMoonRed: EmissaryOfMoonRed{},

	card.EmissaryOfTidesRed: EmissaryOfTidesRed{},

	card.EmissaryOfWindRed: EmissaryOfWindRed{},

	card.EnchantingMelodyRed:    EnchantingMelodyRed{},
	card.EnchantingMelodyYellow: EnchantingMelodyYellow{},
	card.EnchantingMelodyBlue:   EnchantingMelodyBlue{},

	card.EnergyPotionBlue: EnergyPotionBlue{},

	card.EvasiveLeapRed:    EvasiveLeapRed{},
	card.EvasiveLeapYellow: EvasiveLeapYellow{},
	card.EvasiveLeapBlue:   EvasiveLeapBlue{},

	card.EvenBiggerThanThatRed:    EvenBiggerThanThatRed{},
	card.EvenBiggerThanThatYellow: EvenBiggerThanThatYellow{},
	card.EvenBiggerThanThatBlue:   EvenBiggerThanThatBlue{},

	card.ExposedBlue: ExposedBlue{},

	card.FactFindingMissionRed:    FactFindingMissionRed{},
	card.FactFindingMissionYellow: FactFindingMissionYellow{},
	card.FactFindingMissionBlue:   FactFindingMissionBlue{},

	card.FateForeseenRed:    FateForeseenRed{},
	card.FateForeseenYellow: FateForeseenYellow{},
	card.FateForeseenBlue:   FateForeseenBlue{},

	card.FeistyLocalsRed:    FeistyLocalsRed{},
	card.FeistyLocalsYellow: FeistyLocalsYellow{},
	card.FeistyLocalsBlue:   FeistyLocalsBlue{},

	card.FerventForerunnerRed:    FerventForerunnerRed{},
	card.FerventForerunnerYellow: FerventForerunnerYellow{},
	card.FerventForerunnerBlue:   FerventForerunnerBlue{},

	card.FiddlersGreenRed:    FiddlersGreenRed{},
	card.FiddlersGreenYellow: FiddlersGreenYellow{},
	card.FiddlersGreenBlue:   FiddlersGreenBlue{},

	card.FlexRed:    FlexRed{},
	card.FlexYellow: FlexYellow{},
	card.FlexBlue:   FlexBlue{},

	card.FlockOfTheFeatherWalkersRed:    FlockOfTheFeatherWalkersRed{},
	card.FlockOfTheFeatherWalkersYellow: FlockOfTheFeatherWalkersYellow{},
	card.FlockOfTheFeatherWalkersBlue:   FlockOfTheFeatherWalkersBlue{},

	card.FlyingHighRed:    FlyingHighRed{},
	card.FlyingHighYellow: FlyingHighYellow{},
	card.FlyingHighBlue:   FlyingHighBlue{},

	card.FoolsGoldYellow: FoolsGoldYellow{},

	card.ForceSightRed:    ForceSightRed{},
	card.ForceSightYellow: ForceSightYellow{},
	card.ForceSightBlue:   ForceSightBlue{},

	card.FreewheelingRenegadesRed:    FreewheelingRenegadesRed{},
	card.FreewheelingRenegadesYellow: FreewheelingRenegadesYellow{},
	card.FreewheelingRenegadesBlue:   FreewheelingRenegadesBlue{},

	card.FrontlineScoutRed:    FrontlineScoutRed{},
	card.FrontlineScoutYellow: FrontlineScoutYellow{},
	card.FrontlineScoutBlue:   FrontlineScoutBlue{},

	card.FyendalsFightingSpiritRed:    FyendalsFightingSpiritRed{},
	card.FyendalsFightingSpiritYellow: FyendalsFightingSpiritYellow{},
	card.FyendalsFightingSpiritBlue:   FyendalsFightingSpiritBlue{},

	card.GravekeepingRed:    GravekeepingRed{},
	card.GravekeepingYellow: GravekeepingYellow{},
	card.GravekeepingBlue:   GravekeepingBlue{},

	card.HandBehindThePenRed: HandBehindThePenRed{},

	card.HealingBalmRed:    HealingBalmRed{},
	card.HealingBalmYellow: HealingBalmYellow{},
	card.HealingBalmBlue:   HealingBalmBlue{},

	card.HealingPotionBlue: HealingPotionBlue{},

	card.HighStrikerRed:    HighStrikerRed{},
	card.HighStrikerYellow: HighStrikerYellow{},
	card.HighStrikerBlue:   HighStrikerBlue{},

	card.HumbleRed:    HumbleRed{},
	card.HumbleYellow: HumbleYellow{},
	card.HumbleBlue:   HumbleBlue{},

	card.ImperialSealOfCommandRed: ImperialSealOfCommandRed{},

	card.InfectiousHostRed:    InfectiousHostRed{},
	card.InfectiousHostYellow: InfectiousHostYellow{},
	card.InfectiousHostBlue:   InfectiousHostBlue{},

	card.JackBeNimbleRed: JackBeNimbleRed{},

	card.JackBeQuickRed: JackBeQuickRed{},

	card.LayLowYellow: LayLowYellow{},

	card.LeadTheChargeRed:    LeadTheChargeRed{},
	card.LeadTheChargeYellow: LeadTheChargeYellow{},
	card.LeadTheChargeBlue:   LeadTheChargeBlue{},

	card.LifeForALifeRed:    LifeForALifeRed{},
	card.LifeForALifeYellow: LifeForALifeYellow{},
	card.LifeForALifeBlue:   LifeForALifeBlue{},

	card.LifeOfThePartyRed:    LifeOfThePartyRed{},
	card.LifeOfThePartyYellow: LifeOfThePartyYellow{},
	card.LifeOfThePartyBlue:   LifeOfThePartyBlue{},

	card.LookingForAScrapRed:    LookingForAScrapRed{},
	card.LookingForAScrapYellow: LookingForAScrapYellow{},
	card.LookingForAScrapBlue:   LookingForAScrapBlue{},

	card.LookTuffRed: LookTuffRed{},

	card.LungingPressBlue: LungingPressBlue{},

	card.MemorialGroundRed:    MemorialGroundRed{},
	card.MemorialGroundYellow: MemorialGroundYellow{},
	card.MemorialGroundBlue:   MemorialGroundBlue{},

	card.MinnowismRed:    MinnowismRed{},
	card.MinnowismYellow: MinnowismYellow{},
	card.MinnowismBlue:   MinnowismBlue{},

	card.MoneyOrYourLifeRed:    MoneyOrYourLifeRed{},
	card.MoneyOrYourLifeYellow: MoneyOrYourLifeYellow{},
	card.MoneyOrYourLifeBlue:   MoneyOrYourLifeBlue{},

	card.MoneyWhereYaMouthIsRed:    MoneyWhereYaMouthIsRed{},
	card.MoneyWhereYaMouthIsYellow: MoneyWhereYaMouthIsYellow{},
	card.MoneyWhereYaMouthIsBlue:   MoneyWhereYaMouthIsBlue{},

	card.MoonWishRed:    MoonWishRed{},
	card.MoonWishYellow: MoonWishYellow{},
	card.MoonWishBlue:   MoonWishBlue{},

	card.MuscleMuttYellow: MuscleMuttYellow{},

	card.NimbleStrikeRed:    NimbleStrikeRed{},
	card.NimbleStrikeYellow: NimbleStrikeYellow{},
	card.NimbleStrikeBlue:   NimbleStrikeBlue{},

	card.NimblismRed:    NimblismRed{},
	card.NimblismYellow: NimblismYellow{},
	card.NimblismBlue:   NimblismBlue{},

	card.NimbyRed:    NimbyRed{},
	card.NimbyYellow: NimbyYellow{},
	card.NimbyBlue:   NimbyBlue{},

	card.NipAtTheHeelsBlue: NipAtTheHeelsBlue{},

	card.OasisRespiteRed:    OasisRespiteRed{},
	card.OasisRespiteYellow: OasisRespiteYellow{},
	card.OasisRespiteBlue:   OasisRespiteBlue{},

	card.OnAKnifeEdgeYellow: OnAKnifeEdgeYellow{},

	card.OnTheHorizonRed:    OnTheHorizonRed{},
	card.OnTheHorizonYellow: OnTheHorizonYellow{},
	card.OnTheHorizonBlue:   OnTheHorizonBlue{},

	card.OutedRed: OutedRed{},

	card.OutMuscleRed:    OutMuscleRed{},
	card.OutMuscleYellow: OutMuscleYellow{},
	card.OutMuscleBlue:   OutMuscleBlue{},

	card.OverloadRed:    OverloadRed{},
	card.OverloadYellow: OverloadYellow{},
	card.OverloadBlue:   OverloadBlue{},

	card.PeaceOfMindRed:    PeaceOfMindRed{},
	card.PeaceOfMindYellow: PeaceOfMindYellow{},
	card.PeaceOfMindBlue:   PeaceOfMindBlue{},

	card.PerformanceBonusRed:    PerformanceBonusRed{},
	card.PerformanceBonusYellow: PerformanceBonusYellow{},
	card.PerformanceBonusBlue:   PerformanceBonusBlue{},

	card.PickACardAnyCardRed:    PickACardAnyCardRed{},
	card.PickACardAnyCardYellow: PickACardAnyCardYellow{},
	card.PickACardAnyCardBlue:   PickACardAnyCardBlue{},

	card.PilferTheTombBlue: PilferTheTombBlue{},

	card.PlunderRunRed:    PlunderRunRed{},
	card.PlunderRunYellow: PlunderRunYellow{},
	card.PlunderRunBlue:   PlunderRunBlue{},

	card.PotionOfDejaVuBlue: PotionOfDejaVuBlue{},

	card.PotionOfIronhideBlue: PotionOfIronhideBlue{},

	card.PotionOfLuckBlue: PotionOfLuckBlue{},

	card.PotionOfSeeingBlue: PotionOfSeeingBlue{},

	card.PotionOfStrengthBlue: PotionOfStrengthBlue{},

	card.PoundForPoundRed:    PoundForPoundRed{},
	card.PoundForPoundYellow: PoundForPoundYellow{},
	card.PoundForPoundBlue:   PoundForPoundBlue{},

	card.PrimeTheCrowdRed:    PrimeTheCrowdRed{},
	card.PrimeTheCrowdYellow: PrimeTheCrowdYellow{},
	card.PrimeTheCrowdBlue:   PrimeTheCrowdBlue{},

	card.PromiseOfPlentyRed:    PromiseOfPlentyRed{},
	card.PromiseOfPlentyYellow: PromiseOfPlentyYellow{},
	card.PromiseOfPlentyBlue:   PromiseOfPlentyBlue{},

	card.PublicBountyRed:    PublicBountyRed{},
	card.PublicBountyYellow: PublicBountyYellow{},
	card.PublicBountyBlue:   PublicBountyBlue{},

	card.PummelRed:    PummelRed{},
	card.PummelYellow: PummelYellow{},
	card.PummelBlue:   PummelBlue{},

	card.PunchAboveYourWeightRed:    PunchAboveYourWeightRed{},
	card.PunchAboveYourWeightYellow: PunchAboveYourWeightYellow{},
	card.PunchAboveYourWeightBlue:   PunchAboveYourWeightBlue{},

	card.PursueToTheEdgeOfOblivionRed: PursueToTheEdgeOfOblivionRed{},

	card.PursueToThePitsOfDespairRed: PursueToThePitsOfDespairRed{},

	card.PushThePointRed:    PushThePointRed{},
	card.PushThePointYellow: PushThePointYellow{},
	card.PushThePointBlue:   PushThePointBlue{},

	card.PutInContextBlue: PutInContextBlue{},

	card.RagingOnslaughtRed:    RagingOnslaughtRed{},
	card.RagingOnslaughtYellow: RagingOnslaughtYellow{},
	card.RagingOnslaughtBlue:   RagingOnslaughtBlue{},

	card.RallyTheCoastGuardRed:    RallyTheCoastGuardRed{},
	card.RallyTheCoastGuardYellow: RallyTheCoastGuardYellow{},
	card.RallyTheCoastGuardBlue:   RallyTheCoastGuardBlue{},

	card.RallyTheRearguardRed:    RallyTheRearguardRed{},
	card.RallyTheRearguardYellow: RallyTheRearguardYellow{},
	card.RallyTheRearguardBlue:   RallyTheRearguardBlue{},

	card.RansackAndRazeBlue: RansackAndRazeBlue{},

	card.RavenousRabbleRed:    RavenousRabbleRed{},
	card.RavenousRabbleYellow: RavenousRabbleYellow{},
	card.RavenousRabbleBlue:   RavenousRabbleBlue{},

	card.RazorReflexRed:    RazorReflexRed{},
	card.RazorReflexYellow: RazorReflexYellow{},
	card.RazorReflexBlue:   RazorReflexBlue{},

	card.RegainComposureBlue: RegainComposureBlue{},

	card.RegurgitatingSlogRed:    RegurgitatingSlogRed{},
	card.RegurgitatingSlogYellow: RegurgitatingSlogYellow{},
	card.RegurgitatingSlogBlue:   RegurgitatingSlogBlue{},

	card.ReinforceTheLineRed:    ReinforceTheLineRed{},
	card.ReinforceTheLineYellow: ReinforceTheLineYellow{},
	card.ReinforceTheLineBlue:   ReinforceTheLineBlue{},

	card.RelentlessPursuitBlue: RelentlessPursuitBlue{},

	card.RestvineElixirRed: RestvineElixirRed{},

	card.RiftingRed:    RiftingRed{},
	card.RiftingYellow: RiftingYellow{},
	card.RiftingBlue:   RiftingBlue{},

	card.RightBehindYouRed:    RightBehindYouRed{},
	card.RightBehindYouYellow: RightBehindYouYellow{},
	card.RightBehindYouBlue:   RightBehindYouBlue{},

	card.RiseAboveRed:    RiseAboveRed{},
	card.RiseAboveYellow: RiseAboveYellow{},
	card.RiseAboveBlue:   RiseAboveBlue{},

	card.SapwoodElixirRed: SapwoodElixirRed{},

	card.ScarForAScarRed:    ScarForAScarRed{},
	card.ScarForAScarYellow: ScarForAScarYellow{},
	card.ScarForAScarBlue:   ScarForAScarBlue{},

	card.ScourTheBattlescapeRed:    ScourTheBattlescapeRed{},
	card.ScourTheBattlescapeYellow: ScourTheBattlescapeYellow{},
	card.ScourTheBattlescapeBlue:   ScourTheBattlescapeBlue{},

	card.ScoutThePeripheryRed:    ScoutThePeripheryRed{},
	card.ScoutThePeripheryYellow: ScoutThePeripheryYellow{},
	card.ScoutThePeripheryBlue:   ScoutThePeripheryBlue{},

	card.SeekHorizonRed:    SeekHorizonRed{},
	card.SeekHorizonYellow: SeekHorizonYellow{},
	card.SeekHorizonBlue:   SeekHorizonBlue{},

	card.ShatterSorceryBlue: ShatterSorceryBlue{},

	card.SiftRed:    SiftRed{},
	card.SiftYellow: SiftYellow{},
	card.SiftBlue:   SiftBlue{},

	card.SigilOfCyclesBlue: SigilOfCyclesBlue{},

	card.SigilOfFyendalBlue: SigilOfFyendalBlue{},

	card.SigilOfProtectionRed:    SigilOfProtectionRed{},
	card.SigilOfProtectionYellow: SigilOfProtectionYellow{},
	card.SigilOfProtectionBlue:   SigilOfProtectionBlue{},

	card.SigilOfSolaceRed:    SigilOfSolaceRed{},
	card.SigilOfSolaceYellow: SigilOfSolaceYellow{},
	card.SigilOfSolaceBlue:   SigilOfSolaceBlue{},

	card.SinkBelowRed:    SinkBelowRed{},
	card.SinkBelowYellow: SinkBelowYellow{},
	card.SinkBelowBlue:   SinkBelowBlue{},

	card.SirensOfSafeHarborRed:    SirensOfSafeHarborRed{},
	card.SirensOfSafeHarborYellow: SirensOfSafeHarborYellow{},
	card.SirensOfSafeHarborBlue:   SirensOfSafeHarborBlue{},

	card.SloggismRed:    SloggismRed{},
	card.SloggismYellow: SloggismYellow{},
	card.SloggismBlue:   SloggismBlue{},

	card.SmashingGoodTimeRed:    SmashingGoodTimeRed{},
	card.SmashingGoodTimeYellow: SmashingGoodTimeYellow{},
	card.SmashingGoodTimeBlue:   SmashingGoodTimeBlue{},

	card.SmashUpRed: SmashUpRed{},

	card.SnatchRed:    SnatchRed{},
	card.SnatchYellow: SnatchYellow{},
	card.SnatchBlue:   SnatchBlue{},

	card.SoundTheAlarmRed: SoundTheAlarmRed{},

	card.SpringboardSomersaultYellow: SpringboardSomersaultYellow{},

	card.SpringLoadRed:    SpringLoadRed{},
	card.SpringLoadYellow: SpringLoadYellow{},
	card.SpringLoadBlue:   SpringLoadBlue{},

	card.StartingStakeYellow: StartingStakeYellow{},

	card.StonyWoottonhogRed:    StonyWoottonhogRed{},
	card.StonyWoottonhogYellow: StonyWoottonhogYellow{},
	card.StonyWoottonhogBlue:   StonyWoottonhogBlue{},

	card.StrategicPlanningRed:    StrategicPlanningRed{},
	card.StrategicPlanningYellow: StrategicPlanningYellow{},
	card.StrategicPlanningBlue:   StrategicPlanningBlue{},

	card.StrikeGoldRed:    StrikeGoldRed{},
	card.StrikeGoldYellow: StrikeGoldYellow{},
	card.StrikeGoldBlue:   StrikeGoldBlue{},

	card.SunKissRed:    SunKissRed{},
	card.SunKissYellow: SunKissYellow{},
	card.SunKissBlue:   SunKissBlue{},

	card.SurgingMilitiaRed:    SurgingMilitiaRed{},
	card.SurgingMilitiaYellow: SurgingMilitiaYellow{},
	card.SurgingMilitiaBlue:   SurgingMilitiaBlue{},

	card.TalismanOfBalanceBlue: TalismanOfBalanceBlue{},

	card.TalismanOfCremationBlue: TalismanOfCremationBlue{},

	card.TalismanOfDousingYellow: TalismanOfDousingYellow{},

	card.TalismanOfFeatherfootYellow: TalismanOfFeatherfootYellow{},

	card.TalismanOfRecompenseYellow: TalismanOfRecompenseYellow{},

	card.TalismanOfTithesBlue: TalismanOfTithesBlue{},

	card.TalismanOfWarfareYellow: TalismanOfWarfareYellow{},

	card.TestOfStrengthRed: TestOfStrengthRed{},

	card.ThrustRed: ThrustRed{},

	card.TimesnapPotionBlue: TimesnapPotionBlue{},

	card.TipOffRed:    TipOffRed{},
	card.TipOffYellow: TipOffYellow{},
	card.TipOffBlue:   TipOffBlue{},

	card.TitaniumBaubleBlue: TitaniumBaubleBlue{},

	card.TitForTatBlue: TitForTatBlue{},

	card.TongueTiedRed: TongueTiedRed{},

	card.ToughenUpBlue: ToughenUpBlue{},

	card.TradeInRed:    TradeInRed{},
	card.TradeInYellow: TradeInYellow{},
	card.TradeInBlue:   TradeInBlue{},

	card.TremorOfIArathaelRed:    TremorOfIArathaelRed{},
	card.TremorOfIArathaelYellow: TremorOfIArathaelYellow{},
	card.TremorOfIArathaelBlue:   TremorOfIArathaelBlue{},

	card.TrotAlongBlue: TrotAlongBlue{},

	card.UnmovableRed:    UnmovableRed{},
	card.UnmovableYellow: UnmovableYellow{},
	card.UnmovableBlue:   UnmovableBlue{},

	card.VigorRushRed:    VigorRushRed{},
	card.VigorRushYellow: VigorRushYellow{},
	card.VigorRushBlue:   VigorRushBlue{},

	card.VisitTheBlacksmithBlue: VisitTheBlacksmithBlue{},

	card.WageGoldRed:    WageGoldRed{},
	card.WageGoldYellow: WageGoldYellow{},
	card.WageGoldBlue:   WageGoldBlue{},

	card.WalkThePlankRed:    WalkThePlankRed{},
	card.WalkThePlankYellow: WalkThePlankYellow{},
	card.WalkThePlankBlue:   WalkThePlankBlue{},

	card.WarmongersRecitalRed:    WarmongersRecitalRed{},
	card.WarmongersRecitalYellow: WarmongersRecitalYellow{},
	card.WarmongersRecitalBlue:   WarmongersRecitalBlue{},

	card.WaterTheSeedsRed:    WaterTheSeedsRed{},
	card.WaterTheSeedsYellow: WaterTheSeedsYellow{},
	card.WaterTheSeedsBlue:   WaterTheSeedsBlue{},

	card.WhisperOfTheOracleRed:    WhisperOfTheOracleRed{},
	card.WhisperOfTheOracleYellow: WhisperOfTheOracleYellow{},
	card.WhisperOfTheOracleBlue:   WhisperOfTheOracleBlue{},

	card.WoundedBullRed:    WoundedBullRed{},
	card.WoundedBullYellow: WoundedBullYellow{},
	card.WoundedBullBlue:   WoundedBullBlue{},

	card.WoundingBlowRed:    WoundingBlowRed{},
	card.WoundingBlowYellow: WoundingBlowYellow{},
	card.WoundingBlowBlue:   WoundingBlowBlue{},

	card.WreckHavocRed:    WreckHavocRed{},
	card.WreckHavocYellow: WreckHavocYellow{},
	card.WreckHavocBlue:   WreckHavocBlue{},

	card.YintiYantiRed:    YintiYantiRed{},
	card.YintiYantiYellow: YintiYantiYellow{},
	card.YintiYantiBlue:   YintiYantiBlue{},

	card.ZealousBeltingRed:    ZealousBeltingRed{},
	card.ZealousBeltingYellow: ZealousBeltingYellow{},
	card.ZealousBeltingBlue:   ZealousBeltingBlue{},
}

// init eagerly populates package card's chain-step text and DisplayName caches so the
// per-Play hot path is pure cache reads. Done at registration time because the registry
// is the only place that knows the full card set, and the caches are sized for the full
// ID space.
func init() {
	card.WarmChainStepCache(byID)
	card.WarmDisplayNameCache(byID)
}

// byName maps card.DisplayName(c) → ID for reverse lookup. Built once at init. Keyed on
// DisplayName (not bare Name) so each pitch variant gets a distinct entry — Card.Name()
// collapses all three printings to the same base string, so it's not a unique key.
var byName = func() map[string]ID {
	m := make(map[string]ID, len(byID)-1)
	for id, c := range byID {
		if c == nil {
			continue
		}
		m[card.DisplayName(c)] = ID(id)
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

// ByName looks up an ID by the card's DisplayName ("Aether Slash [R]"). Returns
// (Invalid, false) if no such card.
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

// Deckable returns every registered card ID that's legal to put in a real deck. Freshly
// allocated; safe to mutate. The fake card IDs (card.FakeRedAttack, …) are deliberately not
// in the registry, so this is just All() under a different name — kept distinct so callers
// who want "deck-legal cards" stay readable even if the registry ever holds non-deckable
// entries again.
func Deckable() []ID {
	return All()
}
