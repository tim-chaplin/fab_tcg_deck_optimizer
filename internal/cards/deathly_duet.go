// Deathly Duet — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Deathly Duet attacks, if an attack action card was pitched to play it, it gains
// +2{p}. If a 'non-attack' action card was pitched to play it, create 2 Runechant tokens."
//
// Both riders read self.PitchedToPlay (the cards the chain runner attributed to funding
// THIS copy's cost) — they fire independently when those exact cards include an attack
// action / non-attack action respectively.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var deathlyDuetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type DeathlyDuetRed struct{}

func (DeathlyDuetRed) ID() ids.CardID          { return ids.DeathlyDuetRed }
func (DeathlyDuetRed) Name() string            { return "Deathly Duet" }
func (DeathlyDuetRed) Cost(*sim.TurnState) int { return 2 }
func (DeathlyDuetRed) Pitch() int              { return 1 }
func (DeathlyDuetRed) Attack() int             { return 4 }
func (DeathlyDuetRed) Defense() int            { return 3 }
func (DeathlyDuetRed) Types() card.TypeSet     { return deathlyDuetTypes }
func (DeathlyDuetRed) GoAgain() bool           { return false }
func (DeathlyDuetRed) Play(s *sim.TurnState, self *sim.CardState) {
	deathlyDuetApplyRiders(s, self)
}

type DeathlyDuetYellow struct{}

func (DeathlyDuetYellow) ID() ids.CardID          { return ids.DeathlyDuetYellow }
func (DeathlyDuetYellow) Name() string            { return "Deathly Duet" }
func (DeathlyDuetYellow) Cost(*sim.TurnState) int { return 2 }
func (DeathlyDuetYellow) Pitch() int              { return 2 }
func (DeathlyDuetYellow) Attack() int             { return 3 }
func (DeathlyDuetYellow) Defense() int            { return 3 }
func (DeathlyDuetYellow) Types() card.TypeSet     { return deathlyDuetTypes }
func (DeathlyDuetYellow) GoAgain() bool           { return false }
func (DeathlyDuetYellow) Play(s *sim.TurnState, self *sim.CardState) {
	deathlyDuetApplyRiders(s, self)
}

type DeathlyDuetBlue struct{}

func (DeathlyDuetBlue) ID() ids.CardID          { return ids.DeathlyDuetBlue }
func (DeathlyDuetBlue) Name() string            { return "Deathly Duet" }
func (DeathlyDuetBlue) Cost(*sim.TurnState) int { return 2 }
func (DeathlyDuetBlue) Pitch() int              { return 3 }
func (DeathlyDuetBlue) Attack() int             { return 2 }
func (DeathlyDuetBlue) Defense() int            { return 3 }
func (DeathlyDuetBlue) Types() card.TypeSet     { return deathlyDuetTypes }
func (DeathlyDuetBlue) GoAgain() bool           { return false }
func (DeathlyDuetBlue) Play(s *sim.TurnState, self *sim.CardState) {
	deathlyDuetApplyRiders(s, self)
}

// deathlyDuetApplyRiders folds Deathly Duet's two pitch-conditional riders into self and
// state, then emits the chain step:
//   - Attack-action attributed → +2{p} power buff lands on self.BonusAttack so EffectiveAttack
//     and LikelyToHit see the buffed power, and the chain step's (+N) reflects it directly.
//   - Non-attack-action attributed → 2 Runechants enter during Deathly Duet's own attack
//     resolution; the rider lands as a "Created 2 runechants" sub-line under self.
//
// Both riders can stack when self.PitchedToPlay contains both roles.
func deathlyDuetApplyRiders(s *sim.TurnState, self *sim.CardState) {
	var attackPitched, nonAttackActionPitched bool
	for _, p := range self.PitchedToPlay {
		t := p.Types()
		if t.Has(card.TypeAttack) {
			attackPitched = true
		}
		if t.IsNonAttackAction() {
			nonAttackActionPitched = true
		}
	}
	if attackPitched {
		self.BonusAttack += 2
	}
	s.LogChain(self, s.AddValue(self.EffectiveAttack()))
	if nonAttackActionPitched {
		s.AddValue(s.CreateRunechants(2))
		s.LogRider(self, 2, "Created 2 runechants")
	}
}
