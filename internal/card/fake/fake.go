// Package fake provides generic stub Card implementations used by tests in multiple packages (hand,
// sim). These are not real FaB cards — they're deliberately simple attack actions with known stat
// lines so partition/ordering tests have predictable optimal values.
package fake

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

var genericAttackTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)

// BlueAttack is a generic blue attack action: pitches 3, defends 3, attacks 1, costs 1.
type BlueAttack struct{}

func (BlueAttack) ID() card.ID                { return card.FakeBlueAttack }
func (BlueAttack) Name() string               { return "cardtest.BlueAttack" }
func (BlueAttack) Cost(*card.TurnState) int                  { return 1 }
func (BlueAttack) Pitch() int                 { return 3 }
func (BlueAttack) Attack() int                { return 1 }
func (BlueAttack) Defense() int               { return 3 }
func (BlueAttack) Types() card.TypeSet        { return genericAttackTypes }
func (BlueAttack) GoAgain() bool              { return true }
func (c BlueAttack) Play(*card.TurnState, *card.CardState) int { return c.Attack() }

// RedAttack is a generic red attack action: pitches 1, defends 1, attacks 3, costs 1.
type RedAttack struct{}

func (RedAttack) ID() card.ID                { return card.FakeRedAttack }
func (RedAttack) Name() string               { return "cardtest.RedAttack" }
func (RedAttack) Cost(*card.TurnState) int                  { return 1 }
func (RedAttack) Pitch() int                 { return 1 }
func (RedAttack) Attack() int                { return 3 }
func (RedAttack) Defense() int               { return 1 }
func (RedAttack) Types() card.TypeSet        { return genericAttackTypes }
func (RedAttack) GoAgain() bool              { return true }
func (c RedAttack) Play(*card.TurnState, *card.CardState) int { return c.Attack() }

// YellowAttack is a generic yellow attack action: pitches 2, defends 2, attacks 2, costs 1.
type YellowAttack struct{}

func (YellowAttack) ID() card.ID                { return card.FakeYellowAttack }
func (YellowAttack) Name() string               { return "cardtest.YellowAttack" }
func (YellowAttack) Cost(*card.TurnState) int                  { return 1 }
func (YellowAttack) Pitch() int                 { return 2 }
func (YellowAttack) Attack() int                { return 2 }
func (YellowAttack) Defense() int               { return 2 }
func (YellowAttack) Types() card.TypeSet        { return genericAttackTypes }
func (YellowAttack) GoAgain() bool              { return true }
func (c YellowAttack) Play(*card.TurnState, *card.CardState) int { return c.Attack() }

// DrawCantrip is a generic free-cycling attack: cost 0, pitches 1, attacks 1, go again, and
// fires DrawOne on play. Used by tests to exercise mid-turn-draw chains that extend themselves
// (each cantrip plays, draws the next one, which plays, etc.).
type DrawCantrip struct{}

func (DrawCantrip) ID() card.ID                { return card.FakeDrawCantrip }
func (DrawCantrip) Name() string               { return "cardtest.DrawCantrip" }
func (DrawCantrip) Cost(*card.TurnState) int   { return 0 }
func (DrawCantrip) Pitch() int                 { return 1 }
func (DrawCantrip) Attack() int                { return 1 }
func (DrawCantrip) Defense() int               { return 0 }
func (DrawCantrip) Types() card.TypeSet        { return genericAttackTypes }
func (DrawCantrip) GoAgain() bool              { return true }
func (c DrawCantrip) Play(s *card.TurnState, _ *card.CardState) int {
	s.DrawOne()
	return c.Attack()
}

// genericActionTypes is a plain non-attack action (no Attack subtype). Used by CostlyDraw — a
// draw-a-card action card. It isn't a Defense Reaction so it can't be played on the opponent's
// turn; it carries Go again so a drawn-card-as-extension can chain.
var genericActionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeAction)

// CostlyDraw is a 1-cost, pitch-1, no-damage "draw a card, go again" action. Used by the
// mid-turn-draw determinism tests: plays for 0 damage but a chain can continue off its Go again
// onto a drawn-later attack. Implements NoMemo because Play reads the deck via DrawOne.
type CostlyDraw struct{}

func (CostlyDraw) ID() card.ID                { return card.FakeCostlyDraw }
func (CostlyDraw) Name() string               { return "cardtest.CostlyDraw" }
func (CostlyDraw) Cost(*card.TurnState) int   { return 1 }
func (CostlyDraw) Pitch() int                 { return 1 }
func (CostlyDraw) Attack() int                { return 0 }
func (CostlyDraw) Defense() int               { return 0 }
func (CostlyDraw) Types() card.TypeSet        { return genericActionTypes }
func (CostlyDraw) GoAgain() bool              { return true }
func (CostlyDraw) NoMemo()                    {}
func (CostlyDraw) Play(s *card.TurnState, _ *card.CardState) int { s.DrawOne(); return 0 }

// CostlyAttack is a 1-cost, pitch-1, 3-damage attack action — the "deal 3 damage" alternative
// the mid-turn-draw determinism test weighs against CostlyDraw.
type CostlyAttack struct{}

func (CostlyAttack) ID() card.ID                 { return card.FakeCostlyAttack }
func (CostlyAttack) Name() string                { return "cardtest.CostlyAttack" }
func (CostlyAttack) Cost(*card.TurnState) int    { return 1 }
func (CostlyAttack) Pitch() int                  { return 1 }
func (CostlyAttack) Attack() int                 { return 3 }
func (CostlyAttack) Defense() int                { return 0 }
func (CostlyAttack) Types() card.TypeSet         { return genericAttackTypes }
func (CostlyAttack) GoAgain() bool               { return false }
func (c CostlyAttack) Play(*card.TurnState, *card.CardState) int  { return c.Attack() }

var genericDefenseReactionTypes = card.NewTypeSet(card.TypeGeneric, card.TypeDefenseReaction)

// PitchOneDR is a 1-pitch-value defense reaction: cost 0, defense 3. It exists solely so tests
// can pitch it (contributing 1 resource) to fund another 1-cost card without also playing it.
type PitchOneDR struct{}

func (PitchOneDR) ID() card.ID                { return card.FakePitchOneDR }
func (PitchOneDR) Name() string               { return "cardtest.PitchOneDR" }
func (PitchOneDR) Cost(*card.TurnState) int   { return 0 }
func (PitchOneDR) Pitch() int                 { return 1 }
func (PitchOneDR) Attack() int                { return 0 }
func (PitchOneDR) Defense() int               { return 3 }
func (PitchOneDR) Types() card.TypeSet        { return genericDefenseReactionTypes }
func (PitchOneDR) GoAgain() bool              { return false }
func (PitchOneDR) Play(*card.TurnState, *card.CardState) int   { return 0 }

// HugeAttack is a 0-cost "do one million damage" attack. Outrageous on purpose: as the top of
// the deck it makes the CostlyDraw → HugeAttack chain blatantly better than the CostlyAttack
// line, so an evaluator that peeks at the deck will pick a different role for the hand's draw
// card depending on deck order — the determinism test catches that.
type HugeAttack struct{}

const hugeAttackDamage = 1_000_000

func (HugeAttack) ID() card.ID                { return card.FakeHugeAttack }
func (HugeAttack) Name() string               { return "cardtest.HugeAttack" }
func (HugeAttack) Cost(*card.TurnState) int   { return 0 }
func (HugeAttack) Pitch() int                 { return 1 }
func (HugeAttack) Attack() int                { return hugeAttackDamage }
func (HugeAttack) Defense() int               { return 0 }
func (HugeAttack) Types() card.TypeSet        { return genericAttackTypes }
func (HugeAttack) GoAgain() bool              { return false }
func (c HugeAttack) Play(*card.TurnState, *card.CardState) int { return c.Attack() }
