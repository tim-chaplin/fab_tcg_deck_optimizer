package sim

// Print-time replay of the peak turn. The eval loop captures a turnSnapshot every time a
// new best lands; PrintBestTurn re-runs that turn via playOneTurn in replay mode with a
// *gameengine.StreamLogger so every emission streams to the writer — no log accumulation.

import (
	"fmt"
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// turnSnapshot holds the inputs needed to drive a single playSequence call matching the
// eval-time winner. cardsPlayed is the exact resolution order from the winning turn's
// CardsPlayed list. Matchup figures and equipped weapons ride on state.
type turnSnapshot struct {
	state        *gameengine.GameState
	deck         *deck.Deck
	hand         []card.Card
	bestLine     []card.CardAssignment
	cardsPlayed  []card.Card
	swungWeapons []string
	value        int
}

// PrintBestTurn streams the peak turn's printout to w: header + start-of-turn snapshot,
// then one playOneTurn call in replay mode (attack-turn emissions stream inline via a
// StreamLogger), then the end-of-turn snapshot from the post-replay state.
func PrintBestTurn(ev *Evaluator, snap *turnSnapshot, w io.Writer) {
	if snap == nil {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Best turn played (value %d):\n", snap.value)
	writeSnapshotSummary(w, "Start of turn:", snap.state, snap.hand, snap.bestLine)

	snap.state.SetHand(snap.hand)
	summary, _ := playOneTurn(snap.state, snap.deck, ev, snap, gameengine.NewStreamLogger(w))

	writeSnapshotSummary(w, "End of turn:", summary.State, nil, snap.bestLine)
}

// runReplayForTurn drives the attack turn through snapshot's known role assignment + cardsPlayed
// order with logger installed on every per-perm state, skipping Best's partition
// enumeration. Called after start-of-turn auras have fired against snapshot.state / .deck.
func runReplayForTurn(snapshot *turnSnapshot, ev *Evaluator, logger card.Logger) TurnSummary {
	parts := partitionBestLineForDisplay(snapshot.bestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	pitchedCards := append(assignmentCards(defensePitches), assignmentCards(attackPitches)...)
	pitched := wrapPitchedCards(pitchedCards)
	defenders := append(assignmentCards(parts.defenseReactions), assignmentCards(parts.plainBlocks)...)
	held := assignmentCards(parts.held)
	arsenalAtAttackTurnStart := arsenalRoleCard(snapshot.bestLine)
	arsenalCardIn := snapshot.state.Arsenal()
	arsenalInIdx := matchedCardIndex(snapshot.cardsPlayed, arsenalCardIn, snapshot.bestLine, card.Attack)
	arsenalDefenderIdx := matchedCardIndex(defenders, arsenalCardIn, snapshot.bestLine, card.Defend)

	weapons := snapshot.state.Weapons()
	bufs := newAttackBufs(len(snapshot.cardsPlayed), len(weapons), weapons, ev.statePool)
	ctx := newSequenceContext(snapshot.state, weapons, snapshot.cardsPlayed, defenders, pitched, held, snapshot.deck, bufs, 0, arsenalInIdx, arsenalAtAttackTurnStart)
	defer ctx.releaseLeafState()
	ctx.replayLogger = logger
	ctx.attackPitchPerm = wrapPitchedCards(assignmentCards(attackPitches))
	ctx.attackPitchVals = pitchValues(attackPitches)
	// resourceBudget starts at 0 — the pitch pool pops perm entries lazily as costs come
	// in, so pre-filling it would short-circuit pay() and leave perm unconsumed (which
	// playSequenceWithMeta then rejects as "pitched without funding anything").
	ctx.resourceBudget = 0

	if len(defenders) > 0 {
		_, _, ctx.handStart = ctx.runDefense(defenders, pitched, held, snapshot.deck, snapshot.state.IncomingDamage(), noBlockBudgetCap, arsenalDefenderIdx, nil)
	}

	ctx.playSequence(snapshot.cardsPlayed)
	return TurnSummary{
		BestLine:       snapshot.bestLine,
		SwungWeapons:   snapshot.swungWeapons,
		Value:          snapshot.value,
		State:          ctx.permState,
		IncomingDamage: snapshot.state.IncomingDamage(),
	}
}

// writeSnapshotSummary emits a one-section block (header + indented "Hand:" / "Arsenal:"
// / "Auras:" / "Items:" lines). state is the GameState the snapshot reads off; bestLine
// resolves arsenal-role tagging. dealtHand overrides state.Hand() when set (start-of-turn
// uses the snapshot's hand because the post-tick master's hand is nil); when nil, we read
// state.Hand() directly (end-of-turn case).
func writeSnapshotSummary(w io.Writer, header string, state *gameengine.GameState, dealtHand []card.Card, bestLine []card.CardAssignment) {
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
	for _, a := range state.Auras() {
		names = append(names, a.CardName())
	}
	return aurasLineFromNames(names, state.RunechantCount(), state.PonderCount())
}

// itemsSummaryLine renders state's items as "Items: ..." using the per-token counts.
func itemsSummaryLine(state *gameengine.GameState) string {
	return itemsLine(state.GoldCount(), state.SilverCount(), state.CopperCount())
}

// assignmentCards extracts the card.Card field from each assignment.
func assignmentCards(as []card.CardAssignment) []card.Card {
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
func arsenalAssignments(line []card.CardAssignment) []card.CardAssignment {
	var out []card.CardAssignment
	for _, a := range line {
		if a.Role == card.Arsenal {
			out = append(out, a)
		}
	}
	return out
}

// arsenalRoleCard returns the bestLine entry tagged Role=Arsenal, or nil when none.
func arsenalRoleCard(line []card.CardAssignment) card.Card {
	for _, a := range line {
		if a.Role == card.Arsenal {
			return a.Card
		}
	}
	return nil
}

// matchedCardIndex returns the index of arsenalCardIn within cards when the bestLine
// tagged the arsenal-in card with role. -1 when no such assignment exists or no match.
func matchedCardIndex(cards []card.Card, arsenalCardIn card.Card, line []card.CardAssignment, role card.Role) int {
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
func pitchValues(as []card.CardAssignment) []int {
	if len(as) == 0 {
		return nil
	}
	out := make([]int, len(as))
	for i, a := range as {
		out[i] = a.Card.Pitch()
	}
	return out
}
