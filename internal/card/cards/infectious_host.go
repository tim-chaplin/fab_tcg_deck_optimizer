// Infectious Host — Generic Action - Attack. Cost 0. Printed power: Red 4, Yellow 3, Blue 2.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks a hero, if you control a Frailty token, create a Frailty token under
// their control, then repeat for Inertia and Bloodrot Pox."
//
// Each clause gates on us controlling the matching status token and mints the same kind on
// the opponent. CreateXForOpponent credits a flat heuristic value via AddValue; we don't
// track opposing-side status-token state. Today no card grants self-side Frailty / Inertia /
// Bloodrot Pox, so all three gates start false — the framework is in place for future
// self-side granters.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func infectiousHostPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	source := self.Card.DisplayName()
	if ge.FrailtyCount() > 0 {
		ge.CreateFrailtyForOpponent()
		l.AppendPostTrigger(source, "Spread a Frailty token to opponent", 0)
	}
	if ge.InertiaCount() > 0 {
		ge.CreateInertiaForOpponent()
		l.AppendPostTrigger(source, "Spread an Inertia token to opponent", 0)
	}
	if ge.BloodrotPoxCount() > 0 {
		ge.CreateBloodrotPoxForOpponent()
		l.AppendPostTrigger(source, "Spread a Bloodrot Pox token to opponent", 0)
	}
}

func (InfectiousHostRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	infectiousHostPlay(ge, l, self)
}

func (InfectiousHostYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	infectiousHostPlay(ge, l, self)
}

func (InfectiousHostBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	infectiousHostPlay(ge, l, self)
}
