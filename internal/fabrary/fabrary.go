// Package fabrary converts a sim.Deck to and from fabrary.net's plain-text deck format
// (https://fabrary.net/decks?tab=import). The format has a `Name:` / `Hero:` / `Format:` header,
// an "Arena cards" section for equipment and weapons, a "Deck cards" section with pitch cards
// carrying a lowercase color suffix (e.g. "2x Aether Slash (red)"), and an optional
// "Sideboard" section mirroring the Deck section for the user-managed sideboard.
//
// The optimizer models only weapons, not other equipment. Unknown Arena lines are ignored on
// import; on export, modelled weapons are joined by the fixed equipment loadout in
// defaultArenaPackage so the emitted .txt can be pasted into fabrary without hand-editing.
//
// The encoding splits across sibling files in this package: marshal.go (runtime → fabrary
// text), unmarshal.go (fabrary text → runtime), names.go (pitch-suffix case conversion
// shared by both directions).
package fabrary
