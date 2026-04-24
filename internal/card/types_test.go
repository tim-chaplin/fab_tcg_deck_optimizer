package card

import "testing"

// TestNewTypeSet_UnionsAllArguments pins the variadic OR: the returned set contains every type
// passed and no others.
func TestNewTypeSet_UnionsAllArguments(t *testing.T) {
	s := NewTypeSet(TypeAction, TypeAttack, TypeRuneblade)
	for _, want := range []CardType{TypeAction, TypeAttack, TypeRuneblade} {
		if !s.Has(want) {
			t.Errorf("Has(%v) = false, want true", want)
		}
	}
	for _, notWant := range []CardType{TypeAura, TypeWeapon, TypeInstant} {
		if s.Has(notWant) {
			t.Errorf("Has(%v) = true, want false", notWant)
		}
	}
	if NewTypeSet() != 0 {
		t.Errorf("NewTypeSet() with no args = %d, want 0", NewTypeSet())
	}
}

// TestTypeSet_PersistsInPlay pins the solver's zone-routing decision: Aura / Item / Weapon stay
// in play, everything else heads to the graveyard. The bitmask is what makes the post-Play
// graveyard append conditional.
func TestTypeSet_PersistsInPlay(t *testing.T) {
	cases := []struct {
		name string
		set  TypeSet
		want bool
	}{
		{"aura", NewTypeSet(TypeAura), true},
		{"item", NewTypeSet(TypeItem), true},
		{"weapon", NewTypeSet(TypeWeapon), true},
		{"runeblade action aura", NewTypeSet(TypeRuneblade, TypeAction, TypeAura), true},
		{"attack action", NewTypeSet(TypeAttack, TypeAction), false},
		{"defense reaction", NewTypeSet(TypeDefenseReaction), false},
		{"plain action", NewTypeSet(TypeAction), false},
		{"empty", 0, false},
	}
	for _, tc := range cases {
		if got := tc.set.PersistsInPlay(); got != tc.want {
			t.Errorf("%s.PersistsInPlay() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestTypeSet_IsNonAttackAction pins the "non-attack action played" predicate six riders key
// on (Viserai's trigger, Aether Slash's arcane rider, etc.). Pure Action fires; AttackAction
// doesn't; non-Action types never fire regardless of other keywords.
func TestTypeSet_IsNonAttackAction(t *testing.T) {
	cases := []struct {
		name string
		set  TypeSet
		want bool
	}{
		{"plain action", NewTypeSet(TypeAction), true},
		{"action aura", NewTypeSet(TypeAction, TypeAura), true},
		{"runeblade action", NewTypeSet(TypeRuneblade, TypeAction), true},
		{"attack action", NewTypeSet(TypeAction, TypeAttack), false},
		{"pure attack", NewTypeSet(TypeAttack), false},
		{"aura only", NewTypeSet(TypeAura), false},
		{"empty", 0, false},
	}
	for _, tc := range cases {
		if got := tc.set.IsNonAttackAction(); got != tc.want {
			t.Errorf("%s.IsNonAttackAction() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestTypeSet_IsAttackAction pins the predicate every "next attack action card you play this
// turn" rider walks CardsRemaining to find (Come to Fight, Minnowism, Nimblism, Sloggism, Water
// the Seeds, Captain's Call, Flying High, Trot Along, Scout the Periphery). Requires both
// Action and Attack; either alone doesn't qualify.
func TestTypeSet_IsAttackAction(t *testing.T) {
	cases := []struct {
		name string
		set  TypeSet
		want bool
	}{
		{"attack action", NewTypeSet(TypeAction, TypeAttack), true},
		{"runeblade attack action", NewTypeSet(TypeRuneblade, TypeAction, TypeAttack), true},
		{"plain action", NewTypeSet(TypeAction), false},
		{"pure attack", NewTypeSet(TypeAttack), false},
		{"runeblade weapon", NewTypeSet(TypeRuneblade, TypeWeapon, TypeAttack), false},
		{"aura only", NewTypeSet(TypeAura), false},
		{"empty", 0, false},
	}
	for _, tc := range cases {
		if got := tc.set.IsAttackAction(); got != tc.want {
			t.Errorf("%s.IsAttackAction() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestTypeSet_IsRunebladeAttack: every "next Runeblade attack this turn" rider keys on this
// helper. Requires both TypeRuneblade AND (TypeAttack | TypeWeapon). Plain Runeblade auras
// (no attack / weapon) don't qualify.
func TestTypeSet_IsRunebladeAttack(t *testing.T) {
	cases := []struct {
		name string
		set  TypeSet
		want bool
	}{
		{"runeblade attack action", NewTypeSet(TypeRuneblade, TypeAttack, TypeAction), true},
		{"runeblade weapon", NewTypeSet(TypeRuneblade, TypeWeapon), true},
		{"runeblade aura", NewTypeSet(TypeRuneblade, TypeAction, TypeAura), false},
		{"generic attack", NewTypeSet(TypeGeneric, TypeAttack, TypeAction), false},
		{"weapon alone", NewTypeSet(TypeWeapon), false},
		{"empty", 0, false},
	}
	for _, tc := range cases {
		if got := tc.set.IsRunebladeAttack(); got != tc.want {
			t.Errorf("%s.IsRunebladeAttack() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestTypeSet_IsDefenseReaction pins the DR bit read by five solver sites. Any set containing
// TypeDefenseReaction is a DR, regardless of whatever else is on the type line.
func TestTypeSet_IsDefenseReaction(t *testing.T) {
	if !NewTypeSet(TypeDefenseReaction).IsDefenseReaction() {
		t.Error("DR set should report IsDefenseReaction = true")
	}
	if !NewTypeSet(TypeDefenseReaction, TypeRuneblade).IsDefenseReaction() {
		t.Error("Runeblade DR should still report true")
	}
	if NewTypeSet(TypeAttack, TypeAction).IsDefenseReaction() {
		t.Error("AttackAction should not be a DR")
	}
	if TypeSet(0).IsDefenseReaction() {
		t.Error("empty set should not be a DR")
	}
}
