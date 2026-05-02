// Plunder Run — Generic Action. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 2.
//
// Text: "The next time an attack action card you control hits this turn, draw a card. If
// Plunder Run is played from arsenal, the next attack action card you play this turn gains
// +N{p}. **Go again**" (Red N=3, Yellow N=2, Blue N=1.)
//
// Two riders, both routed through existing primitives:
//   - "Next time an attack action hits" — registers a NextAttackActionHitTrigger on
//     TurnState; the chain runner drains the queue inside finalizeActiveAttack on the first
//     attack action that lands. Multiple Plunder Runs queue independent draws and all fire
//     on the same hit.
//   - "From-arsenal +N{p}" — only fires when self.FromArsenal; uses the shared
//     GrantNextAttackActionBonus helper to attach the buff to the next attack action in
//     CardsRemaining.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var plunderRunTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// plunderRunOnHitDraw fires the printed "the next time an attack action card you control
// hits this turn, draw a card" rider. target is the attack action that landed; Source on
// the trigger names the Plunder Run printing for log attribution.
func plunderRunOnHitDraw(s *sim.TurnState, target *sim.CardState, t *sim.NextAttackActionHitTrigger) {
	s.DrawOne()
	s.LogPostTriggerf(sim.DisplayName(target.Card), 0,
		"%s drew a card on attack-action hit", sim.DisplayName(t.Source))
}

func plunderRunPlay(s *sim.TurnState, self *sim.CardState, source sim.Card, n int) {
	s.RegisterNextAttackActionHit(sim.NextAttackActionHitTrigger{
		Fire:   plunderRunOnHitDraw,
		Source: source,
	})
	if self.FromArsenal {
		GrantNextAttackActionBonus(s, n)
	}
	s.Log(self, 0)
}

type PlunderRunRed struct{}

func (PlunderRunRed) ID() ids.CardID          { return ids.PlunderRunRed }
func (PlunderRunRed) Name() string            { return "Plunder Run" }
func (PlunderRunRed) Cost(*sim.TurnState) int { return 0 }
func (PlunderRunRed) Pitch() int              { return 1 }
func (PlunderRunRed) Attack() int             { return 0 }
func (PlunderRunRed) Defense() int            { return 2 }
func (PlunderRunRed) Types() card.TypeSet     { return plunderRunTypes }
func (PlunderRunRed) GoAgain() bool           { return true }
func (PlunderRunRed) NotSilverAgeLegal()      {}
func (c PlunderRunRed) Play(s *sim.TurnState, self *sim.CardState) {
	plunderRunPlay(s, self, c, 3)
}

type PlunderRunYellow struct{}

func (PlunderRunYellow) ID() ids.CardID          { return ids.PlunderRunYellow }
func (PlunderRunYellow) Name() string            { return "Plunder Run" }
func (PlunderRunYellow) Cost(*sim.TurnState) int { return 0 }
func (PlunderRunYellow) Pitch() int              { return 2 }
func (PlunderRunYellow) Attack() int             { return 0 }
func (PlunderRunYellow) Defense() int            { return 2 }
func (PlunderRunYellow) Types() card.TypeSet     { return plunderRunTypes }
func (PlunderRunYellow) GoAgain() bool           { return true }
func (PlunderRunYellow) NotSilverAgeLegal()      {}
func (c PlunderRunYellow) Play(s *sim.TurnState, self *sim.CardState) {
	plunderRunPlay(s, self, c, 2)
}

type PlunderRunBlue struct{}

func (PlunderRunBlue) ID() ids.CardID          { return ids.PlunderRunBlue }
func (PlunderRunBlue) Name() string            { return "Plunder Run" }
func (PlunderRunBlue) Cost(*sim.TurnState) int { return 0 }
func (PlunderRunBlue) Pitch() int              { return 3 }
func (PlunderRunBlue) Attack() int             { return 0 }
func (PlunderRunBlue) Defense() int            { return 2 }
func (PlunderRunBlue) Types() card.TypeSet     { return plunderRunTypes }
func (PlunderRunBlue) GoAgain() bool           { return true }
func (PlunderRunBlue) NotSilverAgeLegal()      {}
func (c PlunderRunBlue) Play(s *sim.TurnState, self *sim.CardState) {
	plunderRunPlay(s, self, c, 1)
}
