// Back Alley Breakline — Generic Action - Attack. Cost 1. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "If an activated ability or action card effect puts Back Alley Breakline face up into a
// zone from your deck, gain 1 action point."
//
// Modelled as a card.FaceUpHook: when another effect calls GameEngine.TurnFaceUp on this
// card's CardState (typically the arsenal slot's, or one scanned out of CardsRemaining),
// the engine fires OnFaceUp and we grant +1 action point. The "from your deck" qualifier
// doesn't translate cleanly — we credit the AP on any face-up flip; in practice no
// implemented effect flips face-up except from-deck contexts so the looser gate is fine.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
)

func backAlleyBreaklineOnFaceUp(ge card.GameEngine, l card.Logger, source string) {
	ge.AddActionPoints(1)
	l.AppendPostTrigger(source, "Face-up trigger granted +1 AP", 0)
}

func (c BackAlleyBreaklineRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (c BackAlleyBreaklineRed) OnFaceUp(ge card.GameEngine, l card.Logger) {
	backAlleyBreaklineOnFaceUp(ge, l, c.DisplayName())
}

func (c BackAlleyBreaklineYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (c BackAlleyBreaklineYellow) OnFaceUp(ge card.GameEngine, l card.Logger) {
	backAlleyBreaklineOnFaceUp(ge, l, c.DisplayName())
}

func (c BackAlleyBreaklineBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {}
func (c BackAlleyBreaklineBlue) OnFaceUp(ge card.GameEngine, l card.Logger) {
	backAlleyBreaklineOnFaceUp(ge, l, c.DisplayName())
}
