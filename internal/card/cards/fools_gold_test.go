package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

// Compile-time: Fool's Gold implements card.OnDiscardHook so the engine's Discard verb
// fires its "create a Gold token" rider.
var _ card.OnDiscardHook = FoolsGoldYellow{}
