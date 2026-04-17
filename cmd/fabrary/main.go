// Command fabrary converts decks between the optimizer's JSON format and fabrary.net's plain-text
// format (the text you get from Export → Plain Text on fabrary, or paste into its Import tab).
//
// Export a deck the optimizer produced:
//
//	fabrary -mode export -in best_deck.json -out best_deck.txt
//
// Import a deck from fabrary.net (paste the "Plain Text" block into a file first):
//
//	fabrary -mode import -in fabrary_paste.txt -out imported_deck.json
//
// When -out is "-" (or omitted), output is written to stdout.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
)

func main() {
	mode := flag.String("mode", "export", "conversion direction: export (json→fabrary) or import (fabrary→json)")
	inPath := flag.String("in", "", "input file path (required)")
	outPath := flag.String("out", "-", "output file path; \"-\" writes to stdout")
	flag.Parse()

	if *inPath == "" {
		fmt.Fprintln(os.Stderr, "fabrary: -in is required")
		flag.Usage()
		os.Exit(2)
	}

	data, err := os.ReadFile(*inPath)
	if err != nil {
		die("read %s: %v", *inPath, err)
	}

	var out []byte
	switch *mode {
	case "export":
		d, err := deckio.Unmarshal(data)
		if err != nil {
			die("parse deck JSON: %v", err)
		}
		out = []byte(fabrary.Marshal(d))
	case "import":
		d, err := fabrary.Unmarshal(string(data))
		if err != nil {
			die("parse fabrary text: %v", err)
		}
		b, err := deckio.Marshal(d)
		if err != nil {
			die("encode deck JSON: %v", err)
		}
		out = append(b, '\n')
	default:
		die("unknown mode %q (want export or import)", *mode)
	}

	if err := writeOut(*outPath, out); err != nil {
		die("write %s: %v", *outPath, err)
	}
}

func writeOut(path string, data []byte) error {
	if path == "-" || path == "" {
		_, err := os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "fabrary: "+format+"\n", args...)
	os.Exit(1)
}
