// Package format defines the constructed gameplay formats — the deck-construction rulesets
// that decide which cards are legal to include in a deck. Today only Silver Age exists; its
// banlist is the hardcoded set in banlist.go.
//
// A Format is consumed by internal/registry to filter the legal-card pool (ANDed with the
// hero-pool class/talent check the registry owns) and by cmd/fabsim to scope a run and name
// its output decks. The package is deliberately narrow: it knows only about a card's name and
// its printed rarity, not its behaviour, so it stays decoupled from internal/card.
package format

import (
	"fmt"
	"sort"
)

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
	// the format's banlist. This is the name-only check; the rarity rule is IsRarityLegal,
	// which the registry ANDs in (it has the card's printed rarity, which a banlist doesn't
	// need). Kept separate so name-only callers — e.g. cmd/parsecarddb's banlist filter — can
	// ask the banlist question without a rarity in hand.
	IsCardLegal(c Card) bool
	// IsRarityLegal reports whether a card printed at the given rarity (its lowest printed
	// rarity — see internal/cardgen) is allowed in this format. Takes the rarity string rather
	// than a Card because rarity lives on the generated card, not the minimal Card view.
	IsRarityLegal(rarity string) bool
	// BannedNames returns the card names this format bans, sorted.
	BannedNames() []string
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

// silverAgeRarities are the printed rarities Silver Age allows: a card is legal only if its
// lowest printed rarity is one of these (so a card printed at both Majestic and Rare qualifies).
var silverAgeRarities = map[string]bool{"Basic": true, "Common": true, "Rare": true}

// IsRarityLegal reports whether rarity is Basic, Common or Rare — the Silver Age rarity rule.
func (silverAge) IsRarityLegal(rarity string) bool {
	return silverAgeRarities[rarity]
}

// BannedNames returns the Silver Age banned card names, sorted.
func (silverAge) BannedNames() []string {
	out := make([]string, 0, len(silverAgeBanlist))
	for name := range silverAgeBanlist {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
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
