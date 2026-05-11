// Sigil of Protection — Generic Action - Aura. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Ward 4** At the beginning of your action phase, destroy Sigil of Protection."
//
// The aura-created flag is set so same-turn aura-readers see the entry.

package notimplemented

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// not implemented: ward (opponent damage prevention)

func (SigilOfProtectionRed) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.SetAuraCreated(true)
}

// not implemented: ward (opponent damage prevention)

func (SigilOfProtectionYellow) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.SetAuraCreated(true)
}

// not implemented: ward (opponent damage prevention)

func (SigilOfProtectionBlue) Play(s card.GameEngine, l card.Logger, self *card.CardState) {
	s.SetAuraCreated(true)
}
