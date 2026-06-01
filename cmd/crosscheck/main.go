// Command crosscheck compares our deck-construction pool against the official Flesh and
// Blood card database (cardvault.fabtcg.com), so we can spot cards the official site says are
// legal that we're missing — or cards in our pool the official site no longer lists.
//
// It pulls every Silver-Age-legal Runeblade and Generic print from the official advanced-
// search API, keeps the no-talent main-deck cards (the ones Viserai can legally run), and
// diffs that set against registry.LegalCardsFor(SilverAge, Viserai). It reports:
//
//   - [A] official-legal cards absent from our pool (annotated: unimplemented vs quarantined);
//   - [B] cards in our pool the official query doesn't list;
//   - [C] entries on our internal/format banlist the official API now lists as legal — stale
//     bans to remove (catches banlist drift automatically).
//
// Network: hits the public cardvault API; run from a machine with outbound HTTPS.
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabtype"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/format"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry"
)

const apiBase = "https://api.cardvault.fabtcg.com/carddb/api/v1/advanced-search/"

// viseraiClasses is Viserai's own hero class. Generic is legal for every hero, so it's handled
// by fabtype.ClassMatches rather than listed here; the talent / non-deck vocabulary also comes
// from fabtype.
var viseraiClasses = map[string]bool{"Runeblade": true}

type apiResponse struct {
	Next    string `json:"next"`
	Results []struct {
		PrintedName    string `json:"printed_name"`
		PrintedTypebox string `json:"printed_typebox"`
	} `json:"results"`
}

type print struct{ name, typebox string }

// fetchPrints paginates the advanced-search API for query q (adding page size / ordering),
// following next links.
func fetchPrints(client *http.Client, q url.Values) ([]print, error) {
	q.Set("page_size", "250")
	q.Set("orderby", "name")
	next := apiBase + "?" + q.Encode()

	var out []print
	for next != "" {
		req, err := http.NewRequest(http.MethodGet, next, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		var r apiResponse
		decErr := json.NewDecoder(resp.Body).Decode(&r)
		resp.Body.Close()
		if decErr != nil {
			return nil, decErr
		}
		for _, c := range r.Results {
			out = append(out, print{name: c.PrintedName, typebox: c.PrintedTypebox})
		}
		next = r.Next
		if next != "" {
			time.Sleep(200 * time.Millisecond)
		}
	}
	return out, nil
}

// typeboxSplit normalizes a printed typebox into its words. Subtype (" - ") and double-faced
// ("||", "//") separators become spaces so each face's class/talent words are checked — a DFC
// whose other face is talented (Burn Up // Shock) is then correctly ruled out.
var typeboxSplit = strings.NewReplacer(" - ", " ", "||", " ", "//", " ")

// dfcSeparator matches the double-faced-card join in either the API's "||" form or our
// banlist's " // " form, with any surrounding spaces.
var dfcSeparator = regexp.MustCompile(`\s*(\|\||//)\s*`)

// normalizeName canonicalizes a card name so the API's DFC form ("Burn Up||Shock") matches our
// banlist's ("Burn Up // Shock").
func normalizeName(s string) string {
	return dfcSeparator.ReplaceAllString(s, " // ")
}

// inScope reports whether a printed typebox is a no-talent Runeblade/Generic main-deck card —
// the cards Viserai can legally run. A card is out of scope if any of its words (across both
// faces) is a talent, a non-deck card-type, or a class other than Runeblade / Generic.
func inScope(typebox string) bool {
	for _, tok := range strings.Fields(typeboxSplit.Replace(typebox)) {
		if fabtype.TalentNames[tok] || fabtype.NonDeckTypes[tok] {
			return false
		}
		if !fabtype.ClassMatches(tok, viseraiClasses) {
			return false
		}
	}
	return true
}

// scanCardYAML walks the card yaml tree and maps each card's printed name to its location:
// "unplayable" / "notimplemented" for the quarantine subpackages (which the registry doesn't
// compile in, so registry.AllCards can't see them), or "main" for the registered cards. Lets
// the diff distinguish "we implemented this but quarantined it" from "we never made it".
func scanCardYAML(root string) (map[string]string, error) {
	loc := map[string]string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return nil
		}
		name, ok := firstYAMLName(path)
		if !ok {
			return nil
		}
		category := "main"
		switch {
		case strings.Contains(filepath.ToSlash(path), "/unplayable/"):
			category = "unplayable"
		case strings.Contains(filepath.ToSlash(path), "/notimplemented/"):
			category = "notimplemented"
		}
		loc[name] = category
		return nil
	})
	return loc, err
}

// firstYAMLName returns the value of the leading `name:` key in a card yaml file.
func firstYAMLName(path string) (string, bool) {
	f, err := os.Open(path)
	if err != nil {
		return "", false
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if rest, ok := strings.CutPrefix(line, "name:"); ok {
			return strings.TrimSpace(rest), true
		}
	}
	return "", false
}

// missingReason explains why an official-legal card is absent from our pool, combining the
// yaml location (quarantined subpackage vs registered) with the registry marker / ban status.
func missingReason(name string, yamlLoc, status map[string]string) string {
	switch yamlLoc[name] {
	case "unplayable":
		return "implemented, quarantined as Unplayable"
	case "notimplemented":
		return "implemented, quarantined as NotImplemented"
	case "main":
		if s, ok := status[name]; ok {
			return s
		}
		return "registered (status unknown)"
	default:
		return "NOT IMPLEMENTED (no card yaml)"
	}
}

// cardStatus explains why a registered card is (or isn't) in the legal pool.
func cardStatus(c card.Card) string {
	if _, ok := c.(registry.NotImplemented); ok {
		return "NotImplemented"
	}
	if _, ok := c.(registry.Unplayable); ok {
		return "Unplayable"
	}
	if !format.SilverAge.IsCardLegal(c) {
		return "banned by our banlist"
	}
	return "in pool"
}

func main() {
	cardsDir := flag.String("cards-dir", "internal/card/cards", "root of the card yaml tree, used to classify missing cards (quarantined vs never implemented)")
	flag.Parse()

	client := &http.Client{Timeout: 60 * time.Second}

	yamlLoc, err := scanCardYAML(*cardsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "scan %s: %v\n", *cardsDir, err)
		os.Exit(1)
	}

	// Viserai diff [A]/[B]: query per class so the API does the class filtering server-side —
	// far more reliable than filtering the full set, which is littered with Events, placeholder
	// cards, and empty typeboxes the inScope blacklist can't catch. officialInScope is the
	// no-talent main-deck subset; officialAllNames explains pool cards it omits.
	officialInScope := map[string]string{} // name -> a representative in-scope typebox
	officialAllNames := map[string]bool{}
	for _, cls := range []string{"Runeblade", "Generic"} {
		q := url.Values{}
		q.Set("classes", cls)
		q.Set("legal_formats", "Silver Age")
		prints, err := fetchPrints(client, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", cls, err)
			os.Exit(1)
		}
		for _, p := range prints {
			officialAllNames[p.name] = true
			if inScope(p.typebox) {
				officialInScope[p.name] = p.typebox
			}
		}
	}

	// Banlist check [C] needs every SA-legal name across all classes (the banlist spans them).
	q := url.Values{}
	q.Set("legal_formats", "Silver Age")
	allPrints, err := fetchPrints(client, q)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch all: %v\n", err)
		os.Exit(1)
	}
	legalNames := map[string]bool{} // normalized SA-legal names
	for _, p := range allPrints {
		legalNames[normalizeName(p.name)] = true
	}

	// Our side: the Viserai pool, plus the status of every registered card (by base name) so we
	// can say why an official card is missing from the pool.
	pool := registry.Registry{}.LegalCardsFor(format.SilverAge, heroes.Viserai)
	poolNames := map[string]bool{}
	for _, c := range pool {
		poolNames[c.Name()] = true
	}
	status := map[string]string{}
	for _, id := range registry.AllCards() {
		c := registry.GetCard(id)
		if _, seen := status[c.Name()]; !seen {
			status[c.Name()] = cardStatus(c)
		}
	}

	// [A] Official-legal but absent from our pool.
	var missing []string
	for name := range officialInScope {
		if !poolNames[name] {
			missing = append(missing, name)
		}
	}
	sort.Strings(missing)

	// [B] In our pool but not in the official in-scope set.
	var extra []string
	for name := range poolNames {
		if _, ok := officialInScope[name]; !ok {
			extra = append(extra, name)
		}
	}
	sort.Strings(extra)

	fmt.Printf("official SA Runeblade/Generic no-talent main-deck cards: %d\n", len(officialInScope))
	fmt.Printf("our LegalCardsFor(SilverAge, Viserai) distinct names:    %d\n\n", len(poolNames))

	fmt.Printf("=== [A] Official-legal but NOT in our pool (%d) ===\n", len(missing))
	tally := map[string]int{}
	for _, name := range missing {
		reason := missingReason(name, yamlLoc, status)
		tally[reason]++
		fmt.Printf("  %-34s %-36s — %s\n", name, "["+officialInScope[name]+"]", reason)
	}
	if len(missing) > 0 {
		fmt.Println("\n  by reason:")
		reasons := make([]string, 0, len(tally))
		for r := range tally {
			reasons = append(reasons, r)
		}
		sort.Strings(reasons)
		for _, r := range reasons {
			fmt.Printf("    %3d  %s\n", tally[r], r)
		}
	}

	fmt.Printf("\n=== [B] In our pool but NOT in official SA Runeblade/Generic no-talent set (%d) ===\n", len(extra))
	for _, name := range extra {
		why := "official query didn't return this name — likely a recent banlist change or a name mismatch"
		if officialAllNames[name] {
			why = "official returned it, but its typebox fell out of scope (talent / non-deck type)"
		}
		fmt.Printf("  %-34s — %s\n", name, why)
	}

	// [C] Our banlist entries the official API now lists as legal — stale bans to remove from
	// internal/format/banlist.go. (The reverse drift — a card we don't ban that's officially
	// banned — surfaces in [B] for cards that are in our pool.)
	var stale []string
	for _, name := range format.SilverAge.BannedNames() {
		if legalNames[normalizeName(name)] {
			stale = append(stale, name)
		}
	}
	fmt.Printf("\n=== [C] On our banlist but official API now lists as LEGAL — stale bans (%d) ===\n", len(stale))
	for _, name := range stale {
		fmt.Printf("  %s\n", name)
	}
	if len(stale) == 0 {
		fmt.Println("  (none — our Silver Age banlist matches the official API)")
	}
}
