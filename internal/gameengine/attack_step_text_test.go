package gameengine

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// attackStepTextFake satisfies card.Card with a configurable id, name, type set, and
// FromArsenal-irrelevant Play. Exercises the verb-selection switch in buildAttackStepText
// and the cache key derived from (id, FromArsenal).
type attackStepTextFake struct {
	id    ids.CardID
	name  string
	types card.TypeSet
}

func (c attackStepTextFake) ID() ids.CardID                                   { return c.id }
func (c attackStepTextFake) Name() string                                     { return c.name }
func (c attackStepTextFake) DisplayName() string                              { return c.name }
func (attackStepTextFake) Cost(card.GameEngine) int                           { return 0 }
func (attackStepTextFake) Pitch() int                                         { return 0 }
func (attackStepTextFake) Attack() int                                        { return 0 }
func (attackStepTextFake) Defense() int                                       { return 0 }
func (c attackStepTextFake) Types(card.GameEngine) card.TypeSet               { return c.types }
func (attackStepTextFake) GoAgain(card.GameEngine) bool                       { return false }
func (attackStepTextFake) Play(card.GameEngine, card.Logger, *card.CardState) {}

// TestAttackStepText_VerbSelection pins each branch of the verb switch so a future type
// reshuffle that breaks one surfaces here rather than inside a downstream golden test.
func TestAttackStepText_VerbSelection(t *testing.T) {
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
				Card:        attackStepTextFake{name: "X", types: tc.types},
				FromArsenal: tc.fromArsenal,
			}
			if got := AttackStepText(pc); got != tc.want {
				t.Errorf("AttackStepText = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestAttackStepText_CacheBackfillsForUnknownIDs: a non-Invalid ID that the cache hasn't
// seen yet is computed once and stored, so the second call is a pure read of the same
// string. Uses a synthetic ID well above the registry range to avoid collisions.
func TestAttackStepText_CacheBackfillsForUnknownIDs(t *testing.T) {
	const syntheticID ids.CardID = 1<<16 - 1
	attackStepCache[attackStepCacheIndex(syntheticID, false)].Store(nil)
	pc := &card.CardState{
		Card: attackStepTextFake{id: syntheticID, name: "Synth", types: card.NewTypeSet(card.TypeAction)},
	}
	first := AttackStepText(pc)
	if want := "Synth: PLAY"; first != want {
		t.Errorf("first call = %q, want %q", first, want)
	}
	cached := attackStepCache[attackStepCacheIndex(syntheticID, false)].Load()
	if cached == nil {
		t.Fatal("first call should backfill the cache")
	}
	if *cached != first {
		t.Errorf("cached entry = %q, want %q", *cached, first)
	}
}
