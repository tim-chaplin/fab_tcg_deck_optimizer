package optimizations

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// TestWarmChainStepCache_PopulatesBothFromArsenalRows: WarmChainStepCache must fill both
// the (id, false) and (id, true) cells for every non-nil card so the runtime hot path is
// pure reads. Sample a known card and confirm both entries are present and produce the
// expected "<DisplayName>: <VERB>[ from arsenal]" string.
func TestWarmChainStepCache_PopulatesBothFromArsenalRows(t *testing.T) {
	c := testutils.NewStubCard("Test").
		WithID(testutils.FakeRedAttack).
		WithTypes(card.NewTypeSet(card.TypeAttack, card.TypeAction))
	chainStepCache[chainStepCacheIndex(c.ID(), false)].Store(nil)
	chainStepCache[chainStepCacheIndex(c.ID(), true)].Store(nil)

	WarmChainStepCache([]sim.Card{c})

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
// Passing it through must not panic and must leave that slot untouched.
func TestWarmChainStepCache_SkipsNil(t *testing.T) {
	WarmChainStepCache([]sim.Card{nil})
	if got := chainStepCache[0].Load(); got != nil {
		t.Errorf("nil entry should leave slot 0 empty, got %q", *got)
	}
}

// TestChainStepText_LazyBackfillForUnregisteredCards: cachedChainStepText is the runtime
// entry point. A card never seen by WarmChainStepCache (test fakes, ad-hoc stubs) must
// still produce the right string and populate the cache so the next call is a hit.
func TestChainStepText_LazyBackfillForUnregisteredCards(t *testing.T) {
	c := testutils.NewStubCard("Unregistered").
		WithID(testutils.FakeHugeAttack).
		WithTypes(card.NewTypeSet(card.TypeAction))
	idx := chainStepCacheIndex(c.ID(), false)
	chainStepCache[idx].Store(nil)

	self := &sim.CardState{Card: c}
	got := cachedChainStepText(self)
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
		types       card.TypeSet
		fromArsenal bool
		want        string
	}{
		{"weapon", card.NewTypeSet(card.TypeWeapon), false, "X: WEAPON ATTACK"},
		{"attack action", card.NewTypeSet(card.TypeAttack, card.TypeAction), false, "X: ATTACK"},
		{"defense reaction", card.NewTypeSet(card.TypeDefenseReaction), false, "X: DEFENSE REACTION"},
		{"non-attack action", card.NewTypeSet(card.TypeAction), false, "X: PLAY"},
		{"from arsenal suffix", card.NewTypeSet(card.TypeAction), true, "X: PLAY from arsenal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := testutils.NewStubCard("X").WithTypes(tc.types)
			self := &sim.CardState{Card: c, FromArsenal: tc.fromArsenal}
			if got := bareChainStepText(self); got != tc.want {
				t.Errorf("bareChainStepText = %q, want %q", got, tc.want)
			}
		})
	}
}
