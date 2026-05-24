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

// seedChainEntryAllowlist names the CardState top-level fields seedChainEntry binds to
// per-permutation values rather than zeroing — i.e. fields that legitimately carry
// information across a chain reset.
//
//   - Card / FromArsenal / Mode: chain-binding identity. Mode is reseeded per modal tuple
//     by the chain runner's enumeration loop.
//   - PerPerm: the embedded scratch struct. Its Reset method owns the per-field zeroing
//     contract and is exercised by TestPerPermReset_ZeroesEveryField in package card.
//   - Role: hand-state field, never touched by chain pcBuf entries; its zero default is
//     left in place by seedChainEntry.
var seedChainEntryAllowlist = map[string]bool{
	"Card":        true,
	"FromArsenal": true,
	"Mode":        true,
	"Role":        true,
	"PerPerm":     true,
}

// TestSeedChainEntry_TopLevelFieldsAllAccountedFor guards against a new top-level
// CardState field slipping in without being either an allowlisted binding or covered by
// PerPerm.Reset. New per-permutation fields belong inside PerPerm so PerPerm.Reset zeroes
// them automatically; new binding fields belong in the allowlist with a comment.
func TestSeedChainEntry_TopLevelFieldsAllAccountedFor(t *testing.T) {
	v := reflect.ValueOf(card.CardState{})
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if !seedChainEntryAllowlist[f.Name] {
			t.Errorf("CardState field %q is neither in seedChainEntryAllowlist nor covered by the PerPerm embedded reset — decide which group it belongs to and update accordingly", f.Name)
		}
	}
}
