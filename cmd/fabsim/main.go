// Command fabsim searches for, evaluates, and iterates on Flesh and Blood decks.
package main

import (
	"errors"
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
	// a long run. Every subcommand reads or writes this directory.
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

// extractSubcommand pulls os.Args[1] as the subcommand name and returns the remaining args
// for the subcommand's own flag.FlagSet. Returns (_, _, false) when no subcommand is given
// or the first arg looks like a flag; the caller prints the subcommand list.
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

// requireFlag dies with a usage error when fs.Parse didn't encounter -name. The flag's own Usage
// string is echoed so the per-flag guidance the caller wrote in the FlagSet (e.g. why the flag
// can't default) shows up alongside the "required" message — no need to duplicate that wording
// here.
func requireFlag(fs *flag.FlagSet, subcommand, name string) {
	seen := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			seen = true
		}
	})
	if seen {
		return
	}
	f := fs.Lookup(name)
	if f == nil {
		die("%s: internal error — required flag -%s is not registered on the FlagSet", subcommand, name)
	}
	die("%s: -%s is required\n  usage: %s", subcommand, name, f.Usage)
}

// parseFlagsAnywhere parses args on fs while tolerating flags that appear before, after, or
// interleaved with positional arguments. Go's stdlib flag package stops at the first
// positional token; every subcommand routes through this helper so flag order never matters
// to the user.
//
// Bool-awareness matters: a non-bool `-name` token consumes the following arg as its value,
// a bool flag (detected via IsBoolFlag) doesn't. Unknown flags pass through untouched so
// fs.Parse emits the canonical "flag provided but not defined" error. A bare `--` acts as
// the end-of-flags terminator, with everything after it treated as positional.
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

// loadExisting reads and deserializes the deck at path. Returns (nil, 0, nil) when the file
// doesn't exist — the caller treats that as "no previous best, generate a fresh deck."
// Returns (nil, 0, err) when the file exists but can't be read or parsed: callers must NOT
// treat that as "missing" or they'd silently overwrite a corrupt file with a random deck
// (looping wrapper scripts would clobber a converged deck after a Ctrl-C mid-write).
func loadExisting(path string) (*deck.Deck, float64, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, 0, nil
		}
		return nil, 0, fmt.Errorf("read %s: %w", path, err)
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		return nil, 0, fmt.Errorf("parse %s: %w (file exists but isn't a valid deck — "+
			"refusing to silently overwrite; inspect the file and delete it manually if you "+
			"want a fresh start)", path, err)
	}
	return d, d.Stats.Mean(), nil
}

// writeDeck persists d as JSON at path plus a sibling fabrary-format .txt ("x.json" →
// "x.txt") so the saved deck is ready to paste into fabrary.net without a second export step.
//
// Both files are written atomically via writeFileAtomic: data lands in <path>.tmp first,
// then os.Rename swaps it into place, so a Ctrl-C mid-write can never leave the destination
// empty or partially written.
func writeDeck(d *deck.Deck, path string) error {
	data, err := deckio.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if err := writeFileAtomic(path, data); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	txtPath := fabraryPathFor(path)
	if err := writeFileAtomic(txtPath, []byte(fabrary.Marshal(d))); err != nil {
		return fmt.Errorf("write %s: %w", txtPath, err)
	}
	return nil
}

// writeFileAtomic writes data to a temp file in the same directory as path and renames it
// over path. The same-directory placement keeps the rename within one filesystem so it stays
// atomic. Removes the temp file on any error so a failed write doesn't leave junk behind.
func writeFileAtomic(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	// Clean up the temp file on any failure path so crashed writes don't litter mydecks/ with
	// .tmp-* files. The rename success path makes this a no-op.
	defer os.Remove(tmpName)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

// fabraryPathFor derives the sibling .txt path. A ".json" extension is swapped for ".txt";
// anything else gets ".txt" appended so non-JSON paths can't be overwritten.
func fabraryPathFor(jsonPath string) string {
	if ext := filepath.Ext(jsonPath); ext == ".json" {
		return strings.TrimSuffix(jsonPath, ext) + ".txt"
	}
	return jsonPath + ".txt"
}

// mustLoadDeck loads the deck at path or dies. For subcommands that always operate on an
// existing deck (eval, print, diff), both "missing" and "corrupt" are fatal. anneal handles
// the distinction itself: "missing" is a valid input ("no deck yet, generate one") while
// "corrupt" needs the loud refusal to overwrite.
func mustLoadDeck(path string) *deck.Deck {
	d, _, err := loadExisting(path)
	if err != nil {
		die("%v", err)
	}
	if d == nil {
		die("could not load deck from %s (file not found)", path)
	}
	return d
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

// printCardList writes the deck's card list in canonical "Card list:" form: one
// grouped-and-sorted count-and-name line per unique card.
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

// printDeckSummary prints the compact score header: min/median/mean/max, hero, weapons, per-cycle
// means, and pitch colour counts. Separated from printBestDeck so eval can emit only this block —
// eval's whole purpose is the freshly-computed score, and on a small terminal the full card list
// scrolls that score off the top.
func printDeckSummary(d *deck.Deck) {
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
}

func printBestDeck(d *deck.Deck) {
	printDeckSummary(d)
	fmt.Println()
	printCardList(d)

	s := d.Stats
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

// maxNameLen returns the length of the longest Name() across cs, or 0 when empty. Used to
// size fixed-width card-name columns in printed tables.
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
