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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/cards"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hand"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/weapon"
)

func main() {
	subcommand, args, ok := extractSubcommand()
	if !ok {
		printSubcommands(os.Stdout)
		return
	}

	// Create mydecks/ up front so downstream WriteFile calls can't fail on a missing dir after
	// a long run. Done once here rather than per-subcommand since every subcommand either reads
	// or writes mydecks/.
	if err := os.MkdirAll(mydecks.Dir, 0o755); err != nil {
		die("mkdir %s: %v", mydecks.Dir, err)
	}

	switch subcommand {
	case "help":
		printSubcommands(os.Stdout)
	case "anneal":
		runAnnealCmd(args)
	case "eval":
		runEvalCmd(args)
	case "print":
		runPrintCmd(args)
	case "diff":
		runDiffCmd(args)
	case "import":
		runImport()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", subcommand)
		printSubcommands(os.Stderr)
		os.Exit(2)
	}
}

// extractSubcommand pulls os.Args[1] as the subcommand name and returns the remaining args for
// the subcommand's own flag.FlagSet to parse. Returns (_, _, false) when no subcommand is given
// or the first arg looks like a flag (bare `fabsim`, `fabsim -help`); the caller prints the
// subcommand list.
func extractSubcommand() (string, []string, bool) {
	if len(os.Args) < 2 {
		return "", nil, false
	}
	first := os.Args[1]
	if strings.HasPrefix(first, "-") {
		return "", nil, false
	}
	return first, os.Args[2:], true
}

// printSubcommands writes the one-liner catalogue shown when no subcommand is given. Flag
// details live behind `fabsim <subcommand> -help`, which each subcommand's own FlagSet renders.
func printSubcommands(w io.Writer) {
	fmt.Fprintln(w, "fabsim: Flesh and Blood goldfishing deck optimizer")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: fabsim <subcommand> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Subcommands:")
	fmt.Fprintln(w, "  anneal    Hill-climb (optionally simulated-annealing) on the saved deck until a local maximum")
	fmt.Fprintln(w, "  eval      Re-score the saved deck at -deep-shuffles without overwriting it (usage: fabsim eval <deck>)")
	fmt.Fprintln(w, "  print     Print the saved deck without simulating (usage: fabsim print <deck>)")
	fmt.Fprintln(w, "  import    Paste a fabrary.net deck into mydecks/<name>.json")
	fmt.Fprintln(w, "  diff      Print the card-count delta between two saved decks (usage: fabsim diff <deck1> <deck2>)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'fabsim <subcommand> -help' for flag details.")
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

// parseFlagsAnywhere parses args on fs while tolerating flags that appear before, after, or
// interleaved with positional arguments. Go's stdlib flag package stops at the first positional
// token, so `fabsim eval <deck> --incoming=100` would otherwise leave --incoming unparsed and
// treat it as a second positional. Every subcommand routes through this helper so flag order
// never matters to the user.
//
// The reorder is aware of each flag's bool-ness (via the optional IsBoolFlag() interface that
// the flag package uses): a non-bool `-name` token consumes the following arg as its value, a
// bool flag doesn't. Unknown flags are passed through untouched so fs.Parse can emit its usual
// "flag provided but not defined" error. A bare `--` is honored as the end-of-flags terminator,
// with everything after it treated as positional.
func parseFlagsAnywhere(fs *flag.FlagSet, args []string) error {
	var flagTokens, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}
		if len(a) < 2 || a[0] != '-' {
			positional = append(positional, a)
			continue
		}
		flagTokens = append(flagTokens, a)
		// -name=value is self-contained; only -name (no =) can consume the next token.
		if strings.Contains(a, "=") {
			continue
		}
		name := strings.TrimLeft(a, "-")
		f := fs.Lookup(name)
		if f == nil {
			// Unknown flag — let fs.Parse produce the canonical error with usage.
			continue
		}
		bf, ok := f.Value.(interface{ IsBoolFlag() bool })
		isBool := ok && bf.IsBoolFlag()
		if !isBool && i+1 < len(args) {
			i++
			flagTokens = append(flagTokens, args[i])
		}
	}
	// Insert `--` before the positional block so fs.Parse doesn't re-interpret positional
	// tokens that happen to start with `-`. That's necessary when the caller supplied an
	// explicit `--` terminator (whose tail ends up in positional) and also future-proofs
	// against deck names that begin with a dash.
	out := append(flagTokens, "--")
	out = append(out, positional...)
	return fs.Parse(out)
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

// resolveDeckPath is the positional-arg counterpart to anneal's -deck flag. Subcommands that
// always operate on an existing deck (eval, print, diff) accept the deck name as a positional
// arg and resolve it to mydecks/<name>.json via mydecks.Path.
func resolveDeckPath(name string) string {
	p, err := mydecks.Path(name)
	if err != nil {
		die("%v", err)
	}
	return p
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
