package sim

// Sequential decision for the anneal hill climb: a candidate mutation is judged against its
// incumbent by streaming paired per-shuffle value deltas (d_i = mutantValue_i − incumbentValue_i,
// common random numbers) into a sequential probability ratio test (SPRT) until the evidence is
// decisive, rather than spending a fixed shuffle budget. Most mutations are clearly neutral and
// bail in a handful of shuffles; only the genuinely close ones run long.

import "math"

// sprtConfig holds the test's two knobs. alpha is the symmetric error rate (false-accept of a
// neutral mutation, false-reject of a real improvement); warmup is the minimum sample count
// before any decision, so the online variance estimate has settled.
type sprtConfig struct {
	alpha  float64
	warmup int
}

// defaultSPRTConfig is the anneal default: a 5% symmetric error rate and an 8-shuffle warmup.
var defaultSPRTConfig = sprtConfig{alpha: 0.05, warmup: 8}

// sprtVerdict is the running test's current call.
type sprtVerdict int

const (
	sprtContinue sprtVerdict = iota
	sprtAccept
	sprtReject
)

// sprtAccumulator streams paired per-shuffle deltas through an SPRT for the hypotheses
//
//	H0: mean delta = 0          (neutral — reject the mutation)
//	H1: mean delta = threshold  (a real improvement — accept it)
//
// Anything strictly between 0 and threshold is the indifference zone: the test still terminates
// there, to whichever side the evidence drifts, which is fine — a sub-threshold gain is a
// don't-care by the same logic min-improvement encodes. Expected sample size is bounded for every
// true mean, so no max-shuffle cap is needed.
//
// The delta variance is estimated online (Welford), making this the standard Gaussian SPRT with
// a plug-in variance rather than a known-variance test — accurate enough at the few-dozen-sample
// scale these decisions resolve at, and it needs no prior on the shuffle noise.
type sprtAccumulator struct {
	n    int
	mean float64
	m2   float64 // running sum of squared deviations; variance = m2/(n-1)
}

// add folds one paired delta into the running mean and variance.
func (s *sprtAccumulator) add(delta float64) {
	s.n++
	d := delta - s.mean
	s.mean += d / float64(s.n)
	s.m2 += d * (delta - s.mean)
}

// variance returns the sample variance of the deltas seen so far, 0 before two samples.
func (s *sprtAccumulator) variance() float64 {
	if s.n < 2 {
		return 0
	}
	return s.m2 / float64(s.n-1)
}

// decision applies the SPRT for H0: mean=0 vs H1: mean=threshold at symmetric error rate
// cfg.alpha, returning sprtContinue until the log-likelihood ratio crosses a boundary or
// sprtAccept / sprtReject once it does. No decision is made before cfg.warmup samples.
func (s *sprtAccumulator) decision(threshold float64, cfg sprtConfig) sprtVerdict {
	if s.n < cfg.warmup {
		return sprtContinue
	}
	sigma2 := s.variance()
	if sigma2 <= 0 {
		// No observed noise: the mean is the exact delta. Decide on the indifference zone's
		// midpoint — the zero-variance limit of the ratio below.
		if s.mean > threshold/2 {
			return sprtAccept
		}
		return sprtReject
	}
	// Boundary on the log-likelihood ratio for the symmetric error rate; ±this is the Wald
	// A / B pair when alpha == beta.
	boundary := math.Log((1 - cfg.alpha) / cfg.alpha)
	// LLR after n iid samples with mean d̄ for N(0,σ²) vs N(threshold,σ²):
	//   (threshold/σ²) · n · (d̄ − threshold/2)
	llr := (threshold / sigma2) * float64(s.n) * (s.mean - threshold/2)
	switch {
	case llr >= boundary:
		return sprtAccept
	case llr <= -boundary:
		return sprtReject
	default:
		return sprtContinue
	}
}
