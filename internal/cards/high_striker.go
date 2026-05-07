// High Striker — Generic Action. Cost 0, Defense 2.
// Text: "The next time an attack you control hits this turn, create N Copper tokens.
// **Go again**" — N = 6 (Red) / 4 (Yellow) / 2 (Blue).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var highStrikerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

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

type HighStrikerRed struct{}

func (HighStrikerRed) ID() ids.CardID          { return ids.HighStrikerRed }
func (HighStrikerRed) Name() string            { return "High Striker" }
func (HighStrikerRed) Cost(*sim.TurnState) int { return 0 }
func (HighStrikerRed) Pitch() int              { return 1 }
func (HighStrikerRed) Attack() int             { return 0 }
func (HighStrikerRed) Defense() int            { return 2 }
func (HighStrikerRed) Types() card.TypeSet     { return highStrikerTypes }
func (HighStrikerRed) GoAgain() bool           { return true }

func (c HighStrikerRed) Play(s *sim.TurnState, self *sim.CardState) {
	highStrikerPlay(s, self, c, highStrikerOnHit6)
}

type HighStrikerYellow struct{}

func (HighStrikerYellow) ID() ids.CardID          { return ids.HighStrikerYellow }
func (HighStrikerYellow) Name() string            { return "High Striker" }
func (HighStrikerYellow) Cost(*sim.TurnState) int { return 0 }
func (HighStrikerYellow) Pitch() int              { return 2 }
func (HighStrikerYellow) Attack() int             { return 0 }
func (HighStrikerYellow) Defense() int            { return 2 }
func (HighStrikerYellow) Types() card.TypeSet     { return highStrikerTypes }
func (HighStrikerYellow) GoAgain() bool           { return true }
func (c HighStrikerYellow) Play(s *sim.TurnState, self *sim.CardState) {
	highStrikerPlay(s, self, c, highStrikerOnHit4)
}

type HighStrikerBlue struct{}

func (HighStrikerBlue) ID() ids.CardID          { return ids.HighStrikerBlue }
func (HighStrikerBlue) Name() string            { return "High Striker" }
func (HighStrikerBlue) Cost(*sim.TurnState) int { return 0 }
func (HighStrikerBlue) Pitch() int              { return 3 }
func (HighStrikerBlue) Attack() int             { return 0 }
func (HighStrikerBlue) Defense() int            { return 2 }
func (HighStrikerBlue) Types() card.TypeSet     { return highStrikerTypes }
func (HighStrikerBlue) GoAgain() bool           { return true }
func (c HighStrikerBlue) Play(s *sim.TurnState, self *sim.CardState) {
	highStrikerPlay(s, self, c, highStrikerOnHit2)
}
