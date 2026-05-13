// Runechant token: aura-flavored token created by ge.CreateRunechants. Fires per attack
// (triggertype.Attack) — flips ArcaneDamageDealt when its count clears the
// damage-likely-to-hit window, then destroys. Damage is credited at creation time inside
// CreateRunechants; this handler is pure state cleanup.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/aura"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/triggertype"
)

// runechantTokenName is the canonical display name. The engine matches by CardName when
// bumping an existing entry's Count or reading a count.
const runechantTokenName = "Runechant"

// NewRunechant returns a fresh Runechant token aura at count n.
func NewRunechant(n int) *aura.Aura {
	return aura.NewToken(runechantTokenName, ids.RunechantTokenID, triggertype.Attack, runechantFire, n)
}

// RunechantToken is the marker struct for FaB's runechant token-aura — exists for
// naming-convention symmetry with the item tokens (GoldToken / SilverToken / CopperToken).
// Runechants don't print an activated-ability card; the behaviour lives in runechantFire.
type RunechantToken struct{}

// runechantEngine is the slice of engine surface runechantFire reaches for —
// *gameengine.GameEngine satisfies it structurally.
type runechantEngine interface {
	SetArcaneDamageDealt(bool)
}

// runechantFire is the triggertype.Attack handler shared by every Runechant aura. Fires
// before each attack / weapon swing.
func runechantFire(engine, _ any, a aura.Ctx) {
	if likelyDamageHits(a.Count(), false) {
		engine.(runechantEngine).SetArcaneDamageDealt(true)
	}
	a.Destroy(false)
}

// likelyDamageHits is the runechant-specific copy of the engine's damage-hits heuristic —
// kept local so this file doesn't reach back into gameengine for it. A typical FaB card is
// worth ~3 points; blocking 1/4/7 with a card over-pays, so those amounts slip past.
// Dominate caps the defender to one block, so 5+ damage always lands at least 2.
func likelyDamageHits(n int, dominate bool) bool {
	if dominate && n >= 5 {
		return true
	}
	return n == 1 || n == 4 || n == 7
}
