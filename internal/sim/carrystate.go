package sim

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// CarryState's reuse / clone helpers. CarryState owns its slice fields, so the methods
// that copy it from various sources or duplicate it for ownership transfer live with the
// type — adding a new persistent field means updating one method per helper.
//
// The four helpers come in two pairs by allocation behavior:
//
//   - Reuse helpers (SnapshotFromTurn, CopyFrom, Reset): operate on the receiver's
//     existing slice backings via append([:0], src...) or [:0] re-slice. Allocation-free
//     after the first sizing.
//   - Clone (returns a fresh CarryState): allocates new slices so the result owns
//     independent backing arrays.

// SnapshotFromTurn captures every persistent GameEngine field into c. Slice fields reuse
// c's backings via append([:0], src...) — mid-chain engine slices alias per-engine scratch
// storage that the next permutation overwrites, so the copy is necessary. Deck gets a
// fresh Copy() so subsequent permutations' deck mutations don't reach back into the
// snapshot.
func (c *CarryState) SnapshotFromTurn(g *gameengine.GameEngine) {
	snap := g.Snapshot()
	c.Hand = append(c.Hand[:0], snap.Hand...)
	c.Deck = snap.Deck
	c.Arsenal = snap.Arsenal
	c.Graveyard = append(c.Graveyard[:0], snap.Graveyard...)
	c.Banish = append(c.Banish[:0], snap.Banished...)
	c.Auras = append(c.Auras[:0], snap.Auras...)
	c.Items = append(c.Items[:0], snap.Items...)
	c.CardsDrawn = snap.CardsDrawn
	c.OpponentMarked = snap.OpponentMarked
	c.Log = append(c.Log[:0], snap.LogEntries...)
}

// CopyFrom copies every field of src into c. Slice fields reuse c's backings via
// append([:0], ...); Deck gets a fresh Copy().
func (c *CarryState) CopyFrom(src *CarryState) {
	c.Hand = append(c.Hand[:0], src.Hand...)
	if src.Deck != nil {
		c.Deck = src.Deck.Copy()
	} else {
		c.Deck = nil
	}
	c.Arsenal = src.Arsenal
	c.Graveyard = append(c.Graveyard[:0], src.Graveyard...)
	c.Banish = append(c.Banish[:0], src.Banish...)
	c.Auras = append(c.Auras[:0], src.Auras...)
	c.Items = append(c.Items[:0], src.Items...)
	c.CardsDrawn = src.CardsDrawn
	c.OpponentMarked = src.OpponentMarked
	c.Log = append(c.Log[:0], src.Log...)
}

// Reset zeros every field of c while preserving slice backing arrays.
func (c *CarryState) Reset() {
	c.Hand = c.Hand[:0]
	c.Deck = nil
	c.Arsenal = nil
	c.Graveyard = c.Graveyard[:0]
	c.Banish = c.Banish[:0]
	c.Auras = c.Auras[:0]
	c.Items = c.Items[:0]
	c.CardsDrawn = 0
	c.OpponentMarked = false
	c.Log = c.Log[:0]
}

// Clone returns a fresh CarryState whose slice / deck fields own independent backing.
func (c CarryState) Clone() CarryState {
	out := CarryState{
		Arsenal:        c.Arsenal,
		CardsDrawn:     c.CardsDrawn,
		OpponentMarked: c.OpponentMarked,
	}
	if len(c.Hand) > 0 {
		out.Hand = append([]card.Card(nil), c.Hand...)
	}
	if c.Deck != nil {
		out.Deck = c.Deck.Copy()
	}
	if len(c.Graveyard) > 0 {
		out.Graveyard = append([]card.Card(nil), c.Graveyard...)
	}
	if len(c.Banish) > 0 {
		out.Banish = append([]card.Card(nil), c.Banish...)
	}
	if len(c.Auras) > 0 {
		out.Auras = append([]gameengine.Aura(nil), c.Auras...)
	}
	if len(c.Items) > 0 {
		out.Items = append([]gameengine.Item(nil), c.Items...)
	}
	if len(c.Log) > 0 {
		out.Log = append([]turnlogger.LogEntry(nil), c.Log...)
	}
	return out
}
