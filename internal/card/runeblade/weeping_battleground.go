// Weeping Battleground — Runeblade Defense Reaction. Cost 0, Defense 3.
// Printed pitch variants: Red 1, Yellow 2, Blue 3.
// Text: "You may banish an aura from your graveyard. If you do, deal 1 arcane damage to target
// hero."
//
// Play routes through banishAuraFromGraveyard: if s.Graveyard has an aura, banish it for 1
// arcane and flip ArcaneDamageDealt. No aura means the banish clause fails and Play returns
// 0 — the printed 3 block still applies via Defense().

package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var weepingBattlegroundTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeDefenseReaction)

// weepingBattlegroundPlay emits the chain step then writes the banish-for-arcane rider as
// a sub-line under self when an aura was successfully banished from the graveyard.
func weepingBattlegroundPlay(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
	if n := banishAuraFromGraveyard(s); n > 0 {
		s.LogRiderOnPlay(self, "Banished an aura, dealt 1 arcane damage", n)
	}
}

type WeepingBattlegroundRed struct{}

func (WeepingBattlegroundRed) ID() card.ID              { return card.WeepingBattlegroundRed }
func (WeepingBattlegroundRed) Name() string             { return "Weeping Battleground" }
func (WeepingBattlegroundRed) Cost(*card.TurnState) int { return 0 }
func (WeepingBattlegroundRed) Pitch() int               { return 1 }
func (WeepingBattlegroundRed) Attack() int              { return 0 }
func (WeepingBattlegroundRed) Defense() int             { return 3 }
func (WeepingBattlegroundRed) Types() card.TypeSet      { return weepingBattlegroundTypes }
func (WeepingBattlegroundRed) GoAgain() bool            { return false }
func (WeepingBattlegroundRed) NoMemo()                  {}
func (WeepingBattlegroundRed) Play(s *card.TurnState, self *card.CardState) {
	weepingBattlegroundPlay(s, self)
}

type WeepingBattlegroundYellow struct{}

func (WeepingBattlegroundYellow) ID() card.ID              { return card.WeepingBattlegroundYellow }
func (WeepingBattlegroundYellow) Name() string             { return "Weeping Battleground" }
func (WeepingBattlegroundYellow) Cost(*card.TurnState) int { return 0 }
func (WeepingBattlegroundYellow) Pitch() int               { return 2 }
func (WeepingBattlegroundYellow) Attack() int              { return 0 }
func (WeepingBattlegroundYellow) Defense() int             { return 3 }
func (WeepingBattlegroundYellow) Types() card.TypeSet      { return weepingBattlegroundTypes }
func (WeepingBattlegroundYellow) GoAgain() bool            { return false }
func (WeepingBattlegroundYellow) NoMemo()                  {}
func (WeepingBattlegroundYellow) Play(s *card.TurnState, self *card.CardState) {
	weepingBattlegroundPlay(s, self)
}

type WeepingBattlegroundBlue struct{}

func (WeepingBattlegroundBlue) ID() card.ID              { return card.WeepingBattlegroundBlue }
func (WeepingBattlegroundBlue) Name() string             { return "Weeping Battleground" }
func (WeepingBattlegroundBlue) Cost(*card.TurnState) int { return 0 }
func (WeepingBattlegroundBlue) Pitch() int               { return 3 }
func (WeepingBattlegroundBlue) Attack() int              { return 0 }
func (WeepingBattlegroundBlue) Defense() int             { return 3 }
func (WeepingBattlegroundBlue) Types() card.TypeSet      { return weepingBattlegroundTypes }
func (WeepingBattlegroundBlue) GoAgain() bool            { return false }
func (WeepingBattlegroundBlue) NoMemo()                  {}
func (WeepingBattlegroundBlue) Play(s *card.TurnState, self *card.CardState) {
	weepingBattlegroundPlay(s, self)
}
