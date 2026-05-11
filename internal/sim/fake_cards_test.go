package sim

// Fake Card / Hero implementations shared by the in-package sim tests. Lives in
// package sim (not testutils) because testutils imports sim — a `package sim` test
// file can't import testutils without cycling. The names mirror the testutils stubs
// the external (sim_test) tests use, just with a Fake prefix instead of Stub.

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
)

// FakeCard is a configurable Card stand-in. Tests construct it via NewFakeCard plus
// the With… builder methods so call sites set only the fields they care about; every
// other field returns its zero value through the Card interface.
type FakeCard struct {
	id      ids.CardID
	name    string
	types   card.TypeSet
	attack  int
	pitch   int
	defense int
	goAgain bool
}

// NewFakeCard returns a FakeCard with just a name set.
func NewFakeCard(name string) FakeCard { return FakeCard{name: name} }

func (c FakeCard) WithID(id ids.CardID) FakeCard     { c.id = id; return c }
func (c FakeCard) WithTypes(t card.TypeSet) FakeCard { c.types = t; return c }
func (c FakeCard) WithAttack(a int) FakeCard         { c.attack = a; return c }
func (c FakeCard) WithPitch(p int) FakeCard          { c.pitch = p; return c }
func (c FakeCard) WithDefense(d int) FakeCard        { c.defense = d; return c }
func (c FakeCard) WithGoAgain() FakeCard             { c.goAgain = true; return c }

func (c FakeCard) ID() ids.CardID                    { return c.id }
func (c FakeCard) Name() string                      { return c.name }
func (c FakeCard) DisplayName() string               { return c.name }
func (FakeCard) Cost(*TurnState) int                 { return 0 }
func (c FakeCard) Pitch() int                        { return c.pitch }
func (c FakeCard) Attack() int                       { return c.attack }
func (c FakeCard) Defense() int                      { return c.defense }
func (c FakeCard) Types() card.TypeSet               { return c.types }
func (c FakeCard) GoAgain() bool                     { return c.goAgain }
func (FakeCard) Play(*TurnState, Logger, *CardState) {}

// FakeRedAttack is a fixed-stat-line attack: 1-cost, 1-pitch, 3-power, 1-defense, no
// Go again. Mirrors testutils.RedAttack so tests that just need a "vanilla red attack"
// in a deck or hand have a direct stand-in.
type FakeRedAttack struct{}

func (FakeRedAttack) ID() ids.CardID      { return ids.InvalidCard }
func (FakeRedAttack) Name() string        { return "FakeRedAttack" }
func (FakeRedAttack) DisplayName() string { return "FakeRedAttack [R]" }
func (FakeRedAttack) Cost(*TurnState) int { return 1 }
func (FakeRedAttack) Pitch() int          { return 1 }
func (FakeRedAttack) Attack() int         { return 3 }
func (FakeRedAttack) Defense() int        { return 1 }
func (FakeRedAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (FakeRedAttack) GoAgain() bool                       { return false }
func (FakeRedAttack) Play(*TurnState, Logger, *CardState) {}

// DominatingFakeCard embeds FakeCard and adds the Dominator marker — exercises the
// printed-Dominate branch of EffectiveDominate / HasDominate.
type DominatingFakeCard struct{ FakeCard }

func (DominatingFakeCard) Dominate() {}

// FakeHero is a minimal Hero with configurable Intelligence and an optional Opt
// strategy. Every other method returns a fixed zero value so tests measure hand /
// deck behaviour in isolation from hero ability contributions.
type FakeHero struct {
	Intel       int
	OptStrategy func(cs []Card) (top, bottom []Card)
}

func (FakeHero) ID() ids.HeroID                            { return ids.InvalidHero }
func (FakeHero) Name() string                              { return "FakeHero" }
func (FakeHero) DisplayName() string                       { return "FakeHero" }
func (FakeHero) Health() int                               { return 20 }
func (h FakeHero) Intelligence() int                       { return h.Intel }
func (FakeHero) Types() card.TypeSet                       { return 0 }
func (FakeHero) Class() card.CardType                      { return 0 }
func (FakeHero) OnCardPlayed(Card, *TurnState, Logger) int { return 0 }

// Opt dispatches to OptStrategy when set; otherwise keeps every revealed card on top
// of the deck in input order (no reshape).
func (h FakeHero) Opt(cs []Card) (top, bottom []Card) {
	if h.OptStrategy != nil {
		return h.OptStrategy(cs)
	}
	return cs, nil
}
