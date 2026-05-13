package gameengine

import (
	"github.com/tim-chaplin/fab-deck-optimizer/v2/card"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/deck"
	"github.com/tim-chaplin/fab-deck-optimizer/v2/turnlogger"
)

// Spec is the test-friendly engine constructor input. Mirrors the scalar / card-typed
// fields of GameState in exported form so callers outside the package don't have to
// chain per-field setters after construction. Build via NewFromSpec.
//
// Spec carries no Auras / Triggers / Items list: callers seed those individually after
// NewFromSpec via CreateAura / CreateTrigger / CreateItem on the underlying GameState.
type Spec struct {
	Hero                  Hero
	Arsenal               card.Card
	Banished              []card.Card
	Graveyard             []card.Card
	Pitched               []card.Card
	Defenders             []card.Card
	CardsPlayed           []card.Card
	CardsRemaining        []*card.CardState
	OpponentMarked        bool
	ArcaneDamageDealt     bool
	AuraCreated           bool
	CardBanished          bool
	NonAttackActionPlayed bool
	ActionPoints          int
	IncomingDamage        int
	ArcaneIncomingDamage  int
	BlockTotal            int
	Value                 int
	TriggeringCard        card.Card
	AttackReactionTarget  *card.CardState
}

// NewFromState wraps s in a *GameEngine. Pass an existing *GameState to drive a
// chain-runner state through the rules-engine API (cards' Card.Play hooks see the
// engine); pass nil to get a fresh empty state under the engine — cacheable=true,
// currentAuraIdx=-1, fresh logger, no other state.
//
// Sim's chain runner holds *GameState directly and wraps via state.Engine() when it
// needs to invoke card hooks; tests and external callers that want a bare engine reach
// here with nil instead.
func NewFromState(s *GameState) *GameEngine {
	if s == nil {
		s = &GameState{
			cacheable:      true,
			currentAuraIdx: -1,
			logger:         turnlogger.New(),
		}
	}
	return &GameEngine{GameState: s}
}

// NewFromSpec builds a *GameEngine from a Spec, sealing the unexported fields inside
// the package. Use when an external caller — turntest, hero test, EvalOneTurnForTesting
// — needs to construct a prior-turn state. Auras / Triggers / Items aren't on Spec;
// callers chain CreateAura / CreateTrigger / CreateItem on the state after NewFromSpec.
func NewFromSpec(spec Spec) *GameEngine {
	s := &GameState{
		hero:                  spec.Hero,
		arsenal:               spec.Arsenal,
		banished:              spec.Banished,
		graveyard:             spec.Graveyard,
		pitched:               spec.Pitched,
		defenders:             spec.Defenders,
		cardsPlayed:           spec.CardsPlayed,
		cardsRemaining:        spec.CardsRemaining,
		opponentMarked:        spec.OpponentMarked,
		arcaneDamageDealt:     spec.ArcaneDamageDealt,
		auraCreated:           spec.AuraCreated,
		cardBanished:          spec.CardBanished,
		nonAttackActionPlayed: spec.NonAttackActionPlayed,
		actionPoints:          spec.ActionPoints,
		incomingDamage:        spec.IncomingDamage,
		arcaneIncomingDamage:  spec.ArcaneIncomingDamage,
		blockTotal:            spec.BlockTotal,
		value:                 spec.Value,
		triggeringCard:        spec.TriggeringCard,
		attackReactionTarget:  spec.AttackReactionTarget,
		cacheable:             true,
		currentAuraIdx:        -1,
		logger:                turnlogger.New(),
	}
	return &GameEngine{GameState: s}
}

// NewFromCards is a test-only constructor that wraps a Card slice in a fresh *deck.Deck
// and seeds the graveyard. Lets card tests build a GameEngine with a hand-rolled deck
// order without each test plumbing the deck construction inline.
func NewFromCards(deckCards, graveyard []card.Card) *GameEngine {
	dc := make([]deck.Card, len(deckCards))
	for i, c := range deckCards {
		dc[i] = c
	}
	g := NewFromState(nil)
	g.deck = deck.New(nil, nil, dc)
	g.graveyard = graveyard
	return g
}
