// Amulet of Assertiveness — Generic Action - Item. Cost 0. Printed pitch variants: Yellow 2.
//
// Text: "**Go again** **Attack Reaction** - Destroy Amulet of Assertiveness: Target attack gains
// "When this hits, banish the top card of your deck. If it's an attack action card, you may play it
// this turn." Activate this ability only if you have 4 or more cards in hand."

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var amuletOfAssertivenessTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type AmuletOfAssertivenessYellow struct{}

func (AmuletOfAssertivenessYellow) ID() card.ID              { return card.AmuletOfAssertivenessYellow }
func (AmuletOfAssertivenessYellow) Name() string             { return "Amulet of Assertiveness" }
func (AmuletOfAssertivenessYellow) Cost(*card.TurnState) int { return 0 }
func (AmuletOfAssertivenessYellow) Pitch() int               { return 2 }
func (AmuletOfAssertivenessYellow) Attack() int              { return 0 }
func (AmuletOfAssertivenessYellow) Defense() int             { return 0 }
func (AmuletOfAssertivenessYellow) Types() card.TypeSet      { return amuletOfAssertivenessTypes }
func (AmuletOfAssertivenessYellow) GoAgain() bool            { return true }

// not implemented: AR grant: target attack 'banish top of deck on hit'; gated on 4+ cards
// in hand
func (AmuletOfAssertivenessYellow) NotImplemented()                              {}
func (AmuletOfAssertivenessYellow) Play(s *card.TurnState, self *card.CardState) { s.LogPlay(self) }
