package card

import (
	"reflect"
	"testing"
)

// TestEphemeralReset_ZeroesEveryField mutates every Ephemeral field non-zero, calls Reset,
// and asserts each one is back to its type's zero value. A new field added to Ephemeral is
// zeroed for free by *e = Ephemeral{}; this test enforces that structurally so the chain
// runner can never re-introduce the "forgot to clear field X" bug class.
//
// OnHit is checked by truncation length, not deep zero, because Reset preserves its
// backing array to keep per-Best reuse allocation-free.
func TestEphemeralReset_ZeroesEveryField(t *testing.T) {
	e := Ephemeral{
		GrantedGoAgain:   true,
		GrantedDominate:  true,
		GrantedOverpower: true,
		GrantedInstant:   true,
		BonusAttack:      99,
		BonusDefense:     99,
		PitchedToPlay:    []Card{nil},
		OnHit:            []OnHitHandler{{N: 1}},
	}
	e.Reset()

	v := reflect.ValueOf(e)
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		if f.Name == "OnHit" {
			if v.Field(i).Len() != 0 {
				t.Errorf("Reset: OnHit not truncated, len = %d", v.Field(i).Len())
			}
			continue
		}
		if !v.Field(i).IsZero() {
			t.Errorf("Reset: %s not cleared (still %v)", f.Name, v.Field(i).Interface())
		}
	}
}
