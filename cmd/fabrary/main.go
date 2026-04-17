// Command fabrary converts decks between the optimizer's JSON format and fabrary.net's plain-text
// format (the text you get from Export → Plain Text on fabrary, or paste into its Import tab).
//
// Export a deck the optimizer produced:
//
//	fabrary -mode export -in best_deck.json -out best_deck.txt
//
// Import a deck from fabrary.net — either pipe a file in or paste interactively:
//
//	fabrary -mode import -in fabrary_paste.txt -out imported_deck.json
//	fabrary -mode import -out imported_deck.json   # paste into stdin, end with Ctrl-Z↵ / Ctrl-D
//
// When -in is empty or "-", stdin is read. When -out is empty or "-", stdout is written.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
)

func main() {
	mode := flag.String("mode", "export", "conversion direction: export (json→fabrary) or import (fabrary→json)")
	inPath := flag.String("in", "", "input file path; \"-\" or empty reads from stdin")
	outPath := flag.String("out", "-", "output file path; \"-\" writes to stdout")
	flag.Parse()

	data, err := readIn(*inPath, *mode)
	if err != nil {
		die("read input: %v", err)
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

// readIn returns the input bytes for the chosen mode. An empty or "-" path means stdin; when
// stdin is a terminal (interactive paste) a short hint is printed to stderr so the user knows how
// to signal end-of-input.
func readIn(path, mode string) ([]byte, error) {
	if path != "" && path != "-" {
		return os.ReadFile(path)
	}
	if isTerminal(os.Stdin) {
		promptPaste(mode)
	}
	return io.ReadAll(os.Stdin)
}

func promptPaste(mode string) {
	what := "fabrary deck text"
	if mode == "export" {
		what = "deck JSON"
	}
	// endKey differs by platform — the Windows console treats Ctrl-Z on a new line as EOF, while
	// Unix shells use Ctrl-D. Listing both avoids guessing the host OS.
	fmt.Fprintf(os.Stderr, "Paste %s and press Enter. End with Ctrl-Z then Enter (Windows) or Ctrl-D (macOS/Linux):\n", what)
}

// isTerminal reports whether f is an interactive character device (as opposed to a pipe or file).
// Uses the portable os.FileInfo mode bits — no syscall import needed.
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
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
