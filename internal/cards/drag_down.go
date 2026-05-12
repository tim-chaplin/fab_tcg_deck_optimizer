// Drag Down — Generic Defense Reaction. Cost 0. Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Defense 0.
//
// Text: "When this defends an attack, it gets -3{p}."
//
// The -3{p} attacker debuff is modelled as a flat 3/2/1{d} block on R/Y/B respectively —
// equivalent to "block N damage" against an attack chain whose damage is consumed by
// IncomingDamage / BlockTotal. The Red→Blue decay matches the standard pitch-vs-body
// trade-off across pitch variants.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (DragDownRed) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (DragDownYellow) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}

func (DragDownBlue) Play(g card.GameEngine, l card.Logger, self *card.CardState) {}
