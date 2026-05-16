package aura

import "github.com/tim-chaplin/fab-deck-optimizer/internal/registry/ids"

// Card is the minimal source-card surface a card-backed aura needs (display name + ID).
// The richer v2/card.Card satisfies it structurally so consumers can pass card values
// through without conversion.
type Card interface {
	DisplayName() string
	ID() ids.CardID
}

// GameEngine is the narrow engine surface the firing aura's context wiring needs — just
// the Destroy hop. *gameengine.GameEngine satisfies it structurally. Handler bodies, which
// receive the engine as `any`, assert to their own (typically richer) engine interface.
type GameEngine interface {
	DestroyAura(addToGraveyard bool)
}
