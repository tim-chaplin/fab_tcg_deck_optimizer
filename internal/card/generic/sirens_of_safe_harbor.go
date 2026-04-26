// Sirens of Safe Harbor — Generic Action - Attack. Cost 2. Printed power: Red 6, Yellow 5, Blue 4.
// Printed pitch variants: Red 1, Yellow 2, Blue 3. Defense 2.
//
// Text: "When this is put into your graveyard from anywhere, gain 1{h}."
//
// Modelling: the card goes to graveyard after it resolves as an attack, so the 1{h} gain
// fires on every Play — credited as +1 damage equivalent (health is valued 1-to-1 with
// damage). Pitched copies go to the bottom of the deck, not the graveyard, so they don't
// trigger the rider; pitched contributions stay at the printed pitch value only.

package generic

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var sirensOfSafeHarborTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

type SirensOfSafeHarborRed struct{}

func (SirensOfSafeHarborRed) ID() card.ID                 { return card.SirensOfSafeHarborRed }
func (SirensOfSafeHarborRed) Name() string                { return "Sirens of Safe Harbor" }
func (SirensOfSafeHarborRed) Cost(*card.TurnState) int                   { return 2 }
func (SirensOfSafeHarborRed) Pitch() int                  { return 1 }
func (SirensOfSafeHarborRed) Attack() int                 { return 6 }
func (SirensOfSafeHarborRed) Defense() int                { return 2 }
func (SirensOfSafeHarborRed) Types() card.TypeSet         { return sirensOfSafeHarborTypes }
func (SirensOfSafeHarborRed) GoAgain() bool               { return false }
func (SirensOfSafeHarborRed) NotSilverAgeLegal()           {}
func (c SirensOfSafeHarborRed) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() + 1 }

type SirensOfSafeHarborYellow struct{}

func (SirensOfSafeHarborYellow) ID() card.ID                 { return card.SirensOfSafeHarborYellow }
func (SirensOfSafeHarborYellow) Name() string                { return "Sirens of Safe Harbor" }
func (SirensOfSafeHarborYellow) Cost(*card.TurnState) int                   { return 2 }
func (SirensOfSafeHarborYellow) Pitch() int                  { return 2 }
func (SirensOfSafeHarborYellow) Attack() int                 { return 5 }
func (SirensOfSafeHarborYellow) Defense() int                { return 2 }
func (SirensOfSafeHarborYellow) Types() card.TypeSet         { return sirensOfSafeHarborTypes }
func (SirensOfSafeHarborYellow) GoAgain() bool               { return false }
func (SirensOfSafeHarborYellow) NotSilverAgeLegal()           {}
func (c SirensOfSafeHarborYellow) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() + 1 }

type SirensOfSafeHarborBlue struct{}

func (SirensOfSafeHarborBlue) ID() card.ID                 { return card.SirensOfSafeHarborBlue }
func (SirensOfSafeHarborBlue) Name() string                { return "Sirens of Safe Harbor" }
func (SirensOfSafeHarborBlue) Cost(*card.TurnState) int                   { return 2 }
func (SirensOfSafeHarborBlue) Pitch() int                  { return 3 }
func (SirensOfSafeHarborBlue) Attack() int                 { return 4 }
func (SirensOfSafeHarborBlue) Defense() int                { return 2 }
func (SirensOfSafeHarborBlue) Types() card.TypeSet         { return sirensOfSafeHarborTypes }
func (SirensOfSafeHarborBlue) GoAgain() bool               { return false }
func (SirensOfSafeHarborBlue) NotSilverAgeLegal()           {}
func (c SirensOfSafeHarborBlue) Play(s *card.TurnState, _ *card.CardState) int { return c.Attack() + 1 }
