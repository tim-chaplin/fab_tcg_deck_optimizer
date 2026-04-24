// Package fabrary converts a deck.Deck to and from fabrary.net's plain-text deck format
// (https://fabrary.net/decks?tab=import). The format has a `Name:` / `Hero:` / `Format:` header,
// an "Arena cards" section for equipment and weapons, a "Deck cards" section with pitch cards
// carrying a lowercase color suffix (e.g. "2x Aether Slash (red)"), and an optional
// "Sideboard" section mirroring the Deck section for the user-managed sideboard.
//
// The optimizer models only weapons, not other equipment. Unknown Arena lines are ignored on
// import; on export, modelled weapons are joined by the fixed equipment loadout in
// defaultArenaPackage so the emitted .txt can be pasted into fabrary without hand-editing.
package fabrary

import (
	"bufio"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// defaultFormat is emitted in the Format: header. Update when a new format comes online.
const defaultFormat = "Silver Age"

// defaultArenaPackage is the fixed equipment loadout the user runs on every exported deck.
// The optimizer doesn't model these slots (only weapons are modelled and varied), so Marshal
// emits them verbatim into the Arena cards section. Modelled weapons (e.g. Nebula Blade)
// are NOT listed here — they reach the Arena section via d.Weapons when the deck uses them.
var defaultArenaPackage = []string{
	"Beckoning Haunt",
	"Blade Beckoner Boots",
	"Blade Beckoner Helm",
	"Blossom of Spring",
	"Crown of Dichotomy",
	"Nullrune Boots",
	"Nullrune Gloves",
	"Runebleed Robe",
}

// Marshal returns fabrary-style deck text for d, suitable for pasting into fabrary.net's
// "Import deck" tab. Weapons go in the Arena section; deck cards in the Deck section with pitch
// color suffix lowercased to match fabrary's own exports. When the deck carries a sideboard a
// trailing "Sideboard" section is appended with the same count-and-name formatting; an empty
// sideboard skips the section entirely so the output stays minimal.
func Marshal(d *deck.Deck) string {
	var b strings.Builder
	name := d.Hero.Name()
	fmt.Fprintf(&b, "Name: %s\n", name)
	fmt.Fprintf(&b, "Hero: %s\n", name)
	fmt.Fprintf(&b, "Format: %s\n\n", defaultFormat)

	b.WriteString("Arena cards\n")
	arena := weaponCounts(d.Weapons)
	for _, name := range defaultArenaPackage {
		arena[name]++
	}
	writeCounts(&b, arena)
	b.WriteString("\n")

	b.WriteString("Deck cards\n")
	writeCounts(&b, cardCountsForExport(d.Cards))

	if len(d.Sideboard) > 0 {
		b.WriteString("\nSideboard\n")
		writeCounts(&b, cardCountsForExport(d.Sideboard))
	}
	return b.String()
}

// Unmarshal parses fabrary-style deck text and returns a *deck.Deck plus a count-keyed map of
// deck cards whose names aren't in the optimizer's registry. Callers should surface the skipped
// map so users aren't surprised by a silently-reduced deck. Stats aren't round-tripped.
//
// Unknown Arena-section lines (non-weapon equipment) are silently skipped and NOT reported — the
// optimizer doesn't model them. A missing hero aborts: the deck can't be constructed without one.
func Unmarshal(text string) (*deck.Deck, map[string]int, error) {
	var (
		heroName  string
		section   string
		weapons   []weapon.Weapon
		cardList  []card.Card
		sideboard []card.Card
		skipped   = map[string]int{}
	)

	sc := bufio.NewScanner(strings.NewReader(text))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if rest, ok := trimHeader(line, "Hero:"); ok {
			heroName = rest
			continue
		}
		if _, ok := trimHeader(line, "Name:"); ok {
			continue
		}
		if _, ok := trimHeader(line, "Format:"); ok {
			continue
		}
		switch line {
		case "Arena cards":
			section = "arena"
			continue
		case "Deck cards":
			section = "deck"
			continue
		case "Sideboard":
			section = "sideboard"
			continue
		}
		if isFooter(line) {
			continue
		}
		qty, name, ok := parseCountedLine(line)
		if !ok {
			continue
		}
		switch section {
		case "arena":
			if w, ok := cards.WeaponByName(name); ok {
				for i := 0; i < qty; i++ {
					weapons = append(weapons, w)
				}
			}
		case "deck", "sideboard":
			canon := fromFabraryCardName(name)
			id, ok := cards.ByName(canon)
			if !ok {
				skipped[canon] += qty
				continue
			}
			c := cards.Get(id)
			for i := 0; i < qty; i++ {
				if section == "sideboard" {
					sideboard = append(sideboard, c)
				} else {
					cardList = append(cardList, c)
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, nil, err
	}
	h, ok := hero.ByName(heroName)
	if !ok {
		return nil, nil, fmt.Errorf("fabrary: unknown hero %q", heroName)
	}
	d := deck.New(h, weapons, cardList)
	d.Sideboard = sideboard
	return d, skipped, nil
}

// countedLine matches "<N>x <name>" — fabrary always uses a lowercase "x" with no spaces around it.
var countedLine = regexp.MustCompile(`^(\d+)x\s+(.+?)\s*$`)

func parseCountedLine(line string) (int, string, bool) {
	m := countedLine.FindStringSubmatch(line)
	if m == nil {
		return 0, "", false
	}
	qty, err := strconv.Atoi(m[1])
	if err != nil || qty <= 0 {
		return 0, "", false
	}
	return qty, m[2], true
}

func trimHeader(line, prefix string) (string, bool) {
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, prefix)), true
}

// isFooter matches fabrary's trailing "Made with ❤️ at the FaBrary" / "See the full deck @ ..."
// lines so pastes with the footer still round-trip cleanly.
func isFooter(line string) bool {
	return strings.HasPrefix(line, "Made with") || strings.HasPrefix(line, "See the full deck")
}

func weaponCounts(ws []weapon.Weapon) map[string]int {
	m := make(map[string]int, len(ws))
	for _, w := range ws {
		m[w.Name()]++
	}
	return m
}

func cardCountsForExport(cs []card.Card) map[string]int {
	m := make(map[string]int, len(cs))
	for _, c := range cs {
		m[toFabraryCardName(c.Name())]++
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

// pitchColors pairs the optimizer's canonical suffix ("(Red)") with fabrary's lowercase form
// ("(red)"). One entry per color is enough — suffixes don't overlap.
var pitchColors = []struct{ canon, fabrary string }{
	{"(Red)", "(red)"},
	{"(Yellow)", "(yellow)"},
	{"(Blue)", "(blue)"},
}

func toFabraryCardName(s string) string {
	for _, p := range pitchColors {
		if strings.HasSuffix(s, p.canon) {
			return strings.TrimSuffix(s, p.canon) + p.fabrary
		}
	}
	return s
}

func fromFabraryCardName(s string) string {
	for _, p := range pitchColors {
		if strings.HasSuffix(s, p.fabrary) {
			return strings.TrimSuffix(s, p.fabrary) + p.canon
		}
	}
	return s
}

