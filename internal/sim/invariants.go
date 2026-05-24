package sim

// Per-turn structural invariants. verifyTurnInvariants runs at the end of every
// playOneTurn under the eval (non-replay) path and panics on assertion failure.
// testing.Testing() gates the checks so production anneal binaries compile them out;
// every turn evaluated under `go test ./...` is verified.

import (
	"fmt"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// verifyTurnInvariants checks the structural invariants playOneTurn promises every turn.
// snap is the pre-turn snapshot, preChainHand is the hand after start-of-turn aura
// processing (the actual input to findBest), and summary is the post-turn result.
func verifyTurnInvariants(snap *turnSnapshot, preChainHand []card.Card, summary TurnSummary) {
	if !testing.Testing() || snap == nil || summary.State == nil {
		return
	}
	checkCardConservation(snap, summary)
	checkHandSorted(summary.State.Hand())
	checkPersistentCounts(summary.State)
	checkActionPoints(summary.State)
	checkBestLineMultiset(preChainHand, snap.state.Arsenal(), summary.BestLine)
}

// checkCardConservation asserts the total card count stays constant turn-over-turn:
// deck + graveyard + hand + banished + arsenal + (each card-backed aura / item entry,
// counted once per entry — the source card lives in the aura / item zone until destroyed,
// at which point it goes to graveyard, so the source counts toward "cards in play").
// Token-backed auras / items (SourceCard()==nil) don't carry a card and are excluded.
func checkCardConservation(snap *turnSnapshot, summary TurnSummary) {
	pre := zoneTotal(snap.deck.Size(), snap.state.Graveyard(), snap.hand, snap.state.Banished(), snap.state.Arsenal())
	pre += cardBackedPersistents(snap.state)
	post := zoneTotal(summary.State.Deck().Size(), summary.State.Graveyard(), summary.State.Hand(), summary.State.Banished(), summary.State.Arsenal())
	post += cardBackedPersistents(summary.State)
	if pre != post {
		panic(fmt.Sprintf("turn invariant: card conservation broken (pre=%d post=%d)\n  pre: deck=%d grav=%d hand=%d ban=%d arsenal=%v card-auras+items=%d\n  post: deck=%d grav=%d hand=%d ban=%d arsenal=%v card-auras+items=%d\n  bestLine=%s",
			pre, post,
			snap.deck.Size(), len(snap.state.Graveyard()), len(snap.hand), len(snap.state.Banished()), snap.state.Arsenal() != nil, cardBackedPersistents(snap.state),
			summary.State.Deck().Size(), len(summary.State.Graveyard()), len(summary.State.Hand()), len(summary.State.Banished()), summary.State.Arsenal() != nil, cardBackedPersistents(summary.State),
			FormatBestLine(summary.BestLine)))
	}
}

// zoneTotal sums the per-zone card counts (treating arsenal as 0/1).
func zoneTotal(deckSize int, grav, hand, banished []card.Card, arsenal card.Card) int {
	n := deckSize + len(grav) + len(hand) + len(banished)
	if arsenal != nil {
		n++
	}
	return n
}

// cardBackedPersistents counts aura + item entries whose SourceCard is a real card (not
// a token). Each such entry holds one card in the aura / item zone; destroying the entry
// sends that card to the graveyard, so conservation must count it on both sides.
func cardBackedPersistents(state *gameengine.GameState) int {
	n := 0
	for _, a := range state.Auras() {
		if a.SourceCard() != nil {
			n++
		}
	}
	for _, it := range state.Items() {
		if it.SourceCard() != nil {
			n++
		}
	}
	return n
}

// checkHandSorted asserts the post-turn hand is sorted by Card.ID(). The chain runner
// and eval cache both depend on the sorted invariant — a drift would silently flip cache
// lookups onto unrelated entries on the next turn's findBest call.
func checkHandSorted(hand []card.Card) {
	for i := 1; i < len(hand); i++ {
		if hand[i].ID() < hand[i-1].ID() {
			ids := make([]int, len(hand))
			for j, c := range hand {
				ids[j] = int(c.ID())
			}
			panic(fmt.Sprintf("turn invariant: hand not sorted by ID at index %d: %v", i, ids))
		}
	}
}

// checkPersistentCounts asserts every aura and item entry has non-negative Count.
// Count == 0 is legal (some cards survive a turn with no counters before their next-turn
// destroy fires). A negative count means a Decrement underflowed without a guard.
func checkPersistentCounts(state *gameengine.GameState) {
	for _, a := range state.Auras() {
		if a.Count() < 0 {
			panic(fmt.Sprintf("turn invariant: aura %q has negative Count=%d", a.CardName(), a.Count()))
		}
	}
	for _, it := range state.Items() {
		if it.Count() < 0 {
			panic(fmt.Sprintf("turn invariant: item %q has negative Count=%d", it.CardName(), it.Count()))
		}
	}
}

// checkActionPoints asserts AP didn't go negative. The chain runner's AP check should
// reject any play that would underflow, so a negative value at end of turn means the
// check was bypassed.
func checkActionPoints(state *gameengine.GameState) {
	if state.ActionPoints() < 0 {
		panic(fmt.Sprintf("turn invariant: ActionPoints went negative: %d", state.ActionPoints()))
	}
}

// checkBestLineMultiset asserts BestLine's card multiset equals preChainHand + arsenal-in.
// Partition role-assignment corruption (cards appearing in BestLine that weren't in the
// pre-chain hand) would surface here. preChainHand is the hand AFTER start-of-turn aura
// processing — auras can reveal / draw into hand, so the hand findBest sees may be
// larger than snap.hand.
func checkBestLineMultiset(preChainHand []card.Card, arsenalIn card.Card, line []card.CardAssignment) {
	if len(line) == 0 {
		return // no-chain fallback path; nothing to check
	}
	expectedLen := len(preChainHand)
	if arsenalIn != nil {
		expectedLen++
	}
	if len(line) != expectedLen {
		panic(fmt.Sprintf("turn invariant: BestLine length %d != hand+arsenal %d", len(line), expectedLen))
	}
	want := map[uint16]int{}
	for _, c := range preChainHand {
		want[uint16(c.ID())]++
	}
	if arsenalIn != nil {
		want[uint16(arsenalIn.ID())]++
	}
	got := map[uint16]int{}
	for _, a := range line {
		got[uint16(a.Card.ID())]++
	}
	for id, n := range want {
		if got[id] != n {
			panic(fmt.Sprintf("turn invariant: BestLine multiset mismatch on ID=%d (want %d, got %d)", id, n, got[id]))
		}
	}
	for id, n := range got {
		if want[id] != n {
			panic(fmt.Sprintf("turn invariant: BestLine has unexpected ID=%d (got %d, want %d)", id, n, want[id]))
		}
	}
}
