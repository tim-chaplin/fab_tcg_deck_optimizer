package sim

// Print-time replay of the peak turn. PrintBestTurn re-runs Best against a snapshot to
// recover the winning chain order from State.CardsPlayed(), then streams the four-section
// printout (start of turn, my turn, opponent's turn, end of turn) to an io.Writer.

import (
	"fmt"
	"io"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/item"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// bestSnapshot holds the inputs required to replay the peak turn at print time.
type bestSnapshot struct {
	master   *gameengine.GameState
	deck     *deck.Deck
	hand     []card.Card
	weapons  []weapon.Weapon
	mp       Matchup
	bestLine []deck.CardAssignment
	value    int
}

// PrintBestTurn replays the peak turn and streams its printout to w — header
// "Best turn played (value N):" plus the four-section body.
func PrintBestTurn(ev *Evaluator, snap *bestSnapshot, w io.Writer) {
	if snap == nil {
		return
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Best turn played (value %d):\n", snap.value)

	masterCopy := snap.master.CopyPersistentState()
	deckCopy := snap.deck.Copy()
	handCopy := append([]card.Card(nil), snap.hand...)
	summary := ev.Best(snap.weapons, handCopy, snap.mp, deckCopy, masterCopy)

	parts := partitionBestLineForDisplay(snap.bestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	startOfTurn := buildStartOfTurnLines(snap.hand, snap.bestLine, snap.master)
	myTurn := buildMyTurnLines(attackPitches, summary.State.CardsPlayed())
	oppTurn := buildOpponentTurnLines(defensePitches, parts)
	endOfTurn := buildEndOfTurnLines(parts.arsenal, summary.State)

	writeSection(w, "Start of turn:", startOfTurn, false)
	writeSection(w, "My turn:", myTurn, true)
	writeSection(w, "Opponent's turn:", oppTurn, true)
	writeSection(w, "End of turn:", endOfTurn, false)
}

// writeSection emits a section header followed by its content lines. numbered=true prefixes
// each line with a step number; numbered=false uses plain indentation. No-op on empty lines.
func writeSection(w io.Writer, header string, lines []string, numbered bool) {
	if len(lines) == 0 {
		return
	}
	fmt.Fprintf(w, "  %s\n", header)
	for i, entry := range lines {
		if numbered {
			fmt.Fprintf(w, "    %d. %s\n", i+1, entry)
		} else {
			fmt.Fprintf(w, "    %s\n", entry)
		}
	}
}

// buildStartOfTurnLines renders the start-of-turn snapshot: dealt hand, arsenal-in card,
// auras / items in play.
func buildStartOfTurnLines(dealtHand []card.Card, line []deck.CardAssignment, master *gameengine.GameState) []string {
	var out []string
	if l := startingHandLine(dealtHand); l != "" {
		out = append(out, l)
	}
	if l := startingArsenalLine(line); l != "" {
		out = append(out, l)
	}
	auras := concreteAuras(master.Auras())
	auraSources := sourceCardsOf(auras)
	runechants := auraCountByName(auras, "Runechant")
	ponders := auraCountByName(auras, "Ponder")
	if l := startingAurasLine(auraSources, runechants, ponders); l != "" {
		out = append(out, l)
	}
	items := concreteItems(master.Items())
	gold := itemCountByName(items, "Gold")
	silver := itemCountByName(items, "Silver")
	copper := itemCountByName(items, "Copper")
	if l := itemsLine(gold, silver, copper); l != "" {
		out = append(out, l)
	}
	return out
}

// buildMyTurnLines emits attack-phase pitch lines followed by one line per card resolved
// in the chain. Chain entries render as "<DisplayName>: <VERB>" via ChainStepText; trigger
// riders and per-step damage deltas are not surfaced.
func buildMyTurnLines(attackPitches []deck.CardAssignment, cardsPlayed []card.Card) []string {
	var out []string
	for _, p := range attackPitches {
		out = append(out, p.Card.DisplayName()+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	for _, c := range cardsPlayed {
		out = append(out, gameengine.ChainStepText(&card.CardState{Card: c}))
	}
	return out
}

// buildOpponentTurnLines emits defense-phase pitches, plain blocks, then Defense Reaction
// rows. DRs render bare ("<name>: DEFENSE REACTION") without per-rider detail.
func buildOpponentTurnLines(defensePitches []deck.CardAssignment, parts bestLineDisplayParts) []string {
	var out []string
	for _, p := range defensePitches {
		out = append(out, p.Card.DisplayName()+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	for _, b := range parts.plainBlocks {
		out = append(out, formatBlockLine(b))
	}
	for _, dr := range parts.defenseReactions {
		out = append(out, dr.Card.DisplayName()+": "+roleLabelWithArsenal(dr, "DEFENSE REACTION"))
	}
	return out
}

// buildEndOfTurnLines renders the end-of-chain snapshot: surviving hand, arsenal-out
// card (with "(stayed)" / "(new)" tag), surviving auras, surviving items.
func buildEndOfTurnLines(arsenal []deck.CardAssignment, state *gameengine.GameState) []string {
	var out []string
	if state == nil {
		if l := endingArsenalLine(arsenal); l != "" {
			out = append(out, l)
		}
		return out
	}
	if l := endingHandLine(state.Hand()); l != "" {
		out = append(out, l)
	}
	if l := endingArsenalLine(arsenal); l != "" {
		out = append(out, l)
	}
	endingAuras := concreteAuras(state.Auras())
	if l := endingAurasLine(endingAuras, auraCountByName(endingAuras, "Runechant"), auraCountByName(endingAuras, "Ponder")); l != "" {
		out = append(out, l)
	}
	endingItems := concreteItems(state.Items())
	if l := itemsLine(itemCountByName(endingItems, "Gold"), itemCountByName(endingItems, "Silver"), itemCountByName(endingItems, "Copper")); l != "" {
		out = append(out, l)
	}
	return out
}

// concreteAuras casts a []gameengine.Aura to []*aura.Aura so the per-aura accessors
// (SourceCard, CardName, Count) become available. Returns nil on empty input.
func concreteAuras(in []gameengine.Aura) []*aura.Aura {
	if len(in) == 0 {
		return nil
	}
	out := make([]*aura.Aura, 0, len(in))
	for _, a := range in {
		out = append(out, a.(*aura.Aura))
	}
	return out
}

// concreteItems is the items counterpart of concreteAuras.
func concreteItems(in []gameengine.Item) []*item.Item {
	if len(in) == 0 {
		return nil
	}
	out := make([]*item.Item, 0, len(in))
	for _, it := range in {
		out = append(out, it.(*item.Item))
	}
	return out
}

// sourceCardsOf returns the SourceCard of each card-style aura. Token auras with nil
// source are skipped.
func sourceCardsOf(in []*aura.Aura) []card.Card {
	if len(in) == 0 {
		return nil
	}
	var out []card.Card
	for _, a := range in {
		if src := a.SourceCard(); src != nil {
			out = append(out, src.(card.Card))
		}
	}
	return out
}
