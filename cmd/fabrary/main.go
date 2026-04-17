// Command fabrary imports a fabrary.net plain-text deck (the thing you get from Export → Plain
// Text on fabrary, or copy out of its Import tab) into the optimizer's JSON format.
//
// Usage:
//
//	fabrary                                 # paste into stdin; interactive name prompt
//	fabrary -in paste.txt                   # read from file; still prompts for a deck name
//	fabrary -in paste.txt -out custom.json  # -out bypasses the mydecks/<name>.json default
//
// When -in is empty or "-" stdin is read. When -out is empty, the deck is saved to
// mydecks/<prompted-name>.json.
//
// Exporting in the other direction isn't a separate command — fabsim writes a sibling .txt in
// fabrary format next to every best_deck.json it saves, so the file is already ready to paste
// into fabrary.net.
package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deckio"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/fabrary"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
)

func main() {
	inPath := flag.String("in", "", "input path; \"-\" or empty reads from stdin")
	outPath := flag.String("out", "", "output path; empty saves to mydecks/<prompted-name>.json")
	flag.Parse()

	data, name, err := readImport(*inPath, *outPath == "")
	if err != nil {
		die("%v", err)
	}
	d, skipped, err := fabrary.Unmarshal(string(data))
	if err != nil {
		die("parse fabrary text: %v", err)
	}
	out, err := deckio.Marshal(d)
	if err != nil {
		die("encode deck JSON: %v", err)
	}
	out = append(out, '\n')

	dest := *outPath
	if dest == "" {
		if err := os.MkdirAll(mydecks.Dir, 0o755); err != nil {
			die("mkdir %s: %v", mydecks.Dir, err)
		}
		p, err := mydecks.Path(name)
		if err != nil {
			die("%v", err)
		}
		dest = p
	}
	if err := writeOut(dest, out); err != nil {
		die("write %s: %v", dest, err)
	}
	if dest != "-" {
		fmt.Fprintf(os.Stderr, "wrote %s\n", dest)
		summarizeDeck(d)
	}
	warnSkipped(skipped)
}

// warnSkipped prints a stderr notice for any fabrary cards the optimizer's registry doesn't yet
// cover. Without this the imported deck would silently be smaller than the user pasted.
func warnSkipped(skipped map[string]int) {
	if len(skipped) == 0 {
		return
	}
	total := 0
	names := make([]string, 0, len(skipped))
	for name, qty := range skipped {
		total += qty
		names = append(names, name)
	}
	sort.Strings(names)
	fmt.Fprintf(os.Stderr, "warning: skipped %d unimplemented card(s):\n", total)
	for _, n := range names {
		fmt.Fprintf(os.Stderr, "  %dx %s\n", skipped[n], n)
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

	// Pasted workflow: name first (one line), paste second (until fabrary footer or EOF). A single
	// bufio.Reader owns stdin for both reads so any bytes buffered past the newline flow into the
	// paste reader instead of being dropped.
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
	data, err := readUntilFabraryFooter(reader)
	if err != nil {
		return nil, "", fmt.Errorf("read stdin: %w", err)
	}
	return data, name, nil
}

// fabraryFooterPrefix is the last line of every fabrary plain-text export. Seeing it means the
// user is done pasting — we stop reading so they don't have to send EOF by hand (Ctrl-Z on
// Windows is especially awkward). EOF is still honored for pastes that have been edited to strip
// the footer.
const fabraryFooterPrefix = "See the full deck"

func readUntilFabraryFooter(r *bufio.Reader) ([]byte, error) {
	var buf bytes.Buffer
	for {
		line, err := r.ReadString('\n')
		buf.WriteString(line)
		if strings.HasPrefix(strings.TrimSpace(line), fabraryFooterPrefix) {
			return buf.Bytes(), nil
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				return buf.Bytes(), nil
			}
			return nil, err
		}
	}
}

func promptPaste() {
	fmt.Fprintln(os.Stderr, "Paste fabrary deck text below — input ends automatically at the 'See the full deck @ …' footer:")
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
	if err := mydecks.ValidateName(name); err != nil {
		return "", err
	}
	return name, nil
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
func summarizeDeck(d *deck.Deck) {
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
