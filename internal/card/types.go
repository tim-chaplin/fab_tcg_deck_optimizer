package card

// Card-type keywords and the TypeSet bitfield that pools them. Every type-line check
// routes through a single bitmask AND — no string comparison or map lookup on the hot path.

// CardType is a card-type descriptor. Each constant corresponds to one keyword from a FaB
// card's type line (e.g. "Runeblade", "Action", "Attack").
type CardType uint64

const (
	TypeAction          CardType = 1 << iota // "Action"
	TypeAttack                               // "Attack"
	TypeAttackReaction                       // "Attack Reaction"
	TypeAura                                 // "Aura"
	TypeBlock                                // "Block"
	TypeBook                                 // "Book"
	TypeClub                                 // "Club"
	TypeDefenseReaction                      // "Defense Reaction"
	TypeGeneric                              // "Generic"
	TypeHammer                               // "Hammer"
	TypeHero                                 // "Hero"
	TypeInstant                              // "Instant"
	TypeItem                                 // "Item"
	TypeLightning                            // "Lightning" — a talent (Aurora et al.)
	TypeOneHand                              // "1H"
	TypeResource                             // "Resource"
	TypeRevered                              // "Revered" — hero-only; Rosetta crowd-cheer keyword
	TypeReviled                              // "Reviled" — hero-only; Rosetta crowd-boo keyword
	TypeRoyal                                // "Royal" — hero-only; Imperial Seal of Command gate
	TypeRuneblade                            // "Runeblade"
	TypeScepter                              // "Scepter"
	TypeSword                                // "Sword"
	TypeThief                                // "Thief"
	TypeTwoHand                              // "2H"
	TypeWeapon                               // "Weapon"
	TypeYoung                                // "Young"
)

// persistsInPlayMask is the set of types that keep a card in its zone after resolving
// rather than heading to the graveyard. Auras and Items linger in the arena until a
// destroy condition fires; weapons stay equipped.
const persistsInPlayMask TypeSet = TypeSet(TypeAura) | TypeSet(TypeItem) | TypeSet(TypeWeapon)

// PersistsInPlay reports whether a card with this type set stays in its zone on resolve
// instead of heading to the graveyard.
func (s TypeSet) PersistsInPlay() bool {
	return s&persistsInPlayMask != 0
}

// TypeSet is a bitfield of CardType values — type checks become a single-word bitmask AND, no
// string hashing or map lookup on the hot path.
type TypeSet uint64

// NewTypeSet returns a TypeSet containing all of the given types.
func NewTypeSet(types ...CardType) TypeSet {
	var s TypeSet
	for _, t := range types {
		s |= TypeSet(t)
	}
	return s
}

// Has reports whether s contains the given type.
func (s TypeSet) Has(t CardType) bool { return s&TypeSet(t) != 0 }

// IsAction reports whether s represents an action card — TypeAction present, covering
// both attack actions and non-attack actions. Backs "action card" riders.
func (s TypeSet) IsAction() bool {
	return s&TypeSet(TypeAction) != 0
}

// IsNonAttackAction reports whether s represents an Action that is not also an Attack —
// the bitmask check behind every "if a non-attack action card was played/pitched" rider.
func (s TypeSet) IsNonAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) == 0
}

// IsAttackAction reports whether s is an attack action card — both Action and Attack.
func (s TypeSet) IsAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) != 0
}

// IsAttack reports whether s represents an attack — an attack action card or a weapon
// swing. Backs riders whose printed text says "your next attack" with no "action card"
// qualifier.
func (s TypeSet) IsAttack() bool {
	return s&TypeSet(TypeAttack) != 0
}

// IsWeaponAttack reports whether s represents a weapon attack — both TypeWeapon and
// TypeAttack present. Backs "weapon attack" riders.
func (s TypeSet) IsWeaponAttack() bool {
	return s&TypeSet(TypeWeapon) != 0 && s&TypeSet(TypeAttack) != 0
}

// IsRunebladeAttack reports whether s is a Runeblade attack — TypeRuneblade and
// TypeAttack both present. Backs "next Runeblade attack this turn" riders.
func (s TypeSet) IsRunebladeAttack() bool {
	return s&TypeSet(TypeRuneblade) != 0 && s&TypeSet(TypeAttack) != 0
}

// IsDefenseReaction reports whether s has the Defense Reaction subtype.
func (s TypeSet) IsDefenseReaction() bool {
	return s&TypeSet(TypeDefenseReaction) != 0
}

// IsAttackReaction reports whether s has the Attack Reaction subtype.
func (s TypeSet) IsAttackReaction() bool {
	return s&TypeSet(TypeAttackReaction) != 0
}

// IsResource reports whether s carries the Resource type. Resource cards have no Action
// subtype so they can't be played, but their printed Defense is still a legal block. The
// post-hoc arsenal promotion skips them (alongside pure Block cards) since neither has a
// legal play from arsenal.
func (s TypeSet) IsResource() bool {
	return s&TypeSet(TypeResource) != 0
}

// IsAttack matches any scheduled attack — action card OR weapon swing — for "your next
// attack" wording. Pairs with helpers that take a `func(GameEngine, *CardState) bool`
// predicate. Uses EffectiveTypes so ModalTypes cards (Tip-Off mode 1) report their mode-
// resolved type-line, not the static printed one.
func IsAttack(ge GameEngine, pc *CardState) bool {
	return pc.EffectiveTypes(ge).IsAttack()
}

// IsAttackAction matches scheduled attack action cards (excludes weapon swings) for "the
// next attack action card" wording.
func IsAttackAction(ge GameEngine, pc *CardState) bool {
	return pc.EffectiveTypes(ge).IsAttackAction()
}

// IsRunebladeAttack matches scheduled Runeblade attacks (action or weapon) for "your next
// Runeblade attack" wording. Engine threaded through so Universal cards fold the active
// hero's class into their Types.
func IsRunebladeAttack(ge GameEngine, pc *CardState) bool {
	return pc.EffectiveTypes(ge).IsRunebladeAttack()
}

// IsWeaponAttack matches scheduled weapon-swing attacks — not attack action cards — for
// "your next weapon attack" wording.
func IsWeaponAttack(ge GameEngine, pc *CardState) bool {
	return pc.EffectiveTypes(ge).IsWeaponAttack()
}

// IsSwordAttack matches scheduled Sword attacks — weapon swings or attack cards carrying
// the Sword type — for "your next sword attack" wording.
func IsSwordAttack(ge GameEngine, pc *CardState) bool {
	t := pc.EffectiveTypes(ge)
	return t.Has(TypeAttack) && t.Has(TypeSword)
}

// IsNonAttackAction matches scheduled non-attack action cards — Actions that are not
// Attacks — for "your next non-attack action card" wording.
func IsNonAttackAction(ge GameEngine, pc *CardState) bool {
	return pc.EffectiveTypes(ge).IsNonAttackAction()
}
