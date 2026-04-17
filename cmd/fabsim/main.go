// Command fabsim generates N random Viserai decks, evaluates each, and reports the best one.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// defaultDeckName is the deck fabsim reads from / writes to when -deck isn't supplied. Matches the
// historical "best_deck" muscle memory.
const defaultDeckName = "best_deck"

func main() {
	mode := flag.String("mode", "random", "run mode: random, iterate, eval, print, or import")
	numDecks := flag.Int("decks", 10000, "number of random decks to generate (phase 1)")
	shallowShuffles := flag.Int("shallow-shuffles", 10, "shuffles per deck in phase 1 (wide search)")
	topN := flag.Int("top-n", 100, "number of top decks to advance to phase 2")
	deepShuffles := flag.Int("deep-shuffles", 1000, "shuffles per deck in phase 2 (deep evaluation)")
	incoming := flag.Int("incoming", 4, "opponent damage per turn")
	deckSize := flag.Int("deck-size", 40, "number of cards per deck")
	maxCopies := flag.Int("max-copies", 2, "maximum copies of any single card printing per deck")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	deckName := flag.String("deck", defaultDeckName, "deck name; resolved to mydecks/<name>.json (\".json\" suffix optional). Ignored by -mode=import, which always prompts interactively.")
	flag.Parse()

	// Create mydecks/ up front so downstream WriteFile calls in the search loops can't fail on
	// a missing dir after a long run. Harmless if it already exists.
	if err := os.MkdirAll(mydecks.Dir, 0o755); err != nil {
		die("mkdir %s: %v", mydecks.Dir, err)
	}

	if *mode == "import" {
		runImport()
		return
	}

	outPath, err := mydecks.Path(*deckName)
	if err != nil {
		die("%v", err)
	}

	cfg := config{
		numDecks:        *numDecks,
		shallowShuffles: *shallowShuffles,
		topN:            *topN,
		deepShuffles:    *deepShuffles,
		incoming:        *incoming,
		deckSize:        *deckSize,
		maxCopies:       *maxCopies,
		seed:            *seed,
		outPath:         outPath,
	}

	switch *mode {
	case "random":
		runRandom(cfg)
	case "iterate":
		runIterate(cfg)
	case "eval":
		runEval(cfg)
	case "print":
		runPrint(cfg.outPath)
	default:
		die("unknown mode %q (want random, iterate, eval, print, or import)", *mode)
	}
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

type config struct {
	numDecks        int
	shallowShuffles int
	topN            int
	deepShuffles    int
	incoming        int
	deckSize        int
	maxCopies       int
	seed            int64
	outPath         string
}

// loadExisting reads and deserializes the deck at path. Returns (nil, 0) if the file doesn't
// exist or can't be parsed — the caller treats that as "no previous best".
func loadExisting(path string) (*deck.Deck, float64) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		return nil, 0
	}
	return d, d.Stats.Avg()
}

// writeDeck persists d as JSON at path plus a sibling fabrary-format .txt ("x.json" → "x.txt"),
// so the saved best deck is ready to paste into fabrary.net without a second export step.
func writeDeck(d *deck.Deck, path string) error {
	data, err := deckio.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	txtPath := fabraryPathFor(path)
	if err := os.WriteFile(txtPath, []byte(fabrary.Marshal(d)), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", txtPath, err)
	}
	return nil
}

// fabraryPathFor derives the sibling .txt path. A ".json" extension is replaced; anything else
// gets ".txt" appended rather than clobbered, so an unusual "-out deck.data" still yields a
// "deck.data.txt" sibling instead of overwriting the JSON.
func fabraryPathFor(jsonPath string) string {
	if ext := filepath.Ext(jsonPath); ext == ".json" {
		return strings.TrimSuffix(jsonPath, ext) + ".txt"
	}
	return jsonPath + ".txt"
}

func printBestDeck(d *deck.Deck) {
	s := d.Stats
	fmt.Printf("Best deck (avg %.3f over %d hands)\n", s.Avg(), s.Hands)
	fmt.Printf("  Hero:    %s\n", d.Hero.Name())
	fmt.Printf("  Weapons: %s\n", weaponNames(d.Weapons))
	fmt.Printf("  Cycle 1 avg: %.3f  (%d hands)\n", s.FirstCycle.Avg(), s.FirstCycle.Hands)
	fmt.Printf("  Cycle 2 avg: %.3f  (%d hands)\n", s.SecondCycle.Avg(), s.SecondCycle.Hands)
	var red, yellow, blue int
	for _, c := range d.Cards {
		switch c.Pitch() {
		case 1:
			red++
		case 2:
			yellow++
		case 3:
			blue++
		}
	}
	fmt.Printf("  Pitch:   %d red / %d yellow / %d blue\n", red, yellow, blue)
	if b := s.Best; b.Hand != nil {
		line := hand.FormatRoles(b.Hand, b.Play.Roles)
		for _, w := range b.Play.Weapons {
			line += ", " + w + ": ATTACK"
		}
		prefix := fmt.Sprintf("  Best hand seen (value %d", b.Play.Value)
		if d.Hero.Types().Has(card.TypeRuneblade) {
			prefix += fmt.Sprintf(", %d carryover runechants", b.StartingRunechants)
		}
		fmt.Printf("%s): %s\n", prefix, line)
	}
	fmt.Println()
	fmt.Println("Card list:")
	counts := map[string]int{}
	for _, c := range d.Cards {
		counts[c.Name()]++
	}
	names := make([]string, 0, len(counts))
	for n := range counts {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("  %dx %s\n", counts[n], n)
	}
}

func weaponNames(ws []weapon.Weapon) string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	return fmt.Sprintf("%v", names)
}
