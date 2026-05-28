package turntests

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/hero/heroes"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/testutils"
)

// Pins per-object counter attribution for two equipped copies of the same 1H weapon. Both
// copies share one Ability() card value, so the old GameEngine.WeaponByID(id) lookup
// resolved every swing to the FIRST equipped copy with a matching CardID — so both swings
// would bump the first object (c0=2, c1=0) no matter which copy the runner actually swung.
// With self.Weapon plumbing each swing mutates the object the attack-turn runner attributed
// to it, so each copy gets exactly one counter (c0=1, c1=1).
//
// The swing has Go again, so the turn's single action point funds the first swing and its
// go-again funds the second — both copies swing in one turn. This is the case WeaponByID
// got wrong.
func TestDuplicateWeapon_CounterAttributedToSwungObject(t *testing.T) {
	d := deck.New(heroes.Viserai,
		[]deck.Weapon{testutils.CounterWeapon{}, testutils.CounterWeapon{}}, nil)
	hand := []card.Card{
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
		testutils.FakeBlueResource(),
	}

	summary := sim.EvalOneTurnForTesting(d,
		gameengine.GameStateBuilder().SetIncomingPhysicalDamage(0).Build(), hand)

	weapons := summary.State.Weapons()
	if len(weapons) != 2 {
		t.Fatalf("equipped weapons = %d, want 2\nBestLine: %s",
			len(weapons), formatBestLine(summary.BestLine))
	}
	c0, c1 := weapons[0].Count(), weapons[1].Count()
	// Both swing (go-again chains them), so each object holds exactly one counter. The old
	// WeaponByID path would put both counters on weapons[0] (c0=2, c1=0).
	if c0 != 1 || c1 != 1 {
		t.Errorf("counters = (c0=%d, c1=%d), want (1, 1): each swung object should bump its own counter, not the first by CardID\nBestLine: %s",
			c0, c1, formatBestLine(summary.BestLine))
	}
}
