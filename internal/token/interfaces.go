package token

// GameEngine is the slice of engine surface the aura-token fire closures need beyond
// what card.GameEngine already provides. Runechant calls SetArcaneDamageDealt and
// ClearOpponentMarked, neither of which is card-facing — fire closures type-assert the
// engine to this narrow interface to reach them. *gameengine.GameEngine satisfies it
// structurally.
type GameEngine interface {
	SetArcaneDamageDealt(bool)
	ClearOpponentMarked()
}
