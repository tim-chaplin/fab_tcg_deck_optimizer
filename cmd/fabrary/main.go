// Command fabrary converts decks between the optimizer's JSON format and fabrary.net's plain-text
// format (the text you get from Export → Plain Text on fabrary, or paste into its Import tab).
//
// Export:
//
//	fabrary -mode export -in best_deck.json            # writes fabrary text to stdout
//	fabrary -mode export -in best_deck.json -out x.txt
//
// Import (default: prompt for a deck name, save to mydecks/<name>.json):
//
//	fabrary -mode import                                 # paste into stdin; interactive name prompt
//	fabrary -mode import -in paste.txt                   # read from file; still prompts for name
//	fabrary -mode import -in paste.txt -out custom.json  # -out bypasses the mydecks/ default
//
// When -in is empty or "-" stdin is read. When -out is empty, export writes to stdout and import
// saves to mydecks/<prompted-name>.json.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
)

// myDecksDir is the directory imported decks default to. Kept relative so the command behaves the
// same regardless of where the user runs it from, matching fabsim's "-out best_deck.json" default.
const myDecksDir = "mydecks"

func main() {
	mode := flag.String("mode", "export", "conversion direction: export (json→fabrary) or import (fabrary→json)")
	inPath := flag.String("in", "", "input path; \"-\" or empty reads from stdin")
	outPath := flag.String("out", "", "output path; empty = stdout for export / mydecks/<name>.json for import")
	flag.Parse()

	switch *mode {
	case "export":
		runExport(*inPath, *outPath)
	case "import":
		runImport(*inPath, *outPath)
	default:
		die("unknown mode %q (want export or import)", *mode)
	}
}

func runExport(inPath, outPath string) {
	data, err := readAllFrom(inPath, "deck JSON")
	if err != nil {
		die("read input: %v", err)
	}
	d, err := deckio.Unmarshal(data)
	if err != nil {
		die("parse deck JSON: %v", err)
	}
	out := []byte(fabrary.Marshal(d))
	dest := outPath
	if dest == "" {
		dest = "-"
	}
	if err := writeOut(dest, out); err != nil {
		die("write %s: %v", dest, err)
	}
}

func runImport(inPath, outPath string) {
	data, name, err := readImport(inPath, outPath == "")
	if err != nil {
		die("%v", err)
	}
	d, err := fabrary.Unmarshal(string(data))
	if err != nil {
		die("parse fabrary text: %v", err)
	}
	out, err := deckio.Marshal(d)
	if err != nil {
		die("encode deck JSON: %v", err)
	}
	out = append(out, '\n')

	dest := outPath
	if dest == "" {
		if err := os.MkdirAll(myDecksDir, 0o755); err != nil {
			die("mkdir %s: %v", myDecksDir, err)
		}
		dest = filepath.Join(myDecksDir, name+".json")
	}
	if err := writeOut(dest, out); err != nil {
		die("write %s: %v", dest, err)
	}
	if dest != "-" {
		fmt.Fprintf(os.Stderr, "wrote %s\n", dest)
		summarizeDeck(d, dest)
	}
}

// readImport reads fabrary text plus — if needDeckName — a deck name from the same stdin. The name
// must be prompted before the paste because Ctrl-Z / Ctrl-D closes stdin permanently, leaving no
// way to ask afterward. When -in points at a file, the name is still prompted interactively so the
// user gets to choose the filename.
func readImport(inPath string, needDeckName bool) ([]byte, string, error) {
	fromStdin := inPath == "" || inPath == "-"

	if !fromStdin {
		data, err := os.ReadFile(inPath)
		if err != nil {
			return nil, "", fmt.Errorf("read %s: %w", inPath, err)
		}
		var name string
		if needDeckName {
			if !isTerminal(os.Stdin) {
				return nil, "", fmt.Errorf("need a deck name but stdin is not a terminal — pass -out to skip the prompt")
			}
			n, err := promptLine(os.Stdin, "Deck name: ")
			if err != nil {
				return nil, "", err
			}
			name = n
		}
		return data, name, nil
	}

	// Pasted workflow: name first (one line), paste second (rest of stdin). A single bufio.Reader
	// owns stdin for both reads so any bytes the line scanner buffers past the newline flow into
	// the io.ReadAll below instead of being dropped.
	reader := bufio.NewReader(os.Stdin)
	var name string
	if needDeckName {
		if !isTerminal(os.Stdin) {
			return nil, "", fmt.Errorf("need a deck name but stdin is not a terminal — pass -out to skip the prompt")
		}
		n, err := promptLineReader(reader, "Deck name: ")
		if err != nil {
			return nil, "", err
		}
		name = n
	}
	if isTerminal(os.Stdin) {
		promptPaste()
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, "", fmt.Errorf("read stdin: %w", err)
	}
	return data, name, nil
}

// readAllFrom reads from a file path, stdin, or "-"-stdin. Shared by export mode where there's no
// interactive name prompt.
func readAllFrom(path, what string) ([]byte, error) {
	if path != "" && path != "-" {
		return os.ReadFile(path)
	}
	if isTerminal(os.Stdin) {
		fmt.Fprintf(os.Stderr, "Paste %s and press Enter. End with Ctrl-Z then Enter (Windows) or Ctrl-D (macOS/Linux):\n", what)
	}
	return io.ReadAll(os.Stdin)
}

func promptPaste() {
	fmt.Fprintln(os.Stderr, "Paste fabrary deck text and press Enter. End with Ctrl-Z then Enter (Windows) or Ctrl-D (macOS/Linux):")
}

func promptLine(f *os.File, prompt string) (string, error) {
	return promptLineReader(bufio.NewReader(f), prompt)
}

// promptLineReader writes the prompt, reads one line, trims whitespace, validates it as a deck
// name, and returns it. Accepts an existing bufio.Reader so the caller can continue draining the
// same buffered stream afterwards.
func promptLineReader(r *bufio.Reader, prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	line, err := r.ReadString('\n')
	if err != nil && line == "" {
		return "", fmt.Errorf("read deck name: %w", err)
	}
	name := strings.TrimSpace(line)
	name = strings.TrimSuffix(name, ".json")
	if err := validateDeckName(name); err != nil {
		return "", err
	}
	return name, nil
}

// validateDeckName rejects names that would escape mydecks/ or otherwise produce an unusable file
// path. Kept conservative — any unusual character the user actually wants can be passed via -out.
func validateDeckName(name string) error {
	if name == "" {
		return fmt.Errorf("deck name is empty")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("deck name %q is reserved", name)
	}
	if strings.ContainsAny(name, `/\:*?"<>|`) {
		return fmt.Errorf("deck name %q contains an invalid character (one of /\\:*?\"<>|)", name)
	}
	return nil
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

// summarizeDeck prints a short confirmation — hero, weapon count, card count — to stderr after a
// successful import, so the user can sanity-check the paste without opening the file.
func summarizeDeck(d *deck.Deck, _ string) {
	weapons := make([]string, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.Name()
	}
	fmt.Fprintf(os.Stderr, "  hero: %s, weapons: %v, cards: %d\n", d.Hero.Name(), weapons, len(d.Cards))
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "fabrary: "+format+"\n", args...)
	os.Exit(1)
}
