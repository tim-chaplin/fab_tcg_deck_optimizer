// Tit for Tat — Generic Action. Cost 0, Defense 2. Printed only in Blue (pitch 3).
// Text: "{t} target hero. {u} another target hero. **Go again**"
//
// The "{t} target hero" clause targets the opposing hero and is ignored — opponent-side
// tap modelling lives outside our deck-evaluation scope. The "{u} another target hero"
// clause untaps our hero. Currently a no-op for the implemented heroes (none get tapped
// during their own action phase), but UntapHero keeps the contract right for future heroes.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func (TitForTatBlue) Play(ge card.GameEngine, _ card.Logger, _ *card.CardState) {
	ge.UntapHero()
}
