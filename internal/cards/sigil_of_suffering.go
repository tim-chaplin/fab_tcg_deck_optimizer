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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfSufferingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

func sigilOfSufferingPlay(s *sim.TurnState, self *sim.CardState) {
	if s.ArcaneDamageDealt || sim.LikelyDamageHits(1, false) {
		self.BonusDefense++
	}
	n := self.DealEffectiveDefense(s)
	s.Log(self, n)
	s.AddValue(s.DealArcaneDamage(1))
	s.LogRider(self, 1, "Dealt 1 arcane damage")
}

type SigilOfSufferingRed struct{}

func (SigilOfSufferingRed) ID() ids.CardID          { return ids.SigilOfSufferingRed }
func (SigilOfSufferingRed) Name() string            { return "Sigil of Suffering" }
func (SigilOfSufferingRed) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSufferingRed) Pitch() int              { return 1 }
func (SigilOfSufferingRed) Attack() int             { return 0 }
func (SigilOfSufferingRed) Defense() int            { return 3 }
func (SigilOfSufferingRed) Types() card.TypeSet     { return sigilOfSufferingTypes }
func (SigilOfSufferingRed) GoAgain() bool           { return false }
func (SigilOfSufferingRed) Play(s *sim.TurnState, self *sim.CardState) {
	sigilOfSufferingPlay(s, self)
}

type SigilOfSufferingYellow struct{}

func (SigilOfSufferingYellow) ID() ids.CardID          { return ids.SigilOfSufferingYellow }
func (SigilOfSufferingYellow) Name() string            { return "Sigil of Suffering" }
func (SigilOfSufferingYellow) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSufferingYellow) Pitch() int              { return 2 }
func (SigilOfSufferingYellow) Attack() int             { return 0 }
func (SigilOfSufferingYellow) Defense() int            { return 2 }
func (SigilOfSufferingYellow) Types() card.TypeSet     { return sigilOfSufferingTypes }
func (SigilOfSufferingYellow) GoAgain() bool           { return false }
func (SigilOfSufferingYellow) Play(s *sim.TurnState, self *sim.CardState) {
	sigilOfSufferingPlay(s, self)
}

type SigilOfSufferingBlue struct{}

func (SigilOfSufferingBlue) ID() ids.CardID          { return ids.SigilOfSufferingBlue }
func (SigilOfSufferingBlue) Name() string            { return "Sigil of Suffering" }
func (SigilOfSufferingBlue) Cost(*sim.TurnState) int { return 0 }
func (SigilOfSufferingBlue) Pitch() int              { return 3 }
func (SigilOfSufferingBlue) Attack() int             { return 0 }
func (SigilOfSufferingBlue) Defense() int            { return 1 }
func (SigilOfSufferingBlue) Types() card.TypeSet     { return sigilOfSufferingTypes }
func (SigilOfSufferingBlue) GoAgain() bool           { return false }
func (SigilOfSufferingBlue) Play(s *sim.TurnState, self *sim.CardState) {
	sigilOfSufferingPlay(s, self)
}
