package sim

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// seedAttackStepEntryAllowlist names the CardState top-level fields seedAttackStepEntry binds to
// per-permutation values rather than zeroing — i.e. fields that legitimately carry
// information across an attack turn reset.
//
//   - Card / FromArsenal / Mode: attack-step-binding identity. Mode is reseeded per modal tuple
//     by the attack-turn runner's enumeration loop.
//   - FromDraw: set only by mid-attack-turn DrawOne inserts; pcBuf entries are planned
//     attackers, so this stays zero.
//   - Ephemeral: the embedded short-lived state. Its Reset method owns the per-field
//     zeroing contract and is exercised by TestEphemeralReset_ZeroesEveryField in package
//     card.
//   - Role: hand-state field, never touched by attack-step pcBuf entries; its zero default is
//     left in place by seedAttackStepEntry.
var seedAttackStepEntryAllowlist = map[string]bool{
	"Card":        true,
	"FromArsenal": true,
	"FromDraw":    true,
	"Mode":        true,
	"Role":        true,
	"Ephemeral":   true,
}

// TestSeedAttackStepEntry_TopLevelFieldsAllAccountedFor guards against a new top-level
// CardState field slipping in without being either an allowlisted binding or covered by
// Ephemeral.Reset. New mutable fields belong inside Ephemeral so Ephemeral.Reset zeroes
// them automatically; new binding fields belong in the allowlist with a comment.
func TestSeedAttackStepEntry_TopLevelFieldsAllAccountedFor(t *testing.T) {
	v := reflect.ValueOf(card.CardState{})
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if !seedAttackStepEntryAllowlist[f.Name] {
			t.Errorf("CardState field %q is neither in seedAttackStepEntryAllowlist nor covered by the Ephemeral embedded reset — decide which group it belongs to and update accordingly", f.Name)
		}
	}
}
