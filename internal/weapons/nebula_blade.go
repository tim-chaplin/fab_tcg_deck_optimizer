// Nebula Blade — Runeblade Weapon - Sword (2H). Cost 2, Power 1.
// Text: "Once per Turn Action - {r}{r}: Attack. If Nebula Blade hits, create a Runechant token. If
// you have played a 'non-attack' action card this turn, Nebula Blade gains +3{p} until end of
// turn."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

var nebulaBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type NebulaBlade struct{}

func (NebulaBlade) ID() ids.WeaponID    { return ids.NebulaBladeID }
func (NebulaBlade) Name() string        { return "Nebula Blade" }
func (NebulaBlade) Types() card.TypeSet { return nebulaBladeTypes }
func (NebulaBlade) Hands() int          { return 2 }
func (NebulaBlade) Ability() card.Card  { return nebulaBladeAbility }

// nebulaBladeAbility caches the Ability() return as a package-level card.Card so the
// chain runner's per-Best-call w.Ability() lookup doesn't re-box the zero-size struct
// in a fresh interface header on every call (would be 5x allocs per anneal otherwise).
var nebulaBladeAbility card.Card = NebulaBladeAbility{}

var nebulaBladeAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type NebulaBladeAbility struct{}

func (NebulaBladeAbility) ID() ids.CardID           { return ids.NebulaBladeAbilityID }
func (NebulaBladeAbility) Name() string             { return "Nebula Blade" }
func (NebulaBladeAbility) DisplayName() string      { return "Nebula Blade" }
func (NebulaBladeAbility) Cost(card.GameEngine) int { return 2 }
func (NebulaBladeAbility) Pitch() int               { return 0 }
func (NebulaBladeAbility) Attack() int              { return 1 }
func (NebulaBladeAbility) Defense() int             { return 0 }
func (NebulaBladeAbility) Types() card.TypeSet      { return nebulaBladeAbilityTypes }
func (NebulaBladeAbility) GoAgain() bool            { return false }
func (NebulaBladeAbility) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	if s.NonAttackActionPlayed() {
		self.BonusAttack += 3
	}
	self.RegisterOnHit(nebulaBladeOnHit)
}

// nebulaBladeOnHit fires the printed "If Nebula Blade hits, create a Runechant token"
// rider. Top-level so registration stays alloc-free.
func nebulaBladeOnHit(s card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	s.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a runechant", 1)
}
