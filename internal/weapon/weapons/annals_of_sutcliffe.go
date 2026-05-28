// Annals of Sutcliffe — Runeblade Weapon - Book (2H). Activation cost {r}{r}{r}.
//
// Text: "**Once per Turn Action** - {r}{r}{r}: Draw a card. If an attack action card and a
// 'non-attack' action card were pitched this way, create a Runechant token."
//
// The activation is a non-attack action — the ability card carries no TypeAttack, so it
// resolves through the attack-turn runner without dealing damage. The conditional Runechant
// reads self.PitchedToPlay (the cards the runner attributed to funding this activation) and
// fires when that set contains both an attack action card and a non-attack action card.

package weapons

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

var annalsOfSutcliffeTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeBook, card.TypeTwoHand)

// AnnalsOfSutcliffe is the platonic weapon card — the equipped permanent. Its activated
// ability is what the attack-turn runner enqueues each turn.
type AnnalsOfSutcliffe struct{}

func (AnnalsOfSutcliffe) ID() ids.CardID                                     { return ids.AnnalsOfSutcliffeID }
func (AnnalsOfSutcliffe) Name() string                                       { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffe) DisplayName() string                                { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffe) Cost() int                                          { return 0 }
func (AnnalsOfSutcliffe) Pitch() int                                         { return 0 }
func (AnnalsOfSutcliffe) Attack() int                                        { return 0 }
func (AnnalsOfSutcliffe) Defense() int                                       { return 0 }
func (AnnalsOfSutcliffe) Types(card.GameEngine) card.TypeSet                 { return annalsOfSutcliffeTypes }
func (AnnalsOfSutcliffe) GoAgain(card.GameEngine) bool                       { return false }
func (AnnalsOfSutcliffe) Play(card.GameEngine, card.Logger, *card.CardState) {}
func (AnnalsOfSutcliffe) Hands() int                                         { return 2 }
func (AnnalsOfSutcliffe) Ability() card.Card                                 { return annalsOfSutcliffeAbility }

var annalsOfSutcliffeAbility card.Card = AnnalsOfSutcliffeAbility{}

var annalsOfSutcliffeAbilityTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon, card.TypeBook, card.TypeTwoHand)

type AnnalsOfSutcliffeAbility struct{}

func (AnnalsOfSutcliffeAbility) ID() ids.CardID      { return ids.AnnalsOfSutcliffeAbilityID }
func (AnnalsOfSutcliffeAbility) Name() string        { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffeAbility) DisplayName() string { return "Annals of Sutcliffe" }
func (AnnalsOfSutcliffeAbility) Cost() int           { return 3 }
func (AnnalsOfSutcliffeAbility) Pitch() int          { return 0 }
func (AnnalsOfSutcliffeAbility) Attack() int         { return 0 }
func (AnnalsOfSutcliffeAbility) Defense() int        { return 0 }
func (AnnalsOfSutcliffeAbility) Types(card.GameEngine) card.TypeSet {
	return annalsOfSutcliffeAbilityTypes
}
func (AnnalsOfSutcliffeAbility) GoAgain(card.GameEngine) bool { return false }

func (AnnalsOfSutcliffeAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.DrawOne()
	l.AppendPostTrigger(self.Card.DisplayName(), "Drew a card", 0)

	pitchedAttack, pitchedNonAttack := false, false
	for _, p := range self.PitchedToPlay {
		t := p.Card.Types(ge)
		if t.IsAttackAction() {
			pitchedAttack = true
		}
		if t.IsNonAttackAction() {
			pitchedNonAttack = true
		}
	}
	if pitchedAttack && pitchedNonAttack {
		ge.CreateRunechants(1)
		l.AppendPostTrigger(self.Card.DisplayName(), "Created a runechant", 1)
	}
}
