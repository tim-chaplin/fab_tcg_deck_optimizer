package deck

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// ReorderTop rebuilds the top of the deck after an Opt-style look at the top n cards: top
// (the kept cards, in the given order) goes back on top, the untouched cards below the top n
// follow, and bottom (the rejected cards, in the given order) is appended to the very bottom.
// len(top)+len(bottom) must equal n.
//
// Takes card.Card — the rich contract the hero's Opt splitter returns — and stores the
// values through the deck's narrower Card view, so the caller avoids a []card.Card ->
// []deck.Card copy on the hot Opt path.
//
// Allocates one fresh backing rather than rewriting in place: a per-permutation deck shares
// its backing with the master leaf via ShallowCopyFrom, so an in-place write would corrupt
// sibling permutations — the same hazard Shuffle's mustNotShuffle guard protects against.
func (d *Deck) ReorderTop(n int, top, bottom []card.Card) {
	rest := d.cards[n:]
	combined := make([]Card, 0, len(top)+len(rest)+len(bottom))
	for _, c := range top {
		combined = append(combined, c)
	}
	combined = append(combined, rest...)
	for _, c := range bottom {
		combined = append(combined, c)
	}
	d.cards = combined
}
