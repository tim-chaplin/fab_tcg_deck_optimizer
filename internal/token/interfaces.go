package token

// GameEngine is the slice of engine surface the aura-token fire closures need beyond
// what card.GameEngine already provides. Runechant calls RegisterArcaneDamage — not
// card-facing because card code reaches the same bookkeeping via DealArcaneDamage.
// *gameengine.GameEngine satisfies it structurally.
type GameEngine interface {
	RegisterArcaneDamage(n int)
}
