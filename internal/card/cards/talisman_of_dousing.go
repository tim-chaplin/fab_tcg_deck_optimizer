// Talisman of Dousing — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Spellvoid 1**"
//
// Modelling: Play creates a card-sourced item subscribed to DamageAboutToBeTaken. The
// handler asks the engine to prevent 1 arcane damage; on any actual prevention the
// talisman is destroyed (the printed Spellvoid charge is consumed). When no arcane damage
// is incoming, PreventArcaneDamage returns 0 and the item survives — so on a physical-only
// damage moment the item is inert.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func dousingHandler(ge card.GameEngine, l card.Logger, it card.Item, _ triggertype.Type) {
	prevented := ge.PreventArcaneDamage(1)
	if prevented <= 0 {
		return
	}
	ge.AddValue(prevented)
	l.AppendPostTriggerf(it.CardName(), prevented, "Prevented %d arcane damage", prevented)
	it.Destroy(true)
}

func (TalismanOfDousingYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItem(self.Card, triggertype.DamageAboutToBeTaken, dousingHandler, false, nil)
}
