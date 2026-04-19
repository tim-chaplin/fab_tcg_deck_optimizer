package main

import (
	"fmt"
	"os"
)

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
