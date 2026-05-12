package gameengine

import "github.com/tim-chaplin/fab-deck-optimizer/v2/card"

// Damage-equivalent tuning constants. The engine owns these — cards reach the values
// through GameEngine accessors so the model's calibration lives in one place.

// DiscardValue is the damage-equivalent credited when the opponent is forced to discard one
// card. A typical FaB card is worth ~3 points of tempo.
const DiscardValue = 3

// GoldTokenValue is the placeholder credit for NotImplemented Gold-creating cards. Zero
// because Gold pays out on activation (see Gold token ability), not at creation.
const GoldTokenValue = 0

// LikelyToHit reports whether self's attack is likely to land past the opponent's blocks.
// Folds self.EffectiveAttack() and self.EffectiveDominate() into the threshold check.
func LikelyToHit(self *card.CardState) bool {
	return LikelyDamageHits(self.EffectiveAttack(), self.EffectiveDominate())
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
