// Command fabsim searches for, evaluates, and iterates on Flesh and Blood decks.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
	fmtpkg "github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

// defaultDeckNameFor returns the deck name when -deck isn't supplied, keyed by hero, format, and
// -incoming. Different regimes produce different optimal decks, so each gets its own file to
// avoid hill-climbing one regime's best under another regime's objective.
func defaultDeckNameFor(h hero.Hero, f fmtpkg.Format, incoming int) string {
	return fmt.Sprintf("%s_%s_%d_incoming", strings.ToLower(h.Name()), f, incoming)
}

func main() {
	subcommand, ok := extractSubcommand()
	if !ok {
		printSubcommands(os.Stdout)
		return
	}

	numDecks := flag.Int("decks", 1000, "number of random decks to generate (phase 1)")
	shallowShuffles := flag.Int("shallow-shuffles", 100, "shuffles per deck in phase 1 (wide search); also used to screen iterate mutations before deep confirmation")
	topN := flag.Int("top-n", 100, "number of top decks to advance to phase 2")
	deepShuffles := flag.Int("deep-shuffles", 10000, "shuffles per deck in phase 2 (deep evaluation); also used to confirm iterate improvements")
	incoming := flag.Int("incoming", 0, "opponent damage per turn")
	deckSize := flag.Int("deck-size", 40, "number of cards per deck")
	maxCopies := flag.Int("max-copies", 2, "maximum copies of any single card printing per deck")
	seed := flag.Int64("seed", time.Now().UnixNano(), "RNG seed")
	deckName := flag.String("deck", "", "deck name; resolved to mydecks/<name>.json (\".json\" suffix optional). Defaults to <hero>_<format>_<incoming>_incoming so different (hero, format, -incoming) regimes keep separate deck files. Ignored by the import subcommand, which always prompts interactively.")
	formatFlag := flag.String("format", string(fmtpkg.SilverAge), "constructed format whose banlist restricts the card pool during search (only \"silver_age\" is supported today)")
	debug := flag.Bool("debug", false, "emit extra diagnostic output (e.g. memo cache size between iterate rounds)")
	reevaluate := flag.Bool("reevaluate", false, "iterate: force re-evaluation of the loaded deck's baseline avg, even if its prior run count already matches -deep-shuffles. Use after adjusting modelling assumptions or fixing bugs that may have shifted the deck's true score.")
	finalize := flag.Bool("finalize", false, "iterate: high-precision pass — overrides -shallow-shuffles to 10000 and -deep-shuffles to 100000. Use on a deck that's already converged to squeeze out the remaining sub-percent improvements.")
	startTemp := flag.Float64("start-temp", 0, "iterate: simulated-annealing starting temperature. 0 (default) runs a pure hill climb. Higher values probabilistically accept worse mutations early; acceptance probability is exp((avg - baseline) / T). Good starting range is ~0.05–0.5 given typical Value units.")
	tempDecay := flag.Float64("temp-decay", 0.95, "iterate: multiplicative cooling per acceptance — T ← T × decay, floored at -min-temp. Unused when -start-temp is 0.")
	minTemp := flag.Float64("min-temp", 0, "iterate: minimum temperature. Once T reaches this floor the climb becomes greedy until a local maximum is found. 0 disables annealing in the converged tail.")
	flag.Parse()
	if *finalize {
		if subcommand != "iterate" {
			die("-finalize is only valid with the iterate subcommand")
		}
		*shallowShuffles = 10000
		*deepShuffles = 100000
	}
	// Reject positional args after the subcommand so `fabsim eval mydeck` errors instead of
	// silently ignoring the deck name. The diff subcommand consumes exactly two positional
	// deck names and is handled below, so skip the rejection for it.
	if subcommand != "diff" && flag.NArg() > 0 {
		die("unexpected positional argument(s): %v (did you mean -deck %s?)", flag.Args(), flag.Args()[0])
	}
	fmtValue, err := fmtpkg.Parse(*formatFlag)
	if err != nil {
		die("%v", err)
	}
	if *deckName == "" {
		*deckName = defaultDeckNameFor(hero.Viserai{}, fmtValue, *incoming)
	}

	// Create mydecks/ up front so downstream WriteFile calls can't fail on a missing dir after
	// a long run.
	if err := os.MkdirAll(mydecks.Dir, 0o755); err != nil {
		die("mkdir %s: %v", mydecks.Dir, err)
	}

	switch subcommand {
	case "help":
		printSubcommands(os.Stdout)
		return
	case "import":
		runImport()
		return
	case "diff":
		args := flag.Args()
		if len(args) != 2 {
			die("diff: need exactly 2 positional deck names (got %d)", len(args))
		}
		runDiff(args[0], args[1])
		return
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
		format:          fmtValue,
		debug:           *debug,
		reevaluate:      *reevaluate,
		startTemp:       *startTemp,
		tempDecay:       *tempDecay,
		minTemp:         *minTemp,
	}

	outPath, err := mydecks.Path(*deckName)
	if err != nil {
		die("%v", err)
	}
	cfg.outPath = outPath

	switch subcommand {
	case "random":
		runRandom(cfg)
	case "iterate":
		// Print the session-level delta (starting best vs final best) on any exit path, then
		// surface abort via a non-zero exit so wrapper scripts (iterate-reanneal.ps1 et al.)
		// can tell Enter-initiated termination from natural convergence and stop looping.
		res := runIterate(cfg)
		fmt.Fprintf(os.Stderr, "\nSession summary: avg %.3f → %.3f (%+.3f)\n",
			res.startingAvg, res.bestEverAvg, res.bestEverAvg-res.startingAvg)
		if res.aborted {
			os.Exit(130)
		}
	case "eval":
		runEval(cfg)
	case "print":
		runPrint(cfg.outPath)
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", subcommand)
		printSubcommands(os.Stderr)
		os.Exit(2)
	}
}

// extractSubcommand pulls os.Args[1] as the subcommand name and rewrites os.Args so flag.Parse
// only sees flags. Returns (_, false) when no subcommand is given or the first arg looks like
// a flag (bare `fabsim`, `fabsim -help`); the caller prints the subcommand list.
func extractSubcommand() (string, bool) {
	if len(os.Args) < 2 {
		return "", false
	}
	first := os.Args[1]
	if strings.HasPrefix(first, "-") {
		return "", false
	}
	os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
	return first, true
}

// printSubcommands writes the one-liner catalogue shown when no subcommand is given. Flag
// details live behind `fabsim <subcommand> -help`.
func printSubcommands(w io.Writer) {
	fmt.Fprintln(w, "fabsim: Flesh and Blood goldfishing deck optimizer")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: fabsim <subcommand> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Subcommands:")
	fmt.Fprintln(w, "  random    Search for a new best deck via two-phase random sampling")
	fmt.Fprintln(w, "  iterate   Hill-climb on the saved deck until a local maximum")
	fmt.Fprintln(w, "  eval      Re-score the saved deck at -deep-shuffles without overwriting it")
	fmt.Fprintln(w, "  print     Print the saved deck without simulating")
	fmt.Fprintln(w, "  import    Paste a fabrary.net deck into mydecks/<name>.json")
	fmt.Fprintln(w, "  diff      Print the card-count delta between two saved decks (usage: fabsim diff <deck1> <deck2>)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'fabsim <subcommand> -help' for flag details.")
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
	format          fmtpkg.Format
	debug           bool
	reevaluate      bool
	// startTemp / tempDecay / minTemp are the simulated-annealing knobs for iterate. startTemp
	// of 0 degenerates to the classical hill-climb (strict > baseline acceptance).
	startTemp float64
	tempDecay float64
	minTemp   float64
}

// legalFilter returns the card-pool predicate for this run's format. fabsim always runs under
// a format, so this is non-nil; the deck package accepts nil for "no filtering" generally.
func (c config) legalFilter() func(card.Card) bool {
	return c.format.IsLegal
}

// loadExisting reads and deserializes the deck at path. Returns (nil, 0) on missing or
// unparsable file — the caller treats that as "no previous best".
func loadExisting(path string) (*deck.Deck, float64) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, 0
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		return nil, 0
	}
	return d, d.Stats.Mean()
}

// writeDeck persists d as JSON at path plus a sibling fabrary-format .txt ("x.json" → "x.txt")
// so the saved deck is ready to paste into fabrary.net without a second export step.
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

// fabraryPathFor derives the sibling .txt path. A ".json" extension is swapped for ".txt";
// anything else gets ".txt" appended so non-JSON paths can't be overwritten.
func fabraryPathFor(jsonPath string) string {
	if ext := filepath.Ext(jsonPath); ext == ".json" {
		return strings.TrimSuffix(jsonPath, ext) + ".txt"
	}
	return jsonPath + ".txt"
}

// printCardList writes the deck's card list in canonical "Card list:" form: one grouped-and-
// sorted count-and-name line per unique card. Shared between printBestDeck and iterate's
// starting-deck banner so both callers render decks the same way.
func printCardList(d *deck.Deck) {
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

func printBestDeck(d *deck.Deck) {
	s := d.Stats
	fmt.Printf("Best deck (min %d, median %.1f, mean %.3f, max %d over %d hands)\n",
		s.Min(), s.Median(), s.Mean(), s.Max(), s.Hands)
	fmt.Printf("  Hero:    %s\n", d.Hero.Name())
	fmt.Printf("  Weapons: %s\n", weaponNames(d.Weapons))
	fmt.Printf("  Cycle 1 mean: %.3f  (%d hands)\n", s.FirstCycle.Mean(), s.FirstCycle.Hands)
	fmt.Printf("  Cycle 2 mean: %.3f  (%d hands)\n", s.SecondCycle.Mean(), s.SecondCycle.Hands)
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
	fmt.Println()
	printCardList(d)

	if b := s.Best; len(b.Summary.BestLine) > 0 {
		fmt.Println()
		header := fmt.Sprintf("Best turn played (value %d", b.Summary.Value)
		if d.Hero.Types().Has(card.TypeRuneblade) {
			header += fmt.Sprintf(", %d carryover runechants", b.StartingRunechants)
		}
		header += "):"
		fmt.Println(header)
		fmt.Println(hand.FormatBestTurn(b.Summary))
	}

	if len(s.PerCard) > 0 {
		printPerCardStats(d)
	}
}

// printPerCardStats renders per-card averages collected by deck.Evaluate: mean per-card
// contribution across hands the card appeared in. Contribution is role-based (attack power on
// attacks, proportional prevented-damage share on defends, Pitch on pitches), so the ranking
// reflects what each card typically does in its hand rather than the hand's total value.
func printPerCardStats(d *deck.Deck) {
	type row struct {
		name           string
		deckCount      int
		plays, pitches int
		avg            float64
	}
	deckCounts := map[card.ID]int{}
	for _, c := range d.Cards {
		deckCounts[c.ID()]++
	}
	rows := make([]row, 0, len(d.Stats.PerCard))
	for id, s := range d.Stats.PerCard {
		rows = append(rows, row{
			name:      cards.Get(id).Name(),
			deckCount: deckCounts[id],
			plays:     s.Plays,
			pitches:   s.Pitches,
			avg:       s.Avg(),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].avg != rows[j].avg {
			return rows[i].avg > rows[j].avg
		}
		ni, nj := rows[i].plays+rows[i].pitches, rows[j].plays+rows[j].pitches
		if ni != nj {
			return ni > nj
		}
		return rows[i].name < rows[j].name
	})

	fmt.Println()
	fmt.Println("Card value (avg contribution per appearance: attack=power, defend=share of block, pitch=resource):")
	for _, r := range rows {
		fmt.Printf("  %-*s avg %6.3f over %4d hands (%4d plays, %4d pitches, %dx in deck)\n",
			maxNameLen(d.Cards), r.name, r.avg, r.plays+r.pitches, r.plays, r.pitches, r.deckCount)
	}
}

// maxNameLen returns the length of the longest Name() across the given cards, or 0 when empty.
// Used to width fixed-width card-name columns in printed tables.
func maxNameLen(cs []card.Card) int {
	m := 0
	for _, c := range cs {
		if n := len(c.Name()); n > m {
			m = n
		}
	}
	return m
}

func weaponNames(ws []weapon.Weapon) string {
	names := make([]string, len(ws))
	for i, w := range ws {
		names[i] = w.Name()
	}
	return fmt.Sprintf("%v", names)
}
