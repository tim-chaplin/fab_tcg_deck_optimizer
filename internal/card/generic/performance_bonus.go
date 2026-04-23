// Performance Bonus — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets **Go
// again**."
//
// The on-hit Gold token is modelled as +1 damage-equivalent (one resource worth), gated on
// card.LikelyToHit. The arsenal-conditional Go again fires via self.GrantedGoAgain when
// self.FromArsenal reports this copy came from the arsenal slot.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var performanceBonusTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// performanceBonusDamage returns the base attack plus the Gold-token rider when the attack is
// likely to land.
func performanceBonusDamage(attack int) int {
	if card.LikelyToHit(attack, false) {
		return attack + card.GoldTokenValue
	}
	return attack
}

// performanceBonusPlay credits the on-hit Gold token and grants self Go again when this
// copy was played from arsenal.
func performanceBonusPlay(c card.Card, self *card.CardState) int {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	return performanceBonusDamage(c.Attack())
}

type PerformanceBonusRed struct{}

func (PerformanceBonusRed) ID() card.ID                  { return card.PerformanceBonusRed }
func (PerformanceBonusRed) Name() string                 { return "Performance Bonus (Red)" }
func (PerformanceBonusRed) Cost(*card.TurnState) int                    { return 0 }
func (PerformanceBonusRed) Pitch() int                   { return 1 }
func (PerformanceBonusRed) Attack() int                  { return 3 }
func (PerformanceBonusRed) Defense() int                 { return 2 }
func (PerformanceBonusRed) Types() card.TypeSet          { return performanceBonusTypes }
func (PerformanceBonusRed) GoAgain() bool                { return false }
func (c PerformanceBonusRed) Play(_ *card.TurnState, self *card.CardState) int { return performanceBonusPlay(c, self) }

type PerformanceBonusYellow struct{}

func (PerformanceBonusYellow) ID() card.ID                  { return card.PerformanceBonusYellow }
func (PerformanceBonusYellow) Name() string                 { return "Performance Bonus (Yellow)" }
func (PerformanceBonusYellow) Cost(*card.TurnState) int                    { return 0 }
func (PerformanceBonusYellow) Pitch() int                   { return 2 }
func (PerformanceBonusYellow) Attack() int                  { return 2 }
func (PerformanceBonusYellow) Defense() int                 { return 2 }
func (PerformanceBonusYellow) Types() card.TypeSet          { return performanceBonusTypes }
func (PerformanceBonusYellow) GoAgain() bool                { return false }
func (c PerformanceBonusYellow) Play(_ *card.TurnState, self *card.CardState) int { return performanceBonusPlay(c, self) }

type PerformanceBonusBlue struct{}

func (PerformanceBonusBlue) ID() card.ID                  { return card.PerformanceBonusBlue }
func (PerformanceBonusBlue) Name() string                 { return "Performance Bonus (Blue)" }
func (PerformanceBonusBlue) Cost(*card.TurnState) int                    { return 0 }
func (PerformanceBonusBlue) Pitch() int                   { return 3 }
func (PerformanceBonusBlue) Attack() int                  { return 1 }
func (PerformanceBonusBlue) Defense() int                 { return 2 }
func (PerformanceBonusBlue) Types() card.TypeSet          { return performanceBonusTypes }
func (PerformanceBonusBlue) GoAgain() bool                { return false }
func (c PerformanceBonusBlue) Play(_ *card.TurnState, self *card.CardState) int { return performanceBonusPlay(c, self) }
