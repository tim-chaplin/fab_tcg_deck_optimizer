// Nebula Blade — Runeblade Weapon - Sword (2H). Cost 2, Power 1.
// Text: "Once per Turn Action - {r}{r}: Attack. If Nebula Blade hits, create a Runechant token. If
// you have played a 'non-attack' action card this turn, Nebula Blade gains +3{p} until end of
// turn."

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

var nebulaBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

// NebulaBlade is the platonic weapon card — the equipped permanent that sits in the arena
// and never enters the attack turn. Its activated ability (NebulaBladeAbility) is what the
// attack-turn runner enqueues each turn.
type NebulaBlade struct{}

func (NebulaBlade) ID() ids.CardID                     { return ids.NebulaBladeID }
func (NebulaBlade) Name() string                       { return "Nebula Blade" }
func (NebulaBlade) DisplayName() string                { return "Nebula Blade" }
func (NebulaBlade) Rarity() string                     { return "Common" }
func (NebulaBlade) Cost() int                          { return 0 }
func (NebulaBlade) Pitch() int                         { return 0 }
func (NebulaBlade) Attack() int                        { return 0 }
func (NebulaBlade) Defense() int                       { return 0 }
func (NebulaBlade) Types(card.GameEngine) card.TypeSet { return nebulaBladeTypes }
func (NebulaBlade) GoAgain(card.GameEngine) bool       { return false }
func (NebulaBlade) Hands() int                         { return 2 }
func (NebulaBlade) Ability() card.Card                 { return nebulaBladeAbility }

func (NebulaBlade) Play(ge card.GameEngine, _ card.Logger, self *card.CardState) {
	ge.CreateWeapon(self.Card, 0, nil, false, nil)
}

var nebulaBladeAbility card.Card = NebulaBladeAbility{}

var nebulaBladeAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand, card.TypeAttack)

type NebulaBladeAbility struct{}

func (NebulaBladeAbility) ID() ids.CardID                     { return ids.NebulaBladeAbilityID }
func (NebulaBladeAbility) Name() string                       { return "Nebula Blade" }
func (NebulaBladeAbility) DisplayName() string                { return "Nebula Blade" }
func (NebulaBladeAbility) Cost() int                          { return 2 }
func (NebulaBladeAbility) Pitch() int                         { return 0 }
func (NebulaBladeAbility) Attack() int                        { return 1 }
func (NebulaBladeAbility) Defense() int                       { return 0 }
func (NebulaBladeAbility) Types(card.GameEngine) card.TypeSet { return nebulaBladeAbilityTypes }
func (NebulaBladeAbility) GoAgain(card.GameEngine) bool       { return false }
func (NebulaBladeAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if ge.NonAttackActionPlayed() {
		self.BonusAttack += 3
	}
	self.RegisterOnHit(nebulaBladeOnHit)
}

// nebulaBladeOnHit fires the printed "If Nebula Blade hits, create a Runechant token" rider.
func nebulaBladeOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	ge.CreateRunechants(1)
	l.AppendPostTrigger(self.Card.DisplayName(), "On-hit created a runechant", 1)
}
