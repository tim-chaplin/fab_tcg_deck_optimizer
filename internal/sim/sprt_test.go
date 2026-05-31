package sim

import (
	"math"
	"math/rand"
	"testing"
)

// runSPRT streams iid N(mean, sd^2) deltas through the accumulator until it decides, capped only
// by a generous safety limit so a bug can't hang the test. Returns the verdict and sample count.
func runSPRT(t *testing.T, mean, sd, threshold float64, seed int64) (sprtVerdict, int) {
	t.Helper()
	rng := rand.New(rand.NewSource(seed))
	var acc sprtAccumulator
	for n := 0; n < 1_000_000; n++ {
		acc.add(mean + sd*rng.NormFloat64())
		if v := acc.decision(threshold, threshold, defaultSPRTConfig); v != sprtContinue {
			return v, acc.n
		}
	}
	t.Fatalf("SPRT did not terminate within the safety limit (mean=%.3f sd=%.3f thr=%.3f)", mean, sd, threshold)
	return sprtContinue, 0
}

// TestSPRT_AcceptsClearImprovement: a stream centred well above the threshold accepts.
func TestSPRT_AcceptsClearImprovement(t *testing.T) {
	const threshold = 0.1
	for seed := int64(1); seed <= 20; seed++ {
		v, n := runSPRT(t, 0.5, 0.5, threshold, seed)
		if v != sprtAccept {
			t.Errorf("seed %d: mean=0.5 (>> %.2f) got verdict %v after %d samples, want accept", seed, threshold, v, n)
		}
	}
}

// TestSPRT_RejectsClearlyWorse: a stream centred below zero rejects every time.
func TestSPRT_RejectsClearlyWorse(t *testing.T) {
	const threshold = 0.1
	for seed := int64(1); seed <= 20; seed++ {
		v, n := runSPRT(t, -0.2, 0.5, threshold, seed)
		if v != sprtReject {
			t.Errorf("seed %d: mean=-0.2 got verdict %v after %d samples, want reject", seed, v, n)
		}
	}
}

// TestSPRT_ErrorRatesControlled checks the empirical false-accept rate at H0 and false-reject rate
// at H1 both stay within a loose bound above the 0.05 nominal (tolerating the plug-in variance).
func TestSPRT_ErrorRatesControlled(t *testing.T) {
	const (
		threshold = 0.1
		trials    = 400
		maxRate   = 0.15
	)
	falseAccept, falseReject := 0, 0
	for seed := int64(1); seed <= trials; seed++ {
		if v, _ := runSPRT(t, 0.0, 0.5, threshold, seed); v == sprtAccept {
			falseAccept++
		}
		if v, _ := runSPRT(t, threshold, 0.5, threshold, seed+trials); v == sprtReject {
			falseReject++
		}
	}
	if rate := float64(falseAccept) / trials; rate > maxRate {
		t.Errorf("false-accept rate at H0 = %.3f (%d/%d), want <= %.2f", rate, falseAccept, trials, maxRate)
	}
	if rate := float64(falseReject) / trials; rate > maxRate {
		t.Errorf("false-reject rate at H1 = %.3f (%d/%d), want <= %.2f", rate, falseReject, trials, maxRate)
	}
}

// TestSPRT_TerminatesAtBoundary: even a stream sitting exactly in the indifference zone
// (mean = threshold/2) terminates — the property that lets the anneal run without a cap.
func TestSPRT_TerminatesAtBoundary(t *testing.T) {
	const threshold = 0.1
	for seed := int64(1); seed <= 20; seed++ {
		v, n := runSPRT(t, threshold/2, 0.5, threshold, seed)
		if v == sprtContinue {
			t.Errorf("seed %d: indifference-zone stream never terminated", seed)
		}
		if n > 100_000 {
			t.Errorf("seed %d: indifference-zone stream took %d samples; expected bounded", seed, n)
		}
	}
}

// TestSPRT_NoDecisionBeforeWarmup: no verdict is returned before the warmup sample count, even
// for an overwhelming stream.
func TestSPRT_NoDecisionBeforeWarmup(t *testing.T) {
	var acc sprtAccumulator
	for i := 0; i < defaultSPRTConfig.warmup-1; i++ {
		acc.add(100.0) // wildly above any threshold
		if v := acc.decision(0.1, 0.1, defaultSPRTConfig); v != sprtContinue {
			t.Fatalf("decided %v at sample %d, before warmup=%d", v, acc.n, defaultSPRTConfig.warmup)
		}
	}
}

// TestSPRT_ZeroVarianceExactDelta: with no observed noise the mean is the exact delta, decided on
// the indifference-zone midpoint.
func TestSPRT_ZeroVarianceExactDelta(t *testing.T) {
	const threshold = 0.1
	cases := []struct {
		delta float64
		want  sprtVerdict
	}{
		{delta: 0.0, want: sprtReject},
		{delta: 0.04, want: sprtReject}, // below midpoint
		{delta: 0.08, want: sprtAccept}, // above midpoint, still sub-threshold (don't-care)
		{delta: 0.5, want: sprtAccept},
	}
	for _, c := range cases {
		var acc sprtAccumulator
		for i := 0; i < defaultSPRTConfig.warmup; i++ {
			acc.add(c.delta)
		}
		if acc.variance() != 0 {
			t.Fatalf("delta=%.3f: expected zero variance, got %v", c.delta, acc.variance())
		}
		if got := acc.decision(threshold, threshold, defaultSPRTConfig); got != c.want {
			t.Errorf("delta=%.3f: verdict %v, want %v", c.delta, got, c.want)
		}
	}
}

// TestSPRT_TauShiftTestsAgainstTau pins the tau/threshold decision: testing raw deltas against
// H0=tau−threshold / H1=tau decides "is mean(delta) > tau?", including a negative tau (an SA step
// accepting a downhill move).
func TestSPRT_TauShiftTestsAgainstTau(t *testing.T) {
	const threshold = 0.1
	cases := []struct {
		dv, tau float64
		want    sprtVerdict
	}{
		{dv: 0.0, tau: -0.5, want: sprtAccept},  // downhill-but-above-tau accepts (SA)
		{dv: -1.0, tau: -0.5, want: sprtReject}, // below tau rejects
		{dv: 0.5, tau: 0.1, want: sprtAccept},   // clear hill-climb improvement
		{dv: -0.2, tau: 0.1, want: sprtReject},  // worse than the hill-climb threshold
	}
	for _, c := range cases {
		rng := rand.New(rand.NewSource(1))
		var acc sprtAccumulator
		got := sprtContinue
		for n := 0; n < 1_000_000 && got == sprtContinue; n++ {
			acc.add(c.dv + 0.5*rng.NormFloat64())
			got = acc.decision(c.tau, threshold, defaultSPRTConfig)
		}
		if got != c.want {
			t.Errorf("dv=%.2f tau=%.2f: verdict %v after %d samples, want %v", c.dv, c.tau, got, acc.n, c.want)
		}
	}
}

// TestSPRT_WelfordMatchesNaive cross-checks the streaming mean/variance against a two-pass
// computation, since the decision math leans on both.
func TestSPRT_WelfordMatchesNaive(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	var acc sprtAccumulator
	var xs []float64
	for i := 0; i < 500; i++ {
		x := 3 + 2*rng.NormFloat64()
		xs = append(xs, x)
		acc.add(x)
	}
	var sum float64
	for _, x := range xs {
		sum += x
	}
	mean := sum / float64(len(xs))
	var ss float64
	for _, x := range xs {
		ss += (x - mean) * (x - mean)
	}
	wantVar := ss / float64(len(xs)-1)
	if math.Abs(acc.mean-mean) > 1e-9 {
		t.Errorf("mean: streaming %.12f vs naive %.12f", acc.mean, mean)
	}
	if math.Abs(acc.variance()-wantVar)/wantVar > 1e-9 {
		t.Errorf("variance: streaming %.12f vs naive %.12f", acc.variance(), wantVar)
	}
}
