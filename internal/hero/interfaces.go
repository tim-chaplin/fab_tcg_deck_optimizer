package hero

// GameEngine is the narrow engine surface hero abilities consume. Today only Viserai's
// "create a Runechant on Runeblade-after-non-attack" needs anything from the engine; new
// hero abilities extend this interface as required. *gameengine.GameEngine satisfies it
// structurally — concrete heroes type-assert when they need a richer surface (e.g.
// played.Types(ge.(card.GameEngine)) to fold Universal class).
type GameEngine interface {
	NonAttackActionPlayed() bool
	CreateRunechants(int)
}

// Logger is the narrow log surface hero abilities use.
type Logger interface {
	AppendPreTrigger(source, text string, n int)
}
