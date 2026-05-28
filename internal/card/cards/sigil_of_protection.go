// Sigil of Protection — Generic Action - Aura. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Ward 4** At the beginning of your action phase, destroy Sigil of Protection."
//
// Ward 4 is modelled as a one-shot DamageAboutToBeTaken absorb capped at 4; firing
// dispatches on the trigger type so the StartOfTurn clause only destroys.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func sigilOfProtectionHandler(ge card.GameEngine, l card.Logger, a card.Aura, firingType triggertype.Type) {
	switch firingType {
	case triggertype.DamageAboutToBeTaken:
		prevented := ge.PreventIncomingDamage(a.Count())
		if prevented > 0 {
			ge.AddValue(prevented)
			l.AppendPostTriggerf(a.CardName(), prevented, "Prevented %d damage", prevented)
		}
		a.Destroy(true)
	case triggertype.StartOfTurn:
		a.Destroy(true)
	}
}

func sigilOfProtectionPlay(ge card.GameEngine, self *card.CardState) {
	ge.CreateAura(self.Card, triggertype.DamageAboutToBeTaken|triggertype.StartOfTurn,
		sigilOfProtectionHandler, 4, false, nil)
}

func (SigilOfProtectionRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sigilOfProtectionPlay(ge, self)
}

func (SigilOfProtectionYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sigilOfProtectionPlay(ge, self)
}

func (SigilOfProtectionBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	sigilOfProtectionPlay(ge, self)
}
