// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to target
// hero."
//
// Play routes through banishAuraFromGraveyard: if s.Graveyard has an aura, banish it for 1
// arcane and flip ArcaneDamageDealt. No aura means the banish clause fails and Play returns
// 0 — the printed 3 block still applies via Defense().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var weepingBattlegroundTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

// weepingBattlegroundPlay emits the chain step then writes the banish-for-arcane rider as
// a sub-line under self when an aura was successfully banished from the graveyard.
func weepingBattlegroundPlay(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
	if n := banishAuraFromGraveyard(s); n > 0 {
		s.ApplyAndLogRiderOnPlay(self, n, "Banished an aura, dealt 1 arcane damage")
	}
}

type WeepingBattlegroundRed struct{}

func (WeepingBattlegroundRed) ID() ids.CardID          { return ids.WeepingBattlegroundRed }
func (WeepingBattlegroundRed) Name() string            { return "Weeping Battleground" }
func (WeepingBattlegroundRed) Cost(*sim.TurnState) int { return 0 }
func (WeepingBattlegroundRed) Pitch() int              { return 1 }
func (WeepingBattlegroundRed) Attack() int             { return 0 }
func (WeepingBattlegroundRed) Defense() int            { return 3 }
func (WeepingBattlegroundRed) Types() card.TypeSet     { return weepingBattlegroundTypes }
func (WeepingBattlegroundRed) GoAgain() bool           { return false }
func (WeepingBattlegroundRed) Play(s *sim.TurnState, self *sim.CardState) {
	weepingBattlegroundPlay(s, self)
}

type WeepingBattlegroundYellow struct{}

func (WeepingBattlegroundYellow) ID() ids.CardID          { return ids.WeepingBattlegroundYellow }
func (WeepingBattlegroundYellow) Name() string            { return "Weeping Battleground" }
func (WeepingBattlegroundYellow) Cost(*sim.TurnState) int { return 0 }
func (WeepingBattlegroundYellow) Pitch() int              { return 2 }
func (WeepingBattlegroundYellow) Attack() int             { return 0 }
func (WeepingBattlegroundYellow) Defense() int            { return 3 }
func (WeepingBattlegroundYellow) Types() card.TypeSet     { return weepingBattlegroundTypes }
func (WeepingBattlegroundYellow) GoAgain() bool           { return false }
func (WeepingBattlegroundYellow) Play(s *sim.TurnState, self *sim.CardState) {
	weepingBattlegroundPlay(s, self)
}

type WeepingBattlegroundBlue struct{}

func (WeepingBattlegroundBlue) ID() ids.CardID          { return ids.WeepingBattlegroundBlue }
func (WeepingBattlegroundBlue) Name() string            { return "Weeping Battleground" }
func (WeepingBattlegroundBlue) Cost(*sim.TurnState) int { return 0 }
func (WeepingBattlegroundBlue) Pitch() int              { return 3 }
func (WeepingBattlegroundBlue) Attack() int             { return 0 }
func (WeepingBattlegroundBlue) Defense() int            { return 3 }
func (WeepingBattlegroundBlue) Types() card.TypeSet     { return weepingBattlegroundTypes }
func (WeepingBattlegroundBlue) GoAgain() bool           { return false }
func (WeepingBattlegroundBlue) Play(s *sim.TurnState, self *sim.CardState) {
	weepingBattlegroundPlay(s, self)
}
