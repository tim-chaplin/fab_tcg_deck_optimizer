package card

import "testing"

// TestDrawOne_AppendsTopAndAdvancesDeck: DrawOne pops the top of Deck and appends it to Hand,
// preserving draw order for downstream effects.
func TestDrawOne_AppendsTopAndAdvancesDeck(t *testing.T) {
	a, b, c := stubCard{name: "a"}, stubCard{name: "b"}, stubCard{name: "c"}
	s := &TurnState{}
	s.SetDeck([]Card{a, b, c})

	s.DrawOne()
	if got := len(s.Deck()); got != 2 {
		t.Fatalf("Deck len = %d, want 2", got)
	}
	if s.Deck()[0] != b {
		t.Errorf("Deck[0] = %v, want b (top advanced past a)", s.Deck()[0])
	}
	if len(s.Hand) != 1 || s.Hand[0] != a {
		t.Errorf("Hand = %v, want [a]", s.Hand)
	}

	s.DrawOne()
	if len(s.Hand) != 2 || s.Hand[1] != b {
		t.Errorf("Hand after second draw = %v, want [a, b]", s.Hand)
	}
}

// TestDrawOne_EmptyDeckIsNoOp: with an empty deck the helper returns silently; Hand stays
// untouched so callers don't see a spurious zero-value card.
func TestDrawOne_EmptyDeckIsNoOp(t *testing.T) {
	s := &TurnState{}
	s.DrawOne()
	if len(s.Hand) != 0 {
		t.Errorf("Hand = %v, want empty on no-deck draw", s.Hand)
	}
	if s.Deck() != nil {
		t.Errorf("Deck = %v, want nil", s.Deck())
	}
}

// TestIsCacheable_DefaultsTrue: a freshly-constructed TurnState reports IsCacheable=true.
// The inverted-storage trick (uncacheable bool, default false) gives us the right default
// without requiring a constructor or per-permutation reset call.
func TestIsCacheable_DefaultsTrue(t *testing.T) {
	var s TurnState
	if !s.IsCacheable() {
		t.Error("zero-valued TurnState should be cacheable")
	}
}

// TestIsCacheable_SetDeckDoesNotFlip: SetDeck is a write — same input + chain operations
// give the same output deck, so writes don't make the result depend on hidden order.
func TestIsCacheable_SetDeckDoesNotFlip(t *testing.T) {
	var s TurnState
	s.SetDeck([]Card{stubCard{name: "a"}})
	if !s.IsCacheable() {
		t.Error("SetDeck should not flip cacheable")
	}
}

// TestIsCacheable_DeckReadFlips: Deck() reads contents, which makes the chain output
// depend on hidden shuffle state — flips IsCacheable to false.
func TestIsCacheable_DeckReadFlips(t *testing.T) {
	var s TurnState
	s.SetDeck([]Card{stubCard{name: "a"}})
	_ = s.Deck()
	if s.IsCacheable() {
		t.Error("Deck() read should flip IsCacheable to false")
	}
}

// TestIsCacheable_GraveyardReadFlips: Graveyard() scan (e.g. aura-banish) reads prior-turn
// hidden state, so the chain becomes uncacheable.
func TestIsCacheable_GraveyardReadFlips(t *testing.T) {
	var s TurnState
	_ = s.Graveyard()
	if s.IsCacheable() {
		t.Error("Graveyard() read should flip IsCacheable to false")
	}
}

// TestIsCacheable_AddToGraveyardDoesNotFlip: append-only graveyard mutations are
// deterministic from inputs + cards played, so they don't make the chain uncacheable.
// This is the per-played-card path the framework uses for every chain-step finisher.
func TestIsCacheable_AddToGraveyardDoesNotFlip(t *testing.T) {
	var s TurnState
	s.AddToGraveyard(stubCard{name: "a"})
	if !s.IsCacheable() {
		t.Error("AddToGraveyard should not flip cacheable (append-only mutation)")
	}
}

// TestIsCacheable_DrawOneFlips: DrawOne reads the deck top, so it flips. Cards calling
// DrawOne (Snatch, Sun Kiss, Drawn to the Dark Dimension) inherit the flip without
// having to set anything explicitly — that's the structural-enforcement guarantee.
func TestIsCacheable_DrawOneFlips(t *testing.T) {
	var s TurnState
	s.SetDeck([]Card{stubCard{name: "a"}})
	s.DrawOne()
	if s.IsCacheable() {
		t.Error("DrawOne should flip IsCacheable to false (reads deck top)")
	}
}

// TestIsCacheable_ClashValueFlips: ClashValue compares deck-top to opponent's; reading
// deck-top flips cacheable. Test of Strength routes through this helper.
func TestIsCacheable_ClashValueFlips(t *testing.T) {
	var s TurnState
	s.SetDeck([]Card{stubCard{attack: 6}})
	_ = ClashValue(&s, 1)
	if s.IsCacheable() {
		t.Error("ClashValue should flip IsCacheable to false (reads deck top)")
	}
}

// TestIsCacheable_PopDeckTopFlips: PopDeckTop is a card-callable mutator that reads the
// deck top to return it; flips IsCacheable since the popped card's identity depends on
// hidden shuffle state.
func TestIsCacheable_PopDeckTopFlips(t *testing.T) {
	var s TurnState
	s.SetDeck([]Card{stubCard{name: "a"}})
	if _, ok := s.PopDeckTop(); !ok {
		t.Fatal("PopDeckTop returned false on a non-empty deck")
	}
	if s.IsCacheable() {
		t.Error("PopDeckTop should flip IsCacheable to false")
	}
}

// TestIsCacheable_PrependToDeckFlips: PrependToDeck mutates deck-top, so subsequent
// readers see this card on top — flips since the post-state depends on the pre-state's
// hidden order.
func TestIsCacheable_PrependToDeckFlips(t *testing.T) {
	var s TurnState
	s.PrependToDeck(stubCard{name: "a"})
	if s.IsCacheable() {
		t.Error("PrependToDeck should flip IsCacheable to false")
	}
}

// TestIsCacheable_TutorFromDeckFlips: TutorFromDeck scans the whole deck — reading the
// contents clearly flips.
func TestIsCacheable_TutorFromDeckFlips(t *testing.T) {
	var s TurnState
	s.SetDeck([]Card{stubCard{name: "a"}})
	s.TutorFromDeck(func(Card) int { return 1 })
	if s.IsCacheable() {
		t.Error("TutorFromDeck should flip IsCacheable to false")
	}
}

// TestIsCacheable_BanishFromGraveyardFlips: BanishFromGraveyard scans the graveyard.
func TestIsCacheable_BanishFromGraveyardFlips(t *testing.T) {
	var s TurnState
	s.AddToGraveyard(stubCard{name: "a"})
	if !s.IsCacheable() {
		t.Fatal("pre: AddToGraveyard alone should keep cacheable")
	}
	s.BanishFromGraveyard(func(Card) bool { return true })
	if s.IsCacheable() {
		t.Error("BanishFromGraveyard should flip IsCacheable to false")
	}
}

// TestAddAuraTrigger_FlipsAuraCreatedAndAppends: AddAuraTrigger MUST flip AuraCreated (so
// same-turn "if you've played or created an aura" riders see the entry) AND push each
// trigger onto s.AuraTriggers in call order. Pairing both in one method is what stops a
// card from registering a trigger without advertising the aura (or vice versa).
func TestAddAuraTrigger_FlipsAuraCreatedAndAppends(t *testing.T) {
	self := stubCard{name: "self"}
	s := &TurnState{}
	if s.AuraCreated {
		t.Fatal("pre: AuraCreated should be false")
	}
	s.AddAuraTrigger(AuraTrigger{Self: self, Type: TriggerStartOfTurn, Count: 2})
	s.AddAuraTrigger(AuraTrigger{Self: self, Type: TriggerStartOfTurn, Count: 1})
	if !s.AuraCreated {
		t.Error("AuraCreated = false, want true")
	}
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

	playedAura := TurnState{CardsPlayed: []Card{stubCard{types: NewTypeSet(TypeAura)}}}
	if !playedAura.HasPlayedOrCreatedAura() {
		t.Error("played aura card → should be true")
	}
}

// TestRecordValue_ClampsNonPositive: the helper sums positive credits into Value and is a
// no-op for n <= 0. Negative grants (debuffs) and zero (no-effect Plays) must not subtract
// from the running total.
func TestRecordValue_ClampsNonPositive(t *testing.T) {
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
			s.RecordValue(tc.bump)
			if s.Value != tc.want {
				t.Errorf("Value = %d, want %d", s.Value, tc.want)
			}
		})
	}
	// Mixed sequence: positives accumulate, non-positives pass through.
	var s TurnState
	s.RecordValue(2)
	s.RecordValue(-10)
	s.RecordValue(0)
	s.RecordValue(5)
	if s.Value != 7 {
		t.Errorf("after mixed sequence Value = %d, want 7 (2+5; -10/0 clamped)", s.Value)
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
	gy := s.Graveyard()
	if len(gy) != 2 || gy[0] != a || gy[1] != b {
		t.Errorf("Graveyard = %v, want [a, b]", gy)
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
			s.SetDeck([]Card{stubCard{attack: tc.topAtk}})
		}
		if got := ClashValue(&s, tc.bonus); got != tc.want {
			t.Errorf("%s: ClashValue = %d, want %d", tc.name, got, tc.want)
		}
	}
}

// TestRegisterStartOfTurn_AutoLogsWithText: when text is non-empty and the handler returns
// n > 0, the wrapper writes a post-trigger log entry "<DisplayName>: text (+n)" attributed
// to self before returning. This is what processTriggersAtStartOfTurn captures as the
// trigger's contribution Text so the rendered turn log names the effect instead of falling
// back to "START OF ACTION PHASE (+N)".
func TestRegisterStartOfTurn_AutoLogsWithText(t *testing.T) {
	self := stubCard{name: "Test Aura", types: NewTypeSet(TypeAura)}
	var s TurnState
	s.RegisterStartOfTurn(self, 1, "Gained 1 health", func(*TurnState) int { return 1 })
	if len(s.AuraTriggers) != 1 {
		t.Fatalf("AuraTriggers len = %d, want 1", len(s.AuraTriggers))
	}

	var fired TurnState
	got := s.AuraTriggers[0].Handler(&fired)
	if got != 1 {
		t.Errorf("handler return = %d, want 1 (passes through inner return)", got)
	}
	if len(fired.Log) != 1 {
		t.Fatalf("Log len = %d, want 1 entry from auto-log", len(fired.Log))
	}
	if want := "Test Aura: Gained 1 health (+1)"; fired.Log[0].Text != want {
		t.Errorf("auto-log text = %q, want %q", fired.Log[0].Text, want)
	}
}

// TestRegisterStartOfTurn_NoLogOnZero: text is set but the handler returned 0 — wrapper
// skips the log entry so a no-op fire (e.g. Sigil of Silphidae's leave trigger when the
// graveyard has no aura to banish) doesn't emit a misleading "Banished an aura" line.
func TestRegisterStartOfTurn_NoLogOnZero(t *testing.T) {
	self := stubCard{name: "Test Aura", types: NewTypeSet(TypeAura)}
	var s TurnState
	s.RegisterStartOfTurn(self, 1, "Did the thing", func(*TurnState) int { return 0 })

	var fired TurnState
	s.AuraTriggers[0].Handler(&fired)
	if len(fired.Log) != 0 {
		t.Errorf("Log = %v, want empty (handler returned 0)", fired.Log)
	}
}

// TestRegisterStartOfTurn_EmptyTextLeavesHandlerAlone: text == "" means the card authors its
// own log line inside the handler (Sigil of the Arknight's "drew X into hand", Silphidae's
// conditional banish line). The wrapper must not append any extra entries — handler logs are
// the only entries written.
func TestRegisterStartOfTurn_EmptyTextLeavesHandlerAlone(t *testing.T) {
	self := stubCard{name: "Test Aura", types: NewTypeSet(TypeAura)}
	var s TurnState
	s.RegisterStartOfTurn(self, 1, "", func(s *TurnState) int {
		s.AddPostTriggerLogEntry("custom handler text", "Test Aura", 0)
		return 0
	})

	var fired TurnState
	s.AuraTriggers[0].Handler(&fired)
	if len(fired.Log) != 1 {
		t.Fatalf("Log len = %d, want exactly 1 (handler-authored only)", len(fired.Log))
	}
	if fired.Log[0].Text != "custom handler text" {
		t.Errorf("Log[0].Text = %q, want handler's own text", fired.Log[0].Text)
	}
}
