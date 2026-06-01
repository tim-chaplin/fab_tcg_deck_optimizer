// Package format defines the constructed gameplay formats — the deck-construction rulesets
// that decide which cards are legal to include in a deck. Today only Silver Age exists; its
// banlist is the hardcoded set in banlist.go.
//
// A Format is consumed by internal/registry to filter the legal-card pool (ANDed with the
// hero-pool class/talent check the registry owns) and by cmd/fabsim to scope a run and name
// its output decks. The package is deliberately narrow: it knows only about a card's name,
// not its behaviour, so it stays decoupled from internal/card.
package format

import "fmt"

// Card is the minimal view a Format needs of a card to judge legality: just its printed
// name, which is what a banlist keys on. Declared locally so format stays decoupled from
// internal/card — mirroring the narrow-interface pattern internal/registry uses. Concrete
// card.Card values satisfy it.
type Card interface {
	Name() string
}

// Format is a constructed gameplay format: a banlist today, plus (eventually) a legal-card
// pool and copy caps. IsCardLegal is the per-card legality predicate the registry ANDs with
// hero-pool compatibility to build the deck-construction pool. A Format is hero-agnostic —
// class/talent legality is a universal deckbuilding rule the registry owns, not a per-format
// rule.
type Format interface {
	// Name is the stable identifier used in CLI flags and output deck filenames
	// (e.g. "silver_age").
	Name() string
	// DisplayName is the human-facing label used in deck headers (e.g. "Silver Age").
	DisplayName() string
	// IsCardLegal reports whether c is legal in this format — today, that its name is not on
	// the format's banlist.
	IsCardLegal(c Card) bool
}

// SilverAge is the current live format. Its banlist is silverAgeBanlist in banlist.go.
var SilverAge Format = silverAge{}

// silverAge implements Format for the Silver Age constructed format.
type silverAge struct{}

func (silverAge) Name() string        { return "silver_age" }
func (silverAge) DisplayName() string { return "Silver Age" }

// IsCardLegal reports whether c is legal in Silver Age — its name is not on the banlist.
func (silverAge) IsCardLegal(c Card) bool {
	_, banned := silverAgeBanlist[c.Name()]
	return !banned
}

// Parse converts a CLI flag value to a Format. Unknown values return an error listing the
// supported formats so the message is self-describing.
func Parse(s string) (Format, error) {
	switch s {
	case SilverAge.Name():
		return SilverAge, nil
	}
	return nil, fmt.Errorf("unknown gameplay format %q (supported: %q)", s, SilverAge.Name())
}
