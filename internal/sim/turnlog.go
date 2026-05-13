package sim

// TurnLog assembly (BuildTurnLog) and rendering (FormatTurnLog).

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// BuildTurnLog converts a TurnSummary into the four-section TurnLog shape.
// startingAuras / startingItems are the carryover aura / item sets entering this turn;
// token counts get pulled from them for the StartOfTurn "Auras: ..." / "Items: ..."
// lines. MyTurn's chain content comes from t.State.Log; pitches and defense lines come
// from BestLine; ending zone state comes from t.State.
func BuildTurnLog(t TurnSummary, startingAuras []*Aura, startingItems []*Item) TurnLog {
	var log TurnLog
	parts := partitionBestLineForDisplay(t.BestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)
	startingRunechants := auraCountByName(startingAuras, "Runechant")
	startingPonders := auraCountByName(startingAuras, "Ponder")
	startingGold := itemCountByName(startingItems, "Gold")
	startingSilver := itemCountByName(startingItems, "Silver")
	startingCopper := itemCountByName(startingItems, "Copper")

	// Start of turn: dealt hand, arsenal-in card, auras / tokens in play. Carryover
	// Aura fires (Sigil reveals, +N damage credits) belong to MyTurn — they're
	// actions resolving at the top of the action phase, not pre-existing state.
	if line := startingHandLine(t.DealtHand); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := startingArsenalLine(t.BestLine); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := startingAurasLine(t.StartOfTurnAuras, startingRunechants, startingPonders); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := itemsLine(startingGold, startingSilver, startingCopper); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}

	// My turn: carryover trigger fires first (resolve at top of action phase), then
	// attack-phase pitches, then the chain log.
	for _, trig := range t.TriggersFromLastTurn {
		if line := startOfTurnTriggerLine(trig); line != "" {
			log.MyTurn = append(log.MyTurn, line)
		}
	}
	for _, p := range attackPitches {
		log.MyTurn = append(log.MyTurn, p.Card.DisplayName()+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	// Chain entries are stored as LogEntry structs on State.Log to defer fmt cost off the
	// per-permutation hot path — they're formatted (and grouped) here, in the
	// once-per-best-turn assembly step. appendGroupedChainEntries clusters trigger lines
	// underneath the chain line of the card that fired them so the printout reads as a tree
	// instead of a strict chronological sequence.
	if t.State != nil {
		log.MyTurn = appendGroupedChainEntries(log.MyTurn, t.State.LogEntries())
	}

	// Opponent's turn: defense pitches, plain blocks, then Defense Reactions.
	for _, p := range defensePitches {
		log.OpponentTurn = append(log.OpponentTurn, p.Card.DisplayName()+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	for _, b := range parts.plainBlocks {
		log.OpponentTurn = append(log.OpponentTurn, formatBlockLine(b))
	}
	defenders := defendersFromParts(parts)
	remaining := t.IncomingDamage
	for _, dr := range parts.defenseReactions {
		log.OpponentTurn, remaining = appendDefenseReactionLines(log.OpponentTurn, dr, defenders, remaining)
	}

	// End of turn: surviving hand cards, arsenal slot's contents, auras still in play.
	if t.State != nil {
		if line := endingHandLine(t.State.Hand()); line != "" {
			log.EndOfTurn = append(log.EndOfTurn, line)
		}
	}
	if line := endingArsenalLine(parts.arsenal); line != "" {
		log.EndOfTurn = append(log.EndOfTurn, line)
	}
	if t.State != nil {
		endingAuras := make([]*Aura, 0, len(t.State.Auras()))
		for _, a := range t.State.Auras() {
			endingAuras = append(endingAuras, a.(*Aura))
		}
		if line := endingAurasLine(endingAuras, auraCountByNameInState(t.State, "Runechant"), auraCountByNameInState(t.State, "Ponder")); line != "" {
			log.EndOfTurn = append(log.EndOfTurn, line)
		}
		if line := itemsLine(itemCountByNameInState(t.State, "Gold"), itemCountByNameInState(t.State, "Silver"), itemCountByNameInState(t.State, "Copper")); line != "" {
			log.EndOfTurn = append(log.EndOfTurn, line)
		}
	}

	return log
}

// startingHandLine builds "Hand: A, B, C, D" from the turn's dealt hand — the cards drawn
// at end of last turn, before any start-of-action-phase reveals (Sigil of the Arknight) or
// mid-chain draws bulked the hand. The deck loop captures this snapshot in TurnSummary
// before processAurasAtStartOfTurn modifies the working hand slice, so reveals show up
// only in MyTurn (where they actually resolve), not in this informational starting-state
// line. Names render in deal order. Returns "" when the dealt hand was empty.
func startingHandLine(dealtHand []card.Card) string {
	if len(dealtHand) == 0 {
		return ""
	}
	names := make([]string, len(dealtHand))
	for i, c := range dealtHand {
		names[i] = c.DisplayName()
	}
	return "Hand: " + strings.Join(names, ", ")
}

// startingArsenalLine returns "Arsenal: cardname" when the turn started with an arsenal-in
// card (BestLine entry with FromArsenal=true), "" otherwise.
func startingArsenalLine(line []CardAssignment) string {
	for _, a := range line {
		if a.FromArsenal {
			return "Arsenal: " + a.Card.DisplayName()
		}
	}
	return ""
}

// startingAurasLine builds "Auras: A, B, 1 Runechant, 2 Ponders" from auras in play at
// the top of the turn plus the per-token carryovers. Card-aura names sort alphabetically;
// token phrases append last in token-declaration order. Returns "" when nothing is in
// play so the caller skips the line entirely.
func startingAurasLine(auras []card.Card, startingRunechants, startingPonders int) string {
	if len(auras) == 0 && startingRunechants == 0 && startingPonders == 0 {
		return ""
	}
	var items []string
	if len(auras) > 0 {
		names := make([]string, len(auras))
		for i, a := range auras {
			names[i] = a.DisplayName()
		}
		sort.Strings(names)
		items = append(items, names...)
	}
	if startingRunechants > 0 {
		items = append(items, runechantPhrase(startingRunechants))
	}
	if startingPonders > 0 {
		items = append(items, ponderPhrase(startingPonders))
	}
	return "Auras: " + strings.Join(items, ", ")
}

// startOfTurnTriggerLine renders one carryover Aura fire as a content line at the
// top of the MyTurn section. When d.Text is set the card owns the wording and renders
// verbatim; otherwise we synthesise "Aura Name: drew X into hand" / "Aura Name: START
// OF ACTION PHASE (+N)" from Damage and Revealed. Returns "" for zero-effect fires so
// the section doesn't pad with bare aura-name lines.
func startOfTurnTriggerLine(d TriggerContribution) string {
	if d.Text != "" {
		return d.Text
	}
	suffix := formatTriggerEffect(d)
	if suffix == "" {
		return ""
	}
	return d.Card.DisplayName() + ": " + suffix
}

// formatBlockLine renders a plain-block content line with the card's effective defense as a
// "(+N)" suffix — printed Defense plus an ArsenalDefenseBonus rider when the blocker came
// from arsenal. The suffix matches the attack chain's "(+N)" convention so the log reads
// symmetrically across attack and defense phases.
func formatBlockLine(a CardAssignment) string {
	def := a.Card.Defense()
	if a.FromArsenal {
		def += card.ArsenalDefenseBonusOf(a.Card)
	}
	return fmt.Sprintf("%s: %s (+%d)", a.Card.DisplayName(), roleLabelWithArsenal(a, "BLOCK"), def)
}

// appendDefenseReactionLines re-Plays the DR with the same fresh state shape defendersDamage
// uses (Graveyard = all defenders so banish-target scans see the right shape; IncomingDamage
// threaded across the DR loop so each defender sees what's left after earlier ones blocked)
// and walks the resulting log: the chain step renders with its own (+N) for the block, and
// any arcane / runechant / +1{d} riders attach as childEntryPrefix-tagged sub-lines via
// appendGroupedChainEntries. Returns the updated remaining-incoming counter so the caller
// can thread it into the next DR.
func appendDefenseReactionLines(out []string, a CardAssignment, defenders []card.Card, remaining int) ([]string, int) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard(append([]card.Card(nil), defenders...)).Build()}
	ge.SetIncomingDamage(remaining)
	cs := card.CardState{Card: a.Card, FromArsenal: a.FromArsenal}
	ge.ResolveChainStep(ge.Logger(), &cs)
	return appendGroupedChainEntries(out, ge.LogEntries()), ge.IncomingDamage()
}

// defendersFromParts collects every card committed to defense — Defense Reactions and plain
// blocks — into a single slice. Mirrors the defenders argument defendersDamage takes during
// partition evaluation so the formatter's per-DR Play call sees the same graveyard.
func defendersFromParts(parts bestLineDisplayParts) []card.Card {
	out := make([]card.Card, 0, len(parts.defenseReactions)+len(parts.plainBlocks))
	for _, a := range parts.defenseReactions {
		out = append(out, a.Card)
	}
	for _, a := range parts.plainBlocks {
		out = append(out, a.Card)
	}
	return out
}

// endingHandLine builds "Hand: A, B" from the cards in hand at end of chain — the partition's
// Held set plus anything tutored / drawn that didn't get played. Returns "" when the hand
// ended empty.
func endingHandLine(handHeld []card.Card) string {
	if len(handHeld) == 0 {
		return ""
	}
	names := make([]string, len(handHeld))
	for i, c := range handHeld {
		names[i] = c.DisplayName()
	}
	return "Hand: " + strings.Join(names, ", ")
}

// endingArsenalLine builds "Arsenal: cardname (stayed)" / "(new)" from BestLine's
// Arsenal-role entries. The tag derives from FromArsenal — the arsenal-in card kept the
// slot when FromArsenal=true ("(stayed)"); any other origin (post-hoc Held promotion)
// reads "(new)". Returns "" when arsenal ended empty.
func endingArsenalLine(arsenal []CardAssignment) string {
	if len(arsenal) == 0 {
		return ""
	}
	parts := make([]string, len(arsenal))
	for i, a := range arsenal {
		label := a.Card.DisplayName()
		if a.FromArsenal {
			label += " (stayed)"
		} else {
			label += " (new)"
		}
		parts[i] = label
	}
	return "Arsenal: " + strings.Join(parts, ", ")
}

// endingAurasLine builds "Auras: A, B, 2 Runechants, 1 Ponder" from the Auras surviving
// into the next turn plus the live token counts. Token auras are filtered out and
// re-rendered as count phrases so pluralisation lives in one place. Card-aura names
// sort alphabetically; token phrases append last in declaration order. Returns "" when
// nothing survived.
func endingAurasLine(triggers []*Aura, runechants, ponders int) string {
	var items []string
	for _, t := range triggers {
		if t.SourceCard() == nil {
			continue // token auras collapse into the token phrases below
		}
		items = append(items, t.CardName())
	}
	sort.Strings(items)
	if runechants > 0 {
		items = append(items, runechantPhrase(runechants))
	}
	if ponders > 0 {
		items = append(items, ponderPhrase(ponders))
	}
	if len(items) == 0 {
		return ""
	}
	return "Auras: " + strings.Join(items, ", ")
}

// runechantPhrase pluralises the Runechant noun: "1 Runechant" vs "N Runechants".
func runechantPhrase(n int) string {
	if n == 1 {
		return "1 Runechant"
	}
	return fmt.Sprintf("%d Runechants", n)
}

// ponderPhrase pluralises the Ponder noun: "1 Ponder" vs "N Ponders".
func ponderPhrase(n int) string {
	if n == 1 {
		return "1 Ponder"
	}
	return fmt.Sprintf("%d Ponders", n)
}

// goldPhrase pluralises the Gold noun: "1 Gold" vs "N Gold" (Gold doesn't pluralise
// in FaB rules text — "5 Gold tokens" reads as bare count + noun).
func goldPhrase(n int) string {
	return fmt.Sprintf("%d Gold", n)
}

// silverPhrase mirrors goldPhrase for Silver tokens.
func silverPhrase(n int) string {
	return fmt.Sprintf("%d Silver", n)
}

// copperPhrase mirrors goldPhrase for Copper tokens.
func copperPhrase(n int) string {
	return fmt.Sprintf("%d Copper", n)
}

// itemsLine builds the "Items: ..." printout line from the per-token-type counts. Returns
// "" when no item is in play so the caller skips the line.
func itemsLine(gold, silver, copper int) string {
	var items []string
	if gold > 0 {
		items = append(items, goldPhrase(gold))
	}
	if silver > 0 {
		items = append(items, silverPhrase(silver))
	}
	if copper > 0 {
		items = append(items, copperPhrase(copper))
	}
	if len(items) == 0 {
		return ""
	}
	return "Items: " + strings.Join(items, ", ")
}

// childEntryPrefix tags MyTurn entries that are trigger lines grouped beneath a parent chain
// entry. FormatTurnLog detects the prefix, strips it, and renders the entry indented under
// the parent without consuming a step number. The tag is in-band on the string so MyTurn
// stays a flat []string for JSON round-tripping.
const childEntryPrefix = "  "

// appendGroupedChainEntries walks t.State.Log and emits each chain entry as a parent
// line followed by its triggers as childEntryPrefix-tagged children. Both pre-triggers
// and post-triggers are buffered until a chain step with a matching Source arrives —
// the sim emits cards' rider lines during Card.Play, but the chain step itself gets
// auto-logged by ResolveChainStep after Play returns, so a card's post-triggers come
// before its chain step in the stream. Buffering both sides keeps the attribution
// correct regardless of stream order. An orphan trigger whose Source matches no chain
// line falls through as a plain top-level entry via FormatLogEntry, keeping the
// "(from <source>)" tail visible so the data isn't silently dropped if a producer
// ever emits a trigger with no parent.
func appendGroupedChainEntries(out []string, log []turnlogger.LogEntry) []string {
	var pendingPre, pendingPost []turnlogger.LogEntry
	i := 0
	for i < len(log) {
		e := log[i]
		switch e.Kind {
		case turnlogger.LogEntryPreTrigger:
			pendingPre = append(pendingPre, e)
			i++
		case turnlogger.LogEntryPostTrigger:
			pendingPost = append(pendingPost, e)
			i++
		default:
			parentName := chainEntryCardName(e.Text)
			out = append(out, formatTextWithDelta(e))
			// Drain pre/post-triggers that arrived BEFORE this chain step (the
			// ResolveChainStep path: Card.Play emits riders, then sim appends the
			// chain step). Pre-triggers always come from BEFORE; we drain them too.
			pendingPre, out = drainMatchingTriggers(pendingPre, parentName, out)
			pendingPost, out = drainMatchingTriggers(pendingPost, parentName, out)
			// Look forward for adjacent post-triggers that arrived AFTER this chain
			// step. Covers manually-built test streams and any producer that appends
			// rider lines after emitting its own chain step.
			j := i + 1
			for j < len(log) && log[j].Kind == turnlogger.LogEntryPostTrigger && log[j].Source == parentName {
				out = append(out, childEntryPrefix+formatTextWithDelta(log[j]))
				j++
			}
			i = j
		}
	}
	for _, p := range pendingPre {
		out = append(out, FormatLogEntry(p))
	}
	for _, p := range pendingPost {
		out = append(out, FormatLogEntry(p))
	}
	return out
}

// drainMatchingTriggers splits pending into entries whose Source matches name (emitted
// as indented children of the chain step) and entries that don't (returned for a later
// chain step to consume). The split preserves arrival order within each bucket so
// pre/post-triggers render in the order their producers fired.
func drainMatchingTriggers(pending []turnlogger.LogEntry, name string, out []string) ([]turnlogger.LogEntry, []string) {
	var remaining []turnlogger.LogEntry
	for _, t := range pending {
		if t.Source == name {
			out = append(out, childEntryPrefix+formatTextWithDelta(t))
		} else {
			remaining = append(remaining, t)
		}
	}
	return remaining, out
}

// formatTextWithDelta renders a LogEntry as just "<Text> (+N)" — the bare-Text form
// suitable for both chain parents (Source=="") and grouped trigger children where the
// indentation already conveys the source. Drops "(+0)" for zero-value entries. Orphan
// triggers go through FormatLogEntry instead so they keep the "(from <source>)" tail.
func formatTextWithDelta(e turnlogger.LogEntry) string {
	if e.N == 0 {
		return e.Text
	}
	return fmt.Sprintf("%s (+%d)", e.Text, e.N)
}

// chainEntryCardName extracts the display name from a chain LogEntry's Text. Chain entries
// the sim writes are shaped "<DisplayName>: <VERB>[ from arsenal]", so the name is
// everything up to the first ": ". An unrecognised shape returns the raw text —
// appendGroupedChainEntries then matches no triggers to it and just emits the parent alone.
func chainEntryCardName(text string) string {
	if i := strings.Index(text, ": "); i >= 0 {
		return text[:i]
	}
	return text
}

// FormatTurnLog renders a TurnLog into a printable string. Section headers and indentation
// come from the formatter; chain events in MyTurn / OpponentTurn get a continuous numbered
// prefix ("    1. ..."), while StartOfTurn / EndOfTurn entries are unnumbered informational
// lines indented two extra spaces beneath the section header. MyTurn entries tagged with
// childEntryPrefix render indented under the preceding numbered entry without their own
// step number — this is the trigger-grouping mode appendGroupedChainEntries produces.
func FormatTurnLog(log TurnLog) string {
	var lines []string
	step := 0
	nextStep := func() int { step++; return step }

	if len(log.StartOfTurn) > 0 {
		lines = append(lines, "  Start of turn:")
		for _, entry := range log.StartOfTurn {
			lines = append(lines, "    "+entry)
		}
	}
	if len(log.MyTurn) > 0 {
		lines = append(lines, "  My turn:")
		for _, entry := range log.MyTurn {
			if rest, isChild := strings.CutPrefix(entry, childEntryPrefix); isChild {
				// 9-space indent visually nests the child under the parent's card name.
				// Aligns under the typical "    N. " parent prefix (4 + 3 = 7) plus 2
				// extra spaces of padding; 2-digit step numbers (10+) shift the parent
				// right one column so the child sits one column shy of perfect alignment
				// — acceptable readability tradeoff for a fixed-width indent.
				lines = append(lines, "         "+rest)
			} else {
				lines = append(lines, fmt.Sprintf("    %d. %s", nextStep(), entry))
			}
		}
	}
	if len(log.OpponentTurn) > 0 {
		lines = append(lines, "  Opponent's turn:")
		for _, entry := range log.OpponentTurn {
			if rest, isChild := strings.CutPrefix(entry, childEntryPrefix); isChild {
				lines = append(lines, "         "+rest)
			} else {
				lines = append(lines, fmt.Sprintf("    %d. %s", nextStep(), entry))
			}
		}
	}
	if len(log.EndOfTurn) > 0 {
		lines = append(lines, "  End of turn:")
		for _, entry := range log.EndOfTurn {
			lines = append(lines, "    "+entry)
		}
	}
	return strings.Join(lines, "\n")
}
