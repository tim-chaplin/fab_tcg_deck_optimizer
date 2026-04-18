// Package format enumerates the deck-construction formats fabsim supports and provides the
// per-format legality filter used by deck generation and mutation. Cards opt out of a format via
// marker interfaces on card.Card (e.g. card.NotSilverAgeLegal); this package translates a
// Format value into the equivalent predicate.
package format

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Format identifies a deck-construction format. The zero value (Unrestricted) skips all
// legality filtering — every implemented card is deckable.
type Format string

const (
	// Unrestricted allows every implemented card. Intended for tests and for callers that want to
	// explore the full card pool regardless of any format's banlist.
	Unrestricted Format = ""
	// SilverAge is the current live format as of 2026; see data_sources/silver_age_banlist.txt for
	// the authoritative banned card list.
	SilverAge Format = "silver_age"
)

// Parse converts a CLI flag value to a Format. Unknown values return an error listing the
// supported formats so the failure message is self-describing.
func Parse(s string) (Format, error) {
	switch Format(s) {
	case Unrestricted, SilverAge:
		return Format(s), nil
	}
	return "", fmt.Errorf("unknown format %q (supported: %q, %q)", s, Unrestricted, SilverAge)
}

// IsLegal reports whether c may appear in a deck built for this format. Unrestricted always
// returns true; Silver Age rejects cards tagged with the card.NotSilverAgeLegal marker.
func (f Format) IsLegal(c card.Card) bool {
	switch f {
	case SilverAge:
		_, banned := c.(card.NotSilverAgeLegal)
		return !banned
	default:
		return true
	}
}
