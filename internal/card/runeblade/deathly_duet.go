// Deathly Duet — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Deathly Duet attacks, if an attack action card was pitched to play it, it gains
// +2{p}. If a 'non-attack' action card was pitched to play it, create 2 Runechant tokens."

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var deathlyDuetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type DeathlyDuetRed struct{}

func (DeathlyDuetRed) ID() card.ID              { return card.DeathlyDuetRed }
func (DeathlyDuetRed) Name() string             { return "Deathly Duet" }
func (DeathlyDuetRed) Cost(*card.TurnState) int { return 2 }
func (DeathlyDuetRed) Pitch() int               { return 1 }
func (DeathlyDuetRed) Attack() int              { return 4 }
func (DeathlyDuetRed) Defense() int             { return 3 }
func (DeathlyDuetRed) Types() card.TypeSet      { return deathlyDuetTypes }
func (DeathlyDuetRed) GoAgain() bool            { return false }

// not implemented: Pitched scan can fire both riders independently of which pitched card paid
// for which play (over-credits when both an attack and a non-attack action are pitched)
func (DeathlyDuetRed) NotImplemented() {}
func (DeathlyDuetRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, deathlyDuetApplyRiders(s, self))
}

type DeathlyDuetYellow struct{}

func (DeathlyDuetYellow) ID() card.ID              { return card.DeathlyDuetYellow }
func (DeathlyDuetYellow) Name() string             { return "Deathly Duet" }
func (DeathlyDuetYellow) Cost(*card.TurnState) int { return 2 }
func (DeathlyDuetYellow) Pitch() int               { return 2 }
func (DeathlyDuetYellow) Attack() int              { return 3 }
func (DeathlyDuetYellow) Defense() int             { return 3 }
func (DeathlyDuetYellow) Types() card.TypeSet      { return deathlyDuetTypes }
func (DeathlyDuetYellow) GoAgain() bool            { return false }

// not implemented: Pitched scan can fire both riders independently of which pitched card paid
// for which play (over-credits when both an attack and a non-attack action are pitched)
func (DeathlyDuetYellow) NotImplemented() {}
func (DeathlyDuetYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, deathlyDuetApplyRiders(s, self))
}

type DeathlyDuetBlue struct{}

func (DeathlyDuetBlue) ID() card.ID              { return card.DeathlyDuetBlue }
func (DeathlyDuetBlue) Name() string             { return "Deathly Duet" }
func (DeathlyDuetBlue) Cost(*card.TurnState) int { return 2 }
func (DeathlyDuetBlue) Pitch() int               { return 3 }
func (DeathlyDuetBlue) Attack() int              { return 2 }
func (DeathlyDuetBlue) Defense() int             { return 3 }
func (DeathlyDuetBlue) Types() card.TypeSet      { return deathlyDuetTypes }
func (DeathlyDuetBlue) GoAgain() bool            { return false }

// not implemented: Pitched scan can fire both riders independently of which pitched card paid
// for which play (over-credits when both an attack and a non-attack action are pitched)
func (DeathlyDuetBlue) NotImplemented() {}
func (DeathlyDuetBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, deathlyDuetApplyRiders(s, self))
}

// deathlyDuetApplyRiders folds Deathly Duet's two pitch-conditional riders into self and
// state and returns the runechant-damage rider for the chain step's (+N) display:
//   - Attack-action pitched → +2{p} power buff lands on self.BonusAttack so EffectiveAttack
//     and LikelyToHit see the buffed power.
//   - Non-attack-action pitched → 2 Runechants enter during Deathly Duet's own attack
//     resolution; their +2 damage credit is the returned rider.
//
// Both riders can stack when both pitched roles are present.
func deathlyDuetApplyRiders(s *card.TurnState, self *card.CardState) int {
	var attackPitched, nonAttackActionPitched bool
	for _, p := range s.Pitched {
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
	rider := 0
	if nonAttackActionPitched {
		rider += s.CreateRunechants(2)
	}
	return rider
}
