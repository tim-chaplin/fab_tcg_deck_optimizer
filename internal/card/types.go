package card

// Card-type keywords and the TypeSet bitfield that pools them. Every type-line check on the
// hot path (attack action? defense reaction? persists in play?) routes through a single
// bitmask AND, so the solver's tight inner loops avoid string comparison / map lookup.

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
	TypeOneHand                              // "1H"
	TypeRuneblade                            // "Runeblade"
	TypeScepter                              // "Scepter"
	TypeSword                                // "Sword"
	TypeTwoHand                              // "2H"
	TypeWeapon                               // "Weapon"
	TypeYoung                                // "Young"
)

// persistsInPlayMask is the set of types that keep a card in its zone after resolving rather
// than heading to the graveyard. Auras (e.g. Sigil of the Arknight: Runeblade, Action, Aura)
// and Items linger in the arena until a destroy condition fires; weapons stay equipped.
const persistsInPlayMask TypeSet = TypeSet(TypeAura) | TypeSet(TypeItem) | TypeSet(TypeWeapon)

// PersistsInPlay reports whether a card with this type set stays in its zone when it resolves
// instead of heading to the graveyard. Used by the solver to decide whether to append a
// just-played card to state.Graveyard.
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

// IsNonAttackAction reports whether s represents an Action that is not also an Attack —
// the bitmask check behind every "if a non-attack action card was played/pitched" rider.
func (s TypeSet) IsNonAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) == 0
}

// IsAttackAction reports whether s is an attack action card — both Action and Attack.
// Single-expression bitmask keeps the "next attack action" peek loops lean.
func (s TypeSet) IsAttackAction() bool {
	return s&TypeSet(TypeAction) != 0 && s&TypeSet(TypeAttack) != 0
}

// IsAttack reports whether s represents an attack — an attack action card OR a weapon
// swing. Used by riders whose printed text says "your next attack" with no "action card"
// qualifier; weapons are eligible alongside attack action cards in that wording.
func (s TypeSet) IsAttack() bool {
	return s&(TypeSet(TypeAttack)|TypeSet(TypeWeapon)) != 0
}

// IsWeaponAttack reports whether s represents a weapon attack — a card with the Weapon
// type. Used by riders whose printed text says "weapon attack" (e.g. Pummel's "club or
// hammer weapon attack"), which gates a weapon swing only and excludes attack action
// cards that happen to share the weapon's type subtag.
func (s TypeSet) IsWeaponAttack() bool {
	return s&TypeSet(TypeWeapon) != 0
}

// IsRunebladeAttack reports whether s is a Runeblade attack — an attack action card OR a
// weapon swing. Used by "next Runeblade attack this turn" riders that peek CardsRemaining.
func (s TypeSet) IsRunebladeAttack() bool {
	return s&TypeSet(TypeRuneblade) != 0 && s&(TypeSet(TypeAttack)|TypeSet(TypeWeapon)) != 0
}

// IsDefenseReaction reports whether s has the Defense Reaction subtype.
func (s TypeSet) IsDefenseReaction() bool {
	return s&TypeSet(TypeDefenseReaction) != 0
}

// IsAttackReaction reports whether s has the Attack Reaction subtype.
func (s TypeSet) IsAttackReaction() bool {
	return s&TypeSet(TypeAttackReaction) != 0
}
