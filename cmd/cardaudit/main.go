// Command cardaudit cross-checks the numeric stats of every implemented card's yaml against
// both the upstream card.csv and the official cardvault.fabtcg.com API, and prints the
// discrepancies. It is report-only — no files are changed. Use it to catch a yaml whose cost /
// power / defense was transcribed wrong, or where card.csv itself disagrees with the printed
// card.
//
// Matching is by (card name, pitch): each yaml variant is looked up in card.csv and the API by
// its printed name and pitch colour. Type lines and rules text aren't compared (their wording
// differs across sources); this is a stats pass.
//
// Network: hits the public cardvault API; run from a machine with outbound HTTPS.
package main

import (
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
	"strconv"
	"strings"
	"time"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cardcsv"
	"gopkg.in/yaml.v3"
)

const apiBase = "https://api.cardvault.fabtcg.com/carddb/api/v1/advanced-search/"

// yamlCard is the slice of a card yaml cardaudit reads: identity plus the per-pitch stats.
type yamlCard struct {
	Name     string `yaml:"name"`
	Cost     any    `yaml:"cost"` // int, "variable", or nil
	Variants []struct {
		ID      string `yaml:"id"`
		Pitch   int    `yaml:"pitch"`
		Attack  int    `yaml:"attack"`
		Defense int    `yaml:"defense"`
	} `yaml:"variants"`
	path string
}

// stats are the comparable numeric fields of one printing.
type stats struct{ cost, power, defense string }

// expected lists stats we deliberately model differently from the printed value, so the audit
// doesn't keep flagging an intentional modeling choice as an error. Keyed by "card name|field".
// Add an entry (with the reason) whenever a card encodes an effect into a stat.
var expected = map[string]string{
	"Drag Down|defense": "modeled as block — the printed effect debuffs the opponent's attack, which we encode as block value",
}

// dfcSeparator normalizes the double-faced-card join across sources ("Burn Up||Shock" vs
// "Burn Up // Shock").
var dfcSeparator = regexp.MustCompile(`\s*(\|\||//)\s*`)

func normName(s string) string { return dfcSeparator.ReplaceAllString(s, " // ") }

// loadYAML walks the card yaml tree and returns every parsed card (those with variants).
func loadYAML(root string) ([]yamlCard, error) {
	var out []yamlCard
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".yaml") {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var c yamlCard
		if err := yaml.Unmarshal(data, &c); err != nil {
			fmt.Fprintf(os.Stderr, "warn: parse %s: %v\n", path, err)
			return nil
		}
		if c.Name != "" && len(c.Variants) > 0 {
			c.path = path
			out = append(out, c)
		}
		return nil
	})
	return out, err
}

type apiResponse struct {
	Next    string `json:"next"`
	Results []struct {
		PrintedName string `json:"printed_name"`
		// The stat fields come back as a JSON string ("0", "") on some records and a number /
		// null on others, so decode them as any and normalize with asStr.
		PrintedPitch   any `json:"printed_pitch"`
		PrintedCost    any `json:"printed_cost"`
		PrintedPower   any `json:"printed_power"`
		PrintedDefense any `json:"printed_defense"`
	} `json:"results"`
}

// asStr normalizes a JSON value (string / number / null) to a trimmed string; a number renders
// as its integer form.
func asStr(v any) string {
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x)
	case float64:
		return strconv.Itoa(int(x))
	default:
		return ""
	}
}

// fetchAPIStats returns printed stats keyed by normalized "name|pitch", across the given
// classes (no legality filter, so banned implemented cards still match).
func fetchAPIStats(client *http.Client, classes []string) (map[string]stats, error) {
	out := map[string]stats{}
	for _, cls := range classes {
		q := url.Values{}
		q.Set("classes", cls)
		q.Set("page_size", "250")
		q.Set("orderby", "name")
		next := apiBase + "?" + q.Encode()
		for next != "" {
			req, _ := http.NewRequest(http.MethodGet, next, nil)
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
			for _, p := range r.Results {
				key := normName(p.PrintedName) + "|" + asStr(p.PrintedPitch)
				if _, seen := out[key]; !seen {
					out[key] = stats{cost: asStr(p.PrintedCost), power: asStr(p.PrintedPower), defense: asStr(p.PrintedDefense)}
				}
			}
			next = r.Next
			if next != "" {
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
	return out, nil
}

// loadCSVStats indexes card.csv stats by normalized "name|pitch".
func loadCSVStats(path string) (map[string]stats, error) {
	rows, err := cardcsv.Load(path)
	if err != nil {
		return nil, err
	}
	out := map[string]stats{}
	for _, r := range rows {
		key := normName(r.Name) + "|" + r.Pitch
		if _, seen := out[key]; !seen {
			out[key] = stats{cost: r.Cost, power: r.Power, defense: r.Defense}
		}
	}
	return out, nil
}

// parseInt parses a source stat string to an int. A blank / non-numeric value (an aura has no
// printed power, a variable-cost card has none) yields ok=false — nothing to compare.
func parseInt(s string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	return n, err == nil
}

// raw renders a source value for display, or "-" when the source has no row for the card.
func raw(s string, present bool) string {
	if !present || s == "" {
		return "-"
	}
	return s
}

func main() {
	cardsDir := flag.String("cards-dir", "internal/card/cards", "root of the card yaml tree")
	csvPath := flag.String("csv", "data_sources/card.csv", "path to card.csv")
	flag.Parse()

	cards, err := loadYAML(*cardsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load yaml: %v\n", err)
		os.Exit(1)
	}
	csvStats, err := loadCSVStats(*csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load csv: %v\n", err)
		os.Exit(1)
	}
	apiStats, err := fetchAPIStats(&http.Client{Timeout: 60 * time.Second}, []string{"Runeblade", "Generic"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch api: %v\n", err)
		os.Exit(1)
	}

	var bothDisagree, oneDisagrees, expectedDev, unmatched []string
	for _, c := range cards {
		yc, costIsInt := c.Cost.(int)
		for _, v := range c.Variants {
			key := normName(c.Name) + "|" + strconv.Itoa(v.Pitch)
			csv, inCSV := csvStats[key]
			api, inAPI := apiStats[key]
			label := fmt.Sprintf("%-32s pitch %d", c.Name, v.Pitch)

			if !inCSV || !inAPI {
				var miss []string
				if !inCSV {
					miss = append(miss, "card.csv")
				}
				if !inAPI {
					miss = append(miss, "API")
				}
				unmatched = append(unmatched, fmt.Sprintf("  %s  — not found in %s", label, strings.Join(miss, " + ")))
			}

			for _, chk := range []struct {
				field          string
				yaml           int
				skip           bool
				csvRaw, apiRaw string
			}{
				{"cost", yc, !costIsInt, csv.cost, api.cost},
				{"attack", v.Attack, false, csv.power, api.power},
				{"defense", v.Defense, false, csv.defense, api.defense},
			} {
				if chk.skip {
					continue
				}
				csvN, csvOk := parseInt(chk.csvRaw)
				apiN, apiOk := parseInt(chk.apiRaw)
				csvDiff := inCSV && csvOk && csvN != chk.yaml
				apiDiff := inAPI && apiOk && apiN != chk.yaml
				if !csvDiff && !apiDiff {
					continue
				}
				line := fmt.Sprintf("  %s  %-8s yaml=%d  csv=%s  api=%s", label, chk.field+":",
					chk.yaml, raw(chk.csvRaw, inCSV), raw(chk.apiRaw, inAPI))
				if reason, ok := expected[c.Name+"|"+chk.field]; ok {
					expectedDev = append(expectedDev, line+"  ("+reason+")")
				} else if csvDiff && apiDiff {
					bothDisagree = append(bothDisagree, line)
				} else {
					oneDisagrees = append(oneDisagrees, line)
				}
			}
		}
	}
	sort.Strings(bothDisagree)
	sort.Strings(oneDisagrees)
	sort.Strings(expectedDev)
	sort.Strings(unmatched)

	fmt.Printf("audited %d cards against card.csv (%d rows) and the cardvault API (%d prints)\n\n",
		len(cards), len(csvStats), len(apiStats))

	fmt.Printf("=== [A] yaml disagrees with BOTH sources — likely a yaml error (%d) ===\n", len(bothDisagree))
	printLines(bothDisagree, "  (none)")

	fmt.Printf("\n=== [B] yaml disagrees with one source only — check which is right (%d) ===\n", len(oneDisagrees))
	printLines(oneDisagrees, "  (none)")

	fmt.Printf("\n=== [C] expected deviations — intentional modeling choices, see the expected allowlist (%d) ===\n", len(expectedDev))
	printLines(expectedDev, "  (none)")

	fmt.Printf("\n=== [D] name+pitch not found in a source (name mismatch / DFC / new card) (%d) ===\n", len(unmatched))
	printLines(unmatched, "  (none — every card matched both sources)")
}

func printLines(lines []string, ifEmpty string) {
	if len(lines) == 0 {
		fmt.Println(ifEmpty)
		return
	}
	for _, l := range lines {
		fmt.Println(l)
	}
}
