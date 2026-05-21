// Regain Composure — Generic Action. Cost 0, Defense 2. Printed only in Blue (pitch 3).
// Text: "Your next attack this turn gets +1{p} and "When this hits, {u} your hero." **Go again**"

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (RegainComposureBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	for _, pc := range ge.CardsRemaining() {
		if !card.IsAttack(ge, pc) {
			continue
		}
		pc.BonusAttack++
		pc.OnHit = append(pc.OnHit, card.OnHitHandler{Fire: regainComposureUntapHero, Source: self.Card})
		return
	}
}

// regainComposureUntapHero is the "When this hits, {u} your hero" clause granted to the
// buffed attack — {u} untaps the hero. Source is Regain Composure, not the attack it rides.
func regainComposureUntapHero(ge card.GameEngine, l card.Logger, _ *card.CardState, h *card.OnHitHandler) {
	ge.UntapHero()
	l.AppendPostTrigger(h.Source.DisplayName(), "Untapped your hero", 0)
}
