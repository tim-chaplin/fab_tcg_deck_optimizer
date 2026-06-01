package textio

// Read-only loader for the upstream the-fab-cube card.csv — the third on-disk text format this
// package handles (alongside the deck JSON and fabrary encodings), used by the card-data tools
// (cmd/parsecarddb, cmd/cardaudit). It's the one place the CSV column mapping lives.

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// CardCSV mirrors the columns of card.csv (csvs/english/card.csv in the-fab-cube's source repo).
// One CSV row = one CardCSV (one printing — a card with three pitch colours is three rows).
type CardCSV struct {
	UniqueID                 string
	Name                     string
	Color                    string
	Pitch                    string
	Cost                     string
	Power                    string
	Defense                  string
	Health                   string
	Intelligence             string
	Arcane                   string
	Types                    string
	Traits                   string
	CardKeywords             string
	AbilitiesAndEffects      string
	AbilityAndEffectKeywords string
	GrantedKeywords          string
	RemovedKeywords          string
	InteractsWithKeywords    string
	FunctionalText           string
	TypeText                 string
	CardPlayedHorizontally   string
	BlitzLegal               string
	CCLegal                  string
	SilverAgeLegal           string
	CommonerLegal            string
	LLLegal                  string
}

// String pretty-prints a CardCSV, omitting blank fields. Implements fmt.Stringer.
func (c CardCSV) String() string {
	var b strings.Builder
	add := func(label, value string) {
		if value == "" {
			return
		}
		fmt.Fprintf(&b, "  %-12s %s\n", label+":", value)
	}
	fmt.Fprintf(&b, "%s", c.Name)
	if c.Color != "" {
		fmt.Fprintf(&b, " (%s)", c.Color)
	}
	b.WriteByte('\n')
	add("Types", c.Types)
	add("Traits", c.Traits)
	add("Pitch", c.Pitch)
	add("Cost", c.Cost)
	add("Power", c.Power)
	add("Defense", c.Defense)
	add("Health", c.Health)
	add("Intelligence", c.Intelligence)
	add("Arcane", c.Arcane)
	add("Keywords", c.CardKeywords)
	add("Text", c.FunctionalText)
	add("Type Text", c.TypeText)
	return b.String()
}

// cardCSVColumns maps CSV header names to the CardCSV field they populate. When a new column
// appears upstream, add the mapping here.
var cardCSVColumns = []struct {
	Header string
	Assign func(*CardCSV, string)
}{
	{"Unique ID", func(c *CardCSV, v string) { c.UniqueID = v }},
	{"Name", func(c *CardCSV, v string) { c.Name = v }},
	{"Color", func(c *CardCSV, v string) { c.Color = v }},
	{"Pitch", func(c *CardCSV, v string) { c.Pitch = v }},
	{"Cost", func(c *CardCSV, v string) { c.Cost = v }},
	{"Power", func(c *CardCSV, v string) { c.Power = v }},
	{"Defense", func(c *CardCSV, v string) { c.Defense = v }},
	{"Health", func(c *CardCSV, v string) { c.Health = v }},
	{"Intelligence", func(c *CardCSV, v string) { c.Intelligence = v }},
	{"Arcane", func(c *CardCSV, v string) { c.Arcane = v }},
	{"Types", func(c *CardCSV, v string) { c.Types = v }},
	{"Traits", func(c *CardCSV, v string) { c.Traits = v }},
	{"Card Keywords", func(c *CardCSV, v string) { c.CardKeywords = v }},
	{"Abilities and Effects", func(c *CardCSV, v string) { c.AbilitiesAndEffects = v }},
	{"Ability and Effect Keywords", func(c *CardCSV, v string) { c.AbilityAndEffectKeywords = v }},
	{"Granted Keywords", func(c *CardCSV, v string) { c.GrantedKeywords = v }},
	{"Removed Keywords", func(c *CardCSV, v string) { c.RemovedKeywords = v }},
	{"Interacts with Keywords", func(c *CardCSV, v string) { c.InteractsWithKeywords = v }},
	{"Functional Text", func(c *CardCSV, v string) { c.FunctionalText = v }},
	{"Type Text", func(c *CardCSV, v string) { c.TypeText = v }},
	{"Card Played Horizontally", func(c *CardCSV, v string) { c.CardPlayedHorizontally = v }},
	{"Blitz Legal", func(c *CardCSV, v string) { c.BlitzLegal = v }},
	{"CC Legal", func(c *CardCSV, v string) { c.CCLegal = v }},
	{"Silver Age Legal", func(c *CardCSV, v string) { c.SilverAgeLegal = v }},
	{"Commoner Legal", func(c *CardCSV, v string) { c.CommonerLegal = v }},
	{"LL Legal", func(c *CardCSV, v string) { c.LLLegal = v }},
}

// LoadCardCSV reads the tab-separated card.csv at path into CardCSV rows, mapping columns by
// header name so an upstream column reorder doesn't break parsing.
func LoadCardCSV(path string) ([]CardCSV, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	colIdx := map[string]int{}
	for i, h := range header {
		colIdx[h] = i
	}
	for _, col := range cardCSVColumns {
		if _, ok := colIdx[col.Header]; !ok {
			return nil, fmt.Errorf("csv missing column %q", col.Header)
		}
	}

	var cards []CardCSV
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		var c CardCSV
		for _, col := range cardCSVColumns {
			i := colIdx[col.Header]
			if i < len(rec) {
				col.Assign(&c, rec[i])
			}
		}
		cards = append(cards, c)
	}
	return cards, nil
}
