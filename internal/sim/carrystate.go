package sim

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
//     independent backing arrays. Used at ownership-transfer points where the caller
//     needs the value to outlive the next reuse-helper call against the same source.
//
// The aliasing rule for reuse helpers: the receiver's slices are SHARED across calls,
// so any value that needs to outlive the next call must Clone first.

// SnapshotFromTurn copies every persistent TurnState field into c, reusing c's slice
// backings via append([:0], src...). The slice copies are intentional: mid-chain
// state.* slices alias attackBufs scratch storage and the next permutation will
// overwrite them. Reads s.deck / s.graveyard directly so the snapshot itself doesn't
// poison cacheable.
func (c *CarryState) SnapshotFromTurn(s *TurnState) {
	c.Hand = append(c.Hand[:0], s.Hand...)
	c.Deck = append(c.Deck[:0], s.deck...)
	c.Arsenal = s.Arsenal
	c.Graveyard = append(c.Graveyard[:0], s.graveyard...)
	c.Banish = append(c.Banish[:0], s.Banish...)
	c.Runechants = s.Runechants
	c.AuraTriggers = append(c.AuraTriggers[:0], s.AuraTriggers...)
	c.Log = append(c.Log[:0], s.turnLog...)
}

// CopyFrom copies every field of src into c, reusing c's slice backings. Used to
// promote one already-built CarryState into a different scratch (e.g.
// bestCarryScratch → findBestCarryScratch when a leaf wins) without paying a fresh
// allocation per promotion.
func (c *CarryState) CopyFrom(src *CarryState) {
	c.Hand = append(c.Hand[:0], src.Hand...)
	c.Deck = append(c.Deck[:0], src.Deck...)
	c.Arsenal = src.Arsenal
	c.Graveyard = append(c.Graveyard[:0], src.Graveyard...)
	c.Banish = append(c.Banish[:0], src.Banish...)
	c.Runechants = src.Runechants
	c.AuraTriggers = append(c.AuraTriggers[:0], src.AuraTriggers...)
	c.Log = append(c.Log[:0], src.Log...)
}

// Reset zeros every field of c while preserving slice backing arrays. Slice lengths
// drop to 0 (backing array kept for reuse via the next append([:0], ...)); scalar /
// pointer fields zero out. Called at the top of an iteration so a stale value from a
// previous run can't leak through when no candidate is promoted.
func (c *CarryState) Reset() {
	c.Hand = c.Hand[:0]
	c.Deck = c.Deck[:0]
	c.Arsenal = nil
	c.Graveyard = c.Graveyard[:0]
	c.Banish = c.Banish[:0]
	c.Runechants = 0
	c.AuraTriggers = c.AuraTriggers[:0]
	c.Log = c.Log[:0]
}

// Clone returns a fresh CarryState whose slice fields own independent backing arrays.
// Used at ownership-transfer points (e.g. the final TurnSummary returned by findBest)
// so the result survives subsequent reuse-helper calls. Empty slices stay nil to keep
// trivial CarryStates allocation-free.
func (c CarryState) Clone() CarryState {
	out := CarryState{
		Arsenal:    c.Arsenal,
		Runechants: c.Runechants,
	}
	if len(c.Hand) > 0 {
		out.Hand = append([]Card(nil), c.Hand...)
	}
	if len(c.Deck) > 0 {
		out.Deck = append([]Card(nil), c.Deck...)
	}
	if len(c.Graveyard) > 0 {
		out.Graveyard = append([]Card(nil), c.Graveyard...)
	}
	if len(c.Banish) > 0 {
		out.Banish = append([]Card(nil), c.Banish...)
	}
	if len(c.AuraTriggers) > 0 {
		out.AuraTriggers = append([]AuraTrigger(nil), c.AuraTriggers...)
	}
	if len(c.Log) > 0 {
		out.Log = append([]LogEntry(nil), c.Log...)
	}
	return out
}
