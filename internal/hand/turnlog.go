package hand

// TurnLog assembly (BuildTurnLog) and rendering (FormatTurnLog).

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// BuildTurnLog converts a TurnSummary into the four-section TurnLog shape. startingRunechants
// is the Runechant carryover entering this turn — surfaced in the StartOfTurn auras line
// alongside any sigils / incantations in play. The chain content for MyTurn comes from
// t.State.Log (the dispatcher's per-event trace); pitches and defense lines come from
// BestLine; ending zone state comes from t.State.{Hand, Arsenal, AuraTriggers, Runechants}.
func BuildTurnLog(t TurnSummary, startingRunechants int) TurnLog {
	var log TurnLog
	parts := partitionBestLineForDisplay(t.BestLine)
	defensePitches, attackPitches := splitPitchesByPhase(parts.pitched, parts.drCost)

	// Start of turn: dealt hand, arsenal-in card, auras / runechants in play. Carryover
	// AuraTrigger fires (Sigil reveals, +N damage credits) belong to MyTurn — they're
	// actions resolving at the top of the action phase, not pre-existing state.
	if line := startingHandLine(t.DealtHand); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := startingArsenalLine(t.BestLine); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := startingAurasLine(t.StartOfTurnAuras, startingRunechants); line != "" {
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
		log.MyTurn = append(log.MyTurn, card.DisplayName(p.Card)+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	// Chain entries are stored as LogEntry structs on State.Log to defer fmt cost off the
	// per-permutation hot path — they're formatted (and grouped) here, in the
	// once-per-best-turn assembly step. appendGroupedChainEntries clusters trigger lines
	// underneath the chain line of the card that fired them so the printout reads as a tree
	// instead of a strict chronological sequence.
	log.MyTurn = appendGroupedChainEntries(log.MyTurn, t.State.Log)

	// Opponent's turn: defense pitches, plain blocks, then Defense Reactions.
	for _, p := range defensePitches {
		log.OpponentTurn = append(log.OpponentTurn, card.DisplayName(p.Card)+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	for _, b := range parts.plainBlocks {
		log.OpponentTurn = append(log.OpponentTurn, formatBlockLine(b))
	}
	for _, dr := range parts.defenseReactions {
		log.OpponentTurn = append(log.OpponentTurn, card.DisplayName(dr.Card)+": "+roleLabelWithArsenal(dr, "DEFENSE REACTION"))
	}

	// End of turn: surviving hand cards, arsenal slot's contents, auras still in play.
	if line := endingHandLine(t.State.Hand); line != "" {
		log.EndOfTurn = append(log.EndOfTurn, line)
	}
	if line := endingArsenalLine(parts.arsenal); line != "" {
		log.EndOfTurn = append(log.EndOfTurn, line)
	}
	if line := endingAurasLine(t.State.AuraTriggers, t.State.Runechants); line != "" {
		log.EndOfTurn = append(log.EndOfTurn, line)
	}

	return log
}

// startingHandLine builds "Hand: A, B, C, D" from the turn's dealt hand — the cards drawn
// at end of last turn, before any start-of-action-phase reveals (Sigil of the Arknight) or
// mid-chain draws bulked the hand. The deck loop captures this snapshot in TurnSummary
// before processTriggersAtStartOfTurn modifies the working hand slice, so reveals show up
// only in MyTurn (where they actually resolve), not in this informational starting-state
// line. Names render in deal order. Returns "" when the dealt hand was empty.
func startingHandLine(dealtHand []card.Card) string {
	if len(dealtHand) == 0 {
		return ""
	}
	names := make([]string, len(dealtHand))
	for i, c := range dealtHand {
		names[i] = card.DisplayName(c)
	}
	return "Hand: " + strings.Join(names, ", ")
}

// startingArsenalLine returns "Arsenal: cardname" when the turn started with an arsenal-in
// card (BestLine entry with FromArsenal=true), "" otherwise.
func startingArsenalLine(line []CardAssignment) string {
	for _, a := range line {
		if a.FromArsenal {
			return "Arsenal: " + card.DisplayName(a.Card)
		}
	}
	return ""
}

// startingAurasLine builds "Auras: A, B, 1 Runechant" from auras in play at the top of the
// turn plus the Runechant carryover. Aura names sort alphabetically; runechants append last.
// Returns "" when both are zero so the caller skips the line entirely.
func startingAurasLine(auras []card.Card, startingRunechants int) string {
	if len(auras) == 0 && startingRunechants == 0 {
		return ""
	}
	var items []string
	if len(auras) > 0 {
		names := make([]string, len(auras))
		for i, a := range auras {
			names[i] = card.DisplayName(a)
		}
		sort.Strings(names)
		items = append(items, names...)
	}
	if startingRunechants > 0 {
		items = append(items, runechantPhrase(startingRunechants))
	}
	return "Auras: " + strings.Join(items, ", ")
}

// startOfTurnTriggerLine renders one carryover AuraTrigger fire as a content line at the top
// of the MyTurn section: "Aura Name: drew X into hand" or "Aura Name: START OF ACTION PHASE
// (+N)" (or both joined when a future trigger does both). Returns "" for zero-effect fires
// so the section doesn't pad with bare aura-name lines.
func startOfTurnTriggerLine(d TriggerContribution) string {
	suffix := formatTriggerEffect(d)
	if suffix == "" {
		return ""
	}
	return card.DisplayName(d.Card) + ": " + suffix
}

// formatBlockLine renders a plain-block content line with the card's effective defense as a
// "(+N)" suffix — printed Defense plus an ArsenalDefenseBonus rider when the blocker came
// from arsenal. The suffix matches the attack chain's "(+N)" convention so the log reads
// symmetrically across attack and defense phases.
func formatBlockLine(a CardAssignment) string {
	def := a.Card.Defense()
	if a.FromArsenal {
		if ab, ok := a.Card.(card.ArsenalDefenseBonus); ok {
			def += ab.ArsenalDefenseBonus()
		}
	}
	return fmt.Sprintf("%s: %s (+%d)", card.DisplayName(a.Card), roleLabelWithArsenal(a, "BLOCK"), def)
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
		names[i] = card.DisplayName(c)
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
		label := card.DisplayName(a.Card)
		if a.FromArsenal {
			label += " (stayed)"
		} else {
			label += " (new)"
		}
		parts[i] = label
	}
	return "Arsenal: " + strings.Join(parts, ", ")
}

// endingAurasLine builds "Auras: A, B, 2 Runechants" from the AuraTriggers surviving into
// the next turn plus the live Runechant count. Aura names sort alphabetically (one entry
// per trigger source — duplicates collapse via the trigger list's natural counts).
// Returns "" when no auras survived and the runechant count is zero.
func endingAurasLine(triggers []card.AuraTrigger, runechants int) string {
	if len(triggers) == 0 && runechants == 0 {
		return ""
	}
	var items []string
	if len(triggers) > 0 {
		names := make([]string, len(triggers))
		for i, t := range triggers {
			names[i] = card.DisplayName(t.Self)
		}
		sort.Strings(names)
		items = append(items, names...)
	}
	if runechants > 0 {
		items = append(items, runechantPhrase(runechants))
	}
	return "Auras: " + strings.Join(items, ", ")
}

// runechantPhrase pluralises the Runechant noun based on count: "1 Runechant" vs
// "N Runechants". Used by both starting and ending auras lines.
func runechantPhrase(n int) string {
	if n == 1 {
		return "1 Runechant"
	}
	return fmt.Sprintf("%d Runechants", n)
}

// childEntryPrefix tags MyTurn entries that are trigger lines grouped beneath a parent chain
// entry. FormatTurnLog detects the prefix, strips it, and renders the entry indented under
// the parent without consuming a step number. The tag is in-band on the string so MyTurn
// stays a flat []string for JSON round-tripping.
const childEntryPrefix = "  "

// appendGroupedChainEntries walks t.State.Log and emits each chain entry as a parent line
// followed by its triggers as childEntryPrefix-tagged children. LogEntry.Kind tells
// pre-triggers from post-triggers — the lookforward step only consumes post-triggers, so
// a pre-trigger whose Source happens to match the previous chain entry's name (two cards
// with the same display name in one chain) gets buffered and attaches to its real parent
// instead. An orphan trigger whose Source matches no chain line falls through as a plain
// top-level entry via FormatLogEntry — keeps "(from <source>)" attribution visible so
// the data isn't silently dropped if the parent-emit invariant ever loosens.
func appendGroupedChainEntries(out []string, log []card.LogEntry) []string {
	var pending []card.LogEntry
	i := 0
	for i < len(log) {
		e := log[i]
		switch e.Kind {
		case card.LogEntryPreTrigger:
			pending = append(pending, e)
			i++
		case card.LogEntryPostTrigger:
			// Orphan post-trigger — no preceding chain step consumed it via the
			// lookforward below. Surface as a top-level line.
			out = append(out, FormatLogEntry(e))
			i++
		default:
			// Chain step: emit parent, attach matching buffered pre-triggers, then look
			// forward for matching post-triggers (only post — pre-triggers belong to a
			// later chain entry).
			parentName := chainEntryCardName(e.Text)
			out = append(out, formatTextWithDelta(e))
			for _, pre := range pending {
				if pre.Source == parentName {
					out = append(out, childEntryPrefix+formatTextWithDelta(pre))
				} else {
					out = append(out, FormatLogEntry(pre))
				}
			}
			pending = pending[:0]
			j := i + 1
			for j < len(log) &&
				log[j].Kind == card.LogEntryPostTrigger &&
				log[j].Source == parentName {
				out = append(out, childEntryPrefix+formatTextWithDelta(log[j]))
				j++
			}
			i = j
		}
	}
	for _, p := range pending {
		out = append(out, FormatLogEntry(p))
	}
	return out
}

// formatTextWithDelta renders a LogEntry as just "<Text> (+N)" — the bare-Text form
// suitable for both chain parents (Source=="") and grouped trigger children where the
// indentation already conveys the source. Drops "(+0)" for zero-value entries. Orphan
// triggers go through FormatLogEntry instead so they keep the "(from <source>)" tail.
func formatTextWithDelta(e card.LogEntry) string {
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
			lines = append(lines, fmt.Sprintf("    %d. %s", nextStep(), entry))
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
