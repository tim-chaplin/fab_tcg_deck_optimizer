// Viserai — Runeblade Hero, Young. Health 20, Intelligence 4.
// Text: "Whenever you play a Runeblade card, if you have played another 'non-attack' action card
// this turn, create a Runechant token."

package heroes

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/trigger"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/triggertype"
)

var viseraiTypes = card.NewTypeSet(card.TypeRuneblade, card.TypeHero, card.TypeYoung)

// viserai embeds trigger.Trigger[card.Hero] so Viserai's Runechant ability is
// dispatched through the engine's normal trigger walk.
type viserai struct {
	trigger.Trigger[card.Hero]
}

// Viserai is Young Viserai.
var Viserai = &viserai{
	Trigger: trigger.FromHero[card.Hero](
		triggertype.CardOrAbility,
		viseraiOnCardPlayed,
		false,
		viseraiTypeFilter,
	),
}

func (*viserai) ID() ids.HeroID       { return ids.ViseraiID }
func (*viserai) Name() string         { return "Viserai" }
func (*viserai) Health() int          { return 20 }
func (*viserai) Intelligence() int    { return 4 }
func (*viserai) Types() card.TypeSet  { return viseraiTypes }
func (*viserai) Class() card.CardType { return card.TypeRuneblade }

func (v *viserai) Fire(engine card.GameEngine, logger card.Logger, ctx card.FireContext) {
	v.Invoke(engine, logger, v, ctx)
}

// viseraiTypeFilter narrows the trigger's firing site to Runeblade-typed cards that
// aren't weapons — equipping or swinging a weapon isn't "playing a card" so a
// Runeblade weapon swing must not credit a Runechant.
func viseraiTypeFilter(t card.TypeSet) bool {
	return t.Has(card.TypeRuneblade) && !t.Has(card.TypeWeapon)
}

// viseraiOnCardPlayed implements Viserai's "Whenever you play a Runeblade card, if you
// have played another non-attack action card this turn, create a Runechant" trigger.
// The type filter above already gates on "Runeblade card, not weapon"; this handler
// adds the non-attack-action precondition and credits the Runechant via the engine.
func viseraiOnCardPlayed(ge card.GameEngine, l card.Logger, _ card.Hero, ctx card.FireContext) {
	if !ge.NonAttackActionPlayed() {
		return
	}
	ge.CreateRunechants(1)
	l.AppendPreTrigger(ctx.TriggeringCard.Card.DisplayName(), "Viserai created a runechant", 1)
}

// Opt is the Viserai-specific Opt heuristic: keep one card per "slot category" and
// bottom anything that would over-fill a slot already covered by an earlier card,
// since a balanced hand is what feeds Viserai's runechant trigger. Slots:
//
//   - Non-attack enabler: an Action card that isn't an Attack — needed to satisfy "if you
//     have played another non-attack action card this turn" before the next Runeblade
//     attack drops a runechant.
//   - Action without Go again: an Action card that doesn't extend the attack turn — one is
//     enough to close out an attack turn; further copies just sit in hand. Uses printed
//     GoAgain() only.
//   - Block-only defender: a card whose only role is defending — Defense Reaction or
//     Block subtype. Most cards carry a non-zero printed Defense value as a secondary
//     option, so Defense > 0 alone is too broad — we only count cards that are
//     defenders first and foremost. One block per turn covers the usual incoming-damage
//     budget.
//   - Blue pitch: any card with Pitch == 3 — one fully funds a 3-cost play; redundant
//     blues stack resources we won't spend.
//
// A card belongs to zero or more slots. It's bottomed when ANY of its slots is already
// covered by an earlier card kept on top — over-filling a slot wastes hand space even
// if the card would also fill some other fresh slot. Cards in zero slots are always
// kept (Runerager Swarm, generic 0-cost go-again attacks, etc.) — Viserai has no
// "balanced hand" signal for them, so multiples are fine.
//
// Opt(1) always tops the only revealed card: with one input the slot tracker starts
// empty, so no slot the card might provide can already be covered.
func (*viserai) Opt(cards []card.Card) (top, bottom []card.Card) {
	var covered viseraiOptSlots
	top = make([]card.Card, 0, len(cards))
	for _, c := range cards {
		slots := viseraiSlotsFor(c)
		if slots.overlaps(covered) {
			bottom = append(bottom, c)
			continue
		}
		top = append(top, c)
		covered = covered.union(slots)
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

// overlaps reports whether s and covered share at least one slot — used to detect
// "this card would over-fill a slot we've already kept a card for".
func (s viseraiOptSlots) overlaps(covered viseraiOptSlots) bool {
	return (s.nonAttackEnabler && covered.nonAttackEnabler) ||
		(s.nonGoAgainAction && covered.nonGoAgainAction) ||
		(s.defender && covered.defender) ||
		(s.bluePitch && covered.bluePitch)
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
func viseraiSlotsFor(c card.Card) viseraiOptSlots {
	// Opt-time heuristic reads printed values only. Types(nil) skips the Universal class
	// fold (class-independent predicates here). GoAgain(nil) returns the printed bool —
	// hero-conditional cards (Life for a Life / Blow for a Blow / Scar for a Scar) fall
	// back to "no go-again" since Viserai isn't a LowerHealthWanter anyway.
	t := c.Types(nil)
	return viseraiOptSlots{
		nonAttackEnabler: t.IsNonAttackAction(),
		nonGoAgainAction: t.Has(card.TypeAction) && !c.GoAgain(nil),
		defender:         t.IsDefenseReaction() || t.Has(card.TypeBlock),
		bluePitch:        c.Pitch() == 3,
	}
}
