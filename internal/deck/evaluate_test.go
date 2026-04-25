package deck

import (
	"math/rand"
	"testing"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card/generic"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// TestEvaluate_PerCardStatsPopulated pins per-card attribution: every card that's played or
// pitched increments Plays+Pitches, and TotalContribution sums role-based per-card credit:
// Attack → Card.Attack(), Defend → proportional share of block, Pitch → Card.Pitch(). Held and
// Arsenal cards don't tick the counters (they didn't contribute to this turn's Value). A
// single-printing deck makes the totals easy to assert against the card's printed stats.
func TestEvaluate_PerCardStatsPopulated(t *testing.T) {
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.PerCard == nil {
		t.Fatalf("PerCard should be initialised after Evaluate")
	}
	stat, ok := d.Stats.PerCard[card.ReadTheRunesRed]
	if !ok {
		t.Fatalf("PerCard missing entry for Read the Runes (Red)")
	}
	// Read the Runes Red has no Go again, so the chain plays at most one per turn. With 4 in a
	// 4-card hand, the solver plays one and the rest fall into Held/Arsenal roles which don't
	// tick Plays or Pitches. Counter should be non-zero (at least one Play) but need not sum
	// to 4.
	if got := stat.Plays + stat.Pitches; got == 0 {
		t.Errorf("Plays+Pitches = 0, want at least 1 (the chosen attacker plays once)")
	}
	// Contributions come from the winning chain replay (Play returns + hero triggers) plus
	// role-based shares for pitch/defend. The exact total depends on rider/trigger damage, so
	// assert the weaker property that it's positive and produces a positive Avg.
	if stat.TotalContribution <= 0 {
		t.Errorf("TotalContribution = %v, want >0 (played Read the Runes deals at least Attack+rider)",
			stat.TotalContribution)
	}
	if stat.Avg() <= 0 {
		t.Errorf("Avg() = %v, want >0", stat.Avg())
	}
}

// TestEvaluate_BestTurnStartingRunechantsIsPreHandCarryover pins the contract of
// BestTurn.StartingRunechants: it's the Runechant count carried in from the previous turn when
// the hand was played, so for the first hand of a run it's always 0 — even if the hand itself
// creates runechants that carry out into the next turn.
func TestEvaluate_BestTurnStartingRunechantsIsPreHandCarryover(t *testing.T) {
	// Viserai has Intelligence 4. A 4-card deck gives exactly one hand per run, so the Best
	// record always reflects that first hand — no previous turn ever existed.
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read})

	// Seed doesn't matter (all cards identical), but fix it for determinism.
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	// Sanity: the hand should have left runechants on the table (otherwise the bug couldn't
	// manifest — pre-hand and post-hand counts would both be 0).
	if d.Stats.Best.Summary.Value == 0 {
		t.Fatalf("expected nonzero Value from a hand of Read the Runes; got 0")
	}
	if d.Stats.Best.StartingRunechants != 0 {
		t.Errorf("StartingRunechants = %d, want 0 (first hand of the run has no previous-turn carryover)",
			d.Stats.Best.StartingRunechants)
	}
}

// TestEvaluate_BestTurnSnapshotsDrawnAndLeftoverRunechants pins the BestTurn snapshot's
// completeness: Drawn (mid-turn-drawn cards with their dispositions) and LeftoverRunechants
// must propagate from play.* into Stats.Best.Summary.* so FormatBestTurn's per-card breakdown
// reconciles with the displayed Value and the header's "carryover runechants" count is real.
// Without the snapshot, drawn-attack extension damage and pitch-from-drawn resource land in
// Value but never show up in the printout, and runechants always read 0.
func TestEvaluate_BestTurnSnapshotsDrawnAndLeftoverRunechants(t *testing.T) {
	// Snatch (cost 0, attack 4) fires on-hit DrawOne — its drawn card lands in summary.Drawn.
	// 4 Snatches keeps Viserai's Intelligence-4 hand full of draw-rider cards on the first
	// turn so at least one Snatch attacks and DrawOne fires.
	snatch := cards.Get(card.SnatchRed)
	d := New(hero.Viserai{}, nil, []card.Card{snatch, snatch, snatch, snatch, snatch, snatch, snatch, snatch})
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after Evaluate")
	}
	if len(d.Stats.Best.Summary.Drawn) == 0 {
		t.Errorf("Stats.Best.Summary.Drawn is empty; want >=1 entry from Snatch's on-hit DrawOne (the snapshot in Evaluate isn't copying play.Drawn)")
	}
}

// TestEvaluate_PerCardMarginalCoversEveryHand pins the marginal-stats invariant: for every
// unique card.ID in the deck, PresentHands + AbsentHands equals Stats.Hands. The bucket sums
// are also non-negative and reflect the per-turn hand-value tally so a regression that
// double-counts (or skips) a hand surfaces immediately. A multi-card deck mixes "present
// every turn" vs "present some turns" so both buckets are exercised.
func TestEvaluate_PerCardMarginalCoversEveryHand(t *testing.T) {
	read := cards.Get(card.ReadTheRunesRed)
	snatch := cards.Get(card.SnatchRed)
	// 4 of each so Snatch isn't pinned to a single hand and the absent bucket gets exercised.
	deckCards := []card.Card{read, read, read, read, snatch, snatch, snatch, snatch}
	d := New(hero.Viserai{}, nil, deckCards)
	d.Evaluate(20, 0, rand.New(rand.NewSource(1)))

	if d.Stats.PerCardMarginal == nil {
		t.Fatalf("PerCardMarginal should be initialised after Evaluate")
	}
	for _, id := range []card.ID{card.ReadTheRunesRed, card.SnatchRed} {
		m, ok := d.Stats.PerCardMarginal[id]
		if !ok {
			t.Errorf("PerCardMarginal missing entry for %s", cards.Get(id).Name())
			continue
		}
		if got := m.PresentHands + m.AbsentHands; got != d.Stats.Hands {
			t.Errorf("%s: PresentHands+AbsentHands = %d, want Stats.Hands = %d (every hand must end up in exactly one bucket)",
				cards.Get(id).Name(), got, d.Stats.Hands)
		}
		if m.PresentHands == 0 {
			t.Errorf("%s: PresentHands = 0 — this card should have been in at least one dealt hand across 20 shuffles",
				cards.Get(id).Name())
		}
	}
}

// TestEvaluate_PerCardMarginalAlwaysPresent pins the singleton case: a deck containing only
// one unique card never registers an absent bucket. PresentMean equals the deck's overall
// Mean, and Marginal is 0 (no comparison possible). Single-card decks are a degenerate but
// realistic test fixture.
func TestEvaluate_PerCardMarginalAlwaysPresent(t *testing.T) {
	read := cards.Get(card.ReadTheRunesRed)
	d := New(hero.Viserai{}, nil, []card.Card{read, read, read, read, read, read, read, read})
	d.Evaluate(5, 0, rand.New(rand.NewSource(1)))

	m := d.Stats.PerCardMarginal[card.ReadTheRunesRed]
	if m.AbsentHands != 0 {
		t.Errorf("AbsentHands = %d, want 0 (single-card deck means card is always present)", m.AbsentHands)
	}
	if m.PresentHands != d.Stats.Hands {
		t.Errorf("PresentHands = %d, want %d (every hand contains the only card in the deck)",
			m.PresentHands, d.Stats.Hands)
	}
	if m.Marginal() != 0 {
		t.Errorf("Marginal() = %v, want 0 (no absent comparison possible)", m.Marginal())
	}
}

// TestEvaluate_HeldCardDefersDrawToNextTurn pins the "up to Intelligence" draw rule plus arsenal
// carryover. Intelligence-1 hero, deck of Toughen Up Blue (DR, cost 2, defense 4): the lone
// card has no legal play (can't pay its 2-cost, can't pitch with nothing on the stack, DRs
// can't Attack). Turn 1 holds then promotes it to Arsenal (empty slot). Turn 2 draws a new DR;
// the arsenal card stays on tie, the new card goes Held, so drawCount = 0 next turn and the
// loop halts at Stats.Hands = 2. Neither turn plays or pitches, so PerCard stays at 0.
func TestEvaluate_HeldCardDefersDrawToNextTurn(t *testing.T) {
	// 40 copies of the DR so we have enough deck to fill many hands if held carryover weren't
	// wired up — the assertion would fail catastrophically (loop or much larger Hands count).
	deckCards := make([]card.Card, 40)
	for i := range deckCards {
		deckCards[i] = generic.ToughenUpBlue{}
	}
	d := New(int1StubHero, nil, deckCards)
	d.Evaluate(1, 0, rand.New(rand.NewSource(1)))

	if d.Stats.Hands != 2 {
		t.Errorf("Stats.Hands = %d, want 2 (turn 1 arsenals the card, turn 2 holds its successor, turn 3 can't draw)", d.Stats.Hands)
	}
	tuStat := d.Stats.PerCard[card.ToughenUpBlue]
	if tuStat.Plays != 0 || tuStat.Pitches != 0 {
		t.Errorf("PerCard[ToughenUpBlue] Plays=%d Pitches=%d, want 0/0 (card was Held/Arsenaled, never played or pitched)",
			tuStat.Plays, tuStat.Pitches)
	}
	// Best captures turn 1 (first hand with a recorded play). That hand's single card got
	// promoted from Held to Arsenal by the post-hoc upgrade.
	if len(d.Stats.Best.Summary.BestLine) == 0 {
		t.Fatalf("expected Best to be populated after at least one hand")
	}
	if d.Stats.Best.Summary.BestLine[0].Role != hand.Arsenal {
		t.Errorf("Best.Play.Roles[0] = %s, want ARSENAL (empty slot on turn 1 → Held promoted)", d.Stats.Best.Summary.BestLine[0].Role)
	}
}

// TestEvaluate_ArsenalPersistsAcrossTurns confirms the arsenal slot threads through Evaluate's
// per-turn loop: a card promoted to Arsenal on one turn becomes arsenalCardIn on the next.
// Intelligence-1 hero, 2-card deck of two Toughen Up Blue. Turn 1 arsenals the drawn TU.
// Turn 2 draws the second TU and against incoming 4 plays the arsenal-in DR, pitching the
// drawn card to fund its 2-cost — Value = 4 (prevents the full attack). Turn 3 re-draws the
// pitched card (returned to deck bottom) and arsenals it again. Loop stops when the deck's
// empty and nothing new can be drawn.
func TestEvaluate_ArsenalPersistsAcrossTurns(t *testing.T) {
	d := New(int1StubHero, nil, []card.Card{generic.ToughenUpBlue{}, generic.ToughenUpBlue{}})
	d.Evaluate(1, 4, rand.New(rand.NewSource(1)))

	// Best captures turn 2 — only turn with Value > 0 (arsenal DR fires).
	if d.Stats.Best.Summary.Value != 4 {
		t.Errorf("Best.Play.Value = %d, want 4 (turn 2 plays arsenal DR, pitches hand DR to pay; prevents 4)", d.Stats.Best.Summary.Value)
	}
	// Turn 1: arsenal the drawn card. Turn 2: play arsenal DR (paid by pitching drawn card).
	// Turn 3: draw the recycled pitched card, arsenal it (deck is then empty). Loop ends.
	if d.Stats.Hands != 3 {
		t.Errorf("Stats.Hands = %d, want 3", d.Stats.Hands)
	}
}

// TestEvaluate_TerminatesAfterTwoCycles pins the infinite-loop guard on Evaluate's per-run loop.
// 40 Toughen Up Blue DRs with Reaping Blade equipped, incoming=0, reaches a steady state after
// turn 1 (pitch one TU, swing Reaping Blade for +3, hold the other 3). From then on every turn
// draws and pitches one card — net deck change zero, hand.Best returns the same TurnSummary.
// Without the cap the loop would spin forever; with it, Stats.Hands halts at 2 × handsPerCycle.
func TestEvaluate_TerminatesAfterTwoCycles(t *testing.T) {
	deckCards := make([]card.Card, 40)
	for i := range deckCards {
		deckCards[i] = generic.ToughenUpBlue{}
	}
	d := New(hero.Viserai{}, []weapon.Weapon{weapon.ReapingBlade{}}, deckCards)
	done := make(chan struct{})
	go func() {
		d.Evaluate(1, 0, rand.New(rand.NewSource(1)))
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("Evaluate did not terminate within 2 seconds — infinite loop regression")
	}
	// Two cycles of a 40-card / 4-hand-size deck is exactly 20 hands.
	handsPerCycle := len(deckCards) / hero.Viserai{}.Intelligence()
	maxHands := 2 * handsPerCycle
	if d.Stats.Hands != maxHands {
		t.Errorf("Stats.Hands = %d, want exactly %d (steady-state pitched-pitch loop hits the cap)",
			d.Stats.Hands, maxHands)
	}
}
