// Moon Wish — Generic Action - Attack. Cost 2. Printed power: Red 5, Yellow 4, Blue 3. Printed
// pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "You may put a card from your hand on top of your deck rather than pay Moon Wish's {r}
// cost. If Moon Wish hits, search your deck for a card named Sun Kiss, reveal it, put it into your
// hand, then shuffle your deck."
//
// Simplification: Alternative hand-on-top cost and Sun Kiss search rider aren't modelled.
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var moonWishTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type MoonWishRed struct{}

func (MoonWishRed) ID() card.ID                 { return card.MoonWishRed }
func (MoonWishRed) Name() string                { return "Moon Wish (Red)" }
func (MoonWishRed) Cost(*card.TurnState) int                   { return 2 }
func (MoonWishRed) Pitch() int                  { return 1 }
func (MoonWishRed) Attack() int                 { return 5 }
func (MoonWishRed) Defense() int                { return 2 }
func (MoonWishRed) Types() card.TypeSet         { return moonWishTypes }
func (MoonWishRed) GoAgain() bool               { return false }
func (c MoonWishRed) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type MoonWishYellow struct{}

func (MoonWishYellow) ID() card.ID                 { return card.MoonWishYellow }
func (MoonWishYellow) Name() string                { return "Moon Wish (Yellow)" }
func (MoonWishYellow) Cost(*card.TurnState) int                   { return 2 }
func (MoonWishYellow) Pitch() int                  { return 2 }
func (MoonWishYellow) Attack() int                 { return 4 }
func (MoonWishYellow) Defense() int                { return 2 }
func (MoonWishYellow) Types() card.TypeSet         { return moonWishTypes }
func (MoonWishYellow) GoAgain() bool               { return false }
func (c MoonWishYellow) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }

type MoonWishBlue struct{}

func (MoonWishBlue) ID() card.ID                 { return card.MoonWishBlue }
func (MoonWishBlue) Name() string                { return "Moon Wish (Blue)" }
func (MoonWishBlue) Cost(*card.TurnState) int                   { return 2 }
func (MoonWishBlue) Pitch() int                  { return 3 }
func (MoonWishBlue) Attack() int                 { return 3 }
func (MoonWishBlue) Defense() int                { return 2 }
func (MoonWishBlue) Types() card.TypeSet         { return moonWishTypes }
func (MoonWishBlue) GoAgain() bool               { return false }
func (c MoonWishBlue) Play(s *card.TurnState, _ *card.PlayedCard) int { return c.Attack() }
