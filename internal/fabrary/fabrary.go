// Package fabrary converts a deck.Deck to and from the plain-text deck format used by fabrary.net
// (https://fabrary.net/decks?tab=import). The format has a `Name:` / `Hero:` / `Format:` header,
// an "Arena cards" section listing equipment and weapons, and a "Deck cards" section listing
// pitch cards with lowercase color suffix (e.g. "2x Aether Slash (red)").
//
// The optimizer only models weapons, not non-weapon equipment (helms, chests, etc.). On import,
// Arena lines that don't name a known weapon are ignored. On export, only weapons appear in the
// Arena section — users will typically re-add equipment in fabrary after pasting.
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

// heroesByName resolves hero names during Unmarshal. Add new heroes here as they're implemented.
var heroesByName = map[string]hero.Hero{
	(hero.Viserai{}).Name(): hero.Viserai{},
}

// defaultFormat is emitted in the Format: header. Silver Age is the current Viserai format;
// update when a new format comes online.
const defaultFormat = "Silver Age"

// Marshal returns fabrary-style deck text for `d`, suitable for pasting into fabrary.net's
// "Import deck" tab. Weapons are listed in the Arena section; deck cards in the Deck section with
// pitch color suffix lowercased to match fabrary's own exports.
func Marshal(d *deck.Deck) string {
	var b strings.Builder
	name := d.Hero.Name()
	fmt.Fprintf(&b, "Name: %s\n", name)
	fmt.Fprintf(&b, "Hero: %s\n", name)
	fmt.Fprintf(&b, "Format: %s\n\n", defaultFormat)

	b.WriteString("Arena cards\n")
	writeCounts(&b, weaponCounts(d.Weapons))
	b.WriteString("\n")

	b.WriteString("Deck cards\n")
	writeCounts(&b, cardCountsForExport(d.Cards))
	return b.String()
}

// Unmarshal parses fabrary-style deck text and returns a *deck.Deck plus a count-keyed map of
// deck cards whose names aren't in the optimizer's registry (typically cards implemented in
// fabrary but not yet in this project). Callers should surface the skipped map so users aren't
// surprised by silently-reduced deck size. Stats aren't round-tripped (they're a simulation
// artifact).
//
// Unknown Arena-section lines (non-weapon equipment) are silently skipped and NOT reported, since
// the optimizer doesn't model non-weapon equipment at all and reporting would be noise.
// Hero/format parse errors still abort — a missing hero means we can't build a deck.
func Unmarshal(text string) (*deck.Deck, map[string]int, error) {
	var (
		heroName string
		section  string
		weapons  []weapon.Weapon
		cardList []card.Card
		skipped  = map[string]int{}
	)
	wReg := weaponsByName()

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
			if w, ok := wReg[name]; ok {
				for i := 0; i < qty; i++ {
					weapons = append(weapons, w)
				}
			}
		case "deck":
			canon := fromFabraryCardName(name)
			id, ok := cards.ByName(canon)
			if !ok {
				skipped[canon] += qty
				continue
			}
			c := cards.Get(id)
			for i := 0; i < qty; i++ {
				cardList = append(cardList, c)
			}
		}
	}
	if err := sc.Err(); err != nil {
		return nil, nil, err
	}
	h, ok := heroesByName[heroName]
	if !ok {
		return nil, nil, fmt.Errorf("fabrary: unknown hero %q", heroName)
	}
	return deck.New(h, weapons, cardList), skipped, nil
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

func weaponsByName() map[string]weapon.Weapon {
	m := make(map[string]weapon.Weapon, len(cards.AllWeapons))
	for _, w := range cards.AllWeapons {
		m[w.Name()] = w
	}
	return m
}
