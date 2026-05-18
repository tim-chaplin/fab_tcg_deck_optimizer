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

// formatBlockLine renders a plain-block content line with the card's effective defense as a
// "(+N)" suffix — printed Defense plus an ArsenalDefenseBonus rider when the blocker came
// from arsenal.
func formatBlockLine(a deck.CardAssignment) string {
	def := a.Card.Defense()
	if a.FromArsenal {
		def += card.ArsenalDefenseBonusOf(a.Card)
	}
	return fmt.Sprintf("%s: %s (+%d)", a.Card.DisplayName(), roleLabelWithArsenal(a, "BLOCK"), def)
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

// endingAurasLine builds "Auras: A, B, 2 Runechants, 1 Ponder" from the auras surviving
// into the next turn plus the live token counts. Token auras are filtered out and
// re-rendered as count phrases so pluralisation lives in one place. Card-aura names
// sort alphabetically; token phrases append last in declaration order. Returns "" when
// nothing survived.
func endingAurasLine(triggers []*aura.Aura, runechants, ponders int) string {
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
