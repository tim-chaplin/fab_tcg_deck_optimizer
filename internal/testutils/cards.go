// Package testutils provides Card stubs and fake card implementations shared by tests in
// multiple packages (card, cards, deck, hand, sim). The configurable Card stub framework
// (Card, GenericAttack, RunebladeAttack, …) builds CardsRemaining / CardsPlayed / Pitched
// lists with specific type, cost, power, and pitch shapes so predicate / lookahead tests
// have predictable inputs. The fixed-stat-line fakes (RedAttack, BlueAttack, DrawCantrip,
// …) are deliberately simple attack actions tests use as deck contents when partition /
// ordering assertions need known optimal values.
package testutils

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

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

func (s Card) ID() card.ID  { return card.Invalid }
func (s Card) Name() string { return s.name }

// WithName returns a copy of s with its display name overridden. Lets cross-package tests
// distinguish multiple Card stubs in log assertions even though `name` is unexported.
func (s Card) WithName(name string) Card { s.name = name; return s }

func (s Card) Cost(*card.TurnState) int            { return s.cost }
func (s Card) Pitch() int                          { return s.pitch }
func (s Card) Attack() int                         { return s.power }
func (s Card) Defense() int                        { return 0 }
func (s Card) Types() card.TypeSet                 { return s.types }
func (s Card) GoAgain() bool                       { return false }
func (Card) Play(*card.TurnState, *card.CardState) {}

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

func (RunebladeAttack) ID() card.ID              { return card.Invalid }
func (RunebladeAttack) Name() string             { return "RunebladeAttack" }
func (RunebladeAttack) Cost(*card.TurnState) int { return 0 }
func (RunebladeAttack) Pitch() int               { return 0 }
func (RunebladeAttack) Attack() int              { return 0 }
func (RunebladeAttack) Defense() int             { return 0 }
func (RunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (RunebladeAttack) GoAgain() bool                         { return true }
func (RunebladeAttack) Play(*card.TurnState, *card.CardState) {}

// RunebladeWeapon is a Runeblade weapon — satisfies "next Runeblade attack" lookaheads
// that include weapons but NOT ones restricted to attack action cards.
type RunebladeWeapon struct{}

func (RunebladeWeapon) ID() card.ID              { return card.Invalid }
func (RunebladeWeapon) Name() string             { return "RunebladeWeapon" }
func (RunebladeWeapon) Cost(*card.TurnState) int { return 0 }
func (RunebladeWeapon) Pitch() int               { return 0 }
func (RunebladeWeapon) Attack() int              { return 0 }
func (RunebladeWeapon) Defense() int             { return 0 }
func (RunebladeWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (RunebladeWeapon) GoAgain() bool                         { return false }
func (RunebladeWeapon) Play(*card.TurnState, *card.CardState) {}

// NonAttack is a non-attack card — covers "attack-typed predicate should reject
// non-attack" cases.
type NonAttack struct{}

func (NonAttack) ID() card.ID                           { return card.Invalid }
func (NonAttack) Name() string                          { return "NonAttack" }
func (NonAttack) Cost(*card.TurnState) int              { return 0 }
func (NonAttack) Pitch() int                            { return 0 }
func (NonAttack) Attack() int                           { return 0 }
func (NonAttack) Defense() int                          { return 0 }
func (NonAttack) Types() card.TypeSet                   { return card.NewTypeSet(card.TypeAction) }
func (NonAttack) GoAgain() bool                         { return false }
func (NonAttack) Play(*card.TurnState, *card.CardState) {}

// NonRunebladeAttack is a Generic Action-Attack — covers Runeblade-gated lookaheads
// rejecting non-Runeblade attacks.
type NonRunebladeAttack struct{}

func (NonRunebladeAttack) ID() card.ID              { return card.Invalid }
func (NonRunebladeAttack) Name() string             { return "NonRunebladeAttack" }
func (NonRunebladeAttack) Cost(*card.TurnState) int { return 0 }
func (NonRunebladeAttack) Pitch() int               { return 0 }
func (NonRunebladeAttack) Attack() int              { return 0 }
func (NonRunebladeAttack) Defense() int             { return 0 }
func (NonRunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (NonRunebladeAttack) GoAgain() bool                         { return true }
func (NonRunebladeAttack) Play(*card.TurnState, *card.CardState) {}

// AttackWithPower is a Runeblade attack-action card with a configurable printed Attack()
// value. Tests set specific numbers to hit/miss the LikelyToHit heuristic (4 lands, 3 blocks).
type AttackWithPower struct {
	Power int
}

func (AttackWithPower) ID() card.ID              { return card.Invalid }
func (AttackWithPower) Name() string             { return "AttackWithPower" }
func (AttackWithPower) Cost(*card.TurnState) int { return 0 }
func (AttackWithPower) Pitch() int               { return 0 }
func (s AttackWithPower) Attack() int            { return s.Power }
func (AttackWithPower) Defense() int             { return 0 }
func (AttackWithPower) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (AttackWithPower) GoAgain() bool                         { return true }
func (AttackWithPower) Play(*card.TurnState, *card.CardState) {}

// Aura is a minimal Aura-typed card — exercises "aura played this turn" checks.
type Aura struct{}

func (Aura) ID() card.ID                           { return card.Invalid }
func (Aura) Name() string                          { return "Aura" }
func (Aura) Cost(*card.TurnState) int              { return 0 }
func (Aura) Pitch() int                            { return 0 }
func (Aura) Attack() int                           { return 0 }
func (Aura) Defense() int                          { return 0 }
func (Aura) Types() card.TypeSet                   { return card.NewTypeSet(card.TypeAura) }
func (Aura) GoAgain() bool                         { return true }
func (Aura) Play(*card.TurnState, *card.CardState) {}

// genericAttackTypes is the type line shared by every attack-action fake below.
var genericAttackTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// BlueAttack is a generic blue attack action: pitches 3, defends 3, attacks 1, costs 1.
type BlueAttack struct{}

func (BlueAttack) ID() card.ID                                  { return card.FakeBlueAttack }
func (BlueAttack) Name() string                                 { return "cardtest.BlueAttack" }
func (BlueAttack) Cost(*card.TurnState) int                     { return 1 }
func (BlueAttack) Pitch() int                                   { return 3 }
func (BlueAttack) Attack() int                                  { return 1 }
func (BlueAttack) Defense() int                                 { return 3 }
func (BlueAttack) Types() card.TypeSet                          { return genericAttackTypes }
func (BlueAttack) GoAgain() bool                                { return true }
func (BlueAttack) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// RedAttack is a generic red attack action: pitches 1, defends 1, attacks 3, costs 1.
type RedAttack struct{}

func (RedAttack) ID() card.ID                                  { return card.FakeRedAttack }
func (RedAttack) Name() string                                 { return "cardtest.RedAttack" }
func (RedAttack) Cost(*card.TurnState) int                     { return 1 }
func (RedAttack) Pitch() int                                   { return 1 }
func (RedAttack) Attack() int                                  { return 3 }
func (RedAttack) Defense() int                                 { return 1 }
func (RedAttack) Types() card.TypeSet                          { return genericAttackTypes }
func (RedAttack) GoAgain() bool                                { return true }
func (RedAttack) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// YellowAttack is a generic yellow attack action: pitches 2, defends 2, attacks 2, costs 1.
type YellowAttack struct{}

func (YellowAttack) ID() card.ID                                  { return card.FakeYellowAttack }
func (YellowAttack) Name() string                                 { return "cardtest.YellowAttack" }
func (YellowAttack) Cost(*card.TurnState) int                     { return 1 }
func (YellowAttack) Pitch() int                                   { return 2 }
func (YellowAttack) Attack() int                                  { return 2 }
func (YellowAttack) Defense() int                                 { return 2 }
func (YellowAttack) Types() card.TypeSet                          { return genericAttackTypes }
func (YellowAttack) GoAgain() bool                                { return true }
func (YellowAttack) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }

// DrawCantrip is a generic free-cycling attack: cost 0, pitches 1, attacks 1, go again, and
// fires DrawOne on play. Used by tests to exercise mid-turn-draw chains that extend themselves
// (each cantrip plays, draws the next one, which plays, etc.).
type DrawCantrip struct{}

func (DrawCantrip) ID() card.ID              { return card.FakeDrawCantrip }
func (DrawCantrip) Name() string             { return "cardtest.DrawCantrip" }
func (DrawCantrip) Cost(*card.TurnState) int { return 0 }
func (DrawCantrip) Pitch() int               { return 1 }
func (DrawCantrip) Attack() int              { return 1 }
func (DrawCantrip) Defense() int             { return 0 }
func (DrawCantrip) Types() card.TypeSet      { return genericAttackTypes }
func (DrawCantrip) GoAgain() bool            { return true }
func (c DrawCantrip) Play(s *card.TurnState, self *card.CardState) {
	s.DrawOne()
	s.ApplyAndLogEffectiveAttack(self)
}

// genericActionTypes is a plain non-attack action (no Attack subtype). Used by CostlyDraw — a
// draw-a-card action card. It isn't a Defense Reaction so it can't be played on the opponent's
// turn; it carries Go again so a drawn-card-as-extension can chain.
var genericActionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// CostlyDraw is a 1-cost, pitch-1, no-damage "draw a card, go again" action. Used by the
// mid-turn-draw determinism tests: plays for 0 damage but a chain can continue off its Go again
// onto a drawn-later attack.
type CostlyDraw struct{}

func (CostlyDraw) ID() card.ID              { return card.FakeCostlyDraw }
func (CostlyDraw) Name() string             { return "cardtest.CostlyDraw" }
func (CostlyDraw) Cost(*card.TurnState) int { return 1 }
func (CostlyDraw) Pitch() int               { return 1 }
func (CostlyDraw) Attack() int              { return 0 }
func (CostlyDraw) Defense() int             { return 0 }
func (CostlyDraw) Types() card.TypeSet      { return genericActionTypes }
func (CostlyDraw) GoAgain() bool            { return true }
func (CostlyDraw) Play(s *card.TurnState, self *card.CardState) {
	s.DrawOne()
	s.LogPlay(self)
}

// CostlyAttack is a 1-cost, pitch-1, 3-damage attack action — the "deal 3 damage" alternative
// the mid-turn-draw determinism test weighs against CostlyDraw.
type CostlyAttack struct{}

func (CostlyAttack) ID() card.ID                                  { return card.FakeCostlyAttack }
func (CostlyAttack) Name() string                                 { return "cardtest.CostlyAttack" }
func (CostlyAttack) Cost(*card.TurnState) int                     { return 1 }
func (CostlyAttack) Pitch() int                                   { return 1 }
func (CostlyAttack) Attack() int                                  { return 3 }
func (CostlyAttack) Defense() int                                 { return 0 }
func (CostlyAttack) Types() card.TypeSet                          { return genericAttackTypes }
func (CostlyAttack) GoAgain() bool                                { return false }
func (CostlyAttack) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }

var genericDefenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)

// PitchOneDR is a 1-pitch-value defense reaction: cost 0, defense 3. It exists solely so tests
// can pitch it (contributing 1 resource) to fund another 1-cost card without also playing it.
type PitchOneDR struct{}

func (PitchOneDR) ID() card.ID              { return card.FakePitchOneDR }
func (PitchOneDR) Name() string             { return "cardtest.PitchOneDR" }
func (PitchOneDR) Cost(*card.TurnState) int { return 0 }
func (PitchOneDR) Pitch() int               { return 1 }
func (PitchOneDR) Attack() int              { return 0 }
func (PitchOneDR) Defense() int             { return 3 }
func (PitchOneDR) Types() card.TypeSet      { return genericDefenseReactionTypes }
func (PitchOneDR) GoAgain() bool            { return false }
func (PitchOneDR) Play(s *card.TurnState, self *card.CardState) {
	s.ApplyAndLogEffectiveDefense(self)
}

// HugeAttack is a 0-cost "do one million damage" attack. Outrageous on purpose: as the top of
// the deck it makes the CostlyDraw → HugeAttack chain blatantly better than the CostlyAttack
// line, so an evaluator that peeks at the deck will pick a different role for the hand's draw
// card depending on deck order — the determinism test catches that.
type HugeAttack struct{}

const hugeAttackDamage = 1_000_000

func (HugeAttack) ID() card.ID                                  { return card.FakeHugeAttack }
func (HugeAttack) Name() string                                 { return "cardtest.HugeAttack" }
func (HugeAttack) Cost(*card.TurnState) int                     { return 0 }
func (HugeAttack) Pitch() int                                   { return 1 }
func (HugeAttack) Attack() int                                  { return hugeAttackDamage }
func (HugeAttack) Defense() int                                 { return 0 }
func (HugeAttack) Types() card.TypeSet                          { return genericAttackTypes }
func (HugeAttack) GoAgain() bool                                { return false }
func (HugeAttack) Play(s *card.TurnState, self *card.CardState) { s.ApplyAndLogEffectiveAttack(self) }
