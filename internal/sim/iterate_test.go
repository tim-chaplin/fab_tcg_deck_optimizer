package sim

import (
	"math/rand"
	"testing"
)

// TestAcceptMutation_StrictGateRequiresMinImprovement: at T==0, a mutation with deepAvg
// equal to or barely above bestAvg fails the gate when minImprovement > 0. Only deepAvg >
// bestAvg + minImprovement passes. Pins the noise-floor guard against shuffle-noise wins
// that would let the hill climb loop forever on near-zero "improvements."
func TestAcceptMutation_StrictGateRequiresMinImprovement(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	const baseline = 10.0
	const minImp = 0.1

	cases := []struct {
		name     string
		deepAvg  float64
		wantPass bool
	}{
		{"equal to baseline", baseline, false},
		{"baseline + 0.05 (under floor)", baseline + 0.05, false},
		{"baseline + minImp exactly (on floor — not strictly above)", baseline + minImp, false},
		{"baseline + minImp + epsilon (just above floor)", baseline + minImp + 0.001, true},
		{"baseline + 1.0 (well above)", baseline + 1.0, true},
	}
	for _, tc := range cases {
		got := acceptMutation(tc.deepAvg, baseline, 0, minImp, rng)
		if got != tc.wantPass {
			t.Errorf("%s: acceptMutation(%.4f, %.1f, T=0, minImp=%.2f) = %v, want %v",
				tc.name, tc.deepAvg, baseline, minImp, got, tc.wantPass)
		}
	}
}

// TestAcceptMutation_AnnealingBypassesMinImprovement: at T>0 the probabilistic gate ignores
// minImprovement so SA can still walk through ties / shallow dips to escape local maxima.
// Verifies by running many trials at deepAvg == baseline (a tie) with a substantial T —
// acceptances should occur regardless of how high minImprovement is set, because the gate
// math reduces to exp(0/T) = 1.
func TestAcceptMutation_AnnealingBypassesMinImprovement(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	const baseline = 10.0
	const minImp = 5.0 // intentionally huge — would block every realistic strict-gate accept.
	accepts := 0
	const trials = 100
	for i := 0; i < trials; i++ {
		// deepAvg == baseline → exp(0 / 1.0) = 1.0 → always passes the probabilistic gate.
		if acceptMutation(baseline, baseline, 1.0, minImp, rng) {
			accepts++
		}
	}
	if accepts != trials {
		t.Errorf("at deepAvg == baseline with T=1.0, accepts = %d/%d, want all (probabilistic gate ignores minImprovement)",
			accepts, trials)
	}
}

// TestAcceptMutation_ZeroMinImprovementAllowsAnyStrictImprovement: with minImprovement=0
// the floor is disabled and any strictly-greater deepAvg passes at T==0; equal deepAvg still
// fails (strict >, not >=). Pins the documented "pass 0 to disable the floor" contract.
func TestAcceptMutation_ZeroMinImprovementAllowsAnyStrictImprovement(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	if !acceptMutation(10.0001, 10.0, 0, 0, rng) {
		t.Error("minImprovement=0: deepAvg = baseline + tiny epsilon should pass at T==0")
	}
	if acceptMutation(10.0, 10.0, 0, 0, rng) {
		t.Error("minImprovement=0: deepAvg == baseline should still fail at T==0 (strict >, not >=)")
	}
}
