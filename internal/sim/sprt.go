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

// defaultSPRTConfig is the anneal default: a 20% symmetric error rate and an 8-shuffle warmup.
// alpha is tuned for minimum compute, not minimum error. On a converged round the per-candidate
// cost is (screen shuffles + P(false-accept)·N confirm), and a higher alpha buys cheaper screens
// at the price of more — confirm-filtered — false-accepts. That total bottoms out across a broad
// alpha ≈ 0.2–0.33 basin (≈1.7× cheaper than alpha=0.05, σ- and N-independent); 0.20 sits at the
// low end, holding the screen's false-reject of a marginal winner to 20% rather than 25–33%. The
// confirm's k-sigma gate, not alpha, decides which accepts actually win — so a loose screen costs
// compute, never correctness; at worst it misses a marginal gain that re-screens next round.
var defaultSPRTConfig = sprtConfig{alpha: 0.20, warmup: 8}

// defaultConfirmSigmas is the accept threshold's width in confirm standard errors: a hill-climb
// mutation must beat the incumbent by more than k·σ̂/√N — k sigma of the N-shuffle confirm — to be
// accepted, i.e. the confirm's interval for the gain excludes zero at this confidence. With no
// min-improvement, this is the sole thing flooring the resolution; the only user knob left is N
// (-shuffles), which sets how small a k-sigma gain is worth chasing.
const defaultConfirmSigmas = 3.0

// sprtVerdict is the running test's current call.
type sprtVerdict int

const (
	sprtContinue sprtVerdict = iota
	sprtAccept
	sprtReject
)

// sprtAccumulator streams paired per-shuffle deltas through an SPRT. decision tests the running
// mean against an accept point tau with an indifference-zone width threshold:
//
//	H0: mean delta = tau − threshold  (reject the mutation)
//	H1: mean delta = tau              (accept it)
//
// Anything strictly inside the zone is a don't-care: the test still terminates there, to whichever
// side the evidence drifts. Expected sample size is bounded for every true mean, so no max-shuffle
// cap is needed. The caller derives threshold per candidate as k·σ̂/√N — the confirm's k-sigma
// resolution — rather than passing a fixed min-improvement, so the gain it can resolve scales with
// the confirm budget N.
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

// decision applies the SPRT for H0: mean = tau−threshold (reject) vs H1: mean = tau (accept) at
// symmetric error rate cfg.alpha, returning sprtContinue until the log-likelihood ratio crosses a
// boundary or sprtAccept / sprtReject once it does. threshold is the indifference-zone width (the
// derived k·σ̂/√N resolution); tau the accept point. tau == threshold reduces to the classic
// H0=0 / H1=threshold test. No decision is made before cfg.warmup samples.
func (s *sprtAccumulator) decision(tau, threshold float64, cfg sprtConfig) sprtVerdict {
	if s.n < cfg.warmup {
		return sprtContinue
	}
	sigma2 := s.variance()
	if sigma2 <= 0 {
		// No observed noise: the mean is the exact delta. Decide on the indifference zone's
		// midpoint, tau − threshold/2 — the zero-variance limit of the ratio below.
		if s.mean > tau-threshold/2 {
			return sprtAccept
		}
		return sprtReject
	}
	// Boundary on the log-likelihood ratio for the symmetric error rate; ±this is the Wald
	// A / B pair when alpha == beta.
	boundary := math.Log((1 - cfg.alpha) / cfg.alpha)
	// LLR after n iid samples with mean d̄ for N(tau−threshold,σ²) vs N(tau,σ²):
	//   (threshold/σ²) · n · (d̄ − tau + threshold/2)
	llr := (threshold / sigma2) * float64(s.n) * (s.mean - tau + threshold/2)
	switch {
	case llr >= boundary:
		return sprtAccept
	case llr <= -boundary:
		return sprtReject
	default:
		return sprtContinue
	}
}
