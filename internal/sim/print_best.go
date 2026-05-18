package sim

// Print-time replay of the peak turn. The eval loop captures a bestSnapshot every time a
// new best lands; PrintBestTurn re-runs that turn via playOneTurn in replay mode with a
// *gameengine.StreamLogger so every emission streams to the writer — no log accumulation.

import (
	"fmt"
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// bestSnapshot holds the inputs needed to drive a single playSequence call matching the
// eval-time winner. cardsPlayed is the exact resolution order from the winning turn's
// CardsPlayed list.
type bestSnapshot struct {
	master       *gameengine.GameState
	deck         *deck.Deck
	hand         []card.Card
	weapons      []weapon.Weapon
	mp           Matchup
	bestLine     []deck.CardAssignment
	cardsPlayed  []card.Card
	swungWeapons []string
	value        int
}

// PrintBestTurn streams the peak turn's printout to w: header + start-of-turn snapshot,
// then one playOneTurn call in replay mode (chain emissions stream inline via a
// StreamLogger), then the end-of-turn snapshot from the post-replay state.
func PrintBestTurn(ev *Evaluator, snap *bestSnapshot, w io.Writer) {
	if snap == nil {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Best turn played (value %d):\n", snap.value)
	writeSnapshotSummary(w, "Start of turn:", snap.master, snap.hand, snap.bestLine)

	summary, _, _, _ := playOneTurn(snap.master, snap.hand, snap.deck, snap.mp, snap.weapons, ev, 0, nil, snap, gameengine.NewStreamLogger(w))

	writeSnapshotSummary(w, "End of turn:", summary.State, nil, snap.bestLine)
}

// runReplayForTurn drives the chain through snapshot's known role assignment + cardsPlayed
// order with logger installed on every per-perm state, skipping Best's partition
// enumeration. Called after start-of-turn auras have fired against snapshot.master / .deck.
func runReplayForTurn(snapshot *bestSnapshot, logger card.Logger) TurnSummary {
	parts := partitionBestLineForDisplay(snapshot.bestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	pitched := append(assignmentCards(defensePitches), assignmentCards(attackPitches)...)
	defenders := append(assignmentCards(parts.defenseReactions), assignmentCards(parts.plainBlocks)...)
	held := assignmentCards(parts.held)
	arsenalAtChainStart := arsenalRoleCard(snapshot.bestLine)
	arsenalCardIn := snapshot.master.Arsenal()
	arsenalInIdx := matchedCardIndex(snapshot.cardsPlayed, arsenalCardIn, snapshot.bestLine, deck.Attack)
	arsenalDefenderIdx := matchedCardIndex(defenders, arsenalCardIn, snapshot.bestLine, deck.Defend)

	bufs := newAttackBufs(len(snapshot.cardsPlayed), len(snapshot.weapons), snapshot.weapons)
	ctx := newSequenceContext(snapshot.master, snapshot.weapons, snapshot.cardsPlayed, defenders, pitched, held, snapshot.deck, bufs, snapshot.mp, 0, arsenalInIdx, arsenalAtChainStart)
	ctx.replayLogger = logger
	ctx.attackPitchPerm = assignmentCards(attackPitches)
	ctx.attackPitchVals = pitchValues(attackPitches)
	// resourceBudget starts at 0 — the pitch pool pops perm entries lazily as costs come
	// in, so pre-filling it would short-circuit pay() and leave perm unconsumed (which
	// playSequenceWithMeta then rejects as "pitched without funding anything").
	ctx.resourceBudget = 0

	if len(defenders) > 0 {
		ctx.runDefense(defenders, pitched, snapshot.deck, snapshot.mp.IncomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	}
	ctx.leafState.SetDeck(nil)
	ctx.seedPoolGravBuf(len(snapshot.cardsPlayed), len(ctx.attackPitchPerm))

	ctx.playSequence(snapshot.cardsPlayed)
	return TurnSummary{
		BestLine:       snapshot.bestLine,
		SwungWeapons:   snapshot.swungWeapons,
		Value:          snapshot.value,
		State:          ctx.permState,
		IncomingDamage: snapshot.mp.IncomingDamage,
	}
}

// writeSnapshotSummary emits a one-section block (header + indented "Hand:" / "Arsenal:"
// / "Auras:" / "Items:" lines). state is the GameState the snapshot reads off; bestLine
// resolves arsenal-role tagging. dealtHand overrides state.Hand() when set (start-of-turn
// uses the snapshot's hand because the post-tick master's hand is nil); when nil, we read
// state.Hand() directly (end-of-turn case).
func writeSnapshotSummary(w io.Writer, header string, state *gameengine.GameState, dealtHand []card.Card, bestLine []deck.CardAssignment) {
	if state == nil {
		return
	}
	var lines []string
	hand := dealtHand
	if hand == nil {
		hand = state.Hand()
	}
	if l := startingHandLine(hand); l != "" {
		lines = append(lines, l)
	}
	if dealtHand != nil {
		if l := startingArsenalLine(bestLine); l != "" {
			lines = append(lines, l)
		}
	} else {
		if l := endingArsenalLine(arsenalAssignments(bestLine)); l != "" {
			lines = append(lines, l)
		}
	}
	if l := aurasSummaryLine(state); l != "" {
		lines = append(lines, l)
	}
	if l := itemsSummaryLine(state); l != "" {
		lines = append(lines, l)
	}
	if len(lines) == 0 {
		return
	}
	fmt.Fprintf(w, "  %s\n", header)
	for _, l := range lines {
		fmt.Fprintf(w, "    %s\n", l)
	}
}

// aurasSummaryLine renders state's surviving auras as "Auras: Name, ..., N Runechants,
// M Ponders". Returns "" when no auras / tokens are in play.
func aurasSummaryLine(state *gameengine.GameState) string {
	var names []string
	runechants := 0
	ponders := 0
	for _, a := range state.Auras() {
		name := a.CardName()
		switch name {
		case "Runechant":
			runechants += a.Count()
		case "Ponder":
			ponders += a.Count()
		default:
			names = append(names, name)
		}
	}
	return aurasLineFromNames(names, runechants, ponders)
}

// itemsSummaryLine renders state's items as "Items: ..." using the per-token counts.
func itemsSummaryLine(state *gameengine.GameState) string {
	gold, silver, copper := 0, 0, 0
	for _, it := range state.Items() {
		switch it.CardName() {
		case "Gold":
			gold += it.Count()
		case "Silver":
			silver += it.Count()
		case "Copper":
			copper += it.Count()
		}
	}
	return itemsLine(gold, silver, copper)
}

// assignmentCards extracts the card.Card field from each assignment.
func assignmentCards(as []deck.CardAssignment) []card.Card {
	if len(as) == 0 {
		return nil
	}
	out := make([]card.Card, len(as))
	for i, a := range as {
		out[i] = a.Card
	}
	return out
}

// arsenalAssignments returns the bestLine entries with Role=Arsenal so endingArsenalLine
// can tag each with "(stayed)" / "(new)".
func arsenalAssignments(line []deck.CardAssignment) []deck.CardAssignment {
	var out []deck.CardAssignment
	for _, a := range line {
		if a.Role == deck.Arsenal {
			out = append(out, a)
		}
	}
	return out
}

// arsenalRoleCard returns the bestLine entry tagged Role=Arsenal, or nil when none.
func arsenalRoleCard(line []deck.CardAssignment) card.Card {
	for _, a := range line {
		if a.Role == deck.Arsenal {
			return a.Card
		}
	}
	return nil
}

// matchedCardIndex returns the index of arsenalCardIn within cards when the bestLine
// tagged the arsenal-in card with role. -1 when no such assignment exists or no match.
func matchedCardIndex(cards []card.Card, arsenalCardIn card.Card, line []deck.CardAssignment, role deck.Role) int {
	if arsenalCardIn == nil {
		return -1
	}
	matched := false
	for _, a := range line {
		if a.FromArsenal && a.Role == role {
			matched = true
			break
		}
	}
	if !matched {
		return -1
	}
	for i, c := range cards {
		if c == arsenalCardIn {
			return i
		}
	}
	return -1
}

// pitchValues returns the Pitch() values of each assignment.
func pitchValues(as []deck.CardAssignment) []int {
	if len(as) == 0 {
		return nil
	}
	out := make([]int, len(as))
	for i, a := range as {
		out[i] = a.Card.Pitch()
	}
	return out
}

// maxInt returns the larger of two ints.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
