package gameengine

import (
	"reflect"
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// TestEphemeralReset_MatchesExpectedShape mutates every ephemeral field to a non-zero
// value, calls reset, and field-by-field asserts the result equals the expected reset
// shape — scalars zero except actionPoints=1, currentHookIdx=-1, cacheable=true,
// logger=NoopLogger{}; slice fields truncated to length 0 but with their backings
// preserved (so per-perm appends like AppendCardsPlayed reuse the cap instead of
// allocating fresh each call).
func TestEphemeralReset_MatchesExpectedShape(t *testing.T) {
	dirty := ephemeral{
		cardsPlayed:           []card.Card{nil, nil, nil},
		cardsRemaining:        []*card.CardState{nil, nil},
		pitched:               []*card.CardState{nil, nil},
		defenders:             []card.Card{nil, nil},
		triggers:              []EphemeralTrigger{nil, nil},
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
	wantCaps := map[string]int{
		"cardsPlayed":    cap(dirty.cardsPlayed),
		"cardsRemaining": cap(dirty.cardsRemaining),
		"pitched":        cap(dirty.pitched),
		"defenders":      cap(dirty.defenders),
		"triggers":       cap(dirty.triggers),
	}
	dirty.reset()

	want := ephemeral{
		actionPoints:           1,
		currentHookIdx:         -1,
		currentFiringTokenAura: -1,
		currentFiringTokenItem: -1,
		cacheable:              true,
		logger:                 NoopLogger{},
	}
	got := reflect.ValueOf(dirty)
	wantV := reflect.ValueOf(want)
	tp := got.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		gotF := got.Field(i)
		if gotF.Kind() == reflect.Slice {
			if gotF.Len() != 0 {
				t.Errorf("reset: %s slice len = %d, want 0", f.Name, gotF.Len())
			}
			if wantCap, ok := wantCaps[f.Name]; ok && gotF.Cap() != wantCap {
				t.Errorf("reset: %s slice cap = %d, want %d (backing should be preserved across reset for cap reuse)", f.Name, gotF.Cap(), wantCap)
			}
			continue
		}
		if !gotF.Equal(wantV.Field(i)) {
			t.Errorf("reset: %s mismatch (got != want); add the field to reset's non-zero defaults if it should reset to a non-zero value, otherwise the wholesale zero already covers it", f.Name)
		}
	}
}
