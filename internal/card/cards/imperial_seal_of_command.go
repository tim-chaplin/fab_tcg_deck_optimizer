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
//   - The item isn't self-destroyed on activation: DestroyOpponentArsenal is idempotent
//     per turn, and leaving the item in play keeps card-conservation accounting honest.

package cards

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

func (c ImperialSealOfCommandRed) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	ge.CreateItemWithAbility(c, ImperialSealOfCommandRedAbility{})
}

// ImperialSealOfCommandRedAbility is Imperial Seal's printed activated ability. The Royal
// gate in Play registers the OnHit trigger that fires DestroyOpponentArsenal.
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
	if !ge.HeroHasType(card.TypeRoyal) {
		return
	}
	ge.CreateTrigger(self.Card, triggertype.Hit, imperialSealOnHitDestroyArsenal, card.TypeSet.IsAttack)
}

func imperialSealOnHitDestroyArsenal(ge card.GameEngine, l card.Logger, t card.EphemeralTrigger, _ *card.CardState, _ triggertype.Type) {
	if v := ge.DestroyOpponentArsenal(); v > 0 {
		l.AppendPostTriggerf(t.CardName(), v, "%s destroyed opposing arsenal", t.CardName())
	}
}
