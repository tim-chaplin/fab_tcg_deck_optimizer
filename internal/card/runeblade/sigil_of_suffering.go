// Sigil of Suffering — Runeblade Defense Reaction. Cost 0, Arcane 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed defense: Red 3, Yellow 2, Blue 1.
// Text: "Deal 1 arcane damage to the attacking hero. If you have dealt arcane damage this turn,
// Sigil of Suffering gains +1{d}."

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sigilOfSufferingTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

// sigilOfSufferingPlay deals 1 arcane damage to the attacking hero and marks ArcaneDamageDealt
// so same-turn triggers reading "if you've dealt arcane damage this turn" see it.
func sigilOfSufferingPlay(s *card.TurnState) int {
	return s.DealArcaneDamage(1)
}

type SigilOfSufferingRed struct{}

func (SigilOfSufferingRed) ID() card.ID                   { return card.SigilOfSufferingRed }
func (SigilOfSufferingRed) Name() string                  { return "Sigil of Suffering" }
func (SigilOfSufferingRed) Cost(*card.TurnState) int                     { return 0 }
func (SigilOfSufferingRed) Pitch() int                    { return 1 }
func (SigilOfSufferingRed) Attack() int                   { return 0 }
func (SigilOfSufferingRed) Defense() int                  { return 4 }
func (SigilOfSufferingRed) Types() card.TypeSet           { return sigilOfSufferingTypes }
func (SigilOfSufferingRed) GoAgain() bool                 { return false }
// not implemented: conditional +1{d} 'if dealt arcane damage' is baked unconditionally into
// Defense — over-credits when the trigger condition isn't met
func (SigilOfSufferingRed) NotImplemented()             {}
func (SigilOfSufferingRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, sigilOfSufferingPlay(s))
}
type SigilOfSufferingYellow struct{}

func (SigilOfSufferingYellow) ID() card.ID                   { return card.SigilOfSufferingYellow }
func (SigilOfSufferingYellow) Name() string                  { return "Sigil of Suffering" }
func (SigilOfSufferingYellow) Cost(*card.TurnState) int                     { return 0 }
func (SigilOfSufferingYellow) Pitch() int                    { return 2 }
func (SigilOfSufferingYellow) Attack() int                   { return 0 }
func (SigilOfSufferingYellow) Defense() int                  { return 3 }
func (SigilOfSufferingYellow) Types() card.TypeSet           { return sigilOfSufferingTypes }
func (SigilOfSufferingYellow) GoAgain() bool                 { return false }
// not implemented: conditional +1{d} 'if dealt arcane damage' is baked unconditionally into
// Defense — over-credits when the trigger condition isn't met
func (SigilOfSufferingYellow) NotImplemented()             {}
func (SigilOfSufferingYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, sigilOfSufferingPlay(s))
}
type SigilOfSufferingBlue struct{}

func (SigilOfSufferingBlue) ID() card.ID                   { return card.SigilOfSufferingBlue }
func (SigilOfSufferingBlue) Name() string                  { return "Sigil of Suffering" }
func (SigilOfSufferingBlue) Cost(*card.TurnState) int                     { return 0 }
func (SigilOfSufferingBlue) Pitch() int                    { return 3 }
func (SigilOfSufferingBlue) Attack() int                   { return 0 }
func (SigilOfSufferingBlue) Defense() int                  { return 2 }
func (SigilOfSufferingBlue) Types() card.TypeSet           { return sigilOfSufferingTypes }
func (SigilOfSufferingBlue) GoAgain() bool                 { return false }
// not implemented: conditional +1{d} 'if dealt arcane damage' is baked unconditionally into
// Defense — over-credits when the trigger condition isn't met
func (SigilOfSufferingBlue) NotImplemented()             {}
func (SigilOfSufferingBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, sigilOfSufferingPlay(s))
}