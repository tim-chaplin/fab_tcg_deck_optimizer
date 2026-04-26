// Deathly Duet — Runeblade Action - Attack. Cost 2, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "When Deathly Duet attacks, if an attack action card was pitched to play it, it gains
// +2{p}. If a 'non-attack' action card was pitched to play it, create 2 Runechant tokens."

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var deathlyDuetTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type DeathlyDuetRed struct{}

func (DeathlyDuetRed) ID() card.ID                 { return card.DeathlyDuetRed }
func (DeathlyDuetRed) Name() string                 { return "Deathly Duet" }
func (DeathlyDuetRed) Cost(*card.TurnState) int                    { return 2 }
func (DeathlyDuetRed) Pitch() int                   { return 1 }
func (DeathlyDuetRed) Attack() int                  { return 4 }
func (DeathlyDuetRed) Defense() int                 { return 3 }
func (DeathlyDuetRed) Types() card.TypeSet       { return deathlyDuetTypes }
func (DeathlyDuetRed) GoAgain() bool                { return false }
// not implemented: Pitched scan can fire both riders independently of which pitched card paid
// for which play (over-credits when both an attack and a non-attack action are pitched)
func (DeathlyDuetRed) NotImplemented()             {}
func (c DeathlyDuetRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, deathlyDuetPlay(c.Attack(), s)-self.Card.Attack())
}
type DeathlyDuetYellow struct{}

func (DeathlyDuetYellow) ID() card.ID                 { return card.DeathlyDuetYellow }
func (DeathlyDuetYellow) Name() string                 { return "Deathly Duet" }
func (DeathlyDuetYellow) Cost(*card.TurnState) int                    { return 2 }
func (DeathlyDuetYellow) Pitch() int                   { return 2 }
func (DeathlyDuetYellow) Attack() int                  { return 3 }
func (DeathlyDuetYellow) Defense() int                 { return 3 }
func (DeathlyDuetYellow) Types() card.TypeSet       { return deathlyDuetTypes }
func (DeathlyDuetYellow) GoAgain() bool                { return false }
// not implemented: Pitched scan can fire both riders independently of which pitched card paid
// for which play (over-credits when both an attack and a non-attack action are pitched)
func (DeathlyDuetYellow) NotImplemented()             {}
func (c DeathlyDuetYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, deathlyDuetPlay(c.Attack(), s)-self.Card.Attack())
}
type DeathlyDuetBlue struct{}

func (DeathlyDuetBlue) ID() card.ID                 { return card.DeathlyDuetBlue }
func (DeathlyDuetBlue) Name() string                 { return "Deathly Duet" }
func (DeathlyDuetBlue) Cost(*card.TurnState) int                    { return 2 }
func (DeathlyDuetBlue) Pitch() int                   { return 3 }
func (DeathlyDuetBlue) Attack() int                  { return 2 }
func (DeathlyDuetBlue) Defense() int                 { return 3 }
func (DeathlyDuetBlue) Types() card.TypeSet       { return deathlyDuetTypes }
func (DeathlyDuetBlue) GoAgain() bool                { return false }
// not implemented: Pitched scan can fire both riders independently of which pitched card paid
// for which play (over-credits when both an attack and a non-attack action are pitched)
func (DeathlyDuetBlue) NotImplemented()             {}
func (c DeathlyDuetBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttackPlus(self, deathlyDuetPlay(c.Attack(), s)-self.Card.Attack())
}
func deathlyDuetPlay(base int, s *card.TurnState) int {
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

	dmg := base
	if attackPitched {
		dmg += 2
	}
	if nonAttackActionPitched {
		// Two Runechants enter during Deathly Duet's own attack resolution. No guard on a
		// following attack existing — Deathly Duet itself is the attack.
		dmg += s.CreateRunechants(2)
	}
	return dmg
}
