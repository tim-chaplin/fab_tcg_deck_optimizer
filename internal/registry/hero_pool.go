package registry

import (
	"fmt"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
)

// Hero-pool legality: the universal class/talent deckbuilding rule, independent of format. A
// hero may include a card only when the card's class is Generic or the hero's own class, and
// every talent on the card is one the hero shares. The registry owns this (the format does
// not) because it holds the same way in every format. See internal/format for the orthogonal
// banlist check that LegalCardsFor ANDs with this.
//
// classMask and talentMask name which card.CardType bits are classes and which are talents.
// They list every class / talent constant that exists today; extend them in lockstep with
// internal/card/types.go as new classes and talents are added (a new talent constant that
// isn't in talentMask would be silently ignored by the subset check below).

// classMask is the set of type bits that denote a card's class line.
var classMask = card.NewTypeSet(card.TypeGeneric, card.TypeRuneblade)

// talentMask is the set of type bits that denote a talent.
var talentMask = card.NewTypeSet(card.TypeLightning)

// heroCanPlay reports whether hero h may include card c in a deck: c's class is Generic or
// h's class, and c's talents are a subset of h's. The hero's talents are read off its type
// line (h.Types() & talentMask) rather than a dedicated accessor, mirroring how a card's are.
func heroCanPlay(h card.Hero, c card.Card) bool {
	heroTalents := h.Types() & talentMask

	// Universal cards adopt the active hero's class at resolution, so they're always
	// class-legal, and their printed type-line carries no talent. Their Types method folds in
	// the hero's class by calling into the engine, which a static legality check has no handle
	// for, so they're matched by marker and accepted directly.
	if _, ok := c.(interface{ Universal() }); ok {
		return true
	}

	ct := c.Types(nil)

	// Talent subset: every talent bit on the card must be one the hero has.
	if ct&talentMask&^heroTalents != 0 {
		return false
	}
	// Class: every class bit on the card must be Generic or the hero's class. A card with no
	// class bit at all (neither Generic nor a class) is unrestricted on this axis.
	allowedClasses := card.NewTypeSet(card.TypeGeneric, h.Class())
	if ct&classMask&^allowedClasses != 0 {
		return false
	}
	return true
}

// asCardHero narrows the deck's opaque hero value to the card.Hero view heroCanPlay needs.
// Every registered hero implements card.Hero (Class / Types); a value that doesn't reach the
// production registry is a programming error at the call site.
func asCardHero(h deck.Hero) card.Hero {
	hc, ok := h.(card.Hero)
	if !ok {
		panic(fmt.Sprintf("registry: hero %T does not implement card.Hero (Class / Types)", h))
	}
	return hc
}
