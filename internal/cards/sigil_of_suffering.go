// Sigil of Suffering — Runeblade Defense Reaction. Cost 0, Arcane 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to the attacking hero. If you have dealt arcane damage this turn,
// Sigil of Suffering gains +1{d}."
//
// Mirrors Hit the High Notes' shape on the defender side: the +1{d} bonus folds into
// BonusDefense before the chain step fires so the (+N) reflects the buffed block, and the
// arcane lands as its own post-trigger sub-line. The Sigil's own printed-1 arcane satisfies
// the conditional via LikelyDamageHits(1, false), so the bonus is credited whenever there's
// IncomingDamage left to absorb it; ApplyAndLogEffectiveDefense's clamp handles over-block.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func sigilOfSufferingPlay(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	if s.ArcaneDamageDealt() || sim.LikelyDamageHits(1, false) {
		self.BonusDefense++
	}
	s.DealArcaneDamage(l, self, 1)
}

func (SigilOfSufferingRed) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sigilOfSufferingPlay(s, l, self)
}

func (SigilOfSufferingYellow) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sigilOfSufferingPlay(s, l, self)
}

func (SigilOfSufferingBlue) Play(s sim.GameEngine, l sim.Logger, self *sim.CardState) {
	sigilOfSufferingPlay(s, l, self)
}
