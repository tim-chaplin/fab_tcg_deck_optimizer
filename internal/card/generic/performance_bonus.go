// Performance Bonus — Generic Action - Attack. Cost 0. Printed power: Red 3, Yellow 2, Blue 1.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this hits, create a Gold token. If this was played from arsenal, it gets **Go
// again**."
//
// Standard played-from-arsenal go-again (docs/dev-standards.md). Gold-token creation isn't
// modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var performanceBonusTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// performanceBonusPlay grants self Go again when this copy was played from arsenal, emits
// the chain step, then writes the on-hit Gold-token rider as a sub-line when LikelyToHit
// fires.
func performanceBonusPlay(s *card.TurnState, self *card.CardState) {
	if self.FromArsenal {
		self.GrantedGoAgain = true
	}
	s.ApplyAndLogEffectiveAttack(self)
	if card.LikelyToHit(self) {
		s.LogRiderOnPlay(self, "On-hit created a gold token", card.GoldTokenValue)
	}
}

type PerformanceBonusRed struct{}

func (PerformanceBonusRed) ID() card.ID              { return card.PerformanceBonusRed }
func (PerformanceBonusRed) Name() string             { return "Performance Bonus" }
func (PerformanceBonusRed) Cost(*card.TurnState) int { return 0 }
func (PerformanceBonusRed) Pitch() int               { return 1 }
func (PerformanceBonusRed) Attack() int              { return 3 }
func (PerformanceBonusRed) Defense() int             { return 2 }
func (PerformanceBonusRed) Types() card.TypeSet      { return performanceBonusTypes }
func (PerformanceBonusRed) GoAgain() bool            { return false }

// not implemented: gold tokens
func (PerformanceBonusRed) NotImplemented() {}
func (PerformanceBonusRed) Play(s *card.TurnState, self *card.CardState) {
	performanceBonusPlay(s, self)
}

type PerformanceBonusYellow struct{}

func (PerformanceBonusYellow) ID() card.ID              { return card.PerformanceBonusYellow }
func (PerformanceBonusYellow) Name() string             { return "Performance Bonus" }
func (PerformanceBonusYellow) Cost(*card.TurnState) int { return 0 }
func (PerformanceBonusYellow) Pitch() int               { return 2 }
func (PerformanceBonusYellow) Attack() int              { return 2 }
func (PerformanceBonusYellow) Defense() int             { return 2 }
func (PerformanceBonusYellow) Types() card.TypeSet      { return performanceBonusTypes }
func (PerformanceBonusYellow) GoAgain() bool            { return false }

// not implemented: gold tokens
func (PerformanceBonusYellow) NotImplemented() {}
func (PerformanceBonusYellow) Play(s *card.TurnState, self *card.CardState) {
	performanceBonusPlay(s, self)
}

type PerformanceBonusBlue struct{}

func (PerformanceBonusBlue) ID() card.ID              { return card.PerformanceBonusBlue }
func (PerformanceBonusBlue) Name() string             { return "Performance Bonus" }
func (PerformanceBonusBlue) Cost(*card.TurnState) int { return 0 }
func (PerformanceBonusBlue) Pitch() int               { return 3 }
func (PerformanceBonusBlue) Attack() int              { return 1 }
func (PerformanceBonusBlue) Defense() int             { return 2 }
func (PerformanceBonusBlue) Types() card.TypeSet      { return performanceBonusTypes }
func (PerformanceBonusBlue) GoAgain() bool            { return false }

// not implemented: gold tokens
func (PerformanceBonusBlue) NotImplemented() {}
func (PerformanceBonusBlue) Play(s *card.TurnState, self *card.CardState) {
	performanceBonusPlay(s, self)
}
