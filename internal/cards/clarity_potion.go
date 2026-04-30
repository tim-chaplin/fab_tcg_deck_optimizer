// Clarity Potion — Generic Action - Item. Cost 0. Printed pitch variants: Blue 3.
//
// Text: "**Instant** - Destroy Clarity Potion: **Opt 2**"
//
// Modelled by crediting Opt 2 immediately on play. The activated ability is a free
// (0 AP) self-destroy that yields Opt 2; optimal play always activates it on the same
// turn it's played, so the on-play credit matches the every-turn outcome. The Item
// persists in the play zone afterward (per the framework's PersistsInPlay rule); no
// production card reads "Clarity Potion in play", so the lingering Item is benign.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var clarityPotionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)

type ClarityPotionBlue struct{}

func (ClarityPotionBlue) ID() ids.CardID          { return ids.ClarityPotionBlue }
func (ClarityPotionBlue) Name() string            { return "Clarity Potion" }
func (ClarityPotionBlue) Cost(*sim.TurnState) int { return 0 }
func (ClarityPotionBlue) Pitch() int              { return 3 }
func (ClarityPotionBlue) Attack() int             { return 0 }
func (ClarityPotionBlue) Defense() int            { return 0 }
func (ClarityPotionBlue) Types() card.TypeSet     { return clarityPotionTypes }
func (ClarityPotionBlue) GoAgain() bool           { return false }
func (ClarityPotionBlue) Play(s *sim.TurnState, self *sim.CardState) {
	s.LogPlay(self)
	s.ApplyAndLogRiderOnPlay(self, "Opt 2", 2*sim.OptValue)
}
