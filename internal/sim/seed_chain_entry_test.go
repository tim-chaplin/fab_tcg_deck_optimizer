package sim

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// seedChainEntryAllowlist names the CardState top-level fields seedChainEntry binds to
// per-permutation values rather than zeroing — i.e. fields that legitimately carry
// information across a chain reset.
//
//   - Card / FromArsenal / Mode: chain-binding identity. Mode is reseeded per modal tuple
//     by the chain runner's enumeration loop.
//   - FromDraw: set only by mid-chain DrawOne inserts; pcBuf entries are planned
//     attackers, so this stays zero.
//   - Ephemeral: the embedded short-lived state. Its Reset method owns the per-field
//     zeroing contract and is exercised by TestEphemeralReset_ZeroesEveryField in package
//     card.
//   - Role: hand-state field, never touched by chain pcBuf entries; its zero default is
//     left in place by seedChainEntry.
var seedChainEntryAllowlist = map[string]bool{
	"Card":        true,
	"FromArsenal": true,
	"FromDraw":    true,
	"Mode":        true,
	"Role":        true,
	"Ephemeral":   true,
}

// TestSeedChainEntry_TopLevelFieldsAllAccountedFor guards against a new top-level
// CardState field slipping in without being either an allowlisted binding or covered by
// Ephemeral.Reset. New mutable fields belong inside Ephemeral so Ephemeral.Reset zeroes
// them automatically; new binding fields belong in the allowlist with a comment.
func TestSeedChainEntry_TopLevelFieldsAllAccountedFor(t *testing.T) {
	v := reflect.ValueOf(card.CardState{})
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if !seedChainEntryAllowlist[f.Name] {
			t.Errorf("CardState field %q is neither in seedChainEntryAllowlist nor covered by the Ephemeral embedded reset — decide which group it belongs to and update accordingly", f.Name)
		}
	}
}
