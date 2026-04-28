package deck

import (
	"math/rand"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// damageTrigger returns a StartOfTurn AuraTrigger crediting the given damage and exhausting
// itself on the first fire. calls is bumped each time the handler runs so tests can assert
// firing counts.
func damageTrigger(self card.Card, damage int, calls *int) card.AuraTrigger {
	return card.AuraTrigger{
		Self:  self,
		Type:  card.TriggerStartOfTurn,
		Count: 1,
		Handler: func(*card.TurnState) int {
			*calls++
			return damage
		},
	}
}

// TestProcessTriggersAtStartOfTurn_FiresEachQueuedTriggerOnce verifies every queued start-of-turn
// trigger's handler is invoked exactly once per pass, contributions are reported, and a
// trigger whose Count hits zero drops out of survivors.
func TestProcessTriggersAtStartOfTurn_FiresEachQueuedTriggerOnce(t *testing.T) {
	aura := testutils.RedAttack{}
	var callsA, callsB int
	queue := []card.AuraTrigger{damageTrigger(aura, 2, &callsA), damageTrigger(aura, 3, &callsB)}
	survivors, contribs, total, _, _, _ := processTriggersAtStartOfTurn(queue, nil)
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

// TestProcessTriggersAtStartOfTurn_EmptyQueue short-circuits: no contribs, no allocation, zero total.
func TestProcessTriggersAtStartOfTurn_EmptyQueue(t *testing.T) {
	survivors, contribs, total, runes, _, _ := processTriggersAtStartOfTurn(nil, nil)
	if total != 0 || runes != 0 {
		t.Errorf("total/runes = %d/%d, want 0/0", total, runes)
	}
	if contribs != nil || len(survivors) != 0 {
		t.Errorf("non-empty outputs on empty input: contribs=%v survivors=%v",
			contribs, survivors)
	}
}

// TestProcessTriggersAtStartOfTurn_GraveyardsExhaustedAura: when a trigger's Count hits zero after
// firing, the sim moves Self into the turn-state graveyard so subsequent handlers (e.g. an
// aura with a graveyard-banish rider) see it. Asserts the contract without relying on any
// specific card to model the "look at graveyard" side.
func TestProcessTriggersAtStartOfTurn_GraveyardsExhaustedAura(t *testing.T) {
	aura := testutils.RedAttack{}
	var seen []card.Card
	// Second trigger's handler records what's currently in the graveyard so we can check the
	// first trigger's destroy happened BEFORE the second fires.
	watcher := card.AuraTrigger{
		Self:  testutils.YellowAttack{},
		Type:  card.TriggerStartOfTurn,
		Count: 1,
		Handler: func(s *card.TurnState) int {
			seen = append([]card.Card(nil), s.Graveyard...)
			return 0
		},
	}
	_, _, _, _, _, _ = processTriggersAtStartOfTurn([]card.AuraTrigger{
		{Self: aura, Type: card.TriggerStartOfTurn, Count: 1, Handler: func(*card.TurnState) int { return 0 }},
		watcher,
	}, nil)
	if len(seen) != 1 || seen[0] != aura {
		t.Errorf("second handler saw Graveyard = %v, want [%v] (first trigger's Self graveyarded first)",
			seen, aura)
	}
}

// TestEvalOneTurn_SigilOfFyendalQueuesTrigger: turn 1 starts with Sigil of Fyendal alone in
// hand. The solver plays it (the beatsBest tiebreaker prefers playing trigger-creating
// auras at equal Value over Held → arsenal promotion). Turn 2's start-of-turn pass fires
// the registered trigger: credit 1 damage-equivalent (the 1{h} gain) and graveyard the
// sigil (Count hits zero).
func TestEvalOneTurn_SigilOfFyendalQueuesTrigger(t *testing.T) {
	sigil := cards.SigilOfFyendalBlue{}
	deckCards := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{sigil})

	sigilPlayed := false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == ids.SigilOfFyendalBlue && a.Role == hand.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play Sigil of Fyendal as Role=Attack: %+v", state.PrevTurnBestLine)
	}
	if state.StartOfTurnTriggerDamage != 1 {
		t.Errorf("StartOfTurnTriggerDamage = %d, want 1 (Fyendal's 1{h} gain fires at start of turn 2)",
			state.StartOfTurnTriggerDamage)
	}
	if len(state.StartOfTurnGraveyard) != 1 || state.StartOfTurnGraveyard[0].ID() != ids.SigilOfFyendalBlue {
		t.Errorf("StartOfTurnGraveyard = %v, want [Sigil of Fyendal] (Count hit zero after firing)",
			state.StartOfTurnGraveyard)
	}
}

// TestProcessTriggersAtStartOfTurn_RevealsAttackActionIntoHand: Sigil of the Arknight's handler
// peeks the post-draw deck top, pops it, and appends to ts.Revealed when it's an attack
// action. The helper surfaces ts.Revealed so the deck loop can forward the revealed card
// into the hand.
func TestProcessTriggersAtStartOfTurn_RevealsAttackActionIntoHand(t *testing.T) {
	var play card.TurnState
	(cards.SigilOfTheArknightBlue{}).Play(&play, &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	slash := cards.AetherSlashRed{}
	_, contribs, total, _, revealed, _ := processTriggersAtStartOfTurn(play.AuraTriggers, []card.Card{slash})
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

// TestProcessTriggersAtStartOfTurn_AttributesRevealedToContribution: the per-trigger
// TriggerContribution carries the card the handler appended so FormatBestTurn can render a
// "drew X into hand" line attributed to the specific aura. Without this attribution the
// printout would know a reveal happened but not which aura caused it.
func TestProcessTriggersAtStartOfTurn_AttributesRevealedToContribution(t *testing.T) {
	var play card.TurnState
	(cards.SigilOfTheArknightBlue{}).Play(&play, &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	slash := cards.AetherSlashRed{}
	_, contribs, _, _, _, _ := processTriggersAtStartOfTurn(play.AuraTriggers, []card.Card{slash})
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	if contribs[0].Revealed == nil || contribs[0].Revealed.ID() != ids.AetherSlashRed {
		t.Errorf("contribs[0].Revealed = %v, want Aether Slash [R]", contribs[0].Revealed)
	}
}

// TestProcessTriggersAtStartOfTurn_CascadingReveals: two Arknight sigil triggers in a row each
// reveal the current top, so the second sees the NEW top after the first pops its card.
func TestProcessTriggersAtStartOfTurn_CascadingReveals(t *testing.T) {
	var play card.TurnState
	(cards.SigilOfTheArknightBlue{}).Play(&play, &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	(cards.SigilOfTheArknightBlue{}).Play(&play, &card.CardState{Card: cards.SigilOfTheArknightBlue{}})
	first := cards.AetherSlashRed{}
	second := cards.ConsumingVolitionRed{}
	_, _, _, _, revealed, _ := processTriggersAtStartOfTurn(play.AuraTriggers, []card.Card{first, second})
	if len(revealed) != 2 {
		t.Fatalf("len(revealed) = %d, want 2 (two cascading reveals)", len(revealed))
	}
	if revealed[0] != first || revealed[1] != second {
		t.Errorf("revealed = %v, want [%v, %v]", revealed, first, second)
	}
}

// TestProcessTriggersAtStartOfTurn_NonAttackActionTopSkipsReveal: the sigil handler peeks a
// non-attack top → no reveal. The top stays on the deck in the real game.
func TestProcessTriggersAtStartOfTurn_NonAttackActionTopSkipsReveal(t *testing.T) {
	var play card.TurnState
	sigil := cards.SigilOfTheArknightBlue{}
	sigil.Play(&play, &card.CardState{Card: sigil})
	// Sigil itself is an Aura (non-attack action) — use it as a convenient non-attack top.
	_, _, total, _, revealed, _ := processTriggersAtStartOfTurn(play.AuraTriggers, []card.Card{sigil})
	if total != 0 {
		t.Errorf("total = %d, want 0 (non-attack top, no credit)", total)
	}
	if revealed != nil {
		t.Errorf("revealed = %v, want nil (non-attack tops aren't moved)", revealed)
	}
}

// TestProcessTriggersAtStartOfTurn_SigilHitAuthorsLogText: Sigil's handler authors a
// "drew X into hand" Text on its TriggerContribution, captured from the trigger's
// TurnState.Log so the format layer can render the line verbatim.
func TestProcessTriggersAtStartOfTurn_SigilHitAuthorsLogText(t *testing.T) {
	var play card.TurnState
	sigil := cards.SigilOfTheArknightBlue{}
	sigil.Play(&play, &card.CardState{Card: sigil})
	_, contribs, _, _, _, _ := processTriggersAtStartOfTurn(play.AuraTriggers, []card.Card{cards.AetherSlashRed{}})
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	want := "Sigil of the Arknight [B] drew Aether Slash [R] into hand"
	if contribs[0].Text != want {
		t.Errorf("contribs[0].Text = %q, want %q", contribs[0].Text, want)
	}
}

// TestProcessTriggersAtStartOfTurn_SigilWhiffStillLogs: the whiff path (top is a non-attack
// action) authors a "revealed X but didn't draw it" Text so the printout names the card
// the player saw on top of the deck, not just the hits.
func TestProcessTriggersAtStartOfTurn_SigilWhiffStillLogs(t *testing.T) {
	var play card.TurnState
	sigil := cards.SigilOfTheArknightBlue{}
	sigil.Play(&play, &card.CardState{Card: sigil})
	// Sigil itself is a non-attack action — convenient whiff top.
	_, contribs, _, _, _, _ := processTriggersAtStartOfTurn(play.AuraTriggers, []card.Card{sigil})
	if len(contribs) != 1 {
		t.Fatalf("contribs = %+v, want one entry", contribs)
	}
	want := "Sigil of the Arknight [B] revealed Sigil of the Arknight [B] but didn't draw it"
	if contribs[0].Text != want {
		t.Errorf("contribs[0].Text = %q, want %q", contribs[0].Text, want)
	}
}

// TestEvalOneTurn_SigilOfTheArknightRevealsIntoHand is the end-to-end 2-turn check: turn 1
// starts with a Sigil of the Arknight as the ONLY card in hand. The solver plays it (the
// beatsBest tiebreaker prefers playing trigger-creating auras at equal Value over Held →
// arsenal promotion, crediting their hidden next-turn payoff). The sigil registers a
// start-of-turn AuraTrigger; on turn 2 the handler peeks the post-draw top (an attack
// action) and moves it into the hand. The returned turn-2 hand should have 5 cards:
// 4 normal refills plus the revealed Aether Slash appended at the tail.
func TestEvalOneTurn_SigilOfTheArknightRevealsIntoHand(t *testing.T) {
	sigil := cards.SigilOfTheArknightBlue{}
	reveal := cards.AetherSlashRed{}
	// Deck layout: positions 0..3 are turn 2's normal refill (Blues), position 4 is the reveal
	// target at the post-draw top, positions 5+ are unused filler.
	deckCards := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{sigil})

	sigilPlayed := false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == ids.SigilOfTheArknightBlue && a.Role == hand.Attack {
			sigilPlayed = true
			break
		}
	}
	if !sigilPlayed {
		t.Errorf("turn 1 BestLine didn't play the sigil as Role=Attack: %+v", state.PrevTurnBestLine)
	}

	// Turn 2: 4 normal draws + 1 revealed = 5 cards. deckCards[0..3] refill turn 2's hand;
	// deckCards[4] is the reveal target appended at the tail.
	wantHand := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		reveal,
	}
	if len(state.Hand) != len(wantHand) {
		t.Fatalf("turn 2 hand size = %d, want %d (4 normal draws + 1 revealed)", len(state.Hand), len(wantHand))
	}
	for i, want := range wantHand {
		if state.Hand[i] != want {
			t.Errorf("turn 2 hand[%d] = %v, want %v", i, state.Hand[i], want)
		}
	}
	if len(state.StartOfTurnGraveyard) != 1 || state.StartOfTurnGraveyard[0].ID() != ids.SigilOfTheArknightBlue {
		t.Errorf("StartOfTurnGraveyard = %v, want [Sigil of the Arknight]", state.StartOfTurnGraveyard)
	}
}

// TestEvalOneTurn_BlessingOfOccultCreatesRunesAtStartOfNextTurn: turn 1's hand has a Red
// Blessing of Occult plus a pitch filler to fund Blessing's 1-cost. Play contributes 0 this
// turn (the 3 Runechants fire at next turn's upkeep via the start-of-turn AuraTrigger), but
// the solver still plays Blessing so the trigger queue picks it up. Turn 2's starting state
// should have 3 Runechants in the carryover.
func TestEvalOneTurn_BlessingOfOccultCreatesRunesAtStartOfNextTurn(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{blessing, pitch})

	if state.PrevTurnValue != 0 {
		t.Errorf("PrevTurnValue = %d, want 0 (Blessing's rune credit is deferred)", state.PrevTurnValue)
	}
	blessingPlayed := false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == ids.BlessingOfOccultRed && a.Role == hand.Attack {
			blessingPlayed = true
			break
		}
	}
	if !blessingPlayed {
		t.Errorf("turn 1 BestLine didn't play Blessing as Role=Attack: %+v", state.PrevTurnBestLine)
	}
	if state.Runechants != 3 {
		t.Errorf("Runechants = %d, want 3 (Blessing's start-of-turn trigger creates 3 tokens)",
			state.Runechants)
	}
	if state.StartOfTurnTriggerDamage != 3 {
		t.Errorf("StartOfTurnTriggerDamage = %d, want 3", state.StartOfTurnTriggerDamage)
	}
	if len(state.StartOfTurnGraveyard) != 1 || state.StartOfTurnGraveyard[0].ID() != ids.BlessingOfOccultRed {
		t.Errorf("StartOfTurnGraveyard = %v, want [Blessing [R]]", state.StartOfTurnGraveyard)
	}
}

// TestEvaluate_TriggersFromLastTurnSurfacesInBest runs a full Evaluate with Red Blessing of
// Occult in the deck and asserts the start-of-turn trigger lands a TriggersFromLastTurn
// entry on at least some hand's TurnSummary. Blessing's +3-rune trigger pads next turn's
// Value directly, so a turn with Blessing queued from the prior turn reliably beats a turn
// without — guaranteeing the best-turn picker selects a trigger-fired hand.
func TestEvaluate_TriggersFromLastTurnSurfacesInBest(t *testing.T) {
	blessing := cards.BlessingOfOccultRed{}
	slash := cards.AetherSlashRed{}
	deckCards := make([]card.Card, 0, 20)
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
	d.Evaluate(20, 0, rng)

	if len(d.Stats.Best.Summary.TriggersFromLastTurn) == 0 {
		t.Errorf("Stats.Best.Summary.TriggersFromLastTurn is empty; Best.Value=%d",
			d.Stats.Best.Summary.Value)
	}
	// The best-turn snapshot must also list Blessing under StartOfTurnAuras — Blessing
	// registers a carryover AuraTrigger on the turn it's played, and the best-scoring turn
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

// TestProcessTriggersAtStartOfTurn_ReArmsOncePerTurnGate: every trigger's FiredThisTurn is
// cleared at every turn boundary regardless of Type, so an AttackAction trigger that
// fired last turn can fire again this turn. Asserts the re-arm contract through the helper
// rather than waiting for the end-to-end multi-turn path to surface a regression.
func TestProcessTriggersAtStartOfTurn_ReArmsOncePerTurnGate(t *testing.T) {
	aura := testutils.RedAttack{}
	exhausted := card.AuraTrigger{
		Self:          aura,
		Type:          card.TriggerAttackAction,
		Count:         2,
		OncePerTurn:   true,
		FiredThisTurn: true,
		Handler:       func(*card.TurnState) int { return 1 },
	}
	survivors, _, _, _, _, _ := processTriggersAtStartOfTurn([]card.AuraTrigger{exhausted}, nil)
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

// TestEvalOneTurn_MaleficIncantationOncePerTurnLimitsToOneRune: turn 1 plays Red Malefic
// (Count=3, OncePerTurn) followed by Red Hocus Pocus (an attack action card). Hocus is the
// only attack action this turn so the trigger fires exactly once — verifying the gate
// doesn't *prevent* the first fire (a separate hand-package test exercises the gate
// closing on a same-turn second attack action). Turn 1 Value breaks down as:
//
//	+3 Hocus Pocus attack (printed power)
//	+1 Hocus Pocus's own Runechant-creation rider
//	+1 Viserai trigger (Hocus is a Runeblade card; Malefic was a prior non-attack action)
//	+1 Malefic AttackAction trigger (creates one Runechant)
//	= 6
//
// Turn 2's start-of-turn pass clears FiredThisTurn so the trigger can fire again next
// turn, but doesn't itself credit damage (Malefic is AttackAction-typed, not StartOfTurn).
// Malefic survives with Count=2.
func TestEvalOneTurn_MaleficIncantationOncePerTurnLimitsToOneRune(t *testing.T) {
	malefic := cards.MaleficIncantationRed{}
	hocus := cards.HocusPocusRed{}
	// Filler deck so turn 2 can be dealt — content doesn't matter for what we assert.
	deckCards := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{malefic, hocus})

	maleficPlayed, hocusPlayed := false, false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == ids.MaleficIncantationRed && a.Role == hand.Attack {
			maleficPlayed = true
		}
		if a.Card.ID() == ids.HocusPocusRed && a.Role == hand.Attack {
			hocusPlayed = true
		}
	}
	if !maleficPlayed {
		t.Errorf("turn 1 BestLine didn't play Malefic as Role=Attack: %+v", state.PrevTurnBestLine)
	}
	if !hocusPlayed {
		t.Errorf("turn 1 BestLine didn't play Hocus Pocus as Role=Attack: %+v", state.PrevTurnBestLine)
	}
	if state.PrevTurnValue != 6 {
		t.Errorf("PrevTurnValue = %d, want 6 (3 Hocus + 1 Hocus rune + 1 Viserai trigger + 1 Malefic trigger)",
			state.PrevTurnValue)
	}
	// Malefic's AttackAction trigger doesn't fire at start of turn — it only ticks on
	// attack actions during the chain. Carry-only at the turn boundary.
	if state.StartOfTurnTriggerDamage != 0 {
		t.Errorf("StartOfTurnTriggerDamage = %d, want 0 (Malefic only fires on attack actions)",
			state.StartOfTurnTriggerDamage)
	}
	if len(state.StartOfTurnGraveyard) != 0 {
		t.Errorf("StartOfTurnGraveyard = %v, want empty (Malefic still has Count>0)",
			state.StartOfTurnGraveyard)
	}
}

// TestEvalOneTurn_RunebloodIncantationTicksAcrossTurns: turn 1 plays Red Runeblood
// Incantation (Count=3 verse counters). Turn 2's start-of-turn pass fires the trigger once
// — credits 1 Runechant, decrements Count to 2, leaves the aura alive. The surviving
// trigger is what carries forward; this test pins the multi-turn fire shape end-to-end at
// the deck-loop boundary.
func TestEvalOneTurn_RunebloodIncantationTicksAcrossTurns(t *testing.T) {
	runeblood := cards.RunebloodIncantationRed{}
	pitch := testutils.PitchOneDR{}
	deckCards := []card.Card{
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
		testutils.BlueAttack{},
	}
	d := New(heroes.Viserai{}, nil, deckCards)
	state := d.EvalOneTurnForTesting(0, nil, []card.Card{runeblood, pitch})

	runebloodPlayed := false
	for _, a := range state.PrevTurnBestLine {
		if a.Card.ID() == ids.RunebloodIncantationRed && a.Role == hand.Attack {
			runebloodPlayed = true
			break
		}
	}
	if !runebloodPlayed {
		t.Errorf("turn 1 BestLine didn't play Runeblood as Role=Attack: %+v", state.PrevTurnBestLine)
	}
	if state.PrevTurnValue != 0 {
		t.Errorf("PrevTurnValue = %d, want 0 (every Runeblood rune is deferred to a future fire)",
			state.PrevTurnValue)
	}
	if state.StartOfTurnTriggerDamage != 1 {
		t.Errorf("StartOfTurnTriggerDamage = %d, want 1 (one tick per turn)", state.StartOfTurnTriggerDamage)
	}
	if state.Runechants != 1 {
		t.Errorf("Runechants = %d, want 1 (one rune per fire)", state.Runechants)
	}
	if len(state.StartOfTurnGraveyard) != 0 {
		t.Errorf("StartOfTurnGraveyard = %v, want empty (Red has Count=3, only one tick fired)",
			state.StartOfTurnGraveyard)
	}
}
