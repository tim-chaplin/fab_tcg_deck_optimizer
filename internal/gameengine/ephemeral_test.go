package gameengine

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestEphemeralReset_MatchesExpectedShape mutates every ephemeral field to a non-zero
// value, calls reset, and field-by-field asserts the result equals the expected reset
// shape — every field zero except actionPoints=1, currentHookIdx=-1, cacheable=true,
// logger=NoopLogger{}. A new ephemeral field added to the struct is zeroed for free by
// *e = ephemeral{}; this test enforces that contract structurally so a forgotten field
// can't silently leak across turn boundaries.
func TestEphemeralReset_MatchesExpectedShape(t *testing.T) {
	dirty := ephemeral{
		cardsPlayed:           []card.Card{nil},
		cardsRemaining:        []*card.CardState{nil},
		pitched:               []card.Card{nil},
		defenders:             []card.Card{nil},
		triggers:              []EphemeralTrigger{nil},
		attackReactionTarget:  &card.CardState{},
		logger:                &StreamLogger{},
		actionPoints:          99,
		value:                 99,
		damageBlocked:         99,
		blockTotal:            99,
		currentHookIdx:        99,
		cardBanished:          true,
		arcaneDamageDealt:     true,
		auraCreated:           true,
		nonAttackActionPlayed: true,
		lastAttackHit:         true,
		currentHookDestroyed:  true,
		currentStepRerouted:   true,
		cacheable:             false,
		heroTapped:            true,
	}
	dirty.reset()

	want := ephemeral{
		actionPoints:   1,
		currentHookIdx: -1,
		cacheable:      true,
		logger:         NoopLogger{},
	}
	got := reflect.ValueOf(dirty)
	wantV := reflect.ValueOf(want)
	tp := got.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		gotF := got.Field(i)
		// Slices aren't reflect.Value.Equal-comparable; IsZero (true when nil & len 0)
		// is the right check for reset. No slice field has a non-zero reset value, so
		// the want side doesn't need a slice-specific path.
		if gotF.Kind() == reflect.Slice {
			if !gotF.IsZero() {
				t.Errorf("reset: %s slice not cleared (len=%d)", f.Name, gotF.Len())
			}
			continue
		}
		if !gotF.Equal(wantV.Field(i)) {
			t.Errorf("reset: %s mismatch (got != want); add the field to reset's non-zero defaults if it should reset to a non-zero value, otherwise the wholesale zero already covers it", f.Name)
		}
	}
}
