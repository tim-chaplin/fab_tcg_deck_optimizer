// Command parsecarddb parses the flesh-and-blood-cards card.csv and prints matching cards in
// either pretty or JSON form. Its class / talent flags express a hero's legal pool (see
// poolFilter), so it can enumerate, say, Aurora's Lightning cards or Viserai's no-talent
// Runeblade/Generic pool.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cardcsv"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabtype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
)

// named adapts a card name to the format.Card interface (just Name()) so the banlist check can
// run against the CSV's Name string.
type named string

func (n named) Name() string { return string(n) }

// poolFilter is the set of predicates a card must satisfy to appear in the output. The class /
// talent fields express a hero's legal pool the way registry.heroCanPlay does (Generic and
// classless cards are always class-legal; a card's talents must be a subset of the hero's), so
// `-classes Runeblade -talents Lightning` enumerates Aurora's pool, and the zero value (no
// class / talent restriction, no-talent-only) is a plain Silver-Age card dump.
type poolFilter struct {
	nameNeedle    string          // substring match on Name (lowercased), "" = any
	typeNeedle    string          // substring match on the raw Types column (lowercased), "" = any
	classes       map[string]bool // allowed non-Generic classes; empty = any class
	talents       map[string]bool // allowed talents (card talents must be a subset)
	anyTalent     bool            // ignore the talent subset check entirely
	require       map[string]bool // card must carry at least one of these talents (empty = none)
	modeled       bool            // keep only card types the optimizer models (deck cards + weapons)
	silverAge     bool            // require Silver Age legality (the card.csv column)
	excludeBanned bool            // drop cards on our internal/format Silver Age banlist
}

// typeWords splits a card.csv Types cell ("Runeblade, Action, Attack") into its words.
func typeWords(types string) []string {
	var out []string
	for _, t := range strings.Split(types, ",") {
		if t = strings.TrimSpace(t); t != "" {
			out = append(out, t)
		}
	}
	return out
}

// matches reports whether c passes every predicate in f.
func (f poolFilter) matches(c cardcsv.Card) bool {
	if f.nameNeedle != "" && !strings.Contains(strings.ToLower(c.Name), f.nameNeedle) {
		return false
	}
	if f.typeNeedle != "" && !strings.Contains(strings.ToLower(c.Types), f.typeNeedle) {
		return false
	}
	// Blank Silver Age column is treated as legal — some rows are missing it.
	if f.silverAge && c.SilverAgeLegal != "Yes" && c.SilverAgeLegal != "" {
		return false
	}
	if f.excludeBanned && !format.SilverAge.IsCardLegal(named(c.Name)) {
		return false
	}

	words := typeWords(c.Types)
	gotTalent := false
	for _, w := range words {
		// The optimizer models deck cards and weapons but not equipment / tokens / etc.
		if f.modeled && fabtype.UnmodeledTypes[w] {
			return false
		}
		// Class: a class word must be one the hero plays (Generic is always legal — built into
		// fabtype.ClassMatches, so -classes lists only the hero's own class).
		if len(f.classes) > 0 && !fabtype.ClassMatches(w, f.classes) {
			return false
		}
		// Talent: every talent word must be one the hero shares (unless -any-talent).
		if fabtype.TalentNames[w] {
			if !f.anyTalent && !f.talents[w] {
				return false
			}
			if f.require[w] {
				gotTalent = true
			}
		}
	}
	if len(f.require) > 0 && !gotTalent {
		return false
	}
	return true
}

// toSet splits a comma-separated flag value into a set, dropping blanks and surrounding space.
func toSet(csv string) map[string]bool {
	out := map[string]bool{}
	for _, s := range strings.Split(csv, ",") {
		if s = strings.TrimSpace(s); s != "" {
			out[s] = true
		}
	}
	return out
}

func main() {
	in := flag.String("in", "data_sources/card.csv", "path to card.csv")
	nameFilter := flag.String("name", "", "only print cards whose name contains this substring (case insensitive)")
	typeFilter := flag.String("type", "", "only print cards whose Types field contains this substring (case insensitive), e.g. 'Aura'")
	classes := flag.String("classes", "", "comma-separated hero classes; a card must be Generic, classless, or one of these (e.g. 'Runeblade'). Empty = any class.")
	talents := flag.String("talents", "", "comma-separated allowed talents; a card's talents must be a subset, so 'Lightning' admits no-talent and Lightning cards. Empty = no-talent cards only.")
	anyTalent := flag.Bool("any-talent", false, "ignore -talents and admit cards carrying any talent")
	requireTalents := flag.String("require-talents", "", "comma-separated talents; admit only cards carrying at least one (e.g. 'Lightning' to list just the Lightning cards in a pool)")
	modeled := flag.Bool("modeled", false, "keep only card types the optimizer models — deck cards and weapons — dropping equipment, heroes, and tokens")
	silverAge := flag.Bool("silver-age", true, "require Silver Age legality (the card.csv 'Silver Age Legal' column)")
	excludeBanned := flag.Bool("exclude-banned", false, "drop cards on our internal/format Silver Age banlist (the cards not worth implementing)")
	outFormat := flag.String("format", "pretty", "output format: pretty | json")
	namesOnly := flag.Bool("names_only", false, "print only unique card names, one per line")
	flag.Parse()

	switch *outFormat {
	case "pretty", "json":
	default:
		log.Fatalf("unknown --format %q (want: pretty | json)", *outFormat)
	}

	filter := poolFilter{
		nameNeedle:    strings.ToLower(*nameFilter),
		typeNeedle:    strings.ToLower(*typeFilter),
		classes:       toSet(*classes),
		talents:       toSet(*talents),
		anyTalent:     *anyTalent,
		require:       toSet(*requireTalents),
		modeled:       *modeled,
		silverAge:     *silverAge,
		excludeBanned: *excludeBanned,
	}

	cards, err := cardcsv.Load(*in)
	if err != nil {
		log.Fatal(err)
	}

	var matched []cardcsv.Card
	for _, c := range cards {
		if filter.matches(c) {
			matched = append(matched, c)
		}
	}

	if *namesOnly {
		seen := map[string]bool{}
		for _, c := range matched {
			if seen[c.Name] {
				continue
			}
			seen[c.Name] = true
			fmt.Println(c.Name)
		}
		return
	}

	switch *outFormat {
	case "pretty":
		for _, c := range matched {
			fmt.Printf("%v\n", c)
		}
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(matched); err != nil {
			log.Fatal(err)
		}
	}
}
