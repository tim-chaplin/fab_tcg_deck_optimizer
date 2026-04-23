package card

import "testing"

// stubCard is a minimal Card implementation for exercising TurnState helpers. Tests only care
// about identity and the Types / Attack / GoAgain hooks some helpers probe; everything else
// returns a zero value.
type stubCard struct {
	name    string
	types   TypeSet
	attack  int
	goAgain bool
}

func (c stubCard) ID() ID                        { return Invalid }
func (c stubCard) Name() string                  { return c.name }
func (stubCard) Cost(*TurnState) int             { return 0 }
func (stubCard) Pitch() int                      { return 0 }
func (c stubCard) Attack() int                   { return c.attack }
func (stubCard) Defense() int                    { return 0 }
func (c stubCard) Types() TypeSet                { return c.types }
func (c stubCard) GoAgain() bool                 { return c.goAgain }
func (stubCard) Play(*TurnState, *CardState) int { return 0 }

// TestDrawOne_AppendsTopAndAdvancesDeck: DrawOne moves the top card from Deck into Drawn and
// preserves draw order for the caller.
func TestDrawOne_AppendsTopAndAdvancesDeck(t *testing.T) {
	a, b, c := stubCard{name: "a"}, stubCard{name: "b"}, stubCard{name: "c"}
	s := &TurnState{Deck: []Card{a, b, c}}

	s.DrawOne()
	if got := len(s.Deck); got != 2 {
		t.Fatalf("Deck len = %d, want 2", got)
	}
	if s.Deck[0] != b {
		t.Errorf("Deck[0] = %v, want b (top advanced past a)", s.Deck[0])
	}
	if len(s.Drawn) != 1 || s.Drawn[0] != a {
		t.Errorf("Drawn = %v, want [a]", s.Drawn)
	}

	s.DrawOne()
	if len(s.Drawn) != 2 || s.Drawn[1] != b {
		t.Errorf("Drawn after second draw = %v, want [a, b]", s.Drawn)
	}
}

// TestDrawOne_EmptyDeckIsNoOp: with an empty deck the helper returns silently; Drawn stays
// nil so callers don't see a spurious zero-value card.
func TestDrawOne_EmptyDeckIsNoOp(t *testing.T) {
	s := &TurnState{}
	s.DrawOne()
	if len(s.Drawn) != 0 {
		t.Errorf("Drawn = %v, want empty on no-deck draw", s.Drawn)
	}
	if s.Deck != nil {
		t.Errorf("Deck = %v, want nil", s.Deck)
	}
}

// TestAddAuraTrigger_AppendsToList: AddAuraTrigger pushes each trigger onto s.AuraTriggers in
// call order — the sim reads that list to decide which handlers to fire on each
// trigger-Type condition.
func TestAddAuraTrigger_AppendsToList(t *testing.T) {
	self := stubCard{name: "self"}
	s := &TurnState{}
	s.AddAuraTrigger(AuraTrigger{Self: self, Type: TriggerStartOfTurn, Count: 2})
	s.AddAuraTrigger(AuraTrigger{Self: self, Type: TriggerStartOfTurn, Count: 1})
	if len(s.AuraTriggers) != 2 {
		t.Fatalf("AuraTriggers len = %d, want 2", len(s.AuraTriggers))
	}
	if s.AuraTriggers[0].Count != 2 || s.AuraTriggers[1].Count != 1 {
		t.Errorf("order broke: got Counts %d,%d want 2,1",
			s.AuraTriggers[0].Count, s.AuraTriggers[1].Count)
	}
	if s.AuraTriggers[0].Self != self {
		t.Errorf("Self = %v, want %v", s.AuraTriggers[0].Self, self)
	}
}

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

// TestCardState_EffectiveGoAgain: printed GoAgain OR a mid-chain grant (Mauvrion Skies et al)
// each qualifies the card for Go again. Neither printed nor granted → false.
func TestCardState_EffectiveGoAgain(t *testing.T) {
	cases := []struct {
		name    string
		printed bool
		granted bool
		want    bool
	}{
		{"neither", false, false, false},
		{"printed only", true, false, true},
		{"granted only", false, true, true},
		{"both", true, true, true},
	}
	for _, tc := range cases {
		p := &CardState{Card: stubCard{name: tc.name, goAgain: tc.printed}, GrantedGoAgain: tc.granted}
		if got := p.EffectiveGoAgain(); got != tc.want {
			t.Errorf("%s: EffectiveGoAgain() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

// TestHasPlayedType_ScansCardsPlayed: returns true when any entry in CardsPlayed has the type
// in its set, false on empty list or no matches. Pins the scan for every "if you've played
// an X this turn" rider.
func TestHasPlayedType_ScansCardsPlayed(t *testing.T) {
	aura := stubCard{name: "aura", types: NewTypeSet(TypeAura)}
	attack := stubCard{name: "attack", types: NewTypeSet(TypeAttack, TypeAction)}

	var s TurnState
	if s.HasPlayedType(TypeAura) {
		t.Error("empty CardsPlayed should return false")
	}
	s.CardsPlayed = []Card{attack, aura}
	if !s.HasPlayedType(TypeAura) {
		t.Error("Aura in CardsPlayed should be detected")
	}
	if !s.HasPlayedType(TypeAttack) {
		t.Error("Attack in CardsPlayed should be detected")
	}
	if s.HasPlayedType(TypeWeapon) {
		t.Error("Weapon not played — should return false")
	}
}

// TestHasAuraInPlay_FlagOrScan: fires on either the AuraCreated flag (Runechant creation,
// token-only auras) OR a played Aura-typed card; returns false when neither.
func TestHasAuraInPlay_FlagOrScan(t *testing.T) {
	var empty TurnState
	if empty.HasAuraInPlay() {
		t.Error("no aura, no flag → should be false")
	}

	flagged := TurnState{AuraCreated: true}
	if !flagged.HasAuraInPlay() {
		t.Error("AuraCreated=true → should be true")
	}

	playedAura := TurnState{CardsPlayed: []Card{stubCard{types: NewTypeSet(TypeAura)}}}
	if !playedAura.HasAuraInPlay() {
		t.Error("played aura card → should be true")
	}
}

// TestCreateRunechants_CountAndFlag: bumps Runechants by n, flips AuraCreated, and returns n
// as the damage-equivalent credit. n=0 is a no-op (no credit, no flag flip).
func TestCreateRunechants_CountAndFlag(t *testing.T) {
	var s TurnState
	if got := s.CreateRunechants(0); got != 0 {
		t.Errorf("CreateRunechants(0) = %d, want 0", got)
	}
	if s.AuraCreated {
		t.Error("AuraCreated should stay false for n=0")
	}
	if s.Runechants != 0 {
		t.Errorf("Runechants = %d, want 0", s.Runechants)
	}

	if got := s.CreateRunechants(3); got != 3 {
		t.Errorf("CreateRunechants(3) = %d, want 3", got)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated should flip to true")
	}
	if s.Runechants != 3 {
		t.Errorf("Runechants = %d, want 3", s.Runechants)
	}

	// Second call accumulates rather than replacing.
	s.CreateRunechants(2)
	if s.Runechants != 5 {
		t.Errorf("Runechants after second call = %d, want 5", s.Runechants)
	}
}

// TestCreateRunechant_IsSingleTokenShorthand: one call = one Runechant + 1 credit.
func TestCreateRunechant_IsSingleTokenShorthand(t *testing.T) {
	var s TurnState
	if got := s.CreateRunechant(); got != 1 {
		t.Errorf("CreateRunechant() = %d, want 1", got)
	}
	if s.Runechants != 1 {
		t.Errorf("Runechants = %d, want 1", s.Runechants)
	}
	if !s.AuraCreated {
		t.Error("AuraCreated should be true after CreateRunechant")
	}
}

// TestAddToGraveyard_AppendsInOrder: graveyard entries appear in append order so downstream
// readers (aura-banish handlers, leave-trigger scanners) see a stable destroy sequence.
func TestAddToGraveyard_AppendsInOrder(t *testing.T) {
	a, b := stubCard{name: "a"}, stubCard{name: "b"}
	var s TurnState
	s.AddToGraveyard(a)
	s.AddToGraveyard(b)
	if len(s.Graveyard) != 2 || s.Graveyard[0] != a || s.Graveyard[1] != b {
		t.Errorf("Graveyard = %v, want [a, b]", s.Graveyard)
	}
}

// TestClashValue_WinTieLose: rule 8.5.45 models the clash by comparing our top card's attack
// to the opponent's (approximated at 5). 6+ wins (credit bonus), 5 ties (0), <5 loses
// (−bonus). Empty deck short-circuits to 0 per rule 8.5.45d.
func TestClashValue_WinTieLose(t *testing.T) {
	cases := []struct {
		name    string
		topAtk  int
		deckLen int
		bonus   int
		want    int
	}{
		{"win on 6", 6, 1, 2, 2},
		{"win on 7", 7, 1, 2, 2},
		{"tie on 5", 5, 1, 2, 0},
		{"lose on 4", 4, 1, 2, -2},
		{"lose on 0", 0, 1, 3, -3},
		{"empty deck short-circuits", 99, 0, 5, 0},
	}
	for _, tc := range cases {
		var s TurnState
		if tc.deckLen > 0 {
			s.Deck = []Card{stubCard{attack: tc.topAtk}}
		}
		if got := ClashValue(&s, tc.bonus); got != tc.want {
			t.Errorf("%s: ClashValue = %d, want %d", tc.name, got, tc.want)
		}
	}
}

// TestLikelyToHit_OnlyAwkwardAmounts: 1 / 4 / 7 damage slip past typical blocks (since cards
// are ~3 points of value, opponents won't over-pay with a 3-block to soak 1 damage, etc.).
// Everything else the sim treats as reliably blockable.
func TestLikelyToHit_OnlyAwkwardAmounts(t *testing.T) {
	for _, n := range []int{1, 4, 7} {
		if !LikelyToHit(n) {
			t.Errorf("LikelyToHit(%d) = false, want true (awkward amount)", n)
		}
	}
	for _, n := range []int{0, 2, 3, 5, 6, 8, 10} {
		if LikelyToHit(n) {
			t.Errorf("LikelyToHit(%d) = true, want false", n)
		}
	}
}
