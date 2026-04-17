package main

import (
	"bufio"
	"bytes"
	"errors"
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

// runImport interactively pastes a fabrary.net plain-text deck from stdin, asks for a deck name,
// and writes the resulting JSON to mydecks/<name>.json. The name is prompted BEFORE the paste
// because fabrary's footer ends stdin — there's no opportunity to ask afterward.
//
// Piping is incidental rather than first-class: the prompts go to stderr, so something like
// `printf "my-deck\n%s" "$(cat paste.txt)" | fabsim -mode import` still works, but there's no
// dedicated -in flag. Exporting in the other direction isn't a mode at all — every random /
// iterate run writes a sibling fabrary .txt next to the JSON, ready to paste into fabrary.net.
func runImport() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Fprint(os.Stderr, "Deck name: ")
	nameLine, err := reader.ReadString('\n')
	if err != nil && nameLine == "" {
		die("read deck name: %v", err)
	}
	name := strings.TrimSpace(nameLine)
	name = strings.TrimSuffix(name, ".json")
	if err := mydecks.ValidateName(name); err != nil {
		die("%v", err)
	}

	fmt.Fprintln(os.Stderr, "Paste fabrary deck text below — input ends automatically at the 'See the full deck @ …' footer:")
	data, err := readUntilFabraryFooter(reader)
	if err != nil {
		die("read stdin: %v", err)
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

	dest, err := mydecks.Path(name)
	if err != nil {
		die("%v", err)
	}
	if err := os.WriteFile(dest, out, 0o644); err != nil {
		die("write %s: %v", dest, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", dest)
	summarizeImportedDeck(d)
	warnSkipped(skipped)
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

// summarizeImportedDeck prints a short stderr confirmation — hero, weapons, card count — so the
// user can sanity-check the paste without opening the file.
func summarizeImportedDeck(d *deck.Deck) {
	weapons := make([]string, len(d.Weapons))
	for i, w := range d.Weapons {
		weapons[i] = w.Name()
	}
	fmt.Fprintf(os.Stderr, "  hero: %s, weapons: %v, cards: %d\n", d.Hero.Name(), weapons, len(d.Cards))
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
