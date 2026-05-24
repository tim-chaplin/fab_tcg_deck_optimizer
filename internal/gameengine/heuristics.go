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

// LikelyDamageHits is the raw-integer threshold check behind LikelyToHit. A typical FaB
// card is worth ~3 points, so blocking 1/4/7 with a pitch or block card over-pays — the
// opponent eats the damage instead. Multiples of 3 are the easy-to-block amounts.
//
// Dominate flips the math for cards printed (or granted) with the Dominate keyword: the
// defender is capped at one blocking card, so any attack of 5+ power slips at least 2
// damage past that single block.
func LikelyDamageHits(n int, dominate bool) bool {
	if dominate && n >= 5 {
		return true
	}
	return n == 1 || n == 4 || n == 7
}

// OptDebug, when true, makes GameEngine.Opt print every Opt outcome to stdout (input
// cards, top split, bottom split). Process-wide toggle set by cmd/fabsim's -debug flag at
// the top of a run. Off in production. Not synchronised; today the sim is single-goroutine
// per evaluator, so a plain bool is fine.
var OptDebug bool
