package cards

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// stubCard is a configurable Card implementation used across generic tests to build
// CardsRemaining / CardsPlayed / Pitched lists with specific type, cost, power, and pitch shapes.
// Zero-value fields mean "don't care" — tests set only what the helper under test predicates on.
type stubCard struct {
	name  string
	cost  int
	power int
	pitch int
	types card.TypeSet
}

func (s stubCard) ID() card.ID                         { return card.Invalid }
func (s stubCard) Name() string                        { return s.name }
func (s stubCard) Cost(*card.TurnState) int            { return s.cost }
func (s stubCard) Pitch() int                          { return s.pitch }
func (s stubCard) Attack() int                         { return s.power }
func (s stubCard) Defense() int                        { return 0 }
func (s stubCard) Types() card.TypeSet                 { return s.types }
func (s stubCard) GoAgain() bool                       { return false }
func (stubCard) Play(*card.TurnState, *card.CardState) {}

// stubGenericAttack returns a Generic Action - Attack stub with the given cost and base power.
// Pitch defaults to 1; override via the pitch field if a test cares.
func stubGenericAttack(cost, power int) stubCard {
	return stubCard{
		name:  "stubGenericAttack",
		cost:  cost,
		power: power,
		pitch: 1,
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack),
	}
}

// stubGenericAttackPitch is stubGenericAttack with an explicit pitch value. Flying High's red
// variant rider reads pitch, so tests that exercise the +1 bonus set this.
func stubGenericAttackPitch(cost, power, pitch int) stubCard {
	s := stubGenericAttack(cost, power)
	s.pitch = pitch
	return s
}

// stubGenericAction returns a Generic Action (non-attack) stub for attack-typed-lookahead
// rejection cases.
func stubGenericAction() stubCard {
	return stubCard{
		name:  "stubGenericAction",
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAction),
	}
}

// stubGenericAura returns a Generic Aura stub — covers Yinti Yanti's HasPlayedType(TypeAura) check.
func stubGenericAura() stubCard {
	return stubCard{
		name:  "stubGenericAura",
		types: card.NewTypeSet(card.TypeGeneric, card.TypeAura),
	}
}

// Shared stub Cards. Each is a zero-value struct with a fixed type line; tests mix and match to
// exercise lookahead / predicate logic on card effects.

// stubRunebladeAttack is a minimal Runeblade Action-Attack card — satisfies "next Runeblade
// attack action card" lookaheads.
type stubRunebladeAttack struct{}

func (stubRunebladeAttack) ID() card.ID              { return card.Invalid }
func (stubRunebladeAttack) Name() string             { return "StubRunebladeAttack" }
func (stubRunebladeAttack) Cost(*card.TurnState) int { return 0 }
func (stubRunebladeAttack) Pitch() int               { return 0 }
func (stubRunebladeAttack) Attack() int              { return 0 }
func (stubRunebladeAttack) Defense() int             { return 0 }
func (stubRunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (stubRunebladeAttack) GoAgain() bool                         { return true }
func (stubRunebladeAttack) Play(*card.TurnState, *card.CardState) {}

// stubRunebladeWeapon is a Runeblade weapon — satisfies "next Runeblade attack" lookaheads
// that include weapons but NOT ones restricted to attack action cards.
type stubRunebladeWeapon struct{}

func (stubRunebladeWeapon) ID() card.ID              { return card.Invalid }
func (stubRunebladeWeapon) Name() string             { return "StubRunebladeWeapon" }
func (stubRunebladeWeapon) Cost(*card.TurnState) int { return 0 }
func (stubRunebladeWeapon) Pitch() int               { return 0 }
func (stubRunebladeWeapon) Attack() int              { return 0 }
func (stubRunebladeWeapon) Defense() int             { return 0 }
func (stubRunebladeWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (stubRunebladeWeapon) GoAgain() bool                         { return false }
func (stubRunebladeWeapon) Play(*card.TurnState, *card.CardState) {}

// stubNonAttack is a non-attack card — covers "attack-typed predicate should reject
// non-attack" cases.
type stubNonAttack struct{}

func (stubNonAttack) ID() card.ID                           { return card.Invalid }
func (stubNonAttack) Name() string                          { return "StubNonAttack" }
func (stubNonAttack) Cost(*card.TurnState) int              { return 0 }
func (stubNonAttack) Pitch() int                            { return 0 }
func (stubNonAttack) Attack() int                           { return 0 }
func (stubNonAttack) Defense() int                          { return 0 }
func (stubNonAttack) Types() card.TypeSet                   { return card.NewTypeSet(card.TypeAction) }
func (stubNonAttack) GoAgain() bool                         { return false }
func (stubNonAttack) Play(*card.TurnState, *card.CardState) {}

// stubNonRunebladeAttack is a Generic Action-Attack — covers Runeblade-gated lookaheads
// rejecting non-Runeblade attacks.
type stubNonRunebladeAttack struct{}

func (stubNonRunebladeAttack) ID() card.ID              { return card.Invalid }
func (stubNonRunebladeAttack) Name() string             { return "StubNonRunebladeAttack" }
func (stubNonRunebladeAttack) Cost(*card.TurnState) int { return 0 }
func (stubNonRunebladeAttack) Pitch() int               { return 0 }
func (stubNonRunebladeAttack) Attack() int              { return 0 }
func (stubNonRunebladeAttack) Defense() int             { return 0 }
func (stubNonRunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubNonRunebladeAttack) GoAgain() bool                         { return true }
func (stubNonRunebladeAttack) Play(*card.TurnState, *card.CardState) {}

// stubAttackWithPower is a Runeblade attack-action card with a configurable printed Attack()
// value. Tests set specific numbers to hit/miss the LikelyToHit heuristic (4 lands, 3 blocks).
type stubAttackWithPower struct {
	power int
}

func (stubAttackWithPower) ID() card.ID              { return card.Invalid }
func (stubAttackWithPower) Name() string             { return "StubAttackWithPower" }
func (stubAttackWithPower) Cost(*card.TurnState) int { return 0 }
func (stubAttackWithPower) Pitch() int               { return 0 }
func (s stubAttackWithPower) Attack() int            { return s.power }
func (stubAttackWithPower) Defense() int             { return 0 }
func (stubAttackWithPower) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (stubAttackWithPower) GoAgain() bool                         { return true }
func (stubAttackWithPower) Play(*card.TurnState, *card.CardState) {}

// stubAura is a minimal Aura-typed card — exercises "aura played this turn" checks.
type stubAura struct{}

func (stubAura) ID() card.ID                           { return card.Invalid }
func (stubAura) Name() string                          { return "StubAura" }
func (stubAura) Cost(*card.TurnState) int              { return 0 }
func (stubAura) Pitch() int                            { return 0 }
func (stubAura) Attack() int                           { return 0 }
func (stubAura) Defense() int                          { return 0 }
func (stubAura) Types() card.TypeSet                   { return card.NewTypeSet(card.TypeAura) }
func (stubAura) GoAgain() bool                         { return true }
func (stubAura) Play(*card.TurnState, *card.CardState) {}