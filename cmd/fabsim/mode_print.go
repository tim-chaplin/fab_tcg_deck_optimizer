package main

import (
	"flag"
	"fmt"
	"os"
)

// runPrintCmd parses print's flags (none today) and dispatches to runPrint. print always
// operates on a specific existing deck, so the deck is a positional arg rather than a -deck
// flag.
func runPrintCmd(args []string) {
	fs := flag.NewFlagSet("print", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "Usage: fabsim print <deck>")
	}
	_ = fs.Parse(args)
	if fs.NArg() != 1 {
		die("print: need exactly one positional <deck> (got %d); try `fabsim print <deck>`", fs.NArg())
	}
	runPrint(resolveDeckPath(fs.Arg(0)))
}

// runPrint loads the deck at path and prints it (card list + persisted stats) without running
// any simulation. Use runEval to re-simulate before printing.
func runPrint(path string) {
	d, _ := loadExisting(path)
	if d == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", path)
		os.Exit(1)
	}
	printBestDeck(d)
}
