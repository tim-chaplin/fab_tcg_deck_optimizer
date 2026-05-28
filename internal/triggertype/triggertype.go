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
	// CardOrAbility fires once as a card or weapon attack is played during the attack turn,
	// before that card's own effect resolves. Subscribers narrow with a typeFilter (e.g.
	// IsAttack for Runechant tokens, IsAttackAction for Malefic Incantation).
	CardOrAbility
	// EndOfTurn fires after the attack turn finishes resolving, before the carry snapshot.
	EndOfTurn
	// Hit fires when an attack hits (post-AR-buff EffectiveAttack survives blocks).
	Hit
	// DamageAboutToBeTaken fires just before DamageTaken, while the unblocked-damage figure
	// is still mutable. Subscribers (Enchanting Melody's "instead destroy and prevent 4
	// damage" rider, Talisman of Dousing's Spellvoid) call PreventGenericDamage /
	// PreventPhysicalDamage / PreventArcaneDamage to absorb part of the swing before the
	// carry-through hits DamageTaken triggers. The engine fires this trigger twice per
	// damage moment: once for the physical side then once for the arcane side. Each fire
	// is gated on the matching damage figure being non-zero — type-specific subscribers
	// (Dousing) call the typed prevention method; type-agnostic subscribers (Melody) call
	// PreventGenericDamage and the engine picks the active type.
	DamageAboutToBeTaken
	// DamageTaken fires at the end of the defense phase when incoming damage got through
	// unblocked. Like DamageAboutToBeTaken, it fires twice per damage moment — physical
	// then arcane — so DamageTaken-subscribed self-destruct auras (Arcane Cussing) pop on
	// either side.
	DamageTaken
	// Pitch fires as each card is pitched to fund a play during the attack phase. The
	// triggering card is the pitched card; a handler reads its Pitch value and may boost
	// the resources it yields via GameEngine.AddResourcePoints.
	Pitch
	// CrowdCheer fires when the crowd cheers your hero, raised by GameState.CrowdCheer.
	// Used by "whenever the crowd cheers you" handlers; "if you've been cheered this turn"
	// gates instead read GameState.HasCrowdCheered.
	CrowdCheer
	// CrowdBoo is the boo-side counterpart to CrowdCheer.
	CrowdBoo
)
