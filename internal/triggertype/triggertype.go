// Package triggertype categorises when an Aura or one-shot Trigger fires. Lives in its
// own micro-package so consumers can name the enum without pulling in the whole engine.
package triggertype

// Type is a bitmask of lifecycle events. Each constant is one bit; an Aura or Trigger
// subscribes to one event or an OR of several, and FireTriggers fires it when the
// dispatched event's bit is set.
type Type int

const (
	// StartOfTurn fires at the start of the owning player's action phase, before the
	// best-line search.
	StartOfTurn Type = 1 << iota
	// CardOrAbility fires once as a card or weapon attack is played during the chain,
	// before that card's own effect resolves. Subscribers narrow with a typeFilter (e.g.
	// IsAttack for Runechant tokens, IsAttackAction for Malefic Incantation).
	CardOrAbility
	// EndOfTurn fires after the chain finishes resolving, before the carry snapshot.
	EndOfTurn
	// Hit fires when an attack hits (post-AR-buff EffectiveAttack survives blocks).
	Hit
	// DamageTaken fires at the end of the defense phase when incoming damage got through
	// unblocked.
	DamageTaken
	// Pitch fires as each card is pitched to fund a play during the attack phase. The
	// triggering card is the pitched card; a handler reads its Pitch value and may boost
	// the resources it yields via GameEngine.AddResourcePoints.
	Pitch
)
