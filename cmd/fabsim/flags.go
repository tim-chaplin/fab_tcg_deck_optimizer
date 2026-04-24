package main

// Flag-parsing helpers shared by every subcommand: requireFlag dies with a usage error when
// a required flag wasn't supplied, parseFlagsAnywhere lets flags appear interleaved with
// positional args so the user never has to care about token order.

import (
	"flag"
	"strings"
)

// requireFlag dies with a usage error when fs.Parse didn't encounter -name. The flag's own Usage
// string is echoed so the per-flag guidance the caller wrote in the FlagSet (e.g. why the flag
// can't default) shows up alongside the "required" message — no need to duplicate that wording
// here.
func requireFlag(fs *flag.FlagSet, subcommand, name string) {
	seen := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			seen = true
		}
	})
	if seen {
		return
	}
	f := fs.Lookup(name)
	if f == nil {
		die("%s: internal error — required flag -%s is not registered on the FlagSet", subcommand, name)
	}
	die("%s: -%s is required\n  usage: %s", subcommand, name, f.Usage)
}

// parseFlagsAnywhere parses args on fs while tolerating flags that appear before, after, or
// interleaved with positional arguments. Go's stdlib flag package stops at the first
// positional token; every subcommand routes through this helper so flag order never matters
// to the user.
//
// Bool-awareness matters: a non-bool `-name` token consumes the following arg as its value,
// a bool flag (detected via IsBoolFlag) doesn't. Unknown flags pass through untouched so
// fs.Parse emits the canonical "flag provided but not defined" error. A bare `--` acts as
// the end-of-flags terminator, with everything after it treated as positional.
func parseFlagsAnywhere(fs *flag.FlagSet, args []string) error {
	var flagTokens, positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--" {
			positional = append(positional, args[i+1:]...)
			break
		}
		if len(a) < 2 || a[0] != '-' {
			positional = append(positional, a)
			continue
		}
		flagTokens = append(flagTokens, a)
		// -name=value is self-contained; only -name (no =) can consume the next token.
		if strings.Contains(a, "=") {
			continue
		}
		name := strings.TrimLeft(a, "-")
		f := fs.Lookup(name)
		if f == nil {
			// Unknown flag — let fs.Parse produce the canonical error with usage.
			continue
		}
		bf, ok := f.Value.(interface{ IsBoolFlag() bool })
		isBool := ok && bf.IsBoolFlag()
		if !isBool && i+1 < len(args) {
			i++
			flagTokens = append(flagTokens, args[i])
		}
	}
	// Insert `--` before the positional block so fs.Parse doesn't re-interpret positional
	// tokens that happen to start with `-`. That's necessary when the caller supplied an
	// explicit `--` terminator (whose tail ends up in positional) and also future-proofs
	// against deck names that begin with a dash.
	out := append(flagTokens, "--")
	out = append(out, positional...)
	return fs.Parse(out)
}
