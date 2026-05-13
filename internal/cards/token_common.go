// Shared helpers backing the Gold / Silver / Copper / Runechant / Ponder token files —
// the type-set every item token's activated-ability card reports, and the narrow engine
// surface those abilities reach for to destroy the consumed token entry.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// tokenAbilityConsumer is the slice of engine surface item-token activated abilities use
// to destroy the consumed entry. *gameengine.GameEngine satisfies it structurally.
type tokenAbilityConsumer interface {
	ConsumeItemByName(name string, n int)
}

// tokenAbilityTypes is the Types bitmask every item-token ability reports — a generic
// item, identical across Gold / Silver / Copper. Pre-built so each ability's Types()
// returns a cached value rather than rebuilding the set per call.
var tokenAbilityTypes = card.NewTypeSet(card.TypeGeneric, card.TypeItem)
