// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// Aura condition checked via s.HasAuraInPlay(). The "when this hits a hero" discard rider fires
// only when the aura clause is satisfied AND this card's printed attack is likely to land
// (1/4/7 per card.LikelyToHit). Runechants alongside don't count — "this hits" is strictly
// about this card's damage reaching the hero. On-hit discard credits +3.

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var reekOfCorruptionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// reekOfCorruptionDamage returns the base attack plus the discard rider when the aura condition
// is satisfied AND this card's printed attack is likely to land on its own.
func reekOfCorruptionDamage(attack int, s *card.TurnState, self *card.CardState) int {
	if s != nil && s.HasAuraInPlay() && card.LikelyToHit(self) {
		return attack + card.DiscardValue
	}
	return attack
}

type ReekOfCorruptionRed struct{}

func (ReekOfCorruptionRed) ID() card.ID                  { return card.ReekOfCorruptionRed }
func (ReekOfCorruptionRed) Name() string                 { return "Reek of Corruption" }
func (ReekOfCorruptionRed) Cost(*card.TurnState) int                    { return 2 }
func (ReekOfCorruptionRed) Pitch() int                   { return 1 }
func (ReekOfCorruptionRed) Attack() int                  { return 4 }
func (ReekOfCorruptionRed) Defense() int                 { return 3 }
func (ReekOfCorruptionRed) Types() card.TypeSet          { return reekOfCorruptionTypes }
func (ReekOfCorruptionRed) GoAgain() bool                { return false }
func (c ReekOfCorruptionRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, reekOfCorruptionDamage(c.Attack(), s, self)-self.Card.Attack())
}
type ReekOfCorruptionYellow struct{}

func (ReekOfCorruptionYellow) ID() card.ID                  { return card.ReekOfCorruptionYellow }
func (ReekOfCorruptionYellow) Name() string                 { return "Reek of Corruption" }
func (ReekOfCorruptionYellow) Cost(*card.TurnState) int                    { return 2 }
func (ReekOfCorruptionYellow) Pitch() int                   { return 2 }
func (ReekOfCorruptionYellow) Attack() int                  { return 3 }
func (ReekOfCorruptionYellow) Defense() int                 { return 3 }
func (ReekOfCorruptionYellow) Types() card.TypeSet          { return reekOfCorruptionTypes }
func (ReekOfCorruptionYellow) GoAgain() bool                { return false }
func (c ReekOfCorruptionYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, reekOfCorruptionDamage(c.Attack(), s, self)-self.Card.Attack())
}
type ReekOfCorruptionBlue struct{}

func (ReekOfCorruptionBlue) ID() card.ID                  { return card.ReekOfCorruptionBlue }
func (ReekOfCorruptionBlue) Name() string                 { return "Reek of Corruption" }
func (ReekOfCorruptionBlue) Cost(*card.TurnState) int                    { return 2 }
func (ReekOfCorruptionBlue) Pitch() int                   { return 3 }
func (ReekOfCorruptionBlue) Attack() int                  { return 2 }
func (ReekOfCorruptionBlue) Defense() int                 { return 3 }
func (ReekOfCorruptionBlue) Types() card.TypeSet          { return reekOfCorruptionTypes }
func (ReekOfCorruptionBlue) GoAgain() bool                { return false }
func (c ReekOfCorruptionBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, reekOfCorruptionDamage(c.Attack(), s, self)-self.Card.Attack())
}