package fabrary

// Runtime Deck → fabrary-style text encoding: Marshal assembles the Name / Hero / Format
// header plus the Arena cards / Deck cards / Sideboard sections. Stats aren't emitted —
// fabrary doesn't consume them and they round-trip through deckio instead.

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapons"
)

// defaultFormat is emitted in the Format: header. Update when a new format comes online.
const defaultFormat = "Silver Age"

// Marshal returns fabrary-style deck text for d, suitable for pasting into fabrary.net's
// "Import deck" tab. The output sections are:
//
//   - Arena cards: weapons + d.Equipment.
//   - Deck cards: d.Cards, pitch color suffix lowercased to match fabrary.
//   - Sideboard: d.Sideboard, lowercased pitch suffix. Empty when d.Sideboard is empty.
//
// Callers that want the hardcoded default equipment / sideboard loadout baked in should run
// d.ApplyDefaults() before Marshal. writeDeck does that automatically so the persisted .txt
// always carries the full loadout.
func Marshal(d *deck.Deck) string {
	var b strings.Builder
	name := d.Hero.Name()
	fmt.Fprintf(&b, "Name: %s\n", name)
	fmt.Fprintf(&b, "Hero: %s\n", name)
	fmt.Fprintf(&b, "Format: %s\n\n", defaultFormat)

	b.WriteString("Arena cards\n")
	arena := weaponCounts(d.Weapons)
	for _, name := range d.Equipment {
		arena[name]++
	}
	writeCounts(&b, arena)
	b.WriteString("\n")

	b.WriteString("Deck cards\n")
	writeCounts(&b, cardCountsForExport(d.Cards))

	sideboardCounts := sideboardCountsForExport(d.Sideboard)
	if len(sideboardCounts) > 0 {
		b.WriteString("\nSideboard\n")
		writeCounts(&b, sideboardCounts)
	}
	return b.String()
}

func weaponCounts(ws []weapons.Weapon) map[string]int {
	m := make(map[string]int, len(ws))
	for _, w := range ws {
		m[w.Name()]++
	}
	return m
}

func cardCountsForExport(cs []card.Card) map[string]int {
	m := make(map[string]int, len(cs))
	for _, c := range cs {
		m[toFabraryCardName(card.DisplayName(c))]++
	}
	return m
}

// sideboardCountsForExport mirrors cardCountsForExport for the Sideboard — a string slice
// rather than []card.Card, since the optimizer doesn't resolve sideboard entries through
// the card registry. Names are converted to fabrary's lowercase-pitch-color form if they
// match a known canonical suffix.
func sideboardCountsForExport(ss []string) map[string]int {
	m := make(map[string]int, len(ss))
	for _, s := range ss {
		m[toFabraryCardName(s)]++
	}
	return m
}

func writeCounts(b *strings.Builder, m map[string]int) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Fprintf(b, "%dx %s\n", m[k], k)
	}
}
