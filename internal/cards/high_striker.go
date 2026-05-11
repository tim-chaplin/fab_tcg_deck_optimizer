// High Striker — Generic Action. Cost 0, Defense 2.
// Text: "The next time an attack you control hits this turn, create N Copper tokens.
// **Go again**" — N = 6 (Red) / 4 (Yellow) / 2 (Blue).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// highStrikerOnHit{6,4,2} fire on the next hit matching the trigger's TypeFilter (any
// attack per the printed wording). One top-level function per variant keeps the handler
// a static function value — no closure allocation per Play.
func highStrikerOnHit6(s card.GameEngine, l card.Logger, t card.Trigger) {
	highStrikerCreate(s, l, t, 6)
}
func highStrikerOnHit4(s card.GameEngine, l card.Logger, t card.Trigger) {
	highStrikerCreate(s, l, t, 4)
}
func highStrikerOnHit2(s card.GameEngine, l card.Logger, t card.Trigger) {
	highStrikerCreate(s, l, t, 2)
}

func highStrikerCreate(s card.GameEngine, l card.Logger, t card.Trigger, n int) {
	s.CreateCopper(n)
	l.AppendPostTriggerf(s.TriggeringCard().DisplayName(), 0,
		"%s created %d copper tokens on attack hit", t.SourceName(), n)
}

func (c HighStrikerRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.AddHitTrigger(self, highStrikerOnHit6, card.TypeSet.IsAttack)
}

func (c HighStrikerYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.AddHitTrigger(self, highStrikerOnHit4, card.TypeSet.IsAttack)
}

func (c HighStrikerBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.AddHitTrigger(self, highStrikerOnHit2, card.TypeSet.IsAttack)
}
