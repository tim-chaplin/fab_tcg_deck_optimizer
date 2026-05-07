package sim_test

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestDrawOne_AppendsTopAndAdvancesDeck: DrawOne pops the top of the deck and appends it to
// Hand, preserving draw order for downstream effects.
func TestDrawOne_AppendsTopAndAdvancesDeck(t *testing.T) {
	a, b, c := testutils.NewStubCard("a"), testutils.NewStubCard("b"), testutils.NewStubCard("c")
	s := NewTurnState([]Card{a, b, c}, nil)

	s.DrawOne()
	if got := len(s.Deck()); got != 2 {
		t.Fatalf("Deck len = %d, want 2", got)
	}
	if s.Deck()[0] != b {
		t.Errorf("Deck[0] = %v, want b (top advanced past a)", s.Deck()[0])
	}
	if h := s.Hand(); len(h) != 1 || h[0] != a {
		t.Errorf("Hand = %v, want [a]", h)
	}

	s.DrawOne()
	if h := s.Hand(); len(h) != 2 || h[1] != b {
		t.Errorf("Hand after second draw = %v, want [a, b]", h)
	}
}

// TestDrawOne_EmptyDeckIsNoOp: with an empty deck the helper returns silently; Hand stays
// untouched so callers don't see a spurious zero-value card.
func TestDrawOne_EmptyDeckIsNoOp(t *testing.T) {
	s := &TurnState{}
	s.DrawOne()
	if h := s.Hand(); len(h) != 0 {
		t.Errorf("Hand = %v, want empty on no-deck draw", h)
	}
	if s.Deck() != nil {
		t.Errorf("Deck = %v, want nil", s.Deck())
	}
}

// Tests that AddAura flips AuraCreated and appends to s.Auras in call order — bundling
// the flip with the append prevents a trigger registered without advertising the aura.
func TestAddAura_FlipsAuraCreatedAndAppends(t *testing.T) {
	self := testutils.NewStubCard("self")
	s := &TurnState{}
	if s.AuraCreated {
		t.Fatal("pre: AuraCreated should be false")
	}
	s.AddAura(Aura{Trigger: Trigger{TriggerType: TriggerStartOfTurn}, Self: CardOrTokenType{Card: self}, Count: 2})
	s.AddAura(Aura{Trigger: Trigger{TriggerType: TriggerStartOfTurn}, Self: CardOrTokenType{Card: self}, Count: 1})
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
	if len(s.Auras) != 2 {
		t.Fatalf("Auras len = %d, want 2", len(s.Auras))
	}
	if s.Auras[0].Count != 2 || s.Auras[1].Count != 1 {
		t.Errorf("order broke: got Counts %d,%d want 2,1",
			s.Auras[0].Count, s.Auras[1].Count)
	}
	if s.Auras[0].Self.Card != self {
		t.Errorf("Self.Card = %v, want %v", s.Auras[0].Self.Card, self)
	}
}

// TestHasPlayedType_ScansCardsPlayed: returns true when any entry in CardsPlayed has the type
// in its set, false on empty list or no matches. Pins the scan for every "if you've played
// an X this turn" rider.
func TestHasPlayedType_ScansCardsPlayed(t *testing.T) {
	aura := testutils.NewStubCard("aura").WithTypes(card.NewTypeSet(card.TypeAura))
	attack := testutils.NewStubCard("attack").WithTypes(card.NewTypeSet(card.TypeAttack, card.TypeAction))

	var s TurnState
	if s.HasPlayedType(card.TypeAura) {
		t.Error("empty CardsPlayed should return false")
	}
	s.CardsPlayed = []Card{attack, aura}
	if !s.HasPlayedType(card.TypeAura) {
		t.Error("Aura in CardsPlayed should be detected")
	}
	if !s.HasPlayedType(card.TypeAttack) {
		t.Error("Attack in CardsPlayed should be detected")
	}
	if s.HasPlayedType(card.TypeWeapon) {
		t.Error("Weapon not played — should return false")
	}
}

// TestHasPlayedOrCreatedAura_FlagOrScan: fires on either the AuraCreated flag (Runechant
// creation, token-only auras) OR a played Aura-typed card; returns false when neither.
func TestHasPlayedOrCreatedAura_FlagOrScan(t *testing.T) {
	var empty TurnState
	if empty.HasPlayedOrCreatedAura() {
		t.Error("no aura, no flag → should be false")
	}

	flagged := TurnState{AuraCreated: true}
	if !flagged.HasPlayedOrCreatedAura() {
		t.Error("AuraCreated=true → should be true")
	}

	playedAura := TurnState{
		CardsPlayed: []Card{testutils.NewStubCard("aura").WithTypes(card.NewTypeSet(card.TypeAura))},
	}
	if !playedAura.HasPlayedOrCreatedAura() {
		t.Error("played aura card → should be true")
	}
}

// TestAddValue_ClampsNonPositive: the helper sums positive credits into Value and is a
// no-op for n <= 0. Negative grants (debuffs) and zero (no-effect Plays) must not subtract
// from the running total.
func TestAddValue_ClampsNonPositive(t *testing.T) {
	cases := []struct {
		name string
		bump int
		want int
	}{
		{"positive accumulates", 3, 3},
		{"zero is no-op", 0, 0},
		{"negative is no-op", -5, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s TurnState
			s.AddValue(tc.bump)
			if s.Value != tc.want {
				t.Errorf("Value = %d, want %d", s.Value, tc.want)
			}
		})
	}
	// Mixed sequence: positives accumulate, non-positives pass through.
	var s TurnState
	s.AddValue(2)
	s.AddValue(-10)
	s.AddValue(0)
	s.AddValue(5)
	if s.Value != 7 {
		t.Errorf("after mixed sequence Value = %d, want 7 (2+5; -10/0 clamped)", s.Value)
	}
}

// TestCreateRunechants_CountAndFlag: bumps the runechant aura by n, flips AuraCreated,
// and credits +n damage to Value. n=0 is a no-op (no aura, no flag flip, no credit).
func TestCreateRunechants_CountAndFlag(t *testing.T) {
	var s TurnState
	s.CreateRunechants(0)
	if s.AuraCreated {
		t.Error("AuraCreated should stay false for n=0")
	}
	if s.Runechants() != 0 {
		t.Errorf("Runechants = %d, want 0", s.Runechants())
	}
	if s.Value != 0 {
		t.Errorf("Value = %d, want 0 (n=0 credits nothing)", s.Value)
	}

	s.CreateRunechants(3)
	if !s.AuraCreated {
		t.Error("AuraCreated should flip to true")
	}
	if s.Runechants() != 3 {
		t.Errorf("Runechants = %d, want 3", s.Runechants())
	}
	if s.Value != 3 {
		t.Errorf("Value = %d, want 3 (creator credits +n)", s.Value)
	}

	// Second call accumulates rather than replacing.
	s.CreateRunechants(2)
	if s.Runechants() != 5 {
		t.Errorf("Runechants after second call = %d, want 5", s.Runechants())
	}
	if s.Value != 5 {
		t.Errorf("Value after second call = %d, want 5", s.Value)
	}
}

// TestCreateRunechants_OneToken: one call = one Runechant + 1 credit.
func TestCreateRunechants_OneToken(t *testing.T) {
	var s TurnState
	s.CreateRunechants(1)
	if s.Runechants() != 1 {
		t.Errorf("Runechants = %d, want 1", s.Runechants())
	}
	if !s.AuraCreated {
		t.Error("AuraCreated should be true after CreateRunechants(1)")
	}
	if s.Value != 1 {
		t.Errorf("Value = %d, want 1", s.Value)
	}
}

// TestCreatePonder_BumpsCountAndFlag: bumps Ponders by n, flips AuraCreated, credits
// no Value, and accumulates on repeat calls. n=0 is a no-op.
func TestCreatePonder_BumpsCountAndFlag(t *testing.T) {
	var s TurnState
	s.CreatePonder(0)
	if s.AuraCreated {
		t.Error("AuraCreated should stay false for n=0")
	}
	if s.Ponders() != 0 {
		t.Errorf("Ponders = %d, want 0", s.Ponders())
	}

	s.CreatePonder(2)
	if !s.AuraCreated {
		t.Error("AuraCreated should flip to true")
	}
	if s.Ponders() != 2 {
		t.Errorf("Ponders = %d, want 2", s.Ponders())
	}
	if s.Value != 0 {
		t.Errorf("Value = %d, want 0 (Ponder is zero-value at creation)", s.Value)
	}

	// Second call accumulates rather than replacing.
	s.CreatePonder(3)
	if s.Ponders() != 5 {
		t.Errorf("Ponders after second call = %d, want 5", s.Ponders())
	}
}

// TestAddToGraveyard_AppendsInOrder: graveyard entries appear in append order so downstream
// readers (aura-banish handlers, leave-trigger scanners) see a stable destroy sequence.
func TestAddToGraveyard_AppendsInOrder(t *testing.T) {
	a, b := testutils.NewStubCard("a"), testutils.NewStubCard("b")
	var s TurnState
	s.AddToGraveyard(a)
	s.AddToGraveyard(b)
	g := s.Graveyard()
	if len(g) != 2 || g[0] != a || g[1] != b {
		t.Errorf("Graveyard = %v, want [a, b]", g)
	}
}

// TestClash_WinTieLose: rule 8.5.45 models the clash by comparing our top card's attack
// to the opponent's (approximated at 5). 6+ fires win, 5 ties (neither callback), <5 fires
// lose. Empty deck short-circuits to neither callback per rule 8.5.45d.
func TestClash_WinTieLose(t *testing.T) {
	cases := []struct {
		name     string
		topAtk   int
		deckLen  int
		wantWin  int
		wantLose int
	}{
		{"win on 6", 6, 1, 1, 0},
		{"win on 7", 7, 1, 1, 0},
		{"tie on 5", 5, 1, 0, 0},
		{"lose on 4", 4, 1, 0, 1},
		{"lose on 0", 0, 1, 0, 1},
		{"empty deck short-circuits", 99, 0, 0, 0},
	}
	for _, tc := range cases {
		var s *TurnState
		if tc.deckLen > 0 {
			s = NewTurnState([]Card{testutils.NewStubCard("top").WithAttack(tc.topAtk)}, nil)
		} else {
			s = &TurnState{}
		}
		var winN, loseN int
		s.Clash(func() { winN++ }, func() { loseN++ })
		if winN != tc.wantWin || loseN != tc.wantLose {
			t.Errorf("%s: win=%d/%d lose=%d/%d", tc.name, winN, tc.wantWin, loseN, tc.wantLose)
		}
	}
}

// Tests that NewTurnState seeds cacheable=true (the zero-value default is false).
func TestIsCacheable_NewTurnStateSeedsCacheable(t *testing.T) {
	s := NewTurnState(nil, nil)
	if !s.IsCacheable() {
		t.Errorf("NewTurnState should seed cacheable=true, got IsCacheable=false")
	}
}

// TestIsCacheable_DeckReadFlips: calling Deck() flips the bit even when only inspecting
// the slice — a card that reads deck contents binds the chain output to hidden shuffle
// order regardless of whether it modifies the deck.
func TestIsCacheable_DeckReadFlips(t *testing.T) {
	s := NewTurnState([]Card{testutils.NewStubCard("x")}, nil)
	_ = s.Deck()
	if s.IsCacheable() {
		t.Error("Deck() read should flip IsCacheable to false")
	}
}

// TestIsCacheable_GraveyardReadFlips: same as Deck — Graveyard() reads prior-turn graveyard
// contents, so the read alone makes the chain output uncacheable.
func TestIsCacheable_GraveyardReadFlips(t *testing.T) {
	s := NewTurnState(nil, []Card{testutils.NewStubCard("g")})
	_ = s.Graveyard()
	if s.IsCacheable() {
		t.Error("Graveyard() read should flip IsCacheable to false")
	}
}

// TestIsCacheable_PopDeckTopFlips: the verb-based mutator flips the bit + mutates atomically
// so a card that pops the top can't sneak past the cacheable signal.
func TestIsCacheable_PopDeckTopFlips(t *testing.T) {
	a := testutils.NewStubCard("a")
	s := NewTurnState([]Card{a, testutils.NewStubCard("b")}, nil)
	got, ok := s.PopDeckTop()
	if !ok || got != a {
		t.Errorf("PopDeckTop = (%v, %v), want (a, true)", got, ok)
	}
	if s.IsCacheable() {
		t.Error("PopDeckTop should flip IsCacheable to false")
	}
}

// TestIsCacheable_PopDeckTopEmptyFlips: even the empty-deck no-op flips — a card that
// reads "is the deck empty" binds the chain output to that information.
func TestIsCacheable_PopDeckTopEmptyFlips(t *testing.T) {
	s := NewTurnState(nil, nil)
	if got, ok := s.PopDeckTop(); got != nil || ok {
		t.Errorf("PopDeckTop on empty = (%v, %v), want (nil, false)", got, ok)
	}
	if s.IsCacheable() {
		t.Error("PopDeckTop on empty should still flip IsCacheable to false")
	}
}

// TestIsCacheable_PrependToDeckFlips: writes to deck order make the next deck-top reader's
// answer depend on this chain step, so the bit flips as soon as any deck mutation lands.
func TestIsCacheable_PrependToDeckFlips(t *testing.T) {
	s := NewTurnState([]Card{testutils.NewStubCard("x")}, nil)
	added := testutils.NewStubCard("y")
	s.PrependToDeck(added)
	if s.IsCacheable() {
		t.Error("PrependToDeck should flip IsCacheable to false")
	}
	if d := s.Deck(); len(d) != 2 || d[0] != added {
		t.Errorf("Deck = %v, want [y, x] after PrependToDeck", d)
	}
}

// TestIsCacheable_TutorFromDeckFlips: tutoring scans deck contents — even a no-match scan
// flips because the result's "no card found" answer depends on shuffle.
func TestIsCacheable_TutorFromDeckFlips(t *testing.T) {
	target := testutils.NewStubCard("target")
	deck := []Card{testutils.NewStubCard("a"), target, testutils.NewStubCard("b")}
	s := NewTurnState(append([]Card(nil), deck...), nil)
	got, ok := s.TutorFromDeck(func(c Card) int {
		if c == target {
			return 1
		}
		return 0
	})
	if !ok || got != target {
		t.Errorf("TutorFromDeck = (%v, %v), want (target, true)", got, ok)
	}
	if s.IsCacheable() {
		t.Error("TutorFromDeck should flip IsCacheable to false")
	}
	// Tutor removes the matched card.
	if d := s.Deck(); len(d) != 2 {
		t.Errorf("Deck after tutor = %v, want 2 cards remaining", d)
	}
}

// TestIsCacheable_TutorFromDeckNoMatchFlips: the score function returning 0 for every entry
// still flips — the scan ran, the answer depended on the deck contents.
func TestIsCacheable_TutorFromDeckNoMatchFlips(t *testing.T) {
	s := NewTurnState([]Card{testutils.NewStubCard("a"), testutils.NewStubCard("b")}, nil)
	if got, ok := s.TutorFromDeck(func(Card) int { return 0 }); got != nil || ok {
		t.Errorf("TutorFromDeck no-match = (%v, %v), want (nil, false)", got, ok)
	}
	if s.IsCacheable() {
		t.Error("TutorFromDeck (no match) should still flip IsCacheable to false")
	}
}

// TestIsCacheable_BanishFromGraveyardFlips: banishing scans the graveyard — even a pred
// that never matches still flips.
func TestIsCacheable_BanishFromGraveyardFlips(t *testing.T) {
	target := testutils.NewStubCard("target")
	s := NewTurnState(nil, []Card{testutils.NewStubCard("a"), target})
	got, ok := s.BanishFromGraveyard(func(c Card) bool { return c == target })
	if !ok || got != target {
		t.Errorf("BanishFromGraveyard = (%v, %v), want (target, true)", got, ok)
	}
	if s.IsCacheable() {
		t.Error("BanishFromGraveyard should flip IsCacheable to false")
	}
	if len(s.Banish) != 1 || s.Banish[0] != target {
		t.Errorf("Banish = %v, want [target]", s.Banish)
	}
}

// TestIsCacheable_BanishFromGraveyardNoMatchFlips: same scan-flip rule on the no-match path.
func TestIsCacheable_BanishFromGraveyardNoMatchFlips(t *testing.T) {
	s := NewTurnState(nil, []Card{testutils.NewStubCard("a")})
	if got, ok := s.BanishFromGraveyard(func(Card) bool { return false }); got != nil || ok {
		t.Errorf("BanishFromGraveyard no-match = (%v, %v), want (nil, false)", got, ok)
	}
	if s.IsCacheable() {
		t.Error("BanishFromGraveyard (no match) should still flip IsCacheable to false")
	}
}

// Tests that AddToGraveyard flips IsCacheable to false (per the "every public accessor
// flips" convention).
func TestIsCacheable_AddToGraveyardFlips(t *testing.T) {
	s := NewTurnState(nil, nil)
	s.AddToGraveyard(testutils.NewStubCard("x"))
	if s.IsCacheable() {
		t.Error("AddToGraveyard should flip IsCacheable to false")
	}
}

// TestIsCacheable_DrawOneFlipsThroughPopDeckTop: DrawOne is a thin wrapper over PopDeckTop,
// so it inherits the flip — pins that the framework helper doesn't cheat by reading the
// private slice directly.
func TestIsCacheable_DrawOneFlipsThroughPopDeckTop(t *testing.T) {
	s := NewTurnState([]Card{testutils.NewStubCard("x")}, nil)
	s.DrawOne()
	if s.IsCacheable() {
		t.Error("DrawOne should flip IsCacheable to false (inherits via PopDeckTop)")
	}
}

// TestIsCacheable_ClashFlipsThroughDeck: Clash reads s.Deck() to peek the top card; the
// call should propagate the flip.
func TestIsCacheable_ClashFlipsThroughDeck(t *testing.T) {
	s := NewTurnState([]Card{testutils.NewStubCard("top").WithAttack(7)}, nil)
	var won bool
	s.Clash(func() { won = true }, nil)
	if !won {
		t.Errorf("Clash didn't fire win on top atk 7")
	}
	if s.IsCacheable() {
		t.Error("Clash should flip IsCacheable to false (inherits via Deck())")
	}
}

// Tests that NewTurnState's seeding doesn't flip IsCacheable — the constructor writes deck
// and graveyard directly, bypassing the accessor flip.
func TestIsCacheable_NewTurnStateSeedingDoesNotFlip(t *testing.T) {
	s := NewTurnState(
		[]Card{testutils.NewStubCard("x")},
		[]Card{testutils.NewStubCard("y")},
	)
	if !s.IsCacheable() {
		t.Error("NewTurnState seeding should not flip IsCacheable")
	}
}
