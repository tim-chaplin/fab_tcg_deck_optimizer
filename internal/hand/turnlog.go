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

	// Start of turn: dealt hand, arsenal-in card, auras / runechants in play, then any
	// carryover AuraTrigger fires.
	if line := startingHandLine(t.BestLine); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := startingArsenalLine(t.BestLine); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	if line := startingAurasLine(t.StartOfTurnAuras, startingRunechants); line != "" {
		log.StartOfTurn = append(log.StartOfTurn, line)
	}
	for _, trig := range t.TriggersFromLastTurn {
		if line := startOfTurnTriggerLine(trig); line != "" {
			log.StartOfTurn = append(log.StartOfTurn, line)
		}
	}

	// My turn: attack-phase pitches followed by the chain log. Chain entries are stored as
	// LogEntry structs on State.Log to defer fmt cost off the per-permutation hot path —
	// they're formatted to strings here, in the once-per-best-turn assembly step.
	for _, p := range attackPitches {
		log.MyTurn = append(log.MyTurn, card.DisplayName(p.Card)+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	for _, e := range t.State.Log {
		log.MyTurn = append(log.MyTurn, FormatLogEntry(e))
	}

	// Opponent's turn: defense pitches, plain blocks, then Defense Reactions.
	for _, p := range defensePitches {
		log.OpponentTurn = append(log.OpponentTurn, card.DisplayName(p.Card)+": "+roleLabelWithArsenal(p, "PITCH"))
	}
	for _, b := range parts.plainBlocks {
		log.OpponentTurn = append(log.OpponentTurn, card.DisplayName(b.Card)+": "+roleLabelWithArsenal(b, "BLOCK"))
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

// startingHandLine builds "Hand: A, B, C, D" from BestLine's hand cards (everything except
// the arsenal-in entry). Names render in BestLine order — that's the canonical post-sort
// order the partition enumerator settled on, which keeps the output deterministic across
// runs of the same hand. Returns "" when no hand cards exist.
func startingHandLine(line []CardAssignment) string {
	var names []string
	for _, a := range line {
		if a.FromArsenal {
			continue
		}
		names = append(names, card.DisplayName(a.Card))
	}
	if len(names) == 0 {
		return ""
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

// startOfTurnTriggerLine renders one carryover AuraTrigger fire as a content line for the
// StartOfTurn section: "Aura Name: drew X into hand" or "Aura Name: START OF ACTION PHASE
// (+N)" (or both joined when a future trigger does both). Returns "" for zero-effect fires
// so the section doesn't pad with bare aura-name lines.
func startOfTurnTriggerLine(d TriggerContribution) string {
	suffix := formatTriggerEffect(d)
	if suffix == "" {
		return ""
	}
	return card.DisplayName(d.Card) + ": " + suffix
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

// FormatTurnLog renders a TurnLog into a printable string. Section headers and indentation
// come from the formatter; chain events in MyTurn / OpponentTurn get a continuous numbered
// prefix ("    1. ..."), while StartOfTurn / EndOfTurn entries are unnumbered informational
// lines indented two extra spaces beneath the section header.
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
			lines = append(lines, fmt.Sprintf("    %d. %s", nextStep(), entry))
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
