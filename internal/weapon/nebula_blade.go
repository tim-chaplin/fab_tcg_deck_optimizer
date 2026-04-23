// Nebula Blade — Runeblade Weapon - Sword (2H). Cost 2, Power 1.
// Text: "Once per Turn Action - {r}{r}: Attack. If Nebula Blade hits, create a Runechant token. If
// you have played a 'non-attack' action card this turn, Nebula Blade gains +3{p} until end of
// turn."
//
// Modelling: the +3 power rider fires when s.NonAttackActionPlayed is set. The runechant rider
// gates on card.LikelyToHit of the effective attack — today's heuristic lets both the base-1 and
// buffed-4 swings qualify, so behavior matches "always create a rune", but gating explicitly
// tracks any future retune of LikelyToHit without a follow-up patch.

package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nebulaBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type NebulaBlade struct{}

func (NebulaBlade) ID() card.ID               { return card.NebulaBladeID }
func (NebulaBlade) Name() string              { return "Nebula Blade" }
func (NebulaBlade) Cost(*card.TurnState) int                 { return 2 }
func (NebulaBlade) Pitch() int                { return 0 }
func (NebulaBlade) Attack() int               { return 1 }
func (NebulaBlade) Defense() int              { return 0 }
func (NebulaBlade) Types() card.TypeSet        { return nebulaBladeTypes }
func (NebulaBlade) GoAgain() bool             { return false }
func (NebulaBlade) Hands() int                { return 2 }
func (c NebulaBlade) Play(s *card.TurnState, _ *card.CardState) int {
	dmg := c.Attack()
	if s.NonAttackActionPlayed {
		dmg += 3
	}
	if card.LikelyToHit(dmg, false) {
		dmg += s.CreateRunechant()
	}
	return dmg
}
