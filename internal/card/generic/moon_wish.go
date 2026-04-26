// Moon Wish — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on top of your deck rather than pay Moon Wish's {r}
// cost. If Moon Wish hits, search your deck for a card named Sun Kiss, reveal it, put it into
// your hand, then shuffle your deck."
//
// Alt cost is implemented via card.VariableCost: Cost(s) returns 0 when the partition has any
// Held card to spend (len(s.Held) > 0), else the printed 2. When the alt cost fires Play pops
// the first Held card, prepends it to s.Deck (top-of-deck placement), and records it on
// s.HeldConsumed so the deck loop suppresses the BestLine[Held] → nextHeld carry that would
// otherwise double-count.
//
// On-hit Sun Kiss tutor: card.LikelyToHit gates the search. The deck is scanned for any Sun
// Kiss printing and the lowest-pitch variant wins (Red over Yellow over Blue) since pitch 1
// is the strongest deck slot. The chosen card is removed from s.Deck — the printed "shuffle
// your deck" reorder is dropped because the simulator doesn't model deck order beyond
// removal. If Moon Wish has go-again at this point (granted by a prior chain card such as
// Flying High), Sun Kiss plays immediately: damage is added to Moon Wish's return, the
// synergy fires (Sun Kiss reads Moon Wish in CardsPlayed via a transient pre-append), and
// Sun Kiss heads to the graveyard. Without go-again Sun Kiss appends to s.Drawn so it
// carries to the next hand as a tutored card.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var moonWishTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// moonWishPrintedCost is the un-discounted resource cost — also the MaxCost bound for the
// VariableCost solver pre-screens.
const moonWishPrintedCost = 2

// moonWishCost returns 0 when the partition has any Held card to spend, else the printed
// cost. Shared across all three pitch variants since the alt cost is identical.
func moonWishCost(s *card.TurnState) int {
	if s != nil && len(s.Held) > 0 {
		return 0
	}
	return moonWishPrintedCost
}

// moonWishPlay applies the alt cost mutation (when a Held card is available), runs the
// printed attack, and on a likely hit tutors a Sun Kiss from the deck. Sun Kiss plays
// immediately when self has go-again granted by a prior chain card; otherwise it carries to
// the next hand via s.Drawn.
func moonWishPlay(c card.Card, attack int, s *card.TurnState, self *card.CardState) int {
	if len(s.Held) > 0 {
		moved := s.Held[0]
		s.Held = s.Held[1:]
		s.HeldConsumed = append(s.HeldConsumed, moved)
		// Prepend to deck — top-of-deck placement so any same-turn deck-top reader (e.g.
		// the Sun Kiss tutor's post-shuffle DrawOne) sees it.
		newDeck := make([]card.Card, 0, len(s.Deck)+1)
		newDeck = append(newDeck, moved)
		newDeck = append(newDeck, s.Deck...)
		s.Deck = newDeck
	}

	if !card.LikelyToHit(self) {
		return attack
	}
	sk := bestSunKissInDeck(s.Deck)
	if sk == nil {
		return attack
	}
	s.Deck = removeFirstByID(s.Deck, sk.ID())

	if !self.EffectiveGoAgain() {
		s.Drawn = append(s.Drawn, sk)
		return attack
	}
	// Go-again means Moon Wish gets a chain extension this turn. Pre-append Moon Wish to
	// CardsPlayed so Sun Kiss's "if you've played Moon Wish" synergy fires; pop after so
	// the solver's normal post-Play append (in playSequenceWithMeta) doesn't double-add.
	s.CardsPlayed = append(s.CardsPlayed, c)
	skSelf := &card.CardState{Card: sk}
	skDmg := sk.Play(s, skSelf)
	s.CardsPlayed = s.CardsPlayed[:len(s.CardsPlayed)-1]
	s.AddToGraveyard(sk)
	return attack + skDmg
}

// bestSunKissInDeck returns the highest-priority Sun Kiss printing present in deck, or nil
// when no Sun Kiss is in the deck. Priority order is Red > Yellow > Blue: the lower the
// pitch value, the more flexible the variant in any future hand.
func bestSunKissInDeck(deck []card.Card) card.Card {
	var pickedRed, pickedYellow, pickedBlue card.Card
	for _, c := range deck {
		switch c.ID() {
		case card.SunKissRed:
			pickedRed = c
		case card.SunKissYellow:
			if pickedYellow == nil {
				pickedYellow = c
			}
		case card.SunKissBlue:
			if pickedBlue == nil {
				pickedBlue = c
			}
		}
	}
	switch {
	case pickedRed != nil:
		return pickedRed
	case pickedYellow != nil:
		return pickedYellow
	default:
		return pickedBlue
	}
}

// removeFirstByID returns deck with the first occurrence of id removed. The returned slice
// shares no backing storage with deck so subsequent mutations on the returned slice can't
// poison the per-leaf deck reference.
func removeFirstByID(deck []card.Card, id card.ID) []card.Card {
	for i, c := range deck {
		if c.ID() == id {
			out := make([]card.Card, 0, len(deck)-1)
			out = append(out, deck[:i]...)
			out = append(out, deck[i+1:]...)
			return out
		}
	}
	return deck
}

type MoonWishRed struct{}

func (MoonWishRed) ID() card.ID                 { return card.MoonWishRed }
func (MoonWishRed) Name() string                { return "Moon Wish (Red)" }
func (MoonWishRed) Cost(s *card.TurnState) int  { return moonWishCost(s) }
func (MoonWishRed) MinCost() int                { return 0 }
func (MoonWishRed) MaxCost() int                { return moonWishPrintedCost }
func (MoonWishRed) Pitch() int                  { return 1 }
func (MoonWishRed) Attack() int                 { return 5 }
func (MoonWishRed) Defense() int                { return 2 }
func (MoonWishRed) Types() card.TypeSet         { return moonWishTypes }
func (MoonWishRed) GoAgain() bool               { return false }
func (MoonWishRed) NoMemo()                     {} // alt-cost mutates Deck; tutor depends on deck contents.
func (c MoonWishRed) Play(s *card.TurnState, self *card.CardState) int {
	return moonWishPlay(c, c.Attack(), s, self)
}

type MoonWishYellow struct{}

func (MoonWishYellow) ID() card.ID                { return card.MoonWishYellow }
func (MoonWishYellow) Name() string               { return "Moon Wish (Yellow)" }
func (MoonWishYellow) Cost(s *card.TurnState) int { return moonWishCost(s) }
func (MoonWishYellow) MinCost() int               { return 0 }
func (MoonWishYellow) MaxCost() int               { return moonWishPrintedCost }
func (MoonWishYellow) Pitch() int                 { return 2 }
func (MoonWishYellow) Attack() int                { return 4 }
func (MoonWishYellow) Defense() int               { return 2 }
func (MoonWishYellow) Types() card.TypeSet        { return moonWishTypes }
func (MoonWishYellow) GoAgain() bool              { return false }
func (MoonWishYellow) NoMemo()                    {}
func (c MoonWishYellow) Play(s *card.TurnState, self *card.CardState) int {
	return moonWishPlay(c, c.Attack(), s, self)
}

type MoonWishBlue struct{}

func (MoonWishBlue) ID() card.ID                { return card.MoonWishBlue }
func (MoonWishBlue) Name() string               { return "Moon Wish (Blue)" }
func (MoonWishBlue) Cost(s *card.TurnState) int { return moonWishCost(s) }
func (MoonWishBlue) MinCost() int               { return 0 }
func (MoonWishBlue) MaxCost() int               { return moonWishPrintedCost }
func (MoonWishBlue) Pitch() int                 { return 3 }
func (MoonWishBlue) Attack() int                { return 3 }
func (MoonWishBlue) Defense() int               { return 2 }
func (MoonWishBlue) Types() card.TypeSet        { return moonWishTypes }
func (MoonWishBlue) GoAgain() bool              { return false }
func (MoonWishBlue) NoMemo()                    {}
func (c MoonWishBlue) Play(s *card.TurnState, self *card.CardState) int {
	return moonWishPlay(c, c.Attack(), s, self)
}
