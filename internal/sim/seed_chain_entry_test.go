package sim

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// fakeCardForSeed is a minimal Card implementation used solely to seed a CardState — its
// behavior is never exercised; we only need a non-nil Card value the reset can write.
type fakeCardForSeed struct{}

func (fakeCardForSeed) ID() ids.CardID                                     { return ids.InvalidCard }
func (fakeCardForSeed) Name() string                                       { return "fakeCardForSeed" }
func (fakeCardForSeed) DisplayName() string                                { return "fakeCardForSeed" }
func (fakeCardForSeed) Cost(card.GameEngine) int                           { return 0 }
func (fakeCardForSeed) Pitch() int                                         { return 0 }
func (fakeCardForSeed) Attack() int                                        { return 0 }
func (fakeCardForSeed) Defense() int                                       { return 0 }
func (fakeCardForSeed) Types(card.GameEngine) card.TypeSet                 { return 0 }
func (fakeCardForSeed) GoAgain(card.GameEngine) bool                       { return false }
func (fakeCardForSeed) Play(card.GameEngine, card.Logger, *card.CardState) {}

// seedChainEntryAllowlist names the CardState fields seedChainEntry binds to per-permutation
// values rather than zeroing — i.e. fields that legitimately carry information across a
// chain reset. Every OTHER field on CardState must be back to its zero value after the
// reset; if you're adding a new field, decide which list it belongs in:
//
//   - Add it here if it carries chain-binding data (like Card / FromArsenal).
//   - Otherwise it MUST be cleared by seedChainEntry — and by the per-permutation reset
//     loop inside playSequenceWithMeta. Leaving a per-permutation field unreset leaks state
//     between bestSequence permutations (see internal/sim/permutation_reset_test.go for the
//     pattern of bugs that catches).
//
// OnHit is also exempt — seedChainEntry resets it via `pc.OnHit = pc.OnHit[:0]` so the
// underlying backing array is preserved across Best calls. The test treats a length-0
// slice as "reset" regardless of capacity.
var seedChainEntryAllowlist = map[string]bool{
	"Card":        true,
	"FromArsenal": true,
}

// TestSeedChainEntry_ResetsEveryPerPermutationField guards against the footgun that
// motivated permutation_reset_test.go: a new mutable CardState field is added, the
// developer remembers seedChainEntry but misses the partial reset inside
// playSequenceWithMeta (or vice versa) — leaving the field to leak across permutations.
//
// The test mutates every CardState field to a non-zero value, runs seedChainEntry, then
// uses reflection to assert each non-allowlisted field is back to its type's zero value.
// New CardState fields force the test to fail until they're either added to seedChainEntry
// or added to the allowlist with a comment explaining why.
func TestSeedChainEntry_ResetsEveryPerPermutationField(t *testing.T) {
	pc := card.CardState{
		Card:             fakeCardForSeed{},
		GrantedGoAgain:   true,
		GrantedDominate:  true,
		GrantedOverpower: true,
		FromArsenal:      true,
		Mode:             5,
		BonusAttack:      99,
		BonusDefense:     99,
		PitchedToPlay:    []card.Card{fakeCardForSeed{}},
		OnHit:            []card.OnHitHandler{{N: 1}},
	}
	ctx := &sequenceContext{arsenalInIdx: -1}
	ctx.seedChainEntry(&pc, fakeCardForSeed{}, 0)

	v := reflect.ValueOf(pc)
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if seedChainEntryAllowlist[f.Name] {
			continue
		}
		// OnHit is reset via slice-truncation, not nil-assignment, so check length rather
		// than IsZero (which is false for a length-0 slice with non-zero capacity).
		if f.Name == "OnHit" {
			if v.Field(i).Len() != 0 {
				t.Errorf("seedChainEntry: OnHit not truncated, len = %d", v.Field(i).Len())
			}
			continue
		}
		if !v.Field(i).IsZero() {
			t.Errorf("seedChainEntry: %s not reset (still %v) — add it to seedChainEntry's reset (and to playSequenceWithMeta's per-permutation reset) or to seedChainEntryAllowlist if it's intentionally preserved", f.Name, v.Field(i).Interface())
		}
	}
}
