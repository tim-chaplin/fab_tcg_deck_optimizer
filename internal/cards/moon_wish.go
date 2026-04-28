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

// moonWishPrintedCost is the un-discounted resource cost — also the MaxCost bound for the
// VariableCost solver pre-screens.
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

// moonWishPlay applies the alt cost mutation (when a hand card is available), emits Moon
// Wish's chain step, and on a likely hit tutors a Sun Kiss from the deck. Sun Kiss plays
// immediately when self has go-again granted by a prior chain card (its own chain step
// emits with its printed heal as a post-Moon Wish entry); otherwise it carries to the
// next hand via s.Hand. Alt-cost fires get a "returned X to top of deck" line; tutor
// outcomes get a "tutored Sun Kiss" / "tutored Sun Kiss and played it" / "found no Sun
// Kiss to tutor" post-trigger child line beneath Moon Wish's chain step.
func moonWishPlay(c sim.Card, s *sim.TurnState, self *sim.CardState) {
	name := sim.DisplayName(c)
	// Alt cost: pop a hand card and prepend to deck. Same-turn deck-top readers (the Sun
	// Kiss tutor's post-resolution DrawOne) see it; the next turn's deal sees it too via
	// the sim's end-of-turn deck snapshot. Routes the prepend through PrependToDeck so the
	// cacheable bit flips — modifying deck order makes the chain depend on hidden state
	// (and downstream cards reading the new top will themselves flip).
	var returned sim.Card
	if len(s.Hand) > 0 {
		returned = s.Hand[0]
		s.Hand = s.Hand[1:]
		s.PrependToDeck(returned)
	}

	// Moon Wish's own chain step lands first so subsequent alt-cost / tutor lines + Sun
	// Kiss's chain step (when go-again fires) follow it in order.
	s.ApplyAndLogEffectiveAttack(self)

	if returned != nil {
		s.AddPostTriggerLogEntry(name+" returned "+sim.DisplayName(returned)+" to top of deck", name, 0)
	}

	if !sim.LikelyToHit(self) {
		return
	}
	sk, ok := s.TutorFromDeck(sunKissTutorPriority)
	if !ok {
		s.AddPostTriggerLogEntry(name+" found no Sun Kiss to tutor", name, 0)
		return
	}

	if !self.EffectiveGoAgain() {
		// Tutor lands the card in hand; carries to next turn via the sim's end-of-turn
		// copy of s.Hand.
		s.Hand = append(s.Hand, sk)
		s.AddPostTriggerLogEntry(name+" tutored "+sim.DisplayName(sk), name, 0)
		return
	}
	// Go-again means Moon Wish gets a chain extension this turn. Pre-append Moon Wish to
	// CardsPlayed so Sun Kiss's "if you've played Moon Wish" synergy fires; pop after so
	// the sim's normal post-Play append doesn't double-add. Sun Kiss authors its own
	// chain step inside its Play call — it appears as a separate top-level entry following
	// Moon Wish's tutor narration.
	s.AddPostTriggerLogEntry(name+" tutored "+sim.DisplayName(sk)+" and played it", name, 0)
	s.CardsPlayed = append(s.CardsPlayed, c)
	skSelf := &sim.CardState{Card: sk}
	sk.Play(s, skSelf)
	s.CardsPlayed = s.CardsPlayed[:len(s.CardsPlayed)-1]
	s.AddToGraveyard(sk)
}

// sunKissTutorPriority is the score function passed to TurnState.TutorFromDeck — picks the
// highest-priority Sun Kiss printing present in the deck. Priority order is Red > Yellow >
// Blue: Red heals the most ({3,2,1}{h} by colour), so the highest-power variant present
// wins. A score of 0 (non-Sun-Kiss card) tells TutorFromDeck to skip the entry.
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
