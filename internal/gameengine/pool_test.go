package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// newGameStateForTest returns a zero-value *GameState — the factory shape NewPool requires.
func newGameStateForTest() *GameState { return new(GameState) }

// Get returns a prewarmed slot from the eagerly populated free list.
func TestPool_GetReturnsPrewarmedSlot(t *testing.T) {
	p := NewPool(4, newGameStateForTest)
	s := p.Get()
	if s == nil {
		t.Fatal("Get returned nil")
	}
}

// Put-then-Get returns the same pointer.
func TestPool_PutGetRoundTrip(t *testing.T) {
	p := NewPool(4, newGameStateForTest)
	original := p.Get()
	p.Put(original)
	if got := p.Get(); got != original {
		t.Errorf("Get returned %p, want %p (the just-Put pointer)", got, original)
	}
}

// Putting nil is a no-op.
func TestPool_PutNilIgnored(t *testing.T) {
	p := NewPool(4, newGameStateForTest)
	p.Put(nil)
	s := p.Get()
	if s == nil {
		t.Fatal("Get after Put(nil) returned nil — Put(nil) should not have populated the free list")
	}
}

// Slice backings on a pooled state survive Put → Get so the next caller reuses the cap.
func TestPool_RetainsSliceBackings(t *testing.T) {
	p := NewPool(4, newGameStateForTest)
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

// Get panics when in-flight checkouts would exceed cap.
func TestPool_GetOverCapPanics(t *testing.T) {
	p := NewPool(2, newGameStateForTest)
	p.Get()
	p.Get()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("Get past cap did not panic")
		}
	}()
	p.Get()
}

// FreeAll resets in-flight accounting to zero; slot reclamation still requires Puts.
func TestPool_FreeAllResetsInFlightAccounting(t *testing.T) {
	p := NewPool(2, newGameStateForTest)
	a := p.Get()
	b := p.Get()
	p.Put(a)
	p.Put(b)
	p.FreeAll()
	p.Get()
	p.Get()
}

// Factory is called eagerly at NewPool and never again for the Pool's lifetime.
func TestPool_FactoryCalledEagerly(t *testing.T) {
	calls := 0
	factory := func() *GameState {
		calls++
		return new(GameState)
	}
	p := NewPool(3, factory)
	if calls != 3 {
		t.Fatalf("NewPool eager prewarm called factory %d times, want 3", calls)
	}
	a := p.Get()
	b := p.Get()
	if calls != 3 {
		t.Errorf("Get reused prewarmed slots; factory call count = %d, want 3", calls)
	}
	p.Put(a)
	p.Put(b)
	p.Get()
	if calls != 3 {
		t.Errorf("Get after Put-back invoked factory; calls = %d, want 3", calls)
	}
}

// HighWaterMark reports the max in-flight ever observed between FreeAll calls.
func TestPool_HighWaterMarkPersistsAcrossFreeAll(t *testing.T) {
	p := NewPool(8, newGameStateForTest)
	p.Get()
	p.Get()
	p.Get()
	p.FreeAll()
	p.Get()
	if got := p.HighWaterMark(); got != 3 {
		t.Errorf("HighWaterMark = %d, want 3 (peak observed in the first cycle)", got)
	}
}
