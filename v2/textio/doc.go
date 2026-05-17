// Package textio holds the durable on-disk text encodings the optimizer reads and writes.
//
// Two formats live here today:
//
//   - JSON deck-plus-stats (deck_*.go) is the canonical persistence format mydecks/*.json
//     uses. MarshalDeck encodes a *deck.Deck + deck.Stats; UnmarshalDeck reverses it.
//     Round-trips the Best turn's structured Log so saved decks render the same chain
//     attribution as live runs.
//
//   - Fabrary plain-text (fabrary_*.go) converts between *deck.Deck and the format
//     fabrary.net's import / export uses (Name / Hero / Format header, Arena cards
//     section, Deck cards section with pitch-color suffixes, optional Sideboard).
//     MarshalFabrary writes the text; UnmarshalFabrary parses it. Stats aren't part of
//     the fabrary format — those round-trip through the JSON encoding above.
//
// Both encodings rely on the v2/registry name-keyed lookups for card / weapon / hero
// resolution; both are pure data formatters with no sim or gameengine dependency, so the
// package can be imported by any consumer that needs to read or write a deck file.
package textio
