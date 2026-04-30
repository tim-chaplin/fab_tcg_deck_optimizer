package heroes

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
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

// Type-set vars that back the slot-classification fixtures below.
var (
	genericActionTypes   = card.NewTypeSet(card.TypeGeneric, card.TypeAction)
	actionAttackTypes    = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
	defenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
)

// nonAttackEnablerCard returns a non-attack action — fills only the non-attack-enabler
// slot (red pitch, no defense, has Go again so it doesn't extend into other slots).
func nonAttackEnablerCard(name string) sim.Card {
	return testutils.NewStubCard(name).WithTypes(genericActionTypes).WithGoAgain()
}

// defenderCard returns a Defense Reaction with positive defense — fills only the
// defender slot (red pitch, no Action subtype).
func defenderCard(name string, defense int) sim.Card {
	return testutils.NewStubCard(name).WithTypes(defenseReactionTypes).WithDefense(defense).WithPitch(1)
}

// bluePitchOnlyCard returns a non-action card with blue pitch — fills only the
// blue-pitch slot.
func bluePitchOnlyCard(name string) sim.Card {
	return testutils.NewStubCard(name).WithTypes(card.NewTypeSet(card.TypeGeneric)).WithPitch(3)
}

// noSlotCard returns an attack action with Go again, red pitch, no defense — none of the
// Viserai slots apply.
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

// Tests that a multi-slot card whose every slot is already covered gets bottomed.
func TestViseraiOpt_MultiSlotCardBottomedWhenAllCovered(t *testing.T) {
	a := nonAttackEnablerCard("a")
	bluePitch := bluePitchOnlyCard("blue")
	// b spans the non-attack-enabler and blue-pitch slots — both already covered.
	b := testutils.NewStubCard("b").WithTypes(genericActionTypes).WithGoAgain().WithPitch(3)
	top, bottom := (Viserai{}).Opt([]sim.Card{a, bluePitch, b})
	if !reflect.DeepEqual(top, []sim.Card{a, bluePitch}) {
		t.Errorf("top = %v, want [%v %v]", top, a, bluePitch)
	}
	if !reflect.DeepEqual(bottom, []sim.Card{b}) {
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
	b := testutils.NewStubCard("b").WithTypes(genericActionTypes).WithGoAgain().WithPitch(3)
	top, bottom := (Viserai{}).Opt([]sim.Card{bluePitch, b})
	if !reflect.DeepEqual(top, []sim.Card{bluePitch}) {
		t.Errorf("top = %v, want [%v]", top, bluePitch)
	}
	if !reflect.DeepEqual(bottom, []sim.Card{b}) {
		t.Errorf("bottom = %v, want [%v]", bottom, b)
	}
}

// Tests that cards with no slot membership stay on top regardless of order.
func TestViseraiOpt_NoSlotCardsStayTop(t *testing.T) {
	a := noSlotCard("a")
	b := noSlotCard("b")
	c := noSlotCard("c")
	top, bottom := (Viserai{}).Opt([]sim.Card{a, b, c})
	if !reflect.DeepEqual(top, []sim.Card{a, b, c}) {
		t.Errorf("top = %v, want [%v %v %v]", top, a, b, c)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty", bottom)
	}
}

// Tests that three Runerager Swarms stay on top — Runerager Swarm has only conditional
// Go again (printed GoAgain() == false) but in Viserai's deck it almost always extends
// the chain via the runechant condition, so we don't bottom it for over-filling a
// "non-go-again" slot. Pinned via the real card; if Runerager Swarm's profile ever
// shifts into a Viserai-Opt slot, this test trips.
func TestViseraiOpt_ThreeRuneragerSwarmsAllKept(t *testing.T) {
	a, b, c := cards.RuneragerSwarmRed{}, cards.RuneragerSwarmRed{}, cards.RuneragerSwarmRed{}
	top, bottom := (Viserai{}).Opt([]sim.Card{a, b, c})
	if len(top) != 3 {
		t.Errorf("top = %v, want 3 Runerager Swarms", top)
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
	cs := []sim.Card{enabler, defender, bluePitch}
	top, bottom := (Viserai{}).Opt(cs)
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

// Tests that the defender slot keys on the type line, not on Defense > 0 — an attack
// action with a printed defense value (most attack actions in the game) doesn't fill
// the slot, so a real defender behind it still gets kept on top.
func TestViseraiOpt_DefenseValueAloneDoesNotFillDefenderSlot(t *testing.T) {
	// attackWithDefense is an attack action with positive Defense — represents the typical
	// FaB attack that doubles as a block. Should not key the defender slot.
	attackWithDefense := testutils.NewStubCard("atkWithDef").
		WithTypes(actionAttackTypes).
		WithDefense(3).
		WithPitch(1)
	defender := defenderCard("dr", 3)
	top, bottom := (Viserai{}).Opt([]sim.Card{attackWithDefense, defender})
	if !reflect.DeepEqual(top, []sim.Card{attackWithDefense, defender}) {
		t.Errorf("top = %v, want [%v %v]", top, attackWithDefense, defender)
	}
	if len(bottom) != 0 {
		t.Errorf("bottom = %v, want empty (DR not eclipsed by attack-with-defense)", bottom)
	}
}

// Tests that a Block-typed card fills the defender slot alongside Defense Reactions.
func TestViseraiOpt_BlockTypeFillsDefenderSlot(t *testing.T) {
	dr := defenderCard("dr", 3)
	blocker := testutils.NewStubCard("block").
		WithTypes(card.NewTypeSet(card.TypeGeneric, card.TypeBlock)).
		WithDefense(3).
		WithPitch(1)
	top, bottom := (Viserai{}).Opt([]sim.Card{dr, blocker})
	if !reflect.DeepEqual(top, []sim.Card{dr}) {
		t.Errorf("top = %v, want [%v]", top, dr)
	}
	if !reflect.DeepEqual(bottom, []sim.Card{blocker}) {
		t.Errorf("bottom = %v, want [%v] (Block-typed card competes with DR for defender slot)",
			bottom, blocker)
	}
}
