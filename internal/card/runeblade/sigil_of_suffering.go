// Sigil of Suffering — Runeblade Defense Reaction. Cost 0, Arcane 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to the attacking hero. If you have dealt arcane damage this turn,
// Sigil of Suffering gains +1{d}."
//
// The chain step credits printed Defense only; the arcane and the +1{d} bonus each land as
// their own post-trigger sub-line under self. The Sigil's own arcane (LikelyDamageHits(1,
// false) is true) satisfies the conditional, so the +1{d} fires whenever there's incoming
// damage left for it to consume — LogDefenseRiderOnPlay caps the rider at the remaining
// IncomingDamage so over-block discards quietly.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSufferingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

func sigilOfSufferingPlay(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
	s.DealAndLogArcaneDamage(self, 1)
	if s.ArcaneDamageDealt {
		s.LogDefenseRiderOnPlay(self, "Gained +1{d} from arcane this turn", 1)
	}
}

type SigilOfSufferingRed struct{}

func (SigilOfSufferingRed) ID() card.ID              { return card.SigilOfSufferingRed }
func (SigilOfSufferingRed) Name() string             { return "Sigil of Suffering" }
func (SigilOfSufferingRed) Cost(*card.TurnState) int { return 0 }
func (SigilOfSufferingRed) Pitch() int               { return 1 }
func (SigilOfSufferingRed) Attack() int              { return 0 }
func (SigilOfSufferingRed) Defense() int             { return 3 }
func (SigilOfSufferingRed) Types() card.TypeSet      { return sigilOfSufferingTypes }
func (SigilOfSufferingRed) GoAgain() bool            { return false }
func (SigilOfSufferingRed) Play(s *card.TurnState, self *card.CardState) {
	sigilOfSufferingPlay(s, self)
}

type SigilOfSufferingYellow struct{}

func (SigilOfSufferingYellow) ID() card.ID              { return card.SigilOfSufferingYellow }
func (SigilOfSufferingYellow) Name() string             { return "Sigil of Suffering" }
func (SigilOfSufferingYellow) Cost(*card.TurnState) int { return 0 }
func (SigilOfSufferingYellow) Pitch() int               { return 2 }
func (SigilOfSufferingYellow) Attack() int              { return 0 }
func (SigilOfSufferingYellow) Defense() int             { return 2 }
func (SigilOfSufferingYellow) Types() card.TypeSet      { return sigilOfSufferingTypes }
func (SigilOfSufferingYellow) GoAgain() bool            { return false }
func (SigilOfSufferingYellow) Play(s *card.TurnState, self *card.CardState) {
	sigilOfSufferingPlay(s, self)
}

type SigilOfSufferingBlue struct{}

func (SigilOfSufferingBlue) ID() card.ID              { return card.SigilOfSufferingBlue }
func (SigilOfSufferingBlue) Name() string             { return "Sigil of Suffering" }
func (SigilOfSufferingBlue) Cost(*card.TurnState) int { return 0 }
func (SigilOfSufferingBlue) Pitch() int               { return 3 }
func (SigilOfSufferingBlue) Attack() int              { return 0 }
func (SigilOfSufferingBlue) Defense() int             { return 1 }
func (SigilOfSufferingBlue) Types() card.TypeSet      { return sigilOfSufferingTypes }
func (SigilOfSufferingBlue) GoAgain() bool            { return false }
func (SigilOfSufferingBlue) Play(s *card.TurnState, self *card.CardState) {
	sigilOfSufferingPlay(s, self)
}
