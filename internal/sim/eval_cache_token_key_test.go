package sim

import (
	"testing"

	"github.com/tim-chaplin/fab-deck-optimizer/internal/gameengine"
)

// Tests that every token-aura / token-item kind known to gameengine is reflected in the
// cache key. A new token kind added to gameengine without a matching entry in makeCacheKey
// silently collides states with different counts onto the same key — the trip-wire that
// caused the replayBest "cached solution is infeasible" panic when Quicken shipped before
// the key array was widened.
func TestMakeCacheKey_DistinguishesEveryTokenKind(t *testing.T) {
	baseline := gameengine.GameStateBuilder().Build()
	baseKey, ok := makeCacheKey(nil, nil, baseline)
	if !ok {
		t.Fatalf("baseline state should produce a valid cache key")
	}

	cases := []struct {
		name string
		fn   func(*gameengine.GameEngine)
	}{
		{"Runechant", func(ge *gameengine.GameEngine) { ge.CreateRunechants(1) }},
		{"Ponder", func(ge *gameengine.GameEngine) { ge.CreatePonders(1) }},
		{"Quicken", func(ge *gameengine.GameEngine) { ge.CreateQuicken(1) }},
		{"Gold", func(ge *gameengine.GameEngine) { ge.CreateGold(1) }},
		{"Silver", func(ge *gameengine.GameEngine) { ge.CreateSilver(1) }},
		{"Copper", func(ge *gameengine.GameEngine) { ge.CreateCopper(1) }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			altered := gameengine.GameStateBuilder().Build()
			ge := &gameengine.GameEngine{GameState: altered}
			tc.fn(ge)
			alteredKey, ok := makeCacheKey(nil, nil, altered)
			if !ok {
				t.Fatalf("altered state (+1 %s) should produce a valid cache key", tc.name)
			}
			if alteredKey == baseKey {
				t.Errorf("makeCacheKey doesn't distinguish %s count — two states differing only in this token count produce identical cache keys, which causes replayBest infeasibility when the cached solution depends on the missing count",
					tc.name)
			}
		})
	}
}
