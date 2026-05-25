package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Get on a fresh Pool returns a non-nil, zero-value *GameState (allocated, not pooled).
func TestPool_GetOnEmptyAllocates(t *testing.T) {
	p := NewPool(4)
	s := p.Get()
	if s == nil {
		t.Fatal("Get on empty Pool returned nil")
	}
}

// Put-then-Get returns the same pointer — verifies the free list semantics.
func TestPool_PutGetRoundTrip(t *testing.T) {
	p := NewPool(4)
	original := p.Get() // increments in-flight so the Put balances out.
	p.Put(original)
	if got := p.Get(); got != original {
		t.Errorf("Get returned %p, want %p (the just-Put pointer)", got, original)
	}
}

// Putting nil is a no-op so callers don't have to gate on the winner-vs-loser branch.
func TestPool_PutNilIgnored(t *testing.T) {
	p := NewPool(4)
	p.Put(nil)
	s := p.Get()
	if s == nil {
		t.Fatal("Get after Put(nil) returned nil — Put(nil) should not have populated the free list")
	}
}

// Slice backings on a pooled state survive Put → Get, so the next caller's
// CopyPersistentStateFrom can reuse them via the resetCardSlice / copyAurasInto cap
// checks. The whole point of pooling here is this backing reuse.
func TestPool_RetainsSliceBackings(t *testing.T) {
	p := NewPool(4)
	s := p.Get()
	s.graveyard = make([]card.Card, 0, 32) // synthetic backing with cap > 0
	originalCap := cap(s.graveyard)
	p.Put(s)
	got := p.Get()
	if cap(got.graveyard) != originalCap {
		t.Errorf("graveyard cap after Put→Get = %d, want %d (backing should survive)",
			cap(got.graveyard), originalCap)
	}
}

// Get panics when in-flight checkouts would exceed cap — the operator-visible signal
// to raise the cap or tighten Put discipline.
func TestPool_GetOverCapPanics(t *testing.T) {
	p := NewPool(2)
	p.Get()
	p.Get()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Get past cap did not panic")
		}
	}()
	p.Get()
}

// FreeAll resets in-flight back to zero so the same Pool can be reused across cycles
// (e.g. between shuffles) without the caller having to Put every escapee individually.
func TestPool_FreeAllResetsInFlight(t *testing.T) {
	p := NewPool(2)
	p.Get()
	p.Get()
	p.FreeAll()
	// After FreeAll a fresh Get should succeed even though we never Put.
	if s := p.Get(); s == nil {
		t.Fatal("Get after FreeAll returned nil")
	}
}

// HighWaterMark reports the max in-flight ever observed between FreeAll calls — the
// number an operator reads to size a future Pool.
func TestPool_HighWaterMarkPersistsAcrossFreeAll(t *testing.T) {
	p := NewPool(8)
	p.Get()
	p.Get()
	p.Get()
	p.FreeAll()
	p.Get()
	if got := p.HighWaterMark(); got != 3 {
		t.Errorf("HighWaterMark = %d, want 3 (peak observed in the first cycle)", got)
	}
}
