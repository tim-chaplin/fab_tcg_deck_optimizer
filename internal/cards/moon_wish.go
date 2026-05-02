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
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var moonWishTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// moonWishPrintedCost is the un-discounted resource cost (also the VariableCost MaxCost bound).
const moonWishPrintedCost = 2

// moonWishCost returns 0 when there's any card left in hand to spend on the alt cost,
// else the printed cost. Shared across all three pitch variants since the alt cost is
// identical.
func moonWishCost(s *sim.TurnState) int {
	if s != nil && len(s.Hand) > 0 {
		return 0
	}
	return moonWishPrintedCost
}

// moonWishPlay pays the alt cost (when a hand card is available), emits the chain step,
// and registers an OnHit that tutors Sun Kiss. Tutored Sun Kiss plays immediately when
// self has go-again granted; otherwise it lands in hand for next turn.
func moonWishPlay(c sim.Card, s *sim.TurnState, self *sim.CardState) {
	name := sim.DisplayName(c)
	// Alt cost: pop a hand card and prepend it to the deck (PrependToDeck flips cacheable).
	var returned sim.Card
	if len(s.Hand) > 0 {
		returned = s.Hand[0]
		s.Hand = s.Hand[1:]
		s.PrependToDeck(returned)
	}

	// Emit Moon Wish's chain step first so the alt-cost / tutor lines follow it in order.
	n := self.DealEffectiveAttack(s)
	s.Log(self, n)

	if returned != nil {
		s.LogPostTriggerf(name, 0, "%s returned %s to top of deck", name, sim.DisplayName(returned))
	}

	self.OnHit = append(self.OnHit, sim.OnHitHandler{Fire: moonWishOnHit})
}

// moonWishOnHit fires the printed "If this hits, search for Sun Kiss" rider. Top-level so
// registration stays alloc-free; reads the Moon Wish printing off self.Card so we don't
// need a captured copy or a Source-field detour (self IS the Moon Wish that registered
// the handler).
func moonWishOnHit(s *sim.TurnState, self *sim.CardState, _ *sim.OnHitHandler) {
	c := self.Card
	name := sim.DisplayName(c)
	sk, ok := s.TutorFromDeck(sunKissTutorPriority)
	if !ok {
		s.LogPostTriggerf(name, 0, "%s found no Sun Kiss to tutor", name)
		return
	}

	if !self.EffectiveGoAgain() {
		// Tutor lands the card in hand for next turn.
		s.Hand = append(s.Hand, sk)
		s.LogPostTriggerf(name, 0, "%s tutored %s", name, sim.DisplayName(sk))
		return
	}
	// Go-again: Sun Kiss plays immediately. Pre-append Moon Wish to CardsPlayed so Sun
	// Kiss's "if you've played Moon Wish" synergy fires; pop after so the sim's normal
	// post-Play append doesn't double-add.
	s.LogPostTriggerf(name, 0, "%s tutored %s and played it", name, sim.DisplayName(sk))
	s.CardsPlayed = append(s.CardsPlayed, c)
	skSelf := &sim.CardState{Card: sk}
	sk.Play(s, skSelf)
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

type MoonWishRed struct{}

func (MoonWishRed) ID() ids.CardID            { return ids.MoonWishRed }
func (MoonWishRed) Name() string              { return "Moon Wish" }
func (MoonWishRed) Cost(s *sim.TurnState) int { return moonWishCost(s) }
func (MoonWishRed) MinCost() int              { return 0 }
func (MoonWishRed) MaxCost() int              { return moonWishPrintedCost }
func (MoonWishRed) Pitch() int                { return 1 }
func (MoonWishRed) Attack() int               { return 5 }
func (MoonWishRed) Defense() int              { return 2 }
func (MoonWishRed) Types() card.TypeSet       { return moonWishTypes }
func (MoonWishRed) GoAgain() bool             { return false }
func (c MoonWishRed) Play(s *sim.TurnState, self *sim.CardState) {
	moonWishPlay(c, s, self)
}

type MoonWishYellow struct{}

func (MoonWishYellow) ID() ids.CardID            { return ids.MoonWishYellow }
func (MoonWishYellow) Name() string              { return "Moon Wish" }
func (MoonWishYellow) Cost(s *sim.TurnState) int { return moonWishCost(s) }
func (MoonWishYellow) MinCost() int              { return 0 }
func (MoonWishYellow) MaxCost() int              { return moonWishPrintedCost }
func (MoonWishYellow) Pitch() int                { return 2 }
func (MoonWishYellow) Attack() int               { return 4 }
func (MoonWishYellow) Defense() int              { return 2 }
func (MoonWishYellow) Types() card.TypeSet       { return moonWishTypes }
func (MoonWishYellow) GoAgain() bool             { return false }
func (c MoonWishYellow) Play(s *sim.TurnState, self *sim.CardState) {
	moonWishPlay(c, s, self)
}

type MoonWishBlue struct{}

func (MoonWishBlue) ID() ids.CardID            { return ids.MoonWishBlue }
func (MoonWishBlue) Name() string              { return "Moon Wish" }
func (MoonWishBlue) Cost(s *sim.TurnState) int { return moonWishCost(s) }
func (MoonWishBlue) MinCost() int              { return 0 }
func (MoonWishBlue) MaxCost() int              { return moonWishPrintedCost }
func (MoonWishBlue) Pitch() int                { return 3 }
func (MoonWishBlue) Attack() int               { return 3 }
func (MoonWishBlue) Defense() int              { return 2 }
func (MoonWishBlue) Types() card.TypeSet       { return moonWishTypes }
func (MoonWishBlue) GoAgain() bool             { return false }
func (c MoonWishBlue) Play(s *sim.TurnState, self *sim.CardState) {
	moonWishPlay(c, s, self)
}
