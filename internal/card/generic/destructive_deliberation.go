// Destructive Deliberation — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue
// 3. Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits a hero, create a Ponder token."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var destructiveDeliberationTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DestructiveDeliberationRed struct{}

func (DestructiveDeliberationRed) ID() card.ID                 { return card.DestructiveDeliberationRed }
func (DestructiveDeliberationRed) Name() string                { return "Destructive Deliberation" }
func (DestructiveDeliberationRed) Cost(*card.TurnState) int                   { return 2 }
func (DestructiveDeliberationRed) Pitch() int                  { return 1 }
func (DestructiveDeliberationRed) Attack() int                 { return 5 }
func (DestructiveDeliberationRed) Defense() int                { return 2 }
func (DestructiveDeliberationRed) Types() card.TypeSet         { return destructiveDeliberationTypes }
func (DestructiveDeliberationRed) GoAgain() bool               { return false }
// not implemented: ponder tokens
func (DestructiveDeliberationRed) NotImplemented()             {}
func (c DestructiveDeliberationRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, destructiveDeliberationDamage(c.Attack(), self)-self.Card.Attack())
}
type DestructiveDeliberationYellow struct{}

func (DestructiveDeliberationYellow) ID() card.ID                 { return card.DestructiveDeliberationYellow }
func (DestructiveDeliberationYellow) Name() string                { return "Destructive Deliberation" }
func (DestructiveDeliberationYellow) Cost(*card.TurnState) int                   { return 2 }
func (DestructiveDeliberationYellow) Pitch() int                  { return 2 }
func (DestructiveDeliberationYellow) Attack() int                 { return 4 }
func (DestructiveDeliberationYellow) Defense() int                { return 2 }
func (DestructiveDeliberationYellow) Types() card.TypeSet         { return destructiveDeliberationTypes }
func (DestructiveDeliberationYellow) GoAgain() bool               { return false }
// not implemented: ponder tokens
func (DestructiveDeliberationYellow) NotImplemented()             {}
func (c DestructiveDeliberationYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, destructiveDeliberationDamage(c.Attack(), self)-self.Card.Attack())
}
type DestructiveDeliberationBlue struct{}

func (DestructiveDeliberationBlue) ID() card.ID                 { return card.DestructiveDeliberationBlue }
func (DestructiveDeliberationBlue) Name() string                { return "Destructive Deliberation" }
func (DestructiveDeliberationBlue) Cost(*card.TurnState) int                   { return 2 }
func (DestructiveDeliberationBlue) Pitch() int                  { return 3 }
func (DestructiveDeliberationBlue) Attack() int                 { return 3 }
func (DestructiveDeliberationBlue) Defense() int                { return 2 }
func (DestructiveDeliberationBlue) Types() card.TypeSet         { return destructiveDeliberationTypes }
func (DestructiveDeliberationBlue) GoAgain() bool               { return false }
// not implemented: ponder tokens
func (DestructiveDeliberationBlue) NotImplemented()             {}
func (c DestructiveDeliberationBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, destructiveDeliberationDamage(c.Attack(), self)-self.Card.Attack())
}
// destructiveDeliberationDamage is a breadcrumb for the on-hit "create a Ponder token" rider —
// Ponder tokens aren't tracked (see TODO.md).
func destructiveDeliberationDamage(attack int, self *card.CardState) int {
	if card.LikelyToHit(self) {
		// TODO: model on-hit Ponder token creation rider.
	}
	return attack
}
