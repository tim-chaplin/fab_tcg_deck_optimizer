// Command cardgen regenerates the <card>_gen.go files and the registry's cardsByID map
// from the <card>.yaml data files, writing the output to disk. The generation logic and
// yaml schema live in internal/cardgen.
//
// Usage:
//
//	go run ./cmd/cardgen [-registry <path>] <dir>...
//
// -registry <path> additionally emits the registry cardsByID file at that path.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/cardgen"
)

func main() {
	registryOut := flag.String("registry", "",
		"if set, emit the registry cardsByID file at this path (built from the cards-package directory)")
	flag.Parse()
	dirs := flag.Args()
	if len(dirs) == 0 {
		fmt.Fprintln(os.Stderr, "usage: cardgen [-registry <path>] <dir>...")
		os.Exit(2)
	}

	files, err := cardgen.Generate(dirs, *registryOut)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cardgen: %v\n", err)
		os.Exit(1)
	}
	for path, content := range files {
		if err := os.WriteFile(path, content, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "cardgen: %v\n", err)
			os.Exit(1)
		}
	}
}
