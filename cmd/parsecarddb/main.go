// Command parsecarddb parses the flesh-and-blood-cards card.csv and prints matching cards in either
// pretty or JSON form.
package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Card mirrors the columns of card.csv (csvs/english/card.csv in
// the-fab-cube/flesh-and-blood-cards). One row = one Card.
type Card struct {
	UniqueID                string
	Name                    string
	Color                   string
	Pitch                   string
	Cost                    string
	Power                   string
	Defense                 string
	Health                  string
	Intelligence            string
	Arcane                  string
	Types                   string
	Traits                  string
	CardKeywords            string
	AbilitiesAndEffects     string
	AbilityAndEffectKeywords string
	GrantedKeywords         string
	RemovedKeywords         string
	InteractsWithKeywords   string
	FunctionalText          string
	TypeText                string
	CardPlayedHorizontally  string
	BlitzLegal              string
	CCLegal                 string
	SilverAgeLegal          string
	CommonerLegal           string
	LLLegal                 string
}

// String pretty-prints a Card, omitting blank fields. Implements fmt.Stringer so fmt.Printf("%v",
// c) uses this format automatically.
func (c Card) String() string {
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
	// add("Blitz", c.BlitzLegal)
	// add("CC", c.CCLegal)
	// add("Silver Age", c.SilverAgeLegal)
	// add("Commoner", c.CommonerLegal)
	// add("LL", c.LLLegal)
	return b.String()
}

// cardCSVColumns maps CSV header names to the Card field they populate. Keeping this adjacent to
// Card makes it obvious when a new column is added upstream — the compiler will complain if the
// field is missing.
var cardCSVColumns = []struct {
	Header string
	Assign func(*Card, string)
}{
	{"Unique ID", func(c *Card, v string) { c.UniqueID = v }},
	{"Name", func(c *Card, v string) { c.Name = v }},
	{"Color", func(c *Card, v string) { c.Color = v }},
	{"Pitch", func(c *Card, v string) { c.Pitch = v }},
	{"Cost", func(c *Card, v string) { c.Cost = v }},
	{"Power", func(c *Card, v string) { c.Power = v }},
	{"Defense", func(c *Card, v string) { c.Defense = v }},
	{"Health", func(c *Card, v string) { c.Health = v }},
	{"Intelligence", func(c *Card, v string) { c.Intelligence = v }},
	{"Arcane", func(c *Card, v string) { c.Arcane = v }},
	{"Types", func(c *Card, v string) { c.Types = v }},
	{"Traits", func(c *Card, v string) { c.Traits = v }},
	{"Card Keywords", func(c *Card, v string) { c.CardKeywords = v }},
	{"Abilities and Effects", func(c *Card, v string) { c.AbilitiesAndEffects = v }},
	{"Ability and Effect Keywords", func(c *Card, v string) { c.AbilityAndEffectKeywords = v }},
	{"Granted Keywords", func(c *Card, v string) { c.GrantedKeywords = v }},
	{"Removed Keywords", func(c *Card, v string) { c.RemovedKeywords = v }},
	{"Interacts with Keywords", func(c *Card, v string) { c.InteractsWithKeywords = v }},
	{"Functional Text", func(c *Card, v string) { c.FunctionalText = v }},
	{"Type Text", func(c *Card, v string) { c.TypeText = v }},
	{"Card Played Horizontally", func(c *Card, v string) { c.CardPlayedHorizontally = v }},
	{"Blitz Legal", func(c *Card, v string) { c.BlitzLegal = v }},
	{"CC Legal", func(c *Card, v string) { c.CCLegal = v }},
	{"Silver Age Legal", func(c *Card, v string) { c.SilverAgeLegal = v }},
	{"Commoner Legal", func(c *Card, v string) { c.CommonerLegal = v }},
	{"LL Legal", func(c *Card, v string) { c.LLLegal = v }},
}

func loadCards(path string) ([]Card, error) {
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

	var cards []Card
	for {
		rec, err := r.Read()
		if err != nil {
			break
		}
		var c Card
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

func main() {
	in := flag.String("in", "data_sources/card.csv", "path to card.csv")
	nameFilter := flag.String("name", "", "only print cards whose name contains this substring (case insensitive)")
	typeFilter := flag.String("type", "", "only print cards whose Types field contains this substring (case insensitive), e.g. 'Aura'")
	format := flag.String("format", "pretty", "output format: pretty | json")
	namesOnly := flag.Bool("names_only", false, "print only unique card names, one per line")
	flag.Parse()

	nameNeedle := strings.ToLower(*nameFilter)
	typeNeedle := strings.ToLower(*typeFilter)
	switch *format {
	case "pretty", "json":
	default:
		log.Fatalf("unknown --format %q (want: pretty | json)", *format)
	}

	cards, err := loadCards(*in)
	if err != nil {
		log.Fatal(err)
	}

	var matched []Card
	for _, c := range cards {
		if nameNeedle != "" && !strings.Contains(strings.ToLower(c.Name), nameNeedle) {
			continue
		}
		if typeNeedle != "" && !strings.Contains(strings.ToLower(c.Types), typeNeedle) {
			continue
		}
		if c.SilverAgeLegal != "Yes" && c.SilverAgeLegal != "" {
			continue
		}
		if !strings.Contains(c.Types, "Runeblade") && !strings.Contains(c.Types, "Generic") {
			continue
		}
		if strings.Contains(c.Types, "Shadow") ||
			strings.Contains(c.Types, "Elemental") ||
			strings.Contains(c.Types, "Lightning") ||
			strings.Contains(c.Types, "Earth") ||
			strings.Contains(c.Types, "Token") ||
			strings.Contains(c.Types, "Equipment") {
			continue
		}
		matched = append(matched, c)
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

	switch *format {
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
