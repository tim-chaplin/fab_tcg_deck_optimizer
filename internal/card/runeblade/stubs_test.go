package runeblade

import "github.com/tim-chaplin/fab-deck-optimizer/internal/card"

// Shared stub Cards used across tests in this package. Each is a zero-value struct with a fixed
// type line that tests mix and match to exercise card effects' lookahead / predicate logic.

// stubRunebladeAttack is a minimal Runeblade Action-Attack card — satisfies lookaheads that look
// for "next Runeblade attack action card" (Runic Reaping, Mauvrion Skies).
type stubRunebladeAttack struct{}

func (stubRunebladeAttack) ID() card.ID  { return card.Invalid }
func (stubRunebladeAttack) Name() string { return "StubRunebladeAttack" }
func (stubRunebladeAttack) Cost() int    { return 0 }
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
func (stubRunebladeWeapon) Cost() int    { return 0 }
func (stubRunebladeWeapon) Pitch() int   { return 0 }
func (stubRunebladeWeapon) Attack() int  { return 0 }
func (stubRunebladeWeapon) Defense() int { return 0 }
func (stubRunebladeWeapon) Types() card.TypeSet {
	return card.NewTypeSet(card.TypeRuneblade, card.TypeWeapon)
}
func (stubRunebladeWeapon) GoAgain() bool            { return false }
func (stubRunebladeWeapon) Play(*card.TurnState) int { return 0 }

// stubNonAttack is a non-attack card — used to confirm attack-typed predicates (Runic Reaping's
// pitched-attack +1{p} rider) do NOT fire on non-attack cards.
type stubNonAttack struct{}

func (stubNonAttack) ID() card.ID              { return card.Invalid }
func (stubNonAttack) Name() string             { return "StubNonAttack" }
func (stubNonAttack) Cost() int                { return 0 }
func (stubNonAttack) Pitch() int               { return 0 }
func (stubNonAttack) Attack() int              { return 0 }
func (stubNonAttack) Defense() int             { return 0 }
func (stubNonAttack) Types() card.TypeSet       { return card.NewTypeSet(card.TypeAction) }
func (stubNonAttack) GoAgain() bool            { return false }
func (stubNonAttack) Play(*card.TurnState) int { return 0 }

// stubAura is a minimal Aura-typed card — used to test "aura played this turn" checks (Shrill of
// Skullform's +3 bonus).
type stubAura struct{}

func (stubAura) ID() card.ID              { return card.Invalid }
func (stubAura) Name() string             { return "StubAura" }
func (stubAura) Cost() int                { return 0 }
func (stubAura) Pitch() int               { return 0 }
func (stubAura) Attack() int              { return 0 }
func (stubAura) Defense() int             { return 0 }
func (stubAura) Types() card.TypeSet       { return card.NewTypeSet(card.TypeAura) }
func (stubAura) GoAgain() bool            { return true }
func (stubAura) Play(*card.TurnState) int { return 0 }
