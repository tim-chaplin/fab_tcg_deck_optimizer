// Package triggertype categorises when an Aura or one-shot Trigger fires. Lives in its
// own micro-package so consumers can name the enum without pulling in the whole engine.
package triggertype

// Type identifies the lifecycle event an Aura / Trigger subscribes to.
type Type int

const (
	// StartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search.
	StartOfTurn Type = iota
	// AttackAction fires each time an attack action card resolves during the chain.
	AttackAction
	// Attack fires when ANY attack resolves — attack action card or weapon swing.
	Attack
	// EndOfTurn fires after the chain finishes resolving, before the carry snapshot.
	EndOfTurn
	// Hit fires when an attack hits (post-AR-buff EffectiveAttack survives blocks).
	Hit
)
