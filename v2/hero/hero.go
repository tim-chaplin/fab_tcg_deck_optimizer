// Package hero defines the Hero interface a FaB hero card satisfies. Sim and sim's external
// consumers (cmd/fabsim, internal/deckio, internal/fabrary) call methods on hero.Hero
// values; concrete heroes live in internal/heroes.
//
// The OnCardPlayed signature uses package-local GameEngine / Logger interfaces — narrow
// surfaces that *gameengine.GameEngine and *turnlogger.TurnLogger satisfy structurally —
// so v2/hero doesn't depend on v2/card.GameEngine or v2/card.Logger. The card-typed
// arguments and constants (CardType / TypeSet / Card) are foundational identification
// types and remain imported from v2/card.
package hero

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
)

// GameEngine is the narrow engine surface hero abilities consume. Today only Viserai's
// "create a Runechant on Runeblade-after-non-attack" needs anything from the engine; new
// hero abilities extend this interface as required. *gameengine.GameEngine satisfies it
// structurally — concrete heroes type-assert when they need a richer surface (e.g.
// played.Types(ge.(card.GameEngine)) to fold Universal class).
type GameEngine interface {
	NonAttackActionPlayed() bool
	CreateRunechants(int)
}

// Logger is the narrow log surface hero abilities use. *turnlogger.TurnLogger satisfies
// it structurally.
type Logger interface {
	AppendPreTrigger(source, text string, n int)
}

// Hero is the view of a hero callers need beyond what the engine itself uses.
type Hero interface {
	ID() ids.HeroID
	Name() string
	Class() card.CardType
	Types() card.TypeSet
	Intelligence() int
	OnCardPlayed(played card.Card, ge GameEngine, l Logger) int
	Opt(cards []card.Card) (top, bottom []card.Card)
}
