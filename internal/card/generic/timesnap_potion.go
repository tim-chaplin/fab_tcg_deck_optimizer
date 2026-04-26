// Timesnap Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Action** - Destroy this: Gain 2 action points."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var timesnapPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type TimesnapPotionBlue struct{}

func (TimesnapPotionBlue) ID() card.ID              { return card.TimesnapPotionBlue }
func (TimesnapPotionBlue) Name() string             { return "Timesnap Potion" }
func (TimesnapPotionBlue) Cost(*card.TurnState) int { return 0 }
func (TimesnapPotionBlue) Pitch() int               { return 3 }
func (TimesnapPotionBlue) Attack() int              { return 0 }
func (TimesnapPotionBlue) Defense() int             { return 0 }
func (TimesnapPotionBlue) Types() card.TypeSet      { return timesnapPotionTypes }
func (TimesnapPotionBlue) GoAgain() bool            { return false }

// not implemented: activated 'gain 2 action points'
func (TimesnapPotionBlue) NotImplemented()                              {}
func (TimesnapPotionBlue) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
