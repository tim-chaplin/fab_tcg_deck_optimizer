// Tip-Off — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3. Printed pitch
// variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "**Instant** - Discard this: **Mark** target opposing hero."
//
// Modelled as a two-mode card where mode 1's type-line is "Generic Instant" instead of the
// printed "Generic Action - Attack" (via card.ModalTypes). The engine reads the per-mode
// type-line: mode 1 pays 0 AP (Instant is a free attack step), credits no attack damage,
// doesn't clear OpponentMarked, and won't match "next attack action card" predicates
// scanning CardsRemaining. Play just marks the opponent — the rest falls out of the type
// dispatch.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

var tipOffInstantTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

func tipOffTypesForMode(_ card.GameEngine, mode int8) card.TypeSet {
	if mode == 0 {
		return tipOffTypes
	}
	return tipOffInstantTypes
}

func tipOffPlay(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if self.Mode == 0 {
		return
	}
	ge.MarkOpponent()
}

func tipOffModalCost(mode int8) int {
	if mode == 0 {
		return 1
	}
	return 0
}

func (TipOffRed) Modes() int              { return 2 }
func (TipOffRed) ModalCost(mode int8) int { return tipOffModalCost(mode) }
func (TipOffRed) TypesForMode(g card.GameEngine, m int8) card.TypeSet {
	return tipOffTypesForMode(g, m)
}
func (TipOffRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tipOffPlay(ge, l, self)
}

func (TipOffYellow) Modes() int              { return 2 }
func (TipOffYellow) ModalCost(mode int8) int { return tipOffModalCost(mode) }
func (TipOffYellow) TypesForMode(g card.GameEngine, m int8) card.TypeSet {
	return tipOffTypesForMode(g, m)
}
func (TipOffYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tipOffPlay(ge, l, self)
}

func (TipOffBlue) Modes() int              { return 2 }
func (TipOffBlue) ModalCost(mode int8) int { return tipOffModalCost(mode) }
func (TipOffBlue) TypesForMode(g card.GameEngine, m int8) card.TypeSet {
	return tipOffTypesForMode(g, m)
}
func (TipOffBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	tipOffPlay(ge, l, self)
}
