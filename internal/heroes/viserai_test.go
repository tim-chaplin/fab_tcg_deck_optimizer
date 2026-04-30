package heroes

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// stubRuneAttack is a minimal Runeblade attack-action card.
type stubRuneAttack struct{}

func (stubRuneAttack) ID() ids.CardID          { return ids.InvalidCard }
func (stubRuneAttack) Name() string            { return "StubRuneAttack" }
func (stubRuneAttack) Cost(*sim.TurnState) int { return 0 }
func (stubRuneAttack) Pitch() int              { return 0 }
func (stubRuneAttack) Attack() int             { return 0 }
func (stubRuneAttack) Defense() int            { return 0 }
func (stubRuneAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (stubRuneAttack) GoAgain() bool                       { return true }
func (stubRuneAttack) Play(*sim.TurnState, *sim.CardState) {}

// stubRuneAura is a minimal Runeblade non-attack action (an Aura).
type stubRuneAura struct{}

func (stubRuneAura) ID() ids.CardID          { return ids.InvalidCard }
func (stubRuneAura) Name() string            { return "StubRuneAura" }
func (stubRuneAura) Cost(*sim.TurnState) int { return 0 }
func (stubRuneAura) Pitch() int              { return 0 }
func (stubRuneAura) Attack() int             { return 0 }
func (stubRuneAura) Defense() int            { return 0 }
func (stubRuneAura) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAura)
}
func (stubRuneAura) GoAgain() bool                       { return true }
func (stubRuneAura) Play(*sim.TurnState, *sim.CardState) {}

// stubNonRuneblade is an Action-Attack with no Runeblade type — should never trigger Viserai.
type stubNonRuneblade struct{}

func (stubNonRuneblade) ID() ids.CardID          { return ids.InvalidCard }
func (stubNonRuneblade) Name() string            { return "StubGeneric" }
func (stubNonRuneblade) Cost(*sim.TurnState) int { return 0 }
func (stubNonRuneblade) Pitch() int              { return 0 }
func (stubNonRuneblade) Attack() int             { return 0 }
func (stubNonRuneblade) Defense() int            { return 0 }
func (stubNonRuneblade) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubNonRuneblade) GoAgain() bool                       { return true }
func (stubNonRuneblade) Play(*sim.TurnState, *sim.CardState) {}
func TestViserai_RunebladeAfterNonAttackActionTriggers(t *testing.T) {
	// Non-attack action played first, then a Runeblade attack. Viserai's OnCardPlayed creates a
	// Runechant token: returns +1 damage (each token credited +1 at creation) and leaves a
	// token on state.Runechants for downstream consume or carryover. NonAttackActionPlayed is
	// maintained by the attack-chain driver as non-attack actions resolve; callers must set it
	// when seeding a TurnState for trigger checks.
	s := sim.TurnState{CardsPlayed: []sim.Card{stubRuneAura{}}, NonAttackActionPlayed: true}
	if got := (Viserai{}).OnCardPlayed(stubRuneAttack{}, &s); got != 1 {
		t.Fatalf("expected +1 damage from OnCardPlayed, got %d", got)
	}
	if s.Runechants != 1 {
		t.Fatalf("expected 1 Runechant on state, got %d", s.Runechants)
	}
}

func TestViserai_NoPriorNonAttackAction(t *testing.T) {
	// Runeblade card, but the only prior play was an attack — no trigger.
	s := sim.TurnState{CardsPlayed: []sim.Card{stubRuneAttack{}}}
	if got := (Viserai{}).OnCardPlayed(stubRuneAttack{}, &s); got != 0 {
		t.Fatalf("expected 0 (no non-attack action in CardsPlayed), got %d", got)
	}
}

func TestViserai_CardStateNotRuneblade(t *testing.T) {
	// Played card isn't Runeblade — Viserai's ability doesn't trigger even if a non-attack action was
	// played earlier.
	s := sim.TurnState{CardsPlayed: []sim.Card{stubRuneAura{}}, NonAttackActionPlayed: true}
	if got := (Viserai{}).OnCardPlayed(stubNonRuneblade{}, &s); got != 0 {
		t.Fatalf("expected 0 (non-Runeblade played), got %d", got)
	}
}

// stubRuneWeapon is a Runeblade weapon — tagged with Types["Weapon"] so Viserai should NOT trigger
// when it swings.
type stubRuneWeapon struct{}

func (stubRuneWeapon) ID() ids.CardID          { return ids.InvalidCard }
func (stubRuneWeapon) Name() string            { return "StubRuneWeapon" }
func (stubRuneWeapon) Cost(*sim.TurnState) int { return 0 }
func (stubRuneWeapon) Pitch() int              { return 0 }
func (stubRuneWeapon) Attack() int             { return 0 }
func (stubRuneWeapon) Defense() int            { return 0 }
func (stubRuneWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (stubRuneWeapon) GoAgain() bool                       { return true }
func (stubRuneWeapon) Play(*sim.TurnState, *sim.CardState) {}
func TestViserai_WeaponSwingDoesNotTrigger(t *testing.T) {
	// Even with a prior non-attack action in CardsPlayed, swinging a Runeblade weapon isn't "playing a
	// card" and must not trigger.
	s := sim.TurnState{CardsPlayed: []sim.Card{stubRuneAura{}}, NonAttackActionPlayed: true}
	if got := (Viserai{}).OnCardPlayed(stubRuneWeapon{}, &s); got != 0 {
		t.Fatalf("expected 0 for weapon swing, got %d", got)
	}
}

func TestViserai_EmptyTurn(t *testing.T) {
	// First card of the turn: no prior plays, nothing to trigger on.
	var s sim.TurnState
	if got := (Viserai{}).OnCardPlayed(stubRuneAura{}, &s); got != 0 {
		t.Fatalf("expected 0 on empty turn, got %d", got)
	}
}

// genericActionTypes / actionAttackTypes / defenseReactionTypes back the slot-classification
// fixtures below — kept terse so the test bodies read as scenario tables, not boilerplate.
var (
	genericActionTypes   = card.NewTypeSet(card.TypeGeneric, card.TypeAction)
	actionAttackTypes    = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
	defenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
)

// nonAttackEnablerCard returns a non-attack action — slot 1 only (no defense, red pitch,
// no go-again so it would also fit slot 2; we suppress slot 2 by giving it Go again).
func nonAttackEnablerCard(name string) sim.Card {
	return testutils.NewStubCard(name).WithTypes(genericActionTypes).WithGoAgain()
}

// nonGoAgainActionCard returns an attack action without Go again — slot 2 only.
func nonGoAgainActionCard(name string) sim.Card {
	return testutils.NewStubCard(name).WithTypes(actionAttackTypes)
}

// defenderCard returns a defense reaction with positive defense — slot 3 only (no Action
// type → not slots 1 or 2; red pitch → not slot 4).
func defenderCard(name string, defense int) sim.Card {
	return testutils.NewStubCard(name).WithTypes(defenseReactionTypes).WithDefense(defense).WithPitch(1)
}

// bluePitchOnlyCard returns a non-action card with blue pitch — slot 4 only.
func bluePitchOnlyCard(name string) sim.Card {
	return testutils.NewStubCard(name).WithTypes(card.NewTypeSet(card.TypeGeneric)).WithPitch(3)
}

// noSlotCard returns an attack action with Go again, red pitch, no defense — none of the
// four Viserai slots apply.
func noSlotCard(name string) sim.Card {
	return testutils.NewStubCard(name).WithTypes(actionAttackTypes).WithGoAgain().WithPitch(1)
}

// Tests that Opt(1) always tops the only revealed card.
func TestViseraiOpt_SingleCardAlwaysTop(t *testing.T) {
	c := defenderCard("d", 3)
	top, bottom := (Viserai{}).Opt([]sim.Card{c})
	if !reflect.DeepEqual(top, []sim.Card{c}) {
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
	top, bottom := (Viserai{}).Opt([]sim.Card{a, b})
	if !reflect.DeepEqual(top, []sim.Card{a}) {
		t.Errorf("top = %v, want [%v]", top, a)
	}
	if !reflect.DeepEqual(bottom, []sim.Card{b}) {
		t.Errorf("bottom = %v, want [%v]", bottom, b)
	}
}

// Tests that two cards in different slots both stay on top.
func TestViseraiOpt_DifferentSlotsBothTop(t *testing.T) {
	a := nonAttackEnablerCard("a")
	b := defenderCard("b", 3)
	top, bottom := (Viserai{}).Opt([]sim.Card{a, b})
	if !reflect.DeepEqual(top, []sim.Card{a, b}) {
		t.Errorf("top = %v, want [%v %v]", top, a, b)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests that a card whose every slot is already covered gets bottomed even when it spans
// multiple slots.
func TestViseraiOpt_MultiSlotCardBottomedWhenAllCovered(t *testing.T) {
	a := nonAttackEnablerCard("a")
	defender := defenderCard("def", 3)
	// b is non-attack-enabler AND defender — both slots already covered by a + defender.
	b := testutils.NewStubCard("b").
		WithTypes(genericActionTypes).
		WithGoAgain().
		WithDefense(3).
		WithPitch(1)
	top, bottom := (Viserai{}).Opt([]sim.Card{a, defender, b})
	if !reflect.DeepEqual(top, []sim.Card{a, defender}) {
		t.Errorf("top = %v, want [%v %v]", top, a, defender)
	}
	if !reflect.DeepEqual(bottom, []sim.Card{b}) {
		t.Errorf("bottom = %v, want [%v]", bottom, b)
	}
}

// Tests that a multi-slot card is kept when it covers at least one new slot, even if some
// of its slots are already covered.
func TestViseraiOpt_MultiSlotCardKeptWhenAnySlotNew(t *testing.T) {
	defender := defenderCard("def", 3)
	// b is non-attack-enabler (uncovered) AND defender (covered). Should still be kept
	// because nonAttackEnabler is fresh.
	b := testutils.NewStubCard("b").
		WithTypes(genericActionTypes).
		WithGoAgain().
		WithDefense(3).
		WithPitch(1)
	top, bottom := (Viserai{}).Opt([]sim.Card{defender, b})
	if !reflect.DeepEqual(top, []sim.Card{defender, b}) {
		t.Errorf("top = %v, want [%v %v]", top, defender, b)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests that cards with no slot membership stay on top regardless of order.
func TestViseraiOpt_NoSlotCardsStayTop(t *testing.T) {
	a := noSlotCard("a")
	b := noSlotCard("b")
	top, bottom := (Viserai{}).Opt([]sim.Card{a, b})
	if !reflect.DeepEqual(top, []sim.Card{a, b}) {
		t.Errorf("top = %v, want [%v %v]", top, a, b)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests Opt 4 with one card per slot category — every card belongs to a distinct slot,
// so all four stay on top.
func TestViseraiOpt_OneCardPerSlotAllKept(t *testing.T) {
	enabler := nonAttackEnablerCard("enabler")
	finisher := nonGoAgainActionCard("finisher")
	defender := defenderCard("defender", 3)
	bluePitch := bluePitchOnlyCard("blue")
	cards := []sim.Card{enabler, finisher, defender, bluePitch}
	top, bottom := (Viserai{}).Opt(cards)
	if !reflect.DeepEqual(top, cards) {
		t.Errorf("top = %v, want %v (one per slot, all kept)", top, cards)
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
	top, bottom := (Viserai{}).Opt([]sim.Card{defA, defB, bluePitchA, bluePitchB})
	if !reflect.DeepEqual(top, []sim.Card{defA, bluePitchA}) {
		t.Errorf("top = %v, want [%v %v]", top, defA, bluePitchA)
	}
	if !reflect.DeepEqual(bottom, []sim.Card{defB, bluePitchB}) {
		t.Errorf("bottom = %v, want [%v %v]", bottom, defB, bluePitchB)
	}
}

// Tests that the empty input returns empty top and bottom.
func TestViseraiOpt_EmptyInput(t *testing.T) {
	top, bottom := (Viserai{}).Opt(nil)
	if len(top) != 0 || len(bottom) != 0 {
		t.Errorf("Opt(nil) = (%v, %v), want both empty", top, bottom)
	}
}
