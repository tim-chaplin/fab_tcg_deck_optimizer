package optimizations

import (
	"sync/atomic"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// bareChainStepText is the uncached card.ChainStepText captured before init swaps in the
// memoised version. The slow path delegates to it.
var bareChainStepText = card.ChainStepText

func init() {
	card.ChainStepText = cachedChainStepText
}

// cachedChainStepText is the memoised ChainStepText installed over card.ChainStepText at
// init. The text depends only on (Card.ID, FromArsenal) — DisplayName, types, and the
// verb selection are all static — so results live in a pre-warmed table; the per-Play
// string concat / DisplayName allocation disappears on the hot path.
func cachedChainStepText(self *card.CardState) string {
	idx := chainStepCacheIndex(self.Card.ID(), self.FromArsenal)
	if s := chainStepCache[idx].Load(); s != nil {
		return *s
	}
	return chainStepTextSlow(self, idx)
}

// chainStepCache memoises chainStepText results keyed by (Card.ID, FromArsenal). The two
// rows per card cover the in-hand and from-arsenal verb suffixes. WarmChainStepCache fills
// every entry once at init; sized for the full uint16 ID space so lookups are direct
// bounds-checked array reads.
const chainStepCacheSize = 1 << 17 // 2 entries per ID × 65536 IDs

var chainStepCache [chainStepCacheSize]atomic.Pointer[string]

// chainStepCacheIndex packs (id, fromArsenal) into a single uint32 cache index. Bit 16 is
// the FromArsenal flag, bits 0-15 are the card ID — keeps the in-hand and from-arsenal
// variants in adjacent slots so the hot path is a plain array read.
func chainStepCacheIndex(id ids.CardID, fromArsenal bool) uint32 {
	idx := uint32(id)
	if fromArsenal {
		idx |= 1 << 16
	}
	return idx
}

// chainStepTextSlow computes the chain-step prefix string and stores it in chainStepCache
// at idx. Used both by WarmChainStepCache (to populate) and by cachedChainStepText's
// miss path (to backfill on first sighting of an unregistered card — fakes / test stubs
// created outside the cards registry).
//
// Multiple goroutines computing the same entry race-safely converge on the first writer's
// string — every writer produces the same value, so reads after a race still match spec.
func chainStepTextSlow(self *card.CardState, idx uint32) string {
	out := bareChainStepText(self)
	chainStepCache[idx].Store(&out)
	return out
}

// WarmChainStepCache populates the chain-step text cache for every non-nil card in cards
// by writing both the in-hand ((id, false)) and from-arsenal ((id, true)) entries.
// Idempotent. The registry package's init calls this once with the registry slice so the
// runtime hot path is pure cache reads — fakes/test stubs created without registration
// still work via cachedChainStepText's lazy backfill.
func WarmChainStepCache(cards []card.Card) {
	var self card.CardState
	for _, c := range cards {
		if c == nil {
			continue
		}
		self.Card = c
		self.FromArsenal = false
		chainStepTextSlow(&self, chainStepCacheIndex(c.ID(), false))
		self.FromArsenal = true
		chainStepTextSlow(&self, chainStepCacheIndex(c.ID(), true))
	}
}
