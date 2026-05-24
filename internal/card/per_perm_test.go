package card

import (
	"reflect"
	"testing"
)

// TestPerPermReset_ZeroesEveryField mutates every PerPerm field to a non-zero value,
// calls Reset, and asserts each one is back to its type's zero value. New mutable per-
// permutation state added inside PerPerm gets zeroed for free by *p = PerPerm{}; this
// test enforces that contract structurally so the chain runner can never re-introduce
// the "forgot to clear field X between permutations" bug class.
//
// OnHit is checked by truncation length, not deep zero, because Reset preserves its
// backing array to keep per-Best reuse allocation-free.
func TestPerPermReset_ZeroesEveryField(t *testing.T) {
	p := PerPerm{
		GrantedGoAgain:   true,
		GrantedDominate:  true,
		GrantedOverpower: true,
		GrantedInstant:   true,
		BonusAttack:      99,
		BonusDefense:     99,
		PitchedToPlay:    []Card{nil},
		OnHit:            []OnHitHandler{{N: 1}},
	}
	p.Reset()

	v := reflect.ValueOf(p)
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
