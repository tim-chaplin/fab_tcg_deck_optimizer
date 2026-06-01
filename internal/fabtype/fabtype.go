// Package fabtype holds the Flesh and Blood class, talent, and non-deck card-type NAMES as
// they appear in printed type lines and the upstream card.csv Types column. It's the text-side
// vocabulary the CSV / typebox tools (cmd/parsecarddb, cmd/crosscheck) use to classify a type
// word as a class, a talent, or a non-deck card type.
//
// This is deliberately separate from the in-engine legality check: internal/registry decides
// what a hero can play with classMask / talentMask (TypeSet bits over the implemented
// card.CardType constants). This package is the broader name list covering classes and talents
// we haven't implemented cards for yet. Keep the two in sync when a new class or talent ships.
package fabtype

// ClassNames is every hero class word that can head a type line. Generic is a class here too
// (it's the "no class" pseudo-class every hero can play).
var ClassNames = map[string]bool{
	"Adjudicator": true, "Assassin": true, "Bard": true, "Brute": true, "Generic": true,
	"Guardian": true, "Illusionist": true, "Mechanologist": true, "Merchant": true,
	"Monarch": true, "Necromancer": true, "Ninja": true, "Pirate": true, "Ranger": true,
	"Runeblade": true, "Shapeshifter": true, "Warrior": true, "Wizard": true,
}

// ClassMatches reports whether a type-line word is class-legal for a hero whose classes are
// heroClasses. Generic is legal for every hero, and a word that isn't a class at all (a talent,
// card type, or subtype) imposes no class constraint — both return true. A class word is legal
// only when it's one the hero plays. Callers pass just the hero's own class(es); Generic is
// built in here, so it never needs listing per hero.
func ClassMatches(word string, heroClasses map[string]bool) bool {
	if word == "Generic" || !ClassNames[word] {
		return true
	}
	return heroClasses[word]
}

// TalentNames is every talent word. A card carrying one is illegal for a hero who lacks that
// talent. Royal is intentionally absent — it's a supertype, not a talent (Royal Generics like
// Imperial Seal of Command stay broadly legal).
var TalentNames = map[string]bool{
	"Chaos": true, "Draconic": true, "Earth": true, "Elemental": true,
	"Ice": true, "Light": true, "Lightning": true, "Shadow": true,
}

// NonDeckTypes is every printed card type that never enters the 40-card deck — equipment slots,
// weapons, the hero, tokens, landmarks, and so on.
var NonDeckTypes = map[string]bool{
	"Weapon": true, "Equipment": true, "Hero": true, "Demi-Hero": true,
	"Token": true, "Landmark": true, "Mentor": true, "Ally": true, "Macro": true,
}

// UnmodeledTypes is NonDeckTypes minus Weapon: the printed card types the optimizer doesn't
// simulate. Weapons are excluded because the deck model *does* carry a weapon slot
// (registry.LegalWeaponsFor) — a Lightning weapon is implementable, a Lightning equipment isn't.
var UnmodeledTypes = map[string]bool{
	"Equipment": true, "Hero": true, "Demi-Hero": true,
	"Token": true, "Landmark": true, "Mentor": true, "Ally": true, "Macro": true,
}
