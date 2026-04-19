// Down But Not Out — Generic Action - Attack. Cost 3. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 3.
//
// Text: "When this attacks a hero, if you have less {h} and control fewer equipment and tokens than
// them, this gets +3{p}, **overpower**, and "When this hits, create an Agility, Might, and Vigor
// token.""
//
// Simplification: Health/equipment/token comparison isn't modelled; none of the riders fire.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var downButNotOutTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type DownButNotOutRed struct{}

func (DownButNotOutRed) ID() card.ID                 { return card.DownButNotOutRed }
func (DownButNotOutRed) Name() string                { return "Down But Not Out (Red)" }
func (DownButNotOutRed) Cost() int                   { return 3 }
func (DownButNotOutRed) Pitch() int                  { return 1 }
func (DownButNotOutRed) Attack() int                 { return 5 }
func (DownButNotOutRed) Defense() int                { return 3 }
func (DownButNotOutRed) Types() card.TypeSet         { return downButNotOutTypes }
func (DownButNotOutRed) GoAgain() bool               { return false }
func (c DownButNotOutRed) Play(s *card.TurnState) int { return downButNotOutDamage(c.Attack()) }

type DownButNotOutYellow struct{}

func (DownButNotOutYellow) ID() card.ID                 { return card.DownButNotOutYellow }
func (DownButNotOutYellow) Name() string                { return "Down But Not Out (Yellow)" }
func (DownButNotOutYellow) Cost() int                   { return 3 }
func (DownButNotOutYellow) Pitch() int                  { return 2 }
func (DownButNotOutYellow) Attack() int                 { return 4 }
func (DownButNotOutYellow) Defense() int                { return 3 }
func (DownButNotOutYellow) Types() card.TypeSet         { return downButNotOutTypes }
func (DownButNotOutYellow) GoAgain() bool               { return false }
func (c DownButNotOutYellow) Play(s *card.TurnState) int { return downButNotOutDamage(c.Attack()) }

type DownButNotOutBlue struct{}

func (DownButNotOutBlue) ID() card.ID                 { return card.DownButNotOutBlue }
func (DownButNotOutBlue) Name() string                { return "Down But Not Out (Blue)" }
func (DownButNotOutBlue) Cost() int                   { return 3 }
func (DownButNotOutBlue) Pitch() int                  { return 3 }
func (DownButNotOutBlue) Attack() int                 { return 3 }
func (DownButNotOutBlue) Defense() int                { return 3 }
func (DownButNotOutBlue) Types() card.TypeSet         { return downButNotOutTypes }
func (DownButNotOutBlue) GoAgain() bool               { return false }
func (c DownButNotOutBlue) Play(s *card.TurnState) int { return downButNotOutDamage(c.Attack()) }

// downButNotOutDamage is a breadcrumb for the conditional "when this hits, create Agility +
// Might + Vigor tokens" rider — gated on a health/equipment/token comparison we don't track
// (see TODO.md).
func downButNotOutDamage(attack int) int {
	if card.LikelyToHit(attack) {
		// TODO: model on-hit status-token creation rider (requires life-total + token tracking).
	}
	return attack
}
