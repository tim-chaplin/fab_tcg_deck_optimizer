package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Damage-equivalent tuning constants. The engine owns these — cards reach the values
// through GameEngine accessors so the model's calibration lives in one place.

// DiscardValue is the damage-equivalent credited when the opponent is forced to discard one
// card. A typical FaB card is worth ~3 points of tempo.
const DiscardValue = 3

// GoldTokenValue is the placeholder credit for NotImplemented Gold-creating cards. Zero
// because Gold pays out on activation (see Gold token ability), not at creation.
const GoldTokenValue = 0

// FrailtyValue / InertiaValue / BloodrotPoxValue are the damage-equivalents credited when
// the matching status token is created under the opponent's control. We don't track
// opposing status-token state — these are flat heuristic stand-ins for the future-turn
// tempo each token costs the opponent.
//
// Inertia: 3, somewhat optimistic — the opponent has to skip arsenal at end of turn, which
// generally costs them a card slot, but doesn't always cash out.
// Bloodrot Pox: 2, assumes the opponent chooses to take the damage rather than its other
// in-game options.
// Frailty: 2, middle-of-the-road — worth more against opponents that attack 3+ times,
// less against opponents that attack 0-1 times.
const (
	FrailtyValue     = 2
	InertiaValue     = 3
	BloodrotPoxValue = 2
)

// LikelyToHit reports whether pc's attack is likely to land past the opponent's blocks.
// Folds pc.EffectiveAttack() and pc.EffectiveDominate() into the threshold check.
func LikelyToHit(pc *card.CardState) bool {
	return LikelyDamageHits(pc.EffectiveAttack(), pc.EffectiveDominate())
}

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit. True iff
// LikelyDamageDealt would credit any damage past blocks.
func LikelyDamageHits(n int, dominate bool) bool {
	return LikelyDamageDealt(n, dominate) > 0
}

// LikelyDamageDealt is the "how much" sibling of LikelyDamageHits: it returns the damage
// expected to land past the opponent's blocks. Heuristic assumes opponents block in
// multiples of 3 (the typical block-card power), so 1/4/7 leaks 1 (opponent over-pays
// blocking these awkward amounts and eats them instead). Multiples of 3 are easy-block
// amounts and deal 0.
//
// Dominate flips the math: the defender is capped at one blocking card (3 power), so any
// attack of 5+ power slips n-3 damage past that single block.
func LikelyDamageDealt(n int, dominate bool) int {
	if dominate && n >= 5 {
		return n - 3
	}
	if n == 1 || n == 4 || n == 7 {
		return 1
	}
	return 0
}

// OptDebug, when true, makes GameEngine.Opt print every Opt outcome to stdout (input
// cards, top split, bottom split). Process-wide toggle set by cmd/fabsim's -debug flag at
// the top of a run. Off in production. Not synchronised; today the sim is single-goroutine
// per evaluator, so a plain bool is fine.
var OptDebug bool
