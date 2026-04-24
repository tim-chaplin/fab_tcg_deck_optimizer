// Command fabsim searches for, evaluates, and iterates on Flesh and Blood decks.
//
// main.go dispatches subcommand args; the per-concern helpers live alongside in
// flags.go (flag-parsing), deckio_helpers.go (deck load / save plumbing), and print.go
// (human-readable renderings). Each subcommand has its own mode_*.go file.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/mydecks"
)

// defaultMaxCopies is the shared fallback for every subcommand's -max-copies flag so the
// anneal hill-climb, eval's sanitize pass, and any future caller all agree on "how many
// copies of one printing is a normal deck allowed to hold" without forking the default.
const defaultMaxCopies = 2

func main() {
	subcommand, args, ok := extractSubcommand()
	if !ok {
		printSubcommands(os.Stdout)
		return
	}

	// Create mydecks/ up front so downstream WriteFile calls can't fail on a missing dir after
	// a long run. Every subcommand reads or writes this directory.
	if err := os.MkdirAll(mydecks.Dir, 0o755); err != nil {
		die("mkdir %s: %v", mydecks.Dir, err)
	}

	switch subcommand {
	case "help":
		printSubcommands(os.Stdout)
	case "anneal":
		runAnnealCmd(args)
	case "eval":
		runEvalCmd(args)
	case "diff":
		runDiffCmd(args)
	case "import":
		runImport()
	default:
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", subcommand)
		printSubcommands(os.Stderr)
		os.Exit(2)
	}
}

// extractSubcommand pulls os.Args[1] as the subcommand name and returns the remaining args
// for the subcommand's own flag.FlagSet. Returns (_, _, false) when no subcommand is given
// or the first arg looks like a flag; the caller prints the subcommand list.
func extractSubcommand() (string, []string, bool) {
	if len(os.Args) < 2 {
		return "", nil, false
	}
	first := os.Args[1]
	if strings.HasPrefix(first, "-") {
		return "", nil, false
	}
	return first, os.Args[2:], true
}

// printSubcommands writes the one-liner catalogue shown when no subcommand is given. Flag
// details live behind `fabsim <subcommand> -help`, which each subcommand's own FlagSet renders.
func printSubcommands(w io.Writer) {
	fmt.Fprintln(w, "fabsim: Flesh and Blood goldfishing deck optimizer")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage: fabsim <subcommand> [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Subcommands:")
	fmt.Fprintln(w, "  anneal    Hill-climb (optionally simulated-annealing) on the saved deck until a local maximum")
	fmt.Fprintln(w, "  eval      Re-score the saved deck at -deep-shuffles and rewrite it; -print-only skips the sim (usage: fabsim eval <deck>)")
	fmt.Fprintln(w, "  import    Paste a fabrary.net deck into mydecks/<name>.json")
	fmt.Fprintln(w, "  diff      Print the card-count delta between two saved decks (usage: fabsim diff <deck1> <deck2>)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Run 'fabsim <subcommand> -help' for flag details.")
}

func die(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
