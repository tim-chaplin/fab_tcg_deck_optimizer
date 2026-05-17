package sim

// Per-section rendering helpers for the peak-turn printout. Each helper builds one line
// of content (or "" when the section is empty).

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/turnlogger"
)

// startingHandLine builds "Hand: A, B, C, D" from the turn's dealt hand. Names render in
// deal order. Returns "" when the dealt hand was empty.
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
func startingArsenalLine(line []deck.CardAssignment) string {
	for _, a := range line {
		if a.FromArsenal {
			return "Arsenal: " + a.Card.DisplayName()
		}
	}
	return ""
}

// startingAurasLine builds "Auras: A, B, 1 Runechant, 2 Ponders" from auras in play at
// the top of the turn plus the per-token carryovers. Card-aura names sort alphabetically;
// token phrases append last. Returns "" when nothing is in play.
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

// formatBlockLine renders a plain-block content line with the card's effective defense as
// a "(+N)" suffix — printed Defense plus an ArsenalDefenseBonus rider when blocked from
// arsenal.
func formatBlockLine(a deck.CardAssignment) string {
	def := a.Card.Defense()
	if a.FromArsenal {
		def += card.ArsenalDefenseBonusOf(a.Card)
	}
	return fmt.Sprintf("%s: %s (+%d)", a.Card.DisplayName(), roleLabelWithArsenal(a, "BLOCK"), def)
}

// appendDefenseReactionLines re-Plays the DR with Graveyard seeded to every defender (so
// banish-target scans see the right shape) and the remaining unblocked damage threaded in.
// Renders the chain step's "(+N)" block plus any rider lines as child entries, and returns
// the updated remaining-damage counter for the next DR.
func appendDefenseReactionLines(out []string, a deck.CardAssignment, defenders []card.Card, remaining int) ([]string, int) {
	ge := &gameengine.GameEngine{GameState: gameengine.GameStateBuilder().SetGraveyard(append([]card.Card(nil), defenders...)).Build()}
	ge.SetIncomingDamage(remaining)
	cs := card.CardState{Card: a.Card, FromArsenal: a.FromArsenal}
	ge.ResolveChainStep(ge.Logger(), &cs)
	return appendGroupedChainEntries(out, ge.LogEntries()), ge.RemainingUnblockedDamage()
}

// defendersFromParts collects every card committed to defense (Defense Reactions plus plain
// blocks) into a single slice — used to seed the per-DR Play call's graveyard.
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

// endingArsenalLine builds "Arsenal: cardname (stayed)" / "(new)" from BestLine's Arsenal-
// role entries. "(stayed)" tags the arsenal-in card that kept the slot (FromArsenal=true);
// "(new)" tags any other origin. Returns "" when arsenal ended empty.
func endingArsenalLine(arsenal []deck.CardAssignment) string {
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
// re-rendered as count phrases so pluralisation lives in one place. Returns "" when
// nothing survived.
func endingAurasLine(triggers []*aura.Aura, runechants, ponders int) string {
	var items []string
	for _, t := range triggers {
		if t.SourceCard() == nil {
			continue // token auras render via the count phrases below
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

// itemsLine builds the "Items: ..." printout line from the per-token-type counts. Returns
// "" when no item is in play so the caller skips the line.
func itemsLine(gold, silver, copper int) string {
	var items []string
	if gold > 0 {
		items = append(items, fmt.Sprintf("%d Gold", gold))
	}
	if silver > 0 {
		items = append(items, fmt.Sprintf("%d Silver", silver))
	}
	if copper > 0 {
		items = append(items, fmt.Sprintf("%d Copper", copper))
	}
	if len(items) == 0 {
		return ""
	}
	return "Items: " + strings.Join(items, ", ")
}

// childEntryPrefix tags chain-log entries that are trigger lines grouped beneath a parent
// chain entry. The rendering pass strips the tag and indents the entry under its parent
// without consuming a step number.
const childEntryPrefix = "  "

// appendGroupedChainEntries walks log and emits each chain entry as a parent line followed
// by its triggers as childEntryPrefix-tagged children. Pre/post-triggers are buffered until
// a chain step with a matching Source arrives — cards emit rider lines during Card.Play
// before ResolveChainStep auto-logs the chain step. Orphan triggers fall through via
// FormatLogEntry so a stray producer with no parent isn't silently dropped.
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
			pendingPre, out = drainMatchingTriggers(pendingPre, parentName, out)
			pendingPost, out = drainMatchingTriggers(pendingPost, parentName, out)
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

// formatTextWithDelta renders a LogEntry as "<Text> (+N)" — bare-Text form for chain
// parents and grouped trigger children. Drops "(+0)" for zero-value entries. Orphan
// triggers go through FormatLogEntry instead so they keep the "(from <source>)" tail.
func formatTextWithDelta(e turnlogger.LogEntry) string {
	if e.N == 0 {
		return e.Text
	}
	return fmt.Sprintf("%s (+%d)", e.Text, e.N)
}

// chainEntryCardName extracts the display name from a chain LogEntry's Text shaped
// "<DisplayName>: <VERB>[ from arsenal]". An unrecognised shape returns the raw text so
// appendGroupedChainEntries simply emits the parent alone.
func chainEntryCardName(text string) string {
	if i := strings.Index(text, ": "); i >= 0 {
		return text[:i]
	}
	return text
}
