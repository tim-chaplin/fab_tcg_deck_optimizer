// Package deckio serializes Deck values (plus their accumulated Stats) to and from JSON.
// Cards, weapons, and heroes are referenced by name; deserialization looks names up in
// package cards and the hero registry.
//
// The concrete work is split across sibling files in this package: types.go (the DeckJSON /
// StatsJSON / BestTurnJSON shapes), marshal.go (runtime → JSON encoding), unmarshal.go
// (JSON → runtime decoding).
package deckio
