// Moon Wish — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on top of your deck rather than pay Moon Wish's {r}
// cost. If Moon Wish hits, search your deck for a card named Sun Kiss, reveal it, put it into
// your hand, then shuffle your deck."
//
// Card-specific quirks:
//   - Tutor priority is Red > Yellow > Blue — the Red printing heals the most (3{h} vs 2 vs 1).
//   - On the go-again branch, Play transiently appends Moon Wish to CardsPlayed before
//     invoking Sun Kiss so Sun Kiss's "if you've played Moon Wish" synergy fires; pops after.
//   - The printed "shuffle your deck" is dropped: deck order isn't modelled beyond removal.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// moonWishPrintedCost is the un-discounted resource cost (also the VariableCost MaxCost bound).
const moonWishPrintedCost = 2

// moonWishCost returns 0 when there's any card left in hand to spend on the alt cost,
// else the printed cost. Shared across all three pitch variants since the alt cost is
// identical.
func moonWishCost(ge card.GameEngine) int {
	if ge != nil && len(ge.HeldHand()) > 0 {
		return 0
	}
	return moonWishPrintedCost
}

// moonWishPlay pays the alt cost (when a Held hand card is available) and registers an
// OnHit that tutors Sun Kiss. Tutored Sun Kiss plays immediately when self has go-again
// granted; otherwise it lands in hand for next turn.
func moonWishPlay(c card.Card, ge card.GameEngine, l card.Logger, self *card.CardState) {
	name := c.DisplayName()
	// Alt cost: pop a Held hand card and prepend it to the deck.
	var returned card.Card
	if len(ge.HeldHand()) > 0 {
		returned = ge.PopHandAt(0)
		ge.PrependToDeck(returned)
	}
	if returned != nil {
		l.AppendPostTriggerf(name, 0, "%s returned %s to top of deck", name, returned.DisplayName())
	}
	self.RegisterOnHit(moonWishOnHit)
}

// moonWishOnHit fires the printed "If this hits, search for Sun Kiss" rider. Reads the
// Moon Wish printing off self.Card — self IS the Moon Wish that registered the handler.
func moonWishOnHit(ge card.GameEngine, l card.Logger, self *card.CardState, _ *card.OnHitHandler) {
	c := self.Card
	name := c.DisplayName()
	sk, ok := ge.TutorFromDeck(sunKissTutorPriority)
	if !ok {
		l.AppendPostTriggerf(name, 0, "%s found no Sun Kiss to tutor", name)
		return
	}

	if !self.EffectiveGoAgain(ge) {
		// Tutor lands the card in hand for next turn.
		ge.AppendHand(sk)
		l.AppendPostTriggerf(name, 0, "%s tutored %s", name, sk.DisplayName())
		return
	}
	// Go-again: Sun Kiss plays immediately. Pre-append Moon Wish to CardsPlayed so Sun
	// Kiss's "if you've played Moon Wish" synergy fires; pop after so the sim's normal
	// post-Play append doesn't double-add.
	l.AppendPostTriggerf(name, 0, "%s tutored %s and played it", name, sk.DisplayName())
	ge.SetCardsPlayed(append(ge.CardsPlayed(), c))
	skSelf := &card.CardState{Card: sk}
	ge.PlayCard(l, skSelf)
	ge.SetCardsPlayed(ge.CardsPlayed()[:len(ge.CardsPlayed())-1])
	ge.AddToGraveyard(sk)
}

// sunKissTutorPriority picks the highest-priority Sun Kiss printing in the deck. Red >
// Yellow > Blue (Red heals the most: {3,2,1}{h} by colour).
func sunKissTutorPriority(c card.Card) int {
	switch c.ID() {
	case ids.SunKissRed:
		return 3
	case ids.SunKissYellow:
		return 2
	case ids.SunKissBlue:
		return 1
	default:
		return 0
	}
}

func (MoonWishRed) Cost(ge card.GameEngine) int { return moonWishCost(ge) }
func (MoonWishRed) MinCost() int                { return 0 }
func (MoonWishRed) MaxCost() int                { return moonWishPrintedCost }

func (c MoonWishRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moonWishPlay(c, ge, l, self)
}

func (MoonWishYellow) Cost(ge card.GameEngine) int { return moonWishCost(ge) }
func (MoonWishYellow) MinCost() int                { return 0 }
func (MoonWishYellow) MaxCost() int                { return moonWishPrintedCost }

func (c MoonWishYellow) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moonWishPlay(c, ge, l, self)
}

func (MoonWishBlue) Cost(ge card.GameEngine) int { return moonWishCost(ge) }
func (MoonWishBlue) MinCost() int                { return 0 }
func (MoonWishBlue) MaxCost() int                { return moonWishPrintedCost }

func (c MoonWishBlue) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	moonWishPlay(c, ge, l, self)
}
