package sim

// Item is a permanent in play with an activated ability the player chooses to play during
// their turn. Each turn the chain runner enqueues the Ability() Card as a playable option
// (1 AP, pays the printed activation cost; the Ability's Play decrements Count and removes
// the entry when Count reaches zero).
//
// Invariant: at most one Item per TokenType per TurnState; helpers bump Count on the
// existing entry rather than appending duplicates.
type Item struct {
	// Self identifies the item — a card or a token type.
	Self CardOrTokenType
	// Count is the number of copies / charges in play. The activated ability decrements
	// this each time it consumes one.
	Count int
	// Ability is the activated-ability Card the chain runner enqueues each turn. The
	// ability's Play calls back into TurnState (e.g. ConsumeItem) to decrement Count and
	// destroy this entry when Count reaches zero. Token items don't head to the graveyard
	// on destroy.
	Ability Card
}

// itemCountIn returns the Count of the item entry matching token type t, or zero.
func itemCountIn(items []Item, t TokenType) int {
	for i := range items {
		if items[i].Self.TokenType == t {
			return items[i].Count
		}
	}
	return 0
}

// ItemCreator is an optional Card marker. Cards whose Play / OnHit creates a token item
// (Strike Gold creating Gold, Performance Bonus creating Gold) declare so via this
// interface; bestAttackWithWeapons unions the declared types with priorItems' types and
// adds one ability instance per resulting type to the wmask. Without the marker the
// optimizer can't enumerate a chain that spends an item created mid-turn — the wmask
// is built once at chain start.
//
// CreatesItem returns TokenTypeNone when the card sometimes creates one and sometimes
// doesn't (e.g. on-hit gated). The marker is conservative: declaring true reserves a
// wmask bit that PlayPrecondition rejects when the token never materialises, so a
// blockable Strike Gold variant doesn't pin a useless ability into every permutation.
type ItemCreator interface {
	CreatesItem() TokenType
}

// tokenItemAbilityFor returns the activated-ability Card for token-item type t, or nil
// when t isn't a token-item type. Single read site for the (TokenType, Ability) mapping
// — bestAttackWithWeapons reads it when assembling itemAbilities for ItemCreator
// declarations that aren't already covered by priorItems' Ability field.
func tokenItemAbilityFor(t TokenType) Card {
	switch t {
	case TokenTypeGold:
		return GoldTokenAbility{}
	}
	return nil
}
