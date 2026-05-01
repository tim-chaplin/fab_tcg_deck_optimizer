// Package testutils provides Card stubs and fake card implementations shared by tests in
// multiple packages (card, cards, deck, hand, sim). The configurable Card stub framework
// (Card, GenericAttack, RunebladeAttack, …) builds CardsRemaining / CardsPlayed / Pitched
// lists with specific type, cost, power, and pitch shapes so predicate / lookahead tests
// have predictable inputs. The fixed-stat-line fakes (RedAttack, BlueAttack, YellowAttack,
// …) are deliberately simple attack actions tests use as deck contents when partition /
// ordering assertions need known optimal values.
package testutils

import (
	"github.com/tim-chaplin/fab-deck-optimizer/internal/card"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"
	"github.com/tim-chaplin/fab-deck-optimizer/internal/sim"
)

// Card is a configurable Card implementation used across tests to build CardsRemaining /
// CardsPlayed / Pitched lists with specific type, cost, power, and pitch shapes. Zero-value
// fields mean "don't care" — tests set only what the helper under test predicates on.
type Card struct {
	name  string
	cost  int
	power int
	pitch int
	types card.TypeSet
}

func (s Card) ID() ids.CardID { return ids.InvalidCard }
func (s Card) Name() string   { return s.name }

// WithName returns a copy of s with its display name overridden. Lets cross-package tests
// distinguish multiple Card stubs in log assertions even though `name` is unexported.
func (s Card) WithName(name string) Card { s.name = name; return s }

func (s Card) Cost(*sim.TurnState) int           { return s.cost }
func (s Card) Pitch() int                        { return s.pitch }
func (s Card) Attack() int                       { return s.power }
func (s Card) Defense() int                      { return 0 }
func (s Card) Types() card.TypeSet               { return s.types }
func (s Card) GoAgain() bool                     { return false }
func (Card) Play(*sim.TurnState, *sim.CardState) {}

// GenericAttack returns a Generic Action - Attack stub with the given cost and base power.
// Pitch defaults to 1; override via GenericAttackPitch if a test cares.
func GenericAttack(cost, power int) Card {
	return Card{
		name:  "GenericAttack",
		cost:  cost,
		power: power,
		pitch: 1,
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack),
	}
}

// GenericAttackPitch is GenericAttack with an explicit pitch value. Flying High's red
// variant rider reads pitch, so tests that exercise the +1 bonus set this.
func GenericAttackPitch(cost, power, pitch int) Card {
	s := GenericAttack(cost, power)
	s.pitch = pitch
	return s
}

// GenericAction returns a Generic Action (non-attack) stub for attack-typed-lookahead
// rejection cases.
func GenericAction() Card {
	return Card{
		name:  "GenericAction",
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAction),
	}
}

// GenericAura returns a Generic Aura stub — covers Yinti Yanti's HasPlayedType(TypeAura) check.
func GenericAura() Card {
	return Card{
		name:  "GenericAura",
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAura),
	}
}

// Shared stub Cards. Each is a zero-value struct with a fixed type line; tests mix and match to
// exercise lookahead / predicate logic on card effects.

// RunebladeAttack is a minimal Runeblade Action-Attack card — satisfies "next Runeblade
// attack action card" lookaheads.
type RunebladeAttack struct{}

func (RunebladeAttack) ID() ids.CardID          { return ids.InvalidCard }
func (RunebladeAttack) Name() string            { return "RunebladeAttack" }
func (RunebladeAttack) Cost(*sim.TurnState) int { return 0 }
func (RunebladeAttack) Pitch() int              { return 0 }
func (RunebladeAttack) Attack() int             { return 0 }
func (RunebladeAttack) Defense() int            { return 0 }
func (RunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (RunebladeAttack) GoAgain() bool                       { return true }
func (RunebladeAttack) Play(*sim.TurnState, *sim.CardState) {}

// RunebladeWeapon is a Runeblade weapon — satisfies "next Runeblade attack" lookaheads
// that include weapons but NOT ones restricted to attack action cards.
type RunebladeWeapon struct{}

func (RunebladeWeapon) ID() ids.CardID          { return ids.InvalidCard }
func (RunebladeWeapon) Name() string            { return "RunebladeWeapon" }
func (RunebladeWeapon) Cost(*sim.TurnState) int { return 0 }
func (RunebladeWeapon) Pitch() int              { return 0 }
func (RunebladeWeapon) Attack() int             { return 0 }
func (RunebladeWeapon) Defense() int            { return 0 }
func (RunebladeWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (RunebladeWeapon) GoAgain() bool                       { return false }
func (RunebladeWeapon) Play(*sim.TurnState, *sim.CardState) {}

// NonAttack is a non-attack card — covers "attack-typed predicate should reject
// non-attack" cases.
type NonAttack struct{}

func (NonAttack) ID() ids.CardID                      { return ids.InvalidCard }
func (NonAttack) Name() string                        { return "NonAttack" }
func (NonAttack) Cost(*sim.TurnState) int             { return 0 }
func (NonAttack) Pitch() int                          { return 0 }
func (NonAttack) Attack() int                         { return 0 }
func (NonAttack) Defense() int                        { return 0 }
func (NonAttack) Types() card.TypeSet                 { return card.NewTypeSet(card.TypeAction) }
func (NonAttack) GoAgain() bool                       { return false }
func (NonAttack) Play(*sim.TurnState, *sim.CardState) {}

// NonRunebladeAttack is a Generic Action-Attack — covers Runeblade-gated lookaheads
// rejecting non-Runeblade attacks.
type NonRunebladeAttack struct{}

func (NonRunebladeAttack) ID() ids.CardID          { return ids.InvalidCard }
func (NonRunebladeAttack) Name() string            { return "NonRunebladeAttack" }
func (NonRunebladeAttack) Cost(*sim.TurnState) int { return 0 }
func (NonRunebladeAttack) Pitch() int              { return 0 }
func (NonRunebladeAttack) Attack() int             { return 0 }
func (NonRunebladeAttack) Defense() int            { return 0 }
func (NonRunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (NonRunebladeAttack) GoAgain() bool                       { return true }
func (NonRunebladeAttack) Play(*sim.TurnState, *sim.CardState) {}

// AttackWithPower is a Runeblade attack-action card with a configurable printed Attack()
// value. Tests set specific numbers to hit/miss the LikelyToHit heuristic (4 lands, 3 blocks).
type AttackWithPower struct {
	Power int
}

func (AttackWithPower) ID() ids.CardID          { return ids.InvalidCard }
func (AttackWithPower) Name() string            { return "AttackWithPower" }
func (AttackWithPower) Cost(*sim.TurnState) int { return 0 }
func (AttackWithPower) Pitch() int              { return 0 }
func (s AttackWithPower) Attack() int           { return s.Power }
func (AttackWithPower) Defense() int            { return 0 }
func (AttackWithPower) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (AttackWithPower) GoAgain() bool                       { return true }
func (AttackWithPower) Play(*sim.TurnState, *sim.CardState) {}

// Aura is a minimal Aura-typed card — exercises "aura played this turn" checks.
type Aura struct{}

func (Aura) ID() ids.CardID                      { return ids.InvalidCard }
func (Aura) Name() string                        { return "Aura" }
func (Aura) Cost(*sim.TurnState) int             { return 0 }
func (Aura) Pitch() int                          { return 0 }
func (Aura) Attack() int                         { return 0 }
func (Aura) Defense() int                        { return 0 }
func (Aura) Types() card.TypeSet                 { return card.NewTypeSet(card.TypeAura) }
func (Aura) GoAgain() bool                       { return true }
func (Aura) Play(*sim.TurnState, *sim.CardState) {}

// genericAttackTypes is the type line shared by every attack-action fake below.
var genericAttackTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// BluePitch is a pure-pitch generic non-attack action: pitches 3, no attack, no defense,
// no go again. Useful as a "blue pitch source" in tests where the optimal line should be
// unambiguous — the optimizer can't repurpose it as an attacker or blocker, so the only
// reasonable use is to pitch it.
type BluePitch struct{}

func (BluePitch) ID() ids.CardID          { return FakeBluePitch }
func (BluePitch) Name() string            { return "cardtest.BluePitch" }
func (BluePitch) Cost(*sim.TurnState) int { return 0 }
func (BluePitch) Pitch() int              { return 3 }
func (BluePitch) Attack() int             { return 0 }
func (BluePitch) Defense() int            { return 0 }
func (BluePitch) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction)
}
func (BluePitch) GoAgain() bool                              { return false }
func (BluePitch) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

// BlueAttack is a generic blue attack action: pitches 3, defends 3, attacks 1, costs 1.
type BlueAttack struct{}

func (BlueAttack) ID() ids.CardID                             { return FakeBlueAttack }
func (BlueAttack) Name() string                               { return "cardtest.BlueAttack" }
func (BlueAttack) Cost(*sim.TurnState) int                    { return 1 }
func (BlueAttack) Pitch() int                                 { return 3 }
func (BlueAttack) Attack() int                                { return 1 }
func (BlueAttack) Defense() int                               { return 3 }
func (BlueAttack) Types() card.TypeSet                        { return genericAttackTypes }
func (BlueAttack) GoAgain() bool                              { return true }
func (BlueAttack) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// RedAttack is a generic red attack action: pitches 1, defends 1, attacks 3, costs 1.
type RedAttack struct{}

func (RedAttack) ID() ids.CardID                             { return FakeRedAttack }
func (RedAttack) Name() string                               { return "cardtest.RedAttack" }
func (RedAttack) Cost(*sim.TurnState) int                    { return 1 }
func (RedAttack) Pitch() int                                 { return 1 }
func (RedAttack) Attack() int                                { return 3 }
func (RedAttack) Defense() int                               { return 1 }
func (RedAttack) Types() card.TypeSet                        { return genericAttackTypes }
func (RedAttack) GoAgain() bool                              { return true }
func (RedAttack) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// YellowAttack is a generic yellow attack action: pitches 2, defends 2, attacks 2, costs 1.
type YellowAttack struct{}

func (YellowAttack) ID() ids.CardID                             { return FakeYellowAttack }
func (YellowAttack) Name() string                               { return "cardtest.YellowAttack" }
func (YellowAttack) Cost(*sim.TurnState) int                    { return 1 }
func (YellowAttack) Pitch() int                                 { return 2 }
func (YellowAttack) Attack() int                                { return 2 }
func (YellowAttack) Defense() int                               { return 2 }
func (YellowAttack) Types() card.TypeSet                        { return genericAttackTypes }
func (YellowAttack) GoAgain() bool                              { return true }
func (YellowAttack) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// genericActionTypes is a plain non-attack action (no Attack subtype). Used by CostlyDraw — a
// draw-a-card action card. It isn't a Defense Reaction so it can't be played on the opponent's
// turn; it carries Go again so a drawn-card-as-extension can chain.
var genericActionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// CostlyDraw is a 1-cost, pitch-1, no-damage "draw a card, go again" action. Used by the
// mid-turn-draw determinism tests: plays for 0 damage but a chain can continue off its Go again
// onto a drawn-later attack.
type CostlyDraw struct{}

func (CostlyDraw) ID() ids.CardID          { return FakeCostlyDraw }
func (CostlyDraw) Name() string            { return "cardtest.CostlyDraw" }
func (CostlyDraw) Cost(*sim.TurnState) int { return 1 }
func (CostlyDraw) Pitch() int              { return 1 }
func (CostlyDraw) Attack() int             { return 0 }
func (CostlyDraw) Defense() int            { return 0 }
func (CostlyDraw) Types() card.TypeSet     { return genericActionTypes }
func (CostlyDraw) GoAgain() bool           { return true }
func (CostlyDraw) Play(s *sim.TurnState, self *sim.CardState) {
	s.DrawOne()
	s.LogPlay(self)
}

// CostlyAttack is a 1-cost, pitch-1, 3-damage attack action — the "deal 3 damage" alternative
// the mid-turn-draw determinism test weighs against CostlyDraw.
type CostlyAttack struct{}

func (CostlyAttack) ID() ids.CardID                             { return FakeCostlyAttack }
func (CostlyAttack) Name() string                               { return "cardtest.CostlyAttack" }
func (CostlyAttack) Cost(*sim.TurnState) int                    { return 1 }
func (CostlyAttack) Pitch() int                                 { return 1 }
func (CostlyAttack) Attack() int                                { return 3 }
func (CostlyAttack) Defense() int                               { return 0 }
func (CostlyAttack) Types() card.TypeSet                        { return genericAttackTypes }
func (CostlyAttack) GoAgain() bool                              { return false }
func (CostlyAttack) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }

var genericDefenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)

// PitchOneDR is a 1-pitch-value defense reaction: cost 0, defense 3. It exists solely so tests
// can pitch it (contributing 1 resource) to fund another 1-cost card without also playing it.
type PitchOneDR struct{}

func (PitchOneDR) ID() ids.CardID          { return FakePitchOneDR }
func (PitchOneDR) Name() string            { return "cardtest.PitchOneDR" }
func (PitchOneDR) Cost(*sim.TurnState) int { return 0 }
func (PitchOneDR) Pitch() int              { return 1 }
func (PitchOneDR) Attack() int             { return 0 }
func (PitchOneDR) Defense() int            { return 3 }
func (PitchOneDR) Types() card.TypeSet     { return genericDefenseReactionTypes }
func (PitchOneDR) GoAgain() bool           { return false }
func (PitchOneDR) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

// HugeAttack is a 0-cost "do one million damage" attack. Outrageous on purpose: as the top of
// the deck it makes the CostlyDraw → HugeAttack chain blatantly better than the CostlyAttack
// line, so an evaluator that peeks at the deck will pick a different role for the hand's draw
// card depending on deck order — the determinism test catches that.
type HugeAttack struct{}

const hugeAttackDamage = 1_000_000

func (HugeAttack) ID() ids.CardID                             { return FakeHugeAttack }
func (HugeAttack) Name() string                               { return "cardtest.HugeAttack" }
func (HugeAttack) Cost(*sim.TurnState) int                    { return 0 }
func (HugeAttack) Pitch() int                                 { return 1 }
func (HugeAttack) Attack() int                                { return hugeAttackDamage }
func (HugeAttack) Defense() int                               { return 0 }
func (HugeAttack) Types() card.TypeSet                        { return genericAttackTypes }
func (HugeAttack) GoAgain() bool                              { return false }
func (HugeAttack) Play(s *sim.TurnState, self *sim.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// StubCard is a minimal sim.Card. Tests construct it via NewStubCard plus the With…
// builder methods so call sites only set the fields they care about — every other field
// returns a zero value through the Card interface methods. ID defaults to InvalidCard;
// tests that reach into ID-keyed caches (cardMetaCache, chainStepCache) should attach a
// distinct ID via WithID(testutils.FakeX) to avoid sharing slot 0.
type StubCard struct {
	id      ids.CardID
	name    string
	types   card.TypeSet
	attack  int
	pitch   int
	defense int
	goAgain bool
}

// NewStubCard returns a StubCard with just a name set. Use the With… mutators below to
// attach additional fields.
func NewStubCard(name string) StubCard { return StubCard{name: name} }

// WithID returns a copy of c with id set.
func (c StubCard) WithID(id ids.CardID) StubCard { c.id = id; return c }

// WithTypes returns a copy of c with types set.
func (c StubCard) WithTypes(t card.TypeSet) StubCard { c.types = t; return c }

// WithAttack returns a copy of c with the printed attack value set.
func (c StubCard) WithAttack(a int) StubCard { c.attack = a; return c }

// WithPitch returns a copy of c with the printed pitch value set (1=red, 2=yellow, 3=blue).
func (c StubCard) WithPitch(p int) StubCard { c.pitch = p; return c }

// WithDefense returns a copy of c with the printed defense value set.
func (c StubCard) WithDefense(d int) StubCard { c.defense = d; return c }

// WithGoAgain returns a copy of c with goAgain=true.
func (c StubCard) WithGoAgain() StubCard { c.goAgain = true; return c }

func (c StubCard) ID() ids.CardID                    { return c.id }
func (c StubCard) Name() string                      { return c.name }
func (StubCard) Cost(*sim.TurnState) int             { return 0 }
func (c StubCard) Pitch() int                        { return c.pitch }
func (c StubCard) Attack() int                       { return c.attack }
func (c StubCard) Defense() int                      { return c.defense }
func (c StubCard) Types() card.TypeSet               { return c.types }
func (c StubCard) GoAgain() bool                     { return c.goAgain }
func (StubCard) Play(*sim.TurnState, *sim.CardState) {}

// DominatingStubCard embeds StubCard and adds the sim.Dominator marker — exercises the
// printed-Dominate branch of EffectiveDominate / HasDominate.
type DominatingStubCard struct{ StubCard }

func (DominatingStubCard) Dominate() {}

// NotImplementedStubCard embeds StubCard and adds the sim.NotImplemented marker —
// exercises the type assertion the deck legal-pool filter keys on.
type NotImplementedStubCard struct{ StubCard }

func (NotImplementedStubCard) NotImplemented() {}

// UnplayableStubCard embeds StubCard and adds the sim.Unplayable marker — exercises the
// second pool-exclusion path the deck legal-pool filter keys on.
type UnplayableStubCard struct{ StubCard }

func (UnplayableStubCard) Unplayable() {}

// InstantStub is a 0-cost, 0-power Generic Action - Instant card with no Go again.
// Tests chain-runner behaviour around the Action Point debit: an Instant after a
// non-Go-again card should still resolve because Instants cost 0 AP.
type InstantStub struct{}

func (InstantStub) ID() ids.CardID          { return FakeInstant }
func (InstantStub) Name() string            { return "InstantStub" }
func (InstantStub) Cost(*sim.TurnState) int { return 0 }
func (InstantStub) Pitch() int              { return 0 }
func (InstantStub) Attack() int             { return 0 }
func (InstantStub) Defense() int            { return 0 }
func (InstantStub) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeInstant)
}
func (InstantStub) GoAgain() bool                              { return false }
func (InstantStub) Play(s *sim.TurnState, self *sim.CardState) { s.LogPlay(self) }

// NoGoAgainAttackStub is a 0-cost, 1-power Generic Action - Attack card with no Go again.
// Tests chain-runner behaviour after the AP pool runs out: a non-Instant follow-up
// should be rejected.
type NoGoAgainAttackStub struct{}

func (NoGoAgainAttackStub) ID() ids.CardID          { return FakeNoGoAgainAttack }
func (NoGoAgainAttackStub) Name() string            { return "NoGoAgainAttack" }
func (NoGoAgainAttackStub) Cost(*sim.TurnState) int { return 0 }
func (NoGoAgainAttackStub) Pitch() int              { return 0 }
func (NoGoAgainAttackStub) Attack() int             { return 1 }
func (NoGoAgainAttackStub) Defense() int            { return 0 }
func (NoGoAgainAttackStub) Types() card.TypeSet     { return genericAttackTypes }
func (NoGoAgainAttackStub) GoAgain() bool           { return false }
func (NoGoAgainAttackStub) Play(s *sim.TurnState, self *sim.CardState) {
	s.ApplyAndLogEffectiveAttack(self)
}

// GrantAll is a Runeblade attack-action card that flips GrantedGoAgain=true on every
// remaining CardState in CardsRemaining when it resolves. Used with GrantSpy to detect
// cross-permutation CardState wrapper leakage in bestSequence: the fresh-wrapper
// invariant must keep grants from bleeding across permutations.
type GrantAll struct{}

func (GrantAll) ID() ids.CardID          { return ids.InvalidCard }
func (GrantAll) Name() string            { return "GrantAll" }
func (GrantAll) Cost(*sim.TurnState) int { return 0 }
func (GrantAll) Pitch() int              { return 0 }
func (GrantAll) Attack() int             { return 0 }
func (GrantAll) Defense() int            { return 0 }
func (GrantAll) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (GrantAll) GoAgain() bool { return true }
func (GrantAll) Play(s *sim.TurnState, self *sim.CardState) {
	for _, pc := range s.CardsRemaining {
		pc.GrantedGoAgain = true
	}
	s.LogPlay(self)
}

// GrantSpy is a Runeblade attack-action card. When it plays first in a permutation it
// records (via *Saw) whether any CardState in CardsRemaining already has
// GrantedGoAgain=true. With per-permutation fresh wrappers that should never happen —
// no prior card in this permutation has run yet. If wrappers leak across permutations,
// a prior permutation's GrantAll Play would still be visible and the spy trips.
type GrantSpy struct{ Saw *bool }

func (GrantSpy) ID() ids.CardID          { return ids.InvalidCard }
func (GrantSpy) Name() string            { return "GrantSpy" }
func (GrantSpy) Cost(*sim.TurnState) int { return 0 }
func (GrantSpy) Pitch() int              { return 0 }
func (GrantSpy) Attack() int             { return 0 }
func (GrantSpy) Defense() int            { return 0 }
func (GrantSpy) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (GrantSpy) GoAgain() bool { return true }
func (g GrantSpy) Play(s *sim.TurnState, self *sim.CardState) {
	defer s.LogPlay(self)
	if len(s.CardsPlayed) != 0 {
		return
	}
	for _, pc := range s.CardsRemaining {
		if pc.GrantedGoAgain {
			*g.Saw = true
		}
	}
}

// CardNames renders a slice of Card names for test failure messages.
func CardNames(cs []sim.Card) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		out[i] = c.Name()
	}
	return out
}
