// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to target
// hero."
//
// Play scans s.Graveyard for any aura. The first one found is removed and appended to s.Banish;
// Play returns 1 arcane and flips ArcaneDamageDealt. No aura means the banish clause fails and
// Play returns 0 — the printed 3 block still applies via Defense().
//
// Source: github.com/the-fab-cube/flesh-and-blood-cards (card.csv).

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var weepingBattlegroundTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

// weepingBattlegroundPlay banishes the first aura in the graveyard and deals 1 arcane; fizzles
// (returns 0) when no aura is there.
func weepingBattlegroundPlay(s *card.TurnState) int {
	for i, c := range s.Graveyard {
		if !c.Types().Has(card.TypeAura) {
			continue
		}
		s.Banish = append(s.Banish, c)
		s.Graveyard = append(s.Graveyard[:i], s.Graveyard[i+1:]...)
		s.ArcaneDamageDealt = true
		return 1
	}
	return 0
}

type WeepingBattlegroundRed struct{}

func (WeepingBattlegroundRed) ID() card.ID                { return card.WeepingBattlegroundRed }
func (WeepingBattlegroundRed) Name() string               { return "Weeping Battleground (Red)" }
func (WeepingBattlegroundRed) Cost(*card.TurnState) int   { return 0 }
func (WeepingBattlegroundRed) Pitch() int                 { return 1 }
func (WeepingBattlegroundRed) Attack() int                { return 0 }
func (WeepingBattlegroundRed) Defense() int               { return 3 }
func (WeepingBattlegroundRed) Types() card.TypeSet        { return weepingBattlegroundTypes }
func (WeepingBattlegroundRed) GoAgain() bool              { return false }
func (WeepingBattlegroundRed) NoMemo()                    {} // Play reads s.Graveyard
func (WeepingBattlegroundRed) Play(s *card.TurnState) int { return weepingBattlegroundPlay(s) }

type WeepingBattlegroundYellow struct{}

func (WeepingBattlegroundYellow) ID() card.ID                { return card.WeepingBattlegroundYellow }
func (WeepingBattlegroundYellow) Name() string               { return "Weeping Battleground (Yellow)" }
func (WeepingBattlegroundYellow) Cost(*card.TurnState) int   { return 0 }
func (WeepingBattlegroundYellow) Pitch() int                 { return 2 }
func (WeepingBattlegroundYellow) Attack() int                { return 0 }
func (WeepingBattlegroundYellow) Defense() int               { return 3 }
func (WeepingBattlegroundYellow) Types() card.TypeSet        { return weepingBattlegroundTypes }
func (WeepingBattlegroundYellow) GoAgain() bool              { return false }
func (WeepingBattlegroundYellow) NoMemo()                    {}
func (WeepingBattlegroundYellow) Play(s *card.TurnState) int { return weepingBattlegroundPlay(s) }

type WeepingBattlegroundBlue struct{}

func (WeepingBattlegroundBlue) ID() card.ID                { return card.WeepingBattlegroundBlue }
func (WeepingBattlegroundBlue) Name() string               { return "Weeping Battleground (Blue)" }
func (WeepingBattlegroundBlue) Cost(*card.TurnState) int   { return 0 }
func (WeepingBattlegroundBlue) Pitch() int                 { return 3 }
func (WeepingBattlegroundBlue) Attack() int                { return 0 }
func (WeepingBattlegroundBlue) Defense() int               { return 3 }
func (WeepingBattlegroundBlue) Types() card.TypeSet        { return weepingBattlegroundTypes }
func (WeepingBattlegroundBlue) GoAgain() bool              { return false }
func (WeepingBattlegroundBlue) NoMemo()                    {}
func (WeepingBattlegroundBlue) Play(s *card.TurnState) int { return weepingBattlegroundPlay(s) }
