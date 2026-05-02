// Nebula Blade — Runeblade Weapon - Sword (2H). Cost 2, Power 1.
// Text: "Once per Turn Action - {r}{r}: Attack. If Nebula Blade hits, create a Runechant token. If
// you have played a 'non-attack' action card this turn, Nebula Blade gains +3{p} until end of
// turn."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var nebulaBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type NebulaBlade struct{}

func (NebulaBlade) ID() ids.WeaponID        { return ids.NebulaBladeID }
func (NebulaBlade) Name() string            { return "Nebula Blade" }
func (NebulaBlade) Cost(*sim.TurnState) int { return 2 }
func (NebulaBlade) Pitch() int              { return 0 }
func (NebulaBlade) Attack() int             { return 1 }
func (NebulaBlade) Defense() int            { return 0 }
func (NebulaBlade) Types() card.TypeSet     { return nebulaBladeTypes }
func (NebulaBlade) GoAgain() bool           { return false }
func (NebulaBlade) Hands() int              { return 2 }
func (c NebulaBlade) Play(s *sim.TurnState, self *sim.CardState) {
	if s.NonAttackActionPlayed {
		self.BonusAttack += 3
	}
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
	self.OnHit = append(self.OnHit, func(state *sim.TurnState) {
		created := state.CreateRunechant()
		state.AddValue(created)
		state.LogRider(self, created, "On-hit created a runechant")
	})
}
