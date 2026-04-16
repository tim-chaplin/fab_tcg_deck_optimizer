// Nebula Blade — Runeblade Weapon - Sword (2H). Cost 2, Power 1.
// Text: "Once per Turn Action - {r}{r}: Attack. If Nebula Blade hits, create a Runechant token. If
// you have played a 'non-attack' action card this turn, Nebula Blade gains +3{p} until end of
// turn."
//
// Simplification: assume the attack always hits (+1 damage for the Runechant, counted as future
// damage). The +3 power rider fires iff a prior card in CardsPlayed has Action && !Attack types.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package weapon

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var nebulaBladeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeSword, card.TypeTwoHand)

type NebulaBlade struct{}

func (NebulaBlade) ID() card.ID               { return card.NebulaBladeID }
func (NebulaBlade) Name() string              { return "Nebula Blade" }
func (NebulaBlade) Cost() int                 { return 2 }
func (NebulaBlade) Pitch() int                { return 0 }
func (NebulaBlade) Attack() int               { return 1 }
func (NebulaBlade) Defense() int              { return 0 }
func (NebulaBlade) Types() card.TypeSet        { return nebulaBladeTypes }
func (NebulaBlade) GoAgain() bool             { return false }
func (NebulaBlade) Hands() int                { return 2 }
func (c NebulaBlade) Play(s *card.TurnState) int {
	dmg := c.Attack() + s.CreateRunechant() // hit creates 1 Runechant (+1 future damage)
	for _, pc := range s.CardsPlayed {
		pt := pc.Types()
		if pt.Has(card.TypeAction) && !pt.Has(card.TypeAttack) {
			dmg += 3
			break
		}
	}
	return dmg
}
