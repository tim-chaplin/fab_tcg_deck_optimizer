// Spellblade Strike — Runeblade Action - Attack. Cost 1, Defense 3.
// Printed power: Red 4, Yellow 3, Blue 2.
// Text: "Create a Runechant token."

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

var spellbladeStrikeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)

type SpellbladeStrikeRed struct{}

func (SpellbladeStrikeRed) ID() ids.CardID           { return ids.SpellbladeStrikeRed }
func (SpellbladeStrikeRed) Name() string             { return "Spellblade Strike" }
func (SpellbladeStrikeRed) Cost(*card.TurnState) int { return 1 }
func (SpellbladeStrikeRed) Pitch() int               { return 1 }
func (SpellbladeStrikeRed) Attack() int              { return 4 }
func (SpellbladeStrikeRed) Defense() int             { return 3 }
func (SpellbladeStrikeRed) Types() card.TypeSet      { return spellbladeStrikeTypes }
func (SpellbladeStrikeRed) GoAgain() bool            { return false }
func (SpellbladeStrikeRed) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}

type SpellbladeStrikeYellow struct{}

func (SpellbladeStrikeYellow) ID() ids.CardID           { return ids.SpellbladeStrikeYellow }
func (SpellbladeStrikeYellow) Name() string             { return "Spellblade Strike" }
func (SpellbladeStrikeYellow) Cost(*card.TurnState) int { return 1 }
func (SpellbladeStrikeYellow) Pitch() int               { return 2 }
func (SpellbladeStrikeYellow) Attack() int              { return 3 }
func (SpellbladeStrikeYellow) Defense() int             { return 3 }
func (SpellbladeStrikeYellow) Types() card.TypeSet      { return spellbladeStrikeTypes }
func (SpellbladeStrikeYellow) GoAgain() bool            { return false }
func (SpellbladeStrikeYellow) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}

type SpellbladeStrikeBlue struct{}

func (SpellbladeStrikeBlue) ID() ids.CardID           { return ids.SpellbladeStrikeBlue }
func (SpellbladeStrikeBlue) Name() string             { return "Spellblade Strike" }
func (SpellbladeStrikeBlue) Cost(*card.TurnState) int { return 1 }
func (SpellbladeStrikeBlue) Pitch() int               { return 3 }
func (SpellbladeStrikeBlue) Attack() int              { return 2 }
func (SpellbladeStrikeBlue) Defense() int             { return 3 }
func (SpellbladeStrikeBlue) Types() card.TypeSet      { return spellbladeStrikeTypes }
func (SpellbladeStrikeBlue) GoAgain() bool            { return false }
func (SpellbladeStrikeBlue) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	s.CreateAndLogRunechantsOnPlay(self, 1)
}
