// Fake-card builders. Every fake is one `Fake` struct configured vian attack turned `With...`
// methods so the call site spells out every attribute the test depends on — no
// cross-referencing a separate type definition to know that "FakeRedAttack" means
// "pitch 1, power 0, no Go again, no Play side effect" by default.
//
// Naming: pitch colour is encoded in the constructor (Red=1, Yellow=2, Blue=3). Card
// shape (Attack / Action / Resource / DR / Instant / Aura) is also in the constructor;
// each shape seeds the right base TypeSet. Defaults are 0 / false / no Play side effect
// for everything else.

package testutils

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/ids"
)

// Fake is the configurable card.Card returned by every FakeXxx... constructor below.
// Construct it via a constructor + follow-up With... methods; the zero value is not
// meaningful (no Name(), no types).
//
// The play hook is held behind a pointer so the struct itself stays comparable — the
// engine compares Card values directly in places like GameState.RemoveFromHand.
type Fake struct {
	name    string
	types   card.TypeSet
	cost    int
	pitch   int
	power   int
	defense int
	goAgain bool
	play    *playBody
}

// playBody wraps a Play func so Fake stays comparable. nil means "no side effect".
type playBody struct {
	fn func(card.GameEngine, card.Logger, *card.CardState)
}

// card.Card implementation.

func (Fake) ID() ids.CardID                       { return ids.InvalidCard }
func (f Fake) Name() string                       { return f.name }
func (f Fake) DisplayName() string                { return f.name }
func (f Fake) Cost() int                          { return f.cost }
func (f Fake) Pitch() int                         { return f.pitch }
func (f Fake) Attack() int                        { return f.power }
func (f Fake) Defense() int                       { return f.defense }
func (f Fake) Types(card.GameEngine) card.TypeSet { return f.types }
func (f Fake) GoAgain(card.GameEngine) bool       { return f.goAgain }
func (f Fake) Play(ge card.GameEngine, l card.Logger, self *card.CardState) {
	if f.play != nil {
		f.play.fn(ge, l, self)
	}
}

// Builder methods. Each returns a copy of f with the named field overridden so chains
// stay value-typed and safe to share.

// WithPower sets the printed attack power.
func (f Fake) WithPower(n int) Fake { f.power = n; return f }

// WithCost sets the action-point cost (default 0).
func (f Fake) WithCost(n int) Fake { f.cost = n; return f }

// WithDefense sets the printed defense.
func (f Fake) WithDefense(n int) Fake { f.defense = n; return f }

// WithPitch overrides the pitch value set by the constructor. Use only for the rare
// test that iterates pitch values dynamically — prefer the colour-named constructor
// (FakeRedAttack / FakeYellowAttack / FakeBlueAttack) when the pitch is known at
// authoring time so the colour reads at the call site.
func (f Fake) WithPitch(n int) Fake { f.pitch = n; return f }

// WithGoAgain sets goAgain=true.
func (f Fake) WithGoAgain() Fake { f.goAgain = true; return f }

// WithTypes ORs the given subtypes into the card's TypeSet — additive, not
// replacement. Use for subtype tags (Runeblade, Hammer, Sword, …) where the base shape
// from the constructor stays.
func (f Fake) WithTypes(types ...card.CardType) Fake {
	f.types = f.types | card.NewTypeSet(types...)
	return f
}

// WithName overrides the default display name. Useful when a test asserts on log
// output and needs to tell two same-shape fakes apart.
func (f Fake) WithName(name string) Fake { f.name = name; return f }

// WithPlay installs a custom Play body for behaviour-spy fakes.
func (f Fake) WithPlay(playFunc func(card.GameEngine, card.Logger, *card.CardState)) Fake {
	f.play = &playBody{fn: playFunc}
	return f
}

// WithDrawOne installs a Play body that draws one card. Convenience for the common
// "this card draws on resolution" test fixture.
func (f Fake) WithDrawOne() Fake {
	f.play = &playBody{fn: func(ge card.GameEngine, _ card.Logger, _ *card.CardState) { ge.DrawOne() }}
	return f
}

// Base type sets shared by the constructors below.
var (
	fakeActionAttackTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
	fakeActionTypes       = card.NewTypeSet(card.TypeGeneric, card.TypeAction)
	fakeResourceTypes     = card.NewTypeSet(card.TypeGeneric, card.TypeResource)
	fakeDRTypes           = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)
	fakeInstantTypes      = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeInstant)
	fakeAuraTypes         = card.NewTypeSet(card.TypeAura)
	fakeWeaponSwingTypes  = card.NewTypeSet(card.TypeGeneric, card.TypeWeapon, card.TypeAttack)
)

// Action-attack constructors. Pitch implied by the colour in the name: Red=1, Yellow=2,
// Blue=3. Defaults: cost=0, power=0, defense=0, no Go again, no Play side effect.

func FakeRedAttack() Fake {
	return Fake{name: "FakeRedAttack", pitch: 1, types: fakeActionAttackTypes}
}
func FakeYellowAttack() Fake {
	return Fake{name: "FakeYellowAttack", pitch: 2, types: fakeActionAttackTypes}
}
func FakeBlueAttack() Fake {
	return Fake{name: "FakeBlueAttack", pitch: 3, types: fakeActionAttackTypes}
}

// Non-attack action constructors.

func FakeRedAction() Fake    { return Fake{name: "FakeRedAction", pitch: 1, types: fakeActionTypes} }
func FakeYellowAction() Fake { return Fake{name: "FakeYellowAction", pitch: 2, types: fakeActionTypes} }
func FakeBlueAction() Fake   { return Fake{name: "FakeBlueAction", pitch: 3, types: fakeActionTypes} }

// Pure-pitch Resource constructors (TypeResource — can only pitch, never plays).

func FakeRedResource() Fake {
	return Fake{name: "FakeRedResource", pitch: 1, types: fakeResourceTypes}
}
func FakeYellowResource() Fake {
	return Fake{name: "FakeYellowResource", pitch: 2, types: fakeResourceTypes}
}
func FakeBlueResource() Fake {
	return Fake{name: "FakeBlueResource", pitch: 3, types: fakeResourceTypes}
}

// Defense reaction constructors.

func FakeRedDR() Fake    { return Fake{name: "FakeRedDR", pitch: 1, types: fakeDRTypes} }
func FakeYellowDR() Fake { return Fake{name: "FakeYellowDR", pitch: 2, types: fakeDRTypes} }
func FakeBlueDR() Fake   { return Fake{name: "FakeBlueDR", pitch: 3, types: fakeDRTypes} }

// Instant constructors.

func FakeRedInstant() Fake {
	return Fake{name: "FakeRedInstant", pitch: 1, types: fakeInstantTypes}
}
func FakeYellowInstant() Fake {
	return Fake{name: "FakeYellowInstant", pitch: 2, types: fakeInstantTypes}
}
func FakeBlueInstant() Fake {
	return Fake{name: "FakeBlueInstant", pitch: 3, types: fakeInstantTypes}
}

// Aura constructors.

func FakeRedAura() Fake    { return Fake{name: "FakeRedAura", pitch: 1, types: fakeAuraTypes} }
func FakeYellowAura() Fake { return Fake{name: "FakeYellowAura", pitch: 2, types: fakeAuraTypes} }
func FakeBlueAura() Fake   { return Fake{name: "FakeBlueAura", pitch: 3, types: fakeAuraTypes} }

// Block-only card constructors — Generic + Block, no Action / Attack / DefenseReaction.
// Used by Opt-strategy tests that need a TypeBlock card competing with DRs for the
// defender slot. Add Defense via WithDefense.

func FakeRedBlocker() Fake {
	return Fake{name: "FakeRedBlocker", pitch: 1, types: card.NewTypeSet(card.TypeGeneric, card.TypeBlock)}
}
func FakeYellowBlocker() Fake {
	return Fake{name: "FakeYellowBlocker", pitch: 2, types: card.NewTypeSet(card.TypeGeneric, card.TypeBlock)}
}
func FakeBlueBlocker() Fake {
	return Fake{name: "FakeBlueBlocker", pitch: 3, types: card.NewTypeSet(card.TypeGeneric, card.TypeBlock)}
}

// FakeWeaponSwing is the activated-ability Card for a hypothetical weapon swing — base
// types: Weapon + Attack. Subtypes like Hammer / Sword are added via WithTypes. Pitch
// defaults to 0 (weapons aren't in the deck so they don't pitch).
func FakeWeaponSwing() Fake {
	return Fake{name: "FakeWeaponSwing", types: fakeWeaponSwingTypes}
}
