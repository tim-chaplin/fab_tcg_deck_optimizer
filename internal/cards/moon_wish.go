// Moon Wish — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on top of your deck rather than pay Moon Wish's {r}
// cost. If Moon Wish hits, search your deck for a card named Sun Kiss, reveal it, put it into
// your hand, then shuffle your deck."
//
// Card-specific quirks:
//   - Tutor priority is Red > Yellow > Blue — the Red printing heals the most (3{h} vs 2 vs
//     1), so the highest-power variant present wins.
//   - The on-hit Sun Kiss tutor wants the synergy ("if you've played Moon Wish") to fire
//     when Sun Kiss resolves immediately (go-again branch), but Moon Wish hasn't been
//     appended to CardsPlayed yet. Play does a transient pre-append + pop around the Sun
//     Kiss invocation so Sun Kiss sees Moon Wish in CardsPlayed without double-adding.
//   - The printed "shuffle your deck" is dropped: deck order isn't modelled beyond removal.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// moonWishPrintedCost is the un-discounted resource cost (also the VariableCost MaxCost bound).
const moonWishPrintedCost = 2

// moonWishCost returns 0 when there's any card left in hand to spend on the alt cost,
// else the printed cost. Shared across all three pitch variants since the alt cost is
// identical.
func moonWishCost(s *sim.TurnState) int {
	if s != nil && len(s.Hand()) > 0 {
		return 0
	}
	return moonWishPrintedCost
}

// moonWishPlay pays the alt cost (when a hand card is available), emits the chain step,
// and registers an OnHit that tutors Sun Kiss. Tutored Sun Kiss plays immediately when
// self has go-again granted; otherwise it lands in hand for next turn.
func moonWishPlay(c sim.Card, s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	name := c.DisplayName()
	// Alt cost: pop a hand card and prepend it to the deck (PrependToDeck flips cacheable).
	var returned sim.Card
	if len(s.Hand()) > 0 {
		returned = s.PopHandAt(0)
		s.PrependToDeck(returned)
	}

	// Emit Moon Wish's chain step first so the alt-cost / tutor lines follow it in order.
	n := self.DealEffectiveAttack(s)
	self.Log(l, n)

	if returned != nil {
		l.AppendPostTriggerf(name, 0, "%s returned %s to top of deck", name, returned.DisplayName())
	}

	self.RegisterOnHit(moonWishOnHit)
}

// moonWishOnHit fires the printed "If this hits, search for Sun Kiss" rider. Top-level so
// registration stays alloc-free; reads the Moon Wish printing off self.Card so we don't
// need a captured copy or a Source-field detour (self IS the Moon Wish that registered
// the handler).
func moonWishOnHit(s *sim.TurnState, l sim.Logger, self *sim.CardState, _ *sim.OnHitHandler) {
	c := self.Card
	name := c.DisplayName()
	sk, ok := s.TutorFromDeck(sunKissTutorPriority)
	if !ok {
		l.AppendPostTriggerf(name, 0, "%s found no Sun Kiss to tutor", name)
		return
	}

	if !self.EffectiveGoAgain() {
		// Tutor lands the card in hand for next turn.
		s.AppendHand(sk)
		l.AppendPostTriggerf(name, 0, "%s tutored %s", name, sk.DisplayName())
		return
	}
	// Go-again: Sun Kiss plays immediately. Pre-append Moon Wish to CardsPlayed so Sun
	// Kiss's "if you've played Moon Wish" synergy fires; pop after so the sim's normal
	// post-Play append doesn't double-add.
	l.AppendPostTriggerf(name, 0, "%s tutored %s and played it", name, sk.DisplayName())
	s.CardsPlayed = append(s.CardsPlayed, c)
	skSelf := &sim.CardState{Card: sk}
	sk.Play(s, l, skSelf)
	s.CardsPlayed = s.CardsPlayed[:len(s.CardsPlayed)-1]
	s.AddToGraveyard(sk)
}

// sunKissTutorPriority picks the highest-priority Sun Kiss printing in the deck. Red >
// Yellow > Blue (Red heals the most: {3,2,1}{h} by colour).
func sunKissTutorPriority(c sim.Card) int {
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

func (MoonWishRed) Cost(s *sim.TurnState) int { return moonWishCost(s) }
func (MoonWishRed) MinCost() int              { return 0 }
func (MoonWishRed) MaxCost() int              { return moonWishPrintedCost }

func (c MoonWishRed) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	moonWishPlay(c, s, l, self)
}

func (MoonWishYellow) Cost(s *sim.TurnState) int { return moonWishCost(s) }
func (MoonWishYellow) MinCost() int              { return 0 }
func (MoonWishYellow) MaxCost() int              { return moonWishPrintedCost }

func (c MoonWishYellow) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	moonWishPlay(c, s, l, self)
}

func (MoonWishBlue) Cost(s *sim.TurnState) int { return moonWishCost(s) }
func (MoonWishBlue) MinCost() int              { return 0 }
func (MoonWishBlue) MaxCost() int              { return moonWishPrintedCost }

func (c MoonWishBlue) Play(s *sim.TurnState, l sim.Logger, self *sim.CardState) {
	moonWishPlay(c, s, l, self)
}
