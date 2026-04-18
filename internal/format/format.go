// Package format enumerates the deck-construction formats fabsim supports and provides the
// per-format legality filter used by deck generation and mutation. Cards opt out of a format via
// marker interfaces on card.Card (e.g. card.NotSilverAgeLegal); this package translates a
// Format value into the equivalent predicate.
package format

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Format identifies a deck-construction format. Every fabsim run is scoped to exactly one
// format — there's no "no format" mode, so Parse rejects the empty string and callers should
// always pass one of the named constants.
type Format string

const (
	// SilverAge is the current live format as of 2026; see data_sources/silver_age_banlist.txt for
	// the authoritative banned card list.
	SilverAge Format = "silver_age"
)

// Parse converts a CLI flag value to a Format. Unknown values return an error listing the
// supported formats so the failure message is self-describing.
func Parse(s string) (Format, error) {
	switch Format(s) {
	case SilverAge:
		return Format(s), nil
	}
	return "", fmt.Errorf("unknown format %q (supported: %q)", s, SilverAge)
}

// IsLegal reports whether c may appear in a deck built for this format. Silver Age rejects
// cards tagged with the card.NotSilverAgeLegal marker; other formats (when added) plug in here.
func (f Format) IsLegal(c card.Card) bool {
	switch f {
	case SilverAge:
		_, banned := c.(card.NotSilverAgeLegal)
		return !banned
	default:
		panic(fmt.Sprintf("format: IsLegal called on unknown Format %q", f))
	}
}
