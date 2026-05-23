package heroes

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	_ "github.com/tim-chaplin/fab-deck-optimizer/internal/sim" // register sim's gameengine builders
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

// fakeRuneAttack is a minimal Runeblade attack-action card.
type fakeRuneAttack struct{}

func (fakeRuneAttack) ID() ids.CardID           { return ids.InvalidCard }
func (fakeRuneAttack) Name() string             { return "StubRuneAttack" }
func (fakeRuneAttack) DisplayName() string      { return "StubRuneAttack" }
func (fakeRuneAttack) Cost(card.GameEngine) int { return 0 }
func (fakeRuneAttack) Pitch() int               { return 0 }
func (fakeRuneAttack) Attack() int              { return 0 }
func (fakeRuneAttack) Defense() int             { return 0 }
func (fakeRuneAttack) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (fakeRuneAttack) GoAgain(card.GameEngine) bool                       { return true }
func (fakeRuneAttack) Play(card.GameEngine, card.Logger, *card.CardState) {}

// fakeRuneAura is a minimal Runeblade non-attack action (an Aura).
type fakeRuneAura struct{}

func (fakeRuneAura) ID() ids.CardID           { return ids.InvalidCard }
func (fakeRuneAura) Name() string             { return "StubRuneAura" }
func (fakeRuneAura) DisplayName() string      { return "StubRuneAura" }
func (fakeRuneAura) Cost(card.GameEngine) int { return 0 }
func (fakeRuneAura) Pitch() int               { return 0 }
func (fakeRuneAura) Attack() int              { return 0 }
func (fakeRuneAura) Defense() int             { return 0 }
func (fakeRuneAura) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)
}
func (fakeRuneAura) GoAgain(card.GameEngine) bool                       { return true }
func (fakeRuneAura) Play(card.GameEngine, card.Logger, *card.CardState) {}

// fakeNonRuneblade is an Action-Attack with no Runeblade type — should never trigger Viserai.
type fakeNonRuneblade struct{}

func (fakeNonRuneblade) ID() ids.CardID           { return ids.InvalidCard }
func (fakeNonRuneblade) Name() string             { return "StubGeneric" }
func (fakeNonRuneblade) DisplayName() string      { return "StubGeneric" }
func (fakeNonRuneblade) Cost(card.GameEngine) int { return 0 }
func (fakeNonRuneblade) Pitch() int               { return 0 }
func (fakeNonRuneblade) Attack() int              { return 0 }
func (fakeNonRuneblade) Defense() int             { return 0 }
func (fakeNonRuneblade) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (fakeNonRuneblade) GoAgain(card.GameEngine) bool                       { return true }
func (fakeNonRuneblade) Play(card.GameEngine, card.Logger, *card.CardState) {}

// fakeRuneWeapon is a Runeblade weapon — tagged with Types["Weapon"] so Viserai should NOT
// trigger when it swings.
type fakeRuneWeapon struct{}

func (fakeRuneWeapon) ID() ids.CardID           { return ids.InvalidCard }
func (fakeRuneWeapon) Name() string             { return "StubRuneWeapon" }
func (fakeRuneWeapon) DisplayName() string      { return "StubRuneWeapon" }
func (fakeRuneWeapon) Cost(card.GameEngine) int { return 0 }
func (fakeRuneWeapon) Pitch() int               { return 0 }
func (fakeRuneWeapon) Attack() int              { return 0 }
func (fakeRuneWeapon) Defense() int             { return 0 }
func (fakeRuneWeapon) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeAttack)
}
func (fakeRuneWeapon) GoAgain(card.GameEngine) bool                       { return true }
func (fakeRuneWeapon) Play(card.GameEngine, card.Logger, *card.CardState) {}

// viseraiTriggerEngine returns a *GameEngine with Viserai installed and the played-this-
// turn flags pre-seeded so a FireTriggers(CardOrAbility, played) exercises the trigger
// in isolation.
func viseraiTriggerEngine(prior []card.Card, nonAttackActionPlayed bool) *gameengine.GameEngine {
	return &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().
		SetHero(Viserai).
		SetCardsPlayed(prior).
		SetNonAttackActionPlayed(nonAttackActionPlayed).
		Build()}
}

// TestViserai_RunebladeAttackAfterNonAttackActionTriggers fires the trigger on a
// Runeblade attack played after a non-attack action and asserts Viserai credits +1 value
// via CreateRunechants. The Runechant aura it creates immediately fires (Runechant's own
// filter is IsAttack and this triggering card is an attack), so the surviving aura count
// is 0 even though the value credit lands.
func TestViserai_RunebladeAttackAfterNonAttackActionTriggers(t *testing.T) {
	ge := viseraiTriggerEngine([]card.Card{fakeRuneAura{}}, true)
	ge.FireTriggers(triggertype.CardOrAbility, fakeRuneAttack{})
	if got := ge.Value(); got != 1 {
		t.Fatalf("expected +1 value from CreateRunechants, got %d", got)
	}
}

// TestViserai_RunebladeNonAttackAfterNonAttackActionTriggers fires the trigger on a
// Runeblade non-attack card. Viserai still creates a Runechant, but the Runechant aura's
// own type filter (IsAttack) blocks it from immediately firing on the non-attack card —
// so the surviving aura count is 1.
func TestViserai_RunebladeNonAttackAfterNonAttackActionTriggers(t *testing.T) {
	ge := viseraiTriggerEngine([]card.Card{fakeRuneAura{}}, true)
	ge.FireTriggers(triggertype.CardOrAbility, fakeRuneAura{})
	if got := ge.Value(); got != 1 {
		t.Fatalf("expected +1 value from CreateRunechants, got %d", got)
	}
	if got := ge.RunechantCount(); got != 1 {
		t.Fatalf("expected 1 surviving Runechant aura (non-attack triggering card), got %d", got)
	}
}

// TestViserai_NoPriorNonAttackAction confirms the handler's NonAttackActionPlayed gate
// suppresses the Runechant when the only prior play was an attack.
func TestViserai_NoPriorNonAttackAction(t *testing.T) {
	ge := viseraiTriggerEngine([]card.Card{fakeRuneAttack{}}, false)
	ge.FireTriggers(triggertype.CardOrAbility, fakeRuneAttack{})
	if got := ge.Value(); got != 0 {
		t.Fatalf("expected 0 value (no non-attack action this turn), got %d", got)
	}
}

// TestViserai_NonRunebladePlayed confirms the type filter blocks non-Runeblade cards
// before the handler runs.
func TestViserai_NonRunebladePlayed(t *testing.T) {
	ge := viseraiTriggerEngine([]card.Card{fakeRuneAura{}}, true)
	ge.FireTriggers(triggertype.CardOrAbility, fakeNonRuneblade{})
	if got := ge.Value(); got != 0 {
		t.Fatalf("expected 0 value (non-Runeblade played), got %d", got)
	}
}

// TestViserai_WeaponSwingFiltered confirms the type filter blocks Runeblade weapon
// swings — equipping or swinging a weapon isn't "playing a card".
func TestViserai_WeaponSwingFiltered(t *testing.T) {
	ge := viseraiTriggerEngine([]card.Card{fakeRuneAura{}}, true)
	ge.FireTriggers(triggertype.CardOrAbility, fakeRuneWeapon{})
	if got := ge.Value(); got != 0 {
		t.Fatalf("expected 0 value for weapon swing, got %d", got)
	}
}

// TestViserai_EmptyTurn confirms the first card of a turn doesn't fire — nothing to
// trigger on.
func TestViserai_EmptyTurn(t *testing.T) {
	ge := viseraiTriggerEngine(nil, false)
	ge.FireTriggers(triggertype.CardOrAbility, fakeRuneAura{})
	if got := ge.Value(); got != 0 {
		t.Fatalf("expected 0 value on empty turn, got %d", got)
	}
}

// Type-set vars that back the slot-classification fixtures below.
var (
	genericActionTypes   = card.NewTypeSet(card.TypeGeneric, card.TypeAction)
	actionAttackTypes    = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
	defenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
)

// nonAttackEnablerCard returns a non-attack action — fills only the non-attack-enabler
// slot (red pitch, no defense, has Go again so it doesn't extend into other slots).
func nonAttackEnablerCard(name string) card.Card {
	return testutils.NewFakeCard(name).WithTypes(genericActionTypes).WithGoAgain()
}

// defenderCard returns a Defense Reaction with positive defense — fills only the
// defender slot (red pitch, no Action subtype).
func defenderCard(name string, defense int) card.Card {
	return testutils.NewFakeCard(name).WithTypes(defenseReactionTypes).WithDefense(defense).WithPitch(1)
}

// bluePitchOnlyCard returns a non-action card with blue pitch — fills only the
// blue-pitch slot.
func bluePitchOnlyCard(name string) card.Card {
	return testutils.NewFakeCard(name).WithTypes(card.NewTypeSet(card.TypeGeneric)).WithPitch(3)
}

// noSlotCard returns an attack action with Go again, red pitch, no defense — none of the
// Viserai slots apply.
func noSlotCard(name string) card.Card {
	return testutils.NewFakeCard(name).WithTypes(actionAttackTypes).WithGoAgain().WithPitch(1)
}

// Tests that Opt(1) always tops the only revealed card.
func TestViseraiOpt_SingleCardAlwaysTop(t *testing.T) {
	c := defenderCard("d", 3)
	top, bottom := Viserai.Opt([]card.Card{c})
	if !reflect.DeepEqual(top, []card.Card{c}) {
		t.Errorf("top = %v, want [%v]", top, c)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests that two cards in the same single slot bottom the second.
func TestViseraiOpt_TwoSameSlotBottomsSecond(t *testing.T) {
	a := defenderCard("a", 3)
	b := defenderCard("b", 2)
	top, bottom := Viserai.Opt([]card.Card{a, b})
	if !reflect.DeepEqual(top, []card.Card{a}) {
		t.Errorf("top = %v, want [%v]", top, a)
	}
	if !reflect.DeepEqual(bottom, []card.Card{b}) {
		t.Errorf("bottom = %v, want [%v]", bottom, b)
	}
}

// Tests that two cards in different slots both stay on top.
func TestViseraiOpt_DifferentSlotsBothTop(t *testing.T) {
	a := nonAttackEnablerCard("a")
	b := defenderCard("b", 3)
	top, bottom := Viserai.Opt([]card.Card{a, b})
	if !reflect.DeepEqual(top, []card.Card{a, b}) {
		t.Errorf("top = %v, want [%v %v]", top, a, b)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests that a multi-slot card whose every slot is already covered gets bottomed.
func TestViseraiOpt_MultiSlotCardBottomedWhenAllCovered(t *testing.T) {
	a := nonAttackEnablerCard("a")
	bluePitch := bluePitchOnlyCard("blue")
	// b spans the non-attack-enabler and blue-pitch slots — both already covered.
	b := testutils.NewFakeCard("b").WithTypes(genericActionTypes).WithGoAgain().WithPitch(3)
	top, bottom := Viserai.Opt([]card.Card{a, bluePitch, b})
	if !reflect.DeepEqual(top, []card.Card{a, bluePitch}) {
		t.Errorf("top = %v, want [%v %v]", top, a, bluePitch)
	}
	if !reflect.DeepEqual(bottom, []card.Card{b}) {
		t.Errorf("bottom = %v, want [%v]", bottom, b)
	}
}

// Tests that a multi-slot card is bottomed if ANY of its slots overlaps a covered slot,
// even when other slots are still uncovered. Over-filling a slot wastes hand space even
// in trade for a fresh slot fill — Viserai prefers redundancy elsewhere.
func TestViseraiOpt_MultiSlotCardBottomedOnAnyOverlap(t *testing.T) {
	bluePitch := bluePitchOnlyCard("blue")
	// b is non-attack-enabler (uncovered) AND blue-pitch (covered). Bottomed because
	// blue-pitch overlaps even though the enabler slot is fresh.
	b := testutils.NewFakeCard("b").WithTypes(genericActionTypes).WithGoAgain().WithPitch(3)
	top, bottom := Viserai.Opt([]card.Card{bluePitch, b})
	if !reflect.DeepEqual(top, []card.Card{bluePitch}) {
		t.Errorf("top = %v, want [%v]", top, bluePitch)
	}
	if !reflect.DeepEqual(bottom, []card.Card{b}) {
		t.Errorf("bottom = %v, want [%v]", bottom, b)
	}
}

// Tests that cards with no slot membership stay on top regardless of order.
func TestViseraiOpt_NoSlotCardsStayTop(t *testing.T) {
	a := noSlotCard("a")
	b := noSlotCard("b")
	c := noSlotCard("c")
	top, bottom := Viserai.Opt([]card.Card{a, b, c})
	if !reflect.DeepEqual(top, []card.Card{a, b, c}) {
		t.Errorf("top = %v, want [%v %v %v]", top, a, b, c)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests Opt 3 with one card per slot category — every card belongs to a distinct slot,
// so all three stay on top.
func TestViseraiOpt_OneCardPerSlotAllKept(t *testing.T) {
	enabler := nonAttackEnablerCard("enabler")
	defender := defenderCard("defender", 3)
	bluePitch := bluePitchOnlyCard("blue")
	cs := []card.Card{enabler, defender, bluePitch}
	top, bottom := Viserai.Opt(cs)
	if !reflect.DeepEqual(top, cs) {
		t.Errorf("top = %v, want %v (one per slot, all kept)", top, cs)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests Opt 4 with two cards per slot category — the second card in each slot bottoms.
func TestViseraiOpt_DoublesInEachSlotBottomedDownToOne(t *testing.T) {
	defA := defenderCard("defA", 3)
	defB := defenderCard("defB", 2)
	bluePitchA := bluePitchOnlyCard("blueA")
	bluePitchB := bluePitchOnlyCard("blueB")
	top, bottom := Viserai.Opt([]card.Card{defA, defB, bluePitchA, bluePitchB})
	if !reflect.DeepEqual(top, []card.Card{defA, bluePitchA}) {
		t.Errorf("top = %v, want [%v %v]", top, defA, bluePitchA)
	}
	if !reflect.DeepEqual(bottom, []card.Card{defB, bluePitchB}) {
		t.Errorf("bottom = %v, want [%v %v]", bottom, defB, bluePitchB)
	}
}

// Tests that the empty input returns empty top and bottom.
func TestViseraiOpt_EmptyInput(t *testing.T) {
	top, bottom := Viserai.Opt(nil)
	if len(top) != 0 || len(bottom) != 0 {
		t.Errorf("Opt(nil) = (%v, %v), want both empty", top, bottom)
	}
}

// Tests that the defender slot keys on the type line, not on Defense > 0 — an attack
// action with a printed defense value (most attack actions in the game) doesn't fill
// the slot, so a real defender behind it still gets kept on top.
func TestViseraiOpt_DefenseValueAloneDoesNotFillDefenderSlot(t *testing.T) {
	// attackWithDefense is an attack action with positive Defense — represents the typical
	// FaB attack that doubles as a block. Should not key the defender slot.
	attackWithDefense := testutils.NewFakeCard("atkWithDef").
		WithTypes(actionAttackTypes).
		WithDefense(3).
		WithPitch(1)
	defender := defenderCard("dr", 3)
	top, bottom := Viserai.Opt([]card.Card{attackWithDefense, defender})
	if !reflect.DeepEqual(top, []card.Card{attackWithDefense, defender}) {
		t.Errorf("top = %v, want [%v %v]", top, attackWithDefense, defender)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty (DR not eclipsed by attack-with-defense)", bottom)
	}
}

// Tests that a Block-typed card fills the defender slot alongside Defense Reactions.
func TestViseraiOpt_BlockTypeFillsDefenderSlot(t *testing.T) {
	dr := defenderCard("dr", 3)
	blocker := testutils.NewFakeCard("block").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeBlock)).
		WithDefense(3).
		WithPitch(1)
	top, bottom := Viserai.Opt([]card.Card{dr, blocker})
	if !reflect.DeepEqual(top, []card.Card{dr}) {
		t.Errorf("top = %v, want [%v]", top, dr)
	}
	if !reflect.DeepEqual(bottom, []card.Card{blocker}) {
		t.Errorf("bottom = %v, want [%v] (Block-typed card competes with DR for defender slot)",
			bottom, blocker)
	}
}
