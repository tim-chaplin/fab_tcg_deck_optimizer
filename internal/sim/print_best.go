package sim

// Print-time replay of the peak turn. The eval loop captures a bestSnapshot every time a
// new best lands; at print time PrintBestTurn rebuilds the sequenceContext from the
// snapshot, installs a *gameengine.StreamLogger, and runs runDefense + playSequence so
// every card emission goes straight to the writer. No log accumulation anywhere — the
// snapshot carries enough state to reconstruct the run; the logger streams as it happens.

import (
	"fmt"
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// bestSnapshot holds the inputs PrintBestTurn needs to drive a single playSequence call
// matching the eval-time winner. cardsPlayed is the exact order in which cards resolved
// during the winning turn (from the eval-time state's CardsPlayed list).
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

// PrintBestTurn streams the peak turn's printout to w. Writes the header + a brief
// start-of-turn snapshot, then drives ONE turn through the chain runner with a
// StreamLogger so every emission lands inline, then writes the end-of-turn snapshot from
// the post-replay state.
func PrintBestTurn(ev *Evaluator, snap *bestSnapshot, w io.Writer) {
	if snap == nil {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Best turn played (value %d):\n", snap.value)
	writeSnapshotSummary(w, "Start of turn:", snap.master, snap.hand, snap.bestLine)

	postState := replayBest(snap, gameengine.NewStreamLogger(w))

	writeSnapshotSummary(w, "End of turn:", postState, nil, snap.bestLine)
}

// replayBest reconstructs the sequenceContext from snap and runs the winning turn through
// existing chain-runner machinery (runDefense + playSequence) with logger installed on
// every per-perm state. Returns the post-turn *GameState the chain runner left behind so
// the caller can summarise end-of-turn state.
func replayBest(snap *bestSnapshot, logger card.Logger) *gameengine.GameState {
	parts := partitionBestLineForDisplay(snap.bestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	pitched := append(assignmentCards(defensePitches), assignmentCards(attackPitches)...)
	defenders := append(assignmentCards(parts.defenseReactions), assignmentCards(parts.plainBlocks)...)
	held := assignmentCards(parts.held)
	arsenalAtChainStart := arsenalRoleCard(snap.bestLine)
	arsenalCardIn := snap.master.Arsenal()
	arsenalInIdx := matchedCardIndex(snap.cardsPlayed, arsenalCardIn, snap.bestLine, deck.Attack)
	arsenalDefenderIdx := matchedCardIndex(defenders, arsenalCardIn, snap.bestLine, deck.Defend)

	bufs := newAttackBufs(maxInt(len(snap.hand), len(snap.cardsPlayed)), len(snap.weapons), snap.weapons)
	ctx := newSequenceContext(snap.master, snap.weapons, snap.cardsPlayed, defenders, pitched, held, snap.deck, bufs, snap.mp, 0, arsenalInIdx, arsenalAtChainStart)
	ctx.replayLogger = logger
	ctx.attackPitchPerm = assignmentCards(attackPitches)
	ctx.attackPitchVals = pitchValues(attackPitches)
	// resourceBudget starts at 0 — the pitch pool pops perm entries lazily as costs come
	// in, so pre-filling it would short-circuit pay() and leave perm unconsumed (which
	// playSequenceWithMeta then rejects as "pitched without funding anything").
	ctx.resourceBudget = 0

	if len(defenders) > 0 {
		ctx.runDefense(defenders, pitched, snap.deck, snap.mp.IncomingDamage, noBlockBudgetCap, arsenalDefenderIdx)
	}
	ctx.leafState.SetDeck(nil)
	ctx.seedPoolGravBuf(len(snap.cardsPlayed), len(ctx.attackPitchPerm))

	ctx.playSequence(snap.cardsPlayed)
	return ctx.permState
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
