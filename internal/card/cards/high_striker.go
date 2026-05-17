// High Striker — Generic Action. Cost 0, Defense 2.
// Text: "The next time an attack you control hits this turn, create N Copper tokens.
// **Go again**" — N = 6 (Red) / 4 (Yellow) / 2 (Blue).

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// highStrikerOnHit{6,4,2} fire on the next hit matching the trigger's TypeFilter (any
// attack per the printed wording). One handler per variant carries N in the function name.
func highStrikerOnHit6(ge card.GameEngine, l card.Logger, t card.Trigger) {
	highStrikerCreate(ge, l, t, 6)
}
func highStrikerOnHit4(ge card.GameEngine, l card.Logger, t card.Trigger) {
	highStrikerCreate(ge, l, t, 4)
}
func highStrikerOnHit2(ge card.GameEngine, l card.Logger, t card.Trigger) {
	highStrikerCreate(ge, l, t, 2)
}

func highStrikerCreate(ge card.GameEngine, l card.Logger, t card.Trigger, n int) {
	ge.CreateCopper(n)
	l.AppendPostTriggerf(ge.TriggeringCard().DisplayName(), 0,
		"%s created %d copper tokens on attack hit", t.CardName(), n)
}

func (c HighStrikerRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.AddHitTrigger(self, highStrikerOnHit6, card.TypeSet.IsAttack)
}

func (c HighStrikerYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.AddHitTrigger(self, highStrikerOnHit4, card.TypeSet.IsAttack)
}

func (c HighStrikerBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.AddHitTrigger(self, highStrikerOnHit2, card.TypeSet.IsAttack)
}
