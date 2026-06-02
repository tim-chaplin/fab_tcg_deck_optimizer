// Imperial Seal of Command — Generic Action - Item. Cost 0. Pitch 1 (Red). Legendary.
//
// Text: "**Legendary** **Action** - Destroy this: Defense reaction cards can't be played this
// turn. If you are Royal, the next time you hit a hero this turn, destroy all cards in their
// arsenal. **Go again**"
//
// Modelling:
//   - Play creates the in-play item; the printed activated ability lives on
//     ImperialSealOfCommandRedAbility, enumerated by the wmask as a 1-AP playable. The
//     ability carries Go again, the Royal gate, and the OnHit-destroy-arsenal trigger.
//   - The "Defense reaction cards can't be played this turn" rider is dropped — the sim's
//     defense pass doesn't model an opt-out window.
//   - "Destroy this" is the activation cost, so the ability destroys its own in-play item
//     (source to graveyard) on every activation, Royal or not, via self.Item.Destroy().

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func (c ImperialSealOfCommandRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItemWithAbility(c, ImperialSealOfCommandRedAbility{})
}

// ImperialSealOfCommandRedAbility is Imperial Seal's printed activated ability. Play pays the
// "Destroy this" cost by destroying the item, then (when Royal) registers the OnHit trigger that
// fires DestroyOpponentArsenal.
type ImperialSealOfCommandRedAbility struct{}

func (ImperialSealOfCommandRedAbility) ID() ids.CardID      { return ids.ImperialSealOfCommandRedAbilityID }
func (ImperialSealOfCommandRedAbility) Name() string        { return "Imperial Seal of Command" }
func (ImperialSealOfCommandRedAbility) DisplayName() string { return "Imperial Seal of Command" }
func (ImperialSealOfCommandRedAbility) Cost() int           { return 0 }
func (ImperialSealOfCommandRedAbility) Pitch() int          { return 0 }
func (ImperialSealOfCommandRedAbility) Attack() int         { return 0 }
func (ImperialSealOfCommandRedAbility) Defense() int        { return 0 }
func (ImperialSealOfCommandRedAbility) Types(card.GameEngine) card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeItem)
}
func (ImperialSealOfCommandRedAbility) GoAgain(card.GameEngine) bool { return true }

func (ImperialSealOfCommandRedAbility) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	self.Item.Destroy(true) // pay "Destroy this" before the Royal early-return
	if !ge.HeroHasType(card.TypeRoyal) {
		return
	}
	ge.CreateTrigger(self.Card, triggertype.Hit, imperialSealOnHitDestroyArsenal, card.TypeSet.IsAttack)
}

func imperialSealOnHitDestroyArsenal(ge card.GameEngine, l card.Logger, t card.EphemeralTrigger, _ card.FireContext) {
	if v := ge.DestroyOpponentArsenal(); v > 0 {
		l.AppendPostTriggerf(t.CardName(), v, "%s destroyed opposing arsenal", t.CardName())
	}
}
