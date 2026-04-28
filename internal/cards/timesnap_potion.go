// Timesnap Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Gain 2 action points."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var timesnapPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TimesnapPotionBlue struct{}

func (TimesnapPotionBlue) ID() ids.CardID          { return ids.TimesnapPotionBlue }
func (TimesnapPotionBlue) Name() string            { return "Timesnap Potion" }
func (TimesnapPotionBlue) Cost(*sim.TurnState) int { return 0 }
func (TimesnapPotionBlue) Pitch() int              { return 3 }
func (TimesnapPotionBlue) Attack() int             { return 0 }
func (TimesnapPotionBlue) Defense() int            { return 0 }
func (TimesnapPotionBlue) Types() card.TypeSet     { return timesnapPotionTypes }
func (TimesnapPotionBlue) GoAgain() bool           { return false }

// not implemented: activated 'gain 2 action points'
func (TimesnapPotionBlue) NotImplemented()                            {}
func (TimesnapPotionBlue) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }
