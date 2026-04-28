// Eirina's Prayer — Generic Instant. Cost 1. Printed pitch variants: Red 1, Yellow 2, Blue 3.
//
// Text: "Reveal the top card of your deck. Prevent the next X arcane damage that would be dealt to
// your hero this turn, where X is 6 minus the pitch value of the card revealed this way."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var eirinasPrayerTypes = card.NewTypeSet(card.TypeGeneric, card.TypeInstant)

type EirinasPrayerRed struct{}

func (EirinasPrayerRed) ID() ids.CardID           { return ids.EirinasPrayerRed }
func (EirinasPrayerRed) Name() string             { return "Eirina's Prayer" }
func (EirinasPrayerRed) Cost(*card.TurnState) int { return 1 }
func (EirinasPrayerRed) Pitch() int               { return 1 }
func (EirinasPrayerRed) Attack() int              { return 0 }
func (EirinasPrayerRed) Defense() int             { return 0 }
func (EirinasPrayerRed) Types() card.TypeSet      { return eirinasPrayerTypes }
func (EirinasPrayerRed) GoAgain() bool            { return false }

// not implemented: Instant prevent X arcane to your hero; X scaled by revealed top-card pitch
func (EirinasPrayerRed) NotImplemented()                              {}
func (EirinasPrayerRed) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type EirinasPrayerYellow struct{}

func (EirinasPrayerYellow) ID() ids.CardID           { return ids.EirinasPrayerYellow }
func (EirinasPrayerYellow) Name() string             { return "Eirina's Prayer" }
func (EirinasPrayerYellow) Cost(*card.TurnState) int { return 1 }
func (EirinasPrayerYellow) Pitch() int               { return 2 }
func (EirinasPrayerYellow) Attack() int              { return 0 }
func (EirinasPrayerYellow) Defense() int             { return 0 }
func (EirinasPrayerYellow) Types() card.TypeSet      { return eirinasPrayerTypes }
func (EirinasPrayerYellow) GoAgain() bool            { return false }

// not implemented: Instant prevent X arcane to your hero; X scaled by revealed top-card pitch
func (EirinasPrayerYellow) NotImplemented()                              {}
func (EirinasPrayerYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }

type EirinasPrayerBlue struct{}

func (EirinasPrayerBlue) ID() ids.CardID           { return ids.EirinasPrayerBlue }
func (EirinasPrayerBlue) Name() string             { return "Eirina's Prayer" }
func (EirinasPrayerBlue) Cost(*card.TurnState) int { return 1 }
func (EirinasPrayerBlue) Pitch() int               { return 3 }
func (EirinasPrayerBlue) Attack() int              { return 0 }
func (EirinasPrayerBlue) Defense() int             { return 0 }
func (EirinasPrayerBlue) Types() card.TypeSet      { return eirinasPrayerTypes }
func (EirinasPrayerBlue) GoAgain() bool            { return false }

// not implemented: Instant prevent X arcane to your hero; X scaled by revealed top-card pitch
func (EirinasPrayerBlue) NotImplemented()                              {}
func (EirinasPrayerBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
