// Jack Be Nimble — Generic Action - Attack. Cost 0, Pitch 1, Power 3, Defense 3. Only printed in
// Red.
//
// Text: "When this attacks, you may banish a Nimblism from your graveyard. If you do, this gets
// +1{p} and **go again**. When this hits a hero, steal an item they control until the end of this
// action phase."
//
// Simplification: Graveyard banish for +1{p}/go-again and on-hit steal aren't modelled.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var jackBeNimbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type JackBeNimbleRed struct{}

func (JackBeNimbleRed) ID() card.ID                 { return card.JackBeNimbleRed }
func (JackBeNimbleRed) Name() string                { return "Jack Be Nimble (Red)" }
func (JackBeNimbleRed) Cost(*card.TurnState) int                   { return 0 }
func (JackBeNimbleRed) Pitch() int                  { return 1 }
func (JackBeNimbleRed) Attack() int                 { return 3 }
func (JackBeNimbleRed) Defense() int                { return 3 }
func (JackBeNimbleRed) Types() card.TypeSet         { return jackBeNimbleTypes }
func (JackBeNimbleRed) GoAgain() bool               { return false }
func (c JackBeNimbleRed) Play(s *card.TurnState, _ *card.CardState) int { return jackBeNimbleDamage(c.Attack()) }

// jackBeNimbleDamage is a breadcrumb for the on-hit "steal an item" rider — not modelled yet
// (see TODO.md). The LikelyToHit call marks where the rider value would plug in.
func jackBeNimbleDamage(attack int) int {
	if card.LikelyToHit(attack, false) {
		// TODO: model on-hit steal-item rider.
	}
	return attack
}
