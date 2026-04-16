package main

import (
	"fmt"
	"os"
)

func runPrintOnly(path string) {
	d, _ := loadExisting(path)
	if d == nil {
		fmt.Fprintf(os.Stderr, "could not load deck from %s\n", path)
		os.Exit(1)
	}
	printBestDeck(d)
}
