// Reek of Corruption — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "If you have played or created an aura this turn, Reek of Corruption gains 'When this
// hits a hero, they discard a card.'"
//
// Aura condition checked via s.HasPlayedOrCreatedAura(). The on-hit discard rider's hit gate
// is owned by ApplyAndLogRiderOnHit. Runechants alongside don't count — "this hits" is strictly
// about this card's damage reaching the hero. On-hit discard credits +3.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var reekOfCorruptionTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

// reekOfCorruptionApplyRider emits the on-hit discard rider as a sub-line under self's
// chain step when the aura precondition is satisfied.
func reekOfCorruptionApplyRider(s *sim.TurnState, self *sim.CardState) {
	if !s.HasPlayedOrCreatedAura() {
		return
	}
	if sim.LikelyToHit(self) {
		s.AddValue(sim.DiscardValue)
		s.LogRider(self, sim.DiscardValue, "On-hit discarded a card")
	}
}

type ReekOfCorruptionRed struct{}

func (ReekOfCorruptionRed) ID() ids.CardID          { return ids.ReekOfCorruptionRed }
func (ReekOfCorruptionRed) Name() string            { return "Reek of Corruption" }
func (ReekOfCorruptionRed) Cost(*sim.TurnState) int { return 2 }
func (ReekOfCorruptionRed) Pitch() int              { return 1 }
func (ReekOfCorruptionRed) Attack() int             { return 4 }
func (ReekOfCorruptionRed) Defense() int            { return 3 }
func (ReekOfCorruptionRed) Types() card.TypeSet     { return reekOfCorruptionTypes }
func (ReekOfCorruptionRed) GoAgain() bool           { return false }
func (ReekOfCorruptionRed) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	reekOfCorruptionApplyRider(s, self)
}

type ReekOfCorruptionYellow struct{}

func (ReekOfCorruptionYellow) ID() ids.CardID          { return ids.ReekOfCorruptionYellow }
func (ReekOfCorruptionYellow) Name() string            { return "Reek of Corruption" }
func (ReekOfCorruptionYellow) Cost(*sim.TurnState) int { return 2 }
func (ReekOfCorruptionYellow) Pitch() int              { return 2 }
func (ReekOfCorruptionYellow) Attack() int             { return 3 }
func (ReekOfCorruptionYellow) Defense() int            { return 3 }
func (ReekOfCorruptionYellow) Types() card.TypeSet     { return reekOfCorruptionTypes }
func (ReekOfCorruptionYellow) GoAgain() bool           { return false }
func (ReekOfCorruptionYellow) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	reekOfCorruptionApplyRider(s, self)
}

type ReekOfCorruptionBlue struct{}

func (ReekOfCorruptionBlue) ID() ids.CardID          { return ids.ReekOfCorruptionBlue }
func (ReekOfCorruptionBlue) Name() string            { return "Reek of Corruption" }
func (ReekOfCorruptionBlue) Cost(*sim.TurnState) int { return 2 }
func (ReekOfCorruptionBlue) Pitch() int              { return 3 }
func (ReekOfCorruptionBlue) Attack() int             { return 2 }
func (ReekOfCorruptionBlue) Defense() int            { return 3 }
func (ReekOfCorruptionBlue) Types() card.TypeSet     { return reekOfCorruptionTypes }
func (ReekOfCorruptionBlue) GoAgain() bool           { return false }
func (ReekOfCorruptionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	reekOfCorruptionApplyRider(s, self)
}
