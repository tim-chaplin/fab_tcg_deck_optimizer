package fabrary

// Fabrary text → runtime Deck decoding: Unmarshal walks the Name / Hero / Format header plus
// each "<N>x <card>" section line, routing entries into weapons / cards / sideboard /
// equipment depending on which section header was last seen. Unknown deck cards are returned
// to the caller as a skipped-count map so a silently-reduced deck surfaces to the user.

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// Unmarshal parses fabrary-style deck text and returns a *deck.Deck plus a count-keyed map of
// deck cards whose names aren't in the optimizer's registry. Callers should surface the skipped
// map so users aren't surprised by a silently-reduced deck. Stats aren't round-tripped.
//
// Arena-section entries split by lookup: weapon names land in d.Weapons, everything else lands
// in d.Equipment (the user-managed arena list) so the round-trip preserves the full loadout.
// A missing hero aborts: the deck can't be constructed without one.
func Unmarshal(text string) (*deck.Deck, map[string]int, error) {
	var (
		heroName  string
		section   string
		weapons   []weapon.Weapon
		cardList  []card.Card
		sideboard []string
		equipment []string
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
			if w, ok := weapon.ByName(name); ok {
				for i := 0; i < qty; i++ {
					weapons = append(weapons, w)
				}
				continue
			}
			// Non-weapon arena lines are equipment (head, chest, arms, legs) — stored as raw
			// names since the optimizer doesn't model them. The fabrary-case suffix isn't
			// stripped: equipment items don't carry pitch colors.
			for i := 0; i < qty; i++ {
				equipment = append(equipment, name)
			}
		case "deck":
			canon := fromFabraryCardName(name)
			id, ok := registry.CardByName(canon)
			if !ok {
				skipped[canon] += qty
				continue
			}
			c := registry.GetCard(id)
			for i := 0; i < qty; i++ {
				cardList = append(cardList, c)
			}
		case "sideboard":
			// Sideboard is a name-only list the sim doesn't touch, so there's no registry
			// lookup — any card or equipment piece the user lists comes back verbatim.
			// fromFabraryCardName maps the lowercase pitch suffix back to the canonical
			// "[R]" form; names without a recognized suffix (e.g. equipment pieces like
			// "Crown of Dichotomy") pass through unchanged.
			canon := fromFabraryCardName(name)
			for i := 0; i < qty; i++ {
				sideboard = append(sideboard, canon)
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
	d.Equipment = equipment
	return d, skipped, nil
}

// countedLine matches "<N>x <name>" — fabrary always uses a lowercase "x" with no surrounding
// spaces.
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
