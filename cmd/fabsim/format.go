package main

import "fmt"

// Format identifies a deck-construction format. Every fabsim run is scoped to one format —
// there's no "no format" mode, so parseFormat rejects the empty string.
//
// Lives in fabsim because (a) the registry already pre-filters its card pool for the
// current format (NotSilverAgeLegal exclusion) so there's no need for a separate "is this
// legal" predicate to plumb through deck-building, and (b) the only consumers are the CLI
// flag and the default deck-name template that bakes the format into output paths.
type Format string

const (
	// SilverAge is the current live format; see data_sources/silver_age_banlist.txt for
	// the authoritative banned card list.
	SilverAge Format = "silver_age"
)

// parseFormat converts a CLI flag value to a Format. Unknown values return an error
// listing the supported formats so the message is self-describing.
func parseFormat(s string) (Format, error) {
	switch Format(s) {
	case SilverAge:
		return Format(s), nil
	}
	return "", fmt.Errorf("unknown format %q (supported: %q)", s, SilverAge)
}
