package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Shared stub Cards. Each is a zero-value struct with a fixed type line; tests mix and match to
// exercise lookahead / predicate logic on card effects.

// stubRunebladeAttack is a minimal Runeblade Action-Attack card — satisfies lookaheads that look
// for "next Runeblade attack action card" (Runic Reaping, Mauvrion Skies).
type stubRunebladeAttack struct{}

func (stubRunebladeAttack) ID() card.ID  { return card.Invalid }
func (stubRunebladeAttack) Name() string { return "StubRunebladeAttack" }
func (stubRunebladeAttack) Cost(*card.TurnState) int    { return 0 }
func (stubRunebladeAttack) Pitch() int   { return 0 }
func (stubRunebladeAttack) Attack() int  { return 0 }
func (stubRunebladeAttack) Defense() int { return 0 }
func (stubRunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (stubRunebladeAttack) GoAgain() bool            { return true }
func (stubRunebladeAttack) Play(*card.TurnState) int { return 0 }

// stubRunebladeWeapon is a Runeblade weapon — satisfies "next Runeblade attack" lookaheads that
// include weapons (Condemn, Oath) but NOT ones restricted to attack action cards (Runic Reaping,
// Mauvrion Skies).
type stubRunebladeWeapon struct{}

func (stubRunebladeWeapon) ID() card.ID  { return card.Invalid }
func (stubRunebladeWeapon) Name() string { return "StubRunebladeWeapon" }
func (stubRunebladeWeapon) Cost(*card.TurnState) int    { return 0 }
func (stubRunebladeWeapon) Pitch() int   { return 0 }
func (stubRunebladeWeapon) Attack() int  { return 0 }
func (stubRunebladeWeapon) Defense() int { return 0 }
func (stubRunebladeWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (stubRunebladeWeapon) GoAgain() bool            { return false }
func (stubRunebladeWeapon) Play(*card.TurnState) int { return 0 }

// stubNonAttack is a non-attack card — covers "attack-typed predicate should reject non-attack"
// cases (e.g. Runic Reaping's pitched-attack +1{p} rider).
type stubNonAttack struct{}

func (stubNonAttack) ID() card.ID              { return card.Invalid }
func (stubNonAttack) Name() string             { return "StubNonAttack" }
func (stubNonAttack) Cost(*card.TurnState) int                { return 0 }
func (stubNonAttack) Pitch() int               { return 0 }
func (stubNonAttack) Attack() int              { return 0 }
func (stubNonAttack) Defense() int             { return 0 }
func (stubNonAttack) Types() card.TypeSet       { return card.NewTypeSet(card.TypeAction) }
func (stubNonAttack) GoAgain() bool            { return false }
func (stubNonAttack) Play(*card.TurnState) int { return 0 }

// stubNonRunebladeAttack is a Generic Action-Attack — covers Runeblade-gated lookaheads
// (Condemn, Oath, Runic Reaping, Mauvrion) rejecting non-Runeblade attacks.
type stubNonRunebladeAttack struct{}

func (stubNonRunebladeAttack) ID() card.ID  { return card.Invalid }
func (stubNonRunebladeAttack) Name() string { return "StubNonRunebladeAttack" }
func (stubNonRunebladeAttack) Cost(*card.TurnState) int    { return 0 }
func (stubNonRunebladeAttack) Pitch() int   { return 0 }
func (stubNonRunebladeAttack) Attack() int  { return 0 }
func (stubNonRunebladeAttack) Defense() int { return 0 }
func (stubNonRunebladeAttack) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeGeneric, card.TypeAction, card.TypeAttack)
}
func (stubNonRunebladeAttack) GoAgain() bool            { return true }
func (stubNonRunebladeAttack) Play(*card.TurnState) int { return 0 }

// stubAttackWithPower is a Runeblade attack-action card with a configurable printed Attack()
// value. Fragile-aura tests set specific numbers to hit/miss the LikelyToHit heuristic (e.g.
// Attack=4 lands, Attack=3 is blockable).
type stubAttackWithPower struct {
	power int
}

func (stubAttackWithPower) ID() card.ID          { return card.Invalid }
func (stubAttackWithPower) Name() string         { return "StubAttackWithPower" }
func (stubAttackWithPower) Cost(*card.TurnState) int            { return 0 }
func (stubAttackWithPower) Pitch() int           { return 0 }
func (s stubAttackWithPower) Attack() int        { return s.power }
func (stubAttackWithPower) Defense() int         { return 0 }
func (stubAttackWithPower) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeAction, card.TypeAttack)
}
func (stubAttackWithPower) GoAgain() bool            { return true }
func (stubAttackWithPower) Play(*card.TurnState) int { return 0 }

// stubAura is a minimal Aura-typed card — exercises "aura played this turn" checks (Shrill of
// Skullform's +3 bonus).
type stubAura struct{}

func (stubAura) ID() card.ID              { return card.Invalid }
func (stubAura) Name() string             { return "StubAura" }
func (stubAura) Cost(*card.TurnState) int                { return 0 }
func (stubAura) Pitch() int               { return 0 }
func (stubAura) Attack() int              { return 0 }
func (stubAura) Defense() int             { return 0 }
func (stubAura) Types() card.TypeSet       { return card.NewTypeSet(card.TypeAura) }
func (stubAura) GoAgain() bool            { return true }
func (stubAura) Play(*card.TurnState) int { return 0 }
