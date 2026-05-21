// Package-level go:generate directive for the cards package. `go generate` regenerates
// every <card>_gen.go file from its <card>.yaml and rewrites internal/registry's cardsByID
// map so a new card is registered for deck construction automatically. Hand-written
// <card>.go files (Play, riders, modal helpers, variable Cost, AR predicates, …) are
// untouched.
//
// New cards: drop a <card>.yaml, run `go generate`, then add the hand-written companion
// file supplying Play().

//go:generate go run ../../../cmd/cardgen -registry ../../registry/cards_gen.go .

package cards
