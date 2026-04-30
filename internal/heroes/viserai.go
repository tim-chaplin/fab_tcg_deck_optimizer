// Viserai — Runeblade Hero, Young. Health 20, Intelligence 4.
// Text: "Whenever you play a Runeblade card, if you have played another 'non-attack' action card
// this turn, create a Runechant token."

package heroes

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

var viseraiTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeHero, card.TypeYoung)

// Viserai is Young Viserai.
type Viserai struct{}

func (Viserai) ID() ids.HeroID      { return ids.ViseraiID }
func (Viserai) Name() string        { return "Viserai" }
func (Viserai) Health() int         { return 20 }
func (Viserai) Intelligence() int   { return 4 }
func (Viserai) Types() card.TypeSet { return viseraiTypes }

// OnCardPlayed implements Viserai's hero ability: whenever a Runeblade card is played, if a
// non-attack action (Action without Attack) has been played this turn, create a Runechant
// token.
func (Viserai) OnCardPlayed(played sim.Card, s *sim.TurnState) int {
	t := played.Types()
	// Weapon swings aren't "playing a card" and don't trigger Viserai.
	if !t.Has(card.TypeRuneblade) || t.Has(card.TypeWeapon) {
		return 0
	}
	if s.NonAttackActionPlayed {
		return s.CreateAndLogRunechants("Viserai", sim.DisplayName(played), 1)
	}
	return 0
}

// Opt is the Viserai-specific Opt heuristic: keep one card per "slot category" and
// bottom the redundant rest, since a balanced hand is what feeds Viserai's runechant
// trigger. Slots:
//
//   - Non-attack enabler: an Action card that isn't an Attack — needed to satisfy "if you
//     have played another non-attack action card this turn" before the next Runeblade
//     attack drops a runechant.
//   - Action without Go again: an Action card that doesn't extend the chain on its own —
//     one is enough to close out a chain; further copies just sit in hand.
//   - Block-only defender: a card whose only role is defending — Defense Reaction or
//     Block subtype. Most cards carry a non-zero printed Defense value as a secondary
//     option, so Defense > 0 alone is too broad — we only count cards that are
//     defenders first and foremost. One block per turn covers the usual incoming-damage
//     budget.
//   - Blue pitch: any card with Pitch == 3 — one fully funds a 3-cost play; redundant
//     blues stack resources we won't spend.
//
// A card belongs to zero or more slots. It's kept on top when at least one of its
// slots is still uncovered (i.e., this is the first card we've seen for that slot);
// otherwise every slot it provides is already covered and we bottom it. Cards that
// belong to no slot at all stay on top — Viserai has no signal that the next hand
// would prefer fewer of them.
//
// Opt(1) always tops the only revealed card: with one input the slot tracker starts
// empty, so any slot the card provides is uncovered.
func (Viserai) Opt(cards []sim.Card) (top, bottom []sim.Card) {
	var covered viseraiOptSlots
	top = make([]sim.Card, 0, len(cards))
	for _, c := range cards {
		slots := viseraiSlotsFor(c)
		// Cards with no slot membership stay on top; we have no redundancy signal.
		keep := slots.empty() || slots.coversNew(covered)
		if keep {
			top = append(top, c)
			covered = covered.union(slots)
		} else {
			bottom = append(bottom, c)
		}
	}
	return top, bottom
}

// viseraiOptSlots is the bitfield of slot categories Viserai's Opt heuristic tracks.
// One bool per slot keeps the helper readable; the small handful of slots doesn't
// justify a packed bitmask.
type viseraiOptSlots struct {
	nonAttackEnabler bool
	nonGoAgainAction bool
	defender         bool
	bluePitch        bool
}

// empty reports whether s sits in no slot at all.
func (s viseraiOptSlots) empty() bool {
	return !s.nonAttackEnabler && !s.nonGoAgainAction && !s.defender && !s.bluePitch
}

// coversNew reports whether s provides at least one slot that's still uncovered in
// covered (i.e., the new card would fill at least one fresh slot).
func (s viseraiOptSlots) coversNew(covered viseraiOptSlots) bool {
	return (s.nonAttackEnabler && !covered.nonAttackEnabler) ||
		(s.nonGoAgainAction && !covered.nonGoAgainAction) ||
		(s.defender && !covered.defender) ||
		(s.bluePitch && !covered.bluePitch)
}

// union returns the OR of two slot sets — used when we keep a card to mark every slot
// it provides as covered.
func (s viseraiOptSlots) union(other viseraiOptSlots) viseraiOptSlots {
	return viseraiOptSlots{
		nonAttackEnabler: s.nonAttackEnabler || other.nonAttackEnabler,
		nonGoAgainAction: s.nonGoAgainAction || other.nonGoAgainAction,
		defender:         s.defender || other.defender,
		bluePitch:        s.bluePitch || other.bluePitch,
	}
}

// viseraiSlotsFor classifies c into Viserai's Opt-heuristic slots.
func viseraiSlotsFor(c sim.Card) viseraiOptSlots {
	t := c.Types()
	return viseraiOptSlots{
		nonAttackEnabler: t.IsNonAttackAction(),
		nonGoAgainAction: t.Has(card.TypeAction) && !c.GoAgain(),
		defender:         t.IsDefenseReaction() || t.Has(card.TypeBlock),
		bluePitch:        c.Pitch() == 3,
	}
}
