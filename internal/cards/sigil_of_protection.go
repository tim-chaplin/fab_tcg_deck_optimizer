// Sigil of Protection — Generic Action - Aura. Cost 1. Printed pitch variants: Red 1, Yellow 2,
// Blue 3. Defense 2.
//
// Text: "**Ward 4** At the beginning of your action phase, destroy Sigil of Protection."
//
// The aura-created flag is set so same-turn aura-readers see the entry.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var sigilOfProtectionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAura)

type SigilOfProtectionRed struct{}

func (SigilOfProtectionRed) ID() ids.CardID          { return ids.SigilOfProtectionRed }
func (SigilOfProtectionRed) Name() string            { return "Sigil of Protection" }
func (SigilOfProtectionRed) Cost(*sim.TurnState) int { return 1 }
func (SigilOfProtectionRed) Pitch() int              { return 1 }
func (SigilOfProtectionRed) Attack() int             { return 0 }
func (SigilOfProtectionRed) Defense() int            { return 2 }
func (SigilOfProtectionRed) Types() card.TypeSet     { return sigilOfProtectionTypes }
func (SigilOfProtectionRed) GoAgain() bool           { return false }

// not implemented: ward (opponent damage prevention)
func (SigilOfProtectionRed) NotImplemented() {}
func (SigilOfProtectionRed) Play(s *sim.TurnState, self *sim.CardState) {
	setAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type SigilOfProtectionYellow struct{}

func (SigilOfProtectionYellow) ID() ids.CardID          { return ids.SigilOfProtectionYellow }
func (SigilOfProtectionYellow) Name() string            { return "Sigil of Protection" }
func (SigilOfProtectionYellow) Cost(*sim.TurnState) int { return 1 }
func (SigilOfProtectionYellow) Pitch() int              { return 2 }
func (SigilOfProtectionYellow) Attack() int             { return 0 }
func (SigilOfProtectionYellow) Defense() int            { return 2 }
func (SigilOfProtectionYellow) Types() card.TypeSet     { return sigilOfProtectionTypes }
func (SigilOfProtectionYellow) GoAgain() bool           { return false }

// not implemented: ward (opponent damage prevention)
func (SigilOfProtectionYellow) NotImplemented() {}
func (SigilOfProtectionYellow) Play(s *sim.TurnState, self *sim.CardState) {
	setAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}

type SigilOfProtectionBlue struct{}

func (SigilOfProtectionBlue) ID() ids.CardID          { return ids.SigilOfProtectionBlue }
func (SigilOfProtectionBlue) Name() string            { return "Sigil of Protection" }
func (SigilOfProtectionBlue) Cost(*sim.TurnState) int { return 1 }
func (SigilOfProtectionBlue) Pitch() int              { return 3 }
func (SigilOfProtectionBlue) Attack() int             { return 0 }
func (SigilOfProtectionBlue) Defense() int            { return 2 }
func (SigilOfProtectionBlue) Types() card.TypeSet     { return sigilOfProtectionTypes }
func (SigilOfProtectionBlue) GoAgain() bool           { return false }

// not implemented: ward (opponent damage prevention)
func (SigilOfProtectionBlue) NotImplemented() {}
func (SigilOfProtectionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	setAuraCreated(s)
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)
}
