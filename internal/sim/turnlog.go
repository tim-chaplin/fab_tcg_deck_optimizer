package sim

// Per-section rendering helpers for the peak-turn printout. Each helper builds one line
// of content (or "" when the section is empty).

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
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

// aurasLineFromNames builds "Auras: A, B, 1 Runechant, 2 Ponders" from card-aura display
// names plus the per-token carryovers. Names sort alphabetically; token phrases append
// last. Returns "" when nothing is in play so the caller skips the line entirely.
func aurasLineFromNames(names []string, runechants, ponders int) string {
	if len(names) == 0 && runechants == 0 && ponders == 0 {
		return ""
	}
	var items []string
	if len(names) > 0 {
		sorted := append([]string(nil), names...)
		sort.Strings(sorted)
		items = append(items, sorted...)
	}
	if runechants > 0 {
		items = append(items, runechantPhrase(runechants))
	}
	if ponders > 0 {
		items = append(items, ponderPhrase(ponders))
	}
	return "Auras: " + strings.Join(items, ", ")
}

// endingArsenalLine builds "Arsenal: cardname (stayed)" / "(new)" from BestLine's
// Arsenal-role entries. The tag derives from FromArsenal — the arsenal-in card kept the
// slot when FromArsenal=true ("(stayed)"); any other origin (post-hoc Held promotion)
// reads "(new)". Returns "" when arsenal ended empty.
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
