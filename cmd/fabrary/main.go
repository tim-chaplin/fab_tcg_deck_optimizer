// Command fabrary imports a fabrary.net plain-text deck (the thing you get from Export → Plain
// Text on fabrary, or copy out of its Import tab) into the optimizer's JSON format.
//
// Every import lands in mydecks/ — the same directory fabsim writes to — so local decks all live
// in one place.
//
// Usage:
//
//	fabrary                          # paste into stdin; interactive name prompt
//	fabrary -deck viserai-v2         # same, but skip the prompt (useful for scripts/pipes)
//	fabrary -in paste.txt            # read from file; still prompts for a name
//	fabrary -in paste.txt -deck foo  # read from file, name from flag, no prompting
//
// When -in is empty or "-" stdin is read. The ".json" suffix on -deck is optional.
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
	deckName := flag.String("deck", "", "deck name; resolved to mydecks/<name>.json. Empty = prompt interactively")
	flag.Parse()

	data, resolvedName, err := readImport(*inPath, *deckName)
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

	if err := os.MkdirAll(mydecks.Dir, 0o755); err != nil {
		die("mkdir %s: %v", mydecks.Dir, err)
	}
	dest, err := mydecks.Path(resolvedName)
	if err != nil {
		die("%v", err)
	}
	if err := os.WriteFile(dest, out, 0o644); err != nil {
		die("write %s: %v", dest, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", dest)
	summarizeDeck(d)
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

// readImport reads fabrary text plus a deck name. When deckFlag is set it's used directly (after
// name validation); otherwise the name is prompted interactively. In the stdin-paste workflow the
// name must be collected before the paste because Ctrl-Z / Ctrl-D closes stdin permanently — if
// the user didn't pass -deck and stdin isn't a terminal (e.g. piped from a file) we error out
// rather than leaving the import without a filename.
func readImport(inPath, deckFlag string) ([]byte, string, error) {
	fromStdin := inPath == "" || inPath == "-"

	if !fromStdin {
		data, err := os.ReadFile(inPath)
		if err != nil {
			return nil, "", fmt.Errorf("read %s: %w", inPath, err)
		}
		name, err := resolveDeckName(deckFlag, os.Stdin)
		if err != nil {
			return nil, "", err
		}
		return data, name, nil
	}

	// Pasted workflow: name first, paste second. A single bufio.Reader owns stdin for both reads
	// so any bytes buffered past the newline flow into the paste reader instead of being dropped.
	reader := bufio.NewReader(os.Stdin)
	name, err := resolveDeckNameFromReader(deckFlag, reader)
	if err != nil {
		return nil, "", err
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

// resolveDeckName returns the deck name to save under: validated -deck flag value, or the result
// of an interactive prompt on src. Errors if -deck is empty and src isn't a terminal (piping
// fabrary text without also passing -deck would leave us with no filename to use).
func resolveDeckName(deckFlag string, src *os.File) (string, error) {
	if deckFlag != "" {
		return sanitizeDeckFlag(deckFlag)
	}
	if !isTerminal(src) {
		return "", fmt.Errorf("no deck name: pass -deck <name> when stdin isn't a terminal")
	}
	return promptLine(src, "Deck name: ")
}

func resolveDeckNameFromReader(deckFlag string, r *bufio.Reader) (string, error) {
	if deckFlag != "" {
		return sanitizeDeckFlag(deckFlag)
	}
	if !isTerminal(os.Stdin) {
		return "", fmt.Errorf("no deck name: pass -deck <name> when stdin isn't a terminal")
	}
	return promptLineReader(r, "Deck name: ")
}

// sanitizeDeckFlag trims the optional .json suffix and validates the rest, matching the same
// rules the interactive prompt applies.
func sanitizeDeckFlag(name string) (string, error) {
	name = strings.TrimSuffix(name, ".json")
	if err := mydecks.ValidateName(name); err != nil {
		return "", err
	}
	return name, nil
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
