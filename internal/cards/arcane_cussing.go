// Arcane Cussing — Runeblade Action - Aura. Cost 1, Defense 2, Go again.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "Go again. When you deal or are dealt damage, destroy this. When this leaves the arena
// during your turn, create N Runechants." (Red N=3, Yellow N=2, Blue N=1.)
//
// Modelled as a fragile aura (fragile_aura.go). The "deal damage" trigger lets any attacker
// type pop the aura, so Play passes attackActionOnly=false.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

func (ArcaneCussingRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	fragileAuraPlay(s, l, self, 3, false)
}

func (ArcaneCussingYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	fragileAuraPlay(s, l, self, 2, false)
}

func (ArcaneCussingBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	fragileAuraPlay(s, l, self, 1, false)
}
