// Package triggertype categorises when an Aura or one-shot Trigger fires. Lives in its
// own micro-package so consumers can name the enum without pulling in the whole engine.
package triggertype

// Type identifies the lifecycle event an Aura / Trigger subscribes to.
type Type int

const (
	// StartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search.
	StartOfTurn Type = iota
	// CardOrAbility fires once as a card or weapon attack is played during the chain,
	// before that card's own effect resolves. Subscribers narrow with a typeFilter (e.g.
	// IsAttack for Runechant tokens, IsAttackAction for Malefic Incantation).
	CardOrAbility
	// EndOfTurn fires after the chain finishes resolving, before the carry snapshot.
	EndOfTurn
	// Hit fires when an attack hits (post-AR-buff EffectiveAttack survives blocks).
	Hit
)
