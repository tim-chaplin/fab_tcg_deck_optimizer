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

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

func (DeathlyDuetRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	deathlyDuetApplyRiders(s, l, self)
}

func (DeathlyDuetYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	deathlyDuetApplyRiders(s, l, self)
}

func (DeathlyDuetBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	deathlyDuetApplyRiders(s, l, self)
}

// deathlyDuetApplyRiders folds Deathly Duet's two pitch-conditional riders into self and
// state, then emits the chain step:
//   - Attack-action attributed → +2{p} power buff lands on self.BonusAttack so EffectiveAttack
//     and LikelyToHit see the buffed power, and the chain step's (+N) reflects it directly.
//   - Non-attack-action attributed → 2 Runechants enter during Deathly Duet's own attack
//     resolution; the rider lands as a "Created 2 runechants" sub-line under self.
//
// Both riders can stack when self.PitchedToPlay contains both roles.
func deathlyDuetApplyRiders(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
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
	n := self.DealEffectiveAttack(s)
	l.Log(self, n)
	if nonAttackActionPitched {
		s.CreateRunechants(2)
		l.LogRider(self, 2, "Created 2 runechants")
	}
}
