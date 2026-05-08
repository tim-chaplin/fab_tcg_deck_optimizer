package sim_test

import (
	. "github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// damageTrigger returns a StartOfTurn Aura crediting the given damage and destroying
// itself on the first fire. calls is bumped each time the handler runs so tests can assert
// firing counts.
func damageTrigger(self Card, damage int, calls *int) Aura {
	return Aura{
		Trigger: Trigger{
			TriggerType: TriggerStartOfTurn,
			Handler: func(s *TurnState, _ *Trigger, a *Aura) {
				*calls++
				s.AddValue(damage)
				s.DestroyAura(a, true)
			},
		},
		Self:  CardOrTokenType{Card: self},
		Count: 1,
	}
}

// TestProcessAurasAtStartOfTurn_FiresEachQueuedTriggerOnce verifies every queued start-of-turn
// trigger's handler is invoked exactly once per pass, contributions are reported, and a
// trigger whose Count hits zero drops out of survivors.
func TestProcessAurasAtStartOfTurn_FiresEachQueuedTriggerOnce(t *testing.T) {
	aura := testutils.RedAttack{}
	var callsA, callsB int
	queue := []Aura{damageTrigger(aura, 2, &callsA), damageTrigger(aura, 3, &callsB)}
	survivors, contribs, total, _, _ := ProcessAurasAtStartOfTurn(queue, nil)
	if total != 5 {
		t.Errorf("total = %d, want 5 (2+3)", total)
	}
	if len(contribs) != 2 || contribs[0].Damage != 2 || contribs[1].Damage != 3 {
		t.Errorf("contribs = %+v, want [{_, 2}, {_, 3}]", contribs)
	}
	if len(survivors) != 0 {
		t.Errorf("survivors = %+v, want empty (both triggers had Count=1)", survivors)
	}
	if callsA != 1 || callsB != 1 {
		t.Errorf("handler call counts = (%d, %d), want (1, 1)", callsA, callsB)
	}
}

// TestProcessAurasAtStartOfTurn_EmptyQueue short-circuits: no contribs, no alloc, zero total.
func TestProcessAurasAtStartOfTurn_EmptyQueue(t *testing.T) {
	survivors, contribs, total, _, _ := ProcessAurasAtStartOfTurn(nil, nil)
	if total != 0 {
		t.Errorf("total = %d, want 0", total)
	}
	if contribs != nil || len(survivors) != 0 {
		t.Errorf("non-empty outputs on empty input: contribs=%v survivors=%v",
			contribs, survivors)
	}
}

// TestProcessAurasAtStartOfTurn_GraveyardsExhaustedAura: a handler that calls DestroyAura
// lands Self in the turn-state graveyard before subsequent handlers run, so a follow-up
// aura with a graveyard-banish rider sees it.
func TestProcessAurasAtStartOfTurn_GraveyardsExhaustedAura(t *testing.T) {
	aura := testutils.RedAttack{}
	var seen []Card
	// Second trigger's handler records what's currently in the graveyard so we can check the
	// first trigger's destroy happened BEFORE the second fires.
	watcher := Aura{
		Trigger: Trigger{
			TriggerType: TriggerStartOfTurn,
			Handler: func(s *TurnState, _ *Trigger, _ *Aura) {
				seen = append([]Card(nil), s.Graveyard()...)
			},
		},
		Self:  CardOrTokenType{Card: testutils.YellowAttack{}},
		Count: 1,
	}
	_, _, _, _, _ = ProcessAurasAtStartOfTurn([]Aura{
		{
			Trigger: Trigger{TriggerType: TriggerStartOfTurn, Handler: func(s *TurnState, _ *Trigger, a *Aura) {
				s.DestroyAura(a, true)
			}},
			Self:  CardOrTokenType{Card: aura},
			Count: 1,
		},
		watcher,
	}, nil)
	if len(seen) != 1 || seen[0] != aura {
		t.Errorf("second handler saw Graveyard = %v, want [%v] (first trigger's Self graveyarded first)",
			seen, aura)
	}
}

// TestProcessAurasAtStartOfTurn_IgnoresPonder: Ponder fires at end-of-turn (not
// start), so a carried-in Ponder is left untouched by the start-of-turn pass — no
// other framework piece should accidentally drain it.
func TestProcessAurasAtStartOfTurn_IgnoresPonder(t *testing.T) {
	survivors, _, _, _, _ := ProcessAurasAtStartOfTurn([]Aura{NewPonderAura(1)}, nil)
	if len(survivors) != 1 || survivors[0].Self.TokenType != TokenTypePonder {
		t.Errorf("survivors = %+v, want one Ponder aura intact", survivors)
	}
}

// TestFireEndOfTurn_PonderPopsDeckTopIntoHand: the end-of-turn fire pops one card
// per Ponder count, appends each to s.hand, and removes the aura. The drawn card lets
// the post-hoc arsenal-promotion step fill an empty arsenal.
func TestFireEndOfTurn_PonderPopsDeckTopIntoHand(t *testing.T) {
	a, b, c := testutils.NewStubCard("a"), testutils.NewStubCard("b"), testutils.NewStubCard("c")
	s := NewTurnState([]Card{a, b, c}, nil)
	s.Auras = append(s.Auras, NewPonderAura(2))

	FireEndOfTurn(s)

	h := s.Hand()
	if len(h) != 2 || h[0] != a || h[1] != b {
		t.Errorf("Hand = %v, want [a, b]", h)
	}
	if got := len(s.Deck()); got != 1 {
		t.Errorf("len(Deck) = %d, want 1 (two cards popped)", got)
	}
	if len(s.Auras) != 0 {
		t.Errorf("Auras = %+v, want empty (Ponder destroyed itself)", s.Auras)
	}
}

// TestFireEndOfTurn_PonderEmptyDeckIsNoOp: the handler skips the draw step
// silently when the deck is empty — destruction still happens.
func TestFireEndOfTurn_PonderEmptyDeckIsNoOp(t *testing.T) {
	s := NewTurnState(nil, nil)
	s.Auras = append(s.Auras, NewPonderAura(1))

	FireEndOfTurn(s)

	if h := s.Hand(); len(h) != 0 {
		t.Errorf("Hand = %v, want empty (no deck to draw from)", h)
	}
	if len(s.Auras) != 0 {
		t.Errorf("Auras = %+v, want empty (Ponder still destroys with empty deck)", s.Auras)
	}
}

// Tests that Sigil of Fyendal plays turn 1 and its start-of-turn trigger fires on turn 2 —
// crediting 1 damage-equivalent and landing the sigil in graveyard.
func TestEvalTwoTurns_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := cards.SigilOfFyendalBlue{}
	deckCards := []Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, []Card{sigil})

	sigilPlayed := false
	for _, a := range got.Turn1.BestLine {
		if a.Card.ID() == ids.SigilOfFyendalBlue && a.Role == Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", got.Turn1.BestLine)
	}
	if len(got.Turn2.TriggersFromLastTurn) != 1 ||
		got.Turn2.TriggersFromLastTurn[0].Card.ID() != ids.SigilOfFyendalBlue ||
		got.Turn2.TriggersFromLastTurn[0].Damage != 1 {
		t.Errorf("Turn2.TriggersFromLastTurn = %+v, want one Fyendal entry with Damage=1",
			got.Turn2.TriggersFromLastTurn)
	}
	foundInGy := false
	for _, c := range got.Turn2.Graveyard {
		if c.ID() == ids.SigilOfFyendalBlue {
			foundInGy = true
			break
		}
	}
	if !foundInGy {
		t.Errorf("Turn2.Graveyard missing Sigil of Fyendal (Count hit zero after firing); got %v",
			got.Turn2.Graveyard)
	}
}

// Tests that Sigil of the Arknight's handler reveals an attack action into ts.Revealed.
func TestProcessAurasAtStartOfTurn_RevealsAttackActionIntoHand(t *testing.T) {
	var play TurnState
	(cards.SigilOfTheArknightBlue{}).Play(&play, &CardState{Card: cards.SigilOfTheArknightBlue{}})
	slash := cards.AetherSlashRed{}
	_, contribs, total, revealed, _ := ProcessAurasAtStartOfTurn(play.Auras, []Card{slash})
	if total != 0 {
		t.Errorf("total = %d, want 0 (reveal contributes via hand, not damage)", total)
	}
	if len(revealed) != 1 || revealed[0] != slash {
		t.Errorf("revealed = %v, want [%v]", revealed, slash)
	}
	if len(contribs) != 1 || contribs[0].Card.ID() != ids.SigilOfTheArknightBlue {
		t.Errorf("contribs = %+v, want one entry for the sigil", contribs)
	}
}

// Tests that the TriggerContribution carries the revealed card so the printout can attribute
// "drew X into hand" to the specific aura.
func TestProcessAurasAtStartOfTurn_AttributesRevealedToContribution(t *testing.T) {
	var play TurnState
	(cards.SigilOfTheArknightBlue{}).Play(&play, &CardState{Card: cards.SigilOfTheArknightBlue{}})
	slash := cards.AetherSlashRed{}
	_, contribs, _, _, _ := ProcessAurasAtStartOfTurn(play.Auras, []Card{slash})
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	if contribs[0].Revealed == nil || contribs[0].Revealed.ID() != ids.AetherSlashRed {
		t.Errorf("contribs[0].Revealed = %v, want Aether Slash [R]", contribs[0].Revealed)
	}
}

// TestProcessAurasAtStartOfTurn_CascadingReveals: two Arknight sigil triggers in a row each
// reveal the current top, so the second sees the NEW top after the first pops its card.
func TestProcessAurasAtStartOfTurn_CascadingReveals(t *testing.T) {
	var play TurnState
	(cards.SigilOfTheArknightBlue{}).Play(&play, &CardState{Card: cards.SigilOfTheArknightBlue{}})
	(cards.SigilOfTheArknightBlue{}).Play(&play, &CardState{Card: cards.SigilOfTheArknightBlue{}})
	first := cards.AetherSlashRed{}
	second := cards.ConsumingVolitionRed{}
	_, _, _, revealed, _ := ProcessAurasAtStartOfTurn(play.Auras, []Card{first, second})
	if len(revealed) != 2 {
		t.Fatalf("len(revealed) = %d, want 2 (two cascading reveals)", len(revealed))
	}
	if revealed[0] != first || revealed[1] != second {
		t.Errorf("revealed = %v, want [%v, %v]", revealed, first, second)
	}
}

// TestProcessAurasAtStartOfTurn_NonAttackActionTopSkipsReveal: the sigil handler peeks a
// non-attack top → no reveal. The top stays on the deck in the real game.
func TestProcessAurasAtStartOfTurn_NonAttackActionTopSkipsReveal(t *testing.T) {
	var play TurnState
	sigil := cards.SigilOfTheArknightBlue{}
	sigil.Play(&play, &CardState{Card: sigil})
	// Sigil itself is an Aura (non-attack action) — use it as a convenient non-attack top.
	_, _, total, revealed, _ := ProcessAurasAtStartOfTurn(play.Auras, []Card{sigil})
	if total != 0 {
		t.Errorf("total = %d, want 0 (non-attack top, no credit)", total)
	}
	if revealed != nil {
		t.Errorf("revealed = %v, want nil (non-attack tops aren't moved)", revealed)
	}
}

// TestProcessAurasAtStartOfTurn_SigilHitAuthorsLogText: Sigil's handler authors a
// "drew X into hand" Text on its TriggerContribution, captured from the trigger's
// TurnState.Log so the format layer can render the line verbatim.
func TestProcessAurasAtStartOfTurn_SigilHitAuthorsLogText(t *testing.T) {
	var play TurnState
	sigil := cards.SigilOfTheArknightBlue{}
	sigil.Play(&play, &CardState{Card: sigil})
	_, contribs, _, _, _ := ProcessAurasAtStartOfTurn(play.Auras, []Card{cards.AetherSlashRed{}})
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	want := "Sigil of the Arknight [B] drew Aether Slash [R] into hand"
	if contribs[0].Text != want {
		t.Errorf("contribs[0].Text = %q, want %q", contribs[0].Text, want)
	}
}

// TestProcessAurasAtStartOfTurn_SigilWhiffStillLogs: the whiff path (top is a non-attack
// action) authors a "revealed X but didn't draw it" Text so the printout names the card
// the player saw on top of the deck, not just the hits.
func TestProcessAurasAtStartOfTurn_SigilWhiffStillLogs(t *testing.T) {
	var play TurnState
	sigil := cards.SigilOfTheArknightBlue{}
	sigil.Play(&play, &CardState{Card: sigil})
	// Sigil itself is a non-attack action — convenient whiff top.
	_, contribs, _, _, _ := ProcessAurasAtStartOfTurn(play.Auras, []Card{sigil})
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	want := "Sigil of the Arknight [B] revealed Sigil of the Arknight [B] but didn't draw it"
	if contribs[0].Text != want {
		t.Errorf("contribs[0].Text = %q, want %q", contribs[0].Text, want)
	}
}

// Tests end-to-end Sigil of the Arknight: plays turn 1, on turn 2 reveals an attack action
// off the deck top into hand alongside the normal refill.
func TestEvalTwoTurns_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := cards.SigilOfTheArknightBlue{}
	reveal := cards.AetherSlashRed{}
	// Deck layout: positions 0..3 are turn 2's normal refill (Blues), position 4 is the reveal
	// target at the post-draw top, positions 5+ are unused filler.
	deckCards := []Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, []Card{sigil})

	sigilPlayed := false
	for _, a := range got.Turn1.BestLine {
		if a.Card.ID() == ids.SigilOfTheArknightBlue && a.Role == Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", got.Turn1.BestLine)
	}
	// Turn 2 surfaces the reveal as the Revealed field on TriggersFromLastTurn.
	if len(got.Turn2.TriggersFromLastTurn) != 1 ||
		got.Turn2.TriggersFromLastTurn[0].Card.ID() != ids.SigilOfTheArknightBlue ||
		got.Turn2.TriggersFromLastTurn[0].Revealed == nil ||
		got.Turn2.TriggersFromLastTurn[0].Revealed.ID() != (cards.AetherSlashRed{}).ID() {
		t.Errorf("Turn2.TriggersFromLastTurn = %+v, want one Sigil entry revealing Aether Slash",
			got.Turn2.TriggersFromLastTurn)
	}
	if len(got.Turn1.StartOfNextTurnHand) != 4 {
		t.Errorf("Turn1.StartOfNextTurnHand size = %d, want 4 (pre-reveal deal)", len(got.Turn1.StartOfNextTurnHand))
	}
	foundInGy := false
	for _, c := range got.Turn2.Graveyard {
		if c.ID() == ids.SigilOfTheArknightBlue {
			foundInGy = true
			break
		}
	}
	if !foundInGy {
		t.Errorf("Turn2.Graveyard missing Sigil of the Arknight; got %v", got.Turn2.Graveyard)
	}
}

// Tests that Blessing of Occult queues a 3-rune start-of-turn trigger that fires on turn 2.
func TestEvalTwoTurns_BlessingOfOccultCreatesRunesAtStartOfNextTurn(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, []Card{blessing, pitch})

	if got.Turn1.Value != 0 {
		t.Errorf("Turn1.Value = %d, want 0 (Blessing's rune credit is deferred)", got.Turn1.Value)
	}
	blessingPlayed := false
	for _, a := range got.Turn1.BestLine {
		if a.Card.ID() == ids.BlessingOfOccultRed && a.Role == Attack {
			blessingPlayed = true
			break
		}
	}
	if !blessingPlayed {
		t.Errorf("turn 1 BestLine didn't play Blessing as Role=Attack: %+v", got.Turn1.BestLine)
	}
	if len(got.Turn2.TriggersFromLastTurn) != 1 ||
		got.Turn2.TriggersFromLastTurn[0].Card.ID() != ids.BlessingOfOccultRed ||
		got.Turn2.TriggersFromLastTurn[0].Damage != 3 {
		t.Errorf("Turn2.TriggersFromLastTurn = %+v, want one Blessing entry with Damage=3",
			got.Turn2.TriggersFromLastTurn)
	}
	foundInGy := false
	for _, c := range got.Turn2.Graveyard {
		if c.ID() == ids.BlessingOfOccultRed {
			foundInGy = true
			break
		}
	}
	if !foundInGy {
		t.Errorf("Turn2.Graveyard missing Blessing [R]; got %v", got.Turn2.Graveyard)
	}
}

// Tests that Evaluate surfaces a start-of-turn trigger as a TriggersFromLastTurn entry on
// the best-scoring hand and lists the source aura under StartOfTurnAuras.
func TestEvaluate_TriggersFromLastTurnSurfacesInBest(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	slash := cards.AetherSlashRed{}
	deckCards := make([]Card, 0, 20)
	for i := 0; i < 8; i++ {
		deckCards = append(deckCards, blessing)
	}
	for i := 0; i < 6; i++ {
		deckCards = append(deckCards, slash)
	}
	for i := 0; i < 6; i++ {
		deckCards = append(deckCards, testutils.BlueAttack{})
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	rng := rand.New(rand.NewSource(42))
	d.Evaluate(20, Matchup{}, rng)

	if len(d.Stats.Best.Summary.TriggersFromLastTurn) == 0 {
		t.Errorf("Stats.Best.Summary.TriggersFromLastTurn is empty; Best.Value=%d",
			d.Stats.Best.Summary.Value)
	}
	// The best-turn snapshot must also list Blessing under StartOfTurnAuras — Blessing
	// registers a carryover Aura on the turn it's played, and the best-scoring turn
	// is the one where that trigger fires, so the aura has to be in play at the top.
	foundBlessing := false
	for _, a := range d.Stats.Best.Summary.StartOfTurnAuras {
		if a.ID() == ids.BlessingOfOccultRed {
			foundBlessing = true
			break
		}
	}
	if !foundBlessing {
		t.Errorf("Stats.Best.Summary.StartOfTurnAuras missing Blessing; got %+v",
			d.Stats.Best.Summary.StartOfTurnAuras)
	}
}

// Tests that the OncePerTurn FiredThisTurn flag is cleared at every turn boundary so the
// trigger can fire again next turn (no Count tick during re-arm).
func TestProcessAurasAtStartOfTurn_ReArmsOncePerTurnGate(t *testing.T) {
	aura := testutils.RedAttack{}
	exhausted := Aura{
		Trigger: Trigger{
			TriggerType: TriggerAttackAction,
			Handler:     func(*TurnState, *Trigger, *Aura) {},
		},
		Self:          CardOrTokenType{Card: aura},
		Count:         2,
		OncePerTurn:   true,
		FiredThisTurn: true,
	}
	survivors, _, _, _, _ := ProcessAurasAtStartOfTurn([]Aura{exhausted}, nil)
	if len(survivors) != 1 {
		t.Fatalf("survivors len = %d, want 1 (AttackAction trigger passes through)", len(survivors))
	}
	if survivors[0].FiredThisTurn {
		t.Errorf("FiredThisTurn = true, want false (turn-boundary reset)")
	}
	if survivors[0].Count != 2 {
		t.Errorf("Count = %d, want 2 (only re-arm; don't tick)", survivors[0].Count)
	}
}

// Tests that Malefic's OncePerTurn AttackAction trigger fires once per chain and survives
// into next turn with Count decremented.
func TestEvalTwoTurns_MaleficIncantationOncePerTurnLimitsToOneRune(t *testing.T) {
	malefic := cards.MaleficIncantationRed{}
	hocus := cards.HocusPocusRed{}
	// Filler deck so turn 2 can be dealt — content doesn't matter for what we assert.
	deckCards := []Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, []Card{malefic, hocus})

	maleficPlayed, hocusPlayed := false, false
	for _, a := range got.Turn1.BestLine {
		if a.Card.ID() == ids.MaleficIncantationRed && a.Role == Attack {
			maleficPlayed = true
		}
		if a.Card.ID() == ids.HocusPocusRed && a.Role == Attack {
			hocusPlayed = true
		}
	}
	if !maleficPlayed {
		t.Errorf("turn 1 BestLine didn't play Malefic as Role=Attack: %+v", got.Turn1.BestLine)
	}
	if !hocusPlayed {
		t.Errorf("turn 1 BestLine didn't play Hocus Pocus as Role=Attack: %+v", got.Turn1.BestLine)
	}
	if got.Turn1.Value != 6 {
		t.Errorf("Turn1.Value = %d, want 6 (3 Hocus + 1 Hocus rune + 1 Viserai trigger + 1 Malefic trigger)",
			got.Turn1.Value)
	}
	for _, c := range got.Turn2.TriggersFromLastTurn {
		if c.Card.ID() == ids.MaleficIncantationRed {
			t.Errorf("Turn2.TriggersFromLastTurn includes Malefic (%+v); Malefic is AttackAction-typed, not start-of-turn", c)
		}
	}
	for _, c := range got.Turn2.Graveyard {
		if c.ID() == ids.MaleficIncantationRed {
			t.Errorf("Turn2.Graveyard includes Malefic; want empty (Malefic still has Count>0)")
		}
	}
}

// Tests that Runeblood Incantation's start-of-turn trigger fires once per turn and survives
// while Count > 0.
func TestEvalTwoTurns_RunebloodIncantationTicksAcrossTurns(t *testing.T) {
	runeblood := cards.RunebloodIncantationRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	got := d.EvalTwoTurnsForTesting(Matchup{IncomingDamage: 0}, TurnState{}, []Card{runeblood, pitch})

	runebloodPlayed := false
	for _, a := range got.Turn1.BestLine {
		if a.Card.ID() == ids.RunebloodIncantationRed && a.Role == Attack {
			runebloodPlayed = true
			break
		}
	}
	if !runebloodPlayed {
		t.Errorf("turn 1 BestLine didn't play Runeblood as Role=Attack: %+v", got.Turn1.BestLine)
	}
	if got.Turn1.Value != 0 {
		t.Errorf("Turn1.Value = %d, want 0 (every Runeblood rune is deferred to a future fire)",
			got.Turn1.Value)
	}
	if len(got.Turn2.TriggersFromLastTurn) != 1 ||
		got.Turn2.TriggersFromLastTurn[0].Card.ID() != ids.RunebloodIncantationRed ||
		got.Turn2.TriggersFromLastTurn[0].Damage != 1 {
		t.Errorf("Turn2.TriggersFromLastTurn = %+v, want one Runeblood entry with Damage=1 (one tick per turn)",
			got.Turn2.TriggersFromLastTurn)
	}
	for _, c := range got.Turn2.Graveyard {
		if c.ID() == ids.RunebloodIncantationRed {
			t.Errorf("Turn2.Graveyard includes Runeblood; want absent (Red has Count=3, only one tick fired)")
		}
	}
}
