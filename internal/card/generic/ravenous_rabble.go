// Ravenous Rabble — Generic Action - Attack. Cost 0. Printed power: Red 5, Yellow 4, Blue 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this attacks, reveal the top card of your deck. This gets -X{p}, where X is the pitch
// value of the card revealed this way. **Go again**"
//
// Peek s.Deck[0].Pitch() and subtract from base power, floored at 0. Opts out of the memo because
// the return depends on deck composition. If the deck is empty, no card is revealed so there's no
// penalty.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var ravenousRabbleTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

func ravenousRabblePlay(basePower int, s *card.TurnState) int {
	p := basePower
	if len(s.Deck) > 0 {
		p -= s.Deck[0].Pitch()
	}
	if p < 0 {
		p = 0
	}
	return p
}

type RavenousRabbleRed struct{}

func (RavenousRabbleRed) ID() card.ID                 { return card.RavenousRabbleRed }
func (RavenousRabbleRed) Name() string                { return "Ravenous Rabble (Red)" }
func (RavenousRabbleRed) Cost(*card.TurnState) int                   { return 0 }
func (RavenousRabbleRed) Pitch() int                  { return 1 }
func (RavenousRabbleRed) Attack() int                 { return 5 }
func (RavenousRabbleRed) Defense() int                { return 2 }
func (RavenousRabbleRed) Types() card.TypeSet         { return ravenousRabbleTypes }
func (RavenousRabbleRed) GoAgain() bool               { return true }
func (RavenousRabbleRed) NoMemo()                     {}
func (c RavenousRabbleRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return ravenousRabblePlay(c.Attack(), s) }

type RavenousRabbleYellow struct{}

func (RavenousRabbleYellow) ID() card.ID                 { return card.RavenousRabbleYellow }
func (RavenousRabbleYellow) Name() string                { return "Ravenous Rabble (Yellow)" }
func (RavenousRabbleYellow) Cost(*card.TurnState) int                   { return 0 }
func (RavenousRabbleYellow) Pitch() int                  { return 2 }
func (RavenousRabbleYellow) Attack() int                 { return 4 }
func (RavenousRabbleYellow) Defense() int                { return 2 }
func (RavenousRabbleYellow) Types() card.TypeSet         { return ravenousRabbleTypes }
func (RavenousRabbleYellow) GoAgain() bool               { return true }
func (RavenousRabbleYellow) NoMemo()                     {}
func (c RavenousRabbleYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return ravenousRabblePlay(c.Attack(), s) }

type RavenousRabbleBlue struct{}

func (RavenousRabbleBlue) ID() card.ID                 { return card.RavenousRabbleBlue }
func (RavenousRabbleBlue) Name() string                { return "Ravenous Rabble (Blue)" }
func (RavenousRabbleBlue) Cost(*card.TurnState) int                   { return 0 }
func (RavenousRabbleBlue) Pitch() int                  { return 3 }
func (RavenousRabbleBlue) Attack() int                 { return 3 }
func (RavenousRabbleBlue) Defense() int                { return 2 }
func (RavenousRabbleBlue) Types() card.TypeSet         { return ravenousRabbleTypes }
func (RavenousRabbleBlue) GoAgain() bool               { return true }
func (RavenousRabbleBlue) NoMemo()                     {}
func (c RavenousRabbleBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return ravenousRabblePlay(c.Attack(), s) }
