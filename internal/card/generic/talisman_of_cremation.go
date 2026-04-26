// Talisman of Cremation — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Go again** When you play a card from your banished zone, destroy Talisman of Cremation
// and name a card. Banish all cards with the chosen name from each opposing hero's graveyard."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var talismanOfCremationTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TalismanOfCremationBlue struct{}

func (TalismanOfCremationBlue) ID() card.ID              { return card.TalismanOfCremationBlue }
func (TalismanOfCremationBlue) Name() string             { return "Talisman of Cremation" }
func (TalismanOfCremationBlue) Cost(*card.TurnState) int { return 0 }
func (TalismanOfCremationBlue) Pitch() int               { return 3 }
func (TalismanOfCremationBlue) Attack() int              { return 0 }
func (TalismanOfCremationBlue) Defense() int             { return 0 }
func (TalismanOfCremationBlue) Types() card.TypeSet      { return talismanOfCremationTypes }
func (TalismanOfCremationBlue) GoAgain() bool            { return true }

// not implemented: self-destroys on play-from-banished → banish a named card from opposing
// graveyards
func (TalismanOfCremationBlue) NotImplemented()                              {}
func (TalismanOfCremationBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
