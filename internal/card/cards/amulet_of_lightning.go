// Amulet of Lightning — Lightning Action - Item. Cost 0. Pitch 3 (Blue). Common. Go again.
//
// Text: "**Go again** / **Instant** - Destroy Amulet of Lightning: Target action card gains
// **go again**. Activate this ability only if you have Lightning **fused** this turn."
//
// Modelling:
//   - The card's printed Go again belongs to the card itself (a 0-cost item play that refunds
//     its action point); it's emitted by cardgen, not on the ability.
//   - Play creates the in-play item; the printed instant ability lives on
//     AmuletOfLightningAbility, enumerated by the wmask. The ability is an Instant (0 action
//     points — no Go again needed) gated by PlayPrecondition on having Lightning fused this
//     turn; on resolve it grants Go again to the next action card still to be played and pays
//     its "Destroy this" cost.
//   - No implemented card sets HasLightningFused, so the gate never opens for the current pool.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

func (c AmuletOfLightningBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItemWithAbility(c, AmuletOfLightningAbility{})
}

// AmuletOfLightningAbility is Amulet of Lightning's printed instant ability. It's an Instant
// (resolves without paying an action point), so it carries no Go again of its own.
type AmuletOfLightningAbility struct{}

func (AmuletOfLightningAbility) ID() ids.CardID      { return ids.AmuletOfLightningAbilityID }
func (AmuletOfLightningAbility) Name() string        { return "Amulet of Lightning" }
func (AmuletOfLightningAbility) DisplayName() string { return "Amulet of Lightning" }
func (AmuletOfLightningAbility) Cost() int           { return 0 }
func (AmuletOfLightningAbility) Pitch() int          { return 0 }
func (AmuletOfLightningAbility) Attack() int         { return 0 }
func (AmuletOfLightningAbility) Defense() int        { return 0 }
func (AmuletOfLightningAbility) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeLightning, card.TypeItem, card.TypeInstant)
}
func (AmuletOfLightningAbility) GoAgain(card.GameEngine) bool { return false }

// PlayPrecondition gates activation on "only if you have Lightning fused this turn".
func (AmuletOfLightningAbility) PlayPrecondition(ge card.GameEngine, _ *card.CardState) bool {
	return ge.HasLightningFused()
}

func (AmuletOfLightningAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	// Grant go again to the next action card still to be played that doesn't already have it —
	// the optimiser orders the sequence so the grant lands where it adds an action point.
	for _, pc := range ge.CardsRemaining() {
		if pc.Card.Types(ge).Has(card.TypeAction) && !pc.EffectiveGoAgain(ge) {
			pc.GrantedGoAgain = true
			break
		}
	}
	self.Item.Destroy(true) // pay the "Destroy Amulet of Lightning" cost
}
