package main

import "fmt"

// GameplayFormat identifies a deck-construction rule set (banlist + legal card pool +
// copy caps). "Gameplay format" disambiguates from the many other "format" senses in
// software (text formatting, file formats, etc.) — at the gameplay level this is the
// constructed-format choice every fabsim run is scoped to. Empty string is rejected;
// parseGameplayFormat insists on a known value.
//
// Lives in fabsim because (a) the registry already pre-filters its card pool for the
// current format (NotSilverAgeLegal exclusion) so there's no need for a separate "is
// this legal" predicate to plumb through deck-building, and (b) the only consumers are
// the CLI flag and the default deck-name template that bakes the format into output
// paths.
type GameplayFormat string

const (
	// SilverAge is the current live format; see data_sources/silver_age_banlist.txt for
	// the authoritative banned card list.
	SilverAge GameplayFormat = "silver_age"
)

// parseGameplayFormat converts a CLI flag value to a GameplayFormat. Unknown values
// return an error listing the supported formats so the message is self-describing.
func parseGameplayFormat(s string) (GameplayFormat, error) {
	switch GameplayFormat(s) {
	case SilverAge:
		return GameplayFormat(s), nil
	}
	return "", fmt.Errorf("unknown gameplay format %q (supported: %q)", s, SilverAge)
}
