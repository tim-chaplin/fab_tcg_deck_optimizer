package optimizations

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Tests that WarmChainStepCache fills both the (id, false) and (id, true) cells per card
// so the runtime hot path is pure reads. Cards return ids.InvalidCard so both rows land in
// slot 0 (in-hand) and slot 1<<16 (from-arsenal); the test asserts both rows are populated.
func TestWarmChainStepCache_PopulatesBothFromArsenalRows(t *testing.T) {
	c := testutils.FakeRedAttack().WithName("Test")
	chainStepCache[chainStepCacheIndex(c.ID(), false)].Store(nil)
	chainStepCache[chainStepCacheIndex(c.ID(), true)].Store(nil)

	WarmChainStepCache([]card.Card{c})

	gotInHand := chainStepCache[chainStepCacheIndex(c.ID(), false)].Load()
	if gotInHand == nil {
		t.Fatal("(id, false) not populated after warm")
	}
	if want := "Test: ATTACK"; *gotInHand != want {
		t.Errorf("(id, false) = %q, want %q", *gotInHand, want)
	}

	gotArsenal := chainStepCache[chainStepCacheIndex(c.ID(), true)].Load()
	if gotArsenal == nil {
		t.Fatal("(id, true) not populated after warm")
	}
	if want := "Test: ATTACK from arsenal"; *gotArsenal != want {
		t.Errorf("(id, true) = %q, want %q", *gotArsenal, want)
	}
}

// TestWarmChainStepCache_SkipsNil: the registry slice has nil at index 0 (Invalid).
// Passing it through must not panic and must leave that slot untouched. Slot 0 is shared
// with every InvalidCard fake so we clear it first to isolate the assertion.
func TestWarmChainStepCache_SkipsNil(t *testing.T) {
	chainStepCache[0].Store(nil)
	WarmChainStepCache([]card.Card{nil})
	if got := chainStepCache[0].Load(); got != nil {
		t.Errorf("nil entry should leave slot 0 empty, got %q", *got)
	}
}

// TestChainStepText_LazyBackfillForUnregisteredCards: cachedChainStepText is the runtime
// entry point. A card never seen by WarmChainStepCache (test fakes, ad-hoc stubs) must
// still produce the right string and populate the cache so the next call is a hit.
func TestChainStepText_LazyBackfillForUnregisteredCards(t *testing.T) {
	c := testutils.FakeRedAction().WithName("Unregistered")
	idx := chainStepCacheIndex(c.ID(), false)
	chainStepCache[idx].Store(nil)

	pc := &card.CardState{Card: c}
	got := cachedChainStepText(pc)
	if want := "Unregistered: PLAY"; got != want {
		t.Errorf("first call = %q, want %q", got, want)
	}
	cached := chainStepCache[idx].Load()
	if cached == nil || *cached != got {
		t.Error("first call should backfill the cache")
	}
}

// TestBuildChainStepText_VerbSelection: the verb-selection switch covers the four
// type buckets the chain-step renderer routes through. Pin each branch so a future type
// reshuffle that breaks one is caught here rather than inside a downstream golden test.
func TestBuildChainStepText_VerbSelection(t *testing.T) {
	cases := []struct {
		name        string
		card        testutils.Fake
		fromArsenal bool
		want        string
	}{
		{"weapon ability", testutils.FakeWeaponSwing().WithName("X"), false, "X: WEAPON ATTACK"},
		{"attack action", testutils.FakeRedAttack().WithName("X"), false, "X: ATTACK"},
		{"defense reaction", testutils.FakeRedDR().WithName("X"), false, "X: DEFENSE REACTION"},
		{"non-attack action", testutils.FakeRedAction().WithName("X"), false, "X: PLAY"},
		{"from arsenal suffix", testutils.FakeRedAction().WithName("X"), true, "X: PLAY from arsenal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			pc := &card.CardState{Card: tc.card, FromArsenal: tc.fromArsenal}
			if got := bareChainStepText(pc); got != tc.want {
				t.Errorf("bareChainStepText = %q, want %q", got, tc.want)
			}
		})
	}
}
