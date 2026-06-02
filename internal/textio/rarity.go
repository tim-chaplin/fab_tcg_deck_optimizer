package textio

// Loader for the upstream rarity data: card-printing.csv (one row per printing, carrying the
// per-printing Rarity shorthand) joined to rarity.csv (shorthand -> full name). Used to apply
// the Silver Age rarity rule — a card is legal only if it has a Basic / Common / Rare printing.

import (
	"encoding/csv"
	"os"
	"strings"
)

// silverAgeRarities are the rarities Silver Age allows.
var silverAgeRarities = map[string]bool{"Basic": true, "Common": true, "Rare": true}

// LoadRarityLegal builds a map from a card's Unique ID to whether it is Silver-Age-rarity-legal
// (has at least one Basic / Common / Rare printing). It joins card-printing.csv (per-printing
// Rarity, keyed to a card by "Card Unique ID") with rarity.csv (shorthand -> full name). Keying
// on Unique ID rather than name avoids upstream name-spelling drift between the CSVs.
func LoadRarityLegal(printingPath, rarityPath string) (map[string]bool, error) {
	shorthand := map[string]string{} // "C" -> "Common"
	if err := forEachTSVRow(rarityPath, func(row map[string]string) {
		shorthand[strings.TrimSpace(row["Shorthand"])] = strings.TrimSpace(row["Text"])
	}); err != nil {
		return nil, err
	}

	legal := map[string]bool{}
	err := forEachTSVRow(printingPath, func(row map[string]string) {
		id := strings.TrimSpace(row["Card Unique ID"])
		if id != "" && silverAgeRarities[shorthand[strings.TrimSpace(row["Rarity"])]] {
			legal[id] = true
		}
	})
	return legal, err
}

// forEachTSVRow reads a tab-separated file and calls fn once per data row with a header->value
// map. Columns are addressed by header name so an upstream reorder doesn't break callers.
func forEachTSVRow(path string, fn func(map[string]string)) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return err
	}
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		row := make(map[string]string, len(header))
		for i, h := range header {
			if i < len(rec) {
				row[h] = rec[i]
			}
		}
		fn(row)
	}
	return nil
}
