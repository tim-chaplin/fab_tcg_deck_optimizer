package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// chainStepTextFake satisfies card.Card with a configurable id, name, type set, and
// FromArsenal-irrelevant Play. Exercises the verb-selection switch in buildChainStepText
// and the cache key derived from (id, FromArsenal).
type chainStepTextFake struct {
	id    ids.CardID
	name  string
	types card.TypeSet
}

func (c chainStepTextFake) ID() ids.CardID                                   { return c.id }
func (c chainStepTextFake) Name() string                                     { return c.name }
func (c chainStepTextFake) DisplayName() string                              { return c.name }
func (chainStepTextFake) Cost(card.GameEngine) int                           { return 0 }
func (chainStepTextFake) Pitch() int                                         { return 0 }
func (chainStepTextFake) Attack() int                                        { return 0 }
func (chainStepTextFake) Defense() int                                       { return 0 }
func (c chainStepTextFake) Types(card.GameEngine) card.TypeSet               { return c.types }
func (chainStepTextFake) GoAgain(card.GameEngine) bool                       { return false }
func (chainStepTextFake) Play(card.GameEngine, card.Logger, *card.CardState) {}

// TestChainStepText_VerbSelection pins each branch of the verb switch so a future type
// reshuffle that breaks one surfaces here rather than inside a downstream golden test.
func TestChainStepText_VerbSelection(t *testing.T) {
	cases := []struct {
		name        string
		types       card.TypeSet
		fromArsenal bool
		want        string
	}{
		{"weapon ability", card.NewTypeSet(card.TypeWeapon, card.TypeAttack), false, "X: WEAPON ATTACK"},
		{"attack action", card.NewTypeSet(card.TypeAction, card.TypeAttack), false, "X: ATTACK"},
		{"defense reaction", card.NewTypeSet(card.TypeDefenseReaction), false, "X: DEFENSE REACTION"},
		{"non-attack action", card.NewTypeSet(card.TypeAction), false, "X: PLAY"},
		{"from arsenal suffix", card.NewTypeSet(card.TypeAction), true, "X: PLAY from arsenal"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// ids.InvalidCard routes through the no-cache path, so each case starts fresh
			// without depending on prior test ordering.
			pc := &card.CardState{
				Card:        chainStepTextFake{name: "X", types: tc.types},
				FromArsenal: tc.fromArsenal,
			}
			if got := ChainStepText(pc); got != tc.want {
				t.Errorf("ChainStepText = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestChainStepText_CacheBackfillsForUnknownIDs: a non-Invalid ID that the cache hasn't
// seen yet is computed once and stored, so the second call is a pure read of the same
// string. Uses a synthetic ID well above the registry range to avoid collisions.
func TestChainStepText_CacheBackfillsForUnknownIDs(t *testing.T) {
	const syntheticID ids.CardID = 1<<16 - 1
	chainStepCache[chainStepCacheIndex(syntheticID, false)].Store(nil)
	pc := &card.CardState{
		Card: chainStepTextFake{id: syntheticID, name: "Synth", types: card.NewTypeSet(card.TypeAction)},
	}
	first := ChainStepText(pc)
	if want := "Synth: PLAY"; first != want {
		t.Errorf("first call = %q, want %q", first, want)
	}
	cached := chainStepCache[chainStepCacheIndex(syntheticID, false)].Load()
	if cached == nil {
		t.Fatal("first call should backfill the cache")
	}
	if *cached != first {
		t.Errorf("cached entry = %q, want %q", *cached, first)
	}
}
