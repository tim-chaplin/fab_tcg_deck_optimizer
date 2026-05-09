// High Striker — Generic Action. Cost 0, Defense 2.
// Text: "The next time an attack you control hits this turn, create N Copper tokens.
// **Go again**" — N = 6 (Red) / 4 (Yellow) / 2 (Blue).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// highStrikerOnHit{6,4,2} fire on the next hit matching the trigger's TypeFilter (any
// attack per the printed wording). One top-level function per variant keeps the handler
// a static function value — no closure allocation per Play.
func highStrikerOnHit6(s *sim.TurnState, t *sim.Trigger, _ *sim.Aura) {
	highStrikerCreate(s, t, 6)
}
func highStrikerOnHit4(s *sim.TurnState, t *sim.Trigger, _ *sim.Aura) {
	highStrikerCreate(s, t, 4)
}
func highStrikerOnHit2(s *sim.TurnState, t *sim.Trigger, _ *sim.Aura) {
	highStrikerCreate(s, t, 2)
}

func highStrikerCreate(s *sim.TurnState, t *sim.Trigger, n int) {
	s.CreateCopper(n)
	s.LogPostTriggerf(sim.DisplayName(s.TriggeringCard), 0,
		"%s created %d copper tokens on attack hit", sim.DisplayName(t.Source), n)
}

func highStrikerPlay(s *sim.TurnState, self *sim.CardState, source sim.Card, handler sim.TriggerHandler) {
	s.AddTrigger(sim.Trigger{
		Source:      source,
		TriggerType: sim.TriggerHit,
		TypeFilter:  card.TypeSet.IsAttack,
		Handler:     handler,
	})
	s.Log(self, 0)
}

func (c HighStrikerRed) Play(s *sim.TurnState, self *sim.CardState) {
	highStrikerPlay(s, self, c, highStrikerOnHit6)
}

func (c HighStrikerYellow) Play(s *sim.TurnState, self *sim.CardState) {
	highStrikerPlay(s, self, c, highStrikerOnHit4)
}

func (c HighStrikerBlue) Play(s *sim.TurnState, self *sim.CardState) {
	highStrikerPlay(s, self, c, highStrikerOnHit2)
}
