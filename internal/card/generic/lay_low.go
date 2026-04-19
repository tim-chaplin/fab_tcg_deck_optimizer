// Lay Low — Generic Defense Reaction. Cost 0, Pitch 2, Defense 3. Only printed in Yellow.
// Text: "If you are marked, you can't play this. If the defending hero is marked, their next attack
// this turn gets -1{p}."
// Simplification: marked-hero state isn't tracked; we assume the defender is never marked (so the
// card is always legal) and ignore the attacker debuff.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

type LayLowYellow struct{}

func (LayLowYellow) ID() card.ID                 { return card.LayLowYellow }
func (LayLowYellow) Name() string             { return "Lay Low (Yellow)" }
func (LayLowYellow) Cost(*card.TurnState) int                { return 0 }
func (LayLowYellow) Pitch() int               { return 2 }
func (LayLowYellow) Attack() int              { return 0 }
func (LayLowYellow) Defense() int             { return 3 }
func (LayLowYellow) Types() card.TypeSet      { return defenseReactionTypes }
func (LayLowYellow) GoAgain() bool            { return false }
func (LayLowYellow) Play(*card.TurnState) int { return 0 }
